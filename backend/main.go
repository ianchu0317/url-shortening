package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"url-shortening/server"

	"github.com/joho/godotenv"
)

// CORS Middleware
func corsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables from .env if present
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found or couldn't load it; falling back to environment variables")
	}

	// Get .env variables values
	shortCodeLen, err := strconv.Atoi(os.Getenv("SHORT_CODE_LEN"))
	if err != nil {
		log.Fatalf("Error loading SHORT_CODE_LEN")
	}
	databaseURL := os.Getenv("DATABASE_URL")

	// Create Server
	app, err := server.CreateServer(databaseURL, shortCodeLen)
	if err != nil {
		log.Fatalf("DB connection error, %v", err)
	}
	defer app.CloseServer()

	// Set up ctrl+c shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server (go routine)
	go func() {
		http.HandleFunc("/shorten", app.CreateURL)
		http.HandleFunc("/{shortCode}", app.HandleShortCode)
		http.HandleFunc("/{shortCode}/stats", app.GetStatsURL)

		log.Println("Server starting on :8080")
		if err := http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux)); err != nil {
			log.Fatal("Failed starting server")
		}
	}()

	<-c
	log.Println("Shutting down server...")

}
