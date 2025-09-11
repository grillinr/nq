package main

import (
	"nq/integrations"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	integrations.FetchSpotifyAuthToken()
}
