package util

import (
	"errors"
	"testing"

	"github.com/sinhaseemant/glofox-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestSendGlobalResponse(t *testing.T) {
	t.Run("should return response with error message when error is provided", func(t *testing.T) {
		err := errors.New("test error")
		status := "failure"
		code := 400

		expectedResponse := models.GlobalResponse{
			Message:    "test error",
			StatusCode: 400,
			Status:     "failure",
			Data:       nil,
		}

		response := SendGlobalResponse(status, nil, code, err)

		assert.Equal(t, expectedResponse, response)
	})

	t.Run("should return response without error message when error is nil", func(t *testing.T) {
		status := "success"
		data := "test data"
		code := 200

		expectedResponse := models.GlobalResponse{
			Message:    "",
			StatusCode: 200,
			Status:     "success",
			Data:       "test data",
		}

		response := SendGlobalResponse(status, data, code, nil)

		assert.Equal(t, expectedResponse, response)
	})
}
