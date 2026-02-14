import React from 'react';
import {
  StyleSheet,
  FlatList,
  RefreshControl,
  useWindowDimensions,
  View,
  Pressable,
  NativeSyntheticEvent,
  NativeScrollEvent,
  Text,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { router } from 'expo-router';
import Animated, { useAnimatedStyle, SharedValue } from 'react-native-reanimated';
import MediaCoverCard from '../components/MediaCoverCard';
import MediaCoverSkeleton from '../components/MediaCoverSkeleton';
import MediaTypeFilter from '../components/MediaTypeFilter';
import PageHeader from '../components/PageHeader';
import { spacing } from '../components/ui/tokens';
import { useTheme } from '../components/ui/ThemeProvider';
import { useHomeMedia } from '../hooks/useHomeMedia';
import { useScrollHeader } from '../hooks/useScrollHeader';
import { Media, MediaType } from '../types';

function HomePage() {
  const { colors } = useTheme();
  const { media, loading, loadMore, hasMore, refresh } = useHomeMedia(18);
  const [refreshing, setRefreshing] = React.useState(false);
  const [showScrollTop, setShowScrollTop] = React.useState(false);
  const [selectedMediaTypes, setSelectedMediaTypes] = React.useState<MediaType[]>([]);
  const { width } = useWindowDimensions();
  const listRef = React.useRef<FlatList<Media | { id: string }>>(null);
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);
  const headerTranslateY = React.useRef<SharedValue<number> | null>(null);

  const onRefresh = React.useCallback(async () => {
    setRefreshing(true);
    await refresh();
    setRefreshing(false);
  }, [refresh]);

  const styles = React.useMemo(() => createStyles(colors), [colors]);
  const itemWidth = React.useMemo(() => calculateItemWidth(width), [width]);

  const onEndReached = React.useCallback(() => {
    if (hasMore) {
      loadMore();
    }
  }, [hasMore, loadMore]);

  const handleScroll = React.useCallback(
    (event: NativeSyntheticEvent<NativeScrollEvent>) => {
      const offsetY = event.nativeEvent.contentOffset.y;
      const shouldShow = offsetY > 600;
      setShowScrollTop(prev => (prev === shouldShow ? prev : shouldShow));
      handleHeaderScroll(event);
    },
    [handleHeaderScroll]
  );

  const scrollToTop = React.useCallback(() => {
    listRef.current?.scrollToOffset({ offset: 0, animated: true });
  }, []);

  const skeletonData = React.useMemo(
    () => Array.from({ length: 12 }, (_, index) => ({ id: `skeleton-${index}` })),
    []
  );
  const paginationSkeletonData = React.useMemo(
    () => Array.from({ length: 6 }, (_, index) => ({ id: `skeleton-page-${index}` })),
    []
  );

  const renderSkeletonItem = React.useCallback(
    () => <MediaCoverSkeleton style={{ width: itemWidth }} />,
    [itemWidth]
  );

  const listData = React.useMemo(() => {
    if (!loading || media.length > 0) return media;
    return [
      ...media,
      ...paginationSkeletonData.map(item => ({ ...item, id: `${item.id}-${media.length}` })),
    ];
  }, [loading, media, paginationSkeletonData]);

  const filteredMedia = React.useMemo(() => {
    if (selectedMediaTypes.length === 0) return listData;
    return listData.filter(item => {
      if (String(item.id).startsWith('skeleton')) return true; // Keep skeletons
      return selectedMediaTypes.includes((item as Media).type);
    });
  }, [listData, selectedMediaTypes]);

  const listRenderItem = React.useCallback(
    ({ item }: { item: Media | { id: string } }) => {
      if (String(item.id).startsWith('skeleton')) {
        return <MediaCoverSkeleton style={{ width: itemWidth }} />;
      }
      return (
        <MediaCoverCard
          title={item.title}
          image={item.image}
          onPress={() => router.push({ pathname: '/media/[id]', params: { id: String(item.id) } })}
          style={[styles.coverCard, { width: itemWidth }]}
        />
      );
    },
    [itemWidth, styles.coverCard]
  );

  const listKeyExtractor = React.useCallback((item: { id: string }) => String(item.id), []);

  const renderSeparator = React.useCallback(
    () => <View style={styles.separator} />,
    [styles.separator]
  );

  const listHeader = React.useMemo(() => <View style={{ paddingTop: 140 }} />, []);

  const emptyStateMessage = React.useMemo(() => {
    if (selectedMediaTypes.length > 0) {
      return 'No results meet filter criteria.';
    }
    return 'No recommendations yet. Add your first title.';
  }, [selectedMediaTypes.length]);

  const filterAnimatedStyle = useAnimatedStyle(() => {
    if (!headerTranslateY.current) {
      return {};
    }
    return {
      transform: [{ translateY: headerTranslateY.current.value }],
    };
  }, []);

  const listNode =
    loading && media.length === 0 ? (
      <FlatList
        ref={listRef}
        style={styles.list}
        contentContainerStyle={styles.listContent}
        data={skeletonData}
        renderItem={renderSkeletonItem}
        keyExtractor={listKeyExtractor}
        ItemSeparatorComponent={renderSeparator}
        ListHeaderComponent={listHeader}
        numColumns={3}
        columnWrapperStyle={styles.row}
        initialNumToRender={6}
        maxToRenderPerBatch={6}
        updateCellsBatchingPeriod={50}
        windowSize={7}
        removeClippedSubviews
        onScroll={handleScroll}
        scrollEventThrottle={16}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.primary}
          />
        }
      />
    ) : (
      <FlatList
        ref={listRef}
        style={styles.list}
        contentContainerStyle={styles.listContent}
        data={filteredMedia}
        renderItem={listRenderItem}
        keyExtractor={listKeyExtractor}
        ItemSeparatorComponent={renderSeparator}
        ListHeaderComponent={listHeader}
        ListEmptyComponent={<Text style={styles.emptyState}>{emptyStateMessage}</Text>}
        stickyHeaderIndices={[0]}
        onEndReached={onEndReached}
        onEndReachedThreshold={0.5}
        numColumns={3}
        columnWrapperStyle={styles.row}
        initialNumToRender={6}
        maxToRenderPerBatch={9}
        updateCellsBatchingPeriod={50}
        windowSize={9}
        removeClippedSubviews
        onScroll={handleScroll}
        scrollEventThrottle={16}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            tintColor={colors.primary}
          />
        }
      />
    );

  return (
    <View style={styles.container}>
      <PageHeader
        title="Recommended"
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
      {listNode}
      {showScrollTop && (
        <Pressable style={styles.fab} onPress={scrollToTop} accessibilityRole="button">
          <Ionicons name="arrow-up" size={18} color={colors.primaryForeground} />
        </Pressable>
      )}
    </View>
  );
}

export default HomePage;

const calculateItemWidth = (windowWidth: number) => {
  const horizontalPadding = spacing[4] * 2;
  const gap = spacing[3] * 2;
  return Math.floor((windowWidth - horizontalPadding - gap) / 3);
};

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
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
    list: {
      flex: 1,
    },
    listContent: {
      padding: spacing[4],
    },
    separator: {
      height: spacing[1],
    },
    row: {
      justifyContent: 'space-between',
      marginBottom: spacing[3],
    },
    emptyState: {
      marginTop: spacing[6],
      textAlign: 'center',
      color: colors.mutedForeground,
      fontSize: 16,
    },
    coverCard: {
      width: '100%',
    },
    fab: {
      position: 'absolute',
      right: spacing[4],
      bottom: spacing[6],
      backgroundColor: colors.primary,
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[3],
      borderRadius: 999,
      shadowColor: colors.border,
      shadowOffset: { width: 0, height: 6 },
      shadowOpacity: 0.2,
      shadowRadius: 8,
      elevation: 4,
    },
  });
