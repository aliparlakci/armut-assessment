package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	To     string             `bson:"to" json:"to"`
	From   string             `bson:"from" json:"from"`
	Body   string             `bson:"body" json:"body"`
	SendAt time.Time          `bson:"send_at" json:"send_at"`
	IsRead bool               `bson:"is_read" json:"is_read"`
}

type NewMessage struct {
	To   string `form:"to" binding:"required"`
	Body string `form:"body" binding:"required"`
}