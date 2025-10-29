package db

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// normalizeName trims, lowercases, strips diacritics and punctuation for deduplication
func normalizeName(s string) string {
	// basic trim + lowercase
	s = strings.ToLower(strings.TrimSpace(s))

	// Decompose and remove diacritic marks (Mn)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, err := transform.String(t, s)
	if err == nil {
		s = res
	}

	// Remove punctuation and symbol characters, collapse whitespace
	var b []rune
	lastWasSpace := false
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			continue
		}
		if unicode.IsSpace(r) {
			if lastWasSpace {
				continue
			}
			lastWasSpace = true
			b = append(b, ' ')
			continue
		}
		lastWasSpace = false
		b = append(b, r)
	}

	return strings.TrimSpace(string(b))
}

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

	movieUUID := uuid.New()

	// Prepare data for nodes: prefer structured credits from metadata when available.
	// Build castData: if input.Cast provided prefer it, else use meta.CastCredits or meta.Cast.
	var castNames []string
	if len(input.Cast) > 0 {
		castNames = input.Cast
	}
	castData := make([]map[string]any, 0)
	if len(castNames) > 0 {
		castData = make([]map[string]any, len(castNames))
		for i, name := range castNames {
			castData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	} else if meta != nil && len(meta.CastCredits) > 0 {
		castData = make([]map[string]any, len(meta.CastCredits))
		for i, c := range meta.CastCredits {
			castData[i] = map[string]any{
				"name":      c.Name,
				"id":        uuid.New().String(),
				"character": c.Character,
				"order":     c.Order,
			}
		}
	} else if meta != nil && len(meta.Cast) > 0 {
		castData = make([]map[string]any, len(meta.Cast))
		for i, name := range meta.Cast {
			castData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	}

	// Build crewData similarly
	var crewNames []string
	if len(input.Crew) > 0 {
		crewNames = input.Crew
	}
	crewData := make([]map[string]any, 0)
	if len(crewNames) > 0 {
		crewData = make([]map[string]any, len(crewNames))
		for i, name := range crewNames {
			crewData[i] = map[string]any{"name": name, "id": uuid.New().String()}
		}
	} else if meta != nil && len(meta.CrewCredits) > 0 {
		crewData = make([]map[string]any, len(meta.CrewCredits))
		for i, c := range meta.CrewCredits {
			crewData[i] = map[string]any{
				"name":       c.Name,
				"id":         uuid.New().String(),
				"job":        c.Job,
				"department": c.Department,
			}
		}
	} else if meta != nil && len(meta.Crew) > 0 {
		crewData = make([]map[string]any, len(meta.Crew))
		for i, name := range meta.Crew {
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
			castData[i]["normalizedName"] = normalizeName(n)
		}
	}
	for i := range crewData {
		if n, ok := crewData[i]["name"].(string); ok {
			crewData[i]["normalizedName"] = normalizeName(n)
		}
	}
	for i := range pcData {
		if n, ok := pcData[i]["name"].(string); ok {
			pcData[i]["normalizedName"] = normalizeName(n)
		}
	}
	for i := range genreData {
		if n, ok := genreData[i]["name"].(string); ok {
			genreData[i]["normalizedName"] = normalizeName(n)
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

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Create the movie node and related nodes/relationships. This write does not RETURN records
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
				MERGE (p:Person {normalizedName: castData.normalizedName})
				ON CREATE SET p.id = castData.id, p.name = castData.name
				MERGE (p)-[r:ACTED_IN]->(m)
				ON CREATE SET r.character = castData.character, r.order = castData.order
			)
			WITH m
			FOREACH (crewData IN $crew |
				MERGE (p:Person {normalizedName: crewData.normalizedName})
				ON CREATE SET p.id = crewData.id, p.name = crewData.name
				MERGE (p)-[r:CREW_ON]->(m)
				ON CREATE SET r.job = crewData.job, r.department = crewData.department
			)

				WITH m
				FOREACH (pcData IN $productionCompanies |
					MERGE (pc:ProductionCompany {normalizedName: pcData.normalizedName})
					ON CREATE SET pc.id = pcData.id, pc.name = pcData.name
					MERGE (pc)-[:PRODUCED]->(m)
				)
				WITH m
				FOREACH (genreData IN $genres |
					MERGE (g:Genre {normalizedName: genreData.normalizedName})
					ON CREATE SET g.id = genreData.id, g.name = genreData.name
					MERGE (g)-[:HAS_MOVIE]->(m)
				)
				WITH m
				FOREACH (pcountryData IN $productionCountries |
					MERGE (pcountry:ProductionCountry {id: pcountryData.id})
					ON CREATE SET pcountry.name = pcountryData.name
					MERGE (pcountry)-[:PRODUCED_IN]->(m)
				)
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

		_, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, fmt.Errorf("failed to run create movie query: %w", err)
		}

		// Return the created movie ID so the caller can fetch the full object via GetMovieByID
		return movieUUID.String(), nil
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
				OPTIONAL MATCH (g:Genre)-[:HAS_MOVIE]->(m)
				OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
				RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
				       m.description as description, m.coverUrl as coverUrl,
				       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
				       collect(DISTINCT cast) as cast,
				       collect(DISTINCT crew) as crew,
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
			movie := &model.Movie{
				ID:                  id,
				Title:               record.AsMap()["title"].(string),
				ReleaseDate:         getStringPointer(record.AsMap()["releaseDate"]),
				Description:         getStringPointer(record.AsMap()["description"]),
				CoverURL:            getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:             getInt32Pointer(record.AsMap()["runtime"]),
				Budget:              getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:           getInt32Pointer(record.AsMap()["boxOffice"]),
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
func (r *Neo4jRepository) GetAllMovies(ctx context.Context) ([]*model.Movie, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
				MATCH (m:Movie)
				OPTIONAL MATCH (m)<-[ract:ACTED_IN]-(cast:Person)
				OPTIONAL MATCH (m)<-[rcrew:CREW_ON]-(crew:Person)
				OPTIONAL MATCH (m)<-[:PRODUCED]-(pc:ProductionCompany)
				OPTIONAL MATCH (g:Genre)-[:HAS_MOVIE]->(m)
				OPTIONAL MATCH (m)<-[:PRODUCED_IN]-(pcountry:ProductionCountry)
				RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
				       m.description as description, m.coverUrl as coverUrl,
				       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice,
				       collect(DISTINCT cast) as cast,
				       collect(DISTINCT crew) as crew,
				       collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
				       collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits,
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
