package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
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
	Themes []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"themes,omitempty"`
	Keywords []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"keywords,omitempty"`
	GameModes []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"game_modes,omitempty"`
	PlayerPerspectives []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"player_perspectives,omitempty"`
	Franchises []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"franchises,omitempty"`
	Platforms []struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"platforms,omitempty"`
	Cover struct {
		ID  int    `json:"id"`
		URL string `json:"url,omitempty"`
	} `json:"cover,omitempty"`
	URL string `json:"url,omitempty"`
}

// GameFetcher implements the Fetcher interface for games
type GameFetcher struct {
	clientID       string
	clientSecret   string
	accessToken    string
	tokenExpiresAt time.Time
	httpClient     *http.Client
	lastRequest    time.Time
	cache          map[string]cachedGameResult
	mu             sync.Mutex
}

type cachedGameResult struct {
	items     []*MediaMetadata
	expiresAt time.Time
}

var ErrIGDBAuthFailed = errors.New("igdb auth failed")

// NewGameFetcher creates a new GameFetcher
func NewGameFetcher() (*GameFetcher, error) {
	clientID := os.Getenv("IGDB_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("IGDB_CLIENT_ID environment variable is not set")
	}
	clientSecret := os.Getenv("IGDB_CLIENT_SECRET")

	accessToken := os.Getenv("IGDB_ACCESS_TOKEN")
	f := &GameFetcher{
		clientID:     clientID,
		clientSecret: clientSecret,
		accessToken:  accessToken,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		cache:        make(map[string]cachedGameResult),
	}

	// If access token isn't set, attempt to obtain one from IGDB_REQUEST_URL if present
	if f.accessToken == "" {
		if reqURL := os.Getenv("IGDB_REQUEST_URL"); reqURL != "" {
			resp, err := f.httpClient.Get(strings.Trim(reqURL, "\""))
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					body, _ := io.ReadAll(resp.Body)
					var tokenResp map[string]any
					if err := json.Unmarshal(body, &tokenResp); err == nil {
						if at, ok := tokenResp["access_token"].(string); ok && at != "" {
							f.accessToken = at
							if exp, ok := tokenResp["expires_in"].(float64); ok && exp > 0 {
								f.tokenExpiresAt = time.Now().Add(time.Duration(exp) * time.Second)
							}
						}
					}
				}
			}
		}
	}

	return f, nil
}

func (f *GameFetcher) ensureAccessToken() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.accessToken != "" && (f.tokenExpiresAt.IsZero() || time.Now().Before(f.tokenExpiresAt.Add(-1*time.Minute))) {
		return nil
	}

	if f.clientSecret == "" {
		return fmt.Errorf("IGDB_CLIENT_SECRET environment variable is not set")
	}

	log.Printf("IGDB: refreshing access token")
	form := url.Values{}
	form.Set("client_id", f.clientID)
	form.Set("client_secret", f.clientSecret)
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("IGDB: token refresh failed status=%d", resp.StatusCode)
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return fmt.Errorf("token response missing access_token")
	}
	f.accessToken = tokenResp.AccessToken
	if tokenResp.ExpiresIn > 0 {
		f.tokenExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	log.Printf("IGDB: token refresh succeeded, expires_in=%d", tokenResp.ExpiresIn)
	return nil
}

