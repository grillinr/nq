package db

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateMovie creates a new movie in the database
func (r *Neo4jRepository) CreateMovie(ctx context.Context, input model.CreateMovieInput) (*model.Movie, error) {

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

	// Generate UUID upfront - will be used if creating new, ignored if matching existing
	movieUUID := uuid.New()

	// Prepare data for nodes: use structured credits from metadata if available, else input.
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

	// Build crewData similarly
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

	// production companies and genres unchanged
	pcData := make([]map[string]any, len(input.ProductionCompanies))
	for i, name := range input.ProductionCompanies {
		pcData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}
	genreData := make([]map[string]any, len(input.Genres))
	for i, name := range input.Genres {
		genreData[i] = map[string]any{"name": name, "id": uuid.New().String()}
	}

	// Add normalizedName for deduplication while preserving original names
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

	// Handle searchDepth
	searchDepth := int32(0)
	if input.SearchDepth != nil {
		searchDepth = *input.SearchDepth
	}

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Use MERGE to atomically check-and-create the movie node, preventing race conditions
		// Match on title + releaseDate to identify duplicates
		// Note: Neo4j doesn't allow null in MERGE keys, so we use empty string for null dates
		releaseDate := ""
		if input.ReleaseDate != nil {
			releaseDate = *input.ReleaseDate
		}

		query := `
		MERGE (m:Movie:Media {title: $title, releaseDate: $releaseDate})
		ON CREATE SET 
			m.id = $id,
			m.description = $description,
			m.coverUrl = $coverUrl,
			m.runtime = $runtime,
			m.budget = $budget,
			m.boxOffice = $boxOffice,
			m.rating = $rating,
			m.url = $url,
			m.searchDepth = $searchDepth,
			m.createdAt = datetime(),
			m.updatedAt = datetime()
		ON MATCH SET
			m.searchDepth = CASE 
				WHEN $searchDepth < m.searchDepth THEN $searchDepth 
				ELSE m.searchDepth 
			END,
			m.updatedAt = datetime(),
			m.description = CASE WHEN m.description IS NULL AND $description IS NOT NULL THEN $description ELSE m.description END,
			m.coverUrl = CASE WHEN m.coverUrl IS NULL AND $coverUrl IS NOT NULL THEN $coverUrl ELSE m.coverUrl END,
			m.runtime = CASE WHEN m.runtime IS NULL AND $runtime IS NOT NULL THEN $runtime ELSE m.runtime END,
			m.budget = CASE WHEN m.budget IS NULL AND $budget IS NOT NULL THEN $budget ELSE m.budget END,
			m.boxOffice = CASE WHEN m.boxOffice IS NULL AND $boxOffice IS NOT NULL THEN $boxOffice ELSE m.boxOffice END,
			m.rating = CASE WHEN m.rating IS NULL AND $rating IS NOT NULL THEN $rating ELSE m.rating END,
			m.url = CASE WHEN m.url IS NULL AND $url IS NOT NULL THEN $url ELSE m.url END
		WITH m
		FOREACH (castData IN $cast |
			MERGE (p:Person {normalizedName: castData.normalizedName})
			ON CREATE SET p.id = castData.id, p.name = castData.name, p.createdAt = datetime()
			ON MATCH SET p.externalID = CASE WHEN castData.externalID IS NOT NULL THEN castData.externalID ELSE p.externalID END
			MERGE (p)-[r:ACTED_IN]->(m)
			ON CREATE SET r.character = castData.character, r.order = castData.order
		)
		WITH m
		FOREACH (crewData IN $crew |
			MERGE (p:Person {normalizedName: crewData.normalizedName})
			ON CREATE SET p.id = crewData.id, p.name = crewData.name, p.createdAt = datetime()
			ON MATCH SET p.externalID = CASE WHEN crewData.externalID IS NOT NULL THEN crewData.externalID ELSE p.externalID END
			MERGE (p)-[r:CREW_ON]->(m)
			ON CREATE SET r.job = crewData.job, r.department = crewData.department
		)
		WITH m
		FOREACH (pcData IN $productionCompanies |
			MERGE (pc:ProductionCompany {normalizedName: pcData.normalizedName})
			ON CREATE SET pc.id = pcData.id, pc.name = pcData.name, pc.createdAt = datetime()
			MERGE (pc)-[:PRODUCED]->(m)
		)
		WITH m
		FOREACH (genreData IN $genres |
			MERGE (t:Tag {type: 'genre', normalizedName: genreData.normalizedName})
			ON CREATE SET t.id = genreData.id, t.name = genreData.name
			MERGE (t)-[:TAGGED]->(m)
		)
		WITH m
		FOREACH (pcountryData IN $productionCountries |
			MERGE (pcountry:ProductionCountry {id: pcountryData.id})
			ON CREATE SET pcountry.name = pcountryData.name
			MERGE (pcountry)-[:PRODUCED_IN]->(m)
		)
		RETURN m.id as movieId
		`

		params := map[string]any{
			"id":                  movieUUID.String(),
			"title":               input.Title,
			"releaseDate":         releaseDate,
			"description":         input.Description,
			"coverUrl":            input.CoverURL,
			"runtime":             input.Runtime,
			"budget":              input.Budget,
			"boxOffice":           input.BoxOffice,
			"rating":              ratingVal,
			"url":                 urlVal,
			"searchDepth":         searchDepth,
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

		// Get the actual movie ID (might be existing or newly created)
		if result.Next(ctx) {
			record := result.Record()
			if movieId, ok := record.Get("movieId"); ok {
				return movieId, nil
			}
			return nil, fmt.Errorf("movieId field not found in result")
		}

		if err = result.Err(); err != nil {
			return nil, fmt.Errorf("error reading movie result: %w", err)
		}

		return nil, fmt.Errorf("failed to get movie ID from result: no records returned")
	})
	if err != nil {
		return nil, err
	}

	// The write returned the created movie ID string; fetch the full movie object
	idStr, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected create movie result type: %T", result)
	}
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created movie id: %w", err)
	}
	return r.GetMovieByID(ctx, parsedID)
}

