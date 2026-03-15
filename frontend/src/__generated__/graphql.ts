import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Date: { input: string; output: string; }
  DateTime: { input: string; output: string; }
  UUID: { input: string; output: string; }
};

export type ActivityStatus = {
  __typename?: 'ActivityStatus';
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
};

export type Book = Media & {
  __typename?: 'Book';
  authors: Array<Creator>;
  averageRating?: Maybe<Scalars['Float']['output']>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  description?: Maybe<Scalars['String']['output']>;
  id: Scalars['UUID']['output'];
  isbn?: Maybe<Scalars['String']['output']>;
  myActivity?: Maybe<UserActivity>;
  pages?: Maybe<Scalars['Int']['output']>;
  platforms: Array<Platform>;
  publisher?: Maybe<Scalars['String']['output']>;
  publishers: Array<Scalars['String']['output']>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  searchDepth: Scalars['Int']['output'];
  subjects: Array<Tag>;
  tags: Array<Tag>;
  title: Scalars['String']['output'];
};


export type BookRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type CastAndCrewResult = {
  __typename?: 'CastAndCrewResult';
  cast: Array<Person>;
  castCredits: Array<PersonCredit>;
  crew: Array<Person>;
  crewCredits: Array<CrewCredit>;
};

export type CreateActivityInput = {
  finishedAt?: InputMaybe<Scalars['DateTime']['input']>;
  mediaId: Scalars['UUID']['input'];
  rating?: InputMaybe<Scalars['Float']['input']>;
  review?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['DateTime']['input']>;
  statusId: Scalars['Int']['input'];
};

export type CreateBookInput = {
  authors?: InputMaybe<Array<Scalars['String']['input']>>;
  coverUrl?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  isbn?: InputMaybe<Scalars['String']['input']>;
  pages?: InputMaybe<Scalars['Int']['input']>;
  publisher?: InputMaybe<Scalars['String']['input']>;
  publishers?: InputMaybe<Array<Scalars['String']['input']>>;
  releaseDate?: InputMaybe<Scalars['Date']['input']>;
  searchDepth?: InputMaybe<Scalars['Int']['input']>;
  subjects?: InputMaybe<Array<Scalars['String']['input']>>;
  title: Scalars['String']['input'];
};

export type CreateGameInput = {
  coverUrl?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  esrbRating?: InputMaybe<Scalars['String']['input']>;
  externalId?: InputMaybe<Scalars['String']['input']>;
  franchises?: InputMaybe<Array<Scalars['String']['input']>>;
  gameModes?: InputMaybe<Array<Scalars['String']['input']>>;
  genre?: InputMaybe<Array<Scalars['String']['input']>>;
  keywords?: InputMaybe<Array<Scalars['String']['input']>>;
  multiplayer?: InputMaybe<Scalars['Boolean']['input']>;
  perspectives?: InputMaybe<Array<Scalars['String']['input']>>;
  platforms?: InputMaybe<Array<Scalars['String']['input']>>;
  releaseDate?: InputMaybe<Scalars['Date']['input']>;
  searchDepth?: InputMaybe<Scalars['Int']['input']>;
  themes?: InputMaybe<Array<Scalars['String']['input']>>;
  title: Scalars['String']['input'];
};

export type CreateMovieInput = {
  boxOffice?: InputMaybe<Scalars['Int']['input']>;
  budget?: InputMaybe<Scalars['Int']['input']>;
  cast?: InputMaybe<Array<Scalars['String']['input']>>;
  coverUrl?: InputMaybe<Scalars['String']['input']>;
  crew?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  externalId?: InputMaybe<Scalars['String']['input']>;
  genres?: InputMaybe<Array<Scalars['String']['input']>>;
  maxConnections?: InputMaybe<Scalars['Int']['input']>;
  productionCompanies?: InputMaybe<Array<Scalars['String']['input']>>;
  releaseDate?: InputMaybe<Scalars['Date']['input']>;
  runtime?: InputMaybe<Scalars['Int']['input']>;
  searchDepth?: InputMaybe<Scalars['Int']['input']>;
  title: Scalars['String']['input'];
};

export type CreateMusicAlbumInput = {
  coverUrl?: InputMaybe<Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  duration?: InputMaybe<Scalars['Int']['input']>;
  label?: InputMaybe<Scalars['String']['input']>;
  releaseDate?: InputMaybe<Scalars['Date']['input']>;
  searchDepth?: InputMaybe<Scalars['Int']['input']>;
  title: Scalars['String']['input'];
  trackCount?: InputMaybe<Scalars['Int']['input']>;
};

