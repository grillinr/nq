package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// IGDBGame represents the game data structure from IGDB API
type IGDBGame struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Summary          string `json:"summary,omitempty"`
	FirstReleaseDate int64  `json:"first_release_date,omitempty"`
	Genres           []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"genres,omitempty"`
	Cover struct {
		ID  int    `json:"id"`
		URL string `json:"url,omitempty"`
	} `json:"cover,omitempty"`
	URL string `json:"url,omitempty"`
}

// GameFetcher implements the Fetcher interface for games
type GameFetcher struct {
	clientID   string
	httpClient *http.Client
}

// NewGameFetcher creates a new GameFetcher
func NewGameFetcher() (*GameFetcher, error) {
	clientID := os.Getenv("IGDB_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("IGDB_CLIENT_ID environment variable is not set")
	}

	return &GameFetcher{
		clientID:   clientID,
		httpClient: &http.Client{},
	}, nil
}

// makeIGDBRequest makes a request to the IGDB API
func (f *GameFetcher) makeIGDBRequest(query string) ([]byte, error) {
	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Client-ID", f.clientID)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("IGDB_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "text/plain")

	// Make the request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to IGDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IGDB API returned status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// Fetch retrieves game metadata from IGDB
func (f *GameFetcher) Fetch(info MediaInfo) (*MediaMetadata, error) {
	if info.ID == "" && info.Title == "" {
		return nil, errors.New("either game ID or title is required")
	}

	// Build the query
	var query string
	if info.ID != "" {
		// Search by ID
		query = fmt.Sprintf(`fields name,summary,first_release_date,genres.name,cover.url,url; where id = %s;`, info.ID)
	} else {
		// Search by title
		query = fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,genres.name,cover.url,url; limit 1;`, info.Title)
	}

	// Make the request
	body, err := f.makeIGDBRequest(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game data: %w", err)
	}

	// Parse the response
	var games []IGDBGame
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("failed to parse game data: %w", err)
	}

	if len(games) == 0 {
		if info.ID != "" {
			return nil, fmt.Errorf("no game found with ID: %s", info.ID)
		} else {
			return nil, fmt.Errorf("no game found with title: %s", info.Title)
		}
	}

	game := games[0]
	metadata := &MediaMetadata{
		Type:        MediaTypeGame,
		Title:       game.Name,
		Description: game.Summary,
		URL:         game.URL,
	}

	// Set the release year if available
	if game.FirstReleaseDate != 0 {
		// Convert Unix timestamp to time.Time
		timeObj := time.Unix(game.FirstReleaseDate, 0)
		metadata.Year = timeObj.Year()
	}

	// Get genres if available
	if len(game.Genres) > 0 {
		metadata.Genres = make([]string, len(game.Genres))
		for i, genre := range game.Genres {
			metadata.Genres[i] = genre.Name
		}
	}

	// Get cover image if available
	if game.Cover.URL != "" {
		// Construct the full URL for the cover image
		// The URL from IGDB is relative, needs to be prefixed with https://images.igdb.com/igdb/image/upload/
		// and suffixed with _cover_big_2x.jpg for a high-res image
		imageID := strings.TrimPrefix(game.Cover.URL, "//")
		metadata.ImageURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big_2x/%s.jpg", imageID)
	}

	return metadata, nil
}
