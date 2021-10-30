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
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		messages, err := getter.GetAllMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.GetAllMessages() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
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
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		messages, err := getter.GetNewMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.GetNewMessages() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
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
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		count, err := getter.CheckNewMessages(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("services.MessageGetter.CheckNewMessages() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
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
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		var message models.NewMessage
		if err := c.Bind(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		if _, err := sender.SendMessage(c.Copy(), message.Body, user.Username, message.To); err == services.ErrNoUser {
			c.JSON(http.StatusBadRequest, gin.H{"result": "user does not exist"})
			return
		} else if err != nil {
			logger.Errorf("services.MessageSender.SendMessage() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"result": "message is successfully sent"})
	}
}

func ReadMessage(reader services.MessageReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)
		var user models.User

		if u, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		messageId := c.Param("id")
		if messageId == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		if err := reader.ReadMessage(c.Copy(), messageId, user.Username); err != nil {
			logger.Errorf("MessageReader.ReadMessage() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}


func ReadMessages(reader services.MessageReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)
		var user models.User

		if u, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		senderUsername := c.Param("username")
		if senderUsername == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		if err := reader.ReadMessagesFromUser(c.Copy(), senderUsername, user.Username); err != nil {
			logger.Errorf("MessageReader.ReadMessage() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}