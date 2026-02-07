import { useQuery } from "@apollo/client/react";
import { useState, useEffect } from "react";
import { Media } from "../types";
import { GET_MOVIES_QUERY } from "../../lib/graphql";

const PAGE_SIZE = 12;

interface MovieData {
  movies: {
    __typename: string;
    id: string;
    title: string;
    coverUrl: string;
    averageRating: number;
    genres: { name: string }[];
    description: string;
    releaseDate: string;
    runtime: number;
  }[];
}

interface MovieVars {
  limit: number;
  offset: number;
}

export function useMovies(limit: number = PAGE_SIZE) {
  const pageSize = normalizePageSize(limit);
  const [movies, setMovies] = useState<Media[]>([]);
  const [hasMore, setHasMore] = useState(true);

  const { data, loading, error, fetchMore, refetch } = useQuery<MovieData, MovieVars>(GET_MOVIES_QUERY, {
    variables: {
      limit: pageSize,
      offset: 0,
    },
  });

  const refresh = async () => {
    setHasMore(true);
    try {
      await refetch({
        limit: pageSize,
        offset: 0,
      });
    } catch (err) {
      console.error("Error refreshing movies:", err);
    }
  };

  const loadMore = async () => {
    if (!hasMore || loading) return;

    try {
      const currentLength = movies.length;
      await fetchMore({
        variables: {
          offset: currentLength,
          limit: pageSize,
        },
        updateQuery: (
          prev: MovieData,
          { fetchMoreResult }: { fetchMoreResult?: MovieData },
        ) => {
          if (!fetchMoreResult) return prev;
          if (fetchMoreResult.movies.length < pageSize) {
            setHasMore(false);
          }

          return {
            ...prev,
            movies: uniqueById([...prev.movies, ...fetchMoreResult.movies]),
          };
        },
      });
    } catch (err) {
      console.error("Error fetching more movies:", err);
    }
  };

  // Sync data with local state whenever query data changes
  useEffect(() => {
    if (data?.movies) {
      setMovies(transformMovies(uniqueById(data.movies)));
      if (data.movies.length < pageSize) {
        setHasMore(false);
      }
    }
  }, [data, pageSize]);

  return {
    movies,
    loading,
    error,
    loadMore,
    hasMore,
    refresh,
  };
}

function normalizePageSize(limit: number) {
  if (limit <= 0) return PAGE_SIZE;
  return Math.ceil(limit / 3) * 3;
}

function transformMovies(gqlMovies: any[]): Media[] {
  return gqlMovies.map((m: any) => ({
    id: m.id,
    title: m.title,
    image: m.coverUrl || "https://placehold.co/400x600?text=No+Image",
    rating: m.averageRating || 0,
    genre: m.genres ? m.genres.map((g: any) => g.name) : [],
    year: m.releaseDate ? parseInt(m.releaseDate.substring(0, 4)) : new Date().getFullYear(),
    duration: m.runtime ? `${Math.floor(m.runtime / 60)}h ${m.runtime % 60}m` : undefined,
    description: m.description || "",
    type: m.__typename === "TVShow" ? "tv" : m.__typename?.toLowerCase() ?? "movie",
  }));
}

function uniqueById<T extends { id: string }>(items: T[]): T[] {
  const seen = new Set<string>();
  const result: T[] = [];
  for (const item of items) {
    if (seen.has(item.id)) continue;
    seen.add(item.id);
    result.push(item);
  }
  return result;
}
