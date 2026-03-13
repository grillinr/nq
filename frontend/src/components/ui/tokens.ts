// Design tokens (originally from web styles, converted to React Native values)

// Shared colors that are identical across light and dark themes
export const sharedColors = {
  primary: '#2663d9',
  primaryTonal: '#5d85d7',
  primaryForeground: '#ffffff',
  secondary: '#f59e0b',
  secondaryTonal: '#f9c873',
  secondaryForeground: '#ffffff',
  input: 'transparent',
  chart3: '#06b6d4',
  chart4: '#f59e0b',
  star: '#FFD700',
};

export const lightColors = {
  // Light theme with improved contrast
  background: '#f8f9fa',
  foreground: '#1e293b',
  card: '#ffffff',
  cardForeground: '#1e293b',
  popover: '#ffffff',
  popoverForeground: '#1e293b',
  primary: sharedColors.primary,
  primaryForeground: sharedColors.primaryForeground,
  secondary: sharedColors.secondary,
  secondaryForeground: sharedColors.secondaryForeground,
  muted: '#e2e8f0',
  mutedForeground: '#52677a',
  destructive: '#dc2626',
  destructiveForeground: '#ffffff',
  destructiveBackground: '#fef2f2',
  warning: '#f59e0b',
  border: '#cbd5e1',
  input: sharedColors.input,
  inputBackground: '#f8f9fa',
  switchBackground: '#cbd5e1',
  chart1: '#3b82f6',
  chart2: '#8b5cf6',
  chart3: sharedColors.chart3,
  chart4: sharedColors.chart4,
  chart5: '#ef4444',
  star: sharedColors.star,
};

export const darkColors = {
  // Dark theme values
  background: '#071328',
  foreground: '#e6eef8',
  card: '#102541',
  cardForeground: '#e6eef8',
  popover: '#071328',
  popoverForeground: '#e6eef8',
  primary: sharedColors.primary,
  primaryForeground: sharedColors.primaryForeground,
  secondary: '#102541',
  secondaryForeground: '#cbd5e1',
  muted: '#0f172a',
  mutedForeground: '#9aa4b2',
  destructive: '#fb7185',
  destructiveForeground: '#071328',
  destructiveBackground: '#0b1220',
  warning: '#fbbf24',
  border: '#666666',
  input: sharedColors.input,
  inputBackground: '#0b1220',
  switchBackground: '#374151',
  chart1: '#60a5fa',
  chart2: '#a78bfa',
  chart3: sharedColors.chart3,
  chart4: sharedColors.chart4,
  chart5: '#fb7185',
  star: sharedColors.star,
};

export type ColorPalette = typeof lightColors;

export const fontWeights = {
  normal: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: 'bold' as const,
};

export const radii = {
  sm: 4,
  md: 6,
  lg: 8,
  xl: 12,
  full: 9999,
};

export const spacing = {
  1: 4,
  2: 8,
  3: 12,
  4: 16,
  5: 20,
  6: 24,
  8: 32,
};

export const fontSize = {
  xs: 12,
  sm: 14,
  base: 16,
  lg: 18,
  xl: 20,
  '2xl': 24,
};

export const lineHeight = {
  sm: 18,
  md: 22,
  lg: 24,
};

export const zIndex = {
  header: 1000,
  modal: 999,
};

// Layout constants
export const layout = {
  headerHeight: 140,
  historyHeaderOffset: 110,
};

// Shadow presets — use createShadows(colors) to get theme-aware shadow styles
export const createShadows = (colors: ColorPalette) => ({
  subtle: {
    shadowColor: colors.foreground,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  card: {
    shadowColor: colors.border,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  raised: {
    shadowColor: colors.border,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.2,
    shadowRadius: 8,
    elevation: 4,
  },
  modal: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 8,
  },
});
