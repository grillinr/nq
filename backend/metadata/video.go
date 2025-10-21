package metadata

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	tmdb "github.com/cyruzin/golang-tmdb"
)

// VideoFetcher implements the Fetcher interface for movies and TV shows
type VideoFetcher struct {
	client *tmdb.Client
}

// NewVideoFetcher creates a new VideoFetcher
func NewVideoFetcher() (*VideoFetcher, error) {
	token := os.Getenv("TMDB_API_READ_ACCESS_TOKEN")
	if token == "" {
		return nil, errors.New("TMDB_API_READ_ACCESS_TOKEN environment variable is not set")
	}

	tmdbClient, err := tmdb.InitV4(token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TMDB client: %w", err)
	}

	// Enable auto-retry for rate limiting
	tmdbClient.SetClientAutoRetry()

	return &VideoFetcher{
		client: tmdbClient,
	}, nil
}

// Fetch retrieves video metadata from TMDB
func (f *VideoFetcher) Fetch(info MediaInfo) (*MediaMetadata, error) {
	if info.ID == "" && info.Title == "" {
		return nil, errors.New("either video ID or title is required")
	}

	// Determine if we're fetching a movie or TV show
	isTV := info.Type == MediaTypeTV

	var metadata *MediaMetadata
	var err error

	if info.ID != "" {
		// Search by ID
		id, err := strconv.Atoi(info.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid video ID: %w", err)
		}

		if isTV {
			metadata, err = f.fetchTVShowByID(id)
		} else {
			metadata, err = f.fetchMovieByID(id)
		}
	} else {
		// Search by title
		if isTV {
			metadata, err = f.searchTVShow(info.Title, info.ReleaseYear)
		} else {
			metadata, err = f.searchMovie(info.Title, info.ReleaseYear)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch video data: %w", err)
	}

	return metadata, nil
}

func (f *VideoFetcher) fetchMovieByID(id int) (*MediaMetadata, error) {
	options := map[string]string{
		"append_to_response": "videos,images",
	}

	movie, err := f.client.GetMovieDetails(id, options)
	if err != nil {
		return nil, err
	}

	metadata := &MediaMetadata{
		MediaInfo: MediaInfo{
			Type:        MediaTypeMovie,
			Title:       movie.Title,
			ReleaseYear: parseYear(movie.ReleaseDate),
		},
		Description: movie.Overview,
		URL:         fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.ID),
	}

	// Set poster image if available
	if movie.PosterPath != "" {
		metadata.ImageURL = tmdb.GetImageURL(movie.PosterPath, tmdb.Original)
	}

	// Get genres
	if len(movie.Genres) > 0 {
		metadata.Genres = make([]string, len(movie.Genres))
		for i, genre := range movie.Genres {
			metadata.Genres[i] = genre.Name
		}
	}

	return metadata, nil
}

func (f *VideoFetcher) fetchTVShowByID(id int) (*MediaMetadata, error) {
	options := map[string]string{
		"append_to_response": "videos,images",
	}

	tvShow, err := f.client.GetTVDetails(id, options)
	if err != nil {
		return nil, err
	}

	metadata := &MediaMetadata{
		MediaInfo: MediaInfo{
			Type:        MediaTypeTV,
			Title:       tvShow.Name,
			ReleaseYear: parseYear(tvShow.FirstAirDate),
		},
		Description: tvShow.Overview,
		URL:         fmt.Sprintf("https://www.themoviedb.org/tv/%d", tvShow.ID),
	}

	// Set poster image if available
	if tvShow.PosterPath != "" {
		metadata.ImageURL = tmdb.GetImageURL(tvShow.PosterPath, tmdb.Original)
	}

	// Get genres
	if len(tvShow.Genres) > 0 {
		metadata.Genres = make([]string, len(tvShow.Genres))
		for i, genre := range tvShow.Genres {
			metadata.Genres[i] = genre.Name
		}
	}

	return metadata, nil
}

func (f *VideoFetcher) searchMovie(title string, year int) (*MediaMetadata, error) {
	options := map[string]string{
		"query": title,
	}

	if year > 0 {
		options["year"] = strconv.Itoa(year)
	}

	result, err := f.client.GetSearchMovies("en-US", options)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no movie found with title: %s", title)
	}

	// Return the first result
	movie := result.Results[0]
	return &MediaMetadata{
		MediaInfo: MediaInfo{
			Type:        MediaTypeMovie,
			Title:       movie.Title,
			ReleaseYear: parseYear(movie.ReleaseDate),
		},
		Description: movie.Overview,
		ImageURL:    tmdb.GetImageURL(movie.PosterPath, tmdb.Original),
		URL:         fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.ID),
	}, nil
}

func (f *VideoFetcher) searchTVShow(title string, year int) (*MediaMetadata, error) {
	options := map[string]string{
		"query": title,
	}

	if year > 0 {
		options["first_air_date_year"] = strconv.Itoa(year)
	}

	result, err := f.client.GetSearchTVShow("en-US", options)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no TV show found with title: %s", title)
	}

	// Return the first result
	tvShow := result.Results[0]
	return &MediaMetadata{
		MediaInfo: MediaInfo{
			Type:        MediaTypeTV,
			Title:       tvShow.Name,
			ReleaseYear: parseYear(tvShow.FirstAirDate),
		},
		Description: tvShow.Overview,
		ImageURL:    tmdb.GetImageURL(tvShow.PosterPath, tmdb.Original),
		URL:         fmt.Sprintf("https://www.themoviedb.org/tv/%d", tvShow.ID),
	}, nil
}

// Helper function to parse year from a date string (YYYY-MM-DD)
func parseYear(dateStr string) int {
	if dateStr == "" {
		return 0
	}

	parts := strings.Split(dateStr, "-")
	if len(parts) > 0 {
		year, _ := strconv.Atoi(parts[0])
		return year
	}
	return 0
}
