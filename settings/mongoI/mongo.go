package mongoI

import (
	"boframe/settings"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var MongoClient *mongo.Client

func Init(conf *settings.MongoConfig) (err error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(conf.MaxPoolSize)))
	if err != nil {
		return err
	}

	timeout, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	err = client.Connect(timeout)
	if err != nil {
		return err
	}

	timeout, cancel = context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	err = client.Ping(timeout, readpref.Primary())
	if err != nil {
		return err
	}
	MongoClient = client
	return nil
}

func Close() error {
	return MongoClient.Disconnect(context.Background())
}

func RankDB() *mongo.Database {
	return MongoClient.Database("rank")
}

func IsErrNoDocuments(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}