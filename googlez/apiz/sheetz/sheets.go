package sheetz

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	// Google Sheets API (go.mod に以下のような行が必要です)
	// google.golang.org/api/sheets/v4

	sheets "google.golang.org/api/sheets/v4"
)

func init() {
	// database/sql にドライバを登録
	sql.Register("gsheets", &GSheetDriver{})
}

var ErrUnsupportedQueryFormat = errors.New("unsupported query format")

// =====================================
// Driver
// =====================================
type GSheetDriver struct{}

var _ driver.Driver = (*GSheetDriver)(nil)

// dsn は「<SpreadsheetID>」のような文字列を想定（本当はシート名も含めるなど工夫が必要）
func (d *GSheetDriver) Open(dsn string) (driver.Conn, error) {
	// ここで認証情報などを設定し、Sheets API クライアントを作成
	// 実際の使用では credentials.json やトークンファイルを読み込むなどが必要になります。

	ctx, cancel := context.WithCancel(context.Background())
	srv, err := sheets.NewService(ctx)
	if err != nil {
		defer cancel()
		return nil, fmt.Errorf("sheets.NewService: %w", err)
	}

	conn := &gSheetConn{
		ctx:       ctx,
		ctxCancel: cancel,
		sheetID:   dsn,
		service:   srv,
	}
	return conn, nil
}

// =====================================
// Conn
// =====================================
type gSheetConn struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	sheetID   string
	service   *sheets.Service
}

var _ driver.Conn = (*gSheetConn)(nil)

// Prepare は SQL 文（query）からステートメントを生成する。
func (c *gSheetConn) Prepare(query string) (driver.Stmt, error) {
	// この例では実装を簡単化するため、そのままステートメントに渡すだけ。
	return &gSheetStmt{
		ctx:       c.ctx,
		ctxCancel: c.ctxCancel,
		conn:      c,
		query:     query,
	}, nil
}

// Close はコネクションをクローズする。
func (c *gSheetConn) Close() error {
	c.ctxCancel()
	return nil
}

// Begin はトランザクション開始メソッドだが、
// 今回のサンプルでは読み取り専用の想定として、未実装エラーを返す。
func (c *gSheetConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

// =====================================
// Stmt
// =====================================
type gSheetStmt struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *gSheetConn
	query     string
}

var _ driver.Stmt = (*gSheetStmt)(nil)

// Close はステートメントをクローズする。
func (s *gSheetStmt) Close() error {
	s.ctxCancel()
	return nil
}

// NumInput はプレースホルダ数を返す。
func (s *gSheetStmt) NumInput() int {
	return strings.Count(s.query, "?")
}

// Exec は更新系のクエリ実行。今回のサンプルでは未対応。
func (s *gSheetStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, driver.ErrSkip
}

// Query は読み取り系クエリ実行
func (s *gSheetStmt) Query(args []driver.Value) (driver.Rows, error) {
	// 1) クエリ文字列をパースして以下を取得
	//    - columns   : SELECT で指定したカラム一覧 (["*"] なら全部)
	//    - sheetName : FROM で指定したテーブル名 (シート名)
	//    - rangeSpec : RANGE で指定した文字列 (例: "A1:Z") (オプション)
	parsed, err := parseSelectQuery(s.query)
	if err != nil {
		return nil, err
	}

	// 2) シートから値を取得する
	readRange := parsed.sheetName
	if parsed.rangeSpec != "" {
		// もし RANGE(...) が指定されていたら "シート名!範囲" として指定
		// 例: "Sheet1!A1:Z"
		readRange = parsed.sheetName + "!" + parsed.rangeSpec
	}

	spreadsheetsValuesGetCall := s.conn.service.Spreadsheets.Values.Get(s.conn.sheetID, readRange)
	resp, err := spreadsheetsValuesGetCall.Do()
	if err != nil {
		return nil, fmt.Errorf("spreadsheetsValuesGetCall.Do: %w", err)
	}

	// 3) 先頭行のコメント行スキップ (複数連続している場合はその分スキップ)
	//    1列目の値が "#" or "--" で始まる行をコメント行とみなし、次に進む
	resp.Values = skipCommentRows(resp.Values)

	// データが無くなった場合
	if len(resp.Values) == 0 {
		return &gSheetRows{
			columns: []string{},
			data:    [][]interface{}{},
			curr:    0,
		}, nil
	}

	// 4) 1行目をヘッダとしてカラム名を確定 (空セルは無視)
	headerRow := resp.Values[0]
	var allColumns []string
	var allColumnsIndexMap []int
	for colIdx, v := range headerRow {
		headerVal, ok := v.(string)
		if !ok || headerVal == "" {
			// 空セル or 文字列じゃないなら無視
			continue
		}
		// 追加: # または -- で始まるなら、このカラムは無視
		if strings.HasPrefix(headerVal, "#") || strings.HasPrefix(headerVal, "--") {
			continue
		}

		allColumns = append(allColumns, headerVal)
		allColumnsIndexMap = append(allColumnsIndexMap, colIdx)
	}

	// 5) SELECT 句で指定したカラムを決定
	var finalColumns []string
	var finalIndexMap []int
	if len(parsed.columns) == 1 && parsed.columns[0] == "*" {
		// 全カラム
		finalColumns = allColumns
		finalIndexMap = allColumnsIndexMap
	} else {
		// 指定カラムのみ
		colMap := make(map[string]int) // カラム名 -> allColumns内のインデックス
		for i, colName := range allColumns {
			colMap[colName] = i
		}

		for _, requestedCol := range parsed.columns {
			idx, ok := colMap[requestedCol]
			if !ok {
				return nil, fmt.Errorf("unknown column name: %q", requestedCol)
			}
			finalColumns = append(finalColumns, requestedCol)
			finalIndexMap = append(finalIndexMap, allColumnsIndexMap[idx])
		}
	}

	// 6) 2行目以降をデータとして取得
	var data [][]interface{}
	for rowIdx := 1; rowIdx < len(resp.Values); rowIdx++ {
		row := resp.Values[rowIdx]
		rowData := make([]interface{}, len(finalIndexMap))
		for i, colIdx := range finalIndexMap {
			if colIdx < len(row) {
				rowData[i] = row[colIdx]
			} else {
				rowData[i] = nil
			}
		}
		data = append(data, rowData)
	}

	return &gSheetRows{
		columns: finalColumns,
		data:    data,
		curr:    0,
	}, nil
}

