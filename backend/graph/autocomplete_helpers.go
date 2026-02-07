package graph

import (
	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"
)

func mapVideoSuggestions(results []*metadata.VideoSearchResult) []*model.MediaSuggestion {
	suggestions := make([]*model.MediaSuggestion, 0, len(results))
	for _, item := range results {
		if item == nil || item.Title == "" {
			continue
		}
		var year *int32
		if item.ReleaseYear > 0 {
			y := int32(item.ReleaseYear)
			year = &y
		}
		suggestions = append(suggestions, &model.MediaSuggestion{
			Title:      item.Title,
			Year:       year,
			ExternalID: stringPointer(item.ID),
			ImageURL:   stringPointer(item.ImageURL),
			Subtitle:   stringPointer(item.Subtitle),
		})
	}
	return suggestions
}

func mapBookSuggestions(results []*metadata.BookSearchResult) []*model.MediaSuggestion {
	suggestions := make([]*model.MediaSuggestion, 0, len(results))
	for _, item := range results {
		if item == nil || item.Title == "" {
			continue
		}
		var year *int32
		if item.ReleaseYear > 0 {
			y := int32(item.ReleaseYear)
			year = &y
		}
		suggestions = append(suggestions, &model.MediaSuggestion{
			Title:      item.Title,
			Year:       year,
			ExternalID: stringPointer(item.ID),
			ImageURL:   stringPointer(item.ImageURL),
			Subtitle:   stringPointer(item.Subtitle),
		})
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
		if item.ReleaseYear > 0 {
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
