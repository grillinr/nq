package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// parseCreatorsFromNeo4j safely parses creators from Neo4j results, skipping invalid entries
func parseCreatorsFromNeo4j(value interface{}) []*model.Creator {
	creators := []*model.Creator{}
	if value == nil {
		return creators
	}

	switch v := value.(type) {
	case []any:
		for _, e := range v {
			if m, ok := e.(map[string]any); ok {
				// Skip entries with null/empty id or name to avoid GraphQL schema violations
				idStr, hasID := m["id"].(string)
				name, hasName := m["name"].(string)
				if !hasID || idStr == "" || !hasName || name == "" {
					continue
				}

				aid, err := uuid.Parse(idStr)
				if err != nil {
					continue // Skip invalid UUIDs
				}

				creators = append(creators, &model.Creator{ID: aid, Name: name})
			}
		}
	}
	return creators
}

// parseTagsFromNeo4j safely parses tags from Neo4j results, skipping invalid entries
func parseTagsFromNeo4j(value interface{}) []*model.Tag {
	tags := []*model.Tag{}
	if value == nil {
		return tags
	}

	switch v := value.(type) {
	case []any:
		for _, e := range v {
			if m, ok := e.(map[string]any); ok {
				// Skip entries with null/empty fields to avoid GraphQL schema violations
				idStr, hasID := m["id"].(string)
				name, hasName := m["name"].(string)
				typeStr, hasType := m["type"].(string)
				if !hasID || idStr == "" || !hasName || name == "" || !hasType || typeStr == "" {
					continue
				}

				sid, err := uuid.Parse(idStr)
				if err != nil {
					continue // Skip invalid UUIDs
				}

				tags = append(tags, &model.Tag{ID: sid, Name: name, Type: typeStr})
			}
		}
	}
	return tags
}

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

	// Check if book already exists (ISBN or title+author+year)
	var year *int
	if input.ReleaseDate != nil {
		if y, err := strconv.Atoi(*input.ReleaseDate); err == nil {
			year = &y
		}
	}
	if existing, err := r.findExistingBook(ctx, input.Title, input.Authors, input.Isbn, year); err == nil && existing != nil {
		inputDepth := int32(0)
		if input.SearchDepth != nil {
			inputDepth = *input.SearchDepth
		}
		if existing.GetSearchDepth() > inputDepth {
			if err := r.UpdateMediaSearchDepth(ctx, existing.GetID(), inputDepth); err != nil {
				return nil, err
			}
			return r.GetBookByID(ctx, existing.GetID())
		}
		if book, ok := existing.(*model.Book); ok {
			return book, nil
		}
		return nil, fmt.Errorf("existing media is not a book")
	}

	bookID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Handle searchDepth
		searchDepth := int32(0)
		if input.SearchDepth != nil {
			searchDepth = *input.SearchDepth
		}

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
						searchDepth: $searchDepth,
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
			"searchDepth": searchDepth,
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
				SearchDepth:   searchDepth,
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

func (r *Neo4jRepository) findExistingBook(ctx context.Context, title string, authors []string, isbn *string, year *int) (model.Media, error) {
	if isbn != nil && *isbn != "" {
		if media, err := r.GetBookByISBN(ctx, *isbn); err == nil && media != nil {
			return media, nil
		}
	}
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	primaryAuthor := ""
	if len(authors) > 0 {
		primaryAuthor = NormalizeName(authors[0])
	}
	return r.GetBookByTitleAuthorYear(ctx, title, primaryAuthor, year)
}

