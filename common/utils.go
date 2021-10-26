package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerWithRequestId(c *gin.Context) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"request_id": c.GetString("request_id"), "ip": c.Request.RemoteAddr})
}
