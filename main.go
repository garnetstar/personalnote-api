package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
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

func main() {
	http.HandleFunc("/", halloHandler)
	addr := ":8080"
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

}
