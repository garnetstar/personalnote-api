package utils

import (
	"net/http"
)

// ValidateHTTPMethod validates that the request uses the specified HTTP method
func ValidateHTTPMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		SendErrorResponse(w, http.StatusMethodNotAllowed,
			"Method not allowed", "Only "+allowedMethod+" requests are accepted")
		return false
	}
	return true
}
