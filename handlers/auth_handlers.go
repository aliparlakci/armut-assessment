package handlers

import (
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Signin(authenticator services.Authenticator, sessions services.SessionCreator, activityLogger services.ActivityLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.AuthForm
		if err := c.Bind(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if _, isLoggedIn := c.Get("user"); isLoggedIn {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already logged in"})
			return
		}

		success, err := authenticator.Authenticate(c.Copy(), creds.Username, creds.Password)
		if err != nil {
			logger.Errorf("Authenticator.Authenticate() raised an error while logging in the user with username: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}

		if !success {
			if err := activityLogger.LogUnsuccesfulSignin(c.Copy(), creds.Username, c.ClientIP()); err != nil {
				logger.Errorf("ActivityLogger.LogUnsuccesfulSignin() raised an error: %v", err.Error())
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "username and password mismatch"})
			return
		}

		sessionId, err := sessions.CreateSession(c.Copy(), creds.Username)
		if err != nil {
			logger.Errorf("SessionService.CreateSession raised an error while creating a new session for user with username: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log in"})
			return
		}

		if err := activityLogger.LogSignin(c.Copy(), creds.Username, c.ClientIP()); err != nil {
			logger.Errorf("ActivityLogger.LogSignin() raised an error: %v", err.Error())
		}

		logger.WithFields(logrus.Fields{"username": creds.Username, "sessionId": sessionId}).Infof("user with username logged in on the session with sessionID")
		c.SetCookie("session", sessionId, 7776000, "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{"result": "logged in"})
		return
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if u, isLoggedIn := c.Get("user"); !isLoggedIn {
			c.String(http.StatusUnauthorized, "")
			return
		} else {
			user = u.(models.User)
		}

		c.JSON(http.StatusOK, gin.H{
			"username": user.Username,
		})
	}
}

func Signout(revoker services.SessionRevoker, activityLogger services.ActivityLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, isLoggedIn := c.Get("user"); !isLoggedIn {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no one is signed in"})
			return
		} else {
			user = u.(models.User)
		}

		sessionId, err := c.Cookie("session")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no one is signed in"})
			return
		}

		if err := revoker.RevokeSession(c.Copy(), sessionId); err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}

		if err := activityLogger.LogSignout(c.Copy(), user.Username, c.ClientIP()); err != nil {
			logger.Errorf("ActivityLogger.LogSignout() raised an error: %v", err.Error())
		}

		logger.WithField("sessionId", sessionId).Infof("user signed out from session with sessionId")
		c.Header("Set-Cookie", fmt.Sprintf("session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;"))
		c.String(http.StatusOK, "")
	}
}