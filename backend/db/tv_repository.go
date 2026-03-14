package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateTVShow creates a new TV show in the database
func (r *Neo4jRepository) CreateTVShow(ctx context.Context, input model.CreateTVShowInput) (*model.TVShow, error) {
	// Try to enrich with metadata if minimal data provided
	var meta *metadata.VideoMetadata
	if r.metadata != nil && shouldEnrichTVShow(input) {
		enrichedInput, m, err := r.enrichTVShowInput(input)
		if err != nil {
			// Log error but continue with original input
			// In production, use proper logging
		} else {
			input = enrichedInput
			meta = m
		}
	}

	// Check if TV show already exists
	var year *int
	if input.ReleaseDate != nil {
		y, err := strconv.Atoi(*input.ReleaseDate)
		if err == nil {
			year = &y
		}
	}
	existing, err := r.FindMediaByTitleTypeYear(ctx, input.Title, "TVShow", year)
	if err == nil && existing != nil {
		// Exists, check if need to update searchDepth
		inputDepth := int32(0)
		if input.SearchDepth != nil {
			inputDepth = *input.SearchDepth
		}
		if existing.GetSearchDepth() > inputDepth {
			// Update to lower depth
			err = r.UpdateMediaSearchDepth(ctx, existing.GetID(), inputDepth)
			if err != nil {
				return nil, err
			}
			// Re-fetch to get updated
			return r.GetTVShowByID(ctx, existing.GetID())
		}
		// Return existing
		if tvShow, ok := existing.(*model.TVShow); ok {
			return tvShow, nil
		}
		return nil, fmt.Errorf("existing media is not a TV show")
	}

	tvShowID := uuid.New()

	// Prepare cast/crew/productionCompanies/genres similar to CreateMovie and use metadata when available
	castData := make([]map[string]any, 0)
	if meta != nil && len(meta.CastCredits) > 0 {
		castData = make([]map[string]any, len(meta.CastCredits))
		for i, c := range meta.CastCredits {
			castData[i] = map[string]any{
				"name":       c.Name,
				"id":         uuid.New().String(),
				"externalID": c.PersonID,
				"character":  c.Character,
				"order":      c.Order,
			}
		}
	} else if len(input.Cast) > 0 {
		castData = make([]map[string]any, len(input.Cast))
		for i, name := range input.Cast {
			castData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	crewData := make([]map[string]any, 0)
	if meta != nil && len(meta.CrewCredits) > 0 {
		crewData = make([]map[string]any, len(meta.CrewCredits))
		for i, c := range meta.CrewCredits {
			crewData[i] = map[string]any{
				"name":       c.Name,
				"id":         uuid.New().String(),
				"externalID": c.PersonID,
				"job":        c.Job,
				"department": c.Department,
			}
		}
	} else if len(input.Crew) > 0 {
		crewData = make([]map[string]any, len(input.Crew))
		for i, name := range input.Crew {
			crewData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	pcData := make([]map[string]any, len(input.ProductionCompanies))
	for i, name := range input.ProductionCompanies {
		pcData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	if len(pcData) == 0 && meta != nil && len(meta.ProductionCompanies) > 0 {
		pcData = make([]map[string]any, len(meta.ProductionCompanies))
		for i, name := range meta.ProductionCompanies {
			pcData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	genreData := make([]map[string]any, len(input.Genres))
	for i, name := range input.Genres {
		genreData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	if len(genreData) == 0 && meta != nil && len(meta.Genres) > 0 {
		genreData = make([]map[string]any, len(meta.Genres))
		for i, name := range meta.Genres {
			genreData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	// Build production country data from metadata (TV input has no productionCountries field)
	productionCountryData := []map[string]any{}
	if meta != nil && len(meta.ProductionCountries) > 0 {
		productionCountryData = make([]map[string]any, len(meta.ProductionCountries))
		for i, name := range meta.ProductionCountries {
			productionCountryData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	// Handle searchDepth
	searchDepth := int32(0)
	if input.SearchDepth != nil {
		searchDepth = *input.SearchDepth
	}

	// Add normalizedName
	for i := range castData {
		if n, ok := castData[i]["name"].(string); ok {
			castData[i]["normalizedName"] = NormalizeName(n)
		}
	}
	for i := range crewData {
		if n, ok := crewData[i]["name"].(string); ok {
			crewData[i]["normalizedName"] = NormalizeName(n)
		}
	}
	for i := range pcData {
		if n, ok := pcData[i]["name"].(string); ok {
			pcData[i]["normalizedName"] = NormalizeName(n)
		}
	}
	for i := range genreData {
		if n, ok := genreData[i]["name"].(string); ok {
			genreData[i]["normalizedName"] = NormalizeName(n)
		}
	}

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
						CREATE (t:TVShow:Media {
								id: $id,
								title: $title,
								releaseDate: $releaseDate,
								description: $description,
								coverUrl: $coverUrl,
								seasons: $seasons,
								episodes: $episodes,
								status: $status,
								searchDepth: $searchDepth,
								createdAt: datetime(),
								updatedAt: datetime()
							})
						WITH t
						FOREACH (castData IN $cast |
								MERGE (p:Person {normalizedName: castData.normalizedName})
								ON CREATE SET p.id = castData.id, p.name = castData.name, p.createdAt = datetime()
								SET p.externalID = castData.externalID
								MERGE (p)-[r:ACTED_IN]->(t)
								ON CREATE SET r.character = castData.character, r.order = castData.order
						)
						WITH t
						FOREACH (crewData IN $crew |
								MERGE (p:Person {normalizedName: crewData.normalizedName})
								ON CREATE SET p.id = crewData.id, p.name = crewData.name, p.createdAt = datetime()
								SET p.externalID = crewData.externalID
								MERGE (p)-[r:CREW_ON]->(t)
								ON CREATE SET r.job = crewData.job, r.department = crewData.department
						)

						WITH t
						FOREACH (pcData IN $productionCompanies |
								MERGE (pc:ProductionCompany {normalizedName: pcData.normalizedName})
								ON CREATE SET pc.id = pcData.id, pc.name = pcData.name, pc.createdAt = datetime()
								MERGE (pc)-[:PRODUCED]->(t)
						)
						WITH t
						FOREACH (genreData IN $genres |
								MERGE (g:Tag {type: 'genre', normalizedName: genreData.normalizedName})
								ON CREATE SET g.id = genreData.id, g.name = genreData.name, g.createdAt = datetime()
								MERGE (g)-[:TAGGED]->(t)
						)
						WITH t
						FOREACH (pcountryData IN $productionCountries |
								MERGE (pcountry:ProductionCountry {id: pcountryData.id})
								ON CREATE SET pcountry.name = pcountryData.name
								MERGE (pcountry)-[:PRODUCED_IN]->(t)
						)
						RETURN t.id
						`

		params := map[string]any{
			"id":                  tvShowID.String(),
			"title":               input.Title,
			"releaseDate":         input.ReleaseDate,
			"description":         input.Description,
			"coverUrl":            input.CoverURL,
			"seasons":             input.Seasons,
			"episodes":            input.Episodes,
			"status":              input.Status,
			"searchDepth":         searchDepth,
			"cast":                castData,
			"crew":                crewData,
			"productionCompanies": pcData,
			"genres":              genreData,
			"productionCountries": productionCountryData,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			idStr, _ := result.Record().Get("t.id")
			if s, ok := idStr.(string); ok {
				parsed, err := uuid.Parse(s)
				if err == nil {
					return parsed.String(), nil
				}
			}
			return nil, fmt.Errorf("failed to parse created TV show id")
		}

		return nil, fmt.Errorf("failed to create TV show")
	})
	if err != nil {
		return nil, err
	}

	idStr, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected create TV show result type: %T", result)
	}
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created TV show id: %w", err)
	}
	return r.GetTVShowByID(ctx, parsedID)
}

// GetTVShowByID retrieves a TV show by its ID
func (r *Neo4jRepository) GetTVShowByID(ctx context.Context, id uuid.UUID) (*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
						MATCH (t:TVShow {id: $id})
						OPTIONAL MATCH (t)<-[ract:ACTED_IN]-(cast:Person)
						OPTIONAL MATCH (t)<-[rcrew:CREW_ON]-(crew:Person)
						OPTIONAL MATCH (t)<-[:PRODUCED]-(pc:ProductionCompany)
						OPTIONAL MATCH (g:Tag)-[:TAGGED]->(t) WHERE g.type = 'genre'
						OPTIONAL MATCH (t)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
						RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
								   t.description as description, t.coverUrl as coverUrl,
					   t.seasons as seasons, t.episodes as episodes, t.status as status,
					   t.searchDepth as searchDepth,
					   collect(DISTINCT cast {.*, externalID: cast.externalID}) as cast,
								   collect(DISTINCT crew {.*, externalID: crew.externalID}) as crew,
								   collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
								   collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits,
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
			castParsed := parsePersons(record.AsMap()["cast"])
			crewParsed := parsePersons(record.AsMap()["crew"])
			castCredits := parsePersonCredits(record.AsMap()["castCredits"])
			crewCredits := parseCrewCredits(record.AsMap()["crewCredits"])
			vShow := &model.TVShow{
				ID:                  id,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:             getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:            getInt32Pointer(record.AsMap()["episodes"]),
				Status:              getStringPointer(record.AsMap()["status"]),
				SearchDepth:         getInt32Value(record.AsMap()["searchDepth"]),
				Cast:                castParsed,
				Crew:                crewParsed,
				CastCredits:         castCredits,
				CrewCredits:         crewCredits,
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
					vShow.ProductionCountries = append(vShow.ProductionCountries, &model.ProductionCountry{ID: pc.ID, Name: pc.Name})
				}
			}
			return vShow, nil
		}

		return nil, fmt.Errorf("TV show not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.TVShow), nil
}

// GetAllTVShows retrieves all TV shows with optional pagination.
func (r *Neo4jRepository) GetAllTVShows(ctx context.Context, limit, offset *int) ([]*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Paginate on bare TVShow nodes first, then expand joins.
		query := `
						MATCH (t:TVShow)
						WITH t ORDER BY t.title
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
						OPTIONAL MATCH (t)<-[ract:ACTED_IN]-(cast:Person)
						OPTIONAL MATCH (t)<-[rcrew:CREW_ON]-(crew:Person)
						OPTIONAL MATCH (t)<-[:PRODUCED]-(pc:ProductionCompany)
						OPTIONAL MATCH (g:Tag)-[:TAGGED]->(t) WHERE g.type = 'genre'
						OPTIONAL MATCH (t)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
						RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
								   t.description as description, t.coverUrl as coverUrl,
					   t.seasons as seasons, t.episodes as episodes, t.status as status,
					   t.searchDepth as searchDepth,
					   collect(DISTINCT cast {.*, externalID: cast.externalID}) as cast,
								   collect(DISTINCT crew {.*, externalID: crew.externalID}) as crew,
								   collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
								   collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits,
								   collect(DISTINCT pc) as productionCompanies,
								   collect(DISTINCT g) as genres,
								   collect(DISTINCT pcountry) as productionCountries
						ORDER BY t.title
						`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var tvShows []*model.TVShow
		for result.Next(ctx) {
			record := result.Record()
			vShowID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			castParsed := parsePersons(record.AsMap()["cast"])
			crewParsed := parsePersons(record.AsMap()["crew"])
			castCredits := parsePersonCredits(record.AsMap()["castCredits"])
			crewCredits := parseCrewCredits(record.AsMap()["crewCredits"])
			vShow := &model.TVShow{
				ID:                  vShowID,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:             getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:            getInt32Pointer(record.AsMap()["episodes"]),
				Status:              getStringPointer(record.AsMap()["status"]),
				SearchDepth:         getInt32Value(record.AsMap()["searchDepth"]),
				Cast:                castParsed,
				Crew:                crewParsed,
				CastCredits:         castCredits,
				CrewCredits:         crewCredits,
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
					vShow.ProductionCountries = append(vShow.ProductionCountries, &model.ProductionCountry{ID: pc.ID, Name: pc.Name})
				}
			}
			tvShows = append(tvShows, vShow)
		}

		return tvShows, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.TVShow), nil
}

// shouldEnrichTVShow determines if a TV show input should be enriched with metadata
func shouldEnrichTVShow(input model.CreateTVShowInput) bool {
	return len(input.Cast) == 0 || input.Description == nil || input.CoverURL == nil || input.Seasons == nil
}

// enrichTVShowInput fetches metadata and merges it with the input
func (r *Neo4jRepository) enrichTVShowInput(input model.CreateTVShowInput) (model.CreateTVShowInput, *metadata.VideoMetadata, error) {
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
		Type:        metadata.MediaTypeTV,
		Title:       input.Title,
		ReleaseYear: year,
		ID:          safeStringValue(input.ExternalID),
	})

	if err != nil {
		return input, nil, err
	}

	meta, ok := metaInterface.(*metadata.VideoMetadata)
	if !ok {
		return input, nil, fmt.Errorf("unexpected metadata type for TV show")
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

	// Populate new fields from metadata if not provided. Prefer structured credits when available.
	if len(enriched.Cast) == 0 {
		if len(meta.CastCredits) > 0 {
			names := make([]string, len(meta.CastCredits))
			for i, c := range meta.CastCredits {
				names[i] = c.Name
			}
			enriched.Cast = names
		} else if len(meta.Cast) > 0 {
			enriched.Cast = meta.Cast
		}
	}

	if len(enriched.Crew) == 0 {
		if len(meta.CrewCredits) > 0 {
			names := make([]string, len(meta.CrewCredits))
			for i, c := range meta.CrewCredits {
				names[i] = c.Name
			}
			enriched.Crew = names
		} else if len(meta.Crew) > 0 {
			enriched.Crew = meta.Crew
		}
	}

	if len(enriched.Genres) == 0 && len(meta.Genres) > 0 {
		enriched.Genres = meta.Genres
	}

	if len(enriched.ProductionCompanies) == 0 && len(meta.ProductionCompanies) > 0 {
		enriched.ProductionCompanies = meta.ProductionCompanies
	}

	return enriched, meta, nil
}
