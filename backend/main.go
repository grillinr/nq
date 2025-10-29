// Package main initializes the NQ backend by starting the GraphQL server
package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, continuing")
	}
	GraphQL()
}
