package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
