import React from 'react';
import { Text, StyleSheet, ViewStyle, TextStyle } from 'react-native';
import { flattenStyles } from './utils';
import { radii, fontSize, fontWeights, spacing } from './tokens';
import { useTheme } from './theme-provider';

type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline';

interface BadgeProps {
  variant?: BadgeVariant;
  children: React.ReactNode;
  style?: ViewStyle;
}

const styles = StyleSheet.create({
  base: {
    borderRadius: radii.sm,
    paddingHorizontal: spacing[2],
    paddingVertical: 2,
    alignSelf: 'flex-start',
  },
  text: {
    fontSize: fontSize.sm,
    fontWeight: fontWeights.semibold,
  },
});

function Badge({ variant = 'default', children, style }: BadgeProps) {
  const { colors } = useTheme();

  const badgeVariants: Record<BadgeVariant, ViewStyle> = {
    default: {
      backgroundColor: colors.primary,
    },
    secondary: {
      backgroundColor: colors.secondary,
    },
    destructive: {
      backgroundColor: colors.destructive,
    },
    outline: {
      backgroundColor: colors.input,
      borderColor: colors.border,
      borderWidth: 1,
    },
  };

  const textVariants: Record<BadgeVariant, TextStyle> = {
    default: {
      color: colors.primaryForeground,
    },
    secondary: {
      color: colors.secondaryForeground,
    },
    destructive: {
      color: colors.destructiveForeground,
    },
    outline: {
      color: colors.foreground,
    },
  };

  const badgeStyle = flattenStyles([styles.base, badgeVariants[variant], style]);
  const textStyle = flattenStyles([styles.text, textVariants[variant]]);

  if (typeof children === 'string') {
    return <Text style={[badgeStyle, textStyle]}>{children}</Text>;
  }
  return <Text style={badgeStyle}>{children}</Text>;
}

export default Badge;
