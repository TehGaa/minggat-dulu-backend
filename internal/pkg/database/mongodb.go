package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

func GetMongoClient(uri string) *mongo.Client {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().
			ApplyURI(uri).
			SetMaxPoolSize(50).
			SetMinPoolSize(1)

		client, err := mongo.Connect(clientOptions)
		if err != nil {
			log.Fatalf("Failed to create MongoDB client: %v", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Println("Successfully connected to MongoDB")
		clientInstance = client

	})
	return clientInstance
}
