import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  useWindowDimensions,
  RefreshControl,
} from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import * as Haptics from 'expo-haptics';
import { useApolloClient, useQuery } from '@apollo/client/react';
import { fontSize, spacing, zIndex, layout, sizes } from '../../src/components/ui/tokens';
import { useTheme } from '../../src/components/ui/theme-provider';
import MediaCoverCard from '../../src/components/MediaCoverCard';
import MediaCoverSkeleton from '../../src/components/MediaCoverSkeleton';
import MediaTypeFilter from '../../src/components/MediaTypeFilter';
import PageHeader, { useHeaderHeight } from '../../src/components/PageHeader';
import { useAuth } from '../../src/lib/AuthContext';
import { ME_ACTIVITIES_QUERY, RECURSIVE_SEARCH_STATUS_QUERY } from '../../src/lib/graphql';
import { Media, MediaType } from '../../src/types';
import { logError } from '../../src/lib/logger';

const FILTER_BAR_HEIGHT = sizes[13]; // 52 — matches the filter pill's rendered height

const calculateItemWidth = (windowWidth: number) => {
  const horizontalPadding = spacing[4] * 2;
  const gap = spacing[3] * 2;
  return Math.floor((windowWidth - horizontalPadding - gap) / 3);
};

const createStyles = (
  colors: ReturnType<typeof useTheme>['colors'],
  filterTop: number,
  headerHeight: number
) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    stickyFilterContainer: {
      position: 'absolute',
      top: filterTop,
      left: spacing[4],
      right: spacing[4],
      zIndex: zIndex.modal,
      alignItems: 'flex-start',
    },
    listContent: {
      padding: spacing[4],
      paddingBottom: layout.tabBarHeight + spacing[4],
    },
    row: {
      justifyContent: 'flex-start',
      gap: spacing[3],
      marginBottom: spacing[3],
    },
    placeholderText: {
      color: colors.mutedForeground,
      textAlign: 'center',
      marginTop: spacing[6],
      fontSize: fontSize.base,
    },
    listHeaderSpacer: {
      height: headerHeight + FILTER_BAR_HEIGHT,
    },
  });

