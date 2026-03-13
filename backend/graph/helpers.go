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

var nonAlnumRE = regexp.MustCompile(`[^a-z0-9]+`)

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

	name = nonAlnumRE.ReplaceAllString(name, " ")
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
	case metadata.MediaTypeGame:
		return "Game", true
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
		case metadata.MediaTypeGame:
			description := m.Description
			coverURL := m.ImageURL
			input := model.CreateGameInput{
				Title:        m.Title,
				ReleaseDate:  &yearStr,
				Description:  &description,
				CoverURL:     &coverURL,
				SearchDepth:  &searchDepth,
				Genre:        m.Genres,
				Themes:       m.Themes,
				Keywords:     m.Keywords,
				GameModes:    m.GameModes,
				Perspectives: m.Perspectives,
				Franchises:   m.Franchises,
				Platforms:    m.Platforms,
			}
			created, err := r.Repo.CreateGame(ctx, input)
			if err != nil {
				log.Printf("Failed to create game %s: %v", m.Title, err)
			} else {
				log.Printf("Created game: %s", m.Title)
				if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
					log.Printf("Failed to link related media for %s: %v", m.Title, linkErr)
				}
			}
		default:
			log.Printf("Skipping unsupported media type for creation: %s", m.Type)
		}
	}
}

// recursiveSearchVideo is a helper that extracts common logic for recursive video searches
func (r *mutationResolver) recursiveSearchVideo(ctx context.Context, title string, id uuid.UUID, releaseDate *string, cast, crew []*model.Person, maxConnections int, mediaType string) {
	log.Printf("Starting recursive search for %s: %s (ID: %s)", mediaType, title, id)

	// Parse release year from ReleaseDate string
	excludeYear := 0
	if releaseDate != nil {
		if year, err := strconv.Atoi(*releaseDate); err == nil {
			excludeYear = year
		}
	}

	// Collect all unique connected media (movies + TV shows)
	uniqueMedia := r.collectUniqueRelatedVideoCredits(ctx, cast, crew, title, excludeYear)

	// Process batch
	r.processMediaBatch(ctx, uniqueMedia, 1, maxConnections, id)

	log.Printf("Completed recursive search for %s: %s", mediaType, title)
}

func collectGenreNames(genres []*model.Genre) []string {
	if len(genres) == 0 {
		return nil
	}
	names := make([]string, 0, len(genres))
	for _, g := range genres {
		if g != nil && g.Name != "" {
			names = append(names, g.Name)
		}
	}
	return names
}

func (r *mutationResolver) recursiveSearchMovies(ctx context.Context, movie *model.Movie, maxConnections int) {
	r.recursiveSearchVideo(ctx, movie.Title, movie.ID, movie.ReleaseDate, movie.Cast, movie.Crew, maxConnections, "movie")

	// Genre-based cross-media linking
	normalizedGenres := normalizeTagNamesFromStrings(collectGenreNames(movie.Genres))
	if len(normalizedGenres) > 0 {
		linked, err := r.Repo.LinkRelatedMediaByTagNames(ctx, movie.ID, normalizedGenres, maxConnections)
		if err != nil {
			log.Printf("Failed to link movie %s by genre tags: %v", movie.Title, err)
		} else {
			log.Printf("Linked %d items by genre tags for movie %s", linked, movie.Title)
		}
	}

	// Title-based cross-media linking (e.g. book ↔ movie adaptation)
	if linked, err := r.Repo.LinkMediaByNormalizedTitle(ctx, movie.ID); err != nil {
		log.Printf("Failed to link movie %s by title: %v", movie.Title, err)
	} else {
		log.Printf("Linked %d items by title for movie %s", linked, movie.Title)
	}
}

