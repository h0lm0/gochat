package models

import (
	"strings"
)

const (
	RoleGuest  = "guest"
	RoleAdmin  = "admin"
	RoleSystem = "system"
)

// type User struct {
// 	Session  ssh.Session
// 	Terminal *term.Terminal
// 	SqlUser  *database.User
// }

// func NewUser(session ssh.Session, terminal *term.Terminal, sql *database.User) User {
// 	return User{
// 		Session:  session,
// 		Terminal: terminal,
// 		SqlUser:  sql,
// 	}
// }

func ListUsers() string {
	var sb strings.Builder

	for _, user := range ConnectedUsers {
		sb.WriteString(SuccessStyle.Sprintf(" " + user.Session.User() + "\n"))
	}

	return sb.String()
}
