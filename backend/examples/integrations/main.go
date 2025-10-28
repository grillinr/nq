package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/grillinr/nq/integrations"
)

func main() {
	// Create integration manager
	manager := integrations.NewManager()

	// Initialize and register integrations
	setupIntegrations(manager)

	// Example: Authenticate and sync data for a user
	userID := uuid.New()
	fmt.Printf("Syncing data for user: %s\n\n", userID)

	ctx := context.Background()

	// Sync data from all authenticated integrations
	results, err := manager.SyncAllUserData(ctx, userID)
	if err != nil {
		log.Printf("Some sync errors occurred: %v\n", err)
	}

	// Display results
	for integrationName, result := range results {
		fmt.Printf("=== %s Results ===\n", integrationName)
		fmt.Printf("Items processed: %d\n", result.ItemsProcessed)
		fmt.Printf("Items added: %d\n", result.ItemsAdded)
		fmt.Printf("Sync time: %s\n", result.SyncedAt.Format("2006-01-02 15:04:05"))

		if len(result.Errors) > 0 {
			fmt.Printf("Errors: %v\n", result.Errors)
		}

		// Show media data summary
		for mediaType, items := range result.MediaData {
			fmt.Printf("%s items: %d\n", mediaType, len(items))

			// Show first item as example
			if len(items) > 0 {
				if itemMap, ok := items[0].(map[string]interface{}); ok {
					if title, ok := itemMap["title"].(string); ok {
						fmt.Printf("  Example: %s\n", title)
					}
				}
			}
		}
		fmt.Println()
	}
}

func setupIntegrations(manager *integrations.Manager) {
	ctx := context.Background()

	// Setup Spotify integration
	if spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID"); spotifyClientID != "" {
		spotify := integrations.NewSpotifyIntegration()
		if err := spotify.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate Spotify: %v", err)
		} else {
			manager.RegisterIntegration(spotify)
			fmt.Println("✅ Spotify integration authenticated")
		}
	}

	// Setup Steam integration
	if steamAPIKey := os.Getenv("STEAM_API_KEY"); steamAPIKey != "" {
		steam := integrations.NewSteamIntegration()
		// Steam requires Steam ID for user-specific data
		if steamID := os.Getenv("STEAM_USER_ID"); steamID != "" {
			credentials := map[string]string{
				"steam_id": steamID,
			}
			if err := steam.Authenticate(ctx, credentials); err != nil {
				log.Printf("Failed to authenticate Steam: %v", err)
			} else {
				manager.RegisterIntegration(steam)
				fmt.Println("✅ Steam integration authenticated")
			}
		}
	}

	// Setup YouTube integration
	if youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY"); youtubeAPIKey != "" {
		youtube := integrations.NewYouTubeIntegration()
		if err := youtube.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate YouTube: %v", err)
		} else {
			manager.RegisterIntegration(youtube)
			fmt.Println("✅ YouTube integration authenticated")
		}
	}

	// Setup YouTube Music integration
	if youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY"); youtubeAPIKey != "" {
		ytMusic := integrations.NewYouTubeMusicIntegration()
		if err := ytMusic.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate YouTube Music: %v", err)
		} else {
			manager.RegisterIntegration(ytMusic)
			fmt.Println("✅ YouTube Music integration authenticated")
		}
	}

	// Setup Twitch integration
	if twitchClientID := os.Getenv("TWITCH_CLIENT_ID"); twitchClientID != "" {
		twitch := integrations.NewTwitchIntegration()
		if err := twitch.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate Twitch: %v", err)
		} else {
			manager.RegisterIntegration(twitch)
			fmt.Println("✅ Twitch integration authenticated")
		}
	}

	// Setup Apple Music integration
	if appleMusicToken := os.Getenv("APPLE_MUSIC_DEVELOPER_TOKEN"); appleMusicToken != "" {
		appleMusic := integrations.NewAppleMusicIntegration()
		if err := appleMusic.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate Apple Music: %v", err)
		} else {
			manager.RegisterIntegration(appleMusic)
			fmt.Println("✅ Apple Music integration authenticated")
		}
	}

	// Setup Instapaper integration
	if instapaperUsername := os.Getenv("INSTAPAPER_USERNAME"); instapaperUsername != "" {
		instapaper := integrations.NewInstapaperIntegration()
		if err := instapaper.Authenticate(ctx, nil); err != nil {
			log.Printf("Failed to authenticate Instapaper: %v", err)
		} else {
			manager.RegisterIntegration(instapaper)
			fmt.Println("✅ Instapaper integration authenticated")
		}
	}

	fmt.Printf("\nRegistered %d integrations\n\n", len(manager.ListIntegrations()))
}
