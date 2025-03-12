package sheetz

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

// newDummyDriver creates a dummy Driver for testing purposes.
func newDummyDriver(newSheetsService func(ctx context.Context) (*sheets.Service, error)) *Driver {
	return &Driver{
		NewContext:       defaultNewContext,
		NewSheetsService: newSheetsService,
	}
}

// fakeClient is a fake implementation of the Client interface that returns a fixed ValueRange.
type fakeClient struct{}

func (f *fakeClient) SpreadsheetsValuesGet(spreadsheetId string, range_ string, opts ...googleapi.CallOption) (*sheets.ValueRange, error) {
	// Given:
	// A fake client that returns a ValueRange with one comment row, one header row and one data row.
	values := [][]interface{}{
		{"# This is a comment"},  // Comment row (skipped)
		{"-- This is a comment"}, // Comment row (skipped)
		{"col1", "col2"},         // Header row
		{"data1", "data2"},       // Data row
	}
	// When: The client is called to get the values.
	// Then: Return the predetermined ValueRange.
	return &sheets.ValueRange{
		Values: values,
	}, nil
}

// ---------------------------------------------------------------------
// Test for Driver.Open
// ---------------------------------------------------------------------
func TestDriver_Open(t *testing.T) {
	tests := []struct {
		name             string
		newSheetsService func(ctx context.Context) (*sheets.Service, error)
		expectedErr      bool
	}{
		{
			name: "success",
			newSheetsService: func(ctx context.Context) (*sheets.Service, error) {
				return &sheets.Service{
					Spreadsheets: &sheets.SpreadsheetsService{},
				}, nil
			},
			expectedErr: false,
		},
		{
			name: "failure",
			newSheetsService: func(ctx context.Context) (*sheets.Service, error) {
				return nil, errors.New("dummy error")
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given: A driver created with the specified newSheetsService function.
			d := newDummyDriver(tc.newSheetsService)
			// When: Opening a connection using a dummy spreadsheet ID.
			conn_, err := d.Open("dummySpreadsheetID")
			// Then: If an error is expected, ensure an error is returned; otherwise check the connection.
			if tc.expectedErr {
				if err == nil {
					t.Errorf("❌: expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("❌: expected nil error, got %v", err)
			}
			c, ok := conn_.(*conn)
			if !ok {
				t.Fatalf("❌: expected *conn, got %T", conn_)
			}
			if c.sheetID != "dummySpreadsheetID" {
				t.Errorf("❌: expected sheetID=%q, got %q", "dummySpreadsheetID", c.sheetID)
			}
			// Then: Clean up by closing the connection.
			conn_.Close()
		})
	}
}

// ---------------------------------------------------------------------
// Test for Conn.Prepare
// ---------------------------------------------------------------------
func TestConn_Prepare(t *testing.T) {
	// Given: A driver with a proper Sheets service.
	d := newDummyDriver(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: Driver.Open failed: %v", err)
	}
	defer connInterface.Close()
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected *conn, got %T", connInterface)
	}

	tests := []struct {
		name  string
		query string
	}{
		{"prepare_valid_query", "SELECT * FROM Sheet1"},
		{"prepare_valid_with_columns", "SELECT col1, col2 FROM Sheet1"},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given: A valid query string.
			// When: Prepare is called with the query.
			stmtInterface, err := c.Prepare(tc.query)
			// Then: The statement should be created without errors.
			if err != nil {
				t.Fatalf("❌: Prepare failed: %v", err)
			}
			defer stmtInterface.Close()
			s, ok := stmtInterface.(*stmt)
			if !ok {
				t.Fatalf("❌: expected *stmt, got %T", stmtInterface)
			}
			if s.query != tc.query {
				t.Errorf("❌: expected query=%q, got %q", tc.query, s.query)
			}
		})
	}
}

