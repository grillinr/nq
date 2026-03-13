import React from 'react';
import { View, TouchableOpacity, StyleSheet, ScrollView, Text, Platform } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { BlurView } from 'expo-blur';
import * as Haptics from 'expo-haptics';
import { useTheme } from './ui/theme-provider';
import { spacing, radii, fontSize, fontWeights } from './ui/tokens';
import { MediaType } from '../types';

interface MediaTypeFilterProps {
  selectedTypes: MediaType[];
  onFilterChange: (types: MediaType[]) => void;
}

const typeOptions = [
  { value: 'movie' as const, icon: 'film-outline' as const, label: 'Movies' },
  { value: 'tv' as const, icon: 'tv-outline' as const, label: 'TV' },
  { value: 'book' as const, icon: 'book-outline' as const, label: 'Books' },
  { value: 'music' as const, icon: 'musical-notes-outline' as const, label: 'Music' },
  { value: 'game' as const, icon: 'game-controller-outline' as const, label: 'Games' },
];

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    wrapper: {
      borderRadius: radii.full,
      overflow: 'hidden',
    },
    androidWrapper: {
      borderRadius: radii.full,
      backgroundColor:
        colors.background === '#000000' ? 'rgba(28,28,30,0.92)' : 'rgba(255,255,255,0.92)',
      overflow: 'hidden',
    },
    blurInner: {
      flexDirection: 'row',
      alignItems: 'center',
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[2],
      gap: spacing[2],
    },
    scrollContent: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[2],
      paddingRight: spacing[1],
    },
    typeButton: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[1],
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[2],
      borderRadius: radii.full,
    },
    typeButtonSelected: {
      backgroundColor: colors.primary,
    },
    typeButtonUnselected: {
      backgroundColor: 'transparent',
    },
    typeLabel: {
      fontSize: fontSize.sm,
      fontWeight: fontWeights.medium,
      color: colors.mutedForeground,
    },
    typeLabelSelected: {
      color: '#ffffff',
    },
    clearButton: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[1],
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[2],
      borderRadius: radii.full,
      borderWidth: 1,
      borderColor: colors.border,
    },
    clearText: {
      fontSize: fontSize.sm,
      fontWeight: fontWeights.medium,
      color: colors.mutedForeground,
    },
    divider: {
      width: 1,
      height: 18,
      backgroundColor: colors.border,
      marginHorizontal: spacing[1],
    },
  });

function MediaTypeFilter({ selectedTypes, onFilterChange }: MediaTypeFilterProps) {
  const { colors, resolved } = useTheme();
  const styles = React.useMemo(() => createStyles(colors), [colors]);

  const handleTypeToggle = (type: MediaType) => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    if (selectedTypes.includes(type)) {
      onFilterChange(selectedTypes.filter(t => t !== type));
    } else {
      onFilterChange([...selectedTypes, type]);
    }
  };

  const handleClearFilters = () => {
    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
    onFilterChange([]);
  };

  const innerContent = (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      contentContainerStyle={styles.scrollContent}
    >
      {typeOptions.map(option => {
        const isSelected = selectedTypes.includes(option.value);
        return (
          <TouchableOpacity
            key={option.value}
            style={[
              styles.typeButton,
              isSelected ? styles.typeButtonSelected : styles.typeButtonUnselected,
            ]}
            onPress={() => handleTypeToggle(option.value)}
            accessibilityRole="button"
            accessibilityLabel={`Filter by ${option.label}`}
            accessibilityState={{ selected: isSelected }}
          >
            <Ionicons
              name={option.icon}
              size={16}
              color={isSelected ? '#ffffff' : colors.mutedForeground}
            />
            <Text style={[styles.typeLabel, isSelected && styles.typeLabelSelected]}>
              {option.label}
            </Text>
          </TouchableOpacity>
        );
      })}
      {selectedTypes.length > 0 && (
        <>
          <View style={styles.divider} />
          <TouchableOpacity
            style={styles.clearButton}
            onPress={handleClearFilters}
            accessibilityRole="button"
            accessibilityLabel="Clear filters"
          >
            <Ionicons name="close" size={14} color={colors.mutedForeground} />
            <Text style={styles.clearText}>Clear</Text>
          </TouchableOpacity>
        </>
      )}
    </ScrollView>
  );

  if (Platform.OS === 'android') {
    return (
      <View style={styles.androidWrapper}>
        <View style={styles.blurInner}>{innerContent}</View>
      </View>
    );
  }

  return (
    <BlurView intensity={70} tint={resolved === 'dark' ? 'dark' : 'light'} style={styles.wrapper}>
      <View style={styles.blurInner}>{innerContent}</View>
    </BlurView>
  );
}

export default MediaTypeFilter;
