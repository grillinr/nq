import React from 'react';
import { Text, StyleSheet, ViewStyle, TextStyle } from 'react-native';
import { cn } from './utils';
import { colors, radii, fontSize } from './tokens';

type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline';

interface BadgeProps {
  variant?: BadgeVariant;
  children: React.ReactNode;
  style?: ViewStyle;
}

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
    backgroundColor: 'transparent',
    borderColor: colors.border,
    borderWidth: 1,
  },
};

const textVariants: Record<BadgeVariant, TextStyle> = {
  default: {
    color: colors['primary-foreground'],
  },
  secondary: {
    color: colors['secondary-foreground'],
  },
  destructive: {
    color: colors['destructive-foreground'],
  },
  outline: {
    color: colors.foreground,
  },
};

function Badge({ variant = 'default', children, style }: BadgeProps) {
  const badgeStyle = cn([styles.base, badgeVariants[variant], style]);
  const textStyle = cn([styles.text, textVariants[variant]]);

  return (
    <Text style={badgeStyle}>
      {typeof children === 'string' ? (
        <Text style={textStyle}>{children}</Text>
      ) : (
        children
      )}
    </Text>
  );
}

export default Badge;

const styles = StyleSheet.create({
  base: {
    borderRadius: radii.sm,
    paddingHorizontal: 8,
    paddingVertical: 2,
    alignSelf: 'flex-start',
  },
  text: {
    fontSize: fontSize.xs,
    fontWeight: '500',
  },
});