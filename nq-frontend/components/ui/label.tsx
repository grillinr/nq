import React from 'react';
import { Text, StyleSheet, TextProps } from 'react-native';
import { fontSize, colors } from './tokens';

interface LabelProps extends TextProps {
  children: React.ReactNode;
}

function Label({ children, style, ...props }: LabelProps) {
  return (
    <Text style={[styles.label, style]} {...props}>
      {children}
    </Text>
  );
}

export default Label;

const styles = StyleSheet.create({
  label: {
    fontSize: fontSize.sm,
    fontWeight: '500',
    color: colors.foreground,
    marginBottom: 4,
  },
});