package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gormless/data"
	"gormless/data/dialect"
	"os"
	"testing"
)

// TestMain is used for setup and teardown of tests
func TestMain(m *testing.M) {
	// Setup code here (if needed)

	// Run tests
	code := m.Run()

	// Teardown code here (if needed)

	os.Exit(code)
}

// TestMainFunction is a unit test for the main function
// This is a bit tricky to test since main() doesn't return anything,
// and we need to mock database connections, so this is more of an example
func TestMainFunction(t *testing.T) {
	// Create a sqlmock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a real Session with the mock database
	session := &data.Session{DB: db}
	session.SQLDialect = dialect.PostgresDialect{}

	// Set up expectations
	mock.ExpectPing()
	mock.ExpectPrepare("CREATE TABLE IF NOT EXISTS").ExpectExec().WillReturnResult(sqlmock.NewResult(0, 0))

	// Call function under test
	err = data.InitDatabaseVersion(session)
	assert.NoError(t, err)

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
