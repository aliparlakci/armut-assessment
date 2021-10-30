package handlers

import (
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetActivities(fetcher services.ActivityFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		activities, err := fetcher.Fetch(c.Copy(), user.Username)
		if err != nil {
			logger.Errorf("ActivityFetcher.Fetch() raised an error while fetching activity of %v: %v", user.Username, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": activities})
	}
}
