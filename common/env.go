package common

import "github.com/aliparlakci/armut-backend-assessment/services"

type Env struct {
	*services.AuthService
	*services.MessagingService
	*services.SessionService
	*services.UserService
}
