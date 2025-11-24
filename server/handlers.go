package server

import (
	"context"
	"encoding/json"
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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get URL from Body and get shorten url (if its unique)
	url := bodyData.Url
	urlInDB, err := s.isUrlInDB(url)
	if err != nil {
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
		http.Error(w, "Internal server error saving to db", http.StatusInternalServerError)
		return
	}
	// Return response to user
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Internal server error converting to JSON", http.StatusInternalServerError)
		return
	}
}
