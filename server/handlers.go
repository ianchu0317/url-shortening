package server

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Server struct and creation

type shortenServer struct {
	mu sync.Mutex
}

func CreateServer() Server {
	return &shortenServer{}
}

// Server Handlers

func (s *shortenServer) CreateURL(w http.ResponseWriter, r *http.Request) {
	// Lock connection resources
	s.mu.Lock()
	defer s.mu.Unlock()

	if r.Method != "POST" {
		http.Error(w, "Only POST Method allowed", http.StatusMethodNotAllowed)
	}

	// Read body content
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	bodyString := string(bodyBytes)

	fmt.Printf("Testing, '%s'\n", bodyString)
}
