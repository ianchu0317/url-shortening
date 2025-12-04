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
	// Read body content
	defer r.Body.Close()
	var bodyData CreateURLData
	if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
		ReturnError(w, err, "Bad Request", http.StatusBadRequest)
		return
	}
	// Get URL from Body and get shorten url (if its unique)
	url := bodyData.Url
	urlInDB, err := s.isUrlInDB(url)
	if err != nil {
		ReturnError(w, err, "Internal server error checking url in db", http.StatusInternalServerError)
		return
	}
	if urlInDB {
		ReturnError(w, err, "URL already in DB", http.StatusConflict)
		return
	}
	shortCode := createShortCode(url)
	// Store new content on db
	responseData, err := s.saveShortenURL(url, shortCode)
	if err != nil {
		ReturnError(w, err, "Internal server error with DB", http.StatusInternalServerError)
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
		ReturnError(w, err, "Error checking short code in db", http.StatusInternalServerError)
		return
	}
	if !shortInDB {
		ReturnError(w, err, "Error short code not in db", http.StatusBadRequest)
		return
	}

	// Update Access Counter / Get original url
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
