package data

import (
	"database/sql"
	"fmt"
	"log"
)

type Db struct {
	*sql.DB
}

func DbSession(dsn string) (*Db, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

func InitDatabaseVersion() error {
	dbVersionType := fmt.Sprintf(PsqlChar, 32)
	versionDateType := PsqlTimestamp

	db, err := DbSession("user=evanravenelle dbname=gotest sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	dbVersionTable := Table{
		Name: "version",
		Columns: &[]Column{
			{Name: "database_version", PrimaryKey: true, Type: &dbVersionType},
			{Name: "version_date", Type: &versionDateType},
		},
	}
	err = CreateTable(*db, dbVersionTable)
	if err != nil {
		return fmt.Errorf("Failed to create User table: %v", err)
	}

	return err
}

func UpsertDbVersion() error {
	return fmt.Errorf("error")
}
