package data

import (
	"database/sql"
	"fmt"
	"log"
)

type Session struct {
	*sql.DB
}

func GetDbSession(dsn string) (*Session, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	session := &Session{db}
	return session, nil
}

func InitDatabaseVersion(session Session) error {
	dbVersionType := fmt.Sprintf(PsqlChar, 32)
	versionDateType := PsqlTimestamp

	err := session.DB.Ping()
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

	err = CreateTable(session, dbVersionTable)
	if err != nil {
		return fmt.Errorf("Failed to create User table: %v", err)
	}

	return err
}

func UpsertDbVersion() error {
	return fmt.Errorf("error")
}
