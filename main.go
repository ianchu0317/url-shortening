package main

import (
	"log"
	"net/http"
	"os"
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

	// Start server
	http.HandleFunc("/shorten", app.CreateURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
