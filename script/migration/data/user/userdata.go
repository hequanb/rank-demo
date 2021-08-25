package main

import (
	"boframe/models/usermodels"
	"boframe/pkg/snowflake"
	"boframe/settings"
	"boframe/settings/mongoI"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func main() {
	// 初始化配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init setting failed: %v \n", err)
		return
	}

	// snowflake
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId); err != nil {
		fmt.Printf("init snowflake failed: %v \n", err)
		return
	}

	// Mongo init
	if err := mongoI.Init(settings.Conf.MongoConfig); err != nil {
		fmt.Printf("init mongo failed: %v \n", err)
		return
	}
	defer mongoI.Close()

	users := []interface{}{
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "用户1",
			UserType:  1,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "用户2",
			UserType:  1,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "用户3",
			UserType:  1,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "用户4",
			UserType:  1,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "用户5",
			UserType:  1,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "主播1",
			UserType:  2,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "主播2",
			UserType:  2,
		},
		&usermodels.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    snowflake.GenId(),
			Name:      "主播3",
			UserType:  2,
		},
	}
	rankDB := mongoI.RankDB()
	results, err := rankDB.Collection("user").InsertMany(context.Background(), users)
	if err != nil {
		log.Fatalf("初始化用户数据，%s \n", err)
		return
	}
	fmt.Println("insert id: ", results)
	fmt.Println("done")
}
