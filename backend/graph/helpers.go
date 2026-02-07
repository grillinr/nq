package graph

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"
)

// Helper methods for recursive media search
func (r *mutationResolver) collectUniqueMovieCredits(ctx context.Context, cast, crew []*model.Person, excludeTitle string, excludeYear int) []*metadata.VideoMetadata {
	uniqueMovies := make(map[string]*metadata.VideoMetadata) // key: "title_year"

	processPersonCredits := func(personID string) {
		log.Printf("Fetching movie credits for person: %s", personID)
		metaSvc := r.Repo.GetMetadata()
		if metaSvc == nil {
			log.Printf("Metadata service not available")
			return
		}

		metadataSvc, ok := metaSvc.(*metadata.Service)
		if !ok {
			log.Printf("Failed to cast metadata service")
			return
		}

		fetchers := metadataSvc.GetFetchers()
		if fetchers == nil {
			log.Printf("No fetchers available")
			return
		}

		id, err := strconv.Atoi(personID)
		if err != nil {
			log.Printf("Invalid person ID: %s, error: %v", personID, err)
			return
		}

		videoFetcher, ok := fetchers[metadata.MediaTypeMovie].(*metadata.VideoFetcher)
		if !ok {
			log.Printf("Video fetcher not available")
			return
		}

		movies, err := videoFetcher.FetchPersonMovieCredits(id)
		if err != nil {
			log.Printf("Failed to fetch movie credits for person %s: %v", personID, err)
			return
		}

		log.Printf("Fetched %d movie credits for person %s", len(movies), personID)

		// Add to unique map, excluding original movie
		excludeKey := fmt.Sprintf("%s_%d", excludeTitle, excludeYear)
		for _, movie := range movies {
			key := fmt.Sprintf("%s_%d", movie.Title, movie.ReleaseYear)
			if key != excludeKey {
				uniqueMovies[key] = movie
			}
		}
	}

	// Process all cast and crew
	for _, person := range append(cast, crew...) {
		if person.ExternalID != nil {
			processPersonCredits(*person.ExternalID)
		}
	}

	// Convert map to slice
	var result []*metadata.VideoMetadata
	for _, movie := range uniqueMovies {
		result = append(result, movie)
	}
	log.Printf("Collected %d unique movie credits (excluding original)", len(result))
	return result
}
func (r *mutationResolver) processMovieBatch(ctx context.Context, movies []*metadata.VideoMetadata, searchDepth int32, maxConnections int, sourceID uuid.UUID) {
	// Limit connections if needed
	if len(movies) > maxConnections {
		log.Printf("Limiting to %d connections (had %d)", maxConnections, len(movies))
		movies = movies[:maxConnections]
	}

	for _, m := range movies {
		log.Printf("Processing movie: %s (%d)", m.Title, m.ReleaseYear)
		// Check if already exists
		existing, err := r.Repo.FindMediaByTitleTypeYear(ctx, m.Title, string(m.Type), &m.ReleaseYear)
		if err == nil && existing != nil {
			log.Printf("Movie %s already exists with depth %d", m.Title, existing.GetSearchDepth())
			// If existing has higher depth, update to lower
			if existing.GetSearchDepth() > searchDepth {
				err = r.Repo.UpdateMediaSearchDepth(ctx, existing.GetID(), searchDepth)
				if err != nil {
					log.Printf("Failed to update search depth for %s: %v", m.Title, err)
				} else {
					log.Printf("Updated search depth for %s to %d", m.Title, searchDepth)
				}
			}
			if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, existing.GetID()); linkErr != nil {
				log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
			}
			continue // Already exists
		}

		// Create the movie
		yearStr := strconv.Itoa(m.ReleaseYear)
		input := model.CreateMovieInput{
			Title:       m.Title,
			ReleaseDate: &yearStr,
			Description: &m.Description,
			CoverURL:    &m.ImageURL,
			SearchDepth: &searchDepth,
		}
		created, err := r.Repo.CreateMovie(ctx, input)
		if err != nil {
			log.Printf("Failed to create movie %s: %v", m.Title, err)
		} else {
			log.Printf("Created movie: %s", m.Title)
			if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
				log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
			}
		}
	}
}
func (r *mutationResolver) recursiveSearchMovies(ctx context.Context, movie *model.Movie, maxConnections int) {
	log.Printf("Starting recursive search for movie: %s (ID: %s)", movie.Title, movie.ID)

	// Parse release year from ReleaseDate string
	excludeYear := 0
	if movie.ReleaseDate != nil {
		if year, err := strconv.Atoi(*movie.ReleaseDate); err == nil {
			excludeYear = year
		}
	}

	// Collect all unique connected movies
	uniqueMovies := r.collectUniqueMovieCredits(ctx, movie.Cast, movie.Crew, movie.Title, excludeYear)

	// Process batch
	r.processMovieBatch(ctx, uniqueMovies, 1, maxConnections, movie.ID)

	log.Printf("Completed recursive search for movie: %s", movie.Title)
}
func (r *mutationResolver) collectUniqueTVShowCredits(ctx context.Context, cast, crew []*model.Person, excludeTitle string, excludeYear int) []*metadata.VideoMetadata {
	uniqueTVShows := make(map[string]*metadata.VideoMetadata) // key: "title_year"

	processPersonCredits := func(personID string) {
		log.Printf("Fetching TV show credits for person: %s", personID)
		metaSvc := r.Repo.GetMetadata()
		if metaSvc == nil {
			log.Printf("Metadata service not available")
			return
		}

		metadataSvc, ok := metaSvc.(*metadata.Service)
		if !ok {
			log.Printf("Failed to cast metadata service")
			return
		}

		fetchers := metadataSvc.GetFetchers()
		if fetchers == nil {
			log.Printf("No fetchers available")
			return
		}

		id, err := strconv.Atoi(personID)
		if err != nil {
			log.Printf("Invalid person ID: %s, error: %v", personID, err)
			return
		}

		videoFetcher, ok := fetchers[metadata.MediaTypeMovie].(*metadata.VideoFetcher)
		if !ok {
			log.Printf("Video fetcher not available")
			return
		}

		tvShows, err := videoFetcher.FetchPersonTVShowCredits(id)
		if err != nil {
			log.Printf("Failed to fetch TV show credits for person %s: %v", personID, err)
			return
		}

		log.Printf("Fetched %d TV show credits for person %s", len(tvShows), personID)

		// Add to unique map, excluding original TV show
		excludeKey := fmt.Sprintf("%s_%d", excludeTitle, excludeYear)
		for _, tvShow := range tvShows {
			key := fmt.Sprintf("%s_%d", tvShow.Title, tvShow.ReleaseYear)
			if key != excludeKey {
				uniqueTVShows[key] = tvShow
			}
		}
	}

	// Process all cast and crew
	for _, person := range append(cast, crew...) {
		if person.ExternalID != nil {
			processPersonCredits(*person.ExternalID)
		}
	}

	// Convert map to slice
	var result []*metadata.VideoMetadata
	for _, tvShow := range uniqueTVShows {
		result = append(result, tvShow)
	}
	log.Printf("Collected %d unique TV show credits (excluding original)", len(result))
	return result
}
func (r *mutationResolver) processTVShowBatch(ctx context.Context, tvShows []*metadata.VideoMetadata, searchDepth int32, maxConnections int, sourceID uuid.UUID) {
	// Limit connections if needed
	if len(tvShows) > maxConnections {
		log.Printf("Limiting to %d connections (had %d)", maxConnections, len(tvShows))
		tvShows = tvShows[:maxConnections]
	}

	for _, m := range tvShows {
		log.Printf("Processing TV show: %s (%d)", m.Title, m.ReleaseYear)
		// Check if already exists
		existing, err := r.Repo.FindMediaByTitleTypeYear(ctx, m.Title, string(m.Type), &m.ReleaseYear)
		if err == nil && existing != nil {
			log.Printf("TV show %s already exists with depth %d", m.Title, existing.GetSearchDepth())
			// If existing has higher depth, update to lower
			if existing.GetSearchDepth() > searchDepth {
				err = r.Repo.UpdateMediaSearchDepth(ctx, existing.GetID(), searchDepth)
				if err != nil {
					log.Printf("Failed to update search depth for %s: %v", m.Title, err)
				} else {
					log.Printf("Updated search depth for %s to %d", m.Title, searchDepth)
				}
			}
			if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, existing.GetID()); linkErr != nil {
				log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
			}
			continue // Already exists
		}

		// Create the TV show
		yearStr := strconv.Itoa(m.ReleaseYear)
		input := model.CreateTVShowInput{
			Title:       m.Title,
			ReleaseDate: &yearStr,
			Description: &m.Description,
			CoverURL:    &m.ImageURL,
			SearchDepth: &searchDepth,
		}
		created, err := r.Repo.CreateTVShow(ctx, input)
		if err != nil {
			log.Printf("Failed to create TV show %s: %v", m.Title, err)
		} else {
			log.Printf("Created TV show: %s", m.Title)
			if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
				log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
			}
		}
	}
}
func (r *mutationResolver) recursiveSearchTVShows(ctx context.Context, tvShow *model.TVShow, maxConnections int) {
	log.Printf("Starting recursive search for TV show: %s (ID: %s)", tvShow.Title, tvShow.ID)

	// Parse release year from ReleaseDate string
	excludeYear := 0
	if tvShow.ReleaseDate != nil {
		if year, err := strconv.Atoi(*tvShow.ReleaseDate); err == nil {
			excludeYear = year
		}
	}

	// Collect all unique connected TV shows
	uniqueTVShows := r.collectUniqueTVShowCredits(ctx, tvShow.Cast, tvShow.Crew, tvShow.Title, excludeYear)

	// Process batch
	r.processTVShowBatch(ctx, uniqueTVShows, 1, maxConnections, tvShow.ID)

	log.Printf("Completed recursive search for TV show: %s", tvShow.Title)
}
