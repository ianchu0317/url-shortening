package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortening/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or couldn't load it; falling back to environment variables")
	}

	// Create Server
	app, err := server.CreateServer(os.Getenv("DATABASE_URL"))
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
		http.HandleFunc("/shorten/{shortCode}", app.HandleShortCode)
		http.HandleFunc("/shorten/{shortCode}/stats", app.GetStatsURL)

		log.Println("Server starting on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("Failed starting server")
		}
	}()

	<-c
	log.Println("Shutting down server...")

}