// SearchRelatedGames searches IGDB by title and returns a small set of candidate games.
func (f *GameFetcher) SearchRelatedGames(title string) ([]*MediaMetadata, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title required")
	}

	cacheKey := strings.ToLower(strings.TrimSpace(title))
	if cached := f.getCachedResults(cacheKey); cached != nil {
		return cached, nil
	}
	query := fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,genres.name,themes.name,keywords.name,game_modes.name,player_perspectives.name,franchises.name,platforms.name,cover.url,url; limit 10;`, title)
	body, err := f.makeIGDBRequest(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch related games: %w", err)
	}

	var games []IGDBGame
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("failed to parse game data: %w", err)
	}
	if len(games) == 0 {
		return nil, nil
	}

	results := make([]*MediaMetadata, 0, len(games))
	for _, game := range games {
		metadata := &MediaMetadata{
			MediaInfo: MediaInfo{
				Type:  MediaTypeGame,
				Title: game.Name,
				ID:    fmt.Sprintf("%d", game.ID),
			},
			Description: game.Summary,
			URL:         game.URL,
		}
		if game.FirstReleaseDate != 0 {
			timeObj := time.Unix(game.FirstReleaseDate, 0)
			metadata.ReleaseYear = timeObj.Year()
		}
		if len(game.Genres) > 0 {
			metadata.Genres = make([]string, len(game.Genres))
			for i, genre := range game.Genres {
				metadata.Genres[i] = genre.Name
			}
		}
		if len(game.Themes) > 0 {
			metadata.Themes = make([]string, len(game.Themes))
			for i, theme := range game.Themes {
				metadata.Themes[i] = theme.Name
			}
		}
		if len(game.Keywords) > 0 {
			metadata.Keywords = make([]string, len(game.Keywords))
			for i, keyword := range game.Keywords {
				metadata.Keywords[i] = keyword.Name
			}
		}
		if len(game.GameModes) > 0 {
			metadata.GameModes = make([]string, len(game.GameModes))
			for i, mode := range game.GameModes {
				metadata.GameModes[i] = mode.Name
			}
		}
		if len(game.PlayerPerspectives) > 0 {
			metadata.Perspectives = make([]string, len(game.PlayerPerspectives))
			for i, perspective := range game.PlayerPerspectives {
				metadata.Perspectives[i] = perspective.Name
			}
		}
		if len(game.Franchises) > 0 {
			metadata.Franchises = make([]string, len(game.Franchises))
			for i, franchise := range game.Franchises {
				metadata.Franchises[i] = franchise.Name
			}
		}
		if len(game.Platforms) > 0 {
			metadata.Platforms = make([]string, len(game.Platforms))
			for i, platform := range game.Platforms {
				metadata.Platforms[i] = platform.Name
			}
		}
		if game.Cover.URL != "" {
			// Normalize protocol-relative URLs and adjust IGDB size segment to t_cover_big_2x.
			coverURL := game.Cover.URL
			if strings.HasPrefix(coverURL, "//") {
				coverURL = "https:" + coverURL
			}

			if parsed, err := url.Parse(coverURL); err == nil {
				pathSegments := strings.Split(parsed.Path, "/")
				for i, seg := range pathSegments {
					if strings.HasPrefix(seg, "t_") {
						pathSegments[i] = "t_cover_big_2x"
						break
					}
				}
				parsed.Path = strings.Join(pathSegments, "/")
				metadata.ImageURL = parsed.String()
			} else {
				// Fall back to the normalized URL if parsing fails.
				metadata.ImageURL = coverURL
			}
		}
		results = append(results, metadata)
	}
	f.setCachedResults(cacheKey, results)
	return results, nil
}

func (f *GameFetcher) getCachedResults(key string) []*MediaMetadata {
	f.mu.Lock()
	defer f.mu.Unlock()

	entry, ok := f.cache[key]
	if !ok {
		return nil
	}
	if time.Now().After(entry.expiresAt) {
		delete(f.cache, key)
		return nil
	}
	return entry.items
}

func (f *GameFetcher) setCachedResults(key string, items []*MediaMetadata) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(items) == 0 {
		return
	}
	entry := cachedGameResult{
		items:     items,
		expiresAt: time.Now().Add(10 * time.Minute),
	}
	f.cache[key] = entry
}

// makeIGDBRequest makes a request to the IGDB API
func (f *GameFetcher) makeIGDBRequest(query string) ([]byte, error) {
	return f.makeIGDBRequestWithRetry(query, true)
}

func (f *GameFetcher) makeIGDBRequestWithRetry(query string, allowRetry bool) ([]byte, error) {
	if err := f.ensureAccessToken(); err != nil {
		// Try environment as a last resort
		if f.accessToken == "" {
			f.accessToken = os.Getenv("IGDB_ACCESS_TOKEN")
		}
		if f.accessToken == "" {
			return nil, fmt.Errorf("IGDB access token not provided")
		}
	}

	f.throttle()

	req, err := http.NewRequest("POST", "https://api.igdb.com/v4/games", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Client-ID", f.clientID)
	req.Header.Set("Authorization", "Bearer "+f.accessToken)
	req.Header.Set("Content-Type", "text/plain")

	// Make the request
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to IGDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			// Attempt one refresh and retry
			if err := f.ensureAccessToken(); err == nil {
				return f.makeIGDBRequestWithRetry(query, false)
			}
			return nil, ErrIGDBAuthFailed
		}
		if resp.StatusCode == http.StatusTooManyRequests && allowRetry {
			log.Printf("IGDB: rate limit hit, backing off")
			time.Sleep(1 * time.Second)
			return f.makeIGDBRequestWithRetry(query, false)
		}
		return nil, fmt.Errorf("IGDB API returned status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (f *GameFetcher) throttle() {
	f.mu.Lock()
	defer f.mu.Unlock()

	minInterval := 300 * time.Millisecond
	now := time.Now()
	if !f.lastRequest.IsZero() {
		elapsed := now.Sub(f.lastRequest)
		if elapsed < minInterval {
			time.Sleep(minInterval - elapsed)
		}
	}
	f.lastRequest = time.Now()
}

// Fetch retrieves game metadata from IGDB
func (f *GameFetcher) Fetch(info MediaInfo, language string) (any, error) {
	if info.ID == "" && info.Title == "" {
		return nil, errors.New("either game ID or title is required")
	}

	// Build the query
	var query string
	if info.ID != "" {
		// Search by ID
		query = fmt.Sprintf(`fields name,summary,first_release_date,genres.name,themes.name,keywords.name,game_modes.name,player_perspectives.name,franchises.name,platforms.name,cover.url,url; where id = %s;`, info.ID)
	} else {
		// Search by title
		query = fmt.Sprintf(`search "%s"; fields name,summary,first_release_date,genres.name,themes.name,keywords.name,game_modes.name,player_perspectives.name,franchises.name,platforms.name,cover.url,url; limit 1;`, info.Title)
	}

	// Add language filtering if specified (IGDB supports language codes)
	if language != "" && language != "en" {
		// Note: IGDB language filtering is more complex, but we'll add it for future enhancement
		// For now, we'll keep the query as-is since language filtering in IGDB requires specific language IDs
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
		MediaInfo: MediaInfo{
			Type:  MediaTypeGame,
			Title: game.Name,
		},
		Description: game.Summary,
		URL:         game.URL,
	}

	// Set the release year if available
	if game.FirstReleaseDate != 0 {
		// Convert Unix timestamp to time.Time
		timeObj := time.Unix(game.FirstReleaseDate, 0)
		metadata.ReleaseYear = timeObj.Year()
	}

	// Get genres if available
	if len(game.Genres) > 0 {
		metadata.Genres = make([]string, len(game.Genres))
		for i, genre := range game.Genres {
			metadata.Genres[i] = genre.Name
		}
	}
	if len(game.Themes) > 0 {
		metadata.Themes = make([]string, len(game.Themes))
		for i, theme := range game.Themes {
			metadata.Themes[i] = theme.Name
		}
	}
	if len(game.Keywords) > 0 {
		metadata.Keywords = make([]string, len(game.Keywords))
		for i, keyword := range game.Keywords {
			metadata.Keywords[i] = keyword.Name
		}
	}
	if len(game.GameModes) > 0 {
		metadata.GameModes = make([]string, len(game.GameModes))
		for i, mode := range game.GameModes {
			metadata.GameModes[i] = mode.Name
		}
	}
	if len(game.PlayerPerspectives) > 0 {
		metadata.Perspectives = make([]string, len(game.PlayerPerspectives))
		for i, perspective := range game.PlayerPerspectives {
			metadata.Perspectives[i] = perspective.Name
		}
	}
	if len(game.Franchises) > 0 {
		metadata.Franchises = make([]string, len(game.Franchises))
		for i, franchise := range game.Franchises {
			metadata.Franchises[i] = franchise.Name
		}
	}
	if len(game.Platforms) > 0 {
		metadata.Platforms = make([]string, len(game.Platforms))
		for i, platform := range game.Platforms {
			metadata.Platforms[i] = platform.Name
		}
	}

	// Get cover image if available
	if game.Cover.URL != "" {
		// Construct the full URL for the cover image
		// The URL from IGDB is often protocol-relative (e.g., //images.igdb.com/...)
		imageID := strings.TrimPrefix(game.Cover.URL, "//")
		metadata.ImageURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big_2x/%s.jpg", imageID)
	}

	return metadata, nil
}
