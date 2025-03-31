package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sinhaseemant/glofox-backend/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewRouter(t *testing.T) {
	// Mock MongoRepository
	mockClient := &mongo.Client{} // Replace with a mock or stub if needed
	mockRepo := &storage.MongoRepository{Client: mockClient}

	// Create the router
	router := NewRouter(mockRepo)

	// Define test cases
	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{
			name:       "Health Check",
			method:     http.MethodGet,
			path:       "/health",
			statusCode: http.StatusOK,
		},
		{
			name:       "Swagger JSON",
			method:     http.MethodGet,
			path:       "/swagger.json",
			statusCode: http.StatusOK,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, rr.Code)
			}
		})
	}
}
