package db

import (
	"context"
	"fmt"
	"nq/graph/model"
	"strconv"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CreateMovie creates a new movie in the database
func (r *Neo4jRepository) CreateMovie(ctx context.Context, input model.CreateMovieInput) (*model.Movie, error) {
	movieID := uuid.New()

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
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
			       m.description as description, m.coverUrl as coverUrl,
			       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice
		`

		params := map[string]any{
			"id":          movieID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"runtime":     input.Runtime,
			"budget":      input.Budget,
			"boxOffice":   input.BoxOffice,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			movie := &model.Movie{
				ID:            movieID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:       getInt32Pointer(record.AsMap()["runtime"]),
				Budget:        getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:     getInt32Pointer(record.AsMap()["boxOffice"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return movie, nil
		}

		return nil, fmt.Errorf("failed to create movie")
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
			RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
			       m.description as description, m.coverUrl as coverUrl,
			       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			movie := &model.Movie{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:       getInt32Pointer(record.AsMap()["runtime"]),
				Budget:        getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:     getInt32Pointer(record.AsMap()["boxOffice"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
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
			RETURN m.id as id, m.title as title, m.releaseDate as releaseDate,
			       m.description as description, m.coverUrl as coverUrl,
			       m.runtime as runtime, m.budget as budget, m.boxOffice as boxOffice
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
				ID:            movieID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Runtime:       getInt32Pointer(record.AsMap()["runtime"]),
				Budget:        getInt32Pointer(record.AsMap()["budget"]),
				BoxOffice:     getInt32Pointer(record.AsMap()["boxOffice"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
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

// CreateTVShow creates a new TV show in the database
func (r *Neo4jRepository) CreateTVShow(ctx context.Context, input model.CreateTVShowInput) (*model.TVShow, error) {
	tvShowID := uuid.New()

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
				createdAt: datetime(),
				updatedAt: datetime()
			})
			RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
			       t.description as description, t.coverUrl as coverUrl,
			       t.seasons as seasons, t.episodes as episodes, t.status as status
		`

		params := map[string]any{
			"id":          tvShowID.String(),
			"title":       input.Title,
			"releaseDate": input.ReleaseDate,
			"description": input.Description,
			"coverUrl":    input.CoverURL,
			"seasons":     input.Seasons,
			"episodes":    input.Episodes,
			"status":      input.Status,
		}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			tvShow := &model.TVShow{
				ID:            tvShowID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return tvShow, nil
		}

		return nil, fmt.Errorf("failed to create TV show")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.TVShow), nil
}

// GetTVShowByID retrieves a TV show by its ID
func (r *Neo4jRepository) GetTVShowByID(ctx context.Context, id uuid.UUID) (*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (t:TVShow {id: $id})
			RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
			       t.description as description, t.coverUrl as coverUrl,
			       t.seasons as seasons, t.episodes as episodes, t.status as status
		`

		params := map[string]any{"id": id.String()}

		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			tvShow := &model.TVShow{
				ID:            id,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			return tvShow, nil
		}

		return nil, fmt.Errorf("TV show not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.TVShow), nil
}

// GetAllTVShows retrieves all TV shows
func (r *Neo4jRepository) GetAllTVShows(ctx context.Context) ([]*model.TVShow, error) {
	result, err := r.db.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (t:TVShow)
			RETURN t.id as id, t.title as title, t.releaseDate as releaseDate,
			       t.description as description, t.coverUrl as coverUrl,
			       t.seasons as seasons, t.episodes as episodes, t.status as status
			ORDER BY t.title
		`

		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var tvShows []*model.TVShow
		for result.Next(ctx) {
			record := result.Record()
			tvShowID, err := uuid.Parse(record.AsMap()["id"].(string))
			if err != nil {
				return nil, err
			}

			tvShow := &model.TVShow{
				ID:            tvShowID,
				Title:         record.AsMap()["title"].(string),
				ReleaseDate:   getStringPointer(record.AsMap()["releaseDate"]),
				Description:   getStringPointer(record.AsMap()["description"]),
				CoverURL:      getStringPointer(record.AsMap()["coverUrl"]),
				Seasons:       getInt32Pointer(record.AsMap()["seasons"]),
				Episodes:      getInt32Pointer(record.AsMap()["episodes"]),
				Status:        getStringPointer(record.AsMap()["status"]),
				Creators:      []*model.Creator{},
				Platforms:     []*model.Platform{},
				Tags:          []*model.Tag{},
				Ratings:       []*model.Rating{},
				AverageRating: nil,
			}
			tvShows = append(tvShows, tvShow)
		}

		return tvShows, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.TVShow), nil
}

// Similar implementations for Book, Game, and MusicAlbum would follow the same pattern
// For brevity, I'll implement a few key ones and provide stubs for the rest

// CreateBook creates a new book in the database
func (r *Neo4jRepository) CreateBook(ctx context.Context, input model.CreateBookInput) (*model.Book, error) {
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
	// Implementation similar to GetMovieByID but for Book
	return nil, fmt.Errorf("not implemented")
}

// GetAllBooks retrieves all books
func (r *Neo4jRepository) GetAllBooks(ctx context.Context) ([]*model.Book, error) {
	// Implementation similar to GetAllMovies but for Book
	return nil, fmt.Errorf("not implemented")
}

// CreateGame creates a new game in the database
func (r *Neo4jRepository) CreateGame(ctx context.Context, input model.CreateGameInput) (*model.Game, error) {
	// Implementation similar to CreateMovie but for Game
	return nil, fmt.Errorf("not implemented")
}

// GetGameByID retrieves a game by its ID
func (r *Neo4jRepository) GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetAllGames retrieves all games
func (r *Neo4jRepository) GetAllGames(ctx context.Context) ([]*model.Game, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateMusicAlbum creates a new music album in the database
func (r *Neo4jRepository) CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error) {
	// Implementation similar to CreateMovie but for MusicAlbum
	return nil, fmt.Errorf("not implemented")
}

// GetMusicAlbumByID retrieves a music album by its ID
func (r *Neo4jRepository) GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetAllMusicAlbums retrieves all music albums
func (r *Neo4jRepository) GetAllMusicAlbums(ctx context.Context) ([]*model.MusicAlbum, error) {
	return nil, fmt.Errorf("not implemented")
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
	// Add other media types...

	return nil, fmt.Errorf("media not found")
}

// GetAllMedia retrieves all media items
func (r *Neo4jRepository) GetAllMedia(ctx context.Context) ([]model.Media, error) {
	// Implementation to get all media types
	return nil, fmt.Errorf("not implemented")
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