export type CreateTvShowInput = {
  cast?: InputMaybe<Array<Scalars['String']['input']>>;
  coverUrl?: InputMaybe<Scalars['String']['input']>;
  crew?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  episodes?: InputMaybe<Scalars['Int']['input']>;
  externalId?: InputMaybe<Scalars['String']['input']>;
  genres?: InputMaybe<Array<Scalars['String']['input']>>;
  maxConnections?: InputMaybe<Scalars['Int']['input']>;
  productionCompanies?: InputMaybe<Array<Scalars['String']['input']>>;
  releaseDate?: InputMaybe<Scalars['Date']['input']>;
  searchDepth?: InputMaybe<Scalars['Int']['input']>;
  seasons?: InputMaybe<Scalars['Int']['input']>;
  status?: InputMaybe<Scalars['String']['input']>;
  title: Scalars['String']['input'];
};

export type CreateUserInput = {
  authProvider?: InputMaybe<Scalars['String']['input']>;
  authSubject?: InputMaybe<Scalars['String']['input']>;
  avatarUrl?: InputMaybe<Scalars['String']['input']>;
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type Creator = {
  __typename?: 'Creator';
  id: Scalars['UUID']['output'];
  mediaItems: Array<Media>;
  name: Scalars['String']['output'];
  role: CreatorRole;
};

export type CreatorRole = {
  __typename?: 'CreatorRole';
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
};

export type CrewCredit = {
  __typename?: 'CrewCredit';
  department?: Maybe<Scalars['String']['output']>;
  job?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  person: Person;
};

export type FriendRequest = {
  __typename?: 'FriendRequest';
  createdAt: Scalars['DateTime']['output'];
  from: User;
  id: Scalars['UUID']['output'];
  to: User;
};

export type Game = Media & {
  __typename?: 'Game';
  averageRating?: Maybe<Scalars['Float']['output']>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  description?: Maybe<Scalars['String']['output']>;
  esrbRating?: Maybe<Scalars['String']['output']>;
  franchises: Array<Scalars['String']['output']>;
  gameModes: Array<Scalars['String']['output']>;
  genre: Array<Scalars['String']['output']>;
  id: Scalars['UUID']['output'];
  keywords: Array<Scalars['String']['output']>;
  multiplayer?: Maybe<Scalars['Boolean']['output']>;
  myActivity?: Maybe<UserActivity>;
  perspectives: Array<Scalars['String']['output']>;
  platforms: Array<Platform>;
  platformsList: Array<Scalars['String']['output']>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  searchDepth: Scalars['Int']['output'];
  tags: Array<Tag>;
  themes: Array<Scalars['String']['output']>;
  title: Scalars['String']['output'];
};


export type GameRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type Genre = {
  __typename?: 'Genre';
  id: Scalars['UUID']['output'];
  movies: Array<Movie>;
  name: Scalars['String']['output'];
};

export type Media = {
  averageRating?: Maybe<Scalars['Float']['output']>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  description?: Maybe<Scalars['String']['output']>;
  id: Scalars['UUID']['output'];
  myActivity?: Maybe<UserActivity>;
  platforms: Array<Platform>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  searchDepth: Scalars['Int']['output'];
  tags: Array<Tag>;
  title: Scalars['String']['output'];
};


export type MediaRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type MediaSuggestion = {
  __typename?: 'MediaSuggestion';
  externalId?: Maybe<Scalars['String']['output']>;
  imageUrl?: Maybe<Scalars['String']['output']>;
  subtitle?: Maybe<Scalars['String']['output']>;
  title: Scalars['String']['output'];
  year?: Maybe<Scalars['Int']['output']>;
};

export type MediaType =
  | 'BOOK'
  | 'GAME'
  | 'MOVIE'
  | 'MUSIC'
  | 'TV';

export type Movie = Media & {
  __typename?: 'Movie';
  averageRating?: Maybe<Scalars['Float']['output']>;
  boxOffice?: Maybe<Scalars['Int']['output']>;
  budget?: Maybe<Scalars['Int']['output']>;
  cast: Array<Person>;
  castCredits: Array<PersonCredit>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  crew: Array<Person>;
  crewCredits: Array<CrewCredit>;
  description?: Maybe<Scalars['String']['output']>;
  genres: Array<Genre>;
  id: Scalars['UUID']['output'];
  myActivity?: Maybe<UserActivity>;
  platforms: Array<Platform>;
  productionCompanies: Array<ProductionCompany>;
  productionCountries: Array<ProductionCountry>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  runtime?: Maybe<Scalars['Int']['output']>;
  searchDepth: Scalars['Int']['output'];
  tags: Array<Tag>;
  title: Scalars['String']['output'];
};


export type MovieRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type MusicAlbum = Media & {
  __typename?: 'MusicAlbum';
  averageRating?: Maybe<Scalars['Float']['output']>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  description?: Maybe<Scalars['String']['output']>;
  duration?: Maybe<Scalars['Int']['output']>;
  id: Scalars['UUID']['output'];
  label?: Maybe<Scalars['String']['output']>;
  myActivity?: Maybe<UserActivity>;
  platforms: Array<Platform>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  searchDepth: Scalars['Int']['output'];
  tags: Array<Tag>;
  title: Scalars['String']['output'];
  trackCount?: Maybe<Scalars['Int']['output']>;
};


export type MusicAlbumRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  acceptFriendRequest: User;
  addToFavorites: Scalars['Boolean']['output'];
  createActivity: UserActivity;
  createBook: Book;
  createGame: Game;
  createMovie: Movie;
  createMusicAlbum: MusicAlbum;
  createTVShow: TvShow;
  createUser: User;
  declineFriendRequest: Scalars['Boolean']['output'];
  deleteUser: Scalars['Boolean']['output'];
  rateMedia: Rating;
  removeFriend: Scalars['Boolean']['output'];
  sendFriendRequest: FriendRequest;
  updateActivity: UserActivity;
  updateUser: User;
};


