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
	ytMusicAPIBaseURL = "https://www.googleapis.com/youtube/v3"
)

// YouTubeMusicIntegration handles YouTube Music API integration
// Note: YouTube Music doesn't have a separate API - it uses YouTube Data API v3
type YouTubeMusicIntegration struct {
	*BaseIntegration
	client *http.Client
}

// YTMusicVideo represents a music video from YouTube
type YTMusicVideo struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		PublishedAt time.Time `json:"publishedAt"`
		Thumbnails  struct {
			High struct {
				URL string `json:"url"`
			} `json:"high"`
			Maxres struct {
				URL string `json:"url"`
			} `json:"maxres"`
		} `json:"thumbnails"`
		ChannelTitle string   `json:"channelTitle"`
		CategoryId   string   `json:"categoryId"`
		Tags         []string `json:"tags"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
}

// YTMusicSearchResponse represents a YouTube Music search response
type YTMusicSearchResponse struct {
	Items         []YTMusicVideo `json:"items"`
	NextPageToken string         `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// YTMusicPlaylist represents a YouTube Music playlist
type YTMusicPlaylist struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			High struct {
				URL string `json:"url"`
			} `json:"high"`
		} `json:"thumbnails"`
		ChannelTitle string `json:"channelTitle"`
	} `json:"snippet"`
	ContentDetails struct {
		ItemCount int `json:"itemCount"`
	} `json:"contentDetails"`
}

// NewYouTubeMusicIntegration creates a new YouTube Music integration
func NewYouTubeMusicIntegration() *YouTubeMusicIntegration {
	base := NewBaseIntegration("youtube_music", []MediaType{MediaTypeMusic})

	return &YouTubeMusicIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with YouTube Music using API key
func (yt *YouTubeMusicIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	apiKey, ok := credentials["api_key"]
	if !ok {
		apiKey = os.Getenv("YOUTUBE_API_KEY")
	}

	if apiKey == "" {
		return fmt.Errorf("youtube API key is required for YouTube Music integration")
	}

	yt.SetCredential("api_key", apiKey)

	// Test the API key
	if err := yt.testCredentials(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with YouTube Music: %w", err)
	}

	yt.SetAuthenticated(true)
	return nil
}

// SyncUserData syncs user's YouTube Music data
func (yt *YouTubeMusicIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !yt.IsAuthenticated() {
		return nil, fmt.Errorf("youtube Music integration not authenticated")
	}

	result := &SyncResult{
		IntegrationName: yt.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// Search for music videos (category ID 10 is Music on YouTube)
	musicVideos, err := yt.searchMusicVideos(ctx, "popular music", 20)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch music videos: %v", err))
	} else {
		musicData := make([]interface{}, 0, len(musicVideos))
		for _, video := range musicVideos {
			videoMap := yt.convertYTMusicVideoToMap(video)
			musicData = append(musicData, videoMap)
			result.ItemsProcessed++
		}
		result.MediaData[MediaTypeMusic] = musicData
		result.ItemsAdded = len(musicVideos)
	}

	return result, nil
}

// testCredentials tests the YouTube API key
func (yt *YouTubeMusicIntegration) testCredentials(ctx context.Context) error {
	apiKey, _ := yt.GetCredential("api_key")

	// Test with a simple search in the Music category
	url := fmt.Sprintf("%s/search?part=snippet&q=music&type=video&categoryId=10&maxResults=1&key=%s",
		ytMusicAPIBaseURL, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := yt.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	return nil
}

// searchMusicVideos searches for music videos on YouTube
func (yt *YouTubeMusicIntegration) searchMusicVideos(ctx context.Context, query string, maxResults int) ([]YTMusicVideo, error) {
	apiKey, _ := yt.GetCredential("api_key")

	// Search in Music category (categoryId=10) with high view count
	url := fmt.Sprintf("%s/search?part=snippet&q=%s&type=video&categoryId=10&maxResults=%d&order=viewCount&key=%s",
		ytMusicAPIBaseURL, query, maxResults, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := yt.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	var response YTMusicSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Get detailed information for the videos
	if len(response.Items) > 0 {
		return yt.getMusicVideoDetails(ctx, response.Items)
	}

	return response.Items, nil
}

// getMusicVideoDetails fetches detailed information about music videos
func (yt *YouTubeMusicIntegration) getMusicVideoDetails(ctx context.Context, videos []YTMusicVideo) ([]YTMusicVideo, error) {
	apiKey, _ := yt.GetCredential("api_key")

	// Collect video IDs
	videoIDs := make([]string, len(videos))
	for i, video := range videos {
		videoIDs[i] = video.ID
	}

	// Join IDs with commas for batch request
	idsParam := strings.Join(videoIDs, ",")

	url := fmt.Sprintf("%s/videos?part=snippet,contentDetails&id=%s&key=%s",
		ytMusicAPIBaseURL, idsParam, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := yt.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	var response YTMusicSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// searchMusicPlaylists searches for music playlists
func (yt *YouTubeMusicIntegration) searchMusicPlaylists(ctx context.Context, query string, maxResults int) ([]YTMusicPlaylist, error) {
	apiKey, _ := yt.GetCredential("api_key")

	url := fmt.Sprintf("%s/search?part=snippet&q=%s music playlist&type=playlist&maxResults=%d&key=%s",
		ytMusicAPIBaseURL, query, maxResults, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := yt.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	var response struct {
		Items []YTMusicPlaylist `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// convertYTMusicVideoToMap converts a YouTube Music video to a generic map for processing
func (yt *YouTubeMusicIntegration) convertYTMusicVideoToMap(video YTMusicVideo) map[string]interface{} {
	data := map[string]interface{}{
		"title":       video.Snippet.Title,
		"description": video.Snippet.Description,
		"external_id": video.ID,
		"source":      "youtube_music",
		"type":        "music_video",
		"url":         fmt.Sprintf("https://music.youtube.com/watch?v=%s", video.ID),
		"artist":      video.Snippet.ChannelTitle,
	}

	if !video.Snippet.PublishedAt.IsZero() {
		data["release_date"] = video.Snippet.PublishedAt.Format("2006-01-02")
	}

	// Add thumbnail - prefer maxres, fallback to high
	if video.Snippet.Thumbnails.Maxres.URL != "" {
		data["cover_url"] = video.Snippet.Thumbnails.Maxres.URL
	} else if video.Snippet.Thumbnails.High.URL != "" {
		data["cover_url"] = video.Snippet.Thumbnails.High.URL
	}

	// Add duration if available
	if video.ContentDetails.Duration != "" {
		data["duration"] = video.ContentDetails.Duration
	}

	// Add tags as genres if available
	if len(video.Snippet.Tags) > 0 {
		data["genres"] = video.Snippet.Tags
	}

	// Try to extract artist and song info from title
	title := video.Snippet.Title
	if strings.Contains(title, " - ") {
		parts := strings.SplitN(title, " - ", 2)
		if len(parts) == 2 {
			data["artist"] = strings.TrimSpace(parts[0])
			data["song_title"] = strings.TrimSpace(parts[1])
		}
	}

	return data
}
