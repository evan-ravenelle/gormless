package data

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gormless/data/dialect"
	"regexp"
	"testing"
	"time"
)

func TestInitDatabaseVersion(t *testing.T) {
	// Create a mock SQL driver and connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	session := &Session{DB: db}
	session.SQLDialect = dialect.PostgresDialect{}
	// Set expectations for the database operations
	mock.ExpectPing()

	expectedSQL := regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS \"version\" (\"database_version\" CHAR(32) PRIMARY KEY, \"version_date\" TIMESTAMP);")
	mock.ExpectPrepare(expectedSQL).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the function being tested
	err = InitDatabaseVersion(session)

	// Assert results
	assert.NoError(t, err)

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
	assert.NoError(t, err)
}

// Integration test - requires a real database
// This test is disabled by default (prefix with _ to enable)
func _TestGetDbSessionIntegration(t *testing.T) {
	// Skip this test in normal test runs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test assumes a local PostgreSQL instance is running
	// with the following credentials (customize as needed)
	dsn := "user=postgres dbname=postgres password=postgres sslmode=disable"

	// Set a timeout for the test
	timeout := time.After(5 * time.Second)
	done := make(chan bool)

	// Run the test in a goroutine
	go func() {
		session, err := GetDbSession(dsn, "postgres")
		if err != nil {
			t.Errorf("Failed to get DB session: %v", err)
			done <- false
			return
		}
		defer session.Close()

		// Try to create the version table
		err = InitDatabaseVersion(session)
		if err != nil {
			t.Errorf("Failed to init database version: %v", err)
			done <- false
			return
		}

		done <- true
	}()

	// Wait for either completion or timeout
	select {
	case <-timeout:
		t.Fatal("Test timed out")
	case success := <-done:
		if !success {
			t.Fatal("Test failed")
		}
	}
}
