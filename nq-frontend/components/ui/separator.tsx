import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors } from './tokens';

interface SeparatorProps {
  orientation?: 'horizontal' | 'vertical';
}

function Separator({ orientation = 'horizontal' }: SeparatorProps) {
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