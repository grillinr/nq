import { useQuery } from "@apollo/client/react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { GET_HOME_MEDIA_QUERY } from "../../lib/graphql";

const PAGE_SIZE = 12;
const PLACEHOLDER_IMAGE = "https://placehold.co/400x600?text=No+Image";

interface HomeMediaData {
  allMedia: {
    __typename: string;
    id: string;
    title: string;
    coverUrl?: string | null;
    averageRating?: number | null;
    description?: string | null;
    releaseDate?: string | null;
    genres?: { name: string }[] | null;
  }[];
}

export function useHomeMedia(limit: number = PAGE_SIZE) {
  const pageSize = normalizePageSize(limit);
  const { data, loading, error, refetch } = useQuery<HomeMediaData>(GET_HOME_MEDIA_QUERY, {
    fetchPolicy: "cache-and-network",
  });
  const [visibleCount, setVisibleCount] = useState(pageSize);

  const mediaItems = useMemo(() => {
    const allMedia = data?.allMedia ?? [];
    return allMedia
      .filter((item) => item.__typename === "Movie" || item.__typename === "TVShow")
      .map((item) => ({
        id: item.id,
        title: item.title ?? "Untitled",
        image: item.coverUrl || PLACEHOLDER_IMAGE,
        rating: item.averageRating || 0,
        genre: item.genres ? item.genres.map((g) => g.name) : [],
        year: item.releaseDate ? parseInt(item.releaseDate.substring(0, 4)) : new Date().getFullYear(),
        duration: undefined,
        description: item.description || "",
        type: item.__typename === "TVShow" ? "tv" : "movie",
      }))
      .sort((a, b) => {
        if (b.rating !== a.rating) return b.rating - a.rating;
        return a.title.localeCompare(b.title);
      });
  }, [data]);

  const media = useMemo(() => mediaItems.slice(0, visibleCount), [mediaItems, visibleCount]);
  const hasMore = visibleCount < mediaItems.length;

  const loadMore = useCallback(() => {
    if (!hasMore) return;
    setVisibleCount((prev) => Math.min(prev + pageSize, mediaItems.length));
  }, [hasMore, mediaItems.length, pageSize]);

  const refresh = useCallback(async () => {
    setVisibleCount(pageSize);
    try {
      await refetch();
    } catch (err) {
      console.error("Error refreshing media:", err);
    }
  }, [pageSize, refetch]);

  useEffect(() => {
    setVisibleCount(pageSize);
  }, [pageSize, data]);

  return {
    media,
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
