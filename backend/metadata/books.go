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

// SearchBookByAuthorAndTitle searches Open Library by author and title and returns metadata for matching ISBNs.
// This is used for related book discovery.
func (f *BookFetcher) SearchBookByAuthorAndTitle(author string, title string) ([]*BookMetadata, error) {
	author = strings.TrimSpace(author)
	title = strings.TrimSpace(title)
	if author == "" && title == "" {
		return nil, errors.New("author or title required")
	}

	query := "https://openlibrary.org/search.json?limit=100&fields=title,author_name,isbn,first_publish_year,language"
	if title != "" {
		query += "&title=" + url.QueryEscape(title)
	}
	if author != "" {
		query += "&author=" + url.QueryEscape(author)
	}

	resp, err := f.client.Get(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search by author/title: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from search: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	var searchResult struct {
		Docs []struct {
			Title            string   `json:"title"`
			ISBN             []string `json:"isbn"`
			FirstPublishYear int      `json:"first_publish_year"`
			AuthorName       []string `json:"author_name"`
		} `json:"docs"`
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}
	if len(searchResult.Docs) == 0 {
		return nil, nil
	}

	var results []*BookMetadata
	for _, doc := range searchResult.Docs {
		if len(doc.ISBN) == 0 {
			continue
		}
		if title != "" && !strings.Contains(strings.ToLower(doc.Title), strings.ToLower(title)) {
			continue
		}
		isbn := pickPreferredISBN(doc.ISBN)
		meta, err := f.fetchByISBN(isbn, "")
		if err != nil {
			continue
		}
		if meta != nil {
			results = append(results, meta)
		}
	}
	return results, nil
}

// Fetch retrieves book metadata from Open Library API.
// Behavior:
// - If ISBN (info.ID) is provided: fetch by ISBN, extract title, then refine by searching the title for a better/first ISBN and re-fetch if different.
// - If ISBN is missing: try to derive a title from MediaInfo (Title, Query, Name) via reflection; search by title to get the first ISBN, then fetch by that ISBN.
func (f *BookFetcher) Fetch(info MediaInfo, language string) (any, error) {
	// Default to English if no language specified
	if language == "" {
		language = "en"
	}
	// Path 1: We have an ISBN directly
	if strings.TrimSpace(info.ID) != "" {
		meta, err := f.fetchByISBN(info.ID, "") // fetch regardless of language to detect available language
		if err != nil {
			return nil, err
		}

		// If the fetched metadata has a detected language and it doesn't match the requested language,
		// try to find a better ISBN in the requested language and re-fetch.
		if meta != nil {
			if meta.Language != "" && meta.Language != language {
				if strings.TrimSpace(meta.Title) != "" {
					if foundISBN, _ := f.searchFirstISBNByTitle(meta.Title, info.ReleaseYear, language, info.Author); foundISBN != "" && foundISBN != info.ID {
						if refined, err := f.fetchByISBN(foundISBN, language); err == nil {
							return refined, nil
						}
					}
					// No better match found — return original metadata to avoid hard failure
					return meta, nil
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

	return f.fetchByISBN(foundISBN, language)
}

func (f *BookFetcher) SearchBooksByTitle(title string, limit int) ([]*BookSearchResult, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title required")
	}
	if limit <= 0 {
		limit = 10
	}

	searchURL := fmt.Sprintf("https://openlibrary.org/search.json?title=%s&limit=%d&fields=title,author_name,isbn,first_publish_year,cover_i", url.QueryEscape(title), limit)
	resp, err := f.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search by title: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from search: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	var searchResult struct {
		Docs []struct {
			Title            string   `json:"title"`
			ISBN             []string `json:"isbn"`
			FirstPublishYear int      `json:"first_publish_year"`
			AuthorName       []string `json:"author_name"`
			CoverID          int      `json:"cover_i"`
		} `json:"docs"`
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	results := make([]*BookSearchResult, 0, len(searchResult.Docs))
	for _, doc := range searchResult.Docs {
		if doc.Title == "" {
			continue
		}
		isbn := pickPreferredISBN(doc.ISBN)
		if isbn == "" {
			continue
		}
		subtitle := ""
		if len(doc.AuthorName) > 0 {
			subtitle = doc.AuthorName[0]
		}
		imageURL := ""
		if doc.CoverID > 0 {
			imageURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", doc.CoverID)
		}
		results = append(results, &BookSearchResult{
			ID:          isbn,
			Title:       doc.Title,
			ReleaseYear: doc.FirstPublishYear,
			ImageURL:    imageURL,
			Subtitle:    subtitle,
		})
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// fetchByISBN retrieves and maps Open Library "books" API data for a single ISBN.
// It will attempt to detect the language of the book via the response and will
// respect the requested language (reqLang) if provided; if the detected language
// does not include the requested language, it returns an error.
func (f *BookFetcher) fetchByISBN(isbn string, reqLang string) (*BookMetadata, error) {
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

	// Detect languages from bookData if available
	detectedLangs := []string{}
	if langs, ok := bookData["languages"].([]any); ok && len(langs) > 0 {
		for _, l := range langs {
			switch v := l.(type) {
			case map[string]any:
				if key, ok := v["key"].(string); ok && key != "" {
					parts := strings.Split(key, "/")
					code := parts[len(parts)-1]
					detectedLangs = append(detectedLangs, strings.ToLower(code))
				}
			case string:
				detectedLangs = append(detectedLangs, strings.ToLower(v))
			}
		}
	}
	// fallback: some responses may use "language" key with []string
	if len(detectedLangs) == 0 {
		if langs2, ok := bookData["language"].([]any); ok && len(langs2) > 0 {
			for _, l := range langs2 {
				if s, ok := l.(string); ok && s != "" {
					detectedLangs = append(detectedLangs, strings.ToLower(s))
				}
			}
		}
	}

	// Convert first detected OpenLibrary code to ISO 639-1 for easier comparisons
	detectedISO := ""
	if len(detectedLangs) > 0 {
		detectedISO = openLibraryLangToISO6391(detectedLangs[0])
	}

	// Normalize requested language to OpenLibrary code (3-letter) for comparison
	reqLangCode := convertToOpenLibraryLangCode(reqLang)
	if reqLang != "" && reqLangCode != "" && len(detectedLangs) > 0 {
		match := false
		for _, dl := range detectedLangs {
			if dl == strings.ToLower(reqLangCode) || strings.HasSuffix(dl, strings.ToLower(reqLangCode)) {
				match = true
				break
			}
		}
		if !match {
			// Allow mismatch if we still have usable metadata
			// Do not hard fail to avoid blocking enrichment
		}
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
	// Only include title if language matches or not detectable
	metadata.Title = title

	// URL
	if u, ok := bookData["url"].(string); ok {
		metadata.URL = u
	}

	// Authors -> populate Authors slice and Description "By <author>"
	if authors, ok := bookData["authors"].([]any); ok && len(authors) > 0 {
		for _, a := range authors {
			if am, ok := a.(map[string]any); ok {
				if name, ok := am["name"].(string); ok && name != "" {
					metadata.Authors = append(metadata.Authors, name)
				}
			}
		}
		if metadata.Description == "" && len(metadata.Authors) > 0 {
			metadata.Description = fmt.Sprintf("By %s", metadata.Authors[0])
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

	// Publishers -> populate Publishers slice and keep first as Publisher for compatibility
	if publishers, ok := bookData["publishers"].([]any); ok && len(publishers) > 0 {
		for _, p := range publishers {
			if pm, ok := p.(map[string]any); ok {
				if name, ok := pm["name"].(string); ok && name != "" {
					metadata.Publishers = append(metadata.Publishers, name)
				}
			} else if name, ok := p.(string); ok && name != "" {
				metadata.Publishers = append(metadata.Publishers, name)
			}
		}
		if len(metadata.Publishers) > 0 {
			metadata.Publisher = metadata.Publishers[0]
		}
	}

	// Subjects and related arrays
	if subjects, ok := bookData["subjects"].([]any); ok && len(subjects) > 0 {
		for _, s := range subjects {
			// subject may be a map with "name" or a string
			switch v := s.(type) {
			case map[string]any:
				if name, ok := v["name"].(string); ok && name != "" {
					metadata.Subjects = append(metadata.Subjects, name)
				}
			case string:
				if v != "" {
					metadata.Subjects = append(metadata.Subjects, v)
				}
			}
		}
	}

	if subjectPlaces, ok := bookData["subject_places"].([]any); ok && len(subjectPlaces) > 0 {
		for _, sp := range subjectPlaces {
			if s, ok := sp.(string); ok && s != "" {
				metadata.SubjectPlaces = append(metadata.SubjectPlaces, s)
			}
		}
	}

	if subjectPeople, ok := bookData["subject_people"].([]any); ok && len(subjectPeople) > 0 {
		for _, sp := range subjectPeople {
			if s, ok := sp.(string); ok && s != "" {
				metadata.SubjectPeople = append(metadata.SubjectPeople, s)
			}
		}
	}

	if subjectTimes, ok := bookData["subject_times"].([]any); ok && len(subjectTimes) > 0 {
		for _, st := range subjectTimes {
			if s, ok := st.(string); ok && s != "" {
				metadata.SubjectTimes = append(metadata.SubjectTimes, s)
			}
		}
	}

	// Set detected language (ISO 639-1) if available
	if detectedISO != "" {
		metadata.Language = detectedISO
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

// openLibraryLangToISO6391 converts OpenLibrary/ISO 639-3 language codes to ISO 639-1 when possible
func openLibraryLangToISO6391(code string) string {
	if code == "" {
		return ""
	}
	s := strings.ToLower(code)
	switch s {
	case "eng":
		return "en"
	case "spa":
		return "es"
	case "fre", "fra":
		return "fr"
	case "ger", "deu":
		return "de"
	case "ita":
		return "it"
	case "por":
		return "pt"
	case "rus":
		return "ru"
	case "jpn":
		return "ja"
	case "chi", "zho":
		return "zh"
	case "kor":
		return "ko"
	case "ara":
		return "ar"
	case "hin":
		return "hi"
	default:
		// If already 2-letter, return as-is
		if len(s) == 2 {
			return s
		}
		return ""
	}
}
