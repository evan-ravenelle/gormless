package data

/*
import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type sqlAware struct {
	id    int
	Table string
	Value string
}

type ISqlAware interface {
	Upsert(db *sql.DB) error
	Get(db *sql.DB, id int)
}

func (s *sqlAware) Upsert(column string, value string) error {
	query := fmt.Sprintf("UPDATE %s SET value = ? WHERE id = ?", sa.Table)
	_, err := db.Exec(query, sa.Value, sa.ID)
	return err
}

func (s *sqlAware) Get(column string, value string) error {
	query := "SELECT value FROM table WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&sa.Value)
	// Update the ID of SQLAware struct on successful loading of data
	if err == nil {
		sa.ID = id
	}

}
*/
