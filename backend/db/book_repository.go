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
						publishers: $publishers,
						createdAt: datetime(),
						updatedAt: datetime()
					})
					WITH b
					FOREACH (author IN CASE WHEN $authors IS NULL THEN [] ELSE $authors END |
						MERGE (p:Person {normalizedName: coalesce(author.normalizedName, toLower(trim(author.name)))})
						ON CREATE SET p.id = randomUUID(), p.name = author.name, p.createdAt = datetime()
						MERGE (p)-[:AUTHORED]->(b)
					)
					FOREACH (subject IN CASE WHEN $subjects IS NULL THEN [] ELSE $subjects END |
						MERGE (t:Tag {type: 'subject', normalizedName: coalesce(subject.normalizedName, toLower(trim(subject.name)))})
						ON CREATE SET t.id = randomUUID(), t.name = subject.name
						MERGE (t)-[:TAGGED]->(b)
					)
					FOREACH (pubName IN CASE WHEN $publishers IS NULL THEN [] ELSE $publishers END |
						MERGE (p:Publisher {name: pubName})
						ON CREATE SET p.id = randomUUID()
						MERGE (p)-[:PUBLISHED]->(b)
					)
					WITH b
					OPTIONAL MATCH (p:Publisher)-[:PUBLISHED]->(b)
					RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
						   b.description as description, b.coverUrl as coverUrl,
					       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
						   collect(DISTINCT p.name) as publisher_nodes,
						b.publishers as publishers
						`
		// Build author and subject param lists with normalized names computed in-app
		var authorsParam []map[string]any
		if input.Authors != nil {
			authorsParam = make([]map[string]any, len(input.Authors))
			for i, a := range input.Authors {
				authorsParam[i] = map[string]any{"name": a, "normalizedName": NormalizeName(a)}
			}
		}
		var subjectsParam []map[string]any
		if input.Subjects != nil {
			subjectsParam = make([]map[string]any, len(input.Subjects))
			for i, s := range input.Subjects {
				subjectsParam[i] = map[string]any{"name": s, "normalizedName": NormalizeName(s)}
			}
		}

		params := map[string]any{
			"id":          bookID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"pages":       input.Pages,
			"isbn":        input.Isbn,
			"publisher":   input.Publisher,
			"publishers":  input.Publishers,
			"authors":     authorsParam,
			"subjects":    subjectsParam,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			// Parse publisher nodes (may be nil) and fallback to property
			var publishers []string
			if pNodes, ok := record.AsMap()["publisher_nodes"]; ok && pNodes != nil {
				switch v := pNodes.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
			} else if p, ok := record.AsMap()["publishers"]; ok && p != nil {
				switch v := p.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
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
				Publishers:    publishers,
				Creators:      []*model.Creator{},
				Authors:       []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Subjects:      []*model.Tag{},
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
				OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
				OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
				RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
					       b.description as description, b.coverUrl as coverUrl,
					       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
					       collect(DISTINCT {id: person.id, name: person.name}) as authors,
					       collect(DISTINCT {id: t.id, name: t.name, type: t.type}) as subjects,
					       collect(DISTINCT pnode.name) as publisher_nodes,
					       b.publishers as publishers
					`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()

			// Parse authors
			authors := []*model.Creator{}
			if a, ok := record.AsMap()["authors"]; ok && a != nil {
				switch v := a.(type) {
				case []any:
					for _, e := range v {
						if m, ok := e.(map[string]any); ok {
							var aid uuid.UUID
							if idStr, _ := m["id"].(string); idStr != "" {
								if parsed, err := uuid.Parse(idStr); err == nil {
									aid = parsed
								}
							}
							name, _ := m["name"].(string)
							authors = append(authors, &model.Creator{ID: aid, Name: name})
						}
					}
				}
			}

			// Parse subjects
			subjects := []*model.Tag{}
			if s, ok := record.AsMap()["subjects"]; ok && s != nil {
				switch v := s.(type) {
				case []any:
					for _, e := range v {
						if m, ok := e.(map[string]any); ok {
							var sid uuid.UUID
							if idStr, _ := m["id"].(string); idStr != "" {
								if parsed, err := uuid.Parse(idStr); err == nil {
									sid = parsed
								}
							}
							name, _ := m["name"].(string)
							typeStr, _ := m["type"].(string)
							subjects = append(subjects, &model.Tag{ID: sid, Name: name, Type: typeStr})
						}
					}
				}
			}

			// Parse publishers (may be nil) prefer publisher_nodes
			var publishers []string
			if pNodes, ok := record.AsMap()["publisher_nodes"]; ok && pNodes != nil {
				switch v := pNodes.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
			} else if p, ok := record.AsMap()["publishers"]; ok && p != nil {
				switch v := p.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
			}

			book := &model.Book{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Pages:         getInt32Pointer(record.AsMap()["pages"]),
				Isbn:          getStringPointer(record.AsMap()["isbn"]),
				Publisher:     getStringPointer(record.AsMap()["publisher"]),
				Publishers:    publishers,
				Creators:      []*model.Creator{},
				Authors:       authors,
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Subjects:      subjects,
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
				OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
				OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
				RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
					       b.description as description, b.coverUrl as coverUrl,
					       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
					       collect(DISTINCT {id: person.id, name: person.name}) as authors,
					       collect(DISTINCT {id: t.id, name: t.name, type: t.type}) as subjects,
					       collect(DISTINCT pnode.name) as publisher_nodes,
					       b.publishers as publishers
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

			// Parse authors
			authors := []*model.Creator{}
			if a, ok := record.AsMap()["authors"]; ok && a != nil {
				switch v := a.(type) {
				case []any:
					for _, e := range v {
						if m, ok := e.(map[string]any); ok {
							var aid uuid.UUID
							if idStr, _ := m["id"].(string); idStr != "" {
								if parsed, err := uuid.Parse(idStr); err == nil {
									aid = parsed
								}
							}
							name, _ := m["name"].(string)
							authors = append(authors, &model.Creator{ID: aid, Name: name})
						}
					}
				}
			}

			// Parse subjects
			subjects := []*model.Tag{}
			if s, ok := record.AsMap()["subjects"]; ok && s != nil {
				switch v := s.(type) {
				case []any:
					for _, e := range v {
						if m, ok := e.(map[string]any); ok {
							var sid uuid.UUID
							if idStr, _ := m["id"].(string); idStr != "" {
								if parsed, err := uuid.Parse(idStr); err == nil {
									sid = parsed
								}
							}
							name, _ := m["name"].(string)
							typeStr, _ := m["type"].(string)
							subjects = append(subjects, &model.Tag{ID: sid, Name: name, Type: typeStr})
						}
					}
				}
			}

			// Parse publishers (may be nil) prefer publisher_nodes
			var publishers []string
			if pNodes, ok := record.AsMap()["publisher_nodes"]; ok && pNodes != nil {
				switch v := pNodes.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
			} else if p, ok := record.AsMap()["publishers"]; ok && p != nil {
				switch v := p.(type) {
				case []any:
					for _, e := range v {
						if s, ok := e.(string); ok {
							publishers = append(publishers, s)
						}
					}
				case []string:
					publishers = v
				}
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
				Publishers:    publishers,
				Creators:      []*model.Creator{},
				Authors:       authors,
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Subjects:      subjects,
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

	// Populate authors, publishers, subjects arrays if present
	if len(meta.Authors) > 0 && len(enriched.Authors) == 0 {
		enriched.Authors = meta.Authors
		log.Printf("Enriched authors: %v", meta.Authors)
	}
	if len(meta.Publishers) > 0 && len(enriched.Publishers) == 0 {
		enriched.Publishers = meta.Publishers
		log.Printf("Enriched publishers: %v", meta.Publishers)
	}
	if len(meta.Subjects) > 0 && len(enriched.Subjects) == 0 {
		enriched.Subjects = meta.Subjects
		log.Printf("Enriched subjects: %v", meta.Subjects)
	}

	if enriched.ReleaseDate == nil && meta.ReleaseYear > 0 {
		releaseDate := fmt.Sprintf("%d", meta.ReleaseYear)
		enriched.ReleaseDate = &releaseDate
		log.Printf("Enriched release date: '%s'", releaseDate)
	}

	return enriched, nil
}
