package data

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gormless/data/dialect"
	fixtures "gormless/data/fixtures"
	"regexp"
	"testing"
)

// MockDB is a mock for the sql.DB

func TestCreateTable(t *testing.T) {
	idType := dialect.PsqlSerial
	nameType := fmt.Sprintf(dialect.PsqlVarChar, 32)
	tests := []struct {
		name          string
		table         Table
		expectError   bool
		errorContains string
	}{

		{
			name: "Valid table with primary key",
			table: Table{
				Name: "test_table",
				Columns: &[]Column{
					{Name: "id", Type: &idType, PrimaryKey: true},
					{Name: "name", Type: &nameType},
				},
			},
			expectError: false,
		},
		{
			name: "Multiple primary keys should fail",
			table: Table{
				Name: "test_table",
				Columns: &[]Column{
					{Name: "id", Type: &idType, PrimaryKey: true},
					{Name: "uuid", Type: stringPtr(dialect.PsqlUuid), PrimaryKey: true},
				},
			},
			expectError:   true,
			errorContains: "multiple primary keys",
		},
		{
			name: "SQL injection in column name should fail",
			table: Table{
				Name: "test_table",
				Columns: &[]Column{
					{Name: "id", Type: &idType, PrimaryKey: true},
					{Name: "name'; DROP TABLE users; --", Type: stringPtr(fmt.Sprintf(dialect.PsqlVarChar, 32))},
				},
			},
			expectError:   true,
			errorContains: "invalid SQL identifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock objects
			mockDB := new(fixtures.MockDB)
			mockStmt := new(fixtures.MockStmt)

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error creating mock database: %v", err)
			}
			defer db.Close()
			expectedSQL := regexp.QuoteMeta(fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\"", tt.table.Name))
			// Only set up expectations if we don't expect an error during validation
			if !tt.expectError {
				// Check if the SQL query is properly formatted
				mock.ExpectPrepare(expectedSQL).
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 0))
			}
			session := &Session{DB: db}
			session.SQLDialect = dialect.PostgresDialect{}
			// Call the function being tested
			err = CreateTable(session, tt.table)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				mockDB.AssertExpectations(t)
				mockStmt.AssertExpectations(t)
			}
		})
	}
}

// Helper function to get string pointer
func stringPtr(s string) *string {
	return &s
}