// GetBookByID retrieves a book by its ID
func (r *Neo4jRepository) GetBookByID(ctx context.Context, id uuid.UUID) (*model.Book, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
				MATCH (b:Book:Media {id: $id})
				OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
				OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
				RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
					       b.description as description, b.coverUrl as coverUrl,
					       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
					       b.searchDepth as searchDepth,
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
			authors := parseCreatorsFromNeo4j(record.AsMap()["authors"])

			// Parse subjects
			subjects := parseTagsFromNeo4j(record.AsMap()["subjects"])

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
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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

// GetBookByISBN retrieves a book by ISBN (if stored)
func (r *Neo4jRepository) GetBookByISBN(ctx context.Context, isbn string) (*model.Book, error) {
	if isbn == "" {
		return nil, fmt.Errorf("isbn required")
	}
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (b:Book:Media)
			WHERE b.isbn = $isbn
			RETURN b.id as id
			LIMIT 1
		`
		params := map[string]any{"isbn": isbn}
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if res.Next(ctx) {
			if idStr, ok := res.Record().Get("id"); ok {
				if s, ok := idStr.(string); ok {
					if id, err := uuid.Parse(s); err == nil {
						return r.GetBookByID(ctx, id)
					}
				}
			}
		}
		return nil, fmt.Errorf("book not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*model.Book), nil
}

// GetBookByTitleAuthorYear attempts to find a book by normalized title, primary author, and year.
// It fetches the full book in a single query — no second round-trip.
// When normalizedAuthor is provided the query JOINs on the author node directly rather than
// collecting all authors and filtering in-memory, which avoids a full-scan anti-pattern.
func (r *Neo4jRepository) GetBookByTitleAuthorYear(ctx context.Context, title string, normalizedAuthor string, year *int) (*model.Book, error) {
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	queryTitle := strings.ToLower(NormalizeName(title))
	yearStr := ""
	if year != nil {
		yearStr = fmt.Sprintf("%d", *year)
	}

	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// When an author is known, join directly on the Person node so Neo4j can use the
		// person_normalized_index and avoid the collect-then-filter anti-pattern.
		var query string
		if normalizedAuthor != "" {
			query = `
				MATCH (b:Book:Media)
				WHERE ($title = "" OR toLower(b.title) CONTAINS $title)
				  AND ($yearStr = "" OR b.releaseDate STARTS WITH $yearStr)
				MATCH (author:Person {normalizedName: $author})-[:AUTHORED]->(b)
				WITH b LIMIT 1
				OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
				OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
				RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
				       b.description as description, b.coverUrl as coverUrl,
				       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
				       b.searchDepth as searchDepth,
				       collect(DISTINCT {id: person.id, name: person.name}) as authors,
				       collect(DISTINCT {id: t.id, name: t.name, type: t.type}) as subjects,
				       collect(DISTINCT pnode.name) as publisher_nodes,
				       b.publishers as publishers
			`
		} else {
			query = `
				MATCH (b:Book:Media)
				WHERE ($title = "" OR toLower(b.title) CONTAINS $title)
				  AND ($yearStr = "" OR b.releaseDate STARTS WITH $yearStr)
				WITH b LIMIT 1
				OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
				OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
				RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
				       b.description as description, b.coverUrl as coverUrl,
				       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
				       b.searchDepth as searchDepth,
				       collect(DISTINCT {id: person.id, name: person.name}) as authors,
				       collect(DISTINCT {id: t.id, name: t.name, type: t.type}) as subjects,
				       collect(DISTINCT pnode.name) as publisher_nodes,
				       b.publishers as publishers
			`
		}

		params := map[string]any{
			"title":   queryTitle,
			"author":  normalizedAuthor,
			"yearStr": yearStr,
		}
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if res.Next(ctx) {
			record := res.Record()

			idRaw, ok := record.Get("id")
			if !ok {
				return nil, fmt.Errorf("book not found")
			}
			idStr, ok := idRaw.(string)
			if !ok {
				return nil, fmt.Errorf("book not found")
			}
			bookID, err := uuid.Parse(idStr)
			if err != nil {
				return nil, err
			}

			authors := parseCreatorsFromNeo4j(record.AsMap()["authors"])
			subjects := parseTagsFromNeo4j(record.AsMap()["subjects"])

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
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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

// GetAllBooks retrieves all books with optional pagination.
func (r *Neo4jRepository) GetAllBooks(ctx context.Context, limit, offset *int) ([]*model.Book, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Paginate on bare Book nodes first, then expand joins — avoids row explosion
		// from cartesian product of authors × subjects × publishers before aggregation.
		query := `
			MATCH (b:Book:Media)
			WITH b ORDER BY b.title
			`

		params := map[string]any{}
		if offset != nil {
			query += " SKIP $offset"
			params["offset"] = *offset
		}
		if limit != nil {
			query += " LIMIT $limit"
			params["limit"] = *limit
		}

		query += `
			OPTIONAL MATCH (pnode:Publisher)-[:PUBLISHED]->(b)
			OPTIONAL MATCH (person:Person)-[:AUTHORED]->(b)
			OPTIONAL MATCH (t:Tag)-[:TAGGED]->(b) WHERE t.type = 'subject'
			RETURN b.id as id, b.title as title, b.releaseDate as releaseDate,
				       b.description as description, b.coverUrl as coverUrl,
				       b.pages as pages, b.isbn as isbn, b.publisher as publisher,
				       b.searchDepth as searchDepth,
				       collect(DISTINCT {id: person.id, name: person.name}) as authors,
				       collect(DISTINCT {id: t.id, name: t.name, type: t.type}) as subjects,
				       collect(DISTINCT pnode.name) as publisher_nodes,
				       b.publishers as publishers
				ORDER BY b.title
			`

		result, err := tx.Run(ctx, query, params)
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
			authors := parseCreatorsFromNeo4j(record.AsMap()["authors"])

			// Parse subjects
			subjects := parseTagsFromNeo4j(record.AsMap()["subjects"])

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
				SearchDepth:   getInt32Value(record.AsMap()["searchDepth"]),
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
