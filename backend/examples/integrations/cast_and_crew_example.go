package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grillinr/nq/db"
	"github.com/grillinr/nq/graph/model"
)

// CastAndCrewExample demonstrates creating a movie and fetching cast & crew via the repository.
func CastAndCrewExample(repo db.Repository) {
	ctx := context.Background()

	input := model.CreateMovieInput{
		Title: "Example Movie for Cast/Crew",
		Cast:  []string{"Alice Example", "Bob Example"},
		Crew:  []string{"Carol Crew", "Dave Crew"},
	}

	movie, err := repo.CreateMovie(ctx, input)
	if err != nil {
		log.Fatalf("failed to create movie: %v", err)
	}

	fmt.Printf("Created movie: %s (ID: %s)\n", movie.Title, movie.ID.String())

	// Fetch cast and crew using the new generic method
	cast, crew, castCredits, crewCredits, err := repo.GetCastAndCrew(ctx, movie.ID)
	if err != nil {
		log.Fatalf("failed to get cast and crew: %v", err)
	}

	fmt.Println("Cast:")
	for _, p := range cast {
		fmt.Printf(" - %s (ID: %s)\n", p.Name, p.ID.String())
	}

	fmt.Println("Crew:")
	for _, p := range crew {
		fmt.Printf(" - %s (ID: %s)\n", p.Name, p.ID.String())
	}

	fmt.Println("Structured Cast Credits:")
	for _, c := range castCredits {
		if c.Character != nil {
			fmt.Printf(" - %s as %s\n", c.Name, *c.Character)
		} else {
			fmt.Printf(" - %s\n", c.Name)
		}
	}

	fmt.Println("Structured Crew Credits:")
	for _, c := range crewCredits {
		if c.Job != nil {
			fmt.Printf(" - %s (%s)\n", c.Name, *c.Job)
		} else {
			fmt.Printf(" - %s\n", c.Name)
		}
	}
}
