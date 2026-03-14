import React from 'react';
import { Text, StyleSheet } from 'react-native';
import { useTheme } from './theme-provider';
import { fontSize, fontWeights } from './tokens';

interface CharacterCounterProps {
  current: number;
  max: number;
  style?: any;
}

const styles = StyleSheet.create({
  counter: {
    fontSize: fontSize.xs,
    fontWeight: fontWeights.medium,
  },
});

export function CharacterCounter({ current, max, style }: CharacterCounterProps) {
  const { colors } = useTheme();

  const isNearLimit = current > max * 0.8;
  const isOverLimit = current > max;

  let color: string;
  if (isOverLimit) {
    color = colors.destructive;
  } else if (isNearLimit) {
    color = colors.warning;
  } else {
    color = colors.mutedForeground;
  }

  return (
    <Text style={[styles.counter, { color }, style]}>
      {current}/{max}
    </Text>
  );
}
