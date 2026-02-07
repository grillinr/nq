package graph

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"
	"golang.org/x/text/unicode/norm"
)

var genericTagStoplist = map[string]struct{}{
	"fiction":     {},
	"novel":       {},
	"books":       {},
	"literature":  {},
	"general":     {},
	"story":       {},
	"stories":     {},
	"nonfiction":  {},
	"non fiction": {},
	"classic":     {},
	"classics":    {},
	"paperback":   {},
	"hardcover":   {},
	"audiobook":   {},
	"audiobooks":  {},
	"young adult": {},
	"children":    {},
	"childrens":   {},
	"juvenile":    {},
	"english":     {},
}

func normalizedTagName(name string) string {
	if name == "" {
		return ""
	}
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return ""
	}

	t := norm.NFD.String(name)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	name = b.String()

	re := regexp.MustCompile(`[^a-z0-9]+`)
	name = re.ReplaceAllString(name, " ")
	name = strings.Join(strings.Fields(name), " ")
	return name
}

func isGenericTag(name string) bool {
	if name == "" {
		return true
	}
	_, ok := genericTagStoplist[normalizedTagName(name)]
	return ok
}

// Helper methods for recursive media search
func (r *mutationResolver) collectUniqueVideoCredits(ctx context.Context, cast, crew []*model.Person, excludeTitle string, excludeYear int, mediaType metadata.MediaType) []*metadata.VideoMetadata {
	uniqueMedia := make(map[string]*metadata.VideoMetadata) // key: "type_title_year"

	processPersonCredits := func(personID string) {
		log.Printf("Fetching %s credits for person: %s", mediaType, personID)
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

		videoFetcher, ok := fetchers[mediaType].(*metadata.VideoFetcher)
		if !ok {
			log.Printf("Video fetcher not available for %s", mediaType)
			return
		}

		var credits []*metadata.VideoMetadata
		switch mediaType {
		case metadata.MediaTypeMovie:
			credits, err = videoFetcher.FetchPersonMovieCredits(id)
		case metadata.MediaTypeTV:
			credits, err = videoFetcher.FetchPersonTVShowCredits(id)
		default:
			log.Printf("Unsupported media type for credits: %s", mediaType)
			return
		}
		if err != nil {
			log.Printf("Failed to fetch %s credits for person %s: %v", mediaType, personID, err)
			return
		}

		log.Printf("Fetched %d %s credits for person %s", len(credits), mediaType, personID)

		// Add to unique map, excluding original media
		excludeKey := fmt.Sprintf("%s_%s_%d", mediaType, excludeTitle, excludeYear)
		for _, item := range credits {
			key := fmt.Sprintf("%s_%s_%d", item.Type, item.Title, item.ReleaseYear)
			if key != excludeKey {
				uniqueMedia[key] = item
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
	for _, item := range uniqueMedia {
		result = append(result, item)
	}
	log.Printf("Collected %d unique %s credits (excluding original)", len(result), mediaType)
	return result
}

func (r *mutationResolver) collectUniqueRelatedVideoCredits(ctx context.Context, cast, crew []*model.Person, excludeTitle string, excludeYear int) []*metadata.VideoMetadata {
	uniqueMedia := make(map[string]*metadata.VideoMetadata)
	addItems := func(items []*metadata.VideoMetadata) {
		for _, item := range items {
			key := fmt.Sprintf("%s_%s_%d", item.Type, item.Title, item.ReleaseYear)
			uniqueMedia[key] = item
		}
	}

	addItems(r.collectUniqueVideoCredits(ctx, cast, crew, excludeTitle, excludeYear, metadata.MediaTypeMovie))
	addItems(r.collectUniqueVideoCredits(ctx, cast, crew, excludeTitle, excludeYear, metadata.MediaTypeTV))

	var result []*metadata.VideoMetadata
	for _, item := range uniqueMedia {
		result = append(result, item)
	}
	log.Printf("Collected %d unique related media credits", len(result))
	return result
}

func mediaLabelFromType(mediaType metadata.MediaType) (string, bool) {
	switch mediaType {
	case metadata.MediaTypeMovie:
		return "Movie", true
	case metadata.MediaTypeTV:
		return "TVShow", true
	case metadata.MediaTypeBook:
		return "Book", true
	default:
		return "", false
	}
}

func (r *mutationResolver) processMediaBatch(ctx context.Context, items []*metadata.VideoMetadata, searchDepth int32, maxConnections int, sourceID uuid.UUID) {
	// Limit connections if needed
	if len(items) > maxConnections {
		log.Printf("Limiting to %d connections (had %d)", maxConnections, len(items))
		items = items[:maxConnections]
	}

	for _, m := range items {
		log.Printf("Processing media: %s (%d, %s)", m.Title, m.ReleaseYear, m.Type)
		label, ok := mediaLabelFromType(m.Type)
		if !ok {
			log.Printf("Skipping unsupported media type: %s", m.Type)
			continue
		}

		// Check if already exists
		existing, err := r.Repo.FindMediaByTitleTypeYear(ctx, m.Title, label, &m.ReleaseYear)
		if err == nil && existing != nil {
			log.Printf("Media %s already exists with depth %d", m.Title, existing.GetSearchDepth())
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

		yearStr := strconv.Itoa(m.ReleaseYear)
		switch m.Type {
		case metadata.MediaTypeMovie:
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
		case metadata.MediaTypeTV:
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
		case metadata.MediaTypeBook:
			input := model.CreateBookInput{
				Title:       m.Title,
				ReleaseDate: &yearStr,
				Description: &m.Description,
				CoverURL:    &m.ImageURL,
				SearchDepth: &searchDepth,
			}
			created, err := r.Repo.CreateBook(ctx, input)
			if err != nil {
				log.Printf("Failed to create book %s: %v", m.Title, err)
			} else {
				log.Printf("Created book: %s", m.Title)
				if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
					log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
				}
			}
		default:
			log.Printf("Skipping unsupported media type for creation: %s", m.Type)
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

	// Collect all unique connected media (movies + TV shows)
	uniqueMedia := r.collectUniqueRelatedVideoCredits(ctx, movie.Cast, movie.Crew, movie.Title, excludeYear)

	// Process batch
	r.processMediaBatch(ctx, uniqueMedia, 1, maxConnections, movie.ID)

	log.Printf("Completed recursive search for movie: %s", movie.Title)
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

	// Collect all unique connected media (TV shows + movies)
	uniqueMedia := r.collectUniqueRelatedVideoCredits(ctx, tvShow.Cast, tvShow.Crew, tvShow.Title, excludeYear)

	// Process batch
	r.processMediaBatch(ctx, uniqueMedia, 1, maxConnections, tvShow.ID)

	log.Printf("Completed recursive search for TV show: %s", tvShow.Title)
}

func (r *mutationResolver) recursiveSearchBooks(ctx context.Context, book *model.Book, maxConnections int) {
	log.Printf("Starting recursive search for book: %s (ID: %s)", book.Title, book.ID)

	excludeYear := 0
	if book.ReleaseDate != nil {
		if year, err := strconv.Atoi(*book.ReleaseDate); err == nil {
			excludeYear = year
		}
	}

	uniqueBooks := r.collectUniqueRelatedBookCredits(ctx, book.Authors, book.Title, excludeYear)
	r.processBookBatch(ctx, uniqueBooks, 1, maxConnections, book.ID)

	// Cross-media linking via shared tags (subjects and genres)
	normalizedTags := collectNormalizedTags(book.Subjects)
	if len(normalizedTags) > 0 {
		linked, err := r.Repo.LinkRelatedMediaByTagNames(ctx, book.ID, normalizedTags, maxConnections)
		if err != nil {
			log.Printf("Failed to link related media by tags for book %s: %v", book.Title, err)
		} else {
			log.Printf("Linked %d related media by tags for book %s", linked, book.Title)
		}
	}

	log.Printf("Completed recursive search for book: %s", book.Title)
}

func collectNormalizedTags(tags []*model.Tag) []string {
	if len(tags) == 0 {
		return nil
	}
	unique := make(map[string]struct{})
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		name := normalizedTagName(tag.Name)
		if name == "" || isGenericTag(name) {
			continue
		}
		unique[name] = struct{}{}
	}
	if len(unique) == 0 {
		return nil
	}
	result := make([]string, 0, len(unique))
	for name := range unique {
		result = append(result, name)
	}
	return result
}

func (r *mutationResolver) collectUniqueRelatedBookCredits(ctx context.Context, authors []*model.Creator, excludeTitle string, excludeYear int) []*metadata.BookMetadata {
	uniqueBooks := make(map[string]*metadata.BookMetadata)

	metaSvc := r.Repo.GetMetadata()
	if metaSvc == nil {
		log.Printf("Metadata service not available")
		return nil
	}

	metadataSvc, ok := metaSvc.(*metadata.Service)
	if !ok {
		log.Printf("Failed to cast metadata service")
		return nil
	}

	fetchers := metadataSvc.GetFetchers()
	if fetchers == nil {
		log.Printf("No fetchers available")
		return nil
	}

	bookFetcher, ok := fetchers[metadata.MediaTypeBook].(*metadata.BookFetcher)
	if !ok {
		log.Printf("Book fetcher not available")
		return nil
	}

	for _, author := range authors {
		if author == nil || author.Name == "" {
			continue
		}
		meta, err := bookFetcher.SearchBookByAuthorAndTitle(author.Name, "", excludeYear)
		if err != nil {
			log.Printf("Failed to fetch books for author %s: %v", author.Name, err)
			continue
		}
		for _, book := range meta {
			if book == nil {
				continue
			}
			if strings.EqualFold(book.Title, excludeTitle) {
				continue
			}
			key := fmt.Sprintf("%s_%d", book.Title, book.ReleaseYear)
			uniqueBooks[key] = book
		}
	}

	var result []*metadata.BookMetadata
	for _, book := range uniqueBooks {
		result = append(result, book)
	}
	log.Printf("Collected %d unique related book credits", len(result))
	return result
}

func (r *mutationResolver) processBookBatch(ctx context.Context, books []*metadata.BookMetadata, searchDepth int32, maxConnections int, sourceID uuid.UUID) {
	if len(books) > maxConnections {
		log.Printf("Limiting to %d connections (had %d)", maxConnections, len(books))
		books = books[:maxConnections]
	}

	for _, b := range books {
		if b == nil {
			continue
		}
		log.Printf("Processing book: %s (%d)", b.Title, b.ReleaseYear)

		yearStr := ""
		if b.ReleaseYear > 0 {
			yearStr = strconv.Itoa(b.ReleaseYear)
		}
		description := b.Description
		coverURL := b.ImageURL
		input := model.CreateBookInput{
			Title:       b.Title,
			ReleaseDate: nil,
			Description: &description,
			CoverURL:    &coverURL,
			SearchDepth: &searchDepth,
			Authors:     b.Authors,
			Publishers:  b.Publishers,
			Subjects:    b.Subjects,
		}
		if yearStr != "" {
			input.ReleaseDate = &yearStr
		}
		if b.Pages > 0 {
			pages := int32(b.Pages)
			input.Pages = &pages
		}
		if b.Publisher != "" {
			input.Publisher = &b.Publisher
		}
		if b.ID != "" {
			input.Isbn = &b.ID
		}
		created, err := r.Repo.CreateBook(ctx, input)
		if err != nil {
			log.Printf("Failed to create book %s: %v", b.Title, err)
			continue
		}
		if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
			log.Printf("Failed to link related media for book %s: %v", b.Title, linkErr)
		}
		normalizedTags := normalizeTagNamesFromStrings(b.Subjects)
		if len(normalizedTags) > 0 {
			if _, err := r.Repo.LinkRelatedMediaByTagNames(ctx, created.ID, normalizedTags, maxConnections); err != nil {
				log.Printf("Failed to link related media by tags for book %s: %v", b.Title, err)
			}
		}
	}
}

func normalizeTagNamesFromStrings(names []string) []string {
	if len(names) == 0 {
		return nil
	}
	unique := make(map[string]struct{})
	for _, name := range names {
		normalized := normalizedTagName(name)
		if normalized == "" || isGenericTag(normalized) {
			continue
		}
		unique[normalized] = struct{}{}
	}
	if len(unique) == 0 {
		return nil
	}
	result := make([]string, 0, len(unique))
	for name := range unique {
		result = append(result, name)
	}
	return result
}
