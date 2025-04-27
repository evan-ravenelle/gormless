package tables

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gormless/data"
	"gormless/data/dialect"
	"regexp"
	"testing"
)

func TestInitUserTable(t *testing.T) {
	// Create a mock SQL driver and connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	session := &data.Session{DB: db}
	session.SQLDialect = dialect.PostgresDialect{}
	// Set expectations for the database operations
	mock.ExpectPing()

	expectedSQL := regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS \"user\" (\"user_id\" SERIAL PRIMARY KEY, \"user_first\" VARCHAR(32), \"user_last\" VARCHAR(32), \"user_email\" VARCHAR(64), \"user_role\" INTEGER, FOREIGN KEY (\"user_role\") REFERENCES \"user_role\"(\"role_id\"));")
	// We'll check that the SQL query for creating the user table is executed
	mock.ExpectPrepare(expectedSQL).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the function being tested
	err = InitUserTable(session)

	// Assert results
	assert.NoError(t, err)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestInitUserRoleTable(t *testing.T) {
	// Create a mock SQL driver and connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	session := &data.Session{DB: db}
	session.SQLDialect = dialect.PostgresDialect{}
	// Set expectations for the database operations
	mock.ExpectPing()

	expectedSQL := regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS \"user_role\" (\"role_id\" SMALLSERIAL PRIMARY KEY, \"role_name\" VARCHAR(32));")
	// We'll check that the SQL query for creating the user_role table is executed
	mock.ExpectPrepare(expectedSQL).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the function being tested
	err = InitUserRoleTable(session)

	// Assert results
	assert.NoError(t, err)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// TestUserTableIntegration is an integration test that requires a real database
// This test is disabled by default (remove the leading underscore to enable)
func _TestUserTableIntegration(t *testing.T) {
	// Skip this test in normal test runs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test assumes a local PostgreSQL instance is running
	// with the following credentials (customize as needed)
	dsn := "user=evan dbname=gotest sslmode=disable"

	// Get a real database connection
	session, err := data.GetDbSession(dsn, "postgres")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer session.Close()

	// Create the user_role table first (since user table depends on it)
	err = InitUserRoleTable(session)
	assert.NoError(t, err)

	// Then create the user table
	err = InitUserTable(session)
	assert.NoError(t, err)

	// Verify the tables exist by querying the information schema
	var userTableExists bool
	err = session.QueryRow(
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user')").
		Scan(&userTableExists)
	assert.NoError(t, err)
	assert.True(t, userTableExists)

	var userRoleTableExists bool
	err = session.QueryRow(
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_role')").
		Scan(&userRoleTableExists)
	assert.NoError(t, err)
	assert.True(t, userRoleTableExists)

	// Clean up - drop the tables
	_, err = session.Exec("DROP TABLE IF EXISTS \"user\"")
	assert.NoError(t, err)
	_, err = session.Exec("DROP TABLE IF EXISTS \"user_role\"")
	assert.NoError(t, err)
}
