package tables

import (
	"fmt"
	data "gormless/data"
	dialect "gormless/data/dialect"

	"log"
)

func UserTable() data.TableDef {
	userIdType := dialect.PsqlSerial
	userFirstType := fmt.Sprintf(dialect.PsqlVarChar, 32)
	userLastType := fmt.Sprintf(dialect.PsqlVarChar, 32)
	userEmailType := fmt.Sprintf(dialect.PsqlVarChar, 64)
	roleIdType := dialect.PsqlInt

	userRoleTable := data.Table{Name: "user_role"}
	userRoleColumn := data.Column{Name: "role_id"}

	userRoleFk := data.ForeignKey{
		Table:  &userRoleTable,
		Column: &userRoleColumn,
	}

	return func() data.Table {
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
		return userTable
	}
}
func InitUserTable(session data.ISession) data.TableInitializer {

	err := session.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return func(def data.TableDef) error {
		userTable := UserTable()
		err = data.CreateTable(session, userTable())
		if err != nil {
			return fmt.Errorf("Failed to create User table: %v", err)
		}
		return err
	}
}
