package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	youtubeAPIBaseURL = "https://www.googleapis.com/youtube/v3"
)

// YouTubeIntegration handles YouTube API integration
type YouTubeIntegration struct {
	*BaseIntegration
	client *http.Client
}

// YouTubeVideo represents a YouTube video
type YouTubeVideo struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		PublishedAt time.Time `json:"publishedAt"`
		Thumbnails  struct {
			High struct {
				URL string `json:"url"`
			} `json:"high"`
			Medium struct {
				URL string `json:"url"`
			} `json:"medium"`
		} `json:"thumbnails"`
		ChannelTitle string `json:"channelTitle"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
	Statistics struct {
		ViewCount    string `json:"viewCount"`
		LikeCount    string `json:"likeCount"`
		CommentCount string `json:"commentCount"`
	} `json:"statistics"`
}

// YouTubePlaylistResponse represents a YouTube playlist response
type YouTubePlaylistResponse struct {
	Items         []YouTubeVideo `json:"items"`
	NextPageToken string         `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

// YouTubeChannel represents a YouTube channel
type YouTubeChannel struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			High struct {
				URL string `json:"url"`
			} `json:"high"`
		} `json:"thumbnails"`
	} `json:"snippet"`
	Statistics struct {
		ViewCount             string `json:"viewCount"`
		SubscriberCount       string `json:"subscriberCount"`
		HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
		VideoCount            string `json:"videoCount"`
	} `json:"statistics"`
}

// YouTubeChannelsResponse represents a YouTube channels response
type YouTubeChannelsResponse struct {
	Items []YouTubeChannel `json:"items"`
}

// NewYouTubeIntegration creates a new YouTube integration
func NewYouTubeIntegration() *YouTubeIntegration {
	base := NewBaseIntegration("youtube", []MediaType{MediaTypeVideo})

	return &YouTubeIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with YouTube using API key
func (y *YouTubeIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	apiKey, ok := credentials["api_key"]
	if !ok {
		apiKey = os.Getenv("YOUTUBE_API_KEY")
	}

	if apiKey == "" {
		return fmt.Errorf("youtube API key is required")
	}

	y.SetCredential("api_key", apiKey)

	// Test the API key by making a simple request
	if err := y.testCredentials(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with YouTube: %w", err)
	}

	y.SetAuthenticated(true)
	return nil
}

// SyncUserData syncs user's YouTube data
// Note: This requires OAuth 2.0 for user-specific data like liked videos, subscriptions
func (y *YouTubeIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !y.IsAuthenticated() {
		return nil, fmt.Errorf("youtube integration not authenticated")
	}

	result := &SyncResult{
		IntegrationName: y.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// For demonstration, we'll search for popular videos
	// In a real implementation, you'd need OAuth to access user's personal data
	videos, err := y.searchPopularVideos(ctx, "music", 10)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch popular videos: %v", err))
	} else {
		videoData := make([]interface{}, 0, len(videos))
		for _, video := range videos {
			videoMap := y.convertYouTubeVideoToMap(video)
			videoData = append(videoData, videoMap)
			result.ItemsProcessed++
		}
		result.MediaData[MediaTypeVideo] = videoData
		result.ItemsAdded = len(videos)
	}

	return result, nil
}

// testCredentials tests the YouTube API key
func (y *YouTubeIntegration) testCredentials(ctx context.Context) error {
	apiKey, _ := y.GetCredential("api_key")

	url := fmt.Sprintf("%s/search?part=snippet&q=test&type=video&maxResults=1&key=%s",
		youtubeAPIBaseURL, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := y.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	return nil
}

// searchPopularVideos searches for popular videos in a category
func (y *YouTubeIntegration) searchPopularVideos(ctx context.Context, query string, maxResults int) ([]YouTubeVideo, error) {
	apiKey, _ := y.GetCredential("api_key")

	url := fmt.Sprintf("%s/search?part=snippet&q=%s&type=video&maxResults=%d&order=viewCount&key=%s",
		youtubeAPIBaseURL, query, maxResults, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	var response YouTubePlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Get additional details for the videos
	if len(response.Items) > 0 {
		return y.getVideoDetails(ctx, response.Items)
	}

	return response.Items, nil
}

// getVideoDetails fetches detailed information about specific videos
func (y *YouTubeIntegration) getVideoDetails(ctx context.Context, videos []YouTubeVideo) ([]YouTubeVideo, error) {
	apiKey, _ := y.GetCredential("api_key")

	// Collect video IDs
	videoIDs := make([]string, len(videos))
	for i, video := range videos {
		videoIDs[i] = video.ID
	}

	// Join IDs with commas for batch request
	idsParam := ""
	for i, id := range videoIDs {
		if i > 0 {
			idsParam += ","
		}
		idsParam += id
	}

	url := fmt.Sprintf("%s/videos?part=snippet,contentDetails,statistics&id=%s&key=%s",
		youtubeAPIBaseURL, idsParam, apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube API returned status %d", resp.StatusCode)
	}

	var response YouTubePlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// getUserLikedVideos fetches user's liked videos (requires OAuth)
func (y *YouTubeIntegration) getUserLikedVideos(ctx context.Context, accessToken string) ([]YouTubeVideo, error) {
	// This would require OAuth 2.0 flow and user access token
	// Implementation would be similar to other methods but with Authorization header
	return nil, fmt.Errorf("getUserLikedVideos requires OAuth implementation")
}

// convertYouTubeVideoToMap converts a YouTube video to a generic map for processing
func (y *YouTubeIntegration) convertYouTubeVideoToMap(video YouTubeVideo) map[string]interface{} {
	data := map[string]interface{}{
		"title":       video.Snippet.Title,
		"description": video.Snippet.Description,
		"external_id": video.ID,
		"source":      "youtube",
		"type":        "video",
		"url":         fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID),
		"channel":     video.Snippet.ChannelTitle,
	}

	if !video.Snippet.PublishedAt.IsZero() {
		data["release_date"] = video.Snippet.PublishedAt.Format("2006-01-02")
	}

	// Add thumbnail
	if video.Snippet.Thumbnails.High.URL != "" {
		data["cover_url"] = video.Snippet.Thumbnails.High.URL
	} else if video.Snippet.Thumbnails.Medium.URL != "" {
		data["cover_url"] = video.Snippet.Thumbnails.Medium.URL
	}

	// Add statistics if available
	if video.Statistics.ViewCount != "" {
		data["view_count"] = video.Statistics.ViewCount
	}
	if video.Statistics.LikeCount != "" {
		data["like_count"] = video.Statistics.LikeCount
	}
	if video.Statistics.CommentCount != "" {
		data["comment_count"] = video.Statistics.CommentCount
	}

	// Add duration if available
	if video.ContentDetails.Duration != "" {
		data["duration"] = video.ContentDetails.Duration
	}

	return data
}