// skipCommentRows は、先頭から順番に「1列目が # or -- で始まる」行を
// コメント行とみなし、すべてスキップして残りを返す。
func skipCommentRows(rows [][]interface{}) [][]interface{} {
	i := 0
	for i < len(rows) {
		if len(rows[i]) == 0 {
			// 空行ならコメント行とは言えないが、今回どう扱うかは設計次第。
			// ここでは「空行はコメント行ではない」扱いとして break。
			break
		}
		firstCell, ok := rows[i][0].(string)
		if !ok {
			// 文字列じゃないならコメント行とは言えないので break。
			break
		}
		if strings.HasPrefix(firstCell, "#") || strings.HasPrefix(firstCell, "--") {
			// コメント行なのでスキップ
			i++
			continue
		}
		// コメント行でなければ break
		break
	}
	return rows[i:]
}

// =====================================
// Rows
// =====================================
type gSheetRows struct {
	columns []string
	data    [][]interface{}
	curr    int
}

func (r *gSheetRows) Columns() []string {
	return r.columns
}

func (r *gSheetRows) Close() error {
	return nil
}

func (r *gSheetRows) Next(dest []driver.Value) error {
	if r.curr >= len(r.data) {
		return io.EOF
	}
	rowData := r.data[r.curr]
	r.curr++
	for i, val := range rowData {
		dest[i] = val
	}
	return nil
}

// =====================================
// SQL パーサ (簡易実装)
//
// 対応例:
//
//	SELECT * FROM Sheet1
//	SELECT column1, column2 FROM Sheet1
//	SELECT * FROM Sheet1 WHERE RANGE("A1:Z")
//	SELECT column1, column2 FROM "Sheet With Space" WHERE RANGE("B1:D10")
//	(末尾のセミコロン ; は省略可)
//
// 取得したいもの:
//
//	columns   []string ("*" or ["colA","colB",...])
//	sheetName string
//	rangeSpec string  (例: "A1:Z")
//
// =====================================
type parsedSelect struct {
	columns   []string
	sheetName string
	rangeSpec string
}

// 今回は "WHERE RANGE(...)" の形式を正規表現で取り出す
//
// 例: SELECT column1, column2 FROM Sheet1 WHERE RANGE("A1:Z")
//
//	^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^^^^ ^^^^^^^^^^^^^^^^^^
var selectRegex = regexp.MustCompile(
	`(?i)^\s*SELECT\s+(.+?)\s+FROM\s+("([^"]+)"|[A-Za-z0-9_]+)(?:\s+WHERE\s+RANGE\("([^"]+)"\))?\s*;?\s*$`,
)

func parseSelectQuery(query string) (*parsedSelect, error) {
	trimmed := strings.TrimSpace(query)
	matches := selectRegex.FindStringSubmatch(trimmed)
	if len(matches) == 0 {
		return nil, fmt.Errorf("query=%q: %w", query, ErrUnsupportedQueryFormat)
	}

	// matches[1] -> カラム部 (例: "*" or "column1, column2")
	// matches[2] -> シート名全体 ("Sheet1" もしくは "Sheet With Space" 等)
	// matches[3] -> シート名のダブルクォート内容 (例: Sheet With Space)
	// matches[4] -> RANGE("A1:Z") の A1:Z 部分

	columnPart := matches[1]
	sheetPart := matches[2]
	inQuote := matches[3]
	rangeSpec := matches[4]

	// シート名を決定
	sheetName := sheetPart
	if inQuote != "" {
		sheetName = inQuote
	}

	// カラムリストを解析
	colList := parseColumnList(columnPart)

	return &parsedSelect{
		columns:   colList,
		sheetName: sheetName,
		rangeSpec: rangeSpec,
	}, nil
}

// parseColumnList は "column1, column2" のような文字列を分割して配列に。
// "*" だけなら ["*"] を返す。
func parseColumnList(cols string) []string {
	trimmed := strings.TrimSpace(cols)
	if trimmed == "*" {
		return []string{"*"}
	}
	parts := strings.Split(trimmed, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		// 万が一空なら "*" と同義にする
		return []string{"*"}
	}
	return result
}
