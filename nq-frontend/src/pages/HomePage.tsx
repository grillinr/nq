import React from "react";
import { View, Text, ScrollView, StyleSheet } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import MediaCard from "../components/MediaCard";
import Badge from "../components/ui/badge";
import { fontSize, spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { Media } from "../types";

interface HomePageProps {
  mediaList: Media[];
}

function HomePage({ mediaList }: HomePageProps) {
  const { colors, resolved } = useTheme();

  // Get top rated items
  const topRated = [...mediaList]
    .sort((a, b) => b.rating - a.rating)
    .slice(0, 6);

  // Get recent items (by year)
  const recentItems = [...mediaList]
    .sort((a, b) => b.year - a.year)
    .slice(0, 6);

  // Get personalized recommendations (mix of different types)
  const recommended = [...mediaList]
    .sort(() => Math.random() - 0.5)
    .slice(0, 6);

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
    },
    hero: {
      backgroundColor: resolved === "light" ? colors.secondary : colors.primary,
      borderRadius: 8,
      padding: spacing[6],
      marginBottom: spacing[8],
    },
    heroContent: {},
    heroHeader: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[3],
      marginBottom: spacing[3],
    },
    heroTitle: {
      fontSize: fontSize.xl,
      fontWeight: "600",
      color:
        resolved === "light"
          ? colors["secondary-foreground"]
          : colors["primary-foreground"],
    },
    heroSubtitle: {
      fontSize: fontSize.base,
      color:
        resolved === "light"
          ? colors["secondary-foreground"]
          : colors["primary-foreground"],
      marginBottom: spacing[6],
      maxWidth: 300,
    },
    heroBadges: {
      flexDirection: "row",
      gap: spacing[2],
    },
    heroBadge: {
      borderColor: colors.ring,
      ...(resolved === "dark" && { backgroundColor: colors.secondary }),
    },
    section: {
      marginBottom: spacing[8],
    },
    sectionHeader: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[2],
      marginBottom: spacing[4],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: "600",
      color: colors.foreground,
    },
    horizontalScroll: {
      paddingLeft: spacing[4],
    },
    cardList: {
      gap: spacing[4],
    },
  });

  const renderSection = (
    title: string,
    icon: keyof typeof Ionicons.glyphMap,
    data: Media[],
  ) => (
    <View style={styles.section}>
      <View style={styles.sectionHeader}>
        <Ionicons name={icon} size={20} color={colors.primary} />
        <Text style={styles.sectionTitle}>{title}</Text>
      </View>
      <View style={styles.cardList}>
        {data.slice(0, 3).map((item) => (
          <MediaCard key={item.id} {...item} />
        ))}
      </View>
    </View>
  );

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      {/* <View style={styles.hero}> */}
      {/*   <View style={styles.heroContent}> */}
      {/*     <View style={styles.heroHeader}> */}
      {/*       <Ionicons */}
      {/*         name="sparkles" */}
      {/*         size={32} */}
      {/*         color={ */}
      {/*           resolved === "light" */}
      {/*             ? colors["secondary-foreground"] */}
      {/*             : colors["primary-foreground"] */}
      {/*         } */}
      {/*       /> */}
      {/*       <Text style={styles.heroTitle}>Discover Your Next Favorite</Text> */}
      {/*     </View> */}
      {/*     <Text style={styles.heroSubtitle}> */}
      {/*       Explore personalized recommendations based on your viewing history */}
      {/*       and preferences */}
      {/*     </Text> */}
      {/*     <View style={styles.heroBadges}> */}
      {/*       <Badge */}
      {/*         variant={resolved === "light" ? "outline" : "secondary"} */}
      {/*         style={styles.heroBadge} */}
      {/*       > */}
      {/*         {mediaList.length} {mediaList.length === 1 ? "item" : "items"} in */}
      {/*         collection */}
      {/*       </Badge> */}
      {/*       <Badge */}
      {/*         variant={resolved === "light" ? "outline" : "secondary"} */}
      {/*         style={styles.heroBadge} */}
      {/*       > */}
      {/*         Updated today */}
      {/*       </Badge> */}
      {/*     </View> */}
      {/*   </View> */}
      {/* </View> */}

      {renderSection("Recommended For You", "sparkles-outline", recommended)}
      {renderSection("Top Rated", "trending-up-outline", topRated)}
      {renderSection("Recently Added", "time-outline", recentItems)}
    </ScrollView>
  );
}

export default HomePage;
