package models

type GlobalResponse struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
