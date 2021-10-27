package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Collection *mongo.Collection
}

type UserGetter interface {
	GetUser(c context.Context, username string) (models.User, error)
	UserExists(c context.Context, username string) (bool, error)
}

type UserCreator interface {
	CreateUser(c context.Context, username, password string) (string, error)
}

func (u *UserService) GetUser(c context.Context, username string) (models.User, error) {
	result := u.Collection.FindOne(c, bson.M{"username": username})

	var user models.User
	if err := result.Err(); err != nil 	{
		return user, err
	}
	err := result.Decode(&user)
	return user, err
}

func (u *UserService) UserExists(c context.Context, username string) (bool, error) {
	_, err := u.GetUser(c, username)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserService) CreateUser(c context.Context, username, password string) (string, error) {
	if exists, err := u.UserExists(c, username); err != nil {
		return "", err
	} else if exists {
		return "", ErrUserAlreadyExists
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("cannot hash the password: %v", err.Error())
	}

	result, err := u.Collection.InsertOne(c, models.User{Username: username, Password: string(bytes)})
	return result.InsertedID.(primitive.ObjectID).String(), err
}

var ErrNoUser error = errors.New("no such user exists")
var ErrUserAlreadyExists error = errors.New("user already exists")