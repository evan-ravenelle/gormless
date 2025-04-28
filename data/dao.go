package data

import (
	"database/sql"
	"fmt"
	"strings"
)

type IDAO interface {
	Upsert() error
	Get() (sql.Row, error)
	GetMultiple() error
	Delete() error
}

type DAO[T any] struct {
	ISession
	id     string
	Row    *sql.Row
	Rows   *sql.Rows
	Column Column
	Value  string
	Table  Table
}

// Upsert handles inserting or updating one or more rows
func (dao *DAO[T]) Upsert(rows ...map[string]interface{}) error {
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
	args := make([]interface{}, 0, len(rows)*len(columns))

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
func (dao *DAO[T]) ToRowMap() map[string]interface{} {
	rowMap := make(map[string]interface{})

	for _, col := range *dao.Table.Columns {
		if col.Value != nil {
			rowMap[col.Name] = *col.Value
		}
	}

	return rowMap
}

func (dao *DAO[T]) Get() error {
	query := dao.ISession.Dialect().Sprintd(
		"SELECT FROM %i WHERE %i = %s",
		dao.Table.Name,
		dao.Column.Name,
		dao.Value)
	stmt, err := dao.ISession.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	err = rows.Scan(&dao.Rows)
	return err
}

func (dao *DAO[T]) GetLike() error {
	query := dao.ISession.Dialect().Sprintd(
		"SELECT * FROM %i WHERE %i LIKE %s",
		dao.Table.Name,
		dao.Column.Name,
		dao.id)
	stmt, err := dao.ISession.Prepare(query)
	stmt.Query()
	return err
}

func (dao *DAO[T]) Delete() error {
	query := dao.ISession.Dialect().Sprintd(
		"DELETE FROM %i WHERE id = %s",
		dao.Table.Name,
		dao.id)
	_, err := dao.ISession.Exec(query, dao.id)
	return err
}
