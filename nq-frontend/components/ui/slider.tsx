import React from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, fontSize } from './tokens';

interface SliderProps {
  value: number[];
  onValueChange: (value: number[]) => void;
  min?: number;
  max?: number;
  step?: number;
}

function Slider({ value, onValueChange, min = 0, max = 10, step = 0.5 }: SliderProps) {
  const currentValue = value[0];

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

const styles = StyleSheet.create({
  slider: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  button: {
    padding: spacing[2],
    backgroundColor: colors.muted,
    borderRadius: 4,
  },
  valueContainer: {
    marginHorizontal: spacing[4],
    minWidth: 40,
    alignItems: 'center',
  },
  value: {
    fontSize: fontSize.base,
    color: colors.foreground,
    fontWeight: '500',
  },
});