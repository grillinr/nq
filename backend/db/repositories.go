package db

import (
	"context"
	"log"

	"github.com/grillinr/nq/graph/model"
	"github.com/grillinr/nq/metadata"

	"github.com/google/uuid"
)

// Repository defines the interface for all database operations
type Repository interface {
	UserRepository
	MediaRepository
	ActivityRepository
	RatingRepository
	RecommendationRepository
	FriendRepository
}

// UserRepository defines operations for user management
type UserRepository interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByAuth(ctx context.Context, provider, subject string) (*model.User, error)
	GetOrCreateUserByAuth(ctx context.Context, provider, subject, email, name string, avatarURL *string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input model.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// AddToFavorites creates a FAVORITES relationship between user and media
	AddToFavorites(ctx context.Context, userID, mediaID uuid.UUID) error
}

// MediaRepository defines operations for media management
type MediaRepository interface {
	// Movie operations
	CreateMovie(ctx context.Context, input model.CreateMovieInput) (*model.Movie, error)
	GetMovieByID(ctx context.Context, id uuid.UUID) (*model.Movie, error)
	GetAllMovies(ctx context.Context, limit, offset *int) ([]*model.Movie, error)

	// TV Show operations
	CreateTVShow(ctx context.Context, input model.CreateTVShowInput) (*model.TVShow, error)
	GetTVShowByID(ctx context.Context, id uuid.UUID) (*model.TVShow, error)
	GetAllTVShows(ctx context.Context, limit, offset *int) ([]*model.TVShow, error)

	// Book operations
	CreateBook(ctx context.Context, input model.CreateBookInput) (*model.Book, error)
	GetBookByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	GetAllBooks(ctx context.Context, limit, offset *int) ([]*model.Book, error)

	// Game operations
	CreateGame(ctx context.Context, input model.CreateGameInput) (*model.Game, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*model.Game, error)
	GetAllGames(ctx context.Context, limit, offset *int) ([]*model.Game, error)

	// Music Album operations
	CreateMusicAlbum(ctx context.Context, input model.CreateMusicAlbumInput) (*model.MusicAlbum, error)
	GetMusicAlbumByID(ctx context.Context, id uuid.UUID) (*model.MusicAlbum, error)
	GetAllMusicAlbums(ctx context.Context, limit, offset *int) ([]*model.MusicAlbum, error)

	// Generic media operations
	GetMediaByID(ctx context.Context, id uuid.UUID) (model.Media, error)
	GetAllMedia(ctx context.Context) ([]model.Media, error)
	// GetCastAndCrew retrieves cast and crew for any media item by ID
	GetCastAndCrew(ctx context.Context, mediaID uuid.UUID) ([]*model.Person, []*model.Person, []*model.PersonCredit, []*model.CrewCredit, error)
	// FindMediaByTitleTypeYear checks if media exists by title, type, and year
	FindMediaByTitleTypeYear(ctx context.Context, title, mediaType string, year *int) (model.Media, error)
	// UpdateMediaSearchDepth updates the searchDepth of a media item
	UpdateMediaSearchDepth(ctx context.Context, id uuid.UUID, searchDepth int32) error
	// LinkRelatedMedia creates a relationship between media items
	LinkRelatedMedia(ctx context.Context, sourceID, relatedID uuid.UUID) error
	// LinkRelatedMediaByTagNames links media by shared tag normalized names
	LinkRelatedMediaByTagNames(ctx context.Context, sourceID uuid.UUID, normalizedNames []string, limit int) (int, error)
	// LinkMediaByNormalizedTitle links all media sharing the same normalized title (cross-type, e.g. book ↔ movie adaptation)
	LinkMediaByNormalizedTitle(ctx context.Context, sourceID uuid.UUID) (int, error)
	// GetRelatedMedia fetches media linked via RELATED_TO edges, limited to top N items
	GetRelatedMedia(ctx context.Context, sourceID uuid.UUID, limit int) ([]model.Media, error)
	// GetMetadata returns the metadata service
	GetMetadata() interface{}
}

// ActivityRepository defines operations for user activities
type ActivityRepository interface {
	CreateActivity(ctx context.Context, userID uuid.UUID, input model.CreateActivityInput) (*model.UserActivity, error)
	GetActivityByID(ctx context.Context, id uuid.UUID) (*model.UserActivity, error)
	GetUserActivities(ctx context.Context, userID uuid.UUID) ([]*model.UserActivity, error)
	GetMediaActivities(ctx context.Context, mediaID uuid.UUID) ([]*model.UserActivity, error)
	GetUserActivityForMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) (*model.UserActivity, error)
	UpdateActivity(ctx context.Context, userID uuid.UUID, id uuid.UUID, input model.UpdateActivityInput) (*model.UserActivity, error)
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
	// BuildRecommendations computes and upserts backend recommendations for a user
	// using tag-graph similarity and friends' activity boost.
	BuildRecommendations(ctx context.Context, userID uuid.UUID) error
}

// FriendRepository defines operations for friend relationships
type FriendRepository interface {
	SendFriendRequest(ctx context.Context, fromID, toID uuid.UUID) (*model.FriendRequest, error)
	AcceptFriendRequest(ctx context.Context, requestID, acceptingUserID uuid.UUID) (*model.User, error)
	DeclineFriendRequest(ctx context.Context, requestID, decliningUserID uuid.UUID) (bool, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) (bool, error)
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*model.User, error)
	GetPendingFriendRequests(ctx context.Context, userID uuid.UUID) ([]*model.FriendRequest, error)
	GetSentFriendRequests(ctx context.Context, userID uuid.UUID) ([]*model.FriendRequest, error)
	SearchUsers(ctx context.Context, query string, excludeUserID uuid.UUID) ([]*model.User, error)
	GetFriendsActivity(ctx context.Context, userID uuid.UUID, limit int) ([]*model.UserActivity, error)
}

// Neo4jRepository implements the Repository interface using Neo4j
type Neo4jRepository struct {
	db       *Database
	metadata *metadata.Service
}

// NewNeo4jRepository creates a new Neo4j repository
func NewNeo4jRepository(db *Database) *Neo4jRepository {
	metadataService, err := metadata.NewService()
	if err != nil {
		log.Printf("Failed to initialize metadata service: %v", err)
		// Don't fail - metadata is optional enhancement
		return &Neo4jRepository{db: db, metadata: nil}
	}

	log.Println("Metadata service initialized successfully")
	return &Neo4jRepository{
		db:       db,
		metadata: metadataService,
	}
}

// GetMetadata returns the metadata service
func (r *Neo4jRepository) GetMetadata() interface{} {
	return r.metadata
}
