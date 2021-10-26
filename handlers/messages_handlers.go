package handlers

import (
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}
}

func GetNewMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}
}

func CheckNewMessages(getter services.MessageGetter) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}
}
func SendMessage(sender services.MessageSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var message models.NewMessage
		if err := c.ShouldBind(&message); err != nil {
			// TODO: logging
			c.String(http.StatusInternalServerError, "")
			return
		}

		if err := sender.SendMessage(message.Body, message.To, message.To); err != nil {
			// TODO: logging
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.String(http.StatusCreated, "")
	}
}
