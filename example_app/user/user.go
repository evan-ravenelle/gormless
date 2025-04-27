package user

import (
	"fmt"
	"gormless/data"
)

type UserRole struct {
	dao      *data.DAO[UserRole]
	RoleName string
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

type User struct {
	dao       *data.DAO[User]
	FirstName string
	LastName  string
	Email     string
	Role      string
}

func (u *User) GetDAO(session data.ISession) *data.DAO[User] {
	dao := data.DAO[User]{
		ISession: session,
		Table:    data.Table{Name: "user"},
	}
	return &dao
}

func (u *UserRole) GetDAO(session *data.Session) *data.DAO[User] {
	dao := data.DAO[User]{
		ISession: session,
		Table:    data.Table{Name: "user_role"},
	}
	return &dao
}

func (user *User) Get() error {

	return user.dao.Get()
}

func isValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

func NewUser(session data.ISession, user_email string, user_first string, user_last string, user_role UserRole) (*User, error) {

	if !isValidRole(user_role.RoleName) {
		return nil, fmt.Errorf("invalid role name")
	}

	user := &User{
		Email:     user_email,
		FirstName: user_first,
		LastName:  user_last,
		Role:      user_role.RoleName,
	}

	user.GetDAO(session)
	return user, nil
}
