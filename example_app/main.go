package example_app

import (
	"fmt"
	"gormless/data"
	"gormless/data/dialect"
	"gormless/example_app/tables"
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

	err = data.InitDatabaseVersion(session)
	if err != nil {
		println("Couldn't init DB:", err.Error())
		return
	}
	fmt.Println("Creating UserRole Table")
	err = tables.InitUserRoleTable(session)
	if err != nil {
		println("Couldn't create UserRole table:", err.Error())
		return
	}
	fmt.Println("Creating User Table")
	err = tables.InitUserTable(session)
	if err != nil {
		println("Couldn't create User table:", err.Error())
		return
	}
}
