package sheetz

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"

	"google.golang.org/api/googleapi"
	sheets "google.golang.org/api/sheets/v4"
)

func init() {
	// Register the driver with database/sql
	sql.Register("gsheets", &Driver{})
}

var DefaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelWarn}))

var (
	ErrUnsupportedQueryFormat             = errors.New("sheetz: unsupported query format")
	ErrWhereClauseIsNotSupportedCurrently = errors.New("sheetz: WHERE clause is not supported currently")
)

// =====================================
// Driver
// =====================================
type Driver struct {
	NewContext       func() context.Context
	NewSheetsService func(ctx context.Context) (*sheets.Service, error)
	ConvertValueFunc func(ctx context.Context, spreadsheetId string, sheetName string, columnName string, value any) (any, error)
	Logger           *slog.Logger
}

var _ driver.Driver = (*Driver)(nil)

var (
	newContext         func() context.Context = context.Background
	newContextMu       sync.RWMutex
	newSheetsService   func(ctx context.Context) (*sheets.Service, error) = func(ctx context.Context) (*sheets.Service, error) { return sheets.NewService(ctx) }
	newSheetsServiceMu sync.RWMutex
)

func defaultNewContext() context.Context {
	return context.Background()
}

func defaultNewSheetsService(ctx context.Context) (*sheets.Service, error) {
	return sheets.NewService(ctx)
}

type Client interface {
	SpreadsheetsValuesGet(spreadsheetId string, range_ string, opts ...googleapi.CallOption) (*sheets.ValueRange, error)
}

type client struct {
	driver        *Driver
	sheetsService *sheets.Service
}

func (c *client) SpreadsheetsValuesGet(spreadsheetId string, range_ string, opts ...googleapi.CallOption) (*sheets.ValueRange, error) {
	c.driver.logger().Debug("SpreadsheetsValuesGet", slog.String("spreadsheetId", spreadsheetId), slog.String("range", range_))
	return c.sheetsService.Spreadsheets.Values.Get(spreadsheetId, range_).Do(opts...)
}

func (d *Driver) convertValue(ctx context.Context, spreadsheetId string, sheetName string, columnName string, value any) (any, error) {
	d.logger().Debug("convertValue", slog.String("spreadsheetId", spreadsheetId), slog.String("sheetName", sheetName), slog.String("columnName", columnName), slog.Any("value", value))
	if d.ConvertValueFunc != nil {
		converted, err := d.ConvertValueFunc(ctx, spreadsheetId, sheetName, columnName, value)
		if err != nil {
			return nil, fmt.Errorf("d.ConvertValueFunc: spreadsheetId=%q, sheetName=%q, columnName=%q, value=%v: %w", spreadsheetId, sheetName, columnName, value, err)
		}
		return converted, nil
	}
	return value, nil
}

func (d *Driver) logger() *slog.Logger {
	if d.Logger != nil {
		return d.Logger
	}
	return DefaultLogger
}

// dsn is expected to be a string in the format "<SpreadsheetID>" (in reality, you might need to include the sheet name as well)
// dsn is expected to be a string in the format "<SpreadsheetID>" (in reality, you might need to include the sheet name as well)
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	// Set up authentication and create a Sheets API client
	// In actual usage, you would typically read credentials from credentials.json or from a token file, etc.
	newContext := defaultNewContext
	if d.NewContext != nil {
		newContext = d.NewContext
	}
	ctx, cancel := context.WithCancel(newContext())

	newSheetsService := defaultNewSheetsService
	if d.NewSheetsService != nil {
		newSheetsService = d.NewSheetsService
	}
	srv, err := newSheetsService(ctx)
	if err != nil {
		defer cancel()
		return nil, fmt.Errorf("newSheetsService: %w", err)
	}
	client := &client{
		driver:        d,
		sheetsService: srv,
	}

	conn := &conn{
		driver:    d,
		ctx:       ctx,
		ctxCancel: cancel,
		sheetID:   dsn,
		client:    client,
	}
	return conn, nil
}

// =====================================
// Conn
// =====================================
type conn struct {
	driver    *Driver
	ctx       context.Context
	ctxCancel context.CancelFunc
	sheetID   string
	client    Client
}

var (
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.Conn               = (*conn)(nil)
)

// PrepareContext creates a statement from an SQL query.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	// In this simplified example, pass the query directly to the statement.
	return &stmt{
		ctx:       ctx,
		ctxCancel: c.ctxCancel,
		conn:      c,
		query:     query,
	}, nil
}

// Prepare creates a statement from an SQL query.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(c.ctx, query)
}

// Close closes the connection.
func (c *conn) Close() error {
	c.ctxCancel()
	return nil
}

// Begin starts a transaction, but since this sample is intended for read-only operations, it returns an unimplemented error.
func (c *conn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

// =====================================
// Stmt
// =====================================
type stmt struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *conn
	query     string
}

var (
	_ driver.StmtQueryContext = (*stmt)(nil)
	_ driver.Stmt             = (*stmt)(nil)
)

