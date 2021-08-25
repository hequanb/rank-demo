package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const EmptyString = ""

type BaseModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
