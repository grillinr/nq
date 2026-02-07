import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { fontSize, spacing } from '../components/ui/tokens';
import { useTheme } from '../components/ui/ThemeProvider';
import MediaCard from '../components/MediaCard';
import { useAuth } from '../../lib/AuthContext';
import { useQuery } from '@apollo/client/react';
import { ME_ACTIVITIES_QUERY } from '../../lib/graphql';
import { Media } from '../types';

function HistoryPage() {
  const { colors } = useTheme();
  const { hasToken } = useAuth();
  const { data, loading } = useQuery(ME_ACTIVITIES_QUERY, {
    fetchPolicy: 'cache-and-network',
    nextFetchPolicy: 'cache-first',
    skip: !hasToken,
  });

  const mediaList: Media[] = (data?.me?.activities ?? [])
    .map((activity: any) => activity.media)
    .filter(Boolean)
    .map((media: any) => {
      const genre = Array.isArray(media.genres)
        ? media.genres.map((g: any) => g.name)
        : Array.isArray(media.genre)
          ? media.genre
          : [];
      return {
        id: media.id,
        title: media.title ?? 'Untitled',
        image: media.coverUrl || 'https://placehold.co/400x600?text=No+Image',
        rating: media.averageRating || 0,
        genre,
        year: media.releaseDate ? parseInt(media.releaseDate.substring(0, 4)) : new Date().getFullYear(),
        duration: media.runtime ? `${Math.floor(media.runtime / 60)}h ${media.runtime % 60}m` : undefined,
        description: media.description || '',
        type: media.__typename === 'TVShow' ? 'tv' : media.__typename?.toLowerCase() ?? 'movie',
      } as Media;
    });

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
    },
    section: {
      marginBottom: spacing[6],
    },
    list: {
      gap: spacing[4],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: '600',
      color: colors.foreground,
      marginBottom: spacing[3],
    },
    placeholderText: {
      color: colors.foreground,
      textAlign: 'center',
      marginTop: spacing[4],
    },
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Recently Viewed</Text>
        <View style={styles.list}>
          {loading ? (
            <Text style={styles.placeholderText}>Loading history…</Text>
          ) : mediaList.length === 0 ? (
            <Text style={styles.placeholderText}>No history yet. Add your first title.</Text>
          ) : (
            mediaList.map((item) => (
              <MediaCard key={String(item.id)} {...item} onPress={() => {}} />
            ))
          )}
        </View>
      </View>
    </ScrollView>
  );
}

export default HistoryPage;
