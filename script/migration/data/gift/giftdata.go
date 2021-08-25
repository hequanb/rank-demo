package main

import (
	"boframe/models/giftmodels"
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

	gifts := []interface{}{
		&giftmodels.Gift{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			GiftId:    snowflake.GenId(),
			Name:      "鱼蛋",
			Score:     1000,
		},
		&giftmodels.Gift{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			GiftId:    snowflake.GenId(),
			Name:      "咖啡",
			Score:     10000,
		},
		&giftmodels.Gift{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			GiftId:    snowflake.GenId(),
			Name:      "奶茶",
			Score:     100000,
		},
		&giftmodels.Gift{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			GiftId:    snowflake.GenId(),
			Name:      "火箭",
			Score:     1000000,
		},
	}
	rankDB := mongoI.RankDB()
	results, err := rankDB.Collection("gift").InsertMany(context.Background(), gifts)
	if err != nil {
		log.Fatalf("初始化礼物数据，%s \n", err)
		return
	}
	fmt.Println("insert id: ", results)
	fmt.Println("done")
}
