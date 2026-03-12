import React, { useMemo } from 'react';
import { TextInput, StyleSheet, ViewStyle } from 'react-native';
import { flattenStyles } from './utils';
import { radii, spacing, fontSize, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface InputProps {
  placeholder?: string;
  value?: string;
  defaultValue?: string;
  onChangeText?: (text: string) => void;
  style?: ViewStyle;
  disabled?: boolean;
  multiline?: boolean;
  numberOfLines?: number;
  keyboardType?: 'default' | 'numeric' | 'email-address' | 'phone-pad';
  maxLength?: number;
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    base: {
      height: 36,
      borderWidth: 1,
      borderColor: colors.border,
      borderRadius: radii.md,
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[1],
      fontSize: fontSize.sm,
      backgroundColor: colors.inputBackground,
      color: colors.foreground,
    },
    disabled: {
      opacity: 0.5,
    },
  });
}

function Input({
  placeholder,
  value,
  defaultValue,
  onChangeText,
  style,
  disabled = false,
  multiline = false,
  numberOfLines = 1,
  keyboardType = 'default',
  ...props
}: InputProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const inputStyle = flattenStyles([styles.base, disabled && styles.disabled, style]);

  return (
    <TextInput
      style={inputStyle}
      placeholder={placeholder}
      value={value}
      defaultValue={defaultValue}
      onChangeText={onChangeText}
      editable={!disabled}
      placeholderTextColor={colors.mutedForeground}
      multiline={multiline}
      numberOfLines={numberOfLines}
      keyboardType={keyboardType}
      {...props}
    />
  );
}

export default Input;
