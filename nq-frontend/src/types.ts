export type MediaType = "movie" | "tv" | "book" | "game" | "music";

export interface Media {
  id: number;
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