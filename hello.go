package main

import (
	"fmt"
	data "main/data"
	tables "main/data/tables"
)

func main() {
	fmt.Println("Hello Woyld")
	fmt.Println("Initializing Database...")
	data.InitDatabaseVersion()
	fmt.Println("Creating UserRole Table")
	tables.InitUserRoleTable()
	fmt.Println("Creating User Table")
	tables.InitUserTable()
}
