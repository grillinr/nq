import React from "react";
import { View, Text, StyleSheet, FlatList, useWindowDimensions } from "react-native";
import { router } from "expo-router";
import { fontSize, spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import MediaCoverCard from "../components/MediaCoverCard";
import MediaCoverSkeleton from "../components/MediaCoverSkeleton";
import { useAuth } from "../../lib/AuthContext";
import { useQuery } from "@apollo/client/react";
import { ME_ACTIVITIES_QUERY } from "../../lib/graphql";
import { Media } from "../types";

function HistoryPage() {
  const { colors } = useTheme();
  const { hasToken } = useAuth();
  const { width } = useWindowDimensions();
  const itemWidth = React.useMemo(() => calculateItemWidth(width), [width]);
  const { data, loading } = useQuery(ME_ACTIVITIES_QUERY, {
    fetchPolicy: "cache-and-network",
    nextFetchPolicy: "cache-first",
    skip: !hasToken,
  });

  const mediaList: Media[] = (data?.me?.activities ?? [])
    .map((activity: any) => activity.media)
    .filter(
      (media: any) =>
        media?.__typename === "Movie" ||
        media?.__typename === "TVShow" ||
        media?.__typename === "Book" ||
        media?.__typename === "Game" ||
        media?.__typename === "MusicAlbum"
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
        title: media.title ?? "Untitled",
        image:
          media.coverUrl ||
          `https://placehold.co/400x600?text=${encodeURIComponent(media.title ?? "Untitled")}`,
        rating: media.averageRating || 0,
        genre,
        year: media.releaseDate
          ? parseInt(media.releaseDate.substring(0, 4))
          : new Date().getFullYear(),
        duration: media.runtime ? `${Math.floor(media.runtime / 60)}h ${media.runtime % 60}m` : undefined,
        description: media.description || "",
        type:
          media.__typename === "TVShow"
            ? "tv"
            : media.__typename === "Book"
              ? "book"
              : media.__typename === "Game"
                ? "game"
                : media.__typename === "MusicAlbum"
                  ? "music"
                  : media.__typename?.toLowerCase() ?? "movie",
      } as Media;
    });

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    listContent: {
      padding: spacing[4],
    },
    row: {
      justifyContent: "space-between",
      marginBottom: spacing[3],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: "600",
      color: colors.foreground,
      marginBottom: spacing[3],
    },
    placeholderText: {
      color: colors.foreground,
      textAlign: "center",
      marginTop: spacing[4],
    },
  });

  const skeletonData = React.useMemo(
    () => Array.from({ length: 12 }, (_, index) => ({ id: `skeleton-${index}` })),
    []
  );

  const renderItem = React.useCallback(
    ({ item }: { item: Media | { id: string } }) => {
      if (String(item.id).startsWith("skeleton")) {
        return <MediaCoverSkeleton style={{ width: itemWidth }} />;
      }
      const mediaItem = item as Media;
      return (
        <MediaCoverCard
          title={mediaItem.title}
          image={mediaItem.image}
          onPress={() =>
            router.push({ pathname: "/media/[id]", params: { id: String(mediaItem.id) } })
          }
          style={{ width: itemWidth }}
        />
      );
    },
    [itemWidth]
  );

  const listData = loading ? skeletonData : mediaList;
  const listEmptyComponent = !loading ? (
    <Text style={styles.placeholderText}>No history yet. Add your first title.</Text>
  ) : null;

  return (
    <View style={styles.container}>
      <FlatList
        data={listData}
        renderItem={renderItem}
        keyExtractor={(item) => String(item.id)}
        contentContainerStyle={styles.listContent}
        numColumns={3}
        columnWrapperStyle={styles.row}
        ListHeaderComponent={<Text style={styles.sectionTitle}>Recently Viewed</Text>}
        ListEmptyComponent={listEmptyComponent}
        showsVerticalScrollIndicator={false}
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
