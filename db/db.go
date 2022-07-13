package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg *MongoInstance

func ConnectDb(mongoUri string, dbName string) (*MongoInstance, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	db := client.Database(dbName)
	if err != nil {
		return nil, err
	}
	Mg = &MongoInstance{Client: client, Db: db}
	return Mg, nil
}
