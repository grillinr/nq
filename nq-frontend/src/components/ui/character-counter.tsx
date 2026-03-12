import React from 'react';
import { Text, StyleSheet } from 'react-native';
import { useTheme } from './theme-provider';
import { fontSize, fontWeights } from './tokens';

interface CharacterCounterProps {
  current: number;
  max: number;
  style?: any;
}

export function CharacterCounter({ current, max, style }: CharacterCounterProps) {
  const { colors } = useTheme();
  
  const isNearLimit = current > max * 0.8;
  const isOverLimit = current > max;
  
  const color = isOverLimit 
    ? colors.destructive 
    : isNearLimit 
    ? colors.warning 
    : colors.mutedForeground;

  return (
    <Text style={[styles.counter, { color }, style]}>
      {current}/{max}
    </Text>
  );
}

const styles = StyleSheet.create({
  counter: {
    fontSize: fontSize.xs,
    fontWeight: fontWeights.medium,
  },
});
