import { StyleSheet } from 'react-native';

// Simple style merger for React Native, replacing clsx + tailwind-merge
// Takes an array of style objects and merges them
export function flattenStyles(styles: (object | undefined | null | false)[]): object {
  return StyleSheet.flatten(styles.filter(Boolean));
}

/**
 * Returns a hex color string with the given opacity applied.
 * `hex` must be a 6-digit hex string (e.g. '#2663d9').
 * `opacity` must be in the range [0, 1].
 */
export function withOpacity(hex: string, opacity: number): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r},${g},${b},${opacity})`;
}
