import React, { useMemo } from 'react';
import { View, StyleSheet, ViewStyle } from 'react-native';
import { flattenStyles } from './utils';
import { createShadows, radii, spacing, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface CardProps {
  children: React.ReactNode;
  style?: ViewStyle;
}

function createStyles(colors: ColorPalette) {
  const shadows = createShadows(colors);
  return StyleSheet.create({
    base: {
      backgroundColor: colors.card,
      borderRadius: radii.lg,
      padding: spacing[4],
      ...shadows.card,
    },
  });
}

function Card({ children, style }: CardProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const cardStyle = flattenStyles([styles.base, style]);
  return <View style={cardStyle}>{children}</View>;
}

export default Card;
