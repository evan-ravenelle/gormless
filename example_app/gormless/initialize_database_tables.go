package gormless

import (
	"fmt"
	"gormless/data"
	"gormless/example_app/gormless/tables"
)

func InitializeDatabaseTables() {
	session, err := GetSession()
	defer func() {
		err = session.Close()
	}()
	if err != nil {
		panic("Couldn't get DB session: " + err.Error())
	}
	err = data.InitDatabaseVersion(session)
	if err != nil {
		panic("Couldn't init DB: " + err.Error())
	}
	
	fmt.Println("Creating UserRole Table")
	initUserRole := tables.InitUserRoleTable(session)
	err = initUserRole(tables.UserRoleTable())
	if err != nil {
		println("Couldn't create UserRole table:", err.Error())
		return
	}
	fmt.Println("Creating User Table")

	initUser := tables.InitUserTable(session)

	err = initUser(tables.UserTable())

	if err != nil {
		err = fmt.Errorf("Couldn't create User table: %v", err)
	}
}
