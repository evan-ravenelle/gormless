package data

import (
	"fmt"
	"log"
)

func InitDatabaseVersion(session ISession) error {
	dbVersionType := session.Dialect().Char(32)
	versionDateType := session.Dialect().Timestamp()

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
		return fmt.Errorf("failed to create database_version table: %v", err)
	}

	return err
}

func UpsertDbVersion() error {

	return fmt.Errorf("Not implemented")
}
