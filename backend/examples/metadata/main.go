package main

import (
	"fmt"
	"log"

	"github.com/grillinr/nq/metadata"
)

func main() {
	// Example 1: Fetch movie metadata by title
	fmt.Println("Example 1: Fetch movie metadata by title")
	movieMetadata, err := metadata.GetMetadataByTitle("movie", "The Matrix", 1999)
	if err != nil {
		log.Printf("Error fetching movie metadata: %v", err)
	} else {
		fmt.Printf("Movie: %s (%d)\n", movieMetadata.Title, movieMetadata.Year)
		fmt.Printf("Description: %s\n", movieMetadata.Description)
		fmt.Printf("Genres: %v\n", movieMetadata.Genres)
		fmt.Printf("Image: %s\n\n", movieMetadata.ImageURL)
	}

	// Example 2: Fetch book metadata by ID (ISBN)
	fmt.Println("Example 2: Fetch book metadata by ISBN")
	bookMetadata, err := metadata.GetMetadataByID("book", "9780439139601") // Harry Potter and the Goblet of Fire
	if err != nil {
		log.Printf("Error fetching book metadata: %v", err)
	} else {
		fmt.Printf("Book: %s\n", bookMetadata.Title)
		fmt.Printf("Description: %s\n", bookMetadata.Description)
		fmt.Printf("URL: %s\n\n", bookMetadata.URL)
	}

	// Example 3: Fetch TV show metadata by title
	fmt.Println("Example 3: Fetch TV show metadata by title")
	tvMetadata, err := metadata.GetMetadataByTitle("tv", "Breaking Bad", 0)
	if err != nil {
		log.Printf("Error fetching TV metadata: %v", err)
	} else {
		fmt.Printf("TV Show: %s (%d)\n", tvMetadata.Title, tvMetadata.Year)
		fmt.Printf("Description: %s\n", tvMetadata.Description)
		fmt.Printf("Genres: %v\n\n", tvMetadata.Genres)
	}

	// Example 4: Using the map-based function (more flexible)
	fmt.Println("Example 4: Using map-based metadata fetching")
	mediaInfo := map[string]interface{}{
		"title": "The Legend of Zelda: Breath of the Wild",
		"year":  2017,
	}
	gameMetadata, err := metadata.GetMetadata("game", mediaInfo)
	if err != nil {
		log.Printf("Error fetching game metadata: %v", err)
	} else {
		fmt.Printf("Game: %s (%d)\n", gameMetadata.Title, gameMetadata.Year)
		fmt.Printf("Description: %s\n", gameMetadata.Description)
		fmt.Printf("Genres: %v\n", gameMetadata.Genres)
	}

	// Example 5: Error handling - invalid media type
	fmt.Println("\nExample 5: Error handling")
	_, err = metadata.GetMetadataByTitle("invalid_type", "Some Title", 2023)
	if err != nil {
		fmt.Printf("Expected error for invalid media type: %v\n", err)
	}

	// Example 6: Error handling - missing required data
	_, err = metadata.GetMetadataByTitle("movie", "", 2023)
	if err != nil {
		fmt.Printf("Expected error for empty title: %v\n", err)
	}
}
