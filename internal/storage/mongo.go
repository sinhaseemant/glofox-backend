package storage

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository handles MongoDB operations
type MongoRepository struct {
	Client *mongo.Client
}

// NewMongoRepository initializes MongoDB connection
func NewMongoRepository() (*MongoRepository, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return &MongoRepository{Client: client}, nil
}

// GetCollection returns a MongoDB collection
func (m *MongoRepository) GetCollection(collectionName string) *mongo.Collection {
	return m.Client.Database("goflox").Collection(collectionName)
}
