package user

import (
	"database/sql"
	"fmt"
	"gormless/data"
	"gormless/example_app/gormless/tables"
)

type Role struct {
	RoleName string
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Role      string
}

type DAO struct {
	data.DAO[User]
	RoleDAO *RoleDAO
}

type RoleDAO struct {
	data.DAO[Role]
}

func NewUserDAO(session data.ISession) *DAO {
	roleDAO := NewUserRoleDAO(session)

	dao := DAO{
		DAO: data.DAO[User]{
			ISession: session,
			Table:    tables.UserTable()(),
		},
		RoleDAO: roleDAO,
	}
	return &dao
}

func NewUserRoleDAO(session data.ISession) *RoleDAO {
	roleDAO := RoleDAO{data.DAO[Role]{
		ISession: session,
		Table:    tables.UserRoleTable()(),
	}}
	return &roleDAO
}

func (D DAO) UpsertUsers(rows ...map[string]any) error {
	err := D.Upsert(rows...)
	if err != nil {
		return err
	}
	return nil
}

func InsertUserRoles(session data.ISession) error {
	roles := map[string]Role{
		"admin": {
			RoleName: RoleAdmin,
		},
		"user": {
			RoleName: RoleUser,
		},
		"guest": {
			RoleName: RoleGuest,
		},
	}

	dao := NewUserRoleDAO(session)

	dao.ToRowMap(roles)

	err := dao.DAO.Upsert()
	if err != nil {
		return err
	}

	return nil
}

func (d *DAO) GetUsersByColumn(column string, value any) (*sql.Rows, error) {

	rows, err := d.GetMany(column, value, false)
	return rows, err
}

func (user *User) GetUserByEmail() (*sql.Row, error) {

	column := "user_email"
	row, err := user.GetUserByColumn(column, user.Email)
	if err != nil {
		return nil, err
	}
	return row, err
}

func (user *User) GetUsersByRole() (*[]User, error) {

	rows, err := user.GetUsersByColumn("user_role", user.Role)
	if err != nil {
	}
	return rows, err
}

func isValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

func NewUser(session data.ISession, user_email string, user_first string, user_last string, user_role Role) (*User, error) {

	if !isValidRole(user_role.RoleName) {
		return nil, fmt.Errorf("invalid role name")
	}

	user := &User{
		Email:     user_email,
		FirstName: user_first,
		LastName:  user_last,
		Role:      user_role.RoleName,
	}

	dao := NewUserDAO(session)
	dao.Table.Columns = &[]data.Column{
		{Name: "user_first", Type: &user.FirstName},
		{Name: "user_last", Type: &user.LastName},
		{Name: "user_email", Type: &user.Email, Indexed: true},
		{Name: "user_role", Type: &user.Role},
	}
	err := dao.Upsert()
	if err != nil {
		return nil, err
	}
	return user, nil
}
