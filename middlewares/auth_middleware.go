package middlewares

import (
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(userGetter services.UserGetter, sessions services.SessionFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		sessionId, err := c.Cookie("session")
		if err == http.ErrNoCookie {
			c.Next()
			return
		} else if err != nil {
			logger.Errorf("cannot extract session cookie from request headers: %v", err)
			c.Next()
			return
		}

		username, err := sessions.FetchSession(c.Copy(), sessionId)
		if err != nil {
			logger.WithField("session_id", sessionId).Errorf("sessions.FetchSession raised an error when fetching session with session_id: %v", err.Error())
			c.Header("Set-Cookie", "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;")
			c.Next()
			return
		}

		user, err := userGetter.GetUser(c.Copy(), username)
		if err == services.ErrNoUser {
			logger.WithField("username", username).Debug("user with username does not exist")
			c.Header("Set-Cookie", "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;")
			c.Next()
			return
		}
		if err != nil {
			logger.WithField("username", username).Errorf("UserGetter.GetUser() raised an error while finding user with username: %v", err.Error())
			c.Header("Set-Cookie", "session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;")
			c.Next()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

