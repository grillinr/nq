import React from 'react';
import { View, StyleSheet, Pressable } from 'react-native';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import Slider from './ui/slider';
import Label from './ui/label';
import Badge from './ui/badge';
import { spacing, radii } from './ui/tokens';
import { useTheme } from './ui/ThemeProvider';

interface FilterPanelProps {
  genres: string[];
  selectedGenres: string[];
  onGenreToggle: (genre: string) => void;
  sortBy: string;
  onSortChange: (value: string) => void;
  minRating: number;
  onRatingChange: (value: number[]) => void;
}

function FilterPanel({
  genres,
  selectedGenres,
  onGenreToggle,
  sortBy,
  onSortChange,
  minRating,
  onRatingChange,
}: FilterPanelProps) {
  const { colors } = useTheme();

  const styles = StyleSheet.create({
    container: {
      backgroundColor: colors.background,
      borderRadius: radii.md,
      borderWidth: 1,
      borderColor: colors.border,
      padding: spacing[6],
      gap: spacing[6],
    },
    section: {
      gap: spacing[3],
    },
    genres: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
    },
    genreBadge: {
      // Additional styles if needed
    },
  });

  return (
    <View style={styles.container}>
      <View style={styles.section}>
        <Label>Sort By</Label>
        <Select value={sortBy} onValueChange={onSortChange}>
          <SelectTrigger>
            <SelectValue placeholder="Select sort option" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="rating">Highest Rated</SelectItem>
            <SelectItem value="year">Newest First</SelectItem>
            <SelectItem value="title">Title A-Z</SelectItem>
          </SelectContent>
        </Select>
      </View>

      <View style={styles.section}>
        <Label>Minimum Rating: {minRating.toFixed(1)}</Label>
        <Slider
          value={[minRating]}
          onValueChange={onRatingChange}
          min={0}
          max={10}
          step={0.5}
        />
      </View>

      <View style={styles.section}>
        <Label>Genres</Label>
        <View style={styles.genres}>
          {genres.map((genre) => {
            const isSelected = selectedGenres.includes(genre);
            return (
              <Pressable key={genre} onPress={() => onGenreToggle(genre)}>
                <Badge
                  variant={isSelected ? "default" : "secondary"}
                  style={styles.genreBadge}
                >
                  {genre}
                </Badge>
              </Pressable>
            );
          })}
        </View>
      </View>
    </View>
  );
}

export default FilterPanel;