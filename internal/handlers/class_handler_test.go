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

// MockClassRepository is a mock implementation of the ClassRepositoryInterface.
type MockClassRepository struct {
	mock.Mock
}

func (m *MockClassRepository) Create(ctx context.Context, class *models.Class) (primitive.ObjectID, error) {
	args := m.Called(ctx, class)
	return primitive.NewObjectID(), args.Error(1)
}

func (m *MockClassRepository) GetAll(ctx context.Context) ([]models.Class, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Class), args.Error(1)
}

func TestCreateClassHandler(t *testing.T) {
	mockRepo := new(MockClassRepository)
	handler := NewClassHandler(mockRepo)

	mockStartDate := models.CustomDate(time.Now())
	mockEndDate := models.CustomDate(time.Now().Add(24 * time.Hour))

	tests := []struct {
		name           string
		requestBody    models.Class
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid class creation",
			requestBody: models.Class{
				Name:      "Yoga Class",
				StartDate: mockStartDate,
				EndDate:   mockEndDate,
				Capacity:  10,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Missing class name",
			requestBody: models.Class{
				StartDate: mockStartDate,
				EndDate:   mockEndDate,
				Capacity:  10,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError == nil {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return("classID", nil)
			} else {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return("", tt.mockError)
			}

			reqBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/create-class", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateClassHandler(rec, req)

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

func TestGetClassesHandler(t *testing.T) {
	mockRepo := new(MockClassRepository)
	handler := NewClassHandler(mockRepo)

	mockStartDate := models.CustomDate(time.Now())
	mockEndDate := models.CustomDate(time.Now().Add(24 * time.Hour))
	tests := []struct {
		name           string
		mockClasses    []models.Class
		mockError      error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "No classes found",
			mockClasses:    []models.Class{},
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedError:  "No classes found",
		},
		{
			name: "Classes retrieved successfully",
			mockClasses: []models.Class{
				{
					ID:        primitive.NewObjectID(),
					Name:      "Yoga Class",
					StartDate: mockStartDate,
					EndDate:   mockEndDate,
					Capacity:  10,
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
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockClasses, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/get-classes", nil)
			rec := httptest.NewRecorder()

			handler.GetClassesHandler(rec, req)

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
