package sheetz

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"google.golang.org/api/sheets/v4"
)

// ---------------------------------------------------------------------
// Test for SetNewContext and GetNewContext
// ---------------------------------------------------------------------

func TestSetNewContext(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		// Set a custom context generator and verify it can be retrieved.
		customCtx := context.WithValue(context.Background(), "key", "value")
		SetNewContext(func() context.Context {
			return customCtx
		})
		got := GetNewContext()()
		if got != customCtx {
			t.Errorf("❌: expected=%v, actual=%v", customCtx, got)
		}
	})

	t.Run("failure,", func(t *testing.T) {
		// Skip because there are no failure cases for this function.
		t.Skipf("🚫: SKIP: SetNewContext always succeeds so no failure case exists")
	})
}

func TestGetNewContext(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		// Verify that the default newContext function returns a non-nil context.
		ctx := GetNewContext()()
		if ctx == nil {
			t.Error("❌: expected non-nil context, actual=nil")
		}
	})
}

// ---------------------------------------------------------------------
// Test for SetNewSheetsService and GetNewSheetsService
// ---------------------------------------------------------------------

func TestSetNewSheetsService(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		dummyFunc := func(ctx context.Context) (*sheets.Service, error) {
			return &sheets.Service{}, nil
		}
		SetNewSheetsService(dummyFunc)
		gotFunc := GetNewSheetsService()
		svc, err := gotFunc(context.Background())
		if err != nil {
			t.Errorf("❌: expected=nil, actual=%v", err)
		}
		if svc == nil {
			t.Error("❌: expected non-nil Sheets service, actual=nil")
		}
	})

	t.Run("failure,", func(t *testing.T) {
		// Skip because there is no failure case.
		t.Skipf("🚫: SKIP: SetNewSheetsService always succeeds so no failure case exists")
	})
}