func (r *mutationResolver) recursiveSearchTVShows(ctx context.Context, tvShow *model.TVShow, maxConnections int) {
	r.recursiveSearchVideo(ctx, tvShow.Title, tvShow.ID, tvShow.ReleaseDate, tvShow.Cast, tvShow.Crew, maxConnections, "TV show")

	// Genre-based cross-media linking
	normalizedGenres := normalizeTagNamesFromStrings(collectGenreNames(tvShow.Genres))
	if len(normalizedGenres) > 0 {
		linked, err := r.Repo.LinkRelatedMediaByTagNames(ctx, tvShow.ID, normalizedGenres, maxConnections)
		if err != nil {
			log.Printf("Failed to link TV show %s by genre tags: %v", tvShow.Title, err)
		} else {
			log.Printf("Linked %d items by genre tags for TV show %s", linked, tvShow.Title)
		}
	}

	// Title-based cross-media linking
	if linked, err := r.Repo.LinkMediaByNormalizedTitle(ctx, tvShow.ID); err != nil {
		log.Printf("Failed to link TV show %s by title: %v", tvShow.Title, err)
	} else {
		log.Printf("Linked %d items by title for TV show %s", linked, tvShow.Title)
	}
}

func (r *mutationResolver) recursiveSearchGames(ctx context.Context, game *model.Game, maxConnections int) {
	log.Printf("Starting recursive search for game: %s (ID: %s)", game.Title, game.ID)

	excludeYear := 0
	if game.ReleaseDate != nil {
		if year, err := strconv.Atoi(*game.ReleaseDate); err == nil {
			excludeYear = year
		}
	}

	uniqueGames := r.collectUniqueRelatedGameCredits(ctx, game.Title, excludeYear)
	r.processGameBatch(ctx, uniqueGames, 1, maxConnections, game.ID)

	normalizedTags := normalizeTagNamesFromStrings(collectGameTags(game))
	if len(normalizedTags) > 0 {
		linked, err := r.Repo.LinkRelatedMediaByTagNames(ctx, game.ID, normalizedTags, maxConnections)
		if err != nil {
			log.Printf("Failed to link related media by tags for game %s: %v", game.Title, err)
		} else {
			log.Printf("Linked %d related media by tags for game %s", linked, game.Title)
		}
	}

	log.Printf("Completed recursive search for game: %s", game.Title)

	// Title-based cross-media linking (e.g. game ↔ movie adaptation)
	if linked, err := r.Repo.LinkMediaByNormalizedTitle(ctx, game.ID); err != nil {
		log.Printf("Failed to link game %s by title: %v", game.Title, err)
	} else {
		log.Printf("Linked %d items by title for game %s", linked, game.Title)
	}
}

func collectGameTags(game *model.Game) []string {
	if game == nil {
		return nil
	}
	var tags []string
	tags = append(tags, game.Genre...)
	tags = append(tags, game.Themes...)
	tags = append(tags, game.Keywords...)
	tags = append(tags, game.GameModes...)
	tags = append(tags, game.Perspectives...)
	tags = append(tags, game.Franchises...)
	tags = append(tags, game.PlatformsList...)
	return tags
}

func (r *mutationResolver) collectUniqueRelatedGameCredits(ctx context.Context, title string, excludeYear int) []*metadata.MediaMetadata {
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

	gameFetcher, ok := fetchers[metadata.MediaTypeGame].(*metadata.GameFetcher)
	if !ok {
		log.Printf("Game fetcher not available")
		return nil
	}

	items, err := gameFetcher.SearchRelatedGames(title)
	if err != nil {
		log.Printf("Failed to fetch related games for title %s: %v", title, err)
		return nil
	}
	unique := make(map[string]*metadata.MediaMetadata)
	for _, item := range items {
		if item == nil || item.Title == "" {
			continue
		}
		if strings.EqualFold(item.Title, title) {
			continue
		}
		key := fmt.Sprintf("%s_%d", item.Title, item.ReleaseYear)
		unique[key] = item
	}
	var result []*metadata.MediaMetadata
	for _, item := range unique {
		result = append(result, item)
	}
	log.Printf("Collected %d unique related games", len(result))
	return result
}

