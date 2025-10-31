import React, { useState, useMemo } from 'react';
import { View, Text, StyleSheet, FlatList } from 'react-native';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import MediaCard from './MediaCard';
import FilterPanel from './FilterPanel';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, fontSize } from './ui/tokens';

interface Media {
  id: number;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: "movie" | "tv" | "book" | "music" | "game";
}

interface HistoryPageProps {
  mediaList: Media[];
}

const genresByType = {
  movie: ["Action", "Sci-Fi", "Drama", "Thriller", "Crime"],
  tv: ["Crime", "Drama", "Sci-Fi", "Horror", "Historical", "Action"],
  book: ["Fiction", "Classic", "Dystopian", "Sci-Fi", "Drama"],
  music: ["Rock", "Pop", "Jazz", "R&B", "Classic", "Progressive"],
  game: ["Action", "RPG", "Strategy", "Puzzle", "Adventure", "Sports"],
};

const typeIcons = {
  movie: 'film',
  tv: 'tv',
  book: 'book',
  music: 'musical-notes',
  game: 'game-controller',
};

function HistoryPage({ mediaList }: HistoryPageProps) {
  const [activeTab, setActiveTab] = useState("movie");
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
  const [sortBy, setSortBy] = useState("rating");
  const [minRating, setMinRating] = useState(0);

  const filteredAndSortedMedia = useMemo(() => {
    let filtered = mediaList.filter((item) => {
      if (item.type !== activeTab) return false;
      if (item.rating < minRating) return false;
      if (selectedGenres.length > 0 && !selectedGenres.some((g) => item.genre.includes(g))) {
        return false;
      }
      return true;
    });

    filtered.sort((a, b) => {
      if (sortBy === "rating") return b.rating - a.rating;
      if (sortBy === "year") return b.year - a.year;
      if (sortBy === "title") return a.title.localeCompare(b.title);
      return 0;
    });

    return filtered;
  }, [mediaList, activeTab, selectedGenres, sortBy, minRating]);

  const handleGenreToggle = (genre: string) => {
    setSelectedGenres((prev) =>
      prev.includes(genre) ? prev.filter((g) => g !== genre) : [...prev, genre]
    );
  };

  const handleTabChange = (value: string) => {
    setActiveTab(value);
    setSelectedGenres([]);
    setMinRating(0);
  };

  const renderMediaItem = ({ item }: { item: Media }) => (
    <MediaCard key={item.id} {...item} />
  );

  return (
    <Tabs value={activeTab} onValueChange={handleTabChange}>
      <TabsList>
        <TabsTrigger value="movie">
          <Ionicons name={typeIcons.movie as any} size={20} />
        </TabsTrigger>
        <TabsTrigger value="tv">
          <Ionicons name={typeIcons.tv as any} size={20} />
        </TabsTrigger>
        <TabsTrigger value="book">
          <Ionicons name={typeIcons.book as any} size={20} />
        </TabsTrigger>
        <TabsTrigger value="music">
          <Ionicons name={typeIcons.music as any} size={20} />
        </TabsTrigger>
        <TabsTrigger value="game">
          <Ionicons name={typeIcons.game as any} size={20} />
        </TabsTrigger>
      </TabsList>

      <View style={styles.filterPanel}>
        <FilterPanel
          genres={genresByType[activeTab as keyof typeof genresByType]}
          selectedGenres={selectedGenres}
          onGenreToggle={handleGenreToggle}
          sortBy={sortBy}
          onSortChange={setSortBy}
          minRating={minRating}
          onRatingChange={(value) => setMinRating(value[0])}
        />
      </View>

      <View style={styles.mediaList}>
        <TabsContent value={activeTab}>
          {filteredAndSortedMedia.length > 0 ? (
            <FlatList
              data={filteredAndSortedMedia}
              renderItem={renderMediaItem}
              keyExtractor={(item) => item.id.toString()}
              numColumns={2}
              contentContainerStyle={styles.flatListContent}
            />
          ) : (
            <View style={styles.noResults}>
              <Text style={styles.noResultsText}>No results found. Try adjusting your filters.</Text>
            </View>
          )}
        </TabsContent>
      </View>
    </Tabs>
  );
}

export default HistoryPage;

const styles = StyleSheet.create({
  filterPanel: {
    marginBottom: spacing[6],
  },
  mediaList: {
    flex: 1,
  },
  flatListContent: {
    padding: spacing[4],
  },
  noResults: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing[8],
  },
  noResultsText: {
    color: colors['muted-foreground'],
    fontSize: fontSize.base,
  },
});