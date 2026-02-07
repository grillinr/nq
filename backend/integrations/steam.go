package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	steamAPIBaseURL = "https://api.steampowered.com"
)

// SteamIntegration handles Steam API integration
type SteamIntegration struct {
	*BaseIntegration
	client *http.Client
}

// SteamGame represents a Steam game
type SteamGame struct {
	AppID                  int    `json:"appid"`
	Name                   string `json:"name"`
	PlaytimeForever        int    `json:"playtime_forever"`
	PlaytimeWindowsForever int    `json:"playtime_windows_forever"`
	PlaytimeMacForever     int    `json:"playtime_mac_forever"`
	PlaytimeLinuxForever   int    `json:"playtime_linux_forever"`
	ImgIconURL             string `json:"img_icon_url"`
	ImgLogoURL             string `json:"img_logo_url"`
}

// SteamOwnedGamesResponse represents the response from Steam's GetOwnedGames API
type SteamOwnedGamesResponse struct {
	Response struct {
		GameCount int         `json:"game_count"`
		Games     []SteamGame `json:"games"`
	} `json:"response"`
}

// SteamAppDetailsResponse represents the response from Steam's app details API
type SteamAppDetailsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		AppID       int    `json:"steam_appid"`
		Description string `json:"detailed_description"`
		ReleaseDate struct {
			Date string `json:"date"`
		} `json:"release_date"`
		HeaderImage string `json:"header_image"`
		Genres      []struct {
			Description string `json:"description"`
		} `json:"genres"`
		Developers []string `json:"developers"`
		Publishers []string `json:"publishers"`
	} `json:"data"`
}

// NewSteamIntegration creates a new Steam integration
func NewSteamIntegration() *SteamIntegration {
	base := NewBaseIntegration("steam", []MediaType{MediaTypeGame})

	return &SteamIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with Steam using API key and Steam ID
func (s *SteamIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	apiKey, ok := credentials["api_key"]
	if !ok {
		apiKey = os.Getenv("STEAM_API_KEY")
	}

	steamID, ok := credentials["steam_id"]
	if !ok {
		return fmt.Errorf("steam ID is required for authentication")
	}

	if apiKey == "" {
		return fmt.Errorf("steam API key is required")
	}

	s.SetCredential("api_key", apiKey)
	s.SetCredential("steam_id", steamID)

	// Test the credentials by making a simple API call
	if err := s.testCredentials(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with Steam: %w", err)
	}

	s.SetAuthenticated(true)
	return nil
}

// SyncUserData syncs user's Steam library
func (s *SteamIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !s.IsAuthenticated() {
		return nil, fmt.Errorf("steam integration not authenticated")
	}

	result := &SyncResult{
		IntegrationName: s.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// Sync owned games
	games, err := s.getUserOwnedGames(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch owned games: %v", err))
		return result, err
	}

	gameData := make([]interface{}, 0, len(games))
	for _, game := range games {
		// Get detailed game info
		details, err := s.getGameDetails(ctx, game.AppID)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch details for game %d: %v", game.AppID, err))
			continue
		}

		gameMap := s.convertSteamGameToMap(game, details)
		gameData = append(gameData, gameMap)
		result.ItemsProcessed++
	}

	result.MediaData[MediaTypeGame] = gameData
	result.ItemsAdded = len(gameData)

	return result, nil
}

// testCredentials tests the Steam API credentials
func (s *SteamIntegration) testCredentials(ctx context.Context) error {
	apiKey, _ := s.GetCredential("api_key")
	steamID, _ := s.GetCredential("steam_id")

	url := fmt.Sprintf("%s/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=false",
		steamAPIBaseURL, apiKey, steamID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam API returned status %d", resp.StatusCode)
	}

	return nil
}

// getUserOwnedGames fetches user's owned games from Steam
func (s *SteamIntegration) getUserOwnedGames(ctx context.Context) ([]SteamGame, error) {
	apiKey, _ := s.GetCredential("api_key")
	steamID, _ := s.GetCredential("steam_id")

	url := fmt.Sprintf("%s/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true&include_played_free_games=true",
		steamAPIBaseURL, apiKey, steamID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned status %d", resp.StatusCode)
	}

	var response SteamOwnedGamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Response.Games, nil
}

// getGameDetails fetches detailed information about a specific game
func (s *SteamIntegration) getGameDetails(ctx context.Context, appID int) (*SteamAppDetailsResponse, error) {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d", appID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam store API returned status %d", resp.StatusCode)
	}

	// The response is a map with the app ID as the key
	var responseMap map[string]SteamAppDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, err
	}

	appIDStr := strconv.Itoa(appID)
	if details, exists := responseMap[appIDStr]; exists && details.Success {
		return &details, nil
	}

	return nil, fmt.Errorf("game details not found for app ID %d", appID)
}

// convertSteamGameToMap converts Steam game data to a generic map for processing
func (s *SteamIntegration) convertSteamGameToMap(game SteamGame, details *SteamAppDetailsResponse) map[string]interface{} {
	data := map[string]interface{}{
		"title":       game.Name,
		"external_id": strconv.Itoa(game.AppID),
		"source":      "steam",
		"playtime":    game.PlaytimeForever, // in minutes
	}

	if details != nil && details.Success {
		if details.Data.Description != "" {
			data["description"] = details.Data.Description
		}

		if details.Data.ReleaseDate.Date != "" {
			data["release_date"] = details.Data.ReleaseDate.Date
		}

		if details.Data.HeaderImage != "" {
			data["cover_url"] = details.Data.HeaderImage
		}

		if len(details.Data.Genres) > 0 {
			genres := make([]string, len(details.Data.Genres))
			for i, genre := range details.Data.Genres {
				genres[i] = genre.Description
			}
			data["genres"] = genres
		}

		if len(details.Data.Developers) > 0 {
			data["developers"] = details.Data.Developers
		}

		if len(details.Data.Publishers) > 0 {
			data["publishers"] = details.Data.Publishers
		}
	}

	// Add Steam-specific metadata
	if game.ImgIconURL != "" {
		data["icon_url"] = fmt.Sprintf("https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", game.AppID, game.ImgIconURL)
	}

	if game.ImgLogoURL != "" {
		data["logo_url"] = fmt.Sprintf("https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg", game.AppID, game.ImgLogoURL)
	}

	return data
}
