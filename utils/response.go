package utils

import (
	"encoding/json"
	"net/http"

	"personalnote.eu/simple-go-api/models"
)

// SendJSONResponse sends a JSON response with the given status code
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendErrorResponse sends an error response
func SendErrorResponse(w http.ResponseWriter, statusCode int, error string, message string) {
	errorResp := models.ErrorResponse{
		Error:   error,
		Message: message,
	}
	SendJSONResponse(w, statusCode, errorResp)
}

// SendSuccessResponse sends a success response
func SendSuccessResponse(w http.ResponseWriter, message string) {
	resp := models.Response{
		Message: message,
	}
	SendJSONResponse(w, http.StatusOK, resp)
}
