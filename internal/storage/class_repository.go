package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/sinhaseemant/glofox-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClassRepositoryInterface defines the contract for ClassRepository
type ClassRepositoryInterface interface {
	Create(ctx context.Context, class *models.Class) (primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]models.Class, error)
}

// ClassRepository struct for MongoDB
type ClassRepository struct {
	Collection *mongo.Collection
}

// NewClassRepository initializes a ClassRepository with MongoDB collection
func NewClassRepository(db *mongo.Database) ClassRepositoryInterface {
	collection := db.Collection("classes")
	return &ClassRepository{Collection: collection}
}

// Create inserts a new class into the MongoDB collection
func (r *ClassRepository) Create(ctx context.Context, class *models.Class) (primitive.ObjectID, error) {
	res, err := r.Collection.InsertOne(ctx, class)
	if err != nil {
		log.Printf("Error inserting class: %v", err)
		return primitive.NilObjectID, fmt.Errorf("failed to insert class: %w", err)
	}
	log.Printf("Inserted class with ID: %v", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID), nil
}

// GetAll retrieves all classes from the MongoDB collection
func (r *ClassRepository) GetAll(ctx context.Context) ([]models.Class, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding classes: %v", err)
		return nil, fmt.Errorf("failed to find classes: %w", err)
	}
	defer cursor.Close(ctx)

	var classes []models.Class
	for cursor.Next(ctx) {
		var class models.Class
		if err := cursor.Decode(&class); err != nil {
			log.Printf("Error decoding class: %v", err)
			return nil, fmt.Errorf("failed to decode class: %w", err)
		}
		classes = append(classes, class)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return classes, nil
}
