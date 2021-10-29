package main

import (
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/handlers"
	"github.com/aliparlakci/armut-backend-assessment/middlewares"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	godotenv.Load()

	logrus.SetLevel(logrus.DebugLevel)
	// Comment out the next line for performance gain
	// logrus.SetReportCaller(true)

	mdb, close := common.InitializeDb(
		os.Getenv("MDB_URI"),
		os.Getenv("MDB_DBNAME"),
		os.Getenv("MDB_USERNAME"),
		os.Getenv("MDB_PASSWORD"),
	)
	defer close()

	rdbUri := os.Getenv("RDB_URI")
	redis := common.RedisInitializer(rdbUri, "")

	env := &common.Env{}
	{
		env.ActivityService = &services.ActivityService{Collection: mdb.Collection("activity")}
		env.AuthService = &services.AuthService{Collection: mdb.Collection("users")}
		env.UserService = &services.UserService{Collection: mdb.Collection("users")}
		env.SessionService = &services.SessionService{Store: redis(0)}
		env.MessagingService = &services.MessagingService{
			Collection:  mdb.Collection("messages"),
			UserService: env.UserService,
		}
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", uuid.New().String())
		c.Next()
	})
	router.Use(middlewares.Logger())
	router.Use(middlewares.AuthMiddleware(env.UserService, env.SessionService))

	api := router.Group("/api")
	{
		api.GET("/messages", middlewares.Protected(handlers.GetAllMessages(env.MessagingService)))
		api.GET("/messages/new", middlewares.Protected(handlers.GetNewMessages(env.MessagingService)))
		api.GET("/messages/check", middlewares.Protected(handlers.CheckNewMessages(env.MessagingService)))
		api.POST("/messages/send", middlewares.Protected(handlers.SendMessage(env.MessagingService)))
		api.PUT("/messages/read/:id", middlewares.Protected(handlers.ReadMessage(env.MessagingService)))
		api.PUT("/messages/user/read/:username/", middlewares.Protected(handlers.ReadMessages(env.MessagingService)))

		api.POST("/signup", handlers.Signup(env.UserService, env.AuthService))

		api.POST("/signin", handlers.Signin(env.AuthService, env.SessionService, env.ActivityService))
		api.POST("/signout", middlewares.Protected(handlers.Signout(env.SessionService, env.ActivityService)))

		api.GET("/me", handlers.Me())

		api.GET("/activity", middlewares.Protected(handlers.GetActivities(env.ActivityService)))
	}

	router.Run(":5000")
}
