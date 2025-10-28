package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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

// Fetch retrieves book metadata from Open Library API.
// Behavior:
// - If ISBN (info.ID) is provided: fetch by ISBN, extract title, then refine by searching the title for a better/first ISBN and re-fetch if different.
// - If ISBN is missing: try to derive a title from MediaInfo (Title, Query, Name) via reflection; search by title to get the first ISBN, then fetch by that ISBN.
func (f *BookFetcher) Fetch(info MediaInfo, language string) (interface{}, error) {
	// Default to English if no language specified
	if language == "" {
		language = "en"
	}
	// Path 1: We have an ISBN directly
	if strings.TrimSpace(info.ID) != "" {
		meta, err := f.fetchByISBN(info.ID)
		if err != nil {
			return nil, err
		}

		// Refine by looking up the first ISBN for the fetched title, then re-fetch if different
		if strings.TrimSpace(meta.Title) != "" {
			if foundISBN, _ := f.searchFirstISBNByTitle(meta.Title, info.ReleaseYear, language, info.Author); foundISBN != "" && foundISBN != info.ID {
				if refined, err := f.fetchByISBN(foundISBN); err == nil {
					return refined, nil
				}
			}
		}
		return meta, nil
	}

	// Path 2: No ISBN provided — try to derive a title and resolve an ISBN from it
	title := extractTitleFromInfo(info)
	if title == "" {
		return nil, errors.New("ISBN or title is required for book lookup")
	}

	foundISBN, err := f.searchFirstISBNByTitle(title, info.ReleaseYear, language, info.Author)
	if err != nil {
		return nil, err
	}
	if foundISBN == "" {
		return nil, fmt.Errorf("no ISBN found for title: %s", title)
	}

	return f.fetchByISBN(foundISBN)
}

// fetchByISBN retrieves and maps Open Library "books" API data for a single ISBN.
func (f *BookFetcher) fetchByISBN(isbn string) (*BookMetadata, error) {
	reqURL := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", url.QueryEscape(isbn))
	resp, err := f.client.Get(reqURL)
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

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no book found with ISBN: %s", isbn)
	}

	// Get the first (usually only) book payload
	var bookData map[string]any
	for _, v := range result {
		if m, ok := v.(map[string]any); ok {
			bookData = m
			break
		}
	}
	if bookData == nil {
		return nil, fmt.Errorf("unexpected data for ISBN: %s", isbn)
	}

	metadata := &BookMetadata{
		MediaMetadata: MediaMetadata{
			MediaInfo: MediaInfo{
				Type: MediaTypeBook,
				ID:   isbn,
			},
		},
	}

	// Title
	title, _ := bookData["title"].(string)
	if title == "" {
		return nil, errors.New("book title not found in response")
	}
	metadata.Title = title

	// URL
	if u, ok := bookData["url"].(string); ok {
		metadata.URL = u
	}

	// Authors -> Description "By <author>"
	if authors, ok := bookData["authors"].([]any); ok && len(authors) > 0 {
		if author, ok := authors[0].(map[string]any); ok {
			if name, ok := author["name"].(string); ok && name != "" {
				metadata.Description = fmt.Sprintf("By %s", name)
			}
		}
	}

	// Cover
	if cover, ok := bookData["cover"].(map[string]any); ok {
		if largeCover, ok := cover["large"].(string); ok && largeCover != "" {
			metadata.ImageURL = largeCover
		} else if mediumCover, ok := cover["medium"].(string); ok && mediumCover != "" {
			metadata.ImageURL = mediumCover
		} else if smallCover, ok := cover["small"].(string); ok && smallCover != "" {
			metadata.ImageURL = smallCover
		}
	}

	// Publish date -> ReleaseYear (best-effort)
	if publishDate, ok := bookData["publish_date"].(string); ok && publishDate != "" {
		var year int
		_, err := fmt.Sscanf(publishDate, "%d", &year)
		if err == nil && year > 0 {
			metadata.ReleaseYear = year
		}
	}

	// Pages
	if pages, ok := bookData["number_of_pages"].(float64); ok {
		metadata.Pages = int(pages)
	}

	// Publisher
	if publishers, ok := bookData["publishers"].([]any); ok && len(publishers) > 0 {
		if pub, ok := publishers[0].(map[string]any); ok {
			if name, ok := pub["name"].(string); ok {
				metadata.Publisher = name
			}
		}
	}

	return metadata, nil
}

