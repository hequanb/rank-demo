package giftmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GiftLog struct {
	Id         primitive.ObjectID `bson:"_id"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	GiftId     int64              `bson:"gift_id"`
	Name       string             `bson:"name"`
	Score      int64              `bson:"score"`
	FromUserId int64              `bson:"from_user_id"`
	ToUserId   int64              `bson:"to_user_id"`
	SendAt     time.Time          `bson:"send_at"`
}
