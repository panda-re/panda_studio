package db

import (
	"context"

	config "github.com/panda-re/panda_studio/internal/configuration"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ObjectID primitive.ObjectID

var mongoDbClient *mongo.Client = nil

func getMongoConnection(ctx context.Context) (*mongo.Client, error) {
	if mongoDbClient != nil {
		return mongoDbClient, nil
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetConfig().Mongo.Uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func GetMongoDatabase(ctx context.Context) (*mongo.Database, error) {
	conn, err := getMongoConnection(ctx)
	if err != nil {
		return nil, err
	}
	
	return conn.Database(config.GetConfig().Mongo.Database), nil
}