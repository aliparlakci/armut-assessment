package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activity struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Event    string             `bson:"event" json:"event"`
	Username string             `bson:"username" json:"username"`
	IP       string             `bson:"ip" json:"ip"`
	When     time.Time          `bson:"when" json:"when"`
}
