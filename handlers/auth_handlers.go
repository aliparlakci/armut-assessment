package handlers

import (
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Signin(authenticator services.Authenticator, sessions services.SessionCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		//logger := common.LoggerWithRequestId(c.Copy())

		var creds models.AuthForm
		if err := c.Bind(&creds); err != nil {
			//logger.Debugf("cannot bind request body to models.SignInRequestForm: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if _, isLoggedIn := c.Get("user"); isLoggedIn {
			//logger.WithField("user_id", user.(models.User).ID).Debug("user with user_id is already logged in on this session")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already logged in"})
			return
		}

		success, err := authenticator.Authenticate(c.Copy(), creds.Username, creds.Password)
		if err != nil {
			// TODO log
			c.String(http.StatusInternalServerError, "")
			return
		}

		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username and password mismatch"})
			return
		}

		sessionId, err := sessions.CreateSession(c.Copy(), creds.Username)
		if err != nil {
			//logger.WithField("user_id", user.ID.Hex()).Errorf("SessionService.CreateSession raised an error while creating a new session for user with user_id: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log in"})
			return
		}

		c.SetCookie("session", sessionId, 7776000, "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{"result": "logged in"})
		return
	}
}

func Signout(revoker services.SessionRevoker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, isLoggedIn := c.Get("user"); !isLoggedIn {
			//logger.WithField("user_id", user.(models.User).ID).Debug("user with user_id is already logged in on this session")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already logged out"})
			return
		}

		sessionId, err := c.Cookie("session")
		if err != nil {
			// logger.Debug("no user is logged in")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already logged out"})
			return
		}

		if err := revoker.RevokeSession(c.Copy(), sessionId); err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}

		c.Header("Set-Cookie", fmt.Sprintf("session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;"))
		c.String(http.StatusOK, "")
	}
}