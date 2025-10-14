package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	appleMusicAPIBaseURL = "https://api.music.apple.com/v1"
)

// AppleMusicIntegration handles Apple Music API integration
type AppleMusicIntegration struct {
	*BaseIntegration
	client *http.Client
}

// AppleMusicAlbum represents an Apple Music album
type AppleMusicAlbum struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        string    `json:"name"`
		ArtistName  string    `json:"artistName"`
		ReleaseDate time.Time `json:"releaseDate"`
		TrackCount  int32     `json:"trackCount"`
		Genre       string    `json:"genreNames"`
		RecordLabel string    `json:"recordLabel"`
		Copyright   string    `json:"copyright"`
		URL         string    `json:"url"`
		Artwork     struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
		} `json:"artwork"`
	} `json:"attributes"`
}

// AppleMusicSearchResponse represents an Apple Music search response
type AppleMusicSearchResponse struct {
	Results struct {
		Albums struct {
			Data []AppleMusicAlbum `json:"data"`
		} `json:"albums"`
		Songs struct {
			Data []AppleMusicSong `json:"data"`
		} `json:"songs"`
	} `json:"results"`
}

// AppleMusicSong represents an Apple Music song
type AppleMusicSong struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name       string `json:"name"`
		ArtistName string `json:"artistName"`
		AlbumName  string `json:"albumName"`
		Duration   int32  `json:"durationInMillis"`
		URL        string `json:"url"`
		Artwork    struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
		} `json:"artwork"`
	} `json:"attributes"`
}

// NewAppleMusicIntegration creates a new Apple Music integration
func NewAppleMusicIntegration() *AppleMusicIntegration {
	base := NewBaseIntegration("apple_music", []MediaType{MediaTypeMusic})

	return &AppleMusicIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with Apple Music using developer token
func (a *AppleMusicIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	// Apple Music uses JWT developer tokens
	developerToken, ok := credentials["developer_token"]
	if !ok {
		developerToken = os.Getenv("APPLE_MUSIC_DEVELOPER_TOKEN")
	}

	if developerToken == "" {
		return fmt.Errorf("apple Music developer token is required")
	}

	a.SetCredential("developer_token", developerToken)

	// Test the developer token
	if err := a.testCredentials(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with Apple Music: %w", err)
	}

	a.SetAuthenticated(true)
	return nil
}

// SyncUserData syncs user's Apple Music data
// Note: User music library requires user token (MusicKit JS or user authentication)
func (a *AppleMusicIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !a.IsAuthenticated() {
		return nil, fmt.Errorf("apple Music integration not authenticated")
	}

	result := &SyncResult{
		IntegrationName: a.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// For demonstration, search for popular albums
	// In a real implementation, you'd need user tokens to access personal library
	albums, err := a.searchAlbums(ctx, "rock", 10)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to search albums: %v", err))
	} else {
		musicData := make([]interface{}, 0, len(albums))
		for _, album := range albums {
			albumMap := a.convertAppleMusicAlbumToMap(album)
			musicData = append(musicData, albumMap)
			result.ItemsProcessed++
		}
		result.MediaData[MediaTypeMusic] = musicData
		result.ItemsAdded = len(albums)
	}

	return result, nil
}

// testCredentials tests the Apple Music developer token
func (a *AppleMusicIntegration) testCredentials(ctx context.Context) error {
	developerToken, _ := a.GetCredential("developer_token")

	// Try to search for a simple query to test the token
	url := fmt.Sprintf("%s/catalog/us/search?term=test&types=albums&limit=1", appleMusicAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+developerToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("apple Music API returned status %d", resp.StatusCode)
	}

	return nil
}

// searchAlbums searches for albums in the Apple Music catalog
func (a *AppleMusicIntegration) searchAlbums(ctx context.Context, query string, limit int) ([]AppleMusicAlbum, error) {
	developerToken, _ := a.GetCredential("developer_token")

	url := fmt.Sprintf("%s/catalog/us/search?term=%s&types=albums&limit=%d",
		appleMusicAPIBaseURL, query, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+developerToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple Music API returned status %d", resp.StatusCode)
	}

	var response AppleMusicSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results.Albums.Data, nil
}

// getUserLibraryAlbums fetches user's library albums (requires user token)
func (a *AppleMusicIntegration) getUserLibraryAlbums(ctx context.Context, userToken string) ([]AppleMusicAlbum, error) {
	developerToken, _ := a.GetCredential("developer_token")

	url := fmt.Sprintf("%s/me/library/albums", appleMusicAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+developerToken)
	req.Header.Set("Music-User-Token", userToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple Music API returned status %d", resp.StatusCode)
	}

	var response struct {
		Data []AppleMusicAlbum `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// convertAppleMusicAlbumToMap converts an Apple Music album to a generic map for processing
func (a *AppleMusicIntegration) convertAppleMusicAlbumToMap(album AppleMusicAlbum) map[string]interface{} {
	data := map[string]interface{}{
		"title":       album.Attributes.Name,
		"external_id": album.ID,
		"source":      "apple_music",
		"type":        "album",
		"url":         album.Attributes.URL,
		"artist":      album.Attributes.ArtistName,
		"track_count": album.Attributes.TrackCount,
	}

	if !album.Attributes.ReleaseDate.IsZero() {
		data["release_date"] = album.Attributes.ReleaseDate.Format("2006-01-02")
	}

	if album.Attributes.RecordLabel != "" {
		data["label"] = album.Attributes.RecordLabel
	}

	if album.Attributes.Genre != "" {
		// Split genre names if they're comma-separated
		genres := strings.Split(album.Attributes.Genre, ",")
		for i, genre := range genres {
			genres[i] = strings.TrimSpace(genre)
		}
		data["genres"] = genres
	}

	if album.Attributes.Artwork.URL != "" {
		// Apple Music artwork URLs have placeholders that need to be replaced
		artworkURL := strings.ReplaceAll(album.Attributes.Artwork.URL, "{w}", "600")
		artworkURL = strings.ReplaceAll(artworkURL, "{h}", "600")
		data["cover_url"] = artworkURL
	}

	if album.Attributes.Copyright != "" {
		data["copyright"] = album.Attributes.Copyright
	}

	data["description"] = fmt.Sprintf("Album by %s", album.Attributes.ArtistName)

	return data
}
