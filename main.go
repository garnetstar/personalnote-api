package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var counter int

func halloHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Counter", strconv.Itoa(counter))
	resp := Response{Message: "Hallo, from Go!"}
	log.Printf("Received request #%d", counter)
	log.Printf("Query params: %v", r.URL.Query())
	log.Printf("ur: %v", r.RequestURI)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResp := ErrorResponse{
			Error:   "Method not allowed",
			Message: "Only POST requests are accepted",
		}
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorResp := ErrorResponse{
			Error:   "Invalid request body",
			Message: "Could not read request body",
		}
		json.NewEncoder(w).Encode(errorResp)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorResp := ErrorResponse{
			Error:   "Invalid JSON",
			Message: "Could not parse JSON request body",
		}
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// Validate the request
	var validationErrors []string
	if strings.TrimSpace(user.Name) == "" {
		validationErrors = append(validationErrors, "name is required")
	}
	if user.ID <= 0 {
		validationErrors = append(validationErrors, "id must be a positive integer")
	}

	if len(validationErrors) > 0 {
		log.Printf("Validation failed: %v", validationErrors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorResp := ErrorResponse{
			Error:   "Validation failed",
			Message: fmt.Sprintf("Validation errors: %s", strings.Join(validationErrors, ", ")),
		}
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// Log the valid request to console
	log.Printf("✅ Received valid user data: Name=%s, ID=%d", user.Name, user.ID)
	log.Printf("📝 Full request body: %s", string(body))

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := Response{
		Message: fmt.Sprintf("User %s with ID %d has been processed successfully", user.Name, user.ID),
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", halloHandler)
	http.HandleFunc("/user", userHandler)
	addr := ":8080"
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

}
