package gormless

import (
	"fmt"
	"gormless/data"
	"gormless/data/dialect"
)

func GetSession() (*data.Session, error) {
	conf, err := data.LoadConfig("example_app/db_config.yml")
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s",
		conf.Database.Username,
		conf.Database.DBName,
		conf.Database.SSLMode)

	session, err := data.GetDbSession(dsn, dialect.POSTGRES)
	if err != nil {
		return session, err
	}

	return session, err
}
