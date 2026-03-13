import React, { useMemo } from 'react';
import { Text, StyleSheet, TextProps } from 'react-native';
import { fontSize, fontWeights, spacing, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface LabelProps extends TextProps {
  children: React.ReactNode;
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    label: {
      fontSize: fontSize.sm,
      fontWeight: fontWeights.medium,
      color: colors.foreground,
      marginBottom: spacing[1],
    },
  });
}

function Label({ children, style, ...props }: LabelProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);

  return (
    <Text style={[styles.label, style]} {...props}>
      {children}
    </Text>
  );
}

export default Label;
