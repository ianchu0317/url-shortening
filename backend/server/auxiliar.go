package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	CHARSET = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// createShortenURL takes a string and returns a shorten version of it (hash)
func createShortCode(shortCodeLen int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, shortCodeLen)
	for i := range b {
		b[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}
	return string(b)
}

// Helpers - Auxiliar functions

// ReturnJSON take a request and respond with JSON
func ReturnJSON(w http.ResponseWriter, r *http.Request, responseData *ResponseCreatedURLData, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Fatalf("Internal server error converting to JSON, %v", err)
		http.Error(w, "Internal server error converting to JSON", http.StatusInternalServerError)
		return
	}
}

// ReturnError takes an error, message and status code and returns error to user / and log
func ReturnError(w http.ResponseWriter, err error, message string, statusCode int) {
	log.Printf("%s: %v", message, err)

	//http.Error(w, message, statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		log.Fatalf("Internal server error converting to JSON, %v", err)
		http.Error(w, "Internal server error converting to JSON", http.StatusInternalServerError)
		return
	}
}
