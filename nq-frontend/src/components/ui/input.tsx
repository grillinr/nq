import React from 'react';
import { TextInput, StyleSheet, ViewStyle } from 'react-native';
import { cn } from './utils';
import { radii, spacing, fontSize } from './tokens';
import { useTheme } from './ThemeProvider';

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

  const computed = StyleSheet.create({
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

  const inputStyle = cn([
    computed.base,
    disabled && computed.disabled,
    style,
  ]);

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
