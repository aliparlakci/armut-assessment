package handlers

import (
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.String(http.StatusUnauthorized, "")
			return
		} else {
			user = u.(models.User)
		}

		messages, err := getter.GetAllMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.GetAllMessages() raised an error: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": messages})
	}
}

func GetNewMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.String(http.StatusUnauthorized, "")
			return
		} else {
			user = u.(models.User)
		}

		messages, err := getter.GetNewMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.GetNewMessages() raised an error: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": messages})
	}
}

func CheckNewMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.String(http.StatusUnauthorized, "")
			return
		} else {
			user = u.(models.User)
		}

		count, err := getter.CheckNewMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.CheckNewMessages() raised an error: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": count})
	}
}
func SendMessage(sender services.MessageSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.String(http.StatusUnauthorized, "")
			return
		} else {
			user = u.(models.User)
		}

		var message models.NewMessage
		if err := c.Bind(&message); err != nil {
			// TODO: logging
			c.String(http.StatusBadRequest, "")
			return
		}

		if _, err := sender.SendMessage(c.Copy(), message.Body, user.Username, message.To); err == services.ErrNoUser {
			c.JSON(http.StatusBadRequest, gin.H{"result": "user does not exist"})
			return
		} else if err != nil {
			logger.Errorf("services.MessageSender.SendMessage() raised an error: %v", err.Error())
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.String(http.StatusCreated, "")
	}
}
