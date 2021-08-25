package giftdao

import (
	"boframe/models/giftmodels"
	"boframe/settings/mongoI"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func OneByGiftId(giftId int64) (*giftmodels.Gift, error) {
	rankDB := mongoI.RankDB()
	result := rankDB.Collection("gift").FindOne(context.Background(), bson.D{
		{
			"gift_id", giftId,
		},
	})
	if err := result.Err(); err != nil {
		return nil, result.Err()
	}
	gift := new(giftmodels.Gift)
	err := result.Decode(gift)
	return gift, err
}
