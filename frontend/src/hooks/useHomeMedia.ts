import { useQuery } from '@apollo/client/react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { GET_RECOMMENDATIONS_QUERY } from '../lib/graphql';
import { GetRecommendationsQuery } from '../__generated__/graphql';
import { logError, logInfo } from '../lib/logger';
import { useAppStateRefetch } from './useAppStateRefetch';

const PAGE_SIZE = 12;

export function useHomeMedia(limit: number = PAGE_SIZE) {
  const pageSize = normalizePageSize(limit);

  const { data, loading, error, refetch } = useQuery<GetRecommendationsQuery>(
    GET_RECOMMENDATIONS_QUERY,
    {
      fetchPolicy: 'cache-and-network',
      nextFetchPolicy: 'cache-first',
      errorPolicy: 'all',
    }
  );

  const [visibleCount, setVisibleCount] = useState(pageSize);

  useEffect(() => {
    if (error) {
      logError('useHomeMedia error:', error);
      if (error.message?.includes('401') || error.message?.includes('Unauthorized')) {
        logError('Authentication error detected - token may be invalid');
      }
    }
    if (data) {
      logInfo('useHomeMedia data received:', {
        recommendationsCount: data.getRecommendations?.length ?? 0,
      });
    }
  }, [data, error]);

  const mediaItems = useMemo(() => {
    const recs = data?.getRecommendations ?? [];
    return recs
      .filter(rec => rec.media != null)
      .map(rec => {
        const item = rec.media!;

        let genre: string[];
        if ('genres' in item && item.genres) {
          genre = (item.genres as { name: string }[]).map(g => g.name);
        } else if ('subjects' in item && item.subjects) {
          genre = (item.subjects as { name: string }[]).map(s => s.name);
        } else if ('genre' in item && item.genre) {
          genre = item.genre as string[];
        } else {
          genre = [];
        }

        let type: string;
        if (item.__typename === 'TVShow') {
          type = 'tv';
        } else if (item.__typename === 'Book') {
          type = 'book';
        } else if (item.__typename === 'Game') {
          type = 'game';
        } else if (item.__typename === 'MusicAlbum') {
          type = 'music';
        } else {
          type = 'movie';
        }

        return {
          id: String(item.id),
          title: item.title ?? 'Untitled',
          image:
            item.coverUrl ||
            `https://placehold.co/400x600?text=${encodeURIComponent(item.title ?? 'Untitled')}`,
          rating: item.averageRating || 0,
          genre,
          year: item.releaseDate
            ? parseInt(item.releaseDate.substring(0, 4), 10)
            : new Date().getFullYear(),
          duration: undefined,
          description: item.description || '',
          type,
        };
      });
  }, [data]);

  const media = useMemo(() => mediaItems.slice(0, visibleCount), [mediaItems, visibleCount]);
  const hasMore = visibleCount < mediaItems.length;

  const loadMore = useCallback(() => {
    if (!hasMore) return;
    setVisibleCount(prev => Math.min(prev + pageSize, mediaItems.length));
  }, [hasMore, mediaItems.length, pageSize]);

  const refresh = useCallback(async () => {
    setVisibleCount(pageSize);
    try {
      await refetch();
    } catch (err: any) {
      logError('Error refreshing media:', err);

      // If it's an Apollo invariant error, it's likely due to cache clearing during account switch
      if (err?.message?.includes('Invariant Violation') || err?.message?.includes('clearStore')) {
        logInfo('Refresh failed due to cache clearing - this is expected during account switching');
        return;
      }

      // If it's an authentication error, log it but don't throw
      if (
        err?.message?.includes('401') ||
        err?.message?.includes('Unauthorized') ||
        err?.networkError?.statusCode === 401
      ) {
        logError('Authentication error during media refresh - user may have been logged out');
      }
    }
  }, [pageSize, refetch]);

  useAppStateRefetch(refresh);

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
