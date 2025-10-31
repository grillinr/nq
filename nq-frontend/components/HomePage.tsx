import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import MediaCard from './MediaCard';
import Badge from './ui/badge';
import { colors, fontSize, spacing } from './ui/tokens';

interface Media {
  id: number;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: 'movie' | 'tv' | 'book' | 'music' | 'game';
}

interface HomePageProps {
  mediaList: Media[];
}

function HomePage({ mediaList }: HomePageProps) {
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

  const renderSection = (title: string, icon: string, data: Media[]) => (
    <View style={styles.section}>
      <View style={styles.sectionHeader}>
        <Ionicons name={icon as any} size={20} color={colors.primary} />
        <Text style={styles.sectionTitle}>{title}</Text>
      </View>
      <ScrollView
        horizontal
        showsHorizontalScrollIndicator={false}
        contentContainerStyle={styles.horizontalScroll}
      >
        {data.map((item) => (
          <View key={item.id} style={styles.cardWrapper}>
            <MediaCard {...item} />
          </View>
        ))}
      </ScrollView>
    </View>
  );

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.hero}>
        <View style={styles.heroContent}>
          <View style={styles.heroHeader}>
            <Ionicons name="sparkles" size={32} color="white" />
            <Text style={styles.heroTitle}>Discover Your Next Favorite</Text>
          </View>
          <Text style={styles.heroSubtitle}>
            Explore personalized recommendations based on your viewing history and preferences
          </Text>
          <View style={styles.heroBadges}>
            <Badge variant="secondary" style={styles.heroBadge}>
              {mediaList.length} items in collection
            </Badge>
            <Badge variant="secondary" style={styles.heroBadge}>
              Updated today
            </Badge>
          </View>
        </View>
      </View>

      {renderSection('Recommended For You', 'sparkles-outline', recommended)}
      {renderSection('Top Rated', 'trending-up-outline', topRated)}
      {renderSection('Recently Added', 'time-outline', recentItems)}
    </ScrollView>
  );
}

export default HomePage;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing[4],
  },
  hero: {
    backgroundColor: colors.primary,
    borderRadius: 8,
    padding: spacing[6],
    marginBottom: spacing[8],
  },
  heroContent: {},
  heroHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
    marginBottom: spacing[3],
  },
  heroTitle: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: 'white',
  },
  heroSubtitle: {
    fontSize: fontSize.base,
    color: '#e9d5ff', // purple-100 equivalent
    marginBottom: spacing[6],
    maxWidth: 300,
  },
  heroBadges: {
    flexDirection: 'row',
    gap: spacing[2],
  },
  heroBadge: {
    backgroundColor: 'rgba(255,255,255,0.2)',
    borderColor: 'rgba(255,255,255,0.3)',
  },
  section: {
    marginBottom: spacing[8],
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[2],
    marginBottom: spacing[4],
  },
  sectionTitle: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    color: colors.foreground,
  },
  horizontalScroll: {
    paddingLeft: spacing[4],
  },
  cardWrapper: {
    width: 180,
    marginRight: spacing[4],
  },
});