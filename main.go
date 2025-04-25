package main

import (
	"fmt"
	data "main/data"
	tables "main/data/tables"
)

func main() {
	fmt.Println("Hello Woyld")
	fmt.Println("Initializing Database...")
	var session, err = data.GetDbSession("user=evan dbname=gotest sslmode=disable")
	if err != nil {
		println("Couldn't get DB Session:", err.Error())
		return
	}
	err = data.InitDatabaseVersion(*session)
	if err != nil {
		println("Couldn't init DB:", err.Error())
		return
	}
	fmt.Println("Creating UserRole Table")
	err = tables.InitUserRoleTable(*session)
	if err != nil {
		println("Couldn't create UserRole table:", err.Error())
		return
	}
	fmt.Println("Creating User Table")
	err = tables.InitUserTable(*session)
	if err != nil {
		println("Couldn't create User table:", err.Error())
		return
	}
}