// GetMovieByID retrieves a movie by its ID
func (r *Neo4jRepository) GetMovieByID(ctx context.Context, id uuid.UUID) (*model.Movie, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
				MATCH (m:Movie {id: $id})
				OPTIONAL MATCH (m)<-[ract:ACTED_IN]-(cast:Person)
				OPTIONAL MATCH (m)<-[rcrew:CREW_ON]-(crew:Person)
				OPTIONAL MATCH (m)<-[:PRODUCED]-(pc:ProductionCompany)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(m) WHERE t.type = 'genre'
				OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
				RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
				       m.description as description, m.coverUrl as coverUrl,
				       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
				       m.searchDepth as searchDepth,
				       collect(DISTINCT cast {.*, externalID: cast.externalID}) as cast,
				       collect(DISTINCT crew {.*, externalID: crew.externalID}) as crew,
				       collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
				       collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits,
				       collect(DISTINCT pc) as productionCompanies,
				       collect(DISTINCT t) as genres,
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
			movie := &model.Movie{
				ID:                  id,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
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
func (r *Neo4jRepository) GetAllMovies(ctx context.Context, limit, offset *int) ([]*model.Movie, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Paginate on bare Movie nodes first, then expand joins — avoids O(N×cast×crew×...)
		// row explosion before aggregation.
		query := `
				MATCH (m:Movie)
				WITH m ORDER BY m.title
			`

		params := make(map[string]any)

		if offset != nil {
			query += " SKIP $offset"
			params["offset"] = *offset
		}
		if limit != nil {
			query += " LIMIT $limit"
			params["limit"] = *limit
		}

		query += `
				OPTIONAL MATCH (m)<-[ract:ACTED_IN]-(cast:Person)
				OPTIONAL MATCH (m)<-[rcrew:CREW_ON]-(crew:Person)
				OPTIONAL MATCH (m)<-[:PRODUCED]-(pc:ProductionCompany)
				OPTIONAL MATCH (t:Tag)-[:TAGGED]->(m) WHERE t.type = 'genre'
				OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
				RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
				       m.description as description, m.coverUrl as coverUrl,
				       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
				       m.searchDepth as searchDepth,
				       collect(DISTINCT cast {.*, externalID: cast.externalID}) as cast,
				       collect(DISTINCT crew {.*, externalID: crew.externalID}) as crew,
				       collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
				       collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits,
				       collect(DISTINCT pc) as productionCompanies,
				       collect(DISTINCT t) as genres,
				       collect(DISTINCT pcountry) as productionCountries
				ORDER BY m.title
			`

		result, err := tx.Run(ctx, query, params)

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

			castCredits := parsePersonCredits(record.AsMap()["castCredits"])
			crewCredits := parseCrewCredits(record.AsMap()["crewCredits"])
			movie := &model.Movie{
				ID:                  movieID,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
				SearchDepth:         getInt32Value(record.AsMap()["searchDepth"]),
				Cast:                parsePersons(record.AsMap()["cast"]),
				Crew:                parsePersons(record.AsMap()["crew"]),
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
	return len(input.Cast) == 0 || input.Description == nil || input.CoverURL == nil || input.Runtime == nil
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
		ID:          safeStringValue(input.ExternalID),
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

	// Populate new fields from metadata if not provided. Prefer structured credits when available.
	if len(enriched.Cast) == 0 {
		if len(meta.CastCredits) > 0 {
			// Convert structured credits to names for the CreateMovieInput's string slices
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

// helper to extract a properties map from various node representations
func extractPropsFromNode(item interface{}) (map[string]any, bool) {
	if item == nil {
		return nil, false
	}

	// Direct map
	if m, ok := item.(map[string]any); ok {
		return m, true
	}
	if m2, ok := item.(map[string]interface{}); ok {
		// convert to map[string]any
		res := make(map[string]any, len(m2))
		for k, v := range m2 {
			res[k] = v
		}
		return res, true
	}

	// If the type implements Props() map[string]any
	type propsProvider interface{ Props() map[string]any }
	if p, ok := item.(propsProvider); ok {
		return p.Props(), true
	}
	// Some implementations use map[string]interface{}
	type propsProviderInterfaceMap interface{ Props() map[string]interface{} }
	if p2, ok := item.(propsProviderInterfaceMap); ok {
		orig := p2.Props()
		res := make(map[string]any, len(orig))
		for k, v := range orig {
			res[k] = v
		}
		return res, true
	}

	// Reflection fallback: look for exported field "Props"
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsValid() && v.Kind() == reflect.Struct {
		f := v.FieldByName("Props")
		if f.IsValid() {
			iface := f.Interface()
			if m, ok := iface.(map[string]any); ok {
				return m, true
			}
			if m2, ok := iface.(map[string]interface{}); ok {
				res := make(map[string]any, len(m2))
				for k, v := range m2 {
					res[k] = v
				}
				return res, true
			}
		}
	}

	return nil, false
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

	for _, item := range arr {
		if nodeMap, ok := extractPropsFromNode(item); ok {
			persons = append(persons, personFromNodeMap(nodeMap))
		}
	}

	return persons
}

// parsePersonCredits converts a list of {person, character, order} maps returned by Cypher into []*model.PersonCredit
func parsePersonCredits(value interface{}) []*model.PersonCredit {
	if value == nil {
		return []*model.PersonCredit{}
	}

	var credits []*model.PersonCredit
	arr, ok := value.([]any)
	if !ok {
		return credits
	}

	for _, item := range arr {
		// item may be a map with keys person (node), character, order
		if m, ok := item.(map[string]any); ok {
			var personNode any
			if p, exists := m["person"]; exists {
				personNode = p
			} else if p, exists := m["cast"]; exists { // some Cypher aliases may differ
				personNode = p
			}
			var person *model.Person
			if nodeMap, ok := extractPropsFromNode(personNode); ok {
				person = personFromNodeMap(nodeMap)
			} else {
				person = &model.Person{ID: uuid.UUID{}, Name: ""}
			}
			name := person.Name
			if n, ok := m["name"].(string); ok && n != "" {
				name = n
			}
			var character *string
			if c, ok := m["character"].(string); ok && c != "" {
				character = &c
			}
			var order *int32
			if o, ok := m["order"]; ok {
				if v, ok := o.(int64); ok {
					i := int32(v)
					order = &i
				} else if v, ok := o.(int32); ok {
					order = &v
				} else if v, ok := o.(float64); ok {
					i := int32(v)
					order = &i
				} else if s, ok := o.(string); ok {
					if parsed, err := strconv.ParseInt(s, 10, 32); err == nil {
						i := int32(parsed)
						order = &i
					}
				}
			}

			credits = append(credits, &model.PersonCredit{Person: person, Name: name, Character: character, Order: order})
		}
	}

	return credits
}

// parseCrewCredits converts a list of {person, job, department} maps returned by Cypher into []*model.CrewCredit
func parseCrewCredits(value interface{}) []*model.CrewCredit {
	if value == nil {
		return []*model.CrewCredit{}
	}

	var credits []*model.CrewCredit
	arr, ok := value.([]any)
	if !ok {
		return credits
	}

	for _, item := range arr {
		if m, ok := item.(map[string]any); ok {
			var personNode any
			if p, exists := m["person"]; exists {
				personNode = p
			} else if p, exists := m["crew"]; exists {
				personNode = p
			}
			var person *model.Person
			if nodeMap, ok := extractPropsFromNode(personNode); ok {
				person = personFromNodeMap(nodeMap)
			} else {
				person = &model.Person{ID: uuid.UUID{}, Name: ""}
			}
			name := person.Name
			if n, ok := m["name"].(string); ok && n != "" {
				name = n
			}
			var job *string
			if j, ok := m["job"].(string); ok && j != "" {
				job = &j
			}
			var department *string
			if d, ok := m["department"].(string); ok && d != "" {
				department = &d
			}

			credits = append(credits, &model.CrewCredit{Person: person, Name: name, Job: job, Department: department})
		}
	}

	return credits
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

	var externalID *string
	if e, ok := nodeMap["externalID"]; ok {
		if s, ok := e.(string); ok && s != "" {
			externalID = &s
		} else if i, ok := e.(int); ok && i != 0 {
			s := strconv.Itoa(i)
			externalID = &s
		} else if i64, ok := e.(int64); ok && i64 != 0 {
			s := strconv.FormatInt(i64, 10)
			externalID = &s
		}
	}

	return &model.Person{
		ID:         id,
		Name:       name,
		ExternalID: externalID,
		ActedIn:    []*model.Movie{}, // Not populated here
		CrewOn:     []*model.Movie{},
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

	for _, item := range arr {
		if nodeMap, ok := extractPropsFromNode(item); ok {
			pcs = append(pcs, productionCompanyFromNodeMap(nodeMap))
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

	for _, item := range arr {
		if nodeMap, ok := extractPropsFromNode(item); ok {
			genres = append(genres, genreFromNodeMap(nodeMap))
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

	for _, item := range arr {
		if nodeMap, ok := extractPropsFromNode(item); ok {
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
		}
	}

	return pcs
}
