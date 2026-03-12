package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grillinr/nq/graph/model"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateMusicAlbum creates a new music album in the database
func (r *Neo4jRepository) CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error) {
	albumID := uuid.New()

	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (a:MusicAlbum:Media {
				id: $id,
				title: $title,
				releaseDate: $releaseDate,
				description: $description,
				coverUrl: $coverUrl,
				trackCount: $trackCount,
				duration: $duration,
				label: $label
			})
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
		`

		params := map[string]any{
			"id":          albumID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"trackCount":  input.TrackCount,
			"duration":    input.Duration,
			"label":       input.Label,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			album := &model.MusicAlbum{
				ID:          albumID,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			return album, nil
		}

		if err := result.Err(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("create music album returned no records")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.MusicAlbum), nil
}

// GetMusicAlbumByID retrieves a music album by its ID
func (r *Neo4jRepository) GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:MusicAlbum {id: $id})
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			album := &model.MusicAlbum{
				ID:          id,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			return album, nil
		}

		return nil, fmt.Errorf("music album not found")
	})
	if err != nil {
		return nil, err
	}

	return result.(*model.MusicAlbum), nil
}

// GetAllMusicAlbums retrieves all music albums with optional pagination.
func (r *Neo4jRepository) GetAllMusicAlbums(ctx context.Context, limit, offset *int) ([]*model.MusicAlbum, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Paginate on bare MusicAlbum nodes first, then project fields — consistent
		// with the pattern used by GetAllMovies, GetAllTVShows, and GetAllBooks.
		query := `
			MATCH (a:MusicAlbum:Media)
			WITH a ORDER BY a.title
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
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
		`

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var albums []*model.MusicAlbum
		for result.Next(ctx) {
			record := result.Record()
			idStr, _ := record.Get("id")
			albumID, _ := uuid.Parse(idStr.(string))
			album := &model.MusicAlbum{
				ID:          albumID,
				Title:       record.AsMap()["title"].(string),
				TrackCount:  getInt32Pointer(record.AsMap()["trackCount"]),
				Duration:    getInt32Pointer(record.AsMap()["duration"]),
				Label:       getStringPointer(record.AsMap()["label"]),
				ReleaseDate: getStringPointer(record.AsMap()["releaseDate"]),
				Description: getStringPointer(record.AsMap()["description"]),
				CoverURL:    getStringPointer(record.AsMap()["coverUrl"]),
			}
			albums = append(albums, album)
		}

		return albums, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]*model.MusicAlbum), nil
}

// GetMediaByID retrieves any media by its ID using a single DB round-trip.
// It inspects the node's labels to dispatch directly to the correct typed fetch,
// eliminating the previous sequential waterfall of up to 5 DB calls.
func (r *Neo4jRepository) GetMediaByID(ctx context.Context, id uuid.UUID) (model.Media, error) {
	mediaType, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (m:Media {id: $id})
			RETURN
				CASE
					WHEN m:Movie       THEN 'Movie'
					WHEN m:TVShow      THEN 'TVShow'
					WHEN m:Book        THEN 'Book'
					WHEN m:Game        THEN 'Game'
					WHEN m:MusicAlbum  THEN 'MusicAlbum'
					ELSE 'Unknown'
				END AS mediaType
		`
		result, err := tx.Run(ctx, query, map[string]any{"id": id.String()})
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			t, _ := result.Record().Get("mediaType")
			return t.(string), nil
		}
		return nil, fmt.Errorf("media not found")
	})
	if err != nil {
		return nil, err
	}

	switch mediaType.(string) {
	case "Movie":
		return r.GetMovieByID(ctx, id)
	case "TVShow":
		return r.GetTVShowByID(ctx, id)
	case "Book":
		return r.GetBookByID(ctx, id)
	case "Game":
		return r.GetGameByID(ctx, id)
	case "MusicAlbum":
		return r.GetMusicAlbumByID(ctx, id)
	default:
		return nil, fmt.Errorf("media not found")
	}
}

