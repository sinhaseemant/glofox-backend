package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sinhaseemant/glofox-backend/api"
	"github.com/sinhaseemant/glofox-backend/internal/handlers"
	"github.com/sinhaseemant/glofox-backend/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter sets up the chi router and registers routes
func NewRouter(repo *storage.MongoRepository) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger) // Logs requests
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	})) // Enable CORS
	// Serve OpenAPI JSON at /swagger.json
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %s", err)
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonSpec, err := swagger.MarshalJSON()
		if err != nil {
			http.Error(w, "Failed to serialize OpenAPI spec", http.StatusInternalServerError)
			return
		}
		w.Write(jsonSpec)
	})

	// Serve Swagger UI at /swagger/
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // Points to OpenAPI JSON spec
	))
	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Inject repository into API handlers
	// Get Database instance
	cr := storage.NewClassRepository(repo.Client.Database("classes"))
	ch := handlers.NewClassHandler(cr)
	br := storage.NewBookingRepository(repo.Client.Database("bookings"))
	bh := handlers.NewBookingHandler(br)
	si := api.NewServerInterface(repo, ch, bh)

	r.Mount("/", api.Handler(si))

	return r
}
