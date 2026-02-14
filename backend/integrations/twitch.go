package integrations

import (
	"context"
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
	twitchTokenURL   = "https://id.twitch.tv/oauth2/token"
	twitchAPIBaseURL = "https://api.twitch.tv/helix"
)

// TwitchIntegration handles Twitch API integration
type TwitchIntegration struct {
	*BaseIntegration
	client      *http.Client
	accessToken string
	tokenExpiry time.Time
}

// TwitchTokenResponse represents the Twitch OAuth token response
type TwitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TwitchGame represents a Twitch game
type TwitchGame struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	BoxArt string `json:"box_art_url"`
}

// TwitchFollowedChannelsResponse represents followed channels from Twitch
type TwitchFollowedChannelsResponse struct {
	Data []struct {
		BroadcasterID   string    `json:"broadcaster_id"`
		BroadcasterName string    `json:"broadcaster_name"`
		FollowedAt      time.Time `json:"followed_at"`
	} `json:"data"`
	Total      int `json:"total"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// TwitchUserResponse represents user information from Twitch
type TwitchUserResponse struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

// NewTwitchIntegration creates a new Twitch integration
func NewTwitchIntegration() *TwitchIntegration {
	base := NewBaseIntegration("twitch", []MediaType{MediaTypeVideo})

	return &TwitchIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with Twitch using client credentials
func (t *TwitchIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	clientID, ok := credentials["client_id"]
	if !ok {
		clientID = os.Getenv("TWITCH_CLIENT_ID")
	}

	clientSecret, ok := credentials["client_secret"]
	if !ok {
		clientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	}

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("twitch client ID and secret are required")
	}

	t.SetCredential("client_id", clientID)
	t.SetCredential("client_secret", clientSecret)

	token, err := t.getAccessToken(clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Twitch: %w", err)
	}

	t.accessToken = token
	t.tokenExpiry = time.Now().Add(55 * 24 * time.Hour) // App access tokens last ~60 days, refresh early
	t.SetAuthenticated(true)

	return nil
}

// SyncUserData syncs user's Twitch data (requires user access token for most endpoints)
func (t *TwitchIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !t.IsAuthenticated() {
		return nil, fmt.Errorf("twitch integration not authenticated")
	}

	// Ensure token is still valid
	if time.Now().After(t.tokenExpiry) {
		clientID, _ := t.GetCredential("client_id")
		clientSecret, _ := t.GetCredential("client_secret")

		token, err := t.getAccessToken(clientID, clientSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh Twitch token: %w", err)
		}

		t.accessToken = token
		t.tokenExpiry = time.Now().Add(55 * 24 * time.Hour)
	}

	result := &SyncResult{
		IntegrationName: t.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// Note: Most user-specific endpoints require user access tokens obtained through OAuth flow
	// For now, we'll fetch top games as an example of what's possible with app access tokens
	games, err := t.getTopGames(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch top games: %v", err))
	} else {
		videoData := make([]interface{}, 0, len(games))
		for _, game := range games {
			gameData := t.convertTwitchGameToMap(game)
			videoData = append(videoData, gameData)
			result.ItemsProcessed++
		}
		result.MediaData[MediaTypeVideo] = videoData
		result.ItemsAdded = len(games)
	}

	return result, nil
}

// getAccessToken obtains an access token using client credentials flow
func (t *TwitchIntegration) getAccessToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", twitchTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("twitch API returned status %d", resp.StatusCode)
	}

	var token TwitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// getTopGames fetches the top games on Twitch (example endpoint that works with app access token)
func (t *TwitchIntegration) getTopGames(ctx context.Context) ([]TwitchGame, error) {
	clientID, _ := t.GetCredential("client_id")

	apiURL := fmt.Sprintf("%s/games/top?first=20", twitchAPIBaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	req.Header.Set("Client-Id", clientID)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API returned status %d", resp.StatusCode)
	}

	var response struct {
		Data []TwitchGame `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// getUserFollowedChannels fetches user's followed channels (requires user access token)
func (t *TwitchIntegration) getUserFollowedChannels(ctx context.Context, userID string) ([]interface{}, error) {
	clientID, _ := t.GetCredential("client_id")

	apiURL := fmt.Sprintf("%s/channels/followed?user_id=%s&first=100", twitchAPIBaseURL, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	req.Header.Set("Client-Id", clientID)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API returned status %d", resp.StatusCode)
	}

	var response TwitchFollowedChannelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	channels := make([]interface{}, len(response.Data))
	for i, channel := range response.Data {
		channels[i] = map[string]interface{}{
			"broadcaster_id":   channel.BroadcasterID,
			"broadcaster_name": channel.BroadcasterName,
			"followed_at":      channel.FollowedAt,
			"type":             "channel",
			"source":           "twitch",
		}
	}

	return channels, nil
}

// convertTwitchGameToMap converts a Twitch game to a generic map for processing
func (t *TwitchIntegration) convertTwitchGameToMap(game TwitchGame) map[string]interface{} {
	data := map[string]interface{}{
		"title":       game.Name,
		"external_id": game.ID,
		"source":      "twitch",
		"type":        "game",
	}

	if game.BoxArt != "" {
		// Twitch box art URLs have template parameters that need to be replaced
		boxArtURL := strings.ReplaceAll(game.BoxArt, "{width}", "285")
		boxArtURL = strings.ReplaceAll(boxArtURL, "{height}", "380")
		data["cover_url"] = boxArtURL
	}

	return data
}
