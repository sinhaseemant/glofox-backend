package api

import (
	"net/http"

	"github.com/sinhaseemant/glofox-backend/internal/handlers"
	"github.com/sinhaseemant/glofox-backend/internal/storage"
)

// NewServerInterface creates and returns a new instance of ServerInterface.
// It initializes and returns a pointer to the serverInterface struct, which
// implements the ServerInterface interface that is present in api.gen.go file.

func NewServerInterface(repo *storage.MongoRepository, classHandler handlers.ClassHandlerInterface, bookingHandler handlers.BookingHandlerInterface) ServerInterface {
	return &serverInterface{repo: repo, ch: classHandler, bh: bookingHandler}
}

type serverInterface struct {
	repo *storage.MongoRepository
	ch   handlers.ClassHandlerInterface
	bh   handlers.BookingHandlerInterface
}

func (s *serverInterface) BookClass(w http.ResponseWriter, r *http.Request) {
	s.bh.BookClassHandler(w, r)
}

func (s *serverInterface) CreateClass(w http.ResponseWriter, r *http.Request) {
	s.ch.CreateClassHandler(w, r)
}

func (s *serverInterface) GetClasses(w http.ResponseWriter, r *http.Request) {
	s.ch.GetClassesHandler(w, r)
}

func (s *serverInterface) GetBookings(w http.ResponseWriter, r *http.Request) {
	s.bh.GetBookingsHandler(w, r)
}
