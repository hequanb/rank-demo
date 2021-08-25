package giftdao

import (
	"boframe/models"
	"boframe/models/giftmodels"
	"boframe/settings/mongoI"
	"boframe/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertOne(log *giftmodels.GiftLog) (id string, err error) {
	rankDB := mongoI.RankDB()
	result, err := rankDB.Collection("gift_log").InsertOne(context.Background(), log)
	if err == nil {
		return result.InsertedID.(primitive.ObjectID).String(), nil
	}
	return models.EmptyString, err
}

func PagerByAnchorIds(page, limit int64, condition, sort bson.D) ([]*giftmodels.GiftLog, error) {
	rankDB := mongoI.RankDB()
	offset := utils.CalPageOffset(page, limit)
	cursor, err := rankDB.Collection("gift_log").Find(context.Background(), condition, options.Find().
		SetLimit(limit).
		SetSort(sort).
		SetSkip(offset).
		SetSort(sort),
	)
	if err != nil {
		return nil, err
	}

	var logs []*giftmodels.GiftLog
	err = cursor.All(context.Background(), &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func CountByCondition(condition bson.D) (int64, error) {
	rankDB := mongoI.RankDB()
	return rankDB.Collection("gift_log").CountDocuments(context.Background(), condition)
}