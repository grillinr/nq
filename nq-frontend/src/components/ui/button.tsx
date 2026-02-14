import React from 'react';
import { Pressable, Text, StyleSheet, ViewStyle, TextStyle } from 'react-native';
import { cn } from './utils';
import { fontWeights, radii, spacing } from './tokens';
import { useTheme } from './ThemeProvider';

type ButtonVariant = 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
type ButtonSize = 'default' | 'sm' | 'lg' | 'icon';

interface ButtonProps {
  variant?: ButtonVariant;
  size?: ButtonSize;
  disabled?: boolean;
  onPress?: () => void;
  children: React.ReactNode;
  style?: ViewStyle;
}

export function Button({
  variant = 'default',
  size = 'default',
  disabled = false,
  onPress,
  children,
  style,
}: ButtonProps) {
  const { colors } = useTheme();

  const buttonVariants: Record<ButtonVariant, ViewStyle> = {
    default: {
      backgroundColor: colors.primary,
    },
    destructive: {
      backgroundColor: colors.destructive,
    },
    outline: {
      backgroundColor: colors.background,
      borderColor: colors.border,
      borderWidth: 1,
    },
    secondary: {
      backgroundColor: colors.secondary,
    },
    ghost: {
      backgroundColor: 'transparent',
    },
    link: {
      backgroundColor: 'transparent',
    },
  };

  const buttonSizes: Record<ButtonSize, ViewStyle> = {
    default: {
      height: 36,
      paddingHorizontal: spacing[4],
      paddingVertical: spacing[2],
    },
    sm: {
      height: 32,
      paddingHorizontal: spacing[3],
    },
    lg: {
      height: 40,
      paddingHorizontal: spacing[6],
    },
    icon: {
      width: 36,
      height: 36,
    },
  };

  const textVariants: Record<ButtonVariant, TextStyle> = {
    default: {
      color: colors.primaryForeground,
    },
    destructive: {
      color: colors.destructiveForeground,
    },
    outline: {
      color: colors.foreground,
    },
    secondary: {
      color: colors.secondaryForeground,
    },
    ghost: {
      color: colors.foreground,
    },
    link: {
      color: colors.primary,
      textDecorationLine: 'underline',
    },
  };

  const textSizes: Record<ButtonSize, TextStyle> = {
    default: {
      fontSize: 14,
      fontWeight: fontWeights.medium,
    },
    sm: {
      fontSize: 14,
    },
    lg: {
      fontSize: 16,
    },
    icon: {
      fontSize: 14,
    },
  };

  const buttonStyle = cn([
    styles.base,
    buttonVariants[variant],
    buttonSizes[size],
    disabled && styles.disabled,
    style,
  ]);

  const textStyle = cn([
    styles.textBase,
    textVariants[variant],
    textSizes[size],
    disabled && styles.disabledText,
  ]);

  return (
    <Pressable
      style={buttonStyle}
      onPress={disabled ? undefined : onPress}
      disabled={disabled}
    >
      {typeof children === 'string' ? (
        <Text style={textStyle}>{children}</Text>
      ) : (
        children
      )}
    </Pressable>
  );
}

export default Button;

const styles = StyleSheet.create({
  base: {
    borderRadius: radii.md,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
  },
  disabled: {
    opacity: 0.5,
  },
  textBase: {
    textAlign: 'center',
  },
  disabledText: {
    opacity: 0.5,
  },
});
