package usermodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserId    int64              `bson:"user_id" json:"user_id,string"`
	Name      string             `bson:"name" json:"name"`
	UserType  int                `bson:"user_type" json:"user_type"`
}
