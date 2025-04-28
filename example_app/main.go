package example_app

import (
	"fmt"
	gl "gormless/example_app/gormless"
	"gormless/example_app/user"
	"time"
)

func Main() {
	go gl.InitializeDatabaseTables()
	session, err := gl.GetSession()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = session.Close()
	}()
	if err != nil {
		panic("Couldn't get DB session: " + err.Error())
	}
	println(fmt.Sprintf("Hello World %s", time.Now()))
	newUser, err := user.NewUser(session, "john@example.com", "John", "Doe", user.UserRole{RoleName: user.RoleAdmin})
	if err != nil {
		return
	}
	println(fmt.Sprintf("Hello World %s", newUser.Email))
	for {
		println(fmt.Sprintf("Hello World %s", time.Now()))
		time.Sleep(5 * time.Minute)
	}
}
