export type MediaType = "movie" | "tv" | "book" | "game" | "music";

export interface Media {
  id: number | string;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: MediaType;
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

