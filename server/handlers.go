package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Server struct and creation

type shortenServer struct {
	mu sync.Mutex
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

	config.MaxConnLifetime = 3
	config.MaxConns = 10

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
	// Lock connection resources
	s.mu.Lock()
	defer s.mu.Unlock()

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

	// Get URL from Body and get shorten url
	url := bodyData.Url
	shortenURL := createShortenURL(url)

	// Store new content on db
	responseData, err := s.saveShortenURL(url, shortenURL)
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
