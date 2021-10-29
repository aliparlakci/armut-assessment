package services

//go:generate mockgen -destination=../mocks/mock_messaging_service.go -package=mocks github.com/aliparlakci/armut-backend-assessment/services MessageSender,MessageReader,MessageGetter

import (
	"context"
	"fmt"
	"github.com/aliparlakci/armut-backend-assessment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MessagingService struct {
	*mongo.Collection
	*UserService
}

type MessageSender interface {
	SendMessage(c context.Context, body, sender, receiver string) (string, error)
}

type MessageGetter interface {
	GetAllMessages(c context.Context, username string) ([]models.Message, error)
	GetNewMessages(c context.Context, username string) ([]models.Message, error)
	CheckNewMessages(c context.Context, username string) (int, error)
}

type MessageReader interface {
	ReadMessage(c context.Context, id, username string) error
	ReadMessagesFromUser(c context.Context, receiver, sender string) error
}

func (m *MessagingService) GetAllMessages(c context.Context, username string) ([]models.Message, error) {
	results := make([]models.Message, 0)

	cursor, err := m.Collection.Find(c, bson.D{
		{"$or",
			bson.A{
				bson.D{{"from", username}},
				bson.D{{"to", username}},
			}},
	}, options.Find().SetSort(bson.D{{"send_at", -1}}))
	if err != nil {
		return results, fmt.Errorf("mongo driver raised an error while fetching new messages: %v", err.Error())
	}

	for cursor.Next(c) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return results, fmt.Errorf("cannot decode the fetched message: %v", err.Error())
		}

		results = append(results, message)
	}

	return results, nil
}

func (m *MessagingService) GetNewMessages(c context.Context, username string) ([]models.Message, error) {
	results := make([]models.Message, 0)

	cursor, err := m.Collection.Find(c, bson.M{"to": username, "is_read": false})
	if err != nil {
		return results, fmt.Errorf("mongo driver raised an error while fetching new messages: %v", err.Error())
	}

	for cursor.Next(c) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return results, fmt.Errorf("cannot decode the fetched message: %v", err.Error())
		}

		results = append(results, message)
	}

	return results, nil
}

func (m *MessagingService) CheckNewMessages(c context.Context, username string) (int, error) {
	result, err := m.Collection.CountDocuments(c, bson.M{"to": username, "is_read": false})
	if err != nil {
		return 0, err
	}

	return int(result), nil
}

func (m *MessagingService) SendMessage(c context.Context, body, sender, receiver string) (string, error) {
	senderExists, err := m.UserExists(c, sender)
	if err != nil {
		return "", err
	}
	receiverExists, err := m.UserExists(c, receiver)
	if err != nil {
		return "", err
	}

	if !senderExists || !receiverExists {
		return "", ErrNoUser
	}

	if result, err := m.Collection.InsertOne(c, models.Message{
		From:   sender,
		To:     receiver,
		Body:   body,
		IsRead: false,
		SendAt: time.Now(),
	}); err != nil {
		return "", fmt.Errorf("mongo driver raised an error while inserting a new message: %v", err.Error())
	} else {
		return result.InsertedID.(primitive.ObjectID).String(), nil
	}
}

func (m *MessagingService) ReadMessage(c context.Context, id, receiver string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("cannot convert id to ObjectID: %v", err.Error())
	}

	_, err = m.Collection.UpdateMany(
		c,
		bson.M{"_id": objID, "to": receiver},
		bson.D{
			{"$set", bson.D{{"is_read", true}}},
		},
	)

	return err
}

func (m *MessagingService) ReadMessagesFromUser(c context.Context, receiver, sender string) error {
	_, err := m.Collection.UpdateMany(
		c,
		bson.M{"from": sender, "to": receiver},
		bson.D{
			{"$set", bson.D{{"is_read", true}}},
		},
	)

	return err
}
