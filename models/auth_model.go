package models

type AuthForm struct {
	Username string `form:"username"`
	Password string	`form:"password"`
}
