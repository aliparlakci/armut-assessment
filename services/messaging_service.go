package services

import (
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessagingService struct {
	Collection *mongo.Collection
}

type MessageSender interface {
	SendMessage(body, sender, receiver string) error
}

type MessageGetter interface {
	GetAllMessages(username string) ([]models.Message, error)
	GetNewMessages(username string) ([]models.Message, error)
	CheckNewMessages(username string) (int, error)
}

type MessageReader interface {
	ReadMessage(id string) error
}

func (m *MessagingService) GetAllMessages(username string) ([]models.Message, error) {
	return []models.Message{}, nil
}

func (m *MessagingService) GetNewMessages(username string) ([]models.Message, error) {
	return []models.Message{}, nil
}

func (m *MessagingService) CheckNewMessages(username string) (int, error) {
	return 0, nil
}

func (m *MessagingService) SendMessage(body, sender, receiver string) error {
	return nil
}