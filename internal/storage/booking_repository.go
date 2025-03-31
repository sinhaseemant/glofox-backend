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

// BookingRepositoryInterface defines the contract for BookingRepository
type BookingRepositoryInterface interface {
	Create(ctx context.Context, booking *models.Booking) (primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]models.Booking, error)
}

// BookingRepository struct for MongoDB
type BookingRepository struct {
	Collection *mongo.Collection
}

// NewBookingRepository initializes a BookingRepository with MongoDB collection
func NewBookingRepository(db *mongo.Database) BookingRepositoryInterface {
	collection := db.Collection("bookings")
	return &BookingRepository{Collection: collection}
}

// Create inserts a new booking into the MongoDB collection
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) (primitive.ObjectID, error) {
	res, err := r.Collection.InsertOne(ctx, booking)
	if err != nil {
		log.Printf("Error inserting booking: %v", err)
		return primitive.NilObjectID, fmt.Errorf("failed to insert booking: %w", err)
	}
	log.Printf("Inserted booking with ID: %v", res.InsertedID)
	return res.InsertedID.(primitive.ObjectID), nil
}

// GetAll retrieves all bookings from the MongoDB collection
func (r *BookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding bookings: %v", err)
		return nil, fmt.Errorf("failed to find bookings: %w", err)
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	for cursor.Next(ctx) {
		var booking models.Booking
		if err := cursor.Decode(&booking); err != nil {
			log.Printf("Error decoding booking: %v", err)
			return nil, fmt.Errorf("failed to decode booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return bookings, nil
}
