export type MediaType = "movie" | "tv" | "book" | "game" | "music";

export interface UserActivity {
  id: string;
  rating?: number | null;
  review?: string | null;
  status: {
    id: number;
    name: string;
  };
  startedAt?: string | null;
  finishedAt?: string | null;
}

export interface Media {
  id: string;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: MediaType;
  externalId?: string;
  isbn?: string;
  myActivity?: UserActivity | null;
}

export type CreateMediaResult = { id: string; title: string } | null;

export type GetMoviesQuery = {
  movies: {
    __typename: "Movie";
    id: string;
    title: string;
    coverUrl: string;
    averageRating: number;
    genres: { name: string }[];
    description: string;
  };
};

export type GetMoviesQueryVariables = Record<string, never>;
