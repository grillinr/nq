package db

import (
	"context"
	"nq/graph/model"

	"github.com/google/uuid"
)

// Repository defines the interface for all database operations
type Repository interface {
	UserRepository
	MediaRepository
	ActivityRepository
	RatingRepository
	RecommendationRepository
}

// UserRepository defines operations for user management
type UserRepository interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input model.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// MediaRepository defines operations for media management
type MediaRepository interface {
	// Movie operations
	CreateMovie(ctx context.Context, input model.CreateMovieInput) (*model.Movie, error)
	GetMovieByID(ctx context.Context, id uuid.UUID) (*model.Movie, error)
	GetAllMovies(ctx context.Context) ([]*model.Movie, error)

	// TV Show operations
	CreateTVShow(ctx context.Context, input model.CreateTVShowInput) (*model.TVShow, error)
	GetTVShowByID(ctx context.Context, id uuid.UUID) (*model.TVShow, error)
	GetAllTVShows(ctx context.Context) ([]*model.TVShow, error)

	// Book operations
	CreateBook(ctx context.Context, input model.CreateBookInput) (*model.Book, error)
	GetBookByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	GetAllBooks(ctx context.Context) ([]*model.Book, error)

	// Game operations
	CreateGame(ctx context.Context, input model.CreateGameInput) (*model.Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	GetAllGames(ctx context.Context) ([]*model.Game, error)

	// Music Album operations
	CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error)
	GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error)
	GetAllMusicAlbums(ctx context.Context) ([]*model.MusicAlbum, error)

	// Generic media operations
	GetMediaByID(ctx context.Context, id uuid.UUID) (model.Media, error)
	GetAllMedia(ctx context.Context) ([]model.Media, error)
}

// ActivityRepository defines operations for user activities
type ActivityRepository interface {
	CreateActivity(ctx context.Context, input model.CreateActivityInput) (*model.UserActivity, error)
	GetActivityByID(ctx context.Context, id uuid.UUID) (*model.UserActivity, error)
	GetUserActivities(ctx context.Context, userID uuid.UUID) ([]*model.UserActivity, error)
	GetMediaActivities(ctx context.Context, mediaID uuid.UUID) ([]*model.UserActivity, error)
	UpdateActivity(ctx context.Context, id uuid.UUID, statusID *int32, rating *float64, review *string, finishedAt *string) (*model.UserActivity, error)
	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

// RatingRepository defines operations for ratings
type RatingRepository interface {
	CreateRating(ctx context.Context, userID, mediaID uuid.UUID, score float64) (*model.Rating, error)
	GetRating(ctx context.Context, userID, mediaID uuid.UUID) (*model.Rating, error)
	GetUserRatings(ctx context.Context, userID uuid.UUID) ([]*model.Rating, error)
	GetMediaRatings(ctx context.Context, mediaID uuid.UUID) ([]*model.Rating, error)
	UpdateRating(ctx context.Context, userID, mediaID uuid.UUID, score float64) (*model.Rating, error)
	DeleteRating(ctx context.Context, userID, mediaID uuid.UUID) error
	GetAverageRating(ctx context.Context, mediaID uuid.UUID) (*float64, error)
}

// RecommendationRepository defines operations for recommendations
type RecommendationRepository interface {
	CreateRecommendation(ctx context.Context, userID, mediaID uuid.UUID, recommenderID *uuid.UUID, source *string, score *float64) (*model.Recommendation, error)
	GetRecommendations(ctx context.Context, userID uuid.UUID) ([]*model.Recommendation, error)
	GetRecommendationByID(ctx context.Context, id uuid.UUID) (*model.Recommendation, error)
	DeleteRecommendation(ctx context.Context, id uuid.UUID) error
}

// Neo4jRepository implements the Repository interface using Neo4j
type Neo4jRepository struct {
	db *Database
}

// NewNeo4jRepository creates a new Neo4j repository
func NewNeo4jRepository(db *Database) *Neo4jRepository {
	return &Neo4jRepository{db: db}
}

// Helper method to get the database instance
func (r *Neo4jRepository) getDB() *Database {
	return r.db
}
