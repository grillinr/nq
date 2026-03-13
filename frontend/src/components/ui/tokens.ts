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
  // Light theme — pure white base, iOS system feel
  background: '#ffffff',
  foreground: '#1e293b',
  card: '#ffffff',
  cardForeground: '#1e293b',
  popover: '#ffffff',
  popoverForeground: '#1e293b',
  primary: sharedColors.primary,
  primaryForeground: sharedColors.primaryForeground,
  secondary: sharedColors.secondary,
  secondaryForeground: sharedColors.secondaryForeground,
  muted: '#f2f2f7',
  mutedForeground: '#52677a',
  destructive: '#dc2626',
  destructiveForeground: '#ffffff',
  destructiveBackground: '#fef2f2',
  warning: '#f59e0b',
  border: '#e5e5ea',
  input: sharedColors.input,
  inputBackground: '#f2f2f7',
  switchBackground: '#cbd5e1',
  chart1: '#3b82f6',
  chart2: '#8b5cf6',
  chart3: sharedColors.chart3,
  chart4: sharedColors.chart4,
  chart5: '#ef4444',
  star: sharedColors.star,
};

export const darkColors = {
  // Dark theme — pure black base, iOS system feel
  background: '#000000',
  foreground: '#e6eef8',
  card: '#1c1c1e',
  cardForeground: '#e6eef8',
  popover: '#1c1c1e',
  popoverForeground: '#e6eef8',
  primary: sharedColors.primary,
  primaryForeground: sharedColors.primaryForeground,
  secondary: '#1c1c1e',
  secondaryForeground: '#cbd5e1',
  muted: '#1c1c1e',
  mutedForeground: '#9aa4b2',
  destructive: '#fb7185',
  destructiveForeground: '#000000',
  destructiveBackground: '#1c1c1e',
  warning: '#fbbf24',
  border: '#38383a',
  input: sharedColors.input,
  inputBackground: '#2c2c2e',
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
  sm: 6,
  md: 10,
  lg: 14,
  xl: 20,
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
  headerHeight: 88, // STATUS_BAR_HEIGHT(44) + BAR_CONTENT_HEIGHT(44) — matches PageHeader
  historyHeaderOffset: 110,
  tabBarHeight: 104, // 64px pill + 24px bottom offset + 16px extra breathing room
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