export type MutationAcceptFriendRequestArgs = {
  requestID: Scalars['UUID']['input'];
};


export type MutationAddToFavoritesArgs = {
  mediaId: Scalars['UUID']['input'];
};


export type MutationCreateActivityArgs = {
  input: CreateActivityInput;
};


export type MutationCreateBookArgs = {
  input: CreateBookInput;
};


export type MutationCreateGameArgs = {
  input: CreateGameInput;
};


export type MutationCreateMovieArgs = {
  input: CreateMovieInput;
};


export type MutationCreateMusicAlbumArgs = {
  input: CreateMusicAlbumInput;
};


export type MutationCreateTvShowArgs = {
  input: CreateTvShowInput;
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


export type MutationDeclineFriendRequestArgs = {
  requestID: Scalars['UUID']['input'];
};


export type MutationDeleteUserArgs = {
  id: Scalars['UUID']['input'];
};


export type MutationRateMediaArgs = {
  mediaId: Scalars['UUID']['input'];
  score: Scalars['Float']['input'];
};


export type MutationRemoveFriendArgs = {
  friendID: Scalars['UUID']['input'];
};


export type MutationSendFriendRequestArgs = {
  toUserID: Scalars['UUID']['input'];
};


export type MutationUpdateActivityArgs = {
  id: Scalars['UUID']['input'];
  input: UpdateActivityInput;
};


export type MutationUpdateUserArgs = {
  id: Scalars['UUID']['input'];
  input: UpdateUserInput;
};

export type Person = {
  __typename?: 'Person';
  actedIn: Array<Movie>;
  crewOn: Array<Movie>;
  externalID?: Maybe<Scalars['String']['output']>;
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
};

export type PersonCredit = {
  __typename?: 'PersonCredit';
  character?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  order?: Maybe<Scalars['Int']['output']>;
  person: Person;
};

export type Platform = {
  __typename?: 'Platform';
  baseUrl?: Maybe<Scalars['String']['output']>;
  id: Scalars['UUID']['output'];
  mediaItems: Array<Media>;
  name: Scalars['String']['output'];
};

export type ProductionCompany = {
  __typename?: 'ProductionCompany';
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  produced: Array<Movie>;
};

export type ProductionCountry = {
  __typename?: 'ProductionCountry';
  id: Scalars['UUID']['output'];
  movies: Array<Movie>;
  name: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  allMedia: Array<Media>;
  autocompleteMedia: Array<MediaSuggestion>;
  books: Array<Book>;
  castAndCrew: CastAndCrewResult;
  friendsActivity: Array<UserActivity>;
  games: Array<Game>;
  getRecommendations: Array<Recommendation>;
  me?: Maybe<User>;
  media?: Maybe<Media>;
  movies: Array<Movie>;
  musicAlbums: Array<MusicAlbum>;
  recursiveSearchStatus: SearchStatus;
  searchUsers: Array<User>;
  tvShows: Array<TvShow>;
  user?: Maybe<User>;
  users: Array<User>;
};


export type QueryAutocompleteMediaArgs = {
  query: Scalars['String']['input'];
  type: MediaType;
};


export type QueryBooksArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryCastAndCrewArgs = {
  mediaID: Scalars['UUID']['input'];
};


export type QueryFriendsActivityArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryGamesArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryMediaArgs = {
  id: Scalars['UUID']['input'];
};


export type QueryMoviesArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryMusicAlbumsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryRecursiveSearchStatusArgs = {
  mediaId: Scalars['UUID']['input'];
};


export type QuerySearchUsersArgs = {
  query: Scalars['String']['input'];
};


export type QueryTvShowsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryUserArgs = {
  id: Scalars['UUID']['input'];
};

export type Rating = {
  __typename?: 'Rating';
  media: Media;
  ratedAt: Scalars['DateTime']['output'];
  score: Scalars['Float']['output'];
  user: User;
};

export type Recommendation = {
  __typename?: 'Recommendation';
  id: Scalars['UUID']['output'];
  media: Media;
  recommender?: Maybe<User>;
  score?: Maybe<Scalars['Float']['output']>;
  source?: Maybe<Scalars['String']['output']>;
  user: User;
};

export type SearchState =
  | 'COMPLETED'
  | 'IDLE'
  | 'RUNNING';

export type SearchStatus = {
  __typename?: 'SearchStatus';
  completedAt?: Maybe<Scalars['DateTime']['output']>;
  state: SearchState;
};

export type TvShow = Media & {
  __typename?: 'TVShow';
  averageRating?: Maybe<Scalars['Float']['output']>;
  cast: Array<Person>;
  castCredits: Array<PersonCredit>;
  coverUrl?: Maybe<Scalars['String']['output']>;
  creators: Array<Creator>;
  crew: Array<Person>;
  crewCredits: Array<CrewCredit>;
  description?: Maybe<Scalars['String']['output']>;
  episodes?: Maybe<Scalars['Int']['output']>;
  genres: Array<Genre>;
  id: Scalars['UUID']['output'];
  myActivity?: Maybe<UserActivity>;
  platforms: Array<Platform>;
  productionCompanies: Array<ProductionCompany>;
  productionCountries: Array<ProductionCountry>;
  ratings: Array<Rating>;
  relatedMedia: Array<Media>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  searchDepth: Scalars['Int']['output'];
  seasons?: Maybe<Scalars['Int']['output']>;
  status?: Maybe<Scalars['String']['output']>;
  tags: Array<Tag>;
  title: Scalars['String']['output'];
};


export type TvShowRelatedMediaArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
};

