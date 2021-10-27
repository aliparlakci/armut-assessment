package handlers

import (
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Signup(usersCreator services.UserCreator, userGetter services.UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.AuthForm
		if err := c.Bind(&creds); err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}

		_, err := usersCreator.CreateUser(c.Copy(), creds.Username, creds.Password);
		if err == services.ErrUserAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		if err != nil {
			logger.WithField("username", creds.Username).Errorf("UserService.CreateUSer raised an error while creating user with username: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusCreated, gin.H{"result": "user created"})
	}
}