func TestGetNewSheetsService(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		gotFunc := GetNewSheetsService()
		svc, err := gotFunc(context.Background())
		if err != nil {
			t.Errorf("❌: expected=nil, actual=%v", err)
		}
		if svc == nil {
			t.Error("❌: expected non-nil Sheets service, actual=nil")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Driver.Open
// ---------------------------------------------------------------------

func TestDriver_Open(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		// For testing: set a function that returns a dummy Sheets service.
		SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
			// This test does not call Query, so return only minimal values.
			return &sheets.Service{
				Spreadsheets: &sheets.SpreadsheetsService{
					// Fields can be set as necessary.
				},
			}, nil
		})
		d := &Driver{}
		conn_, err := d.Open("dummySpreadsheetID")
		if err != nil {
			t.Fatalf("❌: expected=nil, actual=%v", err)
		}
		c, ok := conn_.(*conn)
		if !ok {
			t.Fatalf("❌: expected=*conn, actual=%T", conn_)
		}
		if c.sheetID != "dummySpreadsheetID" {
			t.Errorf("❌: expected=%v, actual=%v", "dummySpreadsheetID", c.sheetID)
		}
		// Cleanup.
		conn_.Close()
	})

	t.Run("failure,", func(t *testing.T) {
		// Set newSheetsService to return an error.
		SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
			return nil, errors.New("dummy error")
		})
		d := &Driver{}
		_, err := d.Open("dummySpreadsheetID")
		if err == nil {
			t.Error("❌: expected=error, actual=nil")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Conn.Prepare
// ---------------------------------------------------------------------

func TestConn_Prepare(t *testing.T) {
	t.Parallel()

	SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})

	d := &Driver{}
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: expected=nil, actual=%v", err)
	}
	defer connInterface.Close()
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected=*conn, actual=%T", connInterface)
	}

	t.Run("success,", func(t *testing.T) {
		stmtInterface, err := c.Prepare("SELECT * FROM Sheet1")
		if err != nil {
			t.Fatalf("❌: expected=nil, actual=%v", err)
		}
		defer stmtInterface.Close()
		s, ok := stmtInterface.(*stmt)
		if !ok {
			t.Fatalf("❌: expected=*stmt, actual=%T", stmtInterface)
		}
		if s.query != "SELECT * FROM Sheet1" {
			t.Errorf("❌: expected=%v, actual=%v", "SELECT * FROM Sheet1", s.query)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Conn.Close
// ---------------------------------------------------------------------

func TestConn_Close(t *testing.T) {
	t.Parallel()

	SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	d := &Driver{}
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: expected=nil, actual=%v", err)
	}
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected=*conn, actual=%T", connInterface)
	}

	t.Run("success,", func(t *testing.T) {
		err := c.Close()
		if err != nil {
			t.Errorf("❌: expected=nil, actual=%v", err)
		}
		select {
		case <-c.ctx.Done():
			// Verify that the context is cancelled.
		default:
			t.Error("❌: expected=context canceled, actual=context not canceled")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Conn.Begin
// ---------------------------------------------------------------------

func TestConn_Begin(t *testing.T) {
	t.Parallel()

	SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	d := &Driver{}
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: expected=nil, actual=%v", err)
	}
	defer connInterface.Close()
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected=*conn, actual=%T", connInterface)
	}

	t.Run("failure,", func(t *testing.T) {
		_, err := c.Begin()
		if err != driver.ErrSkip {
			t.Errorf("❌: expected=driver.ErrSkip, actual=%v", err)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Close
// ---------------------------------------------------------------------

func TestStmt_Close(t *testing.T) {
	t.Parallel()

	SetNewSheetsService(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	d := &Driver{}
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: expected=nil, actual=%v", err)
	}
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected=*conn, actual=%T", connInterface)
	}
	stmtInterface, err := c.Prepare("SELECT * FROM Sheet1")
	if err != nil {
		t.Fatalf("❌: expected=nil, actual=%v", err)
	}
	s, ok := stmtInterface.(*stmt)
	if !ok {
		t.Fatalf("❌: expected=*stmt, actual=%T", stmtInterface)
	}

	t.Run("success,", func(t *testing.T) {
		// Before calling Close(), the context should not be cancelled.
		if s.ctx.Err() != nil {
			t.Error("❌: expected=nil, actual=%v", s.ctx.Err())
		}
		s.Close()
		if s.ctx.Err() == nil {
			t.Error("❌: expected=context canceled, actual=context not canceled")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.NumInput
// ---------------------------------------------------------------------

func TestStmt_NumInput(t *testing.T) {
	t.Parallel()

	fakeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := &stmt{
		ctx:       fakeCtx,
		ctxCancel: cancel,
		query:     "SELECT ? , ? FROM Sheet1",
	}

	t.Run("success,", func(t *testing.T) {
		count := s.NumInput()
		if count != 2 {
			t.Errorf("❌: expected=%d, actual=%d", 2, count)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Exec
// ---------------------------------------------------------------------

func TestStmt_Exec(t *testing.T) {
	t.Parallel()

	fakeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := &stmt{
		ctx:       fakeCtx,
		ctxCancel: cancel,
		query:     "UPDATE Sheet1 SET col1 = ?",
	}

	t.Run("failure,", func(t *testing.T) {
		_, err := s.Exec(nil)
		if err != driver.ErrSkip {
			t.Errorf("❌: expected=driver.ErrSkip, actual=%v", err)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Query
// ---------------------------------------------------------------------

func TestStmt_Query(t *testing.T) {
	t.Parallel()

	t.Run("success,", func(t *testing.T) {
		// Success case:
		// Test if a valid SELECT query ("SELECT * FROM Sheet1") can be processed.
		// Note: In a real test, a fake Sheets service should be implemented to return a ValueRange with expected data from the Sheets API,
		// but due to implementation complexity, this test case is skipped.
		t.Skipf("🚫: SKIP: Skipping success case because a fake Sheets service implementation is required")
	})

	t.Run("failure,", func(t *testing.T) {
		// A SELECT query containing a WHERE clause is not supported and should return an error.
		s := &stmt{
			ctx:   context.Background(),
			query: "SELECT * FROM Sheet1 WHERE col1 = 1",
		}
		_, err := s.Query(nil)
		if !errors.Is(err, ErrWhereClauseIsNotSupportedCurrently) {
			t.Errorf("❌: expected=%v, actual=%v", ErrWhereClauseIsNotSupportedCurrently, err)
		}
	})
}