// GetAllMedia retrieves all media items
func (r *Neo4jRepository) GetAllMedia(ctx context.Context) ([]model.Media, error) {
	var mediaItems []model.Media

	movies, err := r.GetAllMovies(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, m := range movies {
		mediaItems = append(mediaItems, m)
	}

	books, err := r.GetAllBooks(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		mediaItems = append(mediaItems, b)
	}

	games, err := r.GetAllGames(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, g := range games {
		mediaItems = append(mediaItems, g)
	}

	tvShows, err := r.GetAllTVShows(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, tv := range tvShows {
		mediaItems = append(mediaItems, tv)
	}

	albums, err := r.GetAllMusicAlbums(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, a := range albums {
		mediaItems = append(mediaItems, a)
	}

	return mediaItems, nil
}

// GetCastAndCrew retrieves cast and crew for any media item by ID
func (r *Neo4jRepository) GetCastAndCrew(ctx context.Context, mediaID uuid.UUID) ([]*model.Person, []*model.Person, []*model.PersonCredit, []*model.CrewCredit, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (m:Media {id: $id})
			OPTIONAL MATCH (m)<-[ract:ACTED_IN]-(cast:Person)
			OPTIONAL MATCH (m)<-[rcrew:CREW_ON]-(crew:Person)
			RETURN collect(DISTINCT cast) as cast, collect(DISTINCT crew) as crew,
			       collect(DISTINCT {person: cast, character: ract.character, order: ract.order, name: cast.name}) as castCredits,
			       collect(DISTINCT {person: crew, job: rcrew.job, department: rcrew.department, name: crew.name}) as crewCredits
		`

		params := map[string]any{"id": mediaID.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			rec := result.Record()
			castParsed := parsePersons(rec.AsMap()["cast"])
			crewParsed := parsePersons(rec.AsMap()["crew"])
			castCredits := parsePersonCredits(rec.AsMap()["castCredits"])
			crewCredits := parseCrewCredits(rec.AsMap()["crewCredits"])
			return map[string]any{"cast": castParsed, "crew": crewParsed, "castCredits": castCredits, "crewCredits": crewCredits}, nil
		}

		return nil, fmt.Errorf("media not found")
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	m, ok := result.(map[string]any)
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("unexpected result type: %T", result)
	}

	cast, _ := m["cast"].([]*model.Person)
	crew, _ := m["crew"].([]*model.Person)
	castCredits, _ := m["castCredits"].([]*model.PersonCredit)
	crewCredits, _ := m["crewCredits"].([]*model.CrewCredit)
	return cast, crew, castCredits, crewCredits, nil
}

// Helper function to safely get int32 pointer from interface{}
func getInt32Pointer(value interface{}) *int32 {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int32:
		return &v
	case int64:
		i32 := int32(v)
		return &i32
	case float64:
		i32 := int32(v)
		return &i32
	case string:
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			i32 := int32(i)
			return &i32
		}
	}

	return nil
}

// Helper function to safely get int32 value from interface{}, default 0
func getInt32Value(value interface{}) int32 {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(i)
		}
	}

	return 0
}

// FindMediaByTitleTypeYear finds media by title, type, and optional year
func (r *Neo4jRepository) FindMediaByTitleTypeYear(ctx context.Context, title, mediaType string, year *int) (model.Media, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var query string
		var params map[string]any

		if year != nil {
			query = `
				MATCH (m:Media)
				WHERE m.title = $title AND labels(m) = [$mediaType] AND m.releaseDate STARTS WITH $yearStr
				RETURN m.id as id, labels(m) as labels
				LIMIT 1
			`
			params = map[string]any{
				"title":     title,
				"mediaType": mediaType,
				"yearStr":   fmt.Sprintf("%d", *year),
			}
		} else {
			query = `
				MATCH (m:Media)
				WHERE m.title = $title AND labels(m) = [$mediaType]
				RETURN m.id as id, labels(m) as labels
				LIMIT 1
			`
			params = map[string]any{
				"title":     title,
				"mediaType": mediaType,
			}
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			idStr := record.AsMap()["id"].(string)
			id, err := uuid.Parse(idStr)
			if err != nil {
				return nil, err
			}
			return r.GetMediaByID(ctx, id)
		}

		return nil, fmt.Errorf("media not found")
	})
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("media not found")
	}

	return result.(model.Media), nil
}

// UpdateMediaSearchDepth updates the searchDepth of a media item
func (r *Neo4jRepository) UpdateMediaSearchDepth(ctx context.Context, id uuid.UUID, searchDepth int32) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (m:Media {id: $id})
			SET m.searchDepth = $searchDepth
		`
		params := map[string]any{
			"id":          id.String(),
			"searchDepth": searchDepth,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// LinkRelatedMedia creates a relationship between media items.
func (r *Neo4jRepository) LinkRelatedMedia(ctx context.Context, sourceID, relatedID uuid.UUID) error {
	if sourceID == relatedID {
		return nil
	}
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (source:Media {id: $sourceID}), (related:Media {id: $relatedID})
			MERGE (source)-[:RELATED_TO]->(related)
			MERGE (related)-[:RELATED_TO]->(source)
		`
		params := map[string]any{
			"sourceID":  sourceID.String(),
			"relatedID": relatedID.String(),
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// LinkRelatedMediaByTagNames creates RELATED_TO links for media sharing any normalized tag names.
// Returns the number of related media linked.
func (r *Neo4jRepository) LinkRelatedMediaByTagNames(ctx context.Context, sourceID uuid.UUID, normalizedNames []string, limit int) (int, error) {
	if sourceID == uuid.Nil || len(normalizedNames) == 0 {
		return 0, nil
	}
	if limit <= 0 {
		limit = 25
	}
	result, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (source:Media {id: $sourceID})
			MATCH (t:Tag)
			WHERE t.normalizedName IN $normalizedNames
			MATCH (t)-[:TAGGED]->(related:Media)
			WHERE related.id <> $sourceID
			WITH DISTINCT related
			LIMIT $limit
			MERGE (source)-[:RELATED_TO]->(related)
			MERGE (related)-[:RELATED_TO]->(source)
			RETURN count(related) as linkedCount
		`
		params := map[string]any{
			"sourceID":        sourceID.String(),
			"normalizedNames": normalizedNames,
			"limit":           limit,
		}
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if res.Next(ctx) {
			if count, ok := res.Record().Get("linkedCount"); ok {
				switch v := count.(type) {
				case int64:
					return int(v), nil
				case int:
					return v, nil
				case float64:
					return int(v), nil
				}
			}
		}
		if err := res.Err(); err != nil {
			return nil, err
		}
		return 0, nil
	})
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return result.(int), nil
}

// GetRelatedMedia fetches media linked via RELATED_TO edges, up to limit items.
// All fields are returned in a single Cypher query — no per-item round-trips.
func (r *Neo4jRepository) GetRelatedMedia(ctx context.Context, sourceID uuid.UUID, limit int) ([]model.Media, error) {
	if limit <= 0 {
		limit = 12
	}

	results, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (source:Media {id: $sourceID})-[:RELATED_TO]->(related:Media)
			WITH related LIMIT $limit
			RETURN related.id          AS id,
			       related.title       AS title,
			       related.releaseDate AS releaseDate,
			       related.description AS description,
			       related.coverUrl    AS coverUrl,
			       related.searchDepth AS searchDepth,
			       CASE
			           WHEN related:Movie      THEN 'Movie'
			           WHEN related:TVShow     THEN 'TVShow'
			           WHEN related:Book       THEN 'Book'
			           WHEN related:Game       THEN 'Game'
			           WHEN related:MusicAlbum THEN 'MusicAlbum'
			           ELSE 'Unknown'
			       END AS mediaType,
			       related.runtime     AS runtime,
			       related.budget      AS budget,
			       related.boxOffice   AS boxOffice,
			       related.seasons     AS seasons,
			       related.episodes    AS episodes,
			       related.status      AS status,
			       related.pages       AS pages,
			       related.isbn        AS isbn,
			       related.publisher   AS publisher,
			       related.trackCount  AS trackCount,
			       related.duration    AS duration,
			       related.label       AS label,
			       related.genre       AS genre,
			       related.themes      AS themes,
			       related.keywords    AS keywords,
			       related.gameModes   AS gameModes,
			       related.perspectives AS perspectives,
			       related.franchises  AS franchises,
			       related.platformsList AS platformsList,
			       related.esrbRating  AS esrbRating,
			       related.multiplayer AS multiplayer
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"sourceID": sourceID.String(),
			"limit":    limit,
		})
		if err != nil {
			return nil, err
		}

		var media []model.Media
		for result.Next(ctx) {
			rec := result.Record()
			m := rec.AsMap()

			idRaw, _ := m["id"]
			idStr, _ := idRaw.(string)
			id, err := uuid.Parse(idStr)
			if err != nil {
				continue
			}
			mediaType, _ := m["mediaType"].(string)

			switch mediaType {
			case "Movie":
				media = append(media, &model.Movie{
					ID:                  id,
					Title:               m["title"].(string),
					ReleaseDate:         getStringPointer(m["releaseDate"]),
					Description:         getStringPointer(m["description"]),
					CoverURL:            getStringPointer(m["coverUrl"]),
					SearchDepth:         getInt32Value(m["searchDepth"]),
					Runtime:             getInt32Pointer(m["runtime"]),
					Budget:              getInt32Pointer(m["budget"]),
					BoxOffice:           getInt32Pointer(m["boxOffice"]),
					Cast:                []*model.Person{},
					Crew:                []*model.Person{},
					CastCredits:         []*model.PersonCredit{},
					CrewCredits:         []*model.CrewCredit{},
					ProductionCompanies: []*model.ProductionCompany{},
					Genres:              []*model.Genre{},
					ProductionCountries: []*model.ProductionCountry{},
					Creators:            []*model.Creator{},
					Platforms:           []*model.Platform{},
					Tags:                []*model.Tag{},
					Ratings:             []*model.Rating{},
				})
			case "TVShow":
				media = append(media, &model.TVShow{
					ID:                  id,
					Title:               m["title"].(string),
					ReleaseDate:         getStringPointer(m["releaseDate"]),
					Description:         getStringPointer(m["description"]),
					CoverURL:            getStringPointer(m["coverUrl"]),
					SearchDepth:         getInt32Value(m["searchDepth"]),
					Seasons:             getInt32Pointer(m["seasons"]),
					Episodes:            getInt32Pointer(m["episodes"]),
					Status:              getStringPointer(m["status"]),
					Cast:                []*model.Person{},
					Crew:                []*model.Person{},
					CastCredits:         []*model.PersonCredit{},
					CrewCredits:         []*model.CrewCredit{},
					ProductionCompanies: []*model.ProductionCompany{},
					Genres:              []*model.Genre{},
					ProductionCountries: []*model.ProductionCountry{},
					Creators:            []*model.Creator{},
					Platforms:           []*model.Platform{},
					Tags:                []*model.Tag{},
					Ratings:             []*model.Rating{},
				})
			case "Book":
				media = append(media, &model.Book{
					ID:          id,
					Title:       m["title"].(string),
					ReleaseDate: getStringPointer(m["releaseDate"]),
					Description: getStringPointer(m["description"]),
					CoverURL:    getStringPointer(m["coverUrl"]),
					SearchDepth: getInt32Value(m["searchDepth"]),
					Pages:       getInt32Pointer(m["pages"]),
					Isbn:        getStringPointer(m["isbn"]),
					Publisher:   getStringPointer(m["publisher"]),
					Creators:    []*model.Creator{},
					Authors:     []*model.Creator{},
					Platforms:   []*model.Platform{},
					Tags:        []*model.Tag{},
					Subjects:    []*model.Tag{},
					Ratings:     []*model.Rating{},
				})
			case "Game":
				media = append(media, &model.Game{
					ID:            id,
					Title:         m["title"].(string),
					ReleaseDate:   getStringPointer(m["releaseDate"]),
					Description:   getStringPointer(m["description"]),
					CoverURL:      getStringPointer(m["coverUrl"]),
					SearchDepth:   getInt32Value(m["searchDepth"]),
					Genre:         getStringSlice(m["genre"]),
					Themes:        getStringSlice(m["themes"]),
					Keywords:      getStringSlice(m["keywords"]),
					GameModes:     getStringSlice(m["gameModes"]),
					Perspectives:  getStringSlice(m["perspectives"]),
					Franchises:    getStringSlice(m["franchises"]),
					PlatformsList: getStringSlice(m["platformsList"]),
					EsrbRating:    getStringPointer(m["esrbRating"]),
					Multiplayer:   getBoolPointer(m["multiplayer"]),
					Creators:      []*model.Creator{},
					Platforms:     []*model.Platform{},
					Tags:          []*model.Tag{},
					Ratings:       []*model.Rating{},
				})
			case "MusicAlbum":
				media = append(media, &model.MusicAlbum{
					ID:          id,
					Title:       m["title"].(string),
					ReleaseDate: getStringPointer(m["releaseDate"]),
					Description: getStringPointer(m["description"]),
					CoverURL:    getStringPointer(m["coverUrl"]),
					SearchDepth: getInt32Value(m["searchDepth"]),
					TrackCount:  getInt32Pointer(m["trackCount"]),
					Duration:    getInt32Pointer(m["duration"]),
					Label:       getStringPointer(m["label"]),
					Creators:    []*model.Creator{},
					Platforms:   []*model.Platform{},
					Tags:        []*model.Tag{},
					Ratings:     []*model.Rating{},
				})
			}
		}
		return media, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return results.([]model.Media), nil
}
