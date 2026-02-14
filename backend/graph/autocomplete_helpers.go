package graph

import (
	"math"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"
)

// buildSuggestion is a helper that extracts common suggestion building logic
func buildSuggestion(title string, releaseYear int, id, imageURL, subtitle string) *model.MediaSuggestion {
	var year *int32
	if releaseYear > 0 && releaseYear <= math.MaxInt32 {
		y := int32(releaseYear)
		year = &y
	}
	return &model.MediaSuggestion{
		Title:      title,
		Year:       year,
		ExternalID: stringPointer(id),
		ImageURL:   stringPointer(imageURL),
		Subtitle:   stringPointer(subtitle),
	}
}

func mapVideoSuggestions(results []*metadata.VideoSearchResult) []*model.MediaSuggestion {
	suggestions := make([]*model.MediaSuggestion, 0, len(results))
	for _, item := range results {
		if item == nil || item.Title == "" {
			continue
		}
		suggestions = append(suggestions, buildSuggestion(item.Title, item.ReleaseYear, item.ID, item.ImageURL, item.Subtitle))
	}
	return suggestions
}

func mapBookSuggestions(results []*metadata.BookSearchResult) []*model.MediaSuggestion {
	suggestions := make([]*model.MediaSuggestion, 0, len(results))
	for _, item := range results {
		if item == nil || item.Title == "" {
			continue
		}
		suggestions = append(suggestions, buildSuggestion(item.Title, item.ReleaseYear, item.ID, item.ImageURL, item.Subtitle))
	}
	return suggestions
}

func mapGameSuggestions(results []*metadata.MediaMetadata, limit int) []*model.MediaSuggestion {
	if limit <= 0 {
		limit = 10
	}
	suggestions := make([]*model.MediaSuggestion, 0, len(results))
	for _, item := range results {
		if item == nil || item.Title == "" {
			continue
		}
		var year *int32
		if item.ReleaseYear > 0 && item.ReleaseYear <= math.MaxInt32 {
			y := int32(item.ReleaseYear)
			year = &y
		}
		suggestions = append(suggestions, &model.MediaSuggestion{
			Title:      item.Title,
			Year:       year,
			ExternalID: stringPointer(item.ID),
			ImageURL:   stringPointer(item.ImageURL),
			Subtitle:   stringPointer(item.URL),
		})
		if len(suggestions) >= limit {
			break
		}
	}
	return suggestions
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
