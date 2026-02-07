package integrations

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	spotifyTokenURL   = "https://accounts.spotify.com/api/token"
	spotifyAPIBaseURL = "https://api.spotify.com/v1"
)

// SpotifyIntegration handles Spotify API integration
type SpotifyIntegration struct {
	*BaseIntegration
	client      *http.Client
	accessToken string
	tokenExpiry time.Time
}

// SpotifyTokenResponse represents the Spotify OAuth token response
type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SpotifyAlbum represents a Spotify album
type SpotifyAlbum struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Artists     []SpotifyArtist `json:"artists"`
	ReleaseDate string          `json:"release_date"`
	Images      []SpotifyImage  `json:"images"`
	TotalTracks int32           `json:"total_tracks"`
	Label       string          `json:"label,omitempty"`
}

// SpotifyArtist represents a Spotify artist
type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SpotifyImage represents a Spotify image
type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

// SpotifySavedAlbumsResponse represents the response from the saved albums endpoint
type SpotifySavedAlbumsResponse struct {
	Items []struct {
		AddedAt string       `json:"added_at"`
		Album   SpotifyAlbum `json:"album"`
	} `json:"items"`
	Total int     `json:"total"`
	Next  *string `json:"next"`
}

// NewSpotifyIntegration creates a new Spotify integration
func NewSpotifyIntegration() *SpotifyIntegration {
	base := NewBaseIntegration("spotify", []MediaType{MediaTypeMusic})

	return &SpotifyIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with Spotify using client credentials
func (s *SpotifyIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	clientID, ok := credentials["client_id"]
	if !ok {
		clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	}

	clientSecret, ok := credentials["client_secret"]
	if !ok {
		clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	}

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("spotify client ID and secret are required")
	}

	s.SetCredential("client_id", clientID)
	s.SetCredential("client_secret", clientSecret)

	token, err := s.getAccessToken(clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Spotify: %w", err)
	}

	s.accessToken = token
	s.tokenExpiry = time.Now().Add(55 * time.Minute) // Tokens expire in 1 hour, refresh a bit early
	s.SetAuthenticated(true)

	return nil
}

// SyncUserData syncs user's saved albums from Spotify
func (s *SpotifyIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !s.IsAuthenticated() {
		return nil, fmt.Errorf("spotify integration not authenticated")
	}

	// Ensure token is still valid
	if time.Now().After(s.tokenExpiry) {
		clientID, _ := s.GetCredential("client_id")
		clientSecret, _ := s.GetCredential("client_secret")

		token, err := s.getAccessToken(clientID, clientSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh Spotify token: %w", err)
		}

		s.accessToken = token
		s.tokenExpiry = time.Now().Add(55 * time.Minute)
	}

	result := &SyncResult{
		IntegrationName: s.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// Sync saved albums
	albums, err := s.getUserSavedAlbums(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch saved albums: %v", err))
	} else {
		musicData := make([]interface{}, 0, len(albums))
		for _, album := range albums {
			albumData := s.convertSpotifyAlbumToMap(album)
			musicData = append(musicData, albumData)
			result.ItemsProcessed++
		}
		result.MediaData[MediaTypeMusic] = musicData
		result.ItemsAdded = len(albums)
	}

	return result, nil
}

// getAccessToken obtains an access token using client credentials flow
func (s *SpotifyIntegration) getAccessToken(clientID, clientSecret string) (string, error) {
	// Encode clientID:clientSecret to base64
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Prepare request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", spotifyTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify API returned status %d", resp.StatusCode)
	}

	// Decode response
	var token SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// getUserSavedAlbums fetches user's saved albums from Spotify
func (s *SpotifyIntegration) getUserSavedAlbums(ctx context.Context) ([]SpotifyAlbum, error) {
	var allAlbums []SpotifyAlbum
	nextURL := fmt.Sprintf("%s/me/albums?limit=50", spotifyAPIBaseURL)

	for nextURL != "" {
		req, err := http.NewRequestWithContext(ctx, "GET", nextURL, http.NoBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+s.accessToken)

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("spotify API returned status %d", resp.StatusCode)
		}

		var response SpotifySavedAlbumsResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		for _, item := range response.Items {
			allAlbums = append(allAlbums, item.Album)
		}

		if response.Next != nil {
			nextURL = *response.Next
		} else {
			nextURL = ""
		}
	}

	return allAlbums, nil
}

// convertSpotifyAlbumToMap converts a Spotify album to a generic map for processing
func (s *SpotifyIntegration) convertSpotifyAlbumToMap(album SpotifyAlbum) map[string]interface{} {
	data := map[string]interface{}{
		"title":        album.Name,
		"release_date": album.ReleaseDate,
		"track_count":  album.TotalTracks,
		"external_id":  album.ID,
		"source":       "spotify",
	}

	if len(album.Images) > 0 {
		// Get the largest image available
		var largestImage SpotifyImage
		maxSize := 0
		for _, img := range album.Images {
			if img.Height*img.Width > maxSize {
				maxSize = img.Height * img.Width
				largestImage = img
			}
		}
		if largestImage.URL != "" {
			data["cover_url"] = largestImage.URL
		}
	}

	if len(album.Artists) > 0 {
		artistNames := make([]string, len(album.Artists))
		for i, artist := range album.Artists {
			artistNames[i] = artist.Name
		}
		data["artists"] = artistNames
		data["description"] = fmt.Sprintf("By %s", strings.Join(artistNames, ", "))
	}

	if album.Label != "" {
		data["label"] = album.Label
	}

	return data
}
