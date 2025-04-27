package example_app

import (
	"fmt"
	"gormless/data"
	"gormless/data/dialect"
	tables2 "gormless/example_app/gormless/tables"
)

func Main() {
	conf, err := data.LoadConfig("example_app/db_config.yml")
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf(
		"user=%s dbname=%s sslmode=%s",
		conf.Database.Username,
		conf.Database.DBName,
		conf.Database.SSLMode)

	fmt.Println(dsn)
	session, err := data.GetDbSession(dsn, dialect.POSTGRES)
	if err != nil {
		panic(err)
	}

	defer func(session *data.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	err = data.InitDatabaseVersion(session)
	if err != nil {
		println("Couldn't init DB:", err.Error())
		return
	}
	fmt.Println("Creating UserRole Table")
	err = tables2.InitUserRoleTable(session)
	if err != nil {
		println("Couldn't create UserRole table:", err.Error())
		return
	}
	fmt.Println("Creating User Table")
	err = tables2.InitUserTable(session)
	if err != nil {
		println("Couldn't create User table:", err.Error())
		return
	}
}
