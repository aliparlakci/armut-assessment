package main

import (
	"context"
	"github.com/aliparlakci/armut-backend-assessment/common"
	"github.com/aliparlakci/armut-backend-assessment/handlers"
	"github.com/aliparlakci/armut-backend-assessment/middlewares"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDb(uri, name, username, password string) (*mongo.Database, func()) {
	c, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client()
	if username != "" {
		credentials := options.Credential{
			Username: username,
			Password: password,
		}
		clientOptions = clientOptions.SetAuth(credentials)
	}

	client, err := mongo.Connect(c, clientOptions.ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client.Database(name), func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func RedisInitializer(uri, password string) func(db int) *redis.Client {
	return func(db int) *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr: uri,
			DB:   db,
		})
	}
}

func main() {
	godotenv.Load()

	logrus.SetLevel(logrus.DebugLevel)
	// Comment out the next line for performance gain
	logrus.SetReportCaller(true)

	mdb, close := InitializeDb(
		os.Getenv("MDB_URI"),
		os.Getenv("MDB_DBNAME"),
		os.Getenv("MDB_USERNAME"),
		os.Getenv("MDB_PASSWORD"),
	)
	defer close()

	rdbUri := os.Getenv("RDB_URI")
	redis := RedisInitializer(rdbUri, "")

	env := &common.Env{
		AuthService:    &services.AuthService{Collection: mdb.Collection("users")},
		UserService:    &services.UserService{Collection: mdb.Collection("users")},
		SessionService: &services.SessionService{Store: redis(0)},
	}
	env.MessagingService = &services.MessagingService{
		Collection:  mdb.Collection("messages"),
		UserService: env.UserService,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("existingUser", ExistingUserValidator(env.UserService))
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
		api.GET("/messages", handlers.GetAllMessages(env.MessagingService))
		api.GET("/messages/new", handlers.GetNewMessages(env.MessagingService))
		api.GET("/messages/check", handlers.CheckNewMessages(env.MessagingService))
		api.POST("/messages/send", handlers.SendMessage(env.MessagingService))

		api.POST("/signup", handlers.Signup(env.UserService, env.UserService))

		api.POST("/signin", handlers.Signin(env.AuthService, env.SessionService))
		api.POST("/signout", handlers.Signout(env.SessionService))
	}

	router.Run(":5000")
}