// ---------------------------------------------------------------------
// Test for Conn.Close
// ---------------------------------------------------------------------
func TestConn_Close(t *testing.T) {
	// Given: A driver with a valid Sheets service.
	d := newDummyDriver(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: Driver.Open failed: %v", err)
	}
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected *conn, got %T", connInterface)
	}

	t.Run("close_cancels_context", func(t *testing.T) {
		t.Parallel()
		// Given: A valid connection.
		// When: Closing the connection.
		if err := c.Close(); err != nil {
			t.Errorf("❌: Close() returned error: %v", err)
		}
		// Then: The context should be cancelled (verified with a timeout).
		select {
		case <-c.ctx.Done():
			// OK
		case <-time.After(100 * time.Millisecond):
			t.Error("❌: expected context to be cancelled after Close()")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Conn.Begin
// ---------------------------------------------------------------------
func TestConn_Begin(t *testing.T) {
	// Given: A driver with a valid Sheets service.
	d := newDummyDriver(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: Driver.Open failed: %v", err)
	}
	defer connInterface.Close()
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected *conn, got %T", connInterface)
	}

	t.Run("begin_should_return_ErrSkip", func(t *testing.T) {
		t.Parallel()
		// Given: A valid connection.
		// When: Calling Begin() on the connection.
		_, err := c.Begin()
		// Then: It should return driver.ErrSkip.
		if !errors.Is(err, driver.ErrSkip) {
			t.Errorf("❌: expected driver.ErrSkip, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Close
// ---------------------------------------------------------------------
func TestStmt_Close(t *testing.T) {
	// Given: A driver with a valid Sheets service and a prepared statement.
	d := newDummyDriver(func(ctx context.Context) (*sheets.Service, error) {
		return &sheets.Service{
			Spreadsheets: &sheets.SpreadsheetsService{},
		}, nil
	})
	connInterface, err := d.Open("dummySpreadsheetID")
	if err != nil {
		t.Fatalf("❌: Driver.Open failed: %v", err)
	}
	c, ok := connInterface.(*conn)
	if !ok {
		t.Fatalf("❌: expected *conn, got %T", connInterface)
	}
	stmtInterface, err := c.Prepare("SELECT * FROM Sheet1")
	if err != nil {
		t.Fatalf("❌: Prepare failed: %v", err)
	}
	s, ok := stmtInterface.(*stmt)
	if !ok {
		t.Fatalf("❌: expected *stmt, got %T", stmtInterface)
	}

	t.Run("close_statement_cancels_context", func(t *testing.T) {
		t.Parallel()
		// Given: A prepared statement with an active context.
		// When: Closing the statement.
		if s.ctx.Err() != nil {
			t.Errorf("❌: expected context not canceled before Close, got %v", s.ctx.Err())
		}
		s.Close()
		// Then: The statement's context should be cancelled.
		if s.ctx.Err() == nil {
			t.Error("❌: expected context to be canceled after Close")
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.NumInput
// ---------------------------------------------------------------------
func TestStmt_NumInput(t *testing.T) {
	// Given: A statement with two input placeholders.
	fakeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := &stmt{
		ctx:       fakeCtx,
		ctxCancel: cancel,
		query:     "SELECT ? , ? FROM Sheet1",
	}

	t.Run("num_input_count", func(t *testing.T) {
		t.Parallel()
		// When: Counting the number of input placeholders.
		count := s.NumInput()
		// Then: The count should be -1.
		if count != -1 {
			t.Errorf("❌: expected %d, got %d", -1, count)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Exec
// ---------------------------------------------------------------------
func TestStmt_Exec(t *testing.T) {
	// Given: A statement for an update operation that is not supported.
	fakeCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := &stmt{
		ctx:       fakeCtx,
		ctxCancel: cancel,
		query:     "UPDATE Sheet1 SET col1 = ?",
	}

	t.Run("exec_not_supported", func(t *testing.T) {
		t.Parallel()
		// When: Exec() is called on the statement.
		_, err := s.Exec(nil)
		// Then: It should return driver.ErrSkip.
		if !errors.Is(err, driver.ErrSkip) {
			t.Errorf("❌: expected driver.ErrSkip, got %v", err)
		}
	})
}

// ---------------------------------------------------------------------
// Test for Stmt.Query
// ---------------------------------------------------------------------
func TestStmt_Query(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectErr       bool
		expectedColumns []string
		expectedData    [][]interface{}
	}{
		{
			name:            "success_case",
			query:           "SELECT * FROM Sheet1",
			expectErr:       false,
			expectedColumns: []string{"col1", "col2"},
			expectedData:    [][]interface{}{{"data1", "data2"}},
		},
		{
			name:      "failure_case_with_WHERE",
			query:     "SELECT * FROM Sheet1 WHERE col1 = 1",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc // capture variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given: A driver with a dummy Sheets service using a fake client.
			d := newDummyDriver(func(ctx context.Context) (*sheets.Service, error) {
				return &sheets.Service{
					Spreadsheets: &sheets.SpreadsheetsService{},
				}, nil
			})
			connInterface, err := d.Open("dummySpreadsheetID")
			if err != nil {
				t.Fatalf("❌: Driver.Open failed: %v", err)
			}
			c, ok := connInterface.(*conn)
			if !ok {
				t.Fatalf("❌: expected *conn, got %T", connInterface)
			}
			// Given: Setting a fake client to simulate a predefined ValueRange response.
			c.client = &fakeClient{}

			// When: Preparing and executing the query.
			stmtInterface, err := c.Prepare(tc.query)
			if err != nil {
				t.Fatalf("❌: Prepare failed: %v", err)
			}
			defer stmtInterface.Close()
			rows, err := stmtInterface.Query(nil)
			// Then: Validate the error expectation and response.
			if tc.expectErr {
				if err == nil {
					t.Error("❌: expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("❌: expected nil error, got %v", err)
			}
			r, ok := rows.(*gSheetRows)
			if !ok {
				t.Fatalf("❌: expected *gSheetRows, got %T", rows)
			}
			// Then: Check that the columns match expectations.
			if len(r.columns) != len(tc.expectedColumns) {
				t.Errorf("❌: expected columns %v, got %v", tc.expectedColumns, r.columns)
			}
			for i, col := range tc.expectedColumns {
				if r.columns[i] != col {
					t.Errorf("❌: expected column %q, got %q", col, r.columns[i])
				}
			}
			// Then: Check that the data rows match expectations.
			if len(r.data) != len(tc.expectedData) {
				t.Errorf("❌: expected data rows %v, got %v", tc.expectedData, r.data)
			}
		})
	}
}
