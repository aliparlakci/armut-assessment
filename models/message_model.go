package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID     primitive.ObjectID
	To     string
	From   string
	Body   string
	SendAt time.Time
	IsRead bool
}

type NewMessage struct {
	To   string `form:"to"`
	Body string `form:"body"`
}