function HistoryPage() {
  const { colors } = useTheme();
  const { hasToken } = useAuth();
  const { width } = useWindowDimensions();
  const { addedMediaId } = useLocalSearchParams();
  const apolloClient = useApolloClient();
  const itemWidth = React.useMemo(() => calculateItemWidth(width), [width]);
  const { data, loading, refetch } = useQuery(ME_ACTIVITIES_QUERY, {
    fetchPolicy: 'cache-first',
    skip: !hasToken,
  });
  const [refreshing, setRefreshing] = React.useState(false);
  const [hasCompletedSearch, setHasCompletedSearch] = React.useState(false);
  const [enrichingMediaId, setEnrichingMediaId] = React.useState<string | undefined>(undefined);
  const [selectedMediaTypes, setSelectedMediaTypes] = React.useState<MediaType[]>([]);
  const headerHeight = useHeaderHeight();
  const filterTop = headerHeight + spacing[2];

  const mediaList: Media[] = React.useMemo(
    () =>
      ((data as any)?.me?.activities ?? [])
        .map((activity: any) => activity.media)
        .filter(
          (media: any) =>
            media?.__typename === 'Movie' ||
            media?.__typename === 'TVShow' ||
            media?.__typename === 'Book' ||
            media?.__typename === 'Game' ||
            media?.__typename === 'MusicAlbum'
        )
        .map((media: any) => {
          let genre: string[];
          if (Array.isArray(media.genres)) {
            genre = media.genres.map((g: any) => g.name);
          } else if (Array.isArray(media.subjects)) {
            genre = media.subjects.map((s: any) => s.name);
          } else if (Array.isArray(media.genre)) {
            genre = media.genre;
          } else {
            genre = [];
          }

          let type: string;
          if (media.__typename === 'TVShow') {
            type = 'tv';
          } else if (media.__typename === 'Book') {
            type = 'book';
          } else if (media.__typename === 'Game') {
            type = 'game';
          } else if (media.__typename === 'MusicAlbum') {
            type = 'music';
          } else {
            type = media.__typename?.toLowerCase() ?? 'movie';
          }

          return {
            id: media.id,
            title: media.title ?? 'Untitled',
            image:
              media.coverUrl ||
              `https://placehold.co/400x600?text=${encodeURIComponent(media.title ?? 'Untitled')}`,
            rating: media.averageRating || 0,
            genre,
            year: media.releaseDate
              ? parseInt(media.releaseDate.substring(0, 4), 10)
              : new Date().getFullYear(),
            duration: media.runtime
              ? `${Math.floor(media.runtime / 60)}h ${media.runtime % 60}m`
              : undefined,
            description: media.description || '',
            type,
          } as Media;
        }),
    [data]
  );

  const filteredMediaList = React.useMemo(() => {
    if (selectedMediaTypes.length === 0) return mediaList;
    return mediaList.filter(item => selectedMediaTypes.includes(item.type));
  }, [mediaList, selectedMediaTypes]);

  const emptyStateMessage = React.useMemo(() => {
    if (selectedMediaTypes.length > 0) {
      return 'No results meet filter criteria.';
    }
    return 'No history yet. Add your first title.';
  }, [selectedMediaTypes.length]);

  const onRefresh = React.useCallback(async () => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    setRefreshing(true);
    try {
      await refetch();
    } catch (err) {
      logError('Error refreshing history:', err);
    } finally {
      setRefreshing(false);
    }
  }, [refetch]);

  const latestMediaId = React.useMemo(() => {
    if (typeof addedMediaId === 'string' && addedMediaId) {
      return addedMediaId;
    }
    if (!(data as any)?.me?.activities?.length) return undefined;
    return (data as any).me.activities[0]?.media?.id;
  }, [addedMediaId, data]);

  React.useEffect(() => {
    if (!hasToken) return undefined;
    if (!latestMediaId) return undefined;
    if (hasCompletedSearch) return undefined;

    setEnrichingMediaId(latestMediaId as string);

    let isActive = true;
    let interval: ReturnType<typeof setInterval> | undefined;
    const poll = async () => {
      try {
        const result = await apolloClient.query({
          query: RECURSIVE_SEARCH_STATUS_QUERY,
          variables: { mediaId: latestMediaId },
          fetchPolicy: 'no-cache',
        });
        const state = (result.data as any)?.recursiveSearchStatus?.state;
        if ((state === 'COMPLETED' || state === 'IDLE') && isActive) {
          setEnrichingMediaId(undefined);
          setHasCompletedSearch(true);
          if (interval) {
            clearInterval(interval);
          }
        }
      } catch {
        // ignore polling errors
      }
    };

    interval = setInterval(poll, 4000);
    poll();

    return () => {
      isActive = false;
      if (interval) {
        clearInterval(interval);
      }
    };
  }, [apolloClient, hasCompletedSearch, hasToken, latestMediaId]);

  React.useEffect(() => {
    setHasCompletedSearch(false);
  }, [latestMediaId]);

  const styles = React.useMemo(
    () => createStyles(colors, filterTop, headerHeight),
    [colors, filterTop, headerHeight]
  );

  const skeletonData = React.useMemo(
    () => Array.from({ length: 12 }, (_, index) => ({ id: `skeleton-${index}` })),
    []
  );

  const renderItem = React.useCallback(
    // eslint-disable-next-line react/no-unused-prop-types
    ({ item }: { item: Media | { id: string } }) => {
      if (String(item.id).startsWith('skeleton')) {
        return <MediaCoverSkeleton style={{ width: itemWidth }} />;
      }
      const mediaItem = item as Media;
      return (
        <MediaCoverCard
          title={mediaItem.title}
          image={mediaItem.image}
          onPress={() =>
            router.push({ pathname: '/media/[id]', params: { id: String(mediaItem.id) } })
          }
          style={{ width: itemWidth }}
          isEnriching={String(mediaItem.id) === enrichingMediaId}
        />
      );
    },
    [itemWidth, enrichingMediaId]
  );

  const listData = loading ? skeletonData : filteredMediaList;
  const listEmptyComponent = !loading ? (
    <Text style={styles.placeholderText}>{emptyStateMessage}</Text>
  ) : null;

  const listHeader = <View style={styles.listHeaderSpacer} />;

  return (
    <View style={styles.container}>
      <PageHeader title="Recently Viewed" />
      <View style={styles.stickyFilterContainer}>
        <MediaTypeFilter
          selectedTypes={selectedMediaTypes}
          onFilterChange={setSelectedMediaTypes}
        />
      </View>
      <FlatList
        data={listData}
        renderItem={renderItem}
        keyExtractor={item => String(item.id)}
        contentContainerStyle={styles.listContent}
        numColumns={3}
        columnWrapperStyle={styles.row}
        ListHeaderComponent={listHeader}
        ListEmptyComponent={listEmptyComponent}
        showsVerticalScrollIndicator={false}
        scrollEventThrottle={16}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.primary}
          />
        }
      />
    </View>
  );
}

export default HistoryPage;
