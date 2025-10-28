package main

import (
	"fmt"
	"log"

	"github.com/grillinr/nq/metadata"
)

func main() {
	// Example 1: Fetch movie metadata by title
	fmt.Println("Example 1: Fetch movie metadata by title")
	movieMetadataIface, err := metadata.GetMetadataByTitle("movie", "The Matrix", 1999, "en")
	if err != nil {
		log.Printf("Error fetching movie metadata: %v", err)
	} else if movieMetadata, ok := movieMetadataIface.(*metadata.VideoMetadata); !ok {
		log.Printf("Unexpected metadata type for movie: %T", movieMetadataIface)
	} else {
		fmt.Printf("Movie: %s (%d)\n", movieMetadata.Title, movieMetadata.ReleaseYear)
		fmt.Printf("Description: %s\n", movieMetadata.Description)
		fmt.Printf("Genres: %v\n", movieMetadata.Genres)
		fmt.Printf("Image: %s\n\n", movieMetadata.ImageURL)
	}

	// Example 2: Fetch book metadata by ID (ISBN)
	fmt.Println("Example 2: Fetch book metadata by ISBN")
	bookMetadataIface, err := metadata.GetMetadataByID("book", "9780439139601", "en") // Harry Potter and the Goblet of Fire
	if err != nil {
		log.Printf("Error fetching book metadata: %v", err)
	} else if bookMetadata, ok := bookMetadataIface.(*metadata.BookMetadata); !ok {
		log.Printf("Unexpected metadata type for book: %T", bookMetadataIface)
	} else {
		fmt.Printf("Book: %s\n", bookMetadata.Title)
		fmt.Printf("Description: %s\n", bookMetadata.Description)
		fmt.Printf("URL: %s\n\n", bookMetadata.URL)
	}

	// Example 3: Fetch TV show metadata by title
	fmt.Println("Example 3: Fetch TV show metadata by title")
	tvMetadataIface, err := metadata.GetMetadataByTitle("tv", "Breaking Bad", 0, "en")
	if err != nil {
		log.Printf("Error fetching TV metadata: %v", err)
	} else if tvMetadata, ok := tvMetadataIface.(*metadata.VideoMetadata); !ok {
		log.Printf("Unexpected metadata type for tv: %T", tvMetadataIface)
	} else {
		fmt.Printf("TV Show: %s (%d)\n", tvMetadata.Title, tvMetadata.ReleaseYear)
		fmt.Printf("Description: %s\n", tvMetadata.Description)
		fmt.Printf("Genres: %v\n\n", tvMetadata.Genres)
	}

	// Example 4: Using the map-based function (more flexible)
	fmt.Println("Example 4: Using map-based metadata fetching")
	mediaInfo := map[string]interface{}{
		"title": "The Legend of Zelda: Breath of the Wild",
		"year":  2017,
	}
	gameMetadataIface, err := metadata.GetMetadata("game", mediaInfo, "en")
	if err != nil {
		log.Printf("Error fetching game metadata: %v", err)
	} else if gameMetadata, ok := gameMetadataIface.(*metadata.VideoMetadata); !ok {
		log.Printf("Unexpected metadata type for game: %T", gameMetadataIface)
	} else {
		fmt.Printf("Game: %s (%d)\n", gameMetadata.Title, gameMetadata.ReleaseYear)
		fmt.Printf("Description: %s\n", gameMetadata.Description)
		fmt.Printf("Genres: %v\n", gameMetadata.Genres)
	}

	// Example 5: Error handling - invalid media type
	fmt.Println("\nExample 5: Error handling")
	_, err = metadata.GetMetadataByTitle("invalid_type", "Some Title", 2023, "en")
	if err != nil {
		fmt.Printf("Expected error for invalid media type: %v\n", err)
	}

	// Example 6: Error handling - missing required data
	_, err = metadata.GetMetadataByTitle("movie", "", 2023, "en")
	if err != nil {
		fmt.Printf("Expected error for empty title: %v\n", err)
	}

	// Example 7: Using language parameter for international content
	fmt.Println("\nExample 7: Using language parameter")
	movieMetadataSpanishIface, err := metadata.GetMetadataByTitle("movie", "The Matrix", 1999, "es")
	if err != nil {
		log.Printf("Error fetching movie metadata in Spanish: %v", err)
	} else if movieMetadataSpanish, ok := movieMetadataSpanishIface.(*metadata.VideoMetadata); !ok {
		log.Printf("Unexpected metadata type for movie (spanish): %T", movieMetadataSpanishIface)
	} else {
		fmt.Printf("Movie (Spanish): %s (%d)\n", movieMetadataSpanish.Title, movieMetadataSpanish.ReleaseYear)
		fmt.Printf("Description: %s\n", movieMetadataSpanish.Description)
	}

	// Example 8: Using map-based function with language
	fmt.Println("\nExample 8: Using map-based function with language")
	mediaInfoWithLanguage := map[string]interface{}{
		"title": "The Legend of Zelda: Breath of the Wild",
		"year":  2017,
	}
	gameMetadataFrenchIface, err := metadata.GetMetadata("game", mediaInfoWithLanguage, "fr")
	if err != nil {
		log.Printf("Error fetching game metadata in French: %v", err)
	} else if gameMetadataFrench, ok := gameMetadataFrenchIface.(*metadata.VideoMetadata); !ok {
		log.Printf("Unexpected metadata type for game (french): %T", gameMetadataFrenchIface)
	} else {
		fmt.Printf("Game (French): %s (%d)\n", gameMetadataFrench.Title, gameMetadataFrench.ReleaseYear)
		fmt.Printf("Description: %s\n", gameMetadataFrench.Description)
	}
}
