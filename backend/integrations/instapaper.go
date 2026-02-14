package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	instapaperAPIBaseURL = "https://www.instapaper.com/api"
)

// InstapaperIntegration handles Instapaper API integration
type InstapaperIntegration struct {
	*BaseIntegration
	client *http.Client
}

// InstapaperBookmark represents an Instapaper bookmark
type InstapaperBookmark struct {
	BookmarkID        int     `json:"bookmark_id"`
	URL               string  `json:"url"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Time              int64   `json:"time"`
	Starred           int     `json:"starred"`
	PrivateSource     string  `json:"private_source"`
	Hash              string  `json:"hash"`
	Progress          float64 `json:"progress"`
	ProgressTimestamp int64   `json:"progress_timestamp"`
}

// InstapaperFolder represents an Instapaper folder
type InstapaperFolder struct {
	FolderID int    `json:"folder_id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Position int    `json:"position"`
}

// NewInstapaperIntegration creates a new Instapaper integration
func NewInstapaperIntegration() *InstapaperIntegration {
	base := NewBaseIntegration("instapaper", []MediaType{MediaTypeArticle})

	return &InstapaperIntegration{
		BaseIntegration: base,
		client:          &http.Client{Timeout: 30 * time.Second},
	}
}

// Authenticate authenticates with Instapaper using username/password
func (i *InstapaperIntegration) Authenticate(ctx context.Context, credentials map[string]string) error {
	username, ok := credentials["username"]
	if !ok {
		username = os.Getenv("INSTAPAPER_USERNAME")
	}

	password, ok := credentials["password"]
	if !ok {
		password = os.Getenv("INSTAPAPER_PASSWORD")
	}

	if username == "" || password == "" {
		return fmt.Errorf("instapaper username and password are required")
	}

	i.SetCredential("username", username)
	i.SetCredential("password", password)

	// Test credentials by trying to access account info
	if err := i.testCredentials(ctx); err != nil {
		return fmt.Errorf("failed to authenticate with Instapaper: %w", err)
	}

	i.SetAuthenticated(true)
	return nil
}

// SyncUserData syncs user's Instapaper bookmarks
func (i *InstapaperIntegration) SyncUserData(ctx context.Context, userID uuid.UUID) (*SyncResult, error) {
	if !i.IsAuthenticated() {
		return nil, fmt.Errorf("instapaper integration not authenticated")
	}

	result := &SyncResult{
		IntegrationName: i.GetName(),
		UserID:          userID,
		SyncedAt:        time.Now(),
		MediaData:       make(map[MediaType][]interface{}),
	}

	// Sync bookmarks
	bookmarks, err := i.getUserBookmarksListing(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to fetch bookmarks: %v", err))
		return result, err
	}

	articleData := make([]interface{}, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		bookmarkMap := i.convertInstapaperBookmarkToMap(bookmark)
		articleData = append(articleData, bookmarkMap)
		result.ItemsProcessed++
	}

	result.MediaData[MediaTypeArticle] = articleData
	result.ItemsAdded = len(bookmarks)

	return result, nil
}

// testCredentials tests the Instapaper credentials
func (i *InstapaperIntegration) testCredentials(ctx context.Context) error {
	username, _ := i.GetCredential("username")
	password, _ := i.GetCredential("password")

	// Use the account verification endpoint
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, "POST", instapaperAPIBaseURL+"/1/account/verify_credentials", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := i.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("instapaper API returned status %d", resp.StatusCode)
	}

	return nil
}

