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
		// Create the table structure
		userRoleTable := UserRoleTable()
		table := userRoleTable()
		err = data.CreateTable(session, table)
		if err != nil {
			return fmt.Errorf("Failed to create user_role table: %v", err)
		}

		// Create a generic DAO for inserting default values
		dao := data.DAO[any]{
			ISession: session,
			Table:    table,
		}

		// Prepare batch data for default roles
		defaultRoles := []map[string]interface{}{
			{"role_name": "admin"},
			{"role_name": "user"},
			{"role_name": "guest"},
		}

		// Insert all default roles at once
		err = dao.Upsert(defaultRoles...)
		if err != nil {
			return fmt.Errorf("Failed to insert default roles: %v", err)
		}

		return nil
	}
}
