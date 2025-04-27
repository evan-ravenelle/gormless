package data

import (
	"fmt"
	"gormless/data/dialect"
	"log"
)

func InitDatabaseVersion(session ISession) error {
	dbVersionType := fmt.Sprintf(dialect.PsqlChar, 32)
	versionDateType := dialect.PsqlTimestamp

	err := session.Ping()
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
