import React from 'react';
import { View, TouchableOpacity, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTheme } from './theme-provider';
import { fontSize, fontWeights, spacing } from './tokens';

interface StarRatingProps {
  value: number; // 0-5 in 0.5 increments
  onChange?: (value: number) => void;
  size?: 'sm' | 'md' | 'lg';
  readonly?: boolean;
  showValue?: boolean;
  style?: any;
}

export function StarRating({
  value,
  onChange,
  size = 'md',
  readonly = false,
  showValue = false,
  style,
}: StarRatingProps) {
  const { colors } = useTheme();

  const sizeMap = { sm: 16, md: 24, lg: 32 };
  const iconSize = sizeMap[size];

  const handlePress = (starIndex: number, isHalf: boolean) => {
    if (readonly || !onChange) return;
    const newValue = starIndex + (isHalf ? 0.5 : 1);
    onChange(newValue === value ? 0 : newValue); // Toggle off if same
  };

  const renderStar = (index: number) => {
    const filled = value >= index + 1;
    const half = value >= index + 0.5 && value < index + 1;

    const iconName = filled ? 'star' : half ? 'star-half' : 'star-outline';

    if (readonly) {
      return (
        <Ionicons
          key={index}
          name={iconName}
          size={iconSize}
          color={colors.star}
        />
      );
    }

    return (
      <View key={index} style={[styles.starContainer, { width: iconSize, height: iconSize }]}>
        {/* Render the star icon once */}
        <Ionicons
          name={iconName}
          size={iconSize}
          color={colors.star}
          style={styles.starIcon}
          pointerEvents="none"
        />
        {/* Left half touch target */}
        <TouchableOpacity
          onPress={() => handlePress(index, true)}
          style={[styles.halfStar, styles.leftHalf, { width: iconSize / 2, height: iconSize }]}
          activeOpacity={0.7}
        />
        {/* Right half touch target */}
        <TouchableOpacity
          onPress={() => handlePress(index, false)}
          style={[styles.halfStar, styles.rightHalf, { width: iconSize / 2, height: iconSize }]}
          activeOpacity={0.7}
        />
      </View>
    );
  };

  return (
    <View style={[styles.container, style]}>
      <View style={styles.stars}>
        {[0, 1, 2, 3, 4].map(renderStar)}
      </View>
      {showValue && (
        <Text style={[styles.valueText, { color: colors.foreground }]}>
          {value.toFixed(1)}
        </Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[2],
  },
  stars: {
    flexDirection: 'row',
    gap: spacing[1],
  },
  starContainer: {
    position: 'relative',
  },
  starIcon: {
    position: 'absolute',
    left: 0,
    top: 0,
  },
  halfStar: {
    position: 'absolute',
    top: 0,
  },
  leftHalf: {
    left: 0,
  },
  rightHalf: {
    right: 0,
  },
  valueText: {
    fontSize: fontSize.base,
    fontWeight: fontWeights.semibold,
  },
});
