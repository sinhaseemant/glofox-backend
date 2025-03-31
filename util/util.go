package util

import (
	"github.com/sinhaseemant/glofox-backend/models"
)

// SendGlobalResponse -
func SendGlobalResponse(status string, data interface{}, code int, err error) models.GlobalResponse {
	response := models.GlobalResponse{}
	if err != nil {
		response.Message = err.Error()
	}
	response.StatusCode = code
	response.Status = status
	response.Data = data
	return response
}
