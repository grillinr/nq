// Package main initializes the NQ backend by starting the GraphQL server
package main

import (
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	GraphQL()
}
