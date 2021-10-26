package services

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	*mongo.Collection
}

type Authenticator interface {
	Authenticate(c context.Context, username, password string) (bool, error)
}

func (a *AuthService) Authenticate(c context.Context, username, password string) (bool, error) {
	return true, nil
}