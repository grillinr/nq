import React, { useMemo } from 'react';
import { View, StyleSheet } from 'react-native';
import { ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface SeparatorProps {
  orientation?: 'horizontal' | 'vertical';
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    separator: {
      backgroundColor: colors.border,
    },
    horizontal: {
      height: 1,
      width: '100%',
    },
    vertical: {
      width: 1,
      height: '100%',
    },
  });
}

function Separator({ orientation = 'horizontal' }: SeparatorProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);

  return (
    <View
      style={[styles.separator, orientation === 'vertical' ? styles.vertical : styles.horizontal]}
    />
  );
}

export default Separator;
