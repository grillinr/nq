import React from "react";
import { router, useLocalSearchParams } from "expo-router";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Pressable,
  ActivityIndicator,
  FlatList,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import MediaCoverCard from "../../src/components/MediaCoverCard";
import ImageWithFallback from "../../src/components/figma/ImageWithFallback";
import Badge from "../../src/components/ui/badge";
import { useTheme } from "../../src/components/ui/ThemeProvider";
import { fontSize, radii, spacing } from "../../src/components/ui/tokens";
import { useMediaDetails } from "../../src/hooks/useMediaDetails";

const COVER_RATIO = 2 / 3;

export default function MediaDetailsPage() {
  const { id } = useLocalSearchParams<{ id?: string | string[] }>();
  const { colors } = useTheme();
  const mediaId = Array.isArray(id) ? id[0] : id;
  const { details, loading, error } = useMediaDetails(mediaId);
  const styles = React.useMemo(() => createStyles(colors), [colors]);
  React.useEffect(() => {
    if (error) {
      console.error("Failed to load media details:", error);
    }
  }, [error]);
  const backButton = (
    <Pressable
      onPress={() => router.back()}
      style={styles.backButton}
      accessibilityRole="button"
    >
      <Ionicons name="arrow-back" size={18} color={colors.foreground} />
      <Text style={styles.backText}>Back</Text>
    </Pressable>
  );

  if (loading && !details) {
    return (
      <View style={styles.loadingContainer}>
        {backButton}
        <View style={styles.loadingCenter}>
          <ActivityIndicator color={colors.primary} />
        </View>
      </View>
    );
  }

  if (!details) {
    return (
      <View style={styles.loadingContainer}>
        {backButton}
        <View style={styles.loadingCenter}>
          <Text style={styles.emptyText}>
            {error ? "We hit an error loading this title." : "We couldn&apos;t find that title."}
          </Text>
        </View>
      </View>
    );
  }

  const actors = details.actors.slice(0, 10);
  const related = details.related;

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      {backButton}

      <View style={styles.hero}>
        <View style={styles.coverWrap}>
          <ImageWithFallback
            src={details.image}
            alt={details.title}
            style={styles.cover}
          />
        </View>
        <View style={styles.heroText}>
          <Text style={styles.title}>{details.title}</Text>
          <View style={styles.metaRow}>
            <Text style={styles.metaText}>{details.year}</Text>
            {details.duration ? <Text style={styles.metaText}>· {details.duration}</Text> : null}
            <Text style={styles.metaText}>· {details.type.toUpperCase()}</Text>
          </View>
          <View style={styles.ratingRow}>
            <Ionicons name="star" size={14} color={colors["chart-4"]} />
            <Text style={styles.ratingText}>{details.rating.toFixed(1)}</Text>
          </View>
          {details.metaLabel && details.metaValue ? (
            <Text style={styles.metaDetail}>
              {details.metaLabel}: {details.metaValue}
            </Text>
          ) : null}
          {details.genre.length > 0 ? (
            <View style={styles.genreRow}>
              {details.genre.slice(0, 4).map((item) => (
                <Badge key={item} variant="secondary">
                  {item}
                </Badge>
              ))}
            </View>
          ) : null}
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Overview</Text>
        <Text style={styles.bodyText}>
          {details.description || "No description available yet."}
        </Text>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Actors & Creators</Text>
        {actors.length === 0 && details.creators.length === 0 ? (
          <Text style={styles.bodyText}>No cast or creator info yet.</Text>
        ) : (
          <View style={styles.chipRow}>
            {actors.map((actor) => (
              <View key={actor.id} style={styles.chip}>
                <Text style={styles.chipText}>{actor.name}</Text>
              </View>
            ))}
            {details.creators.slice(0, 6).map((creator) => (
              <View key={creator.id} style={styles.chip}>
                <Text style={styles.chipText}>{creator.name}</Text>
              </View>
            ))}
          </View>
        )}
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Related</Text>
        {related.length === 0 ? (
          <Text style={styles.bodyText}>No related titles yet.</Text>
        ) : (
          <FlatList
            data={related}
            horizontal
            showsHorizontalScrollIndicator={false}
            keyExtractor={(item) => String(item.id)}
            contentContainerStyle={styles.relatedList}
            renderItem={({ item }) => (
              <View style={styles.relatedItem}>
                <MediaCoverCard
                  title={item.title}
                  image={item.image}
                  aspectRatio={COVER_RATIO}
                  onPress={() => router.push({ pathname: "/media/[id]", params: { id: item.id } })}
                />
                <Text style={styles.relatedTitle} numberOfLines={1}>
                  {item.title}
                </Text>
                <Text style={styles.relatedMeta} numberOfLines={1}>
                  {item.year} · {item.type.toUpperCase()}
                </Text>
              </View>
            )}
          />
        )}
      </View>
    </ScrollView>
  );
}

