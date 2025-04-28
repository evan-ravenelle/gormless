package data

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gormless/data/sqlsafe"
	"log"
	"strings"
)

type TableDef func() Table

type TableInitializer func(tableDef TableDef) error

type Table struct {
	Name       string
	Columns    *[]Column
	Migrations *[]Migration
}

type ForeignKey struct {
	Table  *Table
	Column *Column
}

type Column struct {

	// Name is a string that represents the name of a column in a database table.
	Name       string
	Type       *string
	Indexed    bool
	PrimaryKey bool
	ForeignKey *ForeignKey
	Value      *string
}

type Migration func(table Table, session ISession) error

func CreateTable(session ISession, table Table) error {
	var stmt strings.Builder
	dialect := session.Dialect()
	countPrimaryKey := 0
	_, err := dialect.Fprintd(&stmt, "CREATE TABLE IF NOT EXISTS %i (", table.Name)
	if err != nil {
		return err
	}
	for i, column := range *table.Columns {
		if column.PrimaryKey {
			countPrimaryKey++
			if countPrimaryKey > 1 {
				return errors.New(fmt.Sprintf("multiple primary keys defined in table: %s", table.Name))
			}
			if !sqlsafe.IsSafeSQLString(column.Name) || !sqlsafe.IsSafeSQLString(*column.Type) {
				return errors.New("invalid SQL identifier found")
			}
		}
		dialect.Fprintd(&stmt, "%i %s", column.Name, *column.Type)
		if column.PrimaryKey {
			dialect.Fprintd(&stmt, " PRIMARY KEY")
		}
		if column.ForeignKey != nil {
			dialect.Fprintd(&stmt, ", FOREIGN KEY (%i) REFERENCES %i(%i)", column.Name, column.ForeignKey.Table.Name, column.ForeignKey.Column.Name)
		}
		if i != len(*table.Columns)-1 {
			dialect.Fprintd(&stmt, ", ")
		}
	}
	dialect.Fprintd(&stmt, ");")
	fmt.Println(stmt.String())
	if !sqlsafe.IsSafeSQLString(stmt.String()) {
		return errors.New("invalid SQL identifier found")
	}
	statement, err := session.Prepare(stmt.String())

	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("execution error: ", err)
	}
	return err
}

func AddColumn(table Table, column Column) Migration {
	return func(table Table, db ISession) error {
		dialect := db.Dialect()
		// Start by adding the column with its type
		query := dialect.Sprintd(
			"ALTER TABLE %i ADD COLUMN %i %s",
			table.Name,
			column.Name,
			*column.Type)
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("adding column: %w", err)
		}

		// Create an index if necessary
		if column.Indexed {
			query = dialect.Sprintd("CREATE INDEX idx_%i_on_%i ON %i (%i)", table.Name, column.Name, table.Name, column.Name)
			_, err := db.Exec(query)
			if err != nil {
				return fmt.Errorf("creating index: %w", err)
			}
		}

		// Set as primary key if necessary
		query = "ALTER TABLE %i ADD COLUMN %i %i"
		if column.PrimaryKey {
			query = "ALTER TABLE %i ADD PRIMARY KEY (%i)"
		}
		_, err = db.Exec(dialect.Sprintd(query, table.Name, column.Name, *column.Type))
		if err != nil {
			return fmt.Errorf("setting primary key: %w", err)
		}

		// Add a foreign key constraint if necessary
		if column.ForeignKey != nil {
			fk := column.ForeignKey
			if fk.Table.Name != "" && fk.Column.Name != "" {
				query = dialect.Sprintd(
					"ALTER TABLE %i ADD FOREIGN KEY (%i) REFERENCES %i (%i)",
					table.Name,
					column.Name,
					fk.Table.Name,
					fk.Column.Name)

				_, err := db.Exec(query)
				if err != nil {
					return fmt.Errorf("setting foreign key: %w", err)
				}
			}
		}

		return nil
	}
}

func RemoveColumn(column Column) Migration {
	return func(table Table, db ISession) error {
		dialect := db.Dialect()
		_, err := db.Exec(
			dialect.Sprintd("ALTER TABLE %i DROP COLUMN %i"),
			table.Name,
			column.Name)
		if err != nil {

			return fmt.Errorf("removing column: %s: %w", table.Name, err)
		}
		return nil

	}
}

func ModifyColumn(table Table, oldColumn Column, newColumn Column) Migration {
	return func(table Table, db ISession) error {
		dialect := db.Dialect()
		// Put rename and type change SQL commands in a transaction.
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		alterTableColumnName := "\"ALTER TABLE %i RENAME COLUMN %i TO %i\""
		alterTableColumnType := "\"ALTER TABLE %i ALTER COLUMN %i TYPE %s\""

		if newColumn.Name != "" {
			_, err = tx.Exec(
				dialect.Sprintd(
					alterTableColumnName,
					table.Name,
					oldColumn.Name,
					newColumn.Name,
				),
			)
			if err != nil {
				return fmt.Errorf("renaming column: %w", err)
			}
		}
		if newColumn.Type != nil {
			_, err = tx.Exec(
				dialect.Sprintd(alterTableColumnType, table.Name, newColumn.Name, *newColumn.Type),
			)
			if err != nil {
				return fmt.Errorf("modifying column type: %w", err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}
}
