import React from 'react';
import { View, StyleSheet, ViewStyle } from 'react-native';
import { flattenStyles } from './utils';
import { radii } from './tokens';
import { useTheme } from './ThemeProvider';

interface CardProps {
  children: React.ReactNode;
  style?: ViewStyle;
}

function Card({ children, style }: CardProps) {
  const { colors } = useTheme();
  const computed = StyleSheet.create({
    base: {
      backgroundColor: colors.card,
      borderRadius: radii.lg,
      padding: 16,
      shadowColor: colors.border,
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.1,
      shadowRadius: 4,
      elevation: 2,
    },
  });

  const cardStyle = flattenStyles([computed.base, style]);

  return <View style={cardStyle}>{children}</View>;
}

export default Card;
