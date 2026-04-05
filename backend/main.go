// Package main initializes the NQ backend by starting the GraphQL server
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Check for health check flag
	if len(os.Args) > 1 && os.Args[1] == "-health-check" {
		healthCheck()
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, continuing")
	}
	GraphQL()
}

// healthCheck performs a simple HTTP health check
func healthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		log.Printf("Health check failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed with status: %d", resp.StatusCode)
		os.Exit(1)
	}

	os.Exit(0)
}
