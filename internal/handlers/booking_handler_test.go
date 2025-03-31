package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sinhaseemant/glofox-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBookingRepository is a mock implementation of the BookingRepositoryInterface.
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(ctx context.Context, booking *models.Booking) (primitive.ObjectID, error) {
	args := m.Called(ctx, booking)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockBookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func TestBookClassHandler(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	handler := NewBookingHandler(mockRepo)
	mockStartDate := models.CustomDate(time.Now())

	tests := []struct {
		name           string
		requestBody    models.Booking
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid booking creation",
			requestBody: models.Booking{
				ClassID:    primitive.NewObjectID(),
				ClassName:  "Yoga Class",
				MemberName: "John Doe",
				Date:       mockStartDate,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Missing class ID",
			requestBody: models.Booking{
				ClassName:  "Yoga Class",
				MemberName: "John Doe",
				Date:       mockStartDate,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError == nil {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
			} else {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(primitive.NilObjectID, tt.mockError)
			}

			reqBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/book-class", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.BookClassHandler(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response models.GlobalResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Message, tt.expectedError)
			}
		})
	}
}

func TestGetBookingsHandler(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	handler := NewBookingHandler(mockRepo)
	mockStartDate := models.CustomDate(time.Now())
	tests := []struct {
		name           string
		mockBookings   []models.Booking
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "No bookings found",
			mockBookings:   []models.Booking{},
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedError:  "No bookings found",
		},
		{
			name: "Bookings retrieved successfully",
			mockBookings: []models.Booking{
				{
					ID:         primitive.NewObjectID(),
					ClassID:    primitive.NewObjectID(),
					ClassName:  "Yoga Class",
					MemberName: "John Doe",
					Date:       mockStartDate,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		//clear the mock expectations
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil

		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockBookings, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/get-bookings", nil)
			rec := httptest.NewRecorder()

			handler.GetBookingsHandler(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var response models.GlobalResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Message, tt.expectedError)
			}
		})
	}
}
