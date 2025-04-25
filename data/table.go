package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Migration func(table Table, db *sql.DB) error

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
}

type Table struct {
	Name       string
	Columns    *[]Column
	Migrations *[]Migration
}

func CreateTable(session Session, table Table) error {
	var stmt strings.Builder
	countPrimaryKey := 0
	fmt.Fprintf(&stmt, "CREATE TABLE IF NOT EXISTS \"%s\" (", table.Name)
	for i, column := range *table.Columns {
		if column.PrimaryKey {
			countPrimaryKey++
			if countPrimaryKey > 1 {
				return fmt.Errorf("multiple primary keys defined in table: %s", table.Name)
			}
			if !isSafeSQLIdentifier(column.Name) || !isSafeSQLIdentifier(*column.Type) {
				return fmt.Errorf("invalid SQL identifier found")
			}
		}
		fmt.Fprintf(&stmt, "\"%s\" %s", column.Name, *column.Type)
		if column.PrimaryKey {
			fmt.Fprintf(&stmt, " PRIMARY KEY")
		}
		if column.ForeignKey != nil {
			fmt.Fprintf(&stmt, ", FOREIGN KEY (\"%s\") REFERENCES \"%s\"(\"%s\")", column.Name, column.ForeignKey.Table.Name, column.ForeignKey.Column.Name)
		}
		if i != len(*table.Columns)-1 {
			fmt.Fprintf(&stmt, ", ")
		}
	}
	fmt.Fprintf(&stmt, ");")
	fmt.Println(stmt.String())

	statement, err := session.DB.Prepare(stmt.String())
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
	return func(table Table, db *sql.DB) error {
		// Start by adding the column with its type
		query := fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN %s %s",
			table.Name,
			column.Name,
			*column.Type)
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("adding column: %w", err)
		}

		// Create an index if necessary
		if column.Indexed {
			query = fmt.Sprintf("CREATE INDEX idx_%s_on_%s ON %s (%s)", table.Name, column.Name, table.Name, column.Name)
			_, err := db.Exec(query)
			if err != nil {
				return fmt.Errorf("creating index: %w", err)
			}
		}

		// Set as primary key if necessary
		query = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, column.Name, column.Type)
		if column.PrimaryKey {
			query = fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (%s)", table.Name, column.Name)
		}
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("setting primary key: %w", err)
		}

		// Add a foreign key constraint if necessary
		if column.ForeignKey != nil {
			fk := column.ForeignKey
			if fk.Table.Name != "" && fk.Column.Name != "" {
				query = fmt.Sprintf(
					"ALTER TABLE %s ADD FOREIGN KEY (%s) REFERENCES %s (%s)",
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
	return func(table Table, db *sql.DB) error {
		_, err := db.Exec(
			"ALTER TABLE %s DROP COLUMN %s",
			table.Name,
			column.Name)
		if err != nil {

			return fmt.Errorf("removing column: %s: %w", table.Name, err)
		}
		return nil

	}
}

func ModifyColumn(table Table, oldColumn Column, newColumn Column) Migration {
	return func(table Table, db *sql.DB) error {
		// Put rename and type change SQL commands in a transaction.
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if newColumn.Name != "" {
			_, err = tx.Exec(
				fmt.Sprintf(
					"ALTER TABLE %s RENAME COLUMN %s TO %s",
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
				fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", table.Name, newColumn.Name, *newColumn.Type),
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

func isSafeSQLIdentifier(name string) bool {
	// Check if name is non-empty and contains only alphanumeric characters or underscores
	// This is simplistic and might not cover all valid cases depending on the SQL dialect and naming conventions
	isSafe := name != "" && !strings.ContainsAny(name, "';--")
	return isSafe
}
