package metadata

import (
	"errors"
	"fmt"
	"os"
	"reflect"
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
func (f *VideoFetcher) Fetch(info MediaInfo, language string) (any, error) {
	if info.ID == "" && info.Title == "" {
		return nil, errors.New("either video ID or title is required")
	}

	// Determine if fetching a movie or TV show
	isTV := info.Type == MediaTypeTV

	var metadata any
	var err error
	var id int

	if info.ID != "" {
		// Search by ID
		id, err = strconv.Atoi(info.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid video ID: %w", err)
		}

		// Fetch by ID
		if isTV {
			metadata, err = f.fetchTVShowByID(id)
		} else {
			metadata, err = f.fetchMovieByID(id)
		}
	} else {
		// Search by title
		if isTV {
			metadata, err = f.searchTVShow(info.Title, info.ReleaseYear, language)
		} else {
			metadata, err = f.searchMovie(info.Title, info.ReleaseYear, language)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch video data: %w", err)
	}

	return metadata, nil
}

func (f *VideoFetcher) fetchMovieByID(id int) (*VideoMetadata, error) {
	options := map[string]string{
		"append_to_response": "videos,images,credits",
	}

	movie, err := f.client.GetMovieDetails(id, options)
	if err != nil {
		return nil, err
	}

	metadata := &VideoMetadata{
		MediaMetadata: MediaMetadata{
			MediaInfo: MediaInfo{
				Type:        MediaTypeMovie,
				Title:       movie.Title,
				ReleaseYear: parseYear(movie.ReleaseDate),
				ID:          fmt.Sprintf("%d", movie.ID),
			},
			Description: movie.Overview,
			URL:         fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.ID),
		},
		Budget:    float64(movie.Budget),
		BoxOffice: float64(movie.Revenue),
		Runtime:   int(movie.Runtime),
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

	// Prefer appended credits when present, using reflection to find the
	// embedded pointer to tmdb.MovieCredits. If not available, fall back to
	// calling the credits endpoint.
	// Reflection helper: find embedded pointer field of the given element type name
	findEmbeddedPtr := func(v any, elemType string) reflect.Value {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return reflect.Value{}
		}
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			t := f.Type().String()
			if t == "*tmdb."+elemType || (f.Kind() == reflect.Ptr && f.Type().Elem().Name() == elemType) {
				return f
			}
		}
		return reflect.Value{}
	}

	// Try to use appended credits if available
	if v := findEmbeddedPtr(movie, "MovieCredits"); v.IsValid() && !v.IsNil() {
		mc := v.Interface().(*tmdb.MovieCredits)
		if mc != nil {
			if len(mc.Cast) > 0 {
				metadata.Cast = make([]string, len(mc.Cast))
				metadata.CastCredits = make([]PersonCredit, len(mc.Cast))
				for i, c := range mc.Cast {
					metadata.Cast[i] = c.Name
					metadata.CastCredits[i] = PersonCredit{PersonID: int(c.ID), Name: c.Name, Character: c.Character, Order: int(c.Order)}
				}
			}
			if len(mc.Crew) > 0 {
				metadata.Crew = make([]string, len(mc.Crew))
				metadata.CrewCredits = make([]CrewCredit, len(mc.Crew))
				for i, c := range mc.Crew {
					metadata.Crew[i] = c.Name
					metadata.CrewCredits[i] = CrewCredit{PersonID: int(c.ID), Name: c.Name, Job: c.Job, Department: c.Department}
				}
			}
		}
	} else {
		// Fallback: call credits endpoint directly
		if credits, err := f.client.GetMovieCredits(id, nil); err == nil {
			if len(credits.Cast) > 0 {
				metadata.Cast = make([]string, len(credits.Cast))
				metadata.CastCredits = make([]PersonCredit, len(credits.Cast))
				for i, c := range credits.Cast {
					metadata.Cast[i] = c.Name
					metadata.CastCredits[i] = PersonCredit{PersonID: int(c.ID), Name: c.Name, Character: c.Character, Order: int(c.Order)}
				}
			}
			if len(credits.Crew) > 0 {
				metadata.Crew = make([]string, len(credits.Crew))
				metadata.CrewCredits = make([]CrewCredit, len(credits.Crew))
				for i, c := range credits.Crew {
					metadata.Crew[i] = c.Name
					metadata.CrewCredits[i] = CrewCredit{PersonID: int(c.ID), Name: c.Name, Job: c.Job, Department: c.Department}
				}
			}
		}
	}

	return metadata, nil
}

