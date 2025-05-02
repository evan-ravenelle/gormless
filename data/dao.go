package data

import (
	"database/sql"
	"fmt"
	"gormless/example_app/user"
	"strings"
)

type IDAO interface {
	Upsert() error
	Get() (sql.Row, error)
	Delete() error
}

// DAO is a generic data access object (DAO) that provides CRUD operations for a specific type of data.
// It is used to abstract away the underlying database operations and provide a consistent interface for
// accessing data.
//
// The DAO is generic over the type of data it manages. This allows the DAO to be used with any type of data,
// such as a user, product, or order.
//
// The DAO is also generic over the type of session it uses to access the database. This allows the DAO to be
// used with any type of session, such as a SQL session or a NoSQL session.
//
// Only ISession and Table are required to create a DAO. The other fields are optional and are used to
// implement the DAO interface.
//
// The DAO is designed to be used with the data package. It is not intended to be used directly.
//
// E.g.,
// 	type User struct {
// 		ID        int
// 		FirstName string
// 		LastName  string
// 		Email     string
// 	}
//
// 	type UserDAO struct {
// 		data.DAO[User]
// 	}
//
// func (dao *UserDAO) GetDAO(session data.ISession) data.DAO[User] {
// 	dao.ISession = session
// 	dao.Table = tables.UserTable()()
// 	return dao
// }
//
// 	func (dao *UserDAO) Get(id int) (*User, error) {
// 		row, err := dao.GetOne("id", id)
// 		if err != nil {
// 			return nil, err
// }

type DAO[T any] struct {
	ISession           // Links the DAO to a specific session
	Table              // Links the DAO to a specific table
	Row      *sql.Row  // Used for GetOne()
	Rows     *sql.Rows // Used for GetMany()
	Columns  []Column  // Used for Upsert()
	Values   []string  // Used for Upsert()
	id       string    // Used for Delete()
}

type DAOFactory struct {
	session ISession
}

// NewDAOFactory creates a new factory with the given session
func NewDAOFactory(session ISession) *DAOFactory {
	return &DAOFactory{session: session}
}

// Upsert handles inserting or updating one or more rows
func (dao *DAO[T]) Upsert(rows ...map[string]any) error {
	if len(rows) == 0 {
		return nil // Nothing to do
	}

	dialect := dao.ISession.Dialect()

	// Get column names from the first row (assuming all rows have the same columns)
	firstRow := rows[0]
	columns := make([]string, 0, len(firstRow))
	for col := range firstRow {
		// Don't skip ID column - let the database handle it during conflict resolution
		columns = append(columns, dialect.QuoteIdentifier(col))
	}

	// Build the bulk insert query
	var builder strings.Builder
	builder.WriteString(dialect.Sprintd("INSERT INTO %i (", dao.Table.Name))
	builder.WriteString(strings.Join(columns, ", "))
	builder.WriteString(") VALUES ")

	// Add placeholders for each row
	placeholderGroups := make([]string, len(rows))
	args := make([]any, 0, len(rows)*len(columns))

	for i, row := range rows {
		placeholders := make([]string, len(columns))

		for j, col := range columns {
			// Remove any quotes added by QuoteIdentifier
			plainCol := strings.Trim(col, "\"'`[]")
			placeholders[j] = dialect.Placeholder(i*len(columns) + j + 1)
			args = append(args, row[plainCol])
		}

		placeholderGroups[i] = "(" + strings.Join(placeholders, ", ") + ")"
	}

	builder.WriteString(strings.Join(placeholderGroups, ", "))

	// Add ON CONFLICT clause for PostgreSQL (upsert)
	// Assuming the first column with PrimaryKey=true is the conflict target
	var primaryKeyCol string
	for _, col := range *dao.Table.Columns {
		if col.PrimaryKey {
			primaryKeyCol = col.Name
			break
		}
	}

	if primaryKeyCol != "" {
		builder.WriteString(dialect.Sprintd(" ON CONFLICT (%i) DO UPDATE SET ", primaryKeyCol))

		updateClauses := make([]string, 0, len(columns))
		for _, col := range columns {
			plainCol := strings.Trim(col, "\"'`[]")
			if plainCol != primaryKeyCol {
				updateClauses = append(updateClauses, dialect.Sprintd(
					"%s = EXCLUDED.%s",
					col,
					dialect.QuoteIdentifier(plainCol),
				))
			}
		}

		builder.WriteString(strings.Join(updateClauses, ", "))
	}

	// Execute the query
	query := builder.String()
	_, err := dao.ISession.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("upsert failed: %w", err)
	}

	return nil
}

// ToRowMap is a helper method to convert DAO columns with values to a row map
func (dao *DAO[T]) ToRowMap(map[string]user.Role) map[string]any {
	rowMap := make(map[string]any)

	for _, col := range *dao.Table.Columns {
		if col.Value != nil {
			rowMap[col.Name] = *col.Value
		}
	}
	return rowMap
}

func (dao *DAO[T]) GetOne(columnName string, value any) (*sql.Row, error) {
	dialect := dao.ISession.Dialect()
	query := dialect.Sprintd(
		"SELECT * FROM %i WHERE %i = %s",
		dao.Table.Name,
		columnName,
		dialect.Placeholder(1))

	stmt, err := dao.ISession.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.QueryRow(value), nil
}

// GetMany retrieves multiple rows by column match or pattern
func (dao *DAO[T]) GetMany(columnName string, value any, patternMatch bool) (*sql.Rows, error) {
	dialect := dao.ISession.Dialect()

	operator := "="
	if patternMatch {
		operator = "LIKE"
	}

	query := dialect.Sprintd(
		"SELECT * FROM %i WHERE %i "+operator+" %s",
		dao.Table.Name,
		columnName,
		dialect.Placeholder(1))

	stmt, err := dao.ISession.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(value)
}

func (dao *DAO[T]) Delete() error {
	query := dao.ISession.Dialect().Sprintd(
		"DELETE FROM %i WHERE id = %s",
		dao.Table.Name,
		dao.id)
	_, err := dao.ISession.Exec(query, dao.id)
	return err
}
