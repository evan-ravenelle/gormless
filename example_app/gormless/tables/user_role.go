package tables

import (
	"fmt"
	"gormless/data"
	"gormless/data/dialect"
	"log"
)

func UserRoleTable() data.TableDef {
	roleIdType := dialect.PsqlSmallSerial
	roleNameType := fmt.Sprintf(dialect.PsqlVarChar, 32)

	return func() data.Table {
		userRoleTable := data.Table{
			Name: "user_role",
			Columns: &[]data.Column{
				{Name: "role_id", PrimaryKey: true, Type: &roleIdType},
				{Name: "role_name", Type: &roleNameType},
			},
		}
		return userRoleTable
	}
}

func InitUserRoleTable(session data.ISession) data.TableInitializer {
	err := session.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return func(def data.TableDef) error {
		userRoleTable := UserRoleTable()
		err = data.CreateTable(session, userRoleTable())
		if err != nil {
			return fmt.Errorf("Failed to create user_role table: %v", err)
		}
		return err
	}

}