func (f *VideoFetcher) fetchTVShowByID(id int) (*VideoMetadata, error) {
	options := map[string]string{
		"append_to_response": "videos,images,credits",
	}

	tvShow, err := f.client.GetTVDetails(id, options)
	if err != nil {
		return nil, err
	}

	metadata := &VideoMetadata{
		MediaMetadata: MediaMetadata{
			MediaInfo: MediaInfo{
				Type:        MediaTypeTV,
				Title:       tvShow.Name,
				ReleaseYear: parseYear(tvShow.FirstAirDate),
				ID:          fmt.Sprintf("%d", tvShow.ID),
			},
			Description: tvShow.Overview,
			URL:         fmt.Sprintf("https://www.themoviedb.org/tv/%d", tvShow.ID),
		},
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

	// Prefer appended credits when present, else fallback to credits endpoint
	// Reuse reflection helper defined in movie fetcher context
	findEmbeddedPtr := func(v any, elemType string) reflect.Value {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return reflect.Value{}
		}
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			if f.Kind() == reflect.Ptr && f.Type().Elem().Name() == elemType {
				return f
			}
		}
		return reflect.Value{}
	}

	if v := findEmbeddedPtr(tvShow, "TVCredits"); v.IsValid() && !v.IsNil() {
		mc := v.Interface().(*tmdb.TVCredits)
		if mc != nil {
			if len(mc.Cast) > 0 {
				metadata.Cast = make([]string, len(mc.Cast))
				metadata.CastCredits = make([]PersonCredit, len(mc.Cast))
				for i, c := range mc.Cast {
					metadata.Cast[i] = c.Name
					metadata.CastCredits[i] = PersonCredit{PersonID: int(c.ID), Name: c.Name, Character: c.Character, Order: int(c.Order)}
				}
			}
			if len(mc.Crew) > 0 {
				metadata.Crew = make([]string, len(mc.Crew))
				metadata.CrewCredits = make([]CrewCredit, len(mc.Crew))
				for i, c := range mc.Crew {
					metadata.Crew[i] = c.Name
					metadata.CrewCredits[i] = CrewCredit{PersonID: int(c.ID), Name: c.Name, Job: c.Job, Department: c.Department}
				}
			}
		}
	} else {
		if credits, err := f.client.GetTVCredits(id, nil); err == nil {
			if len(credits.Cast) > 0 {
				metadata.Cast = make([]string, len(credits.Cast))
				metadata.CastCredits = make([]PersonCredit, len(credits.Cast))
				for i, c := range credits.Cast {
					metadata.Cast[i] = c.Name
					metadata.CastCredits[i] = PersonCredit{PersonID: int(c.ID), Name: c.Name, Character: c.Character, Order: int(c.Order)}
				}
			}
			if len(credits.Crew) > 0 {
				metadata.Crew = make([]string, len(credits.Crew))
				metadata.CrewCredits = make([]CrewCredit, len(credits.Crew))
				for i, c := range credits.Crew {
					metadata.Crew[i] = c.Name
					metadata.CrewCredits[i] = CrewCredit{PersonID: int(c.ID), Name: c.Name, Job: c.Job, Department: c.Department}
				}
			}
		}
	}

	return metadata, nil
}

// searchVideo is a helper that extracts common logic for searching movies and TV shows
func (f *VideoFetcher) searchVideo(title string, year int, language string, isTV bool) (*VideoMetadata, error) {
	options := map[string]string{
		"query": title,
	}

	// Include year if provided (different parameter names for movie vs TV)
	if year > 0 {
		if isTV {
			options["first_air_date_year"] = strconv.Itoa(year)
		} else {
			options["year"] = strconv.Itoa(year)
		}
	}

	// Use provided language or default to "en-US"
	searchLanguage := language
	if searchLanguage == "" {
		searchLanguage = "en-US"
	}

	// Call appropriate search method based on media type
	var resultID int
	if isTV {
		result, err := f.client.GetSearchTVShow(searchLanguage, options)
		if err != nil {
			return nil, err
		}
		if len(result.Results) == 0 {
			return nil, fmt.Errorf("no TV show found with title: %s", title)
		}
		resultID = int(result.Results[0].ID)
	} else {
		result, err := f.client.GetSearchMovies(searchLanguage, options)
		if err != nil {
			return nil, err
		}
		if len(result.Results) == 0 {
			return nil, fmt.Errorf("no movie found with title: %s", title)
		}
		resultID = int(result.Results[0].ID)
	}

	// Fetch full details for complete metadata
	if isTV {
		return f.fetchTVShowByID(resultID)
	}
	return f.fetchMovieByID(resultID)
}

func (f *VideoFetcher) searchMovie(title string, year int, language string) (*VideoMetadata, error) {
	return f.searchVideo(title, year, language, false)
}

