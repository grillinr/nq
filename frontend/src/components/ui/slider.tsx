import React, { useMemo } from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { spacing, fontSize, fontWeights, radii, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface SliderProps {
  value: number[];
  onValueChange: (value: number[]) => void;
  min?: number;
  max?: number;
  step?: number;
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    slider: {
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'center',
    },
    button: {
      padding: spacing[2],
      backgroundColor: colors.muted,
      borderRadius: radii.sm,
    },
    valueContainer: {
      marginHorizontal: spacing[4],
      minWidth: 40,
      alignItems: 'center',
    },
    value: {
      fontSize: fontSize.base,
      color: colors.foreground,
      fontWeight: fontWeights.medium,
    },
  });
}

function Slider({ value, onValueChange, min = 0, max = 10, step = 0.5 }: SliderProps) {
  const currentValue = value[0];
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);

  const decrease = () => {
    const newValue = Math.max(min, currentValue - step);
    onValueChange([newValue]);
  };

  const increase = () => {
    const newValue = Math.min(max, currentValue + step);
    onValueChange([newValue]);
  };

  return (
    <View style={styles.slider}>
      <Pressable onPress={decrease} style={styles.button}>
        <Ionicons name="remove" size={16} color={colors.foreground} />
      </Pressable>
      <View style={styles.valueContainer}>
        <Text style={styles.value}>{currentValue.toFixed(1)}</Text>
      </View>
      <Pressable onPress={increase} style={styles.button}>
        <Ionicons name="add" size={16} color={colors.foreground} />
      </Pressable>
    </View>
  );
}

export default Slider;
