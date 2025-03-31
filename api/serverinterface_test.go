package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sinhaseemant/glofox-backend/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClassHandler is a mock implementation of ClassHandlerInterface.
type MockClassHandler struct {
	mock.Mock
}

func (m *MockClassHandler) CreateClassHandler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockClassHandler) GetClassesHandler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// MockBookingHandler is a mock implementation of BookingHandlerInterface.
type MockBookingHandler struct {
	mock.Mock
}

func (m *MockBookingHandler) BookClassHandler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockBookingHandler) GetBookingsHandler(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestNewServerInterface(t *testing.T) {
	mockRepo := &storage.MongoRepository{}
	mockClassHandler := new(MockClassHandler)
	mockBookingHandler := new(MockBookingHandler)

	server := NewServerInterface(mockRepo, mockClassHandler, mockBookingHandler)
	assert.NotNil(t, server, "NewServerInterface should return a non-nil instance")
}

func TestServerInterfaceMethods(t *testing.T) {
	mockRepo := &storage.MongoRepository{}
	mockClassHandler := new(MockClassHandler)
	mockBookingHandler := new(MockBookingHandler)

	server := NewServerInterface(mockRepo, mockClassHandler, mockBookingHandler)

	tests := []struct {
		name           string
		method         func(s ServerInterface, w http.ResponseWriter, r *http.Request)
		mockHandler    *mock.Mock
		expectedMethod string
	}{
		{
			name: "BookClass calls BookClassHandler",
			method: func(s ServerInterface, w http.ResponseWriter, r *http.Request) {
				s.BookClass(w, r)
			},
			mockHandler:    &mockBookingHandler.Mock,
			expectedMethod: "BookClassHandler",
		},
		{
			name: "CreateClass calls CreateClassHandler",
			method: func(s ServerInterface, w http.ResponseWriter, r *http.Request) {
				s.CreateClass(w, r)
			},
			mockHandler:    &mockClassHandler.Mock,
			expectedMethod: "CreateClassHandler",
		},
		{
			name: "GetClasses calls GetClassesHandler",
			method: func(s ServerInterface, w http.ResponseWriter, r *http.Request) {
				s.GetClasses(w, r)
			},
			mockHandler:    &mockClassHandler.Mock,
			expectedMethod: "GetClassesHandler",
		},
		{
			name: "GetBookings calls GetBookingsHandler",
			method: func(s ServerInterface, w http.ResponseWriter, r *http.Request) {
				s.GetBookings(w, r)
			},
			mockHandler:    &mockBookingHandler.Mock,
			expectedMethod: "GetBookingsHandler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			tt.mockHandler.On(tt.expectedMethod, rec, req).Return()

			tt.method(server, rec, req)

			tt.mockHandler.AssertCalled(t, tt.expectedMethod, rec, req)
		})
	}
}
