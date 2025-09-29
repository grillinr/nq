package metadata

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeBook  MediaType = "book"
	MediaTypeGame  MediaType = "game"
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
)

// MediaInfo contains the basic information needed to look up metadata
type MediaInfo struct {
	Type     MediaType
	Title    string
	Year     int
	ID       string // External ID (e.g., ISBN for books, TMDB ID for movies, etc.)
}

// MediaMetadata contains the full metadata for a media item
type MediaMetadata struct {
	Type        MediaType `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Year        int       `json:"year,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Rating      float32   `json:"rating,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	URL         string    `json:"url,omitempty"`
	// Add more fields as needed for specific media types
}

// Fetcher defines the interface for fetching metadata
// This allows for easy testing and swapping of implementations
type Fetcher interface {
	Fetch(info MediaInfo) (*MediaMetadata, error)
}