// searchFirstISBNByTitle queries Open Library search for a title and returns the first usable ISBN.
// Prefers a 13-digit ISBN if available.
func (f *BookFetcher) searchFirstISBNByTitle(title string, year int, language string, author string) (string, error) {
	// Convert language code to OpenLibrary format (ISO 639-3)
	// OpenLibrary uses 3-letter codes like "eng", "spa", "fre"
	langCode := convertToOpenLibraryLangCode(language)

	searchURL := fmt.Sprintf("https://openlibrary.org/search.json?title=%s&limit=10&fields=title,author_name,isbn,first_publish_year,language", url.QueryEscape(strings.TrimSpace(title)))

	// Add language filter if specified
	if langCode != "" {
		searchURL += "&language=" + langCode
	}
	resp, err := f.client.Get(searchURL)
	if err != nil {
		return "", fmt.Errorf("failed to search by title: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from search: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read search response: %w", err)
	}

	var searchResult struct {
		Docs []struct {
			ISBN             []string `json:"isbn"`
			FirstPublishYear int      `json:"first_publish_year"`
			Language         []string `json:"language"`
			AuthorName       []string `json:"author_name"`
		} `json:"docs"`
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return "", fmt.Errorf("failed to parse search response: %w", err)
	}
	if len(searchResult.Docs) == 0 {
		return "", nil
	}

	// Find the best matching doc: prefer language match, then author match, then year match
	bestIndex := 0
	bestScore := 0 // 0 = no match, 1 = has ISBN, 2 = year match, 3 = language match, 4 = language + year match, 5+ = author boosts

	for i, doc := range searchResult.Docs {
		if len(doc.ISBN) == 0 {
			continue // Skip docs without ISBN
		}

		score := 1 // Base score for having ISBN

		// Check language match
		langMatch := false
		if langCode != "" && len(doc.Language) > 0 {
			for _, docLang := range doc.Language {
				if strings.ToLower(docLang) == langCode {
					langMatch = true
					break
				}
			}
		}
		if langMatch {
			score += 2 // Language match is worth more than year match
		}

		// Check author match
		authorMatch := false
		if author != "" && len(doc.AuthorName) > 0 {
			for _, a := range doc.AuthorName {
				if strings.Contains(strings.ToLower(a), strings.ToLower(author)) {
					authorMatch = true
					break
				}
			}
		}
		if authorMatch {
			score += 3 // Author match is significant
		}

		// Check year match
		if year > 0 {
			yearDiff := abs(doc.FirstPublishYear - year)
			if yearDiff <= 1 { // Within 1 year
				score += 1
			}
		}

		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	if len(searchResult.Docs[bestIndex].ISBN) == 0 {
		return "", nil
	}
	return pickPreferredISBN(searchResult.Docs[bestIndex].ISBN), nil
}

// pickPreferredISBN prefers a 13-digit ISBN if present, otherwise returns the first non-empty.
func pickPreferredISBN(isbns []string) string {
	var first string
	for _, s := range isbns {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if first == "" {
			first = s
		}
		// Prefer ISBN-13 (13 digits)
		digits := 0
		for _, r := range s {
			if r >= '0' && r <= '9' {
				digits++
			}
		}
		if digits == 13 {
			return s
		}
	}
	return first
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// extractTitleFromInfo attempts to read a title-like field from MediaInfo using reflection,
// trying "Title", then "Query", then "Name".
func extractTitleFromInfo(info MediaInfo) string {
	v := reflect.ValueOf(info)
	if v.Kind() != reflect.Struct {
		return ""
	}
	for _, field := range []string{"Title", "Query", "Name"} {
		f := v.FieldByName(field)
		if f.IsValid() && f.Kind() == reflect.String {
			if s := strings.TrimSpace(f.String()); s != "" {
				return normalizeTitle(s)
			}
		}
	}
	return ""
}

// normalizeTitle cleans and shortens the title for better search results
func normalizeTitle(title string) string {
	// For series: subtitle, take the subtitle
	if colonIndex := strings.LastIndex(title, ":"); colonIndex > 0 {
		subtitle := strings.TrimSpace(title[colonIndex+1:])
		if subtitle != "" {
			return subtitle
		}
	}
	// Otherwise, remove after colon
	if colonIndex := strings.Index(title, ":"); colonIndex > 0 {
		title = title[:colonIndex]
	}
	// Trim spaces
	return strings.TrimSpace(title)
}

// convertToOpenLibraryLangCode converts ISO 639-1 language codes to OpenLibrary's ISO 639-3 format
func convertToOpenLibraryLangCode(lang string) string {
	if lang == "" {
		return ""
	}

	// Common conversions from ISO 639-1 to ISO 639-3
	switch strings.ToLower(lang) {
	case "en":
		return "eng"
	case "es", "spa":
		return "spa"
	case "fr", "fre":
		return "fre"
	case "de", "ger":
		return "ger"
	case "it", "ita":
		return "ita"
	case "pt", "por":
		return "por"
	case "ru", "rus":
		return "rus"
	case "ja", "jpn":
		return "jpn"
	case "zh", "chi":
		return "chi"
	case "ko", "kor":
		return "kor"
	case "ar", "ara":
		return "ara"
	case "hi", "hin":
		return "hin"
	default:
		// If it's already a 3-letter code, return as-is
		if len(lang) == 3 {
			return strings.ToLower(lang)
		}
		// For unknown codes, return empty string (no filtering)
		return ""
	}
}
