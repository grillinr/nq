import React from "react";
import { View, StyleSheet, FlatList, ActivityIndicator, RefreshControl } from "react-native";
import MediaCard from "../components/MediaCard";
import { spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { useMovies } from "../hooks/useMovies";

function HomePage() {
  const { colors } = useTheme();
  const { movies, loading, loadMore, hasMore, refresh } = useMovies(20);
  const [refreshing, setRefreshing] = React.useState(false);

  const onRefresh = async () => {
    setRefreshing(true);
    await refresh();
    setRefreshing(false);
  };

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    listContent: {
      padding: spacing[4],
      gap: spacing[4],
    },
    footer: {
      padding: spacing[4],
      alignItems: "center",
      justifyContent: "center",
    },
  });

  const renderFooter = () => {
    if (!loading) return null;
    return (
      <View style={styles.footer}>
        <ActivityIndicator size="small" color={colors.primary} />
      </View>
    );
  };

  return (
    <FlatList
      style={styles.container}
      contentContainerStyle={styles.listContent}
      data={movies}
      renderItem={({ item }) => (
        <MediaCard
          key={item.id}
          // spread item properties to match MediaCardProps
          // MediaCard expects: title, image, rating, genre, year, duration, description
          // item (Media) has: title, image, rating, genre, year, duration, description, type, id
          // So we can just spread item.
          {...item}
          // We can override onPress if needed
          onPress={() => {}}
        />
      )}
      keyExtractor={(item) => String(item.id)}
      onEndReached={() => {
        if (hasMore) {
          loadMore();
        }
      }}
      onEndReachedThreshold={0.5}
      ListFooterComponent={renderFooter}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={onRefresh}
          tintColor={colors.primary}
        />
      }
    />
  );
}

export default HomePage;
