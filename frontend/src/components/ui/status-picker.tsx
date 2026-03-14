import React, { useMemo } from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTheme } from './theme-provider';
import { fontSize, fontWeights, radii, spacing } from './tokens';

export type ActivityStatusId = 1 | 2 | 3;

interface StatusOption {
  id: ActivityStatusId;
  name: string;
  description: string;
  icon: keyof typeof Ionicons.glyphMap;
}

const STATUS_OPTIONS: StatusOption[] = [
  { id: 1, name: 'Planned', description: 'Want to watch/read/play', icon: 'list-outline' },
  { id: 2, name: 'In Progress', description: 'Currently experiencing', icon: 'play-outline' },
  { id: 3, name: 'Completed', description: 'Finished', icon: 'checkmark-circle-outline' },
];

interface StatusPickerProps {
  value: ActivityStatusId;
  onChange: (statusId: ActivityStatusId) => void;
  style?: any;
}

function createStyles() {
  return StyleSheet.create({
    container: {
      gap: spacing[3],
    },
    option: {
      flexDirection: 'row',
      padding: spacing[4],
      borderRadius: radii.lg,
      borderWidth: 2,
      alignItems: 'center',
      gap: spacing[3],
    },
    textContainer: {
      flex: 1,
    },
    name: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.semibold,
      marginBottom: spacing[1],
    },
    description: {
      fontSize: fontSize.sm,
    },
  });
}

export function StatusPicker({ value, onChange, style }: StatusPickerProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(), []);

  return (
    <View style={[styles.container, style]}>
      {STATUS_OPTIONS.map(option => (
        <TouchableOpacity
          key={option.id}
          style={[
            styles.option,
            {
              backgroundColor: value === option.id ? colors.primary : colors.card,
              borderColor: value === option.id ? colors.primary : colors.border,
            },
          ]}
          onPress={() => onChange(option.id)}
          activeOpacity={0.7}
        >
          <Ionicons
            name={option.icon}
            size={24}
            color={value === option.id ? colors.primaryForeground : colors.foreground}
          />
          <View style={styles.textContainer}>
            <Text
              style={[
                styles.name,
                { color: value === option.id ? colors.primaryForeground : colors.foreground },
              ]}
            >
              {option.name}
            </Text>
            <Text
              style={[
                styles.description,
                { color: value === option.id ? colors.primaryForeground : colors.mutedForeground },
              ]}
            >
              {option.description}
            </Text>
          </View>
        </TouchableOpacity>
      ))}
    </View>
  );
}
