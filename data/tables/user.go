package tables

import (
	"fmt"
	"log"
	data "main/data"
)

func InitUserTable(session data.Session) error {

	userIdType := data.PsqlSerial
	userFirstType := fmt.Sprintf(data.PsqlVarChar, 32)
	userLastType := fmt.Sprintf(data.PsqlVarChar, 32)
	userEmailType := fmt.Sprintf(data.PsqlVarChar, 64)
	roleIdType := data.PsqlInt

	err := session.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRoleTable := data.Table{Name: "user_role"}
	userRoleColumn := data.Column{Name: "role_id"}

	userRoleFk := data.ForeignKey{
		Table:  &userRoleTable,
		Column: &userRoleColumn,
	}

	userTable := data.Table{
		Name: "user",
		Columns: &[]data.Column{
			{Name: "user_id", Type: &userIdType, PrimaryKey: true},
			{Name: "user_first", Type: &userFirstType},
			{Name: "user_last", Type: &userLastType},
			{Name: "user_email", Type: &userEmailType, Indexed: true},
			{Name: "user_role", Type: &roleIdType, ForeignKey: &userRoleFk},
		},
	}
	err = data.CreateTable(session, userTable)
	if err != nil {
		return fmt.Errorf("Failed to create User table: %v", err)
	}

	return err
}
