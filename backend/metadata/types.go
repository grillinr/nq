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
	Type        MediaType
	Title       string
	ReleaseYear int
	ID          string // External ID (e.g., ISBN for books, TMDB ID for movies, etc.)
}

// MediaMetadata contains the full metadata for a media item
type MediaMetadata struct {
	MediaInfo
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Rating      float32  `json:"rating,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	URL         string   `json:"url,omitempty"`
}

// BookMetadata contains additional metadata specific to books
type BookMetadata struct {
	MediaMetadata
	Pages         int      `json:"pages,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	Subjects      []string `json:"subjects,omitempty"`
	SubjectPlaces []string `json:"subject_places,omitempty"`
	SubjectPeople []string `json:"subject_people,omitempty"`
	SubjectTimes  []string `json:"subject_times,omitempty"`
}

// VideoMetadata contains additional metadata specific to movies and TV shows
type VideoMetadata struct {
	MediaMetadata
	Budget              float64  `json:"budget,omitempty"`
	Cast                []string `json:"cast,omitempty"`
	Crew                []string `json:"crew,omitempty"`
	Genres              []string `json:"genres,omitempty"`
	ProductionCompanies []string `json:"production_companies,omitempty"`
	ProductionCountries []string `json:"production_countries,omitempty"`
}

// Fetcher defines the interface for fetching metadata
// This allows for easy testing and swapping of implementations
type Fetcher interface {
	Fetch(info MediaInfo) (*MediaMetadata, error)
}