func (f *VideoFetcher) searchTVShow(title string, year int, language string) (*VideoMetadata, error) {
	return f.searchVideo(title, year, language, true)
}

func (f *VideoFetcher) SearchTitles(query string, isTV bool, limit int) ([]*VideoSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("query required")
	}
	if limit <= 0 {
		limit = 10
	}

	options := map[string]string{
		"query": query,
	}

	searchLanguage := "en-US"
	if isTV {
		result, err := f.client.GetSearchTVShow(searchLanguage, options)
		if err != nil {
			return nil, err
		}
		results := make([]*VideoSearchResult, 0, len(result.Results))
		for _, item := range result.Results {
			if item.Name == "" {
				continue
			}
			imageURL := ""
			if item.PosterPath != "" {
				imageURL = tmdb.GetImageURL(item.PosterPath, tmdb.W500)
			}
			results = append(results, &VideoSearchResult{
				ID:          fmt.Sprintf("%d", item.ID),
				Title:       item.Name,
				ReleaseYear: parseYear(item.FirstAirDate),
				ImageURL:    imageURL,
				Subtitle:    "TV Show",
			})
			if len(results) >= limit {
				break
			}
		}
		return results, nil
	}

	result, err := f.client.GetSearchMovies(searchLanguage, options)
	if err != nil {
		return nil, err
	}
	results := make([]*VideoSearchResult, 0, len(result.Results))
	for _, item := range result.Results {
		if item.Title == "" {
			continue
		}
		imageURL := ""
		if item.PosterPath != "" {
			imageURL = tmdb.GetImageURL(item.PosterPath, tmdb.W500)
		}
		results = append(results, &VideoSearchResult{
			ID:          fmt.Sprintf("%d", item.ID),
			Title:       item.Title,
			ReleaseYear: parseYear(item.ReleaseDate),
			ImageURL:    imageURL,
			Subtitle:    "Movie",
		})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
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

// FetchPersonMovieCredits fetches movie credits for a person by TMDB ID
func (f *VideoFetcher) FetchPersonMovieCredits(personID int) ([]*VideoMetadata, error) {
	credits, err := f.client.GetPersonMovieCredits(personID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch person movie credits: %w", err)
	}

	var movies []*VideoMetadata
	for _, credit := range credits.Cast {
		// Skip if no release date or adult
		if credit.ReleaseDate == "" || credit.Adult {
			continue
		}

		metadata := &VideoMetadata{
			MediaMetadata: MediaMetadata{
				MediaInfo: MediaInfo{
					Type:        MediaTypeMovie,
					Title:       credit.Title,
					ReleaseYear: parseYear(credit.ReleaseDate),
					ID:          fmt.Sprintf("%d", credit.ID),
				},
				Description: credit.Overview,
				URL:         fmt.Sprintf("https://www.themoviedb.org/movie/%d", credit.ID),
				Rating:      float32(credit.VoteAverage),
			},
			Runtime:     0, // Not available in credits
			Budget:      0,
			BoxOffice:   0,
			Cast:        nil, // Not available
			Crew:        nil,
			CastCredits: nil,
			CrewCredits: nil,
			Genres:      nil,
		}
		if credit.PosterPath != "" {
			metadata.ImageURL = fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", credit.PosterPath)
		}
		movies = append(movies, metadata)
	}

	return movies, nil
}

// FetchPersonTVShowCredits fetches TV show credits for a person by TMDB ID
func (f *VideoFetcher) FetchPersonTVShowCredits(personID int) ([]*VideoMetadata, error) {
	credits, err := f.client.GetPersonTVCredits(personID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch person TV show credits: %w", err)
	}

	var tvShows []*VideoMetadata
	for _, credit := range credits.Cast {
		// Skip if no first air date
		if credit.FirstAirDate == "" {
			continue
		}

		metadata := &VideoMetadata{
			MediaMetadata: MediaMetadata{
				MediaInfo: MediaInfo{
					Type:        MediaTypeTV,
					Title:       credit.Name,
					ReleaseYear: parseYear(credit.FirstAirDate),
					ID:          fmt.Sprintf("%d", credit.ID),
				},
				Description: credit.Overview,
				URL:         fmt.Sprintf("https://www.themoviedb.org/tv/%d", credit.ID),
				Rating:      float32(credit.VoteAverage),
			},
			Runtime:     0,
			Budget:      0,
			BoxOffice:   0,
			Cast:        nil,
			Crew:        nil,
			CastCredits: nil,
			CrewCredits: nil,
			Genres:      nil,
		}
		if credit.PosterPath != "" {
			metadata.ImageURL = fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", credit.PosterPath)
		}
		tvShows = append(tvShows, metadata)
	}

	return tvShows, nil
}
