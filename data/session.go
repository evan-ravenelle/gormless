package data

import (
	"database/sql"
	"fmt"
	dialect "gormless/data/dialect"
)

// ISession defines the interface for database operations
type ISession interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
	Begin() (*sql.Tx, error)
	Open(dsn string) error
	Ping() error
	Dialect() dialect.Dialect
}
type Session struct {
	DB         *sql.DB
	SQLDialect dialect.Dialect
}

func (s *Session) Dialect() dialect.Dialect {
	return s.SQLDialect
}

func (s *Session) Prepare(query string) (*sql.Stmt, error) {
	return s.DB.Prepare(query)
}

func (s *Session) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args...)
}

func (s *Session) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.DB.QueryRow(query, args...)
}

func (s *Session) Ping() error {
	return s.DB.Ping()
}
func (s *Session) Close() error {
	return s.DB.Close()
}

// Open connects to the database and returns a new session
func (s *Session) Open(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	s.DB = db
	return err
}

func (s *Session) Begin() (*sql.Tx, error) {
	return s.DB.Begin()
}

func (s *Session) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// GetDbSession creates a new database session
func GetDbSession(dsn string, dialectType string) (*Session, error) {
	// Create an empty session
	session := Session{}

	// Set the dialect based on the type
	switch dialectType {
	case "postgres":
		session.SQLDialect = dialect.PostgresDialect{}
	case "mysql":
		session.SQLDialect = dialect.MySQLDialect{}
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", dialectType)
	}

	// Open the connection
	err := session.Open(dsn)
	if err != nil {
		return &session, err
	}

	return &session, nil
}
