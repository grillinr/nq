package db

import (
	"context"
	"fmt"
	"log"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateBook creates a new book in the database
func (r *Neo4jRepository) CreateBook(ctx context.Context, input model.CreateBookInput) (*model.Book, error) {
	// Try to enrich with metadata if minimal data provided
	if r.metadata != nil && shouldEnrichBook(input) {
		log.Printf("Enrichment triggered for book: title='%s'", input.Title)
		enrichedInput, err := r.enrichBookInput(input)
		if err != nil {
			log.Printf("Enrichment failed, continuing with original input: %v", err)
			// Continue with original input
		} else {
			input = enrichedInput
			log.Printf("Enrichment successful, using enriched input")
		}
	} else {
		log.Printf("Enrichment not triggered for book: title='%s' (metadata nil: %v, shouldEnrich: %v)",
			input.Title, r.metadata == nil, shouldEnrichBook(input))
	}

	bookID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (b:Book:Media {
				id: $id,
				title: $title,
				releaseDate: $releaseDate,
				description: $description,
				coverUrl: $coverUrl,
				pages: $pages,
				isbn: $isbn,
				publisher: $publisher,
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
			       b.description as description, b.coverUrl as coverUrl,
			       b.pages as pages, b.isbn as isbn, b.publisher as publisher
		`

		params := map[string]any{
			"id":          bookID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"pages":       input.Pages,
			"isbn":        input.Isbn,
			"publisher":   input.Publisher,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			book := &model.Book{
				ID:            bookID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Pages:         getInt32Pointer(record.AsMap()["pages"]),
				Isbn:          getStringPointer(record.AsMap()["isbn"]),
				Publisher:     getStringPointer(record.AsMap()["publisher"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return book, nil
		}

		return nil, fmt.Errorf("failed to create book")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Book), nil
}

// GetBookByID retrieves a book by its ID
func (r *Neo4jRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*model.Book, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (b:Book {id: $id})
			RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
			       b.description as description, b.coverUrl as coverUrl,
			       b.pages as pages, b.isbn as isbn, b.publisher as publisher
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			book := &model.Book{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Pages:         getInt32Pointer(record.AsMap()["pages"]),
				Isbn:          getStringPointer(record.AsMap()["isbn"]),
				Publisher:     getStringPointer(record.AsMap()["publisher"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return book, nil
		}

		return nil, fmt.Errorf("book not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Book), nil
}

// GetAllBooks retrieves all books
func (r *Neo4jRepository) GetAllBooks(ctx context.Context) ([]*model.Book, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (b:Book)
			RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
			       b.description as description, b.coverUrl as coverUrl,
			       b.pages as pages, b.isbn as isbn, b.publisher as publisher
			ORDER BY b.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var books []*model.Book
		for result.Next(ctx) {
			record := result.Record()
			bookID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			book := &model.Book{
				ID:            bookID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Pages:         getInt32Pointer(record.AsMap()["pages"]),
				Isbn:          getStringPointer(record.AsMap()["isbn"]),
				Publisher:     getStringPointer(record.AsMap()["publisher"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			books = append(books, book)
		}

		return books, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.Book), nil
}

// shouldEnrichBook determines if a book input should be enriched with metadata
func shouldEnrichBook(input model.CreateBookInput) bool {
	return input.Description == nil && input.CoverURL == nil && input.Pages == nil && input.Isbn == nil
}

// enrichBookInput fetches metadata and merges it with the input
func (r *Neo4jRepository) enrichBookInput(input model.CreateBookInput) (model.CreateBookInput, error) {
	if r.metadata == nil {
		log.Println("Metadata service not initialized, skipping enrichment")
		return input, nil
	}

	// Determine year from release date if available
	var year int
	if input.ReleaseDate != nil {
		// Simple year extraction - in production, use proper date parsing
		if len(*input.ReleaseDate) >= 4 {
			fmt.Sscanf(*input.ReleaseDate, "%d", &year)
		}
	}

	log.Printf("Attempting to enrich book: title='%s', year=%d", input.Title, year)

	// Try ISBN first if provided, otherwise search by title
	var metaInterface interface{}
	var err error

	if input.Isbn != nil && *input.Isbn != "" {
		log.Printf("Fetching metadata by ISBN: %s", *input.Isbn)
		metaInterface, err = r.metadata.GetMetadata(metadata.MediaInfo{
			Type: metadata.MediaTypeBook,
			ID:   *input.Isbn,
		})
	} else {
		log.Printf("Fetching metadata by title: '%s' (year: %d)", input.Title, year)
		metaInterface, err = r.metadata.GetMetadata(metadata.MediaInfo{
			Type:        metadata.MediaTypeBook,
			Title:       input.Title,
			ReleaseYear: year,
		})
	}

	if err != nil {
		log.Printf("Metadata fetch failed: %v", err)
		return input, err
	}

	meta, ok := metaInterface.(*metadata.BookMetadata)
	if !ok {
		return input, fmt.Errorf("unexpected metadata type for book")
	}

	log.Printf("Metadata fetched successfully: title='%s', description='%s', image='%s', id='%s'",
		meta.Title, meta.Description, meta.ImageURL, meta.ID)

	// Merge metadata with input (input takes precedence)
	enriched := input

	if enriched.Description == nil && meta.Description != "" {
		enriched.Description = &meta.Description
		log.Printf("Enriched description: '%s'", meta.Description)
	}

	if enriched.CoverURL == nil && meta.ImageURL != "" {
		enriched.CoverURL = &meta.ImageURL
		log.Printf("Enriched cover URL: '%s'", meta.ImageURL)
	}

	// Set ISBN if not provided and found in metadata
	if enriched.Isbn == nil && meta.ID != "" {
		enriched.Isbn = &meta.ID
		log.Printf("Enriched ISBN: '%s'", meta.ID)
	}

	if enriched.Pages == nil && meta.Pages > 0 {
		pages := int32(meta.Pages)
		enriched.Pages = &pages
		log.Printf("Enriched pages: %d", meta.Pages)
	}

	if enriched.Publisher == nil && meta.Publisher != "" {
		enriched.Publisher = &meta.Publisher
		log.Printf("Enriched publisher: '%s'", meta.Publisher)
	}

	if enriched.ReleaseDate == nil && meta.ReleaseYear > 0 {
		releaseDate := fmt.Sprintf("%d", meta.ReleaseYear)
		enriched.ReleaseDate = &releaseDate
		log.Printf("Enriched release date: '%s'", releaseDate)
	}

	return enriched, nil
}
