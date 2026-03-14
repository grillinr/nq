import { useQuery } from '@apollo/client/react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { GET_HOME_MEDIA_QUERY } from '../lib/graphql';
import { capMediaCandidates, scoreMediaFromUser } from '../lib/graphScore';
import { logError, logInfo } from '../lib/logger';
import { useAppStateRefetch } from './useAppStateRefetch';

const PAGE_SIZE = 12;

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
    creators?: { id: string; name: string }[] | null;
    cast?: { id: string; name: string }[] | null;
    authors?: { id: string; name: string }[] | null;
    subjects?: { name: string }[] | null;
    genre?: string[] | null;
    themes?: string[] | null;
    keywords?: string[] | null;
    gameModes?: string[] | null;
    perspectives?: string[] | null;
    franchises?: string[] | null;
    platformsList?: string[] | null;
  }[];
  me?: {
    id: string;
    activities?: { id: string; media?: HomeMediaData['allMedia'][number] | null }[] | null;
  } | null;
}

export function useHomeMedia(limit: number = PAGE_SIZE) {
  const pageSize = normalizePageSize(limit);
  const { data, loading, error, refetch } = useQuery<HomeMediaData>(GET_HOME_MEDIA_QUERY, {
    fetchPolicy: 'cache-and-network',
    nextFetchPolicy: 'cache-first',
    errorPolicy: 'all',
  });
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
        allMediaCount: data.allMedia?.length ?? 0,
        hasMe: !!data.me,
        activitiesCount: data.me?.activities?.length ?? 0,
      });
    }
  }, [data, error]);

  const mediaItems = useMemo(() => {
    const allMedia = data?.allMedia ?? [];
    const activityMedia = (data?.me?.activities ?? [])
      .map(activity => activity.media)
      .filter((item): item is HomeMediaData['allMedia'][number] => Boolean(item));
    const filteredMedia = allMedia.filter(
      item =>
        item.__typename === 'Movie' ||
        item.__typename === 'TVShow' ||
        item.__typename === 'Book' ||
        item.__typename === 'Game' ||
        item.__typename === 'MusicAlbum'
    );
    const candidates = capMediaCandidates(filteredMedia);
    const scores = scoreMediaFromUser({
      candidates,
      activityMedia,
    });
    return candidates
      .map(item => {
        let genre: string[];
        if (item.genres) {
          genre = item.genres.map(g => g.name);
        } else if (item.subjects) {
          genre = item.subjects.map(s => s.name);
        } else if (item.genre) {
          genre = item.genre;
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
          score: scores.get(String(item.id)) ?? 0,
        };
      })
      .sort((a, b) => {
        if (b.score !== a.score) return b.score - a.score;
        if (b.rating !== a.rating) return b.rating - a.rating;
        return a.title.localeCompare(b.title);
      })
      .map(({ score, ...item }) => item);
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
    } catch (err) {
      logError('Error refreshing media:', err);
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
