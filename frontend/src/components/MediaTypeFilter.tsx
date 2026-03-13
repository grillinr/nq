import React from 'react';
import { View, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTheme } from './ui/theme-provider';
import { spacing, radii } from './ui/tokens';
import { MediaType } from '../types';

interface MediaTypeFilterProps {
  selectedTypes: MediaType[];
  onFilterChange: (types: MediaType[]) => void;
}

const typeOptions = [
  { value: 'movie' as const, icon: 'film-outline' as const, label: 'Filter by Movies' },
  { value: 'tv' as const, icon: 'tv-outline' as const, label: 'Filter by TV Shows' },
  { value: 'book' as const, icon: 'book-outline' as const, label: 'Filter by Books' },
  { value: 'music' as const, icon: 'musical-notes-outline' as const, label: 'Filter by Music' },
  { value: 'game' as const, icon: 'game-controller-outline' as const, label: 'Filter by Games' },
];

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
      alignItems: 'center',
    },
    typeButton: {
      padding: spacing[3],
      borderRadius: radii.lg,
      borderWidth: 1,
      borderColor: colors.border,
      backgroundColor: colors.background,
    },
    typeButtonSelected: {
      backgroundColor: colors.primary,
      borderColor: colors.primary,
    },
    clearButton: {
      padding: spacing[3],
      borderRadius: radii.lg,
      borderWidth: 1,
      borderColor: colors.border,
      backgroundColor: colors.background,
    },
  });

function MediaTypeFilter({ selectedTypes, onFilterChange }: MediaTypeFilterProps) {
  const { colors } = useTheme();
  const styles = React.useMemo(() => createStyles(colors), [colors]);

  const handleTypeToggle = (type: MediaType) => {
    if (selectedTypes.includes(type)) {
      // Remove type from selection
      onFilterChange(selectedTypes.filter(t => t !== type));
    } else {
      // Add type to selection
      onFilterChange([...selectedTypes, type]);
    }
  };

  const handleClearFilters = () => {
    onFilterChange([]);
  };

  return (
    <View style={styles.container}>
      {typeOptions.map(option => {
        const isSelected = selectedTypes.includes(option.value);
        return (
          <TouchableOpacity
            key={option.value}
            style={[styles.typeButton, isSelected && styles.typeButtonSelected]}
            onPress={() => handleTypeToggle(option.value)}
            accessibilityRole="button"
            accessibilityLabel={option.label}
            accessibilityState={{ selected: isSelected }}
          >
            <Ionicons
              name={option.icon}
              size={20}
              color={isSelected ? colors.primaryForeground : colors.foreground}
            />
          </TouchableOpacity>
        );
      })}
      {selectedTypes.length > 0 && (
        <TouchableOpacity
          style={styles.clearButton}
          onPress={handleClearFilters}
          accessibilityRole="button"
          accessibilityLabel="Clear filters"
        >
          <Ionicons name="close-circle-outline" size={20} color={colors.foreground} />
        </TouchableOpacity>
      )}
    </View>
  );
}

export default MediaTypeFilter;
