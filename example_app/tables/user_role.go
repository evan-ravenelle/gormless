package tables

import (
	"fmt"
	"gormless/data"
	"gormless/data/dialect"
	"log"
)

func InitUserRoleTable(session data.ISession) error {
	roleIdType := dialect.PsqlSmallSerial
	roleNameType := fmt.Sprintf(dialect.PsqlVarChar, 32)

	err := session.Ping()
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
		return fmt.Errorf("Failed to create user_role table: %v", err)
	}
	
	return err
}
