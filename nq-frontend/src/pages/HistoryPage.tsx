import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import MediaCard from '../components/MediaCard';
import { fontSize, spacing } from '../components/ui/tokens';
import { useTheme } from '../components/ui/ThemeProvider';

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

interface HistoryPageProps {
  mediaList: Media[];
}

function HistoryPage({ mediaList }: HistoryPageProps) {
  const { colors } = useTheme();

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
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Recently Viewed</Text>
        <View style={styles.list}>
          {mediaList.map((item) => (
            <MediaCard key={item.id} {...item} />
          ))}
        </View>
      </View>
    </ScrollView>
  );
}

export default HistoryPage;
