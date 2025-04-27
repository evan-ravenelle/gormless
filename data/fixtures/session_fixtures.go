package fixtures

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Prepare(query string) (*sql.Stmt, error) {
	args := m.Called(query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

// MockSession is a mock for the Session type
type MockSession struct {
	*MockDB
}

func (m *MockSession) Prepare(query string) (*sql.Stmt, error) {
	args := m.Called(query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

func (m *MockSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *MockSession) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSession) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}
func (m *MockSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSession) QueryRow(query string, args ...interface{}) *sql.Row {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(*sql.Row)
}

func (m *MockSession) Query(query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(*sql.Rows), mockArgs.Error(1)
}

func (m *MockSession) Open(dsn string) error {
	args := m.Called()
	return args.Error(0)
}

// MockStmt is a mock for sql.Stmt
type MockStmt struct {
	mock.Mock
}

func (m *MockStmt) Exec(args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(args...)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *MockStmt) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockResult is a mock for sql.Result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
