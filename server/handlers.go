package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Server struct and creation

type shortenServer struct {
	DB *pgxpool.Pool
}

// CreateServer(url) takes a database url and starts connection and initial configuration
// if success, will return the server, if not err != nil
func CreateServer(databaseURL string) (Server, error) {
	server := &shortenServer{}
	// Configuration for psql pool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	config.MaxConnLifetime = 3 * time.Minute
	config.MaxConns = 10
	// Create DB Pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	server.DB = pool
	return server, err
}

func (s *shortenServer) CloseServer() {
	s.DB.Close()
}

// Server Handlers

func (s *shortenServer) CreateURL(w http.ResponseWriter, r *http.Request) {
	// Check method for this endpoint
	if r.Method != "POST" {
		http.Error(w, "Only POST Method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body content
	defer r.Body.Close()
	var bodyData CreateURLData
	if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
		log.Fatalf("Internal server error decoding body, %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get URL from Body and get shorten url (if its unique)
	url := bodyData.Url
	urlInDB, err := s.isUrlInDB(url)
	if err != nil {
		log.Fatalf("Internal server error checking url in db, %v", err)
		http.Error(w, "Internal server error checking url in db", http.StatusInternalServerError)
		return
	}
	if urlInDB {
		http.Error(w, "URL already in DB", http.StatusConflict)
		return
	}
	shortCode := createShortCode(url)

	// Store new content on db
	responseData, err := s.saveShortenURL(url, shortCode)
	if err != nil {
		log.Fatalf("Internal server error saving to db, %v", err)
		http.Error(w, "Internal server error saving to db", http.StatusInternalServerError)
		return
	}
	// Return response to user
	ReturnJSON(w, r, responseData, http.StatusCreated)
}

func (s *shortenServer) RetrieveURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	// Check if url-short code in server
	shortInDB, err := s.isShortCodeInDB(shortCode)
	if err != nil {
		log.Fatalf("Error checking short code in db, %v", err)
		http.Error(w, "Error checking short code in db", http.StatusInternalServerError)
		return
	}
	if !shortInDB {
		http.Error(w, "Error short code not in db", http.StatusBadRequest)
		return
	}

	// Update Access Counter / Get original url
	responseData, err := s.retrieveOriginalURL(shortCode)
	if err != nil {
		log.Fatalf("Error retrieving url from DB, %v", err)
		http.Error(w, "Error retrieving url from DB", http.StatusInternalServerError)
	}

	// Return original url -> Redirect
	http.Redirect(w, r, responseData.URL, http.StatusMovedPermanently)
}

func (s *shortenServer) UpdateURL(w http.ResponseWriter, r *http.Request) {
	// Get short code
	shortCode := r.PathValue("shortCode")

	// Read body data
	defer r.Body.Close()
	var bodyData CreateURLData
	if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
		log.Fatalf("Internal server error decoding body, %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Update data in DB
	responseData, err := s.updateOriginalURL(bodyData.Url, shortCode)
	if err != nil {
		log.Fatalf("Error updating url from DB, %v", err)
		http.Error(w, "Error updating url from DB", http.StatusInternalServerError)
	}

	// Return JSON to user data
	ReturnJSON(w, r, responseData, http.StatusOK)
}

func (s *shortenServer) HandleShortCode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.RetrieveURL(w, r)
	case http.MethodPut:
		s.UpdateURL(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Helpers

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
