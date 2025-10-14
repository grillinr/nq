package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// BookFetcher implements the Fetcher interface for books
type BookFetcher struct {
	client *http.Client
}

// NewBookFetcher creates a new BookFetcher
func NewBookFetcher() (*BookFetcher, error) {
	return &BookFetcher{
		client: &http.Client{},
	}, nil
}

// Fetch retrieves book metadata from Open Library API
func (f *BookFetcher) Fetch(info MediaInfo) (*MediaMetadata, error) {
	if info.ID == "" {
		return nil, errors.New("ISBN is required for book lookup")
	}

	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", info.ID)
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if we got any results
	if len(result) == 0 {
		return nil, fmt.Errorf("no book found with ISBN: %s", info.ID)
	}

	// Get the first (and should be only) result
	var bookData map[string]interface{}
	for _, v := range result {
		bookData = v.(map[string]interface{})
		break
	}

	// Extract relevant fields
	metadata := &MediaMetadata{
		Type: MediaTypeBook,
	}

	// Safely extract title
	if title, ok := bookData["title"].(string); ok {
		metadata.Title = title
	} else {
		return nil, errors.New("book title not found in response")
	}

	// Safely extract URL
	if url, ok := bookData["url"].(string); ok {
		metadata.URL = url
	}

	// Safely extract authors
	if authors, ok := bookData["authors"].([]interface{}); ok && len(authors) > 0 {
		if author, ok := authors[0].(map[string]interface{}); ok {
			if name, ok := author["name"].(string); ok {
				metadata.Description = fmt.Sprintf("By %s", name)
			}
		}
	}

	if cover, ok := bookData["cover"].(map[string]interface{}); ok {
		if largeCover, ok := cover["large"].(string); ok && largeCover != "" {
			metadata.ImageURL = largeCover
		} else if mediumCover, ok := cover["medium"].(string); ok && mediumCover != "" {
			metadata.ImageURL = mediumCover
		} else if smallCover, ok := cover["small"].(string); ok && smallCover != "" {
			metadata.ImageURL = smallCover
		}
	}

	// Try to extract the publish date
	if publishDate, ok := bookData["publish_date"].(string); ok && publishDate != "" {
		// Try to extract the year from the publish date
		var year int
		_, err := fmt.Sscanf(publishDate, "%d", &year)
		if err == nil {
			metadata.Year = year
		}
	}

	return metadata, nil
}
