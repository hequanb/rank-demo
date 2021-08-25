package main

import (
	"boframe/settings/mongoI"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"boframe/settings"
)

func main() {
	// 初始化配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init setting failed: %v \n", err)
		return
	}

	// Mongo init
	if err := mongoI.Init(settings.Conf.MongoConfig); err != nil {
		fmt.Printf("init mongo failed: %v \n", err)
		return
	}
	defer mongoI.Close()

	idx := mongo.IndexModel{
		Keys:    bson.D{{"user_id", 1}},
		Options: options.Index().SetUnique(true),
	}
	rankDD := mongoI.RankDB()
	_, err := rankDD.Collection("user").Indexes().CreateOne(context.Background(), idx, options.CreateIndexes().SetMaxTime(10*time.Second))
	if err != nil {
		log.Fatalf("创建索引失败，%s \n",err)
		return
	}

	idx = mongo.IndexModel{
		Keys:    bson.D{{"gift_id", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = rankDD.Collection("gift").Indexes().CreateOne(context.Background(), idx, options.CreateIndexes().SetMaxTime(10*time.Second))
	if err != nil {
		log.Fatalf("创建索引失败，%s \n",err)
		return
	}

	idxs := []mongo.IndexModel{
		{
			Keys:    bson.D{{"gift_id", 1}},
		},
		{
			Keys:    bson.D{{"to_user_id", 1}},
		},
	}
	_, err = rankDD.Collection("gift_log").Indexes().CreateMany(context.Background(), idxs, options.CreateIndexes().SetMaxTime(10*time.Second))



	fmt.Println("done")
}