// Close closes the statement.
func (s *stmt) Close() error {
	s.ctxCancel()
	return nil
}

// NumInput returns the number of placeholders. It is not supported currently.
//
// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.24.1:src/database/sql/driver/driver.go;l=346-349
// > NumInput may also return -1, if the driver doesn't know
// > its number of placeholders. In that case, the sql package
// > will not sanity check Exec or Query argument counts.
func (s *stmt) NumInput() int {
	return -1
}

// Exec executes an update query. It is not supported currently.
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, driver.ErrSkip
}

// Query executes a read query.
func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 1) Parse the query string to extract the following:
	//    - columns   : list of columns specified in the SELECT clause (if ["*"], then all)
	//    - sheetName : table name specified in the FROM clause (i.e. the sheet name)
	//    - rangeSpec : string specified in the RANGE clause (e.g., "A1:Z") (optional)
	parsedQuery, err := parseSelectQuery(s.query)
	if err != nil {
		return nil, err
	}

	// 2) Retrieve values from the sheet
	resp, err := s.conn.client.SpreadsheetsValuesGet(s.conn.sheetID, parsedQuery.sheetRange)
	if err != nil {
		return nil, fmt.Errorf("spreadsheetsValuesGetCall.Do: %w", err)
	}

	// 3) Skip comment rows at the beginning (if there are multiple consecutive comment rows, skip them all).
	//    Rows whose first cell starts with "#" or "--" are considered comments.
	resp.Values = skipCommentRows(resp.Values)

	// If no data remains, return an empty set of rows.
	if len(resp.Values) == 0 {
		return &gSheetRows{
			columns: []string{},
			data:    [][]interface{}{},
			curr:    0,
		}, nil
	}

	// 4) Use the first row as the header to determine column names (ignore empty cells)
	headerRow := resp.Values[0]
	var columnsCommentFiltered []string
	var columnsCommentFilteredIndexMap []int
	for colIdx, v := range headerRow {
		headerVal, ok := v.(string)
		if !ok || headerVal == "" {
			// Ignore empty cells or non-string values
			continue
		}
		// Additionally, if the value starts with "#" or "--" in the header row, ignore this column.
		if strings.HasPrefix(headerVal, "#") || strings.HasPrefix(headerVal, "--") {
			continue
		}

		columnsCommentFiltered = append(columnsCommentFiltered, headerVal)
		columnsCommentFilteredIndexMap = append(columnsCommentFilteredIndexMap, colIdx)
	}

	// 5) Determine the columns specified in the SELECT clause.
	parsedColumns := parsedQuery.columns
	var columnsQueried []string
	var columnsQueriedIndexMap []int
	if len(parsedColumns) == 1 && parsedColumns[0] == "*" {
		// All columns
		columnsQueried = columnsCommentFiltered
		columnsQueriedIndexMap = columnsCommentFilteredIndexMap
	} else {
		// Only the specified columns
		colMap := make(map[string]int) // Map column names to their index in allColumns.
		for i, colName := range columnsCommentFiltered {
			colMap[colName] = i
		}

		for _, requestedCol := range parsedColumns {
			idx, ok := colMap[requestedCol]
			if !ok {
				return nil, fmt.Errorf("unknown column name: sheet=%q, column=%q", parsedQuery.sheetName, requestedCol)
			}
			columnsQueried = append(columnsQueried, requestedCol)
			columnsQueriedIndexMap = append(columnsQueriedIndexMap, columnsCommentFilteredIndexMap[idx])
		}
	}

	// 6) Retrieve data starting from the second row.
	var data [][]interface{}
	for rowIdx := 1; rowIdx < len(resp.Values); rowIdx++ {
		row := resp.Values[rowIdx]
		rowData := make([]interface{}, len(columnsQueriedIndexMap))
		for i, colIdx := range columnsQueriedIndexMap {
			if colIdx < len(row) {
				rowData[i] = row[colIdx]
			} else {
				rowData[i] = nil
			}
		}
		data = append(data, rowData)
	}

	// 7) Convert values to the appropriate type.
	for rowIdx, row := range data {
		for colIdx, colName := range columnsQueried {
			converted, err := s.conn.driver.convertValue(s.conn.ctx, s.conn.sheetID, parsedQuery.sheetName, colName, row[colIdx])
			if err != nil {
				return nil, fmt.Errorf("s.conn.driver.convertValue: %w", err)
			}
			data[rowIdx][colIdx] = converted
		}
	}

	return &gSheetRows{
		columns: columnsQueried,
		data:    data,
		curr:    0,
	}, nil
}

// skipCommentRows treats rows whose first cell starts with "#" or "--" as comments,
// skipping them all and returning the remaining rows.
func skipCommentRows(rows [][]interface{}) [][]interface{} {
	i := 0
	for i < len(rows) {
		if len(rows[i]) == 0 {
			// An empty row cannot be considered a comment row.
			// In this design, empty rows are treated as non-comment rows, so break.
			break
		}
		firstCell, ok := rows[i][0].(string)
		if !ok {
			// If the cell is not a string, it is not considered a comment row, so break.
			break
		}
		if strings.HasPrefix(firstCell, "#") || strings.HasPrefix(firstCell, "--") {
			// This is a comment row, so skip it.
			i++
			continue
		}
		// If not a comment row, break.
		break
	}
	return rows[i:]
}