func (r *mutationResolver) processGameBatch(ctx context.Context, games []*metadata.MediaMetadata, searchDepth int32, maxConnections int, sourceID uuid.UUID) {
	if len(games) > maxConnections {
		log.Printf("Limiting to %d connections (had %d)", maxConnections, len(games))
		games = games[:maxConnections]
	}

	for _, g := range games {
		if g == nil {
			continue
		}
		log.Printf("Processing game: %s (%d)", g.Title, g.ReleaseYear)

		yearStr := ""
		if g.ReleaseYear > 0 {
			yearStr = strconv.Itoa(g.ReleaseYear)
		}
		description := g.Description
		coverURL := g.ImageURL
		input := model.CreateGameInput{
			Title:        g.Title,
			ReleaseDate:  nil,
			Description:  &description,
			CoverURL:     &coverURL,
			SearchDepth:  &searchDepth,
			Genre:        g.Genres,
			Themes:       g.Themes,
			Keywords:     g.Keywords,
			GameModes:    g.GameModes,
			Perspectives: g.Perspectives,
			Franchises:   g.Franchises,
			Platforms:    g.Platforms,
		}
		if yearStr != "" {
			input.ReleaseDate = &yearStr
		}
		created, err := r.Repo.CreateGame(ctx, input)
		if err != nil {
			log.Printf("Failed to create game %s: %v", g.Title, err)
			continue
		}
		if linkErr := r.Repo.LinkRelatedMedia(ctx, sourceID, created.ID); linkErr != nil {
			log.Printf("Failed to link related media for game %s: %v", g.Title, linkErr)
		}
		normalizedTags := normalizeTagNamesFromStrings(collectMetadataTags(g))
		if len(normalizedTags) > 0 {
			if _, err := r.Repo.LinkRelatedMediaByTagNames(ctx, created.ID, normalizedTags, maxConnections); err != nil {
				log.Printf("Failed to link related media by tags for game %s: %v", g.Title, err)
			}
		}
	}
}

func collectMetadataTags(gameMeta *metadata.MediaMetadata) []string {
	if gameMeta == nil {
		return nil
	}
	var tags []string
	tags = append(tags, gameMeta.Genres...)
	tags = append(tags, gameMeta.Themes...)
	tags = append(tags, gameMeta.Keywords...)
	tags = append(tags, gameMeta.GameModes...)
	tags = append(tags, gameMeta.Perspectives...)
	tags = append(tags, gameMeta.Franchises...)
	tags = append(tags, gameMeta.Platforms...)
	return tags
}

func (r *mutationResolver) recursiveSearchBooks(ctx context.Context, book *model.Book, maxConnections int) {
	log.Printf("Starting recursive search for book: %s (ID: %s)", book.Title, book.ID)

	uniqueBooks := r.collectUniqueRelatedBookCredits(ctx, book.Authors, book.Title)
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

	// Title-based cross-media linking (e.g. book ↔ movie adaptation)
	if linked, err := r.Repo.LinkMediaByNormalizedTitle(ctx, book.ID); err != nil {
		log.Printf("Failed to link book %s by title: %v", book.Title, err)
	} else {
		log.Printf("Linked %d items by title for book %s", linked, book.Title)
	}
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

func (r *mutationResolver) collectUniqueRelatedBookCredits(ctx context.Context, authors []*model.Creator, excludeTitle string) []*metadata.BookMetadata {
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
		meta, err := bookFetcher.SearchBookByAuthorAndTitle(author.Name, "")
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

func (r *Resolver) getMyActivityForMedia(ctx context.Context, mediaID uuid.UUID) (*model.UserActivity, error) {
	// Get authenticated user from context
	currentUser, err := CurrentUser(ctx)
	if err != nil {
		// Not authenticated - return nil (not an error, just no activity)
		return nil, nil
	}

	// Get user's activity for this media
	activity, err := r.Repo.GetUserActivityForMedia(ctx, currentUser.ID, mediaID)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
