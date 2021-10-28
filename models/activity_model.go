package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activity struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Event    string             `bson:"event"`
	Username string             `bson:"username"`
	IP       string             `bson:"ip"`
	When     time.Time          `bson:"when"`
}
