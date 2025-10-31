import React from 'react';
import { Text, StyleSheet, TextProps } from 'react-native';
import { fontSize } from './tokens';
import { useTheme } from './ThemeProvider';

interface LabelProps extends TextProps {
  children: React.ReactNode;
}

function Label({ children, style, ...props }: LabelProps) {
  const { colors } = useTheme();
  const styles = StyleSheet.create({
    label: {
      fontSize: fontSize.sm,
      fontWeight: '500',
      color: colors.foreground,
      marginBottom: 4,
    },
  });

  return (
    <Text style={[styles.label, style]} {...props}>
      {children}
    </Text>
  );
}

export default Label;
