import React from 'react';
import { View, StyleSheet, ViewStyle } from 'react-native';
import { cn } from './utils';
import { colors, radii } from './tokens';

interface CardProps {
  children: React.ReactNode;
  style?: ViewStyle;
}

function Card({ children, style }: CardProps) {
  const cardStyle = cn([styles.base, style]);

  return <View style={cardStyle}>{children}</View>;
}

export default Card;

const styles = StyleSheet.create({
  base: {
    backgroundColor: colors.card,
    borderRadius: radii.lg,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
});