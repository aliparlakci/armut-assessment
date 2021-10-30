package services

//go:generate mockgen -destination=../mocks/mock_session_service.go -package=mocks github.com/aliparlakci/armut-backend-assessment/services SessionFetcher,SessionCreator,SessionRevoker

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type SessionService struct {
	Store *redis.Client
}

type SessionFetcher interface {
	FetchSession(c context.Context, sessionId string) (string, error)
}

type SessionCreator interface {
	CreateSession(c context.Context, data string) (string, error)
}

type SessionRevoker interface {
	RevokeSession(c context.Context, sessionId string) error
}

func (s *SessionService) FetchSession(c context.Context, sessionId string) (string, error) {
	session, err := s.Store.Get(c, sessionId).Result()
	if err == redis.Nil {
		return session, ErrNoSession
	} else if err != nil {
		return session, err
	}
	return session, nil
}

func (s *SessionService) CreateSession(c context.Context, data string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("cannot create new uuid: %v", err.Error())
	}
	// TODO: Set expiration
	if err := s.Store.Set(c, id.String(), data, 0).Err(); err != nil {
		return "", fmt.Errorf("cannot create a new session: %v", err.Error())
	}

	return id.String(), nil
}

func (s *SessionService) RevokeSession(c context.Context, sessionId string) error {
	return s.Store.Del(c, sessionId).Err()
}

var ErrNoSession error = fmt.Errorf("session does not exist")