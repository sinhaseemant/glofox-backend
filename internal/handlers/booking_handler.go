package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sinhaseemant/glofox-backend/internal/storage"
	"github.com/sinhaseemant/glofox-backend/models"
	"github.com/sinhaseemant/glofox-backend/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookingHandlerInterface defines the contract for BookingHandler
type BookingHandlerInterface interface {
	BookClassHandler(w http.ResponseWriter, r *http.Request)
	GetBookingsHandler(w http.ResponseWriter, r *http.Request)
}

type BookingHandler struct {
	Repo storage.BookingRepositoryInterface
}

// NewClassHandler initializes a handler with DI
func NewBookingHandler(repo storage.BookingRepositoryInterface) BookingHandlerInterface {
	return &BookingHandler{Repo: repo}
}

// BookClassHandler handles class bookings
func (h *BookingHandler) BookClassHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Booking
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusBadRequest, err))
		http.Error(w, string(resStr), http.StatusBadRequest)
		return
	}

	// Validate the booking data
	var validationErrors []string = []string{}
	if req.ClassID == primitive.NilObjectID {
		validationErrors = append(validationErrors, "Class ID is required")

	}
	if req.MemberName == "" {
		validationErrors = append(validationErrors, "Member name is required")
	}

	if req.Date.IsZero() {
		validationErrors = append(validationErrors, "Booking date is required")
	}
	if req.ClassName == "" {
		validationErrors = append(validationErrors, "Class name is required")
	}
	if len(validationErrors) > 0 {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, validationErrors, http.StatusBadRequest, errors.New("Validation failed")))
		http.Error(w, string(resStr), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Insert booking into MongoDB
	id, err := h.Repo.Create(ctx, &req)
	if err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusInternalServerError, err))
		http.Error(w, string(resStr), http.StatusInternalServerError)
		return
	}
	req.ID = id

	log.Info().Msgf("Booking created: %v", req)
	resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusSuccess, req, http.StatusCreated, nil))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resStr)
}

// GetBookingsHandler retrieves all bookings
func (h *BookingHandler) GetBookingsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	bookings, err := h.Repo.GetAll(ctx)
	if err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusInternalServerError, err))
		http.Error(w, string(resStr), http.StatusInternalServerError)
		return
	}

	if len(bookings) == 0 {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusNotFound, errors.New("No bookings found")))
		http.Error(w, string(resStr), http.StatusNotFound)
		return
	}

	resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusSuccess, bookings, http.StatusOK, nil))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resStr)

}
