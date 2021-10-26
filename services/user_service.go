package services

import (
	"context"
	"errors"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Collection *mongo.Collection
}

type UserGetter interface {
	GetUser(c context.Context, username string) (models.User, error)
	UserExists(c context.Context, username string) (bool, error)
}

type UserCreator interface {
	CreateUser(c context.Context, username, password string) error
}

func (u *UserService) GetUser(c context.Context, username string) (models.User, error) {
	return models.User{}, nil
}

func (u *UserService) UserExists(c context.Context, username string) (bool, error) {
	return false, nil
}

func (u *UserService) CreateUser(c context.Context, username, password string) error {
	return nil
}

var ErrNoUser error = errors.New("no such user exists")