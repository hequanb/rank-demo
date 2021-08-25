package giftmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Gift struct {
	Id        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	GiftId    int64              `bson:"gift_id" json:"gift_id,string"`
	Name      string             `bson:"name"`
	Score     int64              `bson:"score"`
}
