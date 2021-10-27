package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	To     string             `bson:"to"`
	From   string             `bson:"from"`
	Body   string             `bson:"body"`
	SendAt time.Time          `bson:"send_at"`
	IsRead bool               `bson:"is_read"`
}

type NewMessage struct {
	To   string `form:"to"`
	Body string `form:"body"`
}
