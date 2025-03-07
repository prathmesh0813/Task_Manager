package dao

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var Chats *mongo.Database

func ConnectMongoDB() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err!=nil{
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	err= client.Connect(ctx)
	if err!=nil{
		return err
	}

	MongoClient = client
	Chats = client.Database("chats")
	return nil
}
