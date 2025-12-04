package server

import "net/http"

type Server interface {
	// CloseServer closes database connection
	CloseServer()

	// CreateURL handles posts requests from POST /shorten
	CreateURL(w http.ResponseWriter, r *http.Request)

	// RetrieveURL retrieves the original url with shortCode
	// Handles GET requests from GET /shorten/{shortCode}
	RetrieveURL(w http.ResponseWriter, r *http.Request)

	// UpdateURL updates the original url given a shortCode
	// Handles PUT requests from PUT /shorten/{shortCode}
	UpdateURL(w http.ResponseWriter, r *http.Request)

	// DeleteURL deletes the original url and associated shortCode
	// Handles DELETE requests from DELETE /shorten/{shortCode}
	DeleteURL(w http.ResponseWriter, r *http.Request)

	// Handles Requests from endpoint /shorten/{shortCode}
	HandleShortCode(w http.ResponseWriter, r *http.Request)

	// Handles requests for getting URL stats
	GetStatsURL(w http.ResponseWriter, r *http.Request)
}
