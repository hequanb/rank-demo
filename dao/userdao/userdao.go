package userdao

import (
	"boframe/models/usermodels"
	"boframe/settings/mongoI"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func OneByUserId(userId int64) (*usermodels.User, error) {
	rankDB := mongoI.RankDB()
	result := rankDB.Collection("user").FindOne(context.Background(), bson.D{
		{
			"user_id", userId,
		},
	})
	if err := result.Err(); err != nil {
		return nil, result.Err()
	}
	user := new(usermodels.User)
	err := result.Decode(user)
	return user, err
}

func ListByUserIds(usersId []int64) ([]*usermodels.User, error) {
	res := make([]*usermodels.User, 0, len(usersId))
	if len(usersId) <= 0 {
		return res, nil
	}
	rankDB := mongoI.RankDB()
	cursor, err := rankDB.Collection("user").Find(context.Background(),
		bson.M{
			"user_id": bson.M{"$in": usersId},
		},
	)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func MapByUserIds(usersId []int64) (map[int64]*usermodels.User, error) {
	res := make(map[int64]*usermodels.User, len(usersId))
	if len(usersId) <= 0 {
		return res, nil
	}
	list, err := ListByUserIds(usersId)
	if err != nil {
		return nil, err
	}
	for _, elem := range list {
		res[elem.UserId] = elem
	}
	return res, nil
}
