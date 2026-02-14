import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  useWindowDimensions,
  TouchableOpacity,
  RefreshControl,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { router, useLocalSearchParams } from 'expo-router';
import Animated, { useAnimatedStyle, SharedValue } from 'react-native-reanimated';
import { fontSize, spacing } from '../components/ui/tokens';
import { useTheme } from '../components/ui/ThemeProvider';
import MediaCoverCard from '../components/MediaCoverCard';
import MediaCoverSkeleton from '../components/MediaCoverSkeleton';
import MediaTypeFilter from '../components/MediaTypeFilter';
import PageHeader from '../components/PageHeader';
import { useAuth } from '../../lib/AuthContext';
import { useApolloClient, useQuery } from '@apollo/client/react';
import { ME_ACTIVITIES_QUERY, RECURSIVE_SEARCH_STATUS_QUERY } from '../../lib/graphql';
import { useScrollHeader } from '../hooks/useScrollHeader';
import { Media, MediaType } from '../types';

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
  const [showStatusBanner, setShowStatusBanner] = React.useState(false);
  const [statusMessage, setStatusMessage] = React.useState('Related titles are ready.');
  const [hasCompletedSearch, setHasCompletedSearch] = React.useState(false);
  const [selectedMediaTypes, setSelectedMediaTypes] = React.useState<MediaType[]>([]);
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);
  const headerTranslateY = React.useRef<SharedValue<number> | null>(null);

  const mediaList: Media[] = (data?.me?.activities ?? [])
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
          ? parseInt(media.releaseDate.substring(0, 4))
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
    });

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
    if (!data?.me?.activities?.length) return undefined;
    return data.me.activities[0]?.media?.id;
  }, [addedMediaId, data?.me?.activities]);

  React.useEffect(() => {
    if (!hasToken) return;
    if (!latestMediaId) return;
    if (hasCompletedSearch) return;

    let isActive = true;
    let interval: ReturnType<typeof setInterval> | undefined;
    const poll = async () => {
      try {
        const result = await apolloClient.query({
          query: RECURSIVE_SEARCH_STATUS_QUERY,
          variables: { mediaId: latestMediaId },
          fetchPolicy: 'no-cache',
        });
        const state = result.data?.recursiveSearchStatus?.state;
        if (state === 'COMPLETED' && isActive) {
          setStatusMessage('Related titles are ready.');
          setShowStatusBanner(true);
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
    if (!showStatusBanner) return;
    const timeout = setTimeout(() => setShowStatusBanner(false), 5000);
    return () => clearTimeout(timeout);
  }, [showStatusBanner]);

  React.useEffect(() => {
    setHasCompletedSearch(false);
  }, [latestMediaId]);

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    stickyFilterContainer: {
      position: 'absolute',
      top: 80,
      left: 0,
      right: 0,
      backgroundColor: 'transparent',
      paddingHorizontal: spacing[4],
      paddingVertical: spacing[3],
      zIndex: 999,
    },
    listContent: {
      padding: spacing[4],
    },
    row: {
      justifyContent: 'space-between',
      marginBottom: spacing[3],
    },
    placeholderText: {
      color: colors.mutedForeground,
      textAlign: 'center',
      marginTop: spacing[6],
      fontSize: 16,
    },
    statusBanner: {
      paddingHorizontal: spacing[4],
      paddingVertical: spacing[3],
      backgroundColor: colors.primary,
      marginHorizontal: 0,
      borderRadius: 10,
      marginTop: 110,
      marginBottom: spacing[3],
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'space-between',
    },
    statusBannerText: {
      color: colors.primaryForeground,
      fontSize: fontSize.sm,
      fontWeight: '600',
      flex: 1,
      marginRight: spacing[3],
    },
    statusBannerDismiss: {
      color: colors.primaryForeground,
    },
  });

  const skeletonData = React.useMemo(
    () => Array.from({ length: 12 }, (_, index) => ({ id: `skeleton-${index}` })),
    []
  );

  const renderItem = React.useCallback(
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
        />
      );
    },
    [itemWidth]
  );

  const listData = loading ? skeletonData : filteredMediaList;
  const listEmptyComponent = !loading ? (
    <Text style={styles.placeholderText}>{emptyStateMessage}</Text>
  ) : null;

  const filterAnimatedStyle = useAnimatedStyle(() => {
    if (!headerTranslateY.current) {
      return {};
    }
    return {
      transform: [{ translateY: headerTranslateY.current.value }],
    };
  }, []);

  const listHeader = (
    <View>
      {showStatusBanner ? (
        <View style={styles.statusBanner}>
          <Text style={styles.statusBannerText}>{statusMessage}</Text>
          <TouchableOpacity onPress={() => setShowStatusBanner(false)}>
            <Ionicons name="close" size={18} color={styles.statusBannerDismiss.color} />
          </TouchableOpacity>
        </View>
      ) : (
        <View style={{ height: 140 }} />
      )}
    </View>
  );

  return (
    <View style={styles.container}>
      <PageHeader
        title="Recently Viewed"
        visible={isHeaderVisible}
        onTranslateYChange={translateY => {
          headerTranslateY.current = translateY;
        }}
      />
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

const calculateItemWidth = (windowWidth: number) => {
  const horizontalPadding = spacing[4] * 2;
  const gap = spacing[3] * 2;
  return Math.floor((windowWidth - horizontalPadding - gap) / 3);
};
