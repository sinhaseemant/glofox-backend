package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/sinhaseemant/glofox-backend/internal/storage"
	"github.com/sinhaseemant/glofox-backend/models"
	"github.com/sinhaseemant/glofox-backend/util"
)

// ClassHandlerInterface defines the contract for ClassHandler
type ClassHandlerInterface interface {
	CreateClassHandler(w http.ResponseWriter, r *http.Request)
	GetClassesHandler(w http.ResponseWriter, r *http.Request)
}

// ClassHandler struct for dependency injection
type ClassHandler struct {
	Repo storage.ClassRepositoryInterface
}

// NewClassHandler initializes a handler with DI
func NewClassHandler(repo storage.ClassRepositoryInterface) ClassHandlerInterface {
	return &ClassHandler{Repo: repo}
}

// CreateClassHandler handles class creation
func (h *ClassHandler) CreateClassHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Class
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusBadRequest, err))
		http.Error(w, string(resStr), http.StatusBadRequest)
		return
	}

	// Validate the class data
	var validationErrors []string

	if req.Name == "" {
		validationErrors = append(validationErrors, "Class name is required")
	}
	if req.StartDate.IsZero() {
		validationErrors = append(validationErrors, "Start date is required")
	}
	if req.EndDate.IsZero() {
		validationErrors = append(validationErrors, "End date is required")
	}
	if req.Capacity <= 0 {
		validationErrors = append(validationErrors, "Capacity must be greater than 0")
	}

	// If there are validation errors, return them in the response
	if len(validationErrors) > 0 {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, validationErrors, http.StatusBadRequest, errors.New("Validation failed")))
		http.Error(w, string(resStr), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	// Insert class into MongoDB
	id, err := h.Repo.Create(ctx, &req)
	if err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusInternalServerError, err))
		http.Error(w, string(resStr), http.StatusInternalServerError)
		return
	}
	req.ID = id

	log.Info().Msgf("Class created: %+v", req)
	resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusSuccess, req, http.StatusCreated, nil))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resStr)
}

func (h *ClassHandler) GetClassesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Fetch classes from MongoDB
	classes, err := h.Repo.GetAll(ctx)
	if err != nil {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusInternalServerError, err))
		http.Error(w, string(resStr), http.StatusInternalServerError)
		return
	}

	// If no classes are found, return a 404 response
	if len(classes) == 0 {
		resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusFailed, nil, http.StatusNotFound, errors.New("No classes found")))
		http.Error(w, string(resStr), http.StatusNotFound)
		return
	}

	// Return the list of classes
	w.Header().Set("Content-Type", "application/json")
	resStr, _ := json.Marshal(util.SendGlobalResponse(util.StatusSuccess, classes, http.StatusOK, nil))
	w.WriteHeader(http.StatusOK)
	w.Write(resStr)
}
