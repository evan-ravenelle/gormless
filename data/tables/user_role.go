package tables

import (
	"fmt"
	"log"
	"main/data"
)

func InitUserRoleTable(session data.Session) error {
	roleIdType := data.PsqlSmallSerial
	roleNameType := fmt.Sprintf(data.PsqlVarChar, 32)

	err := session.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRoleTable := data.Table{
		Name: "user_role",
		Columns: &[]data.Column{
			{Name: "role_id", PrimaryKey: true, Type: &roleIdType},
			{Name: "role_name", Type: &roleNameType},
		},
	}
	err = data.CreateTable(session, userRoleTable)
	if err != nil {
		return fmt.Errorf("Failed to create User table: %v", err)
	}

	return err
}
