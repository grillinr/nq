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

// GetAllMusicAlbums retrieves all music albums
func (r *Neo4jRepository) GetAllMusicAlbums(ctx context.Context) ([]*model.MusicAlbum, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:MusicAlbum:Media)
			RETURN a.id as id, a.title as title, a.releaseDate as releaseDate,
			       a.description as description, a.coverUrl as coverUrl,
			       a.trackCount as trackCount, a.duration as duration, a.label as label
			ORDER BY a.title
		`

		result, err := tx.Run(ctx, query, nil)
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

// GetMediaByID retrieves any media by its ID
func (r *Neo4jRepository) GetMediaByID(ctx context.Context, id uuid.UUID) (model.Media, error) {
	// Try each media type
	if movie, err := r.GetMovieByID(ctx, id); err == nil {
		return movie, nil
	}
	if tvShow, err := r.GetTVShowByID(ctx, id); err == nil {
		return tvShow, nil
	}
	if book, err := r.GetBookByID(ctx, id); err == nil {
		return book, nil
	}
	if game, err := r.GetGameByID(ctx, id); err == nil {
		return game, nil
	}
	if album, err := r.GetMusicAlbumByID(ctx, id); err == nil {
		return album, nil
	}

	return nil, fmt.Errorf("media not found")
}

// GetAllMedia retrieves all media items
func (r *Neo4jRepository) GetAllMedia(ctx context.Context) ([]model.Media, error) {
	var mediaItems []model.Media

	movies, err := r.GetAllMovies(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range movies {
		mediaItems = append(mediaItems, m)
	}

	books, err := r.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		mediaItems = append(mediaItems, b)
	}

	games, err := r.GetAllGames(ctx)
	if err != nil {
		return nil, err
	}
	for _, g := range games {
		mediaItems = append(mediaItems, g)
	}

	tvShows, err := r.GetAllTVShows(ctx)
	if err != nil {
		return nil, err
	}
	for _, tv := range tvShows {
		mediaItems = append(mediaItems, tv)
	}

	albums, err := r.GetAllMusicAlbums(ctx)
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
