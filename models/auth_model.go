package models

type AuthForm struct {
	Username string `form:"username" binding:"required"`
	Password string	`form:"password" binding:"required"`
}
