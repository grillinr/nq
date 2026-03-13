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
import Animated, { useAnimatedStyle, useSharedValue, withTiming } from 'react-native-reanimated';
import { useApolloClient, useQuery } from '@apollo/client/react';
import { fontSize, spacing, zIndex, layout } from '../../src/components/ui/tokens';
import { useTheme } from '../../src/components/ui/theme-provider';
import MediaCoverCard from '../../src/components/MediaCoverCard';
import MediaCoverSkeleton from '../../src/components/MediaCoverSkeleton';
import MediaTypeFilter from '../../src/components/MediaTypeFilter';
import PageHeader, { useHeaderHeight } from '../../src/components/PageHeader';
import { useAuth } from '../../src/lib/AuthContext';
import { ME_ACTIVITIES_QUERY, RECURSIVE_SEARCH_STATUS_QUERY } from '../../src/lib/graphql';
import { useScrollHeader } from '../../src/hooks/useScrollHeader';
import { Media, MediaType } from '../../src/types';

const calculateItemWidth = (windowWidth: number) => {
  const horizontalPadding = spacing[4] * 2;
  const gap = spacing[3] * 2;
  return Math.floor((windowWidth - horizontalPadding - gap) / 3);
};

const createStyles = (colors: ReturnType<typeof useTheme>['colors'], filterTop: number) =>
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
      justifyContent: 'space-between',
      marginBottom: spacing[3],
    },
    placeholderText: {
      color: colors.mutedForeground,
      textAlign: 'center',
      marginTop: spacing[6],
      fontSize: fontSize.base,
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
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);

  const headerHeight = useHeaderHeight();
  const filterTop = headerHeight + spacing[2];

  const filterTranslateY = useSharedValue(0);
  React.useEffect(() => {
    filterTranslateY.value = withTiming(isHeaderVisible ? 0 : -120, { duration: 300 });
  }, [isHeaderVisible, filterTranslateY]);

  const filterAnimatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateY: filterTranslateY.value }],
  }));

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
          const genre = Array.isArray(media.genres)
            ? media.genres.map((g: any) => g.name)
            : Array.isArray(media.subjects)
              ? media.subjects.map((s: any) => s.name)
              : Array.isArray(media.genre)
                ? media.genre
                : [];
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
            type:
              media.__typename === 'TVShow'
                ? 'tv'
                : media.__typename === 'Book'
                  ? 'book'
                  : media.__typename === 'Game'
                    ? 'game'
                    : media.__typename === 'MusicAlbum'
                      ? 'music'
                      : (media.__typename?.toLowerCase() ?? 'movie'),
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
      console.error('Error refreshing history:', err);
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
        if (state === 'COMPLETED' && isActive) {
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

  const styles = React.useMemo(() => createStyles(colors, filterTop), [colors, filterTop]);

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

  const listHeader = <View style={{ height: headerHeight + 52 }} />;

  return (
    <View style={styles.container}>
      <PageHeader title="Recently Viewed" visible={isHeaderVisible} />
      <Animated.View style={[styles.stickyFilterContainer, filterAnimatedStyle]}>
        <MediaTypeFilter
          selectedTypes={selectedMediaTypes}
          onFilterChange={setSelectedMediaTypes}
        />
      </Animated.View>
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
        onScroll={handleHeaderScroll}
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
