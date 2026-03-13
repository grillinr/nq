import { StyleSheet } from 'react-native';

// Simple style merger for React Native, replacing clsx + tailwind-merge
// Takes an array of style objects and merges them
export function flattenStyles(styles: (object | undefined | null | false)[]): object {
  return StyleSheet.flatten(styles.filter(Boolean));
}
