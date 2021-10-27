package services

import (
	"context"
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	*mongo.Collection
}

type Authenticator interface {
	Authenticate(c context.Context, username, password string) (bool, error)
}

func (a *AuthService) Authenticate(c context.Context, username, password string) (bool, error) {
	result := a.Collection.FindOne(c, bson.M{"username": username})

	var user models.User
	if err := result.Err(); err != nil 	{
		return false, fmt.Errorf("mongodb driver raised an error while fetching the user: %v", err.Error())
	}

	if err := result.Decode(&user); err != nil {
		return false, fmt.Errorf("cannot decode user: %v", err.Error())
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil, nil
}