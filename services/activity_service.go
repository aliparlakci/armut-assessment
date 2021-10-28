package services

import (
	"context"
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ActivityService struct {
	*mongo.Collection
}

type ActivityLogger interface {
	LogSignin(c context.Context, username, ip string) error
	LogSignout(c context.Context, username, ip string) error
	LogUnsuccesfulSignin(c context.Context, username, ip string) error
}

type ActivityFetcher interface {
	Fetch(c context.Context, username string) ([]models.Activity, error)
}

func (a *ActivityService) LogSignin(c context.Context, username, ip string) error {
	return a.log(c, username, ip, "signin")
}

func (a *ActivityService) LogSignout(c context.Context, username, ip string) error {
	return a.log(c, username, ip, "signout")
}

func (a *ActivityService) LogUnsuccesfulSignin(c context.Context, username, ip string) error {
	return a.log(c, username, ip, "fail_signin")
}

func (a *ActivityService) log(c context.Context, username, ip string, event string) error {
	_, err := a.InsertOne(c, models.Activity{
		Username: username,
		Event:    event,
		When:     time.Now(),
		IP:       ip,
	})
	if err != nil {
		return fmt.Errorf("mongo driver raised an error while logging an %v event: %v", event, err.Error())
	}

	return nil
}

func (a *ActivityService) Fetch(c context.Context, username string) ([]models.Activity, error) {
	results := make([]models.Activity, 0)

	cursor, err := a.Collection.Find(
		c,
		bson.M{"username": username},
		options.Find().SetSort(bson.D{{"when", -1}}))
	if err == mongo.ErrNoDocuments {
		return results, nil
	} else if err != nil {
		return nil, fmt.Errorf("mongo driver raised an error while fetching new messages: %v", err.Error())
	}

	for cursor.Next(c) {
		var activity models.Activity
		if err := cursor.Decode(&activity); err != nil {
			return nil, fmt.Errorf("cannot decode the activity: %v", err.Error())
		}

		results = append(results, activity)
	}

	return results, nil
}
