package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestUnimplementedServer(t *testing.T) {
	server := Unimplemented{}

	t.Run("GetBookings", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bookings", nil)
		rec := httptest.NewRecorder()

		server.GetBookings(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})

	t.Run("BookClass", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/bookings", nil)
		rec := httptest.NewRecorder()

		server.BookClass(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})

	t.Run("GetClasses", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/classes", nil)
		rec := httptest.NewRecorder()

		server.GetClasses(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})

	t.Run("CreateClass", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/classes", nil)
		rec := httptest.NewRecorder()

		server.CreateClass(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})
}

func TestHandler(t *testing.T) {
	t.Run("HandlerWithDefaultOptions", func(t *testing.T) {
		server := Unimplemented{}
		handler := Handler(&server)

		req, _ := http.NewRequest("GET", "/bookings", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})

	t.Run("HandlerWithCustomOptions", func(t *testing.T) {
		server := Unimplemented{}
		router := chi.NewRouter()
		handler := HandlerWithOptions(&server, ChiServerOptions{
			BaseRouter: router,
		})

		req, _ := http.NewRequest("GET", "/classes", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	})
}

func TestDecodeSpec(t *testing.T) {
	t.Run("DecodeSpecSuccess", func(t *testing.T) {
		data, err := decodeSpec()

		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("DecodeSpecCached", func(t *testing.T) {
		cachedFunc := decodeSpecCached()
		data, err := cachedFunc()

		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestGetSwagger(t *testing.T) {
	t.Run("GetSwaggerSuccess", func(t *testing.T) {
		swagger, err := GetSwagger()

		assert.NoError(t, err)
		assert.NotNil(t, swagger)
	})
}