const createStyles = (colors: ReturnType<typeof useTheme>["colors"]) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
      paddingBottom: spacing[8],
    },
    loadingContainer: {
      flex: 1,
      backgroundColor: colors.background,
      padding: spacing[4],
    },
    loadingCenter: {
      flex: 1,
      alignItems: "center",
      justifyContent: "center",
    },
    backButton: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[2],
      marginBottom: spacing[4],
    },
    backText: {
      color: colors.foreground,
      fontSize: fontSize.base,
    },
    hero: {
      flexDirection: "row",
      gap: spacing[4],
      marginBottom: spacing[6],
    },
    coverWrap: {
      width: 140,
      aspectRatio: COVER_RATIO,
      borderRadius: radii.lg,
      overflow: "hidden",
      backgroundColor: colors["input-background"],
    },
    cover: {
      width: "100%",
      height: "100%",
    },
    heroText: {
      flex: 1,
    },
    title: {
      fontSize: fontSize["2xl"],
      fontWeight: "600",
      color: colors.foreground,
      marginBottom: spacing[2],
    },
    metaRow: {
      flexDirection: "row",
      flexWrap: "wrap",
      gap: spacing[2],
      marginBottom: spacing[2],
    },
    metaText: {
      color: colors["muted-foreground"],
      fontSize: fontSize.sm,
    },
    ratingRow: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[1],
      marginBottom: spacing[2],
    },
    ratingText: {
      color: colors.foreground,
      fontSize: fontSize.base,
      fontWeight: "500",
    },
    metaDetail: {
      color: colors["muted-foreground"],
      fontSize: fontSize.sm,
      marginBottom: spacing[2],
    },
    genreRow: {
      flexDirection: "row",
      flexWrap: "wrap",
      gap: spacing[2],
    },
    section: {
      marginBottom: spacing[6],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: "600",
      color: colors.foreground,
      marginBottom: spacing[3],
    },
    bodyText: {
      fontSize: fontSize.base,
      color: colors["muted-foreground"],
      lineHeight: 22,
    },
    chipRow: {
      flexDirection: "row",
      flexWrap: "wrap",
      gap: spacing[2],
    },
    chip: {
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[1],
      borderRadius: radii.md,
      backgroundColor: colors.secondary,
    },
    chipText: {
      fontSize: fontSize.sm,
      color: colors["secondary-foreground"],
    },
    relatedList: {
      gap: spacing[3],
    },
    relatedItem: {
      width: 120,
    },
    relatedTitle: {
      marginTop: spacing[2],
      color: colors.foreground,
      fontSize: fontSize.sm,
      fontWeight: "500",
    },
    relatedMeta: {
      color: colors["muted-foreground"],
      fontSize: fontSize.xs,
      marginTop: 2,
    },
    emptyText: {
      color: colors["muted-foreground"],
      fontSize: fontSize.base,
    },
  });
