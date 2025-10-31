import React from 'react';
import { View, StyleSheet } from 'react-native';
import { useTheme } from './ThemeProvider';

interface SeparatorProps {
  orientation?: 'horizontal' | 'vertical';
}

function Separator({ orientation = 'horizontal' }: SeparatorProps) {
  const { colors } = useTheme();
  const styles = StyleSheet.create({
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

  return (
    <View
      style={[
        styles.separator,
        orientation === 'vertical' ? styles.vertical : styles.horizontal,
      ]}
    />
  );
}

export default Separator;
