package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
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
	// Read body content
	defer r.Body.Close()
	// 400 Bad Request (JSON Bad request)
	var bodyData CreateURLData
	if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
		ReturnError(w, err, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verify URL
	_url := bodyData.Url
	// 400 Bad Request (invalid URL)
	_, err := url.ParseRequestURI(_url)
	if err != nil {
		ReturnError(w, err, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Create Random URL Short code
	shortInDB := true
	var shortCode string
	for shortInDB {
		shortCode = createShortCode()
		shortInDB, err = s.isShortCodeInDB(shortCode)
		// ERROR with database connection
		if err != nil {
			ReturnError(w, err, "Error checking short code in db", http.StatusInternalServerError)
			return
		}
	}

	// Store new content on db
	responseData, err := s.saveShortenURL(_url, shortCode)
	if err != nil {
		ReturnError(w, err, "Internal server error with DB", http.StatusInternalServerError)
		return
	}

	// Return response to user (STATUS OK)
	ReturnJSON(w, r, responseData, http.StatusCreated)
}

func (s *shortenServer) RetrieveURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	// Check if url-short code in server
	shortInDB, err := s.isShortCodeInDB(shortCode)
	if err != nil {
		ReturnError(w, err, "Error checking short code in db", http.StatusInternalServerError)
		return
	}
	if !shortInDB {
		ReturnError(w, err, "Error short code not in db", http.StatusBadRequest)
		return
	}

	// Update Access Counter
	err = s.updateAccessCount(shortCode)
	if err != nil {
		ReturnError(w, err, "Error updating url count from DB", http.StatusInternalServerError)
		return
	}
	// Retrieve Original URL
	responseData, err := s.retrieveOriginalURL(shortCode)
	if err != nil {
		ReturnError(w, err, "Error retrieving url from DB", http.StatusInternalServerError)
		return
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
		ReturnError(w, err, "Internal server error decoding body", http.StatusBadRequest)
		return
	}

	// Update data in DB
	responseData, err := s.updateOriginalURL(bodyData.Url, shortCode)
	if err != nil {
		ReturnError(w, err, "Error updating url from DB", http.StatusInternalServerError)
		return
	}

	// Return JSON to user data
	ReturnJSON(w, r, responseData, http.StatusOK)
}

func (s *shortenServer) DeleteURL(w http.ResponseWriter, r *http.Request) {
	// Get shorten Code
	shortCode := r.PathValue("shortCode")

	// Check if shorten Code is in DB
	shortInDB, err := s.isShortCodeInDB(shortCode)
	if err != nil {
		ReturnError(w, err, "Error accessing DB", http.StatusInternalServerError)
		return
	}
	if !shortInDB {
		ReturnError(w, nil, "Invalid short code", http.StatusNotFound)
		return
	}

	// Delete from DB
	err = s.deleteShortCode(shortCode)
	if err != nil {
		ReturnError(w, err, "Error deleting url from DB", http.StatusInternalServerError)
		return
	}
	// Return No Content if success
	w.WriteHeader(http.StatusNoContent)
}

func (s *shortenServer) GetStatsURL(w http.ResponseWriter, r *http.Request) {
	// Get Short code
	shortCode := r.PathValue("shortCode")

	// Check short code in DB
	shortInDB, err := s.isShortCodeInDB(shortCode)
	if err != nil {
		ReturnError(w, err, "Error accessing DB", http.StatusInternalServerError)
	}
	if !shortInDB {
		ReturnError(w, nil, "No shortCode in DB", http.StatusBadRequest)
	}

	// Retrieve short code stats
	responseData, err := s.retrieveOriginalURL(shortCode)
	if err != nil {
		ReturnError(w, err, "Error accessing DB", http.StatusInternalServerError)
	}
	ReturnJSON(w, r, responseData, http.StatusOK)
}

func (s *shortenServer) HandleShortCode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.RetrieveURL(w, r)
	case http.MethodPut:
		s.UpdateURL(w, r)
	case http.MethodDelete:
		s.DeleteURL(w, r)
	default:
		ReturnError(w, nil, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