export type Tag = {
  __typename?: 'Tag';
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type UpdateActivityInput = {
  finishedAt?: InputMaybe<Scalars['DateTime']['input']>;
  rating?: InputMaybe<Scalars['Float']['input']>;
  review?: InputMaybe<Scalars['String']['input']>;
  startedAt?: InputMaybe<Scalars['DateTime']['input']>;
  statusId?: InputMaybe<Scalars['Int']['input']>;
};

export type UpdateUserInput = {
  avatarUrl?: InputMaybe<Scalars['String']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type User = {
  __typename?: 'User';
  activities: Array<UserActivity>;
  authProvider?: Maybe<Scalars['String']['output']>;
  authSubject?: Maybe<Scalars['String']['output']>;
  avatarUrl?: Maybe<Scalars['String']['output']>;
  email: Scalars['String']['output'];
  favorites: Array<Media>;
  friends: Array<User>;
  id: Scalars['UUID']['output'];
  name: Scalars['String']['output'];
  pendingFriendRequests: Array<FriendRequest>;
  ratings: Array<Rating>;
  recommendations: Array<Recommendation>;
  sentFriendRequests: Array<FriendRequest>;
};

export type UserActivity = {
  __typename?: 'UserActivity';
  finishedAt?: Maybe<Scalars['DateTime']['output']>;
  id: Scalars['UUID']['output'];
  media: Media;
  rating?: Maybe<Scalars['Float']['output']>;
  review?: Maybe<Scalars['String']['output']>;
  sourcePlatform?: Maybe<Platform>;
  startedAt?: Maybe<Scalars['DateTime']['output']>;
  status: ActivityStatus;
  user: User;
};

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename?: 'Query', me?: { __typename?: 'User', id: string, name: string, email: string, authProvider?: string | null, authSubject?: string | null, avatarUrl?: string | null } | null };

export type MeActivitiesQueryVariables = Exact<{ [key: string]: never; }>;


export type MeActivitiesQuery = { __typename?: 'Query', me?: { __typename?: 'User', id: string, activities: Array<{ __typename?: 'UserActivity', id: string, media:
        | { __typename: 'Book', pages?: number | null, id: string, title: string, description?: string | null, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, subjects: Array<{ __typename?: 'Tag', name: string }> }
        | { __typename: 'Game', genre: Array<string>, themes: Array<string>, keywords: Array<string>, gameModes: Array<string>, perspectives: Array<string>, franchises: Array<string>, platformsList: Array<string>, id: string, title: string, description?: string | null, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null }
        | { __typename: 'Movie', runtime?: number | null, id: string, title: string, description?: string | null, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, genres: Array<{ __typename?: 'Genre', name: string }> }
        | { __typename: 'MusicAlbum', id: string, title: string, description?: string | null, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null }
        | { __typename: 'TVShow', id: string, title: string, description?: string | null, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, genres: Array<{ __typename?: 'Genre', name: string }> }
       }> } | null };

export type GetHomeMediaQueryVariables = Exact<{ [key: string]: never; }>;


export type GetHomeMediaQuery = { __typename?: 'Query', allMedia: Array<
    | { __typename: 'Book', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, subjects: Array<{ __typename?: 'Tag', name: string }>, authors: Array<{ __typename?: 'Creator', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
    | { __typename: 'Game', genre: Array<string>, themes: Array<string>, keywords: Array<string>, gameModes: Array<string>, perspectives: Array<string>, franchises: Array<string>, platformsList: Array<string>, id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
    | { __typename: 'Movie', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
    | { __typename: 'MusicAlbum', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
    | { __typename: 'TVShow', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
  >, me?: { __typename?: 'User', id: string, activities: Array<{ __typename?: 'UserActivity', id: string, media:
        | { __typename: 'Book', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, subjects: Array<{ __typename?: 'Tag', name: string }>, authors: Array<{ __typename?: 'Creator', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
        | { __typename: 'Game', genre: Array<string>, themes: Array<string>, keywords: Array<string>, gameModes: Array<string>, perspectives: Array<string>, franchises: Array<string>, platformsList: Array<string>, id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
        | { __typename: 'Movie', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
        | { __typename: 'MusicAlbum', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
        | { __typename: 'TVShow', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
       }> } | null };

export type GetMediaDetailsQueryVariables = Exact<{
  id: Scalars['UUID']['input'];
}>;


export type GetMediaDetailsQuery = { __typename?: 'Query', media?:
    | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, description?: string | null, coverUrl?: string | null, averageRating?: number | null, authors: Array<{ __typename?: 'Creator', id: string, name: string }>, relatedMedia: Array<
        | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Game', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Movie', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'MusicAlbum', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'TVShow', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
      >, creators: Array<{ __typename?: 'Creator', id: string, name: string }>, tags: Array<{ __typename?: 'Tag', id: string, name: string }>, myActivity?: { __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, status: { __typename?: 'ActivityStatus', id: number, name: string } } | null }
    | { __typename: 'Game', genre: Array<string>, themes: Array<string>, keywords: Array<string>, gameModes: Array<string>, perspectives: Array<string>, franchises: Array<string>, platformsList: Array<string>, id: string, title: string, releaseDate?: string | null, description?: string | null, coverUrl?: string | null, averageRating?: number | null, relatedMedia: Array<
        | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Game', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Movie', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'MusicAlbum', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'TVShow', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
      >, creators: Array<{ __typename?: 'Creator', id: string, name: string }>, tags: Array<{ __typename?: 'Tag', id: string, name: string }>, myActivity?: { __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, status: { __typename?: 'ActivityStatus', id: number, name: string } } | null }
    | { __typename: 'Movie', runtime?: number | null, id: string, title: string, releaseDate?: string | null, description?: string | null, coverUrl?: string | null, averageRating?: number | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, crew: Array<{ __typename?: 'Person', id: string, name: string }>, relatedMedia: Array<
        | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Game', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Movie', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'MusicAlbum', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'TVShow', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
      >, creators: Array<{ __typename?: 'Creator', id: string, name: string }>, tags: Array<{ __typename?: 'Tag', id: string, name: string }>, myActivity?: { __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, status: { __typename?: 'ActivityStatus', id: number, name: string } } | null }
    | { __typename: 'MusicAlbum', label?: string | null, id: string, title: string, releaseDate?: string | null, description?: string | null, coverUrl?: string | null, averageRating?: number | null, relatedMedia: Array<
        | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Game', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Movie', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'MusicAlbum', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'TVShow', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
      >, creators: Array<{ __typename?: 'Creator', id: string, name: string }>, tags: Array<{ __typename?: 'Tag', id: string, name: string }>, myActivity?: { __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, status: { __typename?: 'ActivityStatus', id: number, name: string } } | null }
    | { __typename: 'TVShow', seasons?: number | null, episodes?: number | null, id: string, title: string, releaseDate?: string | null, description?: string | null, coverUrl?: string | null, averageRating?: number | null, genres: Array<{ __typename?: 'Genre', name: string }>, cast: Array<{ __typename?: 'Person', id: string, name: string }>, crew: Array<{ __typename?: 'Person', id: string, name: string }>, relatedMedia: Array<
        | { __typename: 'Book', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Game', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'Movie', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'MusicAlbum', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
        | { __typename: 'TVShow', id: string, title: string, releaseDate?: string | null, coverUrl?: string | null, averageRating?: number | null, description?: string | null }
      >, creators: Array<{ __typename?: 'Creator', id: string, name: string }>, tags: Array<{ __typename?: 'Tag', id: string, name: string }>, myActivity?: { __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, status: { __typename?: 'ActivityStatus', id: number, name: string } } | null }
   | null };

export type AutocompleteMediaQueryVariables = Exact<{
  type: MediaType;
  query: Scalars['String']['input'];
}>;


export type AutocompleteMediaQuery = { __typename?: 'Query', autocompleteMedia: Array<{ __typename?: 'MediaSuggestion', title: string, year?: number | null, externalId?: string | null, imageUrl?: string | null, subtitle?: string | null }> };

export type RecursiveSearchStatusQueryVariables = Exact<{
  mediaId: Scalars['UUID']['input'];
}>;


export type RecursiveSearchStatusQuery = { __typename?: 'Query', recursiveSearchStatus: { __typename?: 'SearchStatus', state: SearchState, completedAt?: string | null } };

export type UpdateUserMutationVariables = Exact<{
  id: Scalars['UUID']['input'];
  input: UpdateUserInput;
}>;


export type UpdateUserMutation = { __typename?: 'Mutation', updateUser: { __typename?: 'User', id: string, name: string, email: string, avatarUrl?: string | null } };

export type SearchUsersQueryVariables = Exact<{
  query: Scalars['String']['input'];
}>;


export type SearchUsersQuery = { __typename?: 'Query', searchUsers: Array<{ __typename?: 'User', id: string, name: string, avatarUrl?: string | null }> };

export type MeFriendsQueryVariables = Exact<{ [key: string]: never; }>;


export type MeFriendsQuery = { __typename?: 'Query', me?: { __typename?: 'User', id: string, friends: Array<{ __typename?: 'User', id: string, name: string, avatarUrl?: string | null }>, pendingFriendRequests: Array<{ __typename?: 'FriendRequest', id: string, createdAt: string, from: { __typename?: 'User', id: string, name: string, avatarUrl?: string | null } }>, sentFriendRequests: Array<{ __typename?: 'FriendRequest', id: string, createdAt: string, to: { __typename?: 'User', id: string, name: string, avatarUrl?: string | null } }> } | null };

export type FriendsActivityQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
}>;


export type FriendsActivityQuery = { __typename?: 'Query', friendsActivity: Array<{ __typename?: 'UserActivity', id: string, rating?: number | null, review?: string | null, startedAt?: string | null, finishedAt?: string | null, user: { __typename?: 'User', id: string, name: string, avatarUrl?: string | null }, media:
      | { __typename: 'Book', id: string, title: string, coverUrl?: string | null, averageRating?: number | null }
      | { __typename: 'Game', id: string, title: string, coverUrl?: string | null, averageRating?: number | null }
      | { __typename: 'Movie', id: string, title: string, coverUrl?: string | null, averageRating?: number | null }
      | { __typename: 'MusicAlbum', id: string, title: string, coverUrl?: string | null, averageRating?: number | null }
      | { __typename: 'TVShow', id: string, title: string, coverUrl?: string | null, averageRating?: number | null }
    , status: { __typename?: 'ActivityStatus', id: number, name: string } }> };

export type GetRecommendationsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetRecommendationsQuery = { __typename?: 'Query', getRecommendations: Array<{ __typename?: 'Recommendation', id: string, score?: number | null, source?: string | null, media:
      | { __typename: 'Book', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, subjects: Array<{ __typename?: 'Tag', name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
      | { __typename: 'Game', genre: Array<string>, themes: Array<string>, id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
      | { __typename: 'Movie', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
      | { __typename: 'MusicAlbum', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
      | { __typename: 'TVShow', id: string, title: string, coverUrl?: string | null, averageRating?: number | null, description?: string | null, releaseDate?: string | null, genres: Array<{ __typename?: 'Genre', name: string }>, creators: Array<{ __typename?: 'Creator', id: string, name: string }> }
     }> };

export type SendFriendRequestMutationVariables = Exact<{
  toUserID: Scalars['UUID']['input'];
}>;


export type SendFriendRequestMutation = { __typename?: 'Mutation', sendFriendRequest: { __typename?: 'FriendRequest', id: string, createdAt: string, to: { __typename?: 'User', id: string, name: string, avatarUrl?: string | null } } };

export type AcceptFriendRequestMutationVariables = Exact<{
  requestID: Scalars['UUID']['input'];
}>;


export type AcceptFriendRequestMutation = { __typename?: 'Mutation', acceptFriendRequest: { __typename?: 'User', id: string, name: string, avatarUrl?: string | null } };

export type DeclineFriendRequestMutationVariables = Exact<{
  requestID: Scalars['UUID']['input'];
}>;


export type DeclineFriendRequestMutation = { __typename?: 'Mutation', declineFriendRequest: boolean };

export type RemoveFriendMutationVariables = Exact<{
  friendID: Scalars['UUID']['input'];
}>;


export type RemoveFriendMutation = { __typename?: 'Mutation', removeFriend: boolean };


export const MeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"authProvider"}},{"kind":"Field","name":{"kind":"Name","value":"authSubject"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]} as unknown as DocumentNode<MeQuery, MeQueryVariables>;
export const MeActivitiesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"MeActivities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"activities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"media"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Movie"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"runtime"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"TVShow"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Book"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"subjects"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pages"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Game"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genre"}},{"kind":"Field","name":{"kind":"Name","value":"themes"}},{"kind":"Field","name":{"kind":"Name","value":"keywords"}},{"kind":"Field","name":{"kind":"Name","value":"gameModes"}},{"kind":"Field","name":{"kind":"Name","value":"perspectives"}},{"kind":"Field","name":{"kind":"Name","value":"franchises"}},{"kind":"Field","name":{"kind":"Name","value":"platformsList"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<MeActivitiesQuery, MeActivitiesQueryVariables>;
export const GetHomeMediaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetHomeMedia"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"allMedia"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"creators"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Movie"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"TVShow"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Book"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"subjects"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"authors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Game"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genre"}},{"kind":"Field","name":{"kind":"Name","value":"themes"}},{"kind":"Field","name":{"kind":"Name","value":"keywords"}},{"kind":"Field","name":{"kind":"Name","value":"gameModes"}},{"kind":"Field","name":{"kind":"Name","value":"perspectives"}},{"kind":"Field","name":{"kind":"Name","value":"franchises"}},{"kind":"Field","name":{"kind":"Name","value":"platformsList"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"activities"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"media"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"creators"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Movie"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"TVShow"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Book"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"subjects"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"authors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Game"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genre"}},{"kind":"Field","name":{"kind":"Name","value":"themes"}},{"kind":"Field","name":{"kind":"Name","value":"keywords"}},{"kind":"Field","name":{"kind":"Name","value":"gameModes"}},{"kind":"Field","name":{"kind":"Name","value":"perspectives"}},{"kind":"Field","name":{"kind":"Name","value":"franchises"}},{"kind":"Field","name":{"kind":"Name","value":"platformsList"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetHomeMediaQuery, GetHomeMediaQueryVariables>;
export const GetMediaDetailsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetMediaDetails"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"media"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"creators"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"tags"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"myActivity"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"rating"}},{"kind":"Field","name":{"kind":"Name","value":"review"}},{"kind":"Field","name":{"kind":"Name","value":"status"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"startedAt"}},{"kind":"Field","name":{"kind":"Name","value":"finishedAt"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Movie"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"runtime"}},{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"crew"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"relatedMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"IntValue","value":"12"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"TVShow"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"seasons"}},{"kind":"Field","name":{"kind":"Name","value":"episodes"}},{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"cast"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"crew"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"relatedMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"IntValue","value":"12"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Book"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authors"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"relatedMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"IntValue","value":"12"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Game"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genre"}},{"kind":"Field","name":{"kind":"Name","value":"themes"}},{"kind":"Field","name":{"kind":"Name","value":"keywords"}},{"kind":"Field","name":{"kind":"Name","value":"gameModes"}},{"kind":"Field","name":{"kind":"Name","value":"perspectives"}},{"kind":"Field","name":{"kind":"Name","value":"franchises"}},{"kind":"Field","name":{"kind":"Name","value":"platformsList"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"IntValue","value":"12"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"MusicAlbum"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"label"}},{"kind":"Field","name":{"kind":"Name","value":"relatedMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"IntValue","value":"12"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetMediaDetailsQuery, GetMediaDetailsQueryVariables>;
export const AutocompleteMediaDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"AutocompleteMedia"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"type"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"MediaType"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"autocompleteMedia"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"type"},"value":{"kind":"Variable","name":{"kind":"Name","value":"type"}}},{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"year"}},{"kind":"Field","name":{"kind":"Name","value":"externalId"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"subtitle"}}]}}]}}]} as unknown as DocumentNode<AutocompleteMediaQuery, AutocompleteMediaQueryVariables>;
export const RecursiveSearchStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"RecursiveSearchStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"mediaId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"recursiveSearchStatus"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"mediaId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"mediaId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"state"}},{"kind":"Field","name":{"kind":"Name","value":"completedAt"}}]}}]}}]} as unknown as DocumentNode<RecursiveSearchStatusQuery, RecursiveSearchStatusQueryVariables>;
export const UpdateUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateUserInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]} as unknown as DocumentNode<UpdateUserMutation, UpdateUserMutationVariables>;
export const SearchUsersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"SearchUsers"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"query"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"searchUsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"query"},"value":{"kind":"Variable","name":{"kind":"Name","value":"query"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]} as unknown as DocumentNode<SearchUsersQuery, SearchUsersQueryVariables>;
export const MeFriendsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"MeFriends"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"friends"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pendingFriendRequests"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"from"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"sentFriendRequests"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"to"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]}}]}}]} as unknown as DocumentNode<MeFriendsQuery, MeFriendsQueryVariables>;
export const FriendsActivityDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FriendsActivity"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"limit"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"friendsActivity"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"Variable","name":{"kind":"Name","value":"limit"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}},{"kind":"Field","name":{"kind":"Name","value":"media"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}}]}},{"kind":"Field","name":{"kind":"Name","value":"status"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"rating"}},{"kind":"Field","name":{"kind":"Name","value":"review"}},{"kind":"Field","name":{"kind":"Name","value":"startedAt"}},{"kind":"Field","name":{"kind":"Name","value":"finishedAt"}}]}}]}}]} as unknown as DocumentNode<FriendsActivityQuery, FriendsActivityQueryVariables>;
export const GetRecommendationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetRecommendations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"getRecommendations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"score"}},{"kind":"Field","name":{"kind":"Name","value":"source"}},{"kind":"Field","name":{"kind":"Name","value":"media"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"coverUrl"}},{"kind":"Field","name":{"kind":"Name","value":"averageRating"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"releaseDate"}},{"kind":"Field","name":{"kind":"Name","value":"creators"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Movie"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"TVShow"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genres"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Book"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"subjects"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Game"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"genre"}},{"kind":"Field","name":{"kind":"Name","value":"themes"}}]}}]}}]}}]}}]} as unknown as DocumentNode<GetRecommendationsQuery, GetRecommendationsQueryVariables>;
export const SendFriendRequestDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SendFriendRequest"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"toUserID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"sendFriendRequest"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"toUserID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"toUserID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"to"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]}}]} as unknown as DocumentNode<SendFriendRequestMutation, SendFriendRequestMutationVariables>;
export const AcceptFriendRequestDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"AcceptFriendRequest"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"requestID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"acceptFriendRequest"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"requestID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"requestID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"avatarUrl"}}]}}]}}]} as unknown as DocumentNode<AcceptFriendRequestMutation, AcceptFriendRequestMutationVariables>;
export const DeclineFriendRequestDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeclineFriendRequest"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"requestID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"declineFriendRequest"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"requestID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"requestID"}}}]}]}}]} as unknown as DocumentNode<DeclineFriendRequestMutation, DeclineFriendRequestMutationVariables>;
export const RemoveFriendDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RemoveFriend"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"friendID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"removeFriend"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"friendID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"friendID"}}}]}]}}]} as unknown as DocumentNode<RemoveFriendMutation, RemoveFriendMutationVariables>;