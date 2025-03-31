package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sinhaseemant/glofox-backend/internal/storage"
	"github.com/sinhaseemant/glofox-backend/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Connect to MongoDB
	mr, err := storage.NewMongoRepository()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	router := routes.NewRouter(mr)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