func toNamedValues(args []driver.Value) []driver.NamedValue {
	namedArgs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		namedArgs[i] = driver.NamedValue{
			Ordinal: i,
			Value:   arg,
		}
	}
	return namedArgs
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(s.ctx, toNamedValues(args))
}

// =====================================
// Rows
// =====================================
type gSheetRows struct {
	columns []string
	data    [][]interface{}
	curr    int
}

var _ driver.Rows = (*gSheetRows)(nil)

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
// SQL Parser (Experimental Implementation)
//
// Supported examples:
//
//	SELECT * FROM Sheet1
//	SELECT column1, column2 FROM Sheet1
//	SELECT * FROM Sheet1!A1:Z
//	SELECT column1, column2 FROM "Sheet With Space"!A1:Z
//	(Trailing semicolon is optional.)
//
// Desired output:
//
//	columns    []string ("*" or ["colA","colB",...])
//	sheetName  string
//	sheetRange string
//
// =====================================
type parsedSelect struct {
	columns    []string
	sheetName  string
	rangeSpec  string
	sheetRange string
}

// Use a regular expression to extract the following:
//
//	columns   : list of columns specified in the SELECT clause (if ["*"], then all)
//	sheetName : table name specified in the FROM clause (i.e. the sheet name)
//	rangeSpec : string specified in the RANGE clause (e.g., "A1:Z") (optional)
//
// Example: SELECT column1, column2 FROM Sheet1!A1:Z
// ->              ^^^^^^^^^^^^^^^^      ^^^^^^^^^^^
var (
	identRegexStr = `("([^"]+)"|'([^']+)'|` + "`([^`]+)`" + `|(\S+))`
	selectRegex   = regexp.MustCompile(
		`(?i)^\s*SELECT\s+` + // SELECT
			`(?P<columns>.+?)\s+` + // column part
			`FROM\s+` + // FROM
			`(?P<sheetRange>` + // sheet range
			`(?P<sheetName>` + // sheet name
			`"(?P<sheetNameInQuote>[^"]+?)"|` +
			`'(?P<sheetNameInQuote>[^']+?)'|` +
			"`(?P<sheetNameInQuote>[^`]+?)`|" +
			`\S+?` +
			`)` + // sheet name
			`(!(?P<rangeSpec>\S+?))?` + // range spec
			`)` + // sheet range
			`(\s+(?P<wherePart>WHERE\s+` + // WHERE clause
			`(.+?)` + // TODO: implement
			`))?` + // WHERE clause
			`\s*;?\s*$`,
	)
)

func parseSelectQuery(query string) (*parsedSelect, error) {
	trimmed := strings.TrimSpace(query)
	matches := selectRegex.FindStringSubmatch(trimmed)
	if len(matches) == 0 {
		return nil, fmt.Errorf("query=%q: %w", query, ErrUnsupportedQueryFormat)
	}

	result := make(map[string]string)
	for i, name := range selectRegex.SubexpNames() {
		if name == "" {
			continue
		}
		r := result[name]
		if r == "" {
			result[name] = matches[i]
		}
	}

	// matches["columns"]          -> column part (e.g., "*" or "column1, column2")
	// matches["sheetRange"]       -> sheet range (e.g., "A1:Z")
	// matches["sheetName"]        -> full sheet name (e.g., "Sheet1" or "Sheet With Space")
	// matches["sheetNameInQuote"] -> double-quoted content of the sheet name (e.g., Sheet With Space)
	// matches["rangeSpec"]        -> range spec (e.g., "A1:Z")
	// matches["wherePart"]        -> WHERE clause
	columnsPart := result["columns"]
	sheetNamePart := result["sheetName"]
	sheetNameInQuotePart := result["sheetNameInQuote"]
	rangeSpecPart := result["rangeSpec"]
	wherePart := result["wherePart"]

	if wherePart != "" {
		return nil, fmt.Errorf("query=%q: %w", query, ErrWhereClauseIsNotSupportedCurrently)
	}

	// Parse the column list.
	colList := parseColumnList(columnsPart)

	// Determine the sheet name.
	sheetName := sheetNamePart
	if sheetNameInQuotePart != "" {
		sheetName = sheetNameInQuotePart
	}

	sheetRange := sheetName
	if rangeSpecPart != "" {
		sheetRange = sheetName + "!" + rangeSpecPart
	}

	return &parsedSelect{
		columns:    colList,
		sheetName:  sheetName,
		rangeSpec:  rangeSpecPart,
		sheetRange: sheetRange,
	}, nil
}

// parseColumnList splits a string like "column1, column2" into a slice.
// If the string is "*" only, it returns ["*"].
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
		// If empty, default to "*"
		return []string{"*"}
	}
	return result
}
