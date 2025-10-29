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

// CreateMovie creates a new movie in the database
func (r *Neo4jRepository) CreateMovie(ctx context.Context, input model.CreateMovieInput) (*model.Movie, error) {
	log.Printf("CreateMovie called: title=%s, metadata_present=%v", input.Title, r.metadata != nil)

	// Try to enrich with metadata if minimal data provided
	var meta *metadata.VideoMetadata
	if r.metadata != nil && shouldEnrichMovie(input) {
		enrichedInput, m, err := r.enrichMovieInput(input)
		if err != nil {
			// Log error but continue with original input
			// In production, use proper logging
		} else {
			input = enrichedInput
			meta = m
		}
	}

	movieUUID := uuid.New()

	// Prepare data for nodes
	castData := make([]map[string]any, len(input.Cast))
	for i, name := range input.Cast {
		castData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	crewData := make([]map[string]any, len(input.Crew))
	for i, name := range input.Crew {
		crewData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	pcData := make([]map[string]any, len(input.ProductionCompanies))
	for i, name := range input.ProductionCompanies {
		pcData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	genreData := make([]map[string]any, len(input.Genres))
	for i, name := range input.Genres {
		genreData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}

	// Use metadata returned from enrichMovieInput (if any) for extra fields
	var ratingVal any = nil
	var urlVal any = nil
	productionCountryData := []map[string]any{}
	if meta != nil {
		if meta.Rating > 0 {
			ratingVal = float64(meta.Rating)
		}
		if meta.URL != "" {
			urlVal = meta.URL
		}
		if len(meta.ProductionCountries) > 0 {
			productionCountryData = make([]map[string]any, len(meta.ProductionCountries))
			for i, name := range meta.ProductionCountries {
				productionCountryData[i] = map[string]any{"name": name, "id": uuid.New().String()}
			}
		}
	}

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
						CREATE (m:Movie:Media {
							id: $id,
							title: $title,
							releaseDate: $releaseDate,
							description: $description,
							coverUrl: $coverUrl,
							runtime: $runtime,
							budget: $budget,
							boxOffice: $boxOffice,
							rating: $rating,
							url: $url
						})
						WITH m
						FOREACH (castData IN $cast |
							MERGE (p:Person {id: castData.id})
							ON CREATE SET p.name = castData.name
							MERGE (p)-[:ACTED_IN]->(m)
						)
						WITH m
						FOREACH (crewData IN $crew |
							MERGE (p:Person {id: crewData.id})
							ON CREATE SET p.name = crewData.name
							MERGE (p)-[:CREW_ON]->(m)
						)
						WITH m
						FOREACH (pcData IN $productionCompanies |
							MERGE (pc:ProductionCompany {id: pcData.id})
							ON CREATE SET pc.name = pcData.name
							MERGE (pc)-[:PRODUCED]->(m)
						)
						WITH m
			FOREACH (genreData IN $genres |
				MERGE (g:Genre {id: genreData.id})
				ON CREATE SET g.name = genreData.name
				MERGE (g)-[:HAS_MOVIE]->(m)
			)

						WITH m
						FOREACH (pcountryData IN $productionCountries |
							MERGE (pcountry:ProductionCountry {id: pcountryData.id})
							ON CREATE SET pcountry.name = pcountryData.name
							MERGE (pcountry)-[:PRODUCED_IN]->(m)
						)
						RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
							   m.description as description, m.coverUrl as coverUrl,
							   m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
							   m.rating as rating, m.url as url
						`

		params := map[string]any{
			"id":                  movieUUID.String(),
			"title":               input.Title,
			"releaseDate":         input.ReleaseDate,
			"description":         input.Description,
			"coverUrl":            input.CoverURL,
			"runtime":             input.Runtime,
			"budget":              input.Budget,
			"boxOffice":           input.BoxOffice,
			"rating":              ratingVal,
			"url":                 urlVal,
			"cast":                castData,
			"crew":                crewData,
			"productionCompanies": pcData,
			"genres":              genreData,
			"productionCountries": productionCountryData,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, fmt.Errorf("failed to run create movie query: %w", err)
		}

		if result.Next(ctx) {
			record := result.Record()
			movie := &model.Movie{
				ID:                  movieUUID,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
				Cast:                []*model.Person{}, // Will be populated in GetMovieByID
				Crew:                []*model.Person{},
				ProductionCompanies: []*model.ProductionCompany{},
				Genres:              []*model.Genre{},
				ProductionCountries: []*model.ProductionCountry{},
				Creators:            []*model.Creator{},
				Platforms:           []*model.Platform{},
				Tags:                []*model.Tag{},
				Ratings:             []*model.Rating{},
				AverageRating:       nil,
			}
			return movie, nil
		}

		if err := result.Err(); err != nil {
			return nil, fmt.Errorf("create movie result error: %w", err)
		}

		return nil, fmt.Errorf("create movie returned no records")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Movie), nil
}

// GetMovieByID retrieves a movie by its ID
func (r *Neo4jRepository) GetMovieByID(ctx context.Context, id uuid.UUID) (*model.Movie, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
					MATCH (m:Movie {id: $id})
					OPTIONAL MATCH (m)<-[:ACTED_IN]-(cast:Person)
					OPTIONAL MATCH (m)<-[:CREW_ON]-(crew:Person)
					OPTIONAL MATCH (m)<-[:PRODUCED]-(pc:ProductionCompany)
					OPTIONAL MATCH (g:Genre)-[:HAS_MOVIE]->(m)
					OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
					RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
					       m.description as description, m.coverUrl as coverUrl,
					       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
					       collect(DISTINCT cast) as cast,
					       collect(DISTINCT crew) as crew,
					       collect(DISTINCT pc) as productionCompanies,
					       collect(DISTINCT g) as genres,
					       collect(DISTINCT pcountry) as productionCountries
				`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			movie := &model.Movie{
				ID:                  id,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
				Cast:                parsePersons(record.AsMap()["cast"]),
				Crew:                parsePersons(record.AsMap()["crew"]),
				ProductionCompanies: parseProductionCompanies(record.AsMap()["productionCompanies"]),
				Genres:              parseGenres(record.AsMap()["genres"]),
				ProductionCountries: []*model.ProductionCountry{},
				Creators:            []*model.Creator{},
				Platforms:           []*model.Platform{},
				Tags:                []*model.Tag{},
				Ratings:             []*model.Rating{},
				AverageRating:       nil,
			}
			if pcs := parseProductionCountries(record.AsMap()["productionCountries"]); len(pcs) > 0 {
				for _, pc := range pcs {
					movie.ProductionCountries = append(movie.ProductionCountries, &model.ProductionCountry{ID: pc.ID, Name: pc.Name})
				}
			}
			return movie, nil
		}

		return nil, fmt.Errorf("movie not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.Movie), nil
}

// GetAllMovies retrieves all movies
func (r *Neo4jRepository) GetAllMovies(ctx context.Context) ([]*model.Movie, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (m:Movie)
			OPTIONAL MATCH (m)<-[:ACTED_IN]-(cast:Person)
			OPTIONAL MATCH (m)<-[:CREW_ON]-(crew:Person)
			OPTIONAL MATCH (m)<-[:PRODUCED]-(pc:ProductionCompany)
			OPTIONAL MATCH (g:Genre)-[:HAS_MOVIE]->(m)
			OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
			RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
			       m.description as description, m.coverUrl as coverUrl,
			       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
			       collect(DISTINCT cast) as cast,
			       collect(DISTINCT crew) as crew,
			       collect(DISTINCT pc) as productionCompanies,
			       collect(DISTINCT g) as genres,
			       collect(DISTINCT pcountry) as productionCountries
			ORDER BY m.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var movies []*model.Movie
		for result.Next(ctx) {
			record := result.Record()
			movieID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			movie := &model.Movie{
				ID:                  movieID,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
				Cast:                parsePersons(record.AsMap()["cast"]),
				Crew:                parsePersons(record.AsMap()["crew"]),
				ProductionCompanies: parseProductionCompanies(record.AsMap()["productionCompanies"]),
				Genres:              parseGenres(record.AsMap()["genres"]),
				ProductionCountries: []*model.ProductionCountry{},
				Creators:            []*model.Creator{},
				Platforms:           []*model.Platform{},
				Tags:                []*model.Tag{},
				Ratings:             []*model.Rating{},
				AverageRating:       nil,
			}
			if pcs := parseProductionCountries(record.AsMap()["productionCountries"]); len(pcs) > 0 {
				for _, pc := range pcs {
					movie.ProductionCountries = append(movie.ProductionCountries, &model.ProductionCountry{ID: pc.ID, Name: pc.Name})
				}
			}
			movies = append(movies, movie)
		}

		return movies, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.Movie), nil
}

// shouldEnrichMovie determines if a movie input should be enriched with metadata
func shouldEnrichMovie(input model.CreateMovieInput) bool {
	return input.Description == nil && input.CoverURL == nil && input.Runtime == nil
}

// enrichMovieInput fetches metadata and merges it with the input
func (r *Neo4jRepository) enrichMovieInput(input model.CreateMovieInput) (model.CreateMovieInput, *metadata.VideoMetadata, error) {
	if r.metadata == nil {
		return input, nil, nil
	}

	// Determine year from release date if available
	var year int
	if input.ReleaseDate != nil {
		// Simple year extraction - in production, use proper date parsing
		if len(*input.ReleaseDate) >= 4 {
			fmt.Sscanf(*input.ReleaseDate, "%d", &year)
		}
	}

	metaInterface, err := r.metadata.GetMetadata(metadata.MediaInfo{
		Type:        metadata.MediaTypeMovie,
		Title:       input.Title,
		ReleaseYear: year,
	})

	if err != nil {
		return input, nil, err
	}

	meta, ok := metaInterface.(*metadata.VideoMetadata)
	if !ok {
		return input, nil, fmt.Errorf("unexpected metadata type for movie")
	}

	// Merge metadata with input (input takes precedence)
	enriched := input

	if enriched.Description == nil && meta.Description != "" {
		enriched.Description = &meta.Description
	}

	if enriched.CoverURL == nil && meta.ImageURL != "" {
		enriched.CoverURL = &meta.ImageURL
	}

	if enriched.ReleaseDate == nil && meta.ReleaseYear > 0 {
		releaseDate := fmt.Sprintf("%d", meta.ReleaseYear)
		enriched.ReleaseDate = &releaseDate
	}

	if enriched.Runtime == nil && meta.Runtime > 0 {
		runtime := int32(meta.Runtime)
		enriched.Runtime = &runtime
	}

	if enriched.Budget == nil && meta.Budget > 0 {
		budget := int32(meta.Budget)
		enriched.Budget = &budget
	}

	if enriched.BoxOffice == nil && meta.BoxOffice > 0 {
		boxOffice := int32(meta.BoxOffice)
		enriched.BoxOffice = &boxOffice
	}

	// Populate new fields from metadata if not provided
	if len(enriched.Cast) == 0 && len(meta.Cast) > 0 {
		enriched.Cast = meta.Cast
	}

	if len(enriched.Crew) == 0 && len(meta.Crew) > 0 {
		enriched.Crew = meta.Crew
	}

	if len(enriched.Genres) == 0 && len(meta.Genres) > 0 {
		enriched.Genres = meta.Genres
	}

	if len(enriched.ProductionCompanies) == 0 && len(meta.ProductionCompanies) > 0 {
		enriched.ProductionCompanies = meta.ProductionCompanies
	}

	return enriched, meta, nil
}

// parsePersons converts a value returned by collect(DISTINCT person) into []*model.Person
func parsePersons(value interface{}) []*model.Person {
	if value == nil {
		return []*model.Person{}
	}

	var persons []*model.Person

	arr, ok := value.([]any)
	if !ok {
		return persons
	}

	// Generic interface for types that provide Props() map
	type propsProvider interface{ Props() map[string]any }

	for _, item := range arr {
		if nodeMap, ok := item.(map[string]any); ok {
			persons = append(persons, personFromNodeMap(nodeMap))
			continue
		}
		if p, ok := item.(propsProvider); ok {
			persons = append(persons, personFromNodeMap(p.Props()))
			continue
		}
	}

	return persons
}

// personFromNodeMap converts a node map into a Person
func personFromNodeMap(nodeMap map[string]any) *model.Person {
	var id uuid.UUID
	if s, ok := nodeMap["id"].(string); ok && s != "" {
		if parsed, err := uuid.Parse(s); err == nil {
			id = parsed
		}
	}

	name := ""
	if n, ok := nodeMap["name"].(string); ok {
		name = n
	}

	return &model.Person{
		ID:      id,
		Name:    name,
		ActedIn: []*model.Movie{}, // Not populated here
		CrewOn:  []*model.Movie{},
	}
}

// parseProductionCompanies converts a value returned by collect(DISTINCT pc) into []*model.ProductionCompany
func parseProductionCompanies(value interface{}) []*model.ProductionCompany {
	if value == nil {
		return []*model.ProductionCompany{}
	}

	var pcs []*model.ProductionCompany

	arr, ok := value.([]any)
	if !ok {
		return pcs
	}

	type propsProvider interface{ Props() map[string]any }

	for _, item := range arr {
		if nodeMap, ok := item.(map[string]any); ok {
			pcs = append(pcs, productionCompanyFromNodeMap(nodeMap))
			continue
		}
		if p, ok := item.(propsProvider); ok {
			pcs = append(pcs, productionCompanyFromNodeMap(p.Props()))
			continue
		}
	}

	return pcs
}

// productionCompanyFromNodeMap converts a node map into a ProductionCompany
func productionCompanyFromNodeMap(nodeMap map[string]any) *model.ProductionCompany {
	var id uuid.UUID
	if s, ok := nodeMap["id"].(string); ok && s != "" {
		if parsed, err := uuid.Parse(s); err == nil {
			id = parsed
		}
	}

	name := ""
	if n, ok := nodeMap["name"].(string); ok {
		name = n
	}

	return &model.ProductionCompany{
		ID:       id,
		Name:     name,
		Produced: []*model.Movie{}, // Not populated here
	}
}

// parseGenres converts a value returned by collect(DISTINCT g) into []*model.Genre
func parseGenres(value interface{}) []*model.Genre {
	if value == nil {
		return []*model.Genre{}
	}

	var genres []*model.Genre

	arr, ok := value.([]any)
	if !ok {
		return genres
	}

	type propsProvider interface{ Props() map[string]any }

	for _, item := range arr {
		if nodeMap, ok := item.(map[string]any); ok {
			genres = append(genres, genreFromNodeMap(nodeMap))
			continue
		}
		if p, ok := item.(propsProvider); ok {
			genres = append(genres, genreFromNodeMap(p.Props()))
			continue
		}
	}

	return genres
}

// genreFromNodeMap converts a node map into a Genre
func genreFromNodeMap(nodeMap map[string]any) *model.Genre {
	var id uuid.UUID
	if s, ok := nodeMap["id"].(string); ok && s != "" {
		if parsed, err := uuid.Parse(s); err == nil {
			id = parsed
		}
	}

	name := ""
	if n, ok := nodeMap["name"].(string); ok {
		name = n
	}

	return &model.Genre{
		ID:     id,
		Name:   name,
		Movies: []*model.Movie{}, // Not populated here
	}
}

// parseProductionCountries converts a value returned by collect(DISTINCT pcountry) into []*model.ProductionCountry
func parseProductionCountries(value interface{}) []*model.ProductionCountry {
	if value == nil {
		return []*model.ProductionCountry{}
	}

	var pcs []*model.ProductionCountry

	arr, ok := value.([]any)
	if !ok {
		return pcs
	}

	type propsProvider interface{ Props() map[string]any }

	for _, item := range arr {
		if nodeMap, ok := item.(map[string]any); ok {
			var id uuid.UUID
			if s, ok := nodeMap["id"].(string); ok && s != "" {
				if parsed, err := uuid.Parse(s); err == nil {
					id = parsed
				}
			}
			name := ""
			if n, ok := nodeMap["name"].(string); ok {
				name = n
			}
			pcs = append(pcs, &model.ProductionCountry{ID: id, Name: name})
			continue
		}
		if p, ok := item.(propsProvider); ok {
			nodeMap := p.Props()
			var id uuid.UUID
			if s, ok := nodeMap["id"].(string); ok && s != "" {
				if parsed, err := uuid.Parse(s); err == nil {
					id = parsed
				}
			}
			name := ""
			if n, ok := nodeMap["name"].(string); ok {
				name = n
			}
			pcs = append(pcs, &model.ProductionCountry{ID: id, Name: name})
			continue
		}
	}

	return pcs
}
