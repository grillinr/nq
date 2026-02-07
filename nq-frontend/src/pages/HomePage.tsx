import React from "react";
import { StyleSheet, FlatList, RefreshControl, useWindowDimensions, View, Pressable, Text, NativeSyntheticEvent, NativeScrollEvent } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import MediaCard from "../components/MediaCard";
import MediaCardSkeleton from "../components/MediaCardSkeleton";
import { fontSize, spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { useMovies } from "../hooks/useMovies";
import { Media } from "../types";

function HomePage() {
  const { colors } = useTheme();
  const { movies, loading, loadMore, hasMore, refresh } = useMovies(20);
  const [refreshing, setRefreshing] = React.useState(false);
  const [showScrollTop, setShowScrollTop] = React.useState(false);
  const { width } = useWindowDimensions();
  const listRef = React.useRef<FlatList<Media | { id: string }>>(null);

  const onRefresh = React.useCallback(async () => {
    setRefreshing(true);
    await refresh();
    setRefreshing(false);
  }, [refresh]);

  const styles = React.useMemo(() => createStyles(colors), [colors]);
  const cardHeight = React.useMemo(
    () => calculateCardHeight(width),
    [width]
  );

  const onEndReached = React.useCallback(() => {
    if (hasMore) {
      loadMore();
    }
  }, [hasMore, loadMore]);

  const handleScroll = React.useCallback((event: NativeSyntheticEvent<NativeScrollEvent>) => {
    const offsetY = event.nativeEvent.contentOffset.y;
    const shouldShow = offsetY > 600;
    setShowScrollTop((prev) => (prev === shouldShow ? prev : shouldShow));
  }, []);

  const scrollToTop = React.useCallback(() => {
    listRef.current?.scrollToOffset({ offset: 0, animated: true });
  }, []);

  const skeletonData = React.useMemo(
    () => Array.from({ length: 8 }, (_, index) => ({ id: `skeleton-${index}` })),
    []
  );
  const paginationSkeletonData = React.useMemo(
    () => Array.from({ length: 3 }, (_, index) => ({ id: `skeleton-page-${index}` })),
    []
  );

  const renderSkeletonItem = React.useCallback(() => <MediaCardSkeleton />, []);

  const listData = React.useMemo(() => {
    if (!loading || movies.length === 0) return movies;
    return [...movies, ...paginationSkeletonData.map((item) => ({ ...item, id: `${item.id}-${movies.length}` }))];
  }, [loading, movies, paginationSkeletonData]);


  const listRenderItem = React.useCallback(
    ({ item }: { item: Media | { id: string } }) => {
      if (String(item.id).startsWith("skeleton")) {
        return <MediaCardSkeleton />;
      }
      return <MediaCard {...item} />;
    },
    []
  );

  const listKeyExtractor = React.useCallback((item: { id: string }) => String(item.id), []);

  const getItemLayout = React.useCallback(
    (_: any, index: number) => ({
      length: cardHeight + spacing[4],
      offset: (cardHeight + spacing[4]) * index,
      index,
    }),
    [cardHeight]
  );

  const renderSeparator = React.useCallback(
    () => <View style={styles.separator} />,
    [styles.separator]
  );

  const listNode = loading && movies.length === 0 ? (
    <FlatList
      ref={listRef}
      style={styles.list}
      contentContainerStyle={styles.listContent}
      data={skeletonData}
      renderItem={renderSkeletonItem}
      keyExtractor={listKeyExtractor}
      getItemLayout={getItemLayout}
      ItemSeparatorComponent={renderSeparator}
      initialNumToRender={6}
      maxToRenderPerBatch={6}
      updateCellsBatchingPeriod={50}
      windowSize={7}
      removeClippedSubviews
      onScroll={handleScroll}
      scrollEventThrottle={16}
      showsVerticalScrollIndicator={false}
    />
  ) : (
    <FlatList
      ref={listRef}
      style={styles.list}
      contentContainerStyle={styles.listContent}
      data={listData}
      renderItem={listRenderItem}
      keyExtractor={listKeyExtractor}
      ItemSeparatorComponent={renderSeparator}
      onEndReached={onEndReached}
      onEndReachedThreshold={0.5}
      getItemLayout={getItemLayout}
      initialNumToRender={6}
      maxToRenderPerBatch={8}
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
      {listNode}
      {showScrollTop && (
        <Pressable style={styles.fab} onPress={scrollToTop} accessibilityRole="button">
          <Ionicons name="arrow-up" size={18} color={colors["primary-foreground"]} />
        </Pressable>
      )}
    </View>
  );
}

export default HomePage;

const CARD_PADDING = spacing[4];
const CONTENT_PADDING = spacing[4];
const TITLE_HEIGHT = fontSize.base;
const TITLE_MARGIN_BOTTOM = spacing[2];
const META_HEIGHT = 12;
const META_MARGIN_BOTTOM = spacing[3];
const GENRES_HEIGHT = fontSize.sm + 4;
const GENRES_MARGIN_BOTTOM = spacing[3];
const DESCRIPTION_HEIGHT = 36;

const calculateCardHeight = (windowWidth: number) => {
  const cardWidth = windowWidth - spacing[4] * 2;
  const innerWidth = cardWidth - CARD_PADDING * 2;
  const imageHeight = innerWidth * 1.5;
  const contentHeight =
    CONTENT_PADDING * 2 +
    TITLE_HEIGHT +
    TITLE_MARGIN_BOTTOM +
    META_HEIGHT +
    META_MARGIN_BOTTOM +
    GENRES_HEIGHT +
    GENRES_MARGIN_BOTTOM +
    DESCRIPTION_HEIGHT;
  return Math.round(imageHeight + contentHeight + CARD_PADDING * 2);
};

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    list: {
      flex: 1,
    },
    listContent: {
      padding: spacing[4],
    },
    separator: {
      height: spacing[4],
    },
    fab: {
      position: "absolute",
      right: spacing[4],
      bottom: spacing[6],
      backgroundColor: colors.primary,
      paddingHorizontal: spacing[4],
      paddingVertical: spacing[3],
      borderRadius: 999,
      shadowColor: colors.border,
      shadowOffset: { width: 0, height: 6 },
      shadowOpacity: 0.2,
      shadowRadius: 8,
      elevation: 4,
    },
  });