// getUserBookmarksListing fetches user's bookmarks from Instapaper
func (i *InstapaperIntegration) getUserBookmarksListing(ctx context.Context) ([]InstapaperBookmark, error) {
	username, _ := i.GetCredential("username")
	password, _ := i.GetCredential("password")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("limit", "100") // Maximum allowed by API

	req, err := http.NewRequestWithContext(ctx, "POST", instapaperAPIBaseURL+"/1/bookmarks/list", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instapaper API returned status %d", resp.StatusCode)
	}

	// Instapaper API returns a JSON array where the first element is metadata
	// and subsequent elements are bookmark objects
	var responseArray []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseArray); err != nil {
		return nil, err
	}

	var bookmarksRaw []interface{}
	if len(responseArray) > 1 {
		bookmarksRaw = responseArray[1:]
	}

	// Convert to bookmark structs
	bookmarks := make([]InstapaperBookmark, 0, len(bookmarksRaw))
	for _, bookmarkRaw := range bookmarksRaw {
		if bookmarkMap, ok := bookmarkRaw.(map[string]interface{}); ok {
			bookmark := i.parseBookmarkFromMap(bookmarkMap)
			bookmarks = append(bookmarks, bookmark)
		}
	}

	return bookmarks, nil
}

// getUserFolders fetches user's folders from Instapaper
func (i *InstapaperIntegration) getUserFolders(ctx context.Context) ([]InstapaperFolder, error) {
	username, _ := i.GetCredential("username")
	password, _ := i.GetCredential("password")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, "POST", instapaperAPIBaseURL+"/1/folders/list", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instapaper API returned status %d", resp.StatusCode)
	}

	var folders []InstapaperFolder
	if err := json.NewDecoder(resp.Body).Decode(&folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// parseBookmarkFromMap converts a map to InstapaperBookmark struct
func (i *InstapaperIntegration) parseBookmarkFromMap(bookmarkMap map[string]interface{}) InstapaperBookmark {
	bookmark := InstapaperBookmark{}

	if id, ok := bookmarkMap["bookmark_id"].(float64); ok {
		bookmark.BookmarkID = int(id)
	}

	if bookmarkURL, ok := bookmarkMap["url"].(string); ok {
		bookmark.URL = bookmarkURL
	}

	if title, ok := bookmarkMap["title"].(string); ok {
		bookmark.Title = title
	}

	if description, ok := bookmarkMap["description"].(string); ok {
		bookmark.Description = description
	}

	if timestamp, ok := bookmarkMap["time"].(float64); ok {
		bookmark.Time = int64(timestamp)
	}

	if starred, ok := bookmarkMap["starred"].(float64); ok {
		bookmark.Starred = int(starred)
	}

	if privateSource, ok := bookmarkMap["private_source"].(string); ok {
		bookmark.PrivateSource = privateSource
	}

	if hash, ok := bookmarkMap["hash"].(string); ok {
		bookmark.Hash = hash
	}

	if progress, ok := bookmarkMap["progress"].(float64); ok {
		bookmark.Progress = progress
	}

	if progressTimestamp, ok := bookmarkMap["progress_timestamp"].(float64); ok {
		bookmark.ProgressTimestamp = int64(progressTimestamp)
	}

	return bookmark
}

// convertInstapaperBookmarkToMap converts an Instapaper bookmark to a generic map for processing
func (i *InstapaperIntegration) convertInstapaperBookmarkToMap(bookmark InstapaperBookmark) map[string]interface{} {
	data := map[string]interface{}{
		"title":       bookmark.Title,
		"url":         bookmark.URL,
		"external_id": strconv.Itoa(bookmark.BookmarkID),
		"source":      "instapaper",
		"type":        "article",
		"starred":     bookmark.Starred == 1,
		"progress":    bookmark.Progress,
	}

	if bookmark.Description != "" {
		data["description"] = bookmark.Description
	}

	if bookmark.Time > 0 {
		bookmarkTime := time.Unix(bookmark.Time, 0)
		data["added_at"] = bookmarkTime.Format(time.RFC3339)
		data["release_date"] = bookmarkTime.Format("2006-01-02")
	}

	if bookmark.ProgressTimestamp > 0 {
		progressTime := time.Unix(bookmark.ProgressTimestamp, 0)
		data["last_read_at"] = progressTime.Format(time.RFC3339)
	}

	if bookmark.PrivateSource != "" {
		data["private_source"] = bookmark.PrivateSource
	}

	// Calculate reading status based on progress
	if bookmark.Progress >= 1.0 {
		data["reading_status"] = "completed"
	} else if bookmark.Progress > 0 {
		data["reading_status"] = "in_progress"
	} else {
		data["reading_status"] = "unread"
	}

	return data
}
