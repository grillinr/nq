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
  card: '#f2f2f2',
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

export const sizes = {
  1: 4,
  2: 8,
  3: 12,
  4: 16,
  5: 20,
  6: 24,
  7: 28,
  8: 32,
  9: 36,
  10: 40,
  11: 44,
  12: 48,
  13: 52,
  14: 56,
  16: 64,
};

export const zIndex = {
  header: 1000,
  modal: 999,
};

// Layout constants
export const layout = {
  headerHeight: 88, // STATUS_BAR_HEIGHT(44) + BAR_CONTENT_HEIGHT(44) — matches PageHeader
  historyHeaderOffset: 110,
  // tab button (sizes[11]) + equal margin top/bottom (sizes[1] each) = pill height (sizes[13])
  // pill (sizes[13]) + bottom offset (spacing[6]) + breathing room (spacing[4]) = tabBarHeight
  tabBarHeight: sizes[13] + spacing[6] + spacing[4],
};

// Shimmer gradient tuples — used by skeleton components for the animated highlight pass.
// The tuple is [transparent edge, highlight peak, transparent edge].
export const shimmerColors = {
  dark: ['rgba(255,255,255,0)', 'rgba(255,255,255,0.12)', 'rgba(255,255,255,0)'] as const,
  light: ['rgba(0,0,0,0)', 'rgba(0,0,0,0.06)', 'rgba(0,0,0,0)'] as const,
  // Slightly more subtle variant used in text-line skeletons
  darkSubtle: ['rgba(255,255,255,0)', 'rgba(255,255,255,0.1)', 'rgba(255,255,255,0)'] as const,
  lightSubtle: ['rgba(0,0,0,0)', 'rgba(0,0,0,0.05)', 'rgba(0,0,0,0)'] as const,
};

// Android glass-pill background — replicates the iOS BlurView effect on Android.
// Uses card/background colors with opacity so they stay in sync with the theme tokens.
export const androidGlassBg = {
  dark: 'rgba(28,28,30,0.95)' as const, // darkColors.card (#1c1c1e) at 95 %
  light: 'rgba(255,255,255,0.95)' as const, // lightColors.background (#ffffff) at 95 %
  // Slightly lower opacity used by the MediaTypeFilter pill
  darkFilter: 'rgba(28,28,30,0.92)' as const,
  lightFilter: 'rgba(255,255,255,0.92)' as const,
};

// Tab-bar active highlight background (primary brand color with low opacity).
export const tabActiveBg = {
  dark: 'rgba(38,99,217,0.25)' as const, // sharedColors.primary at 25 %
  light: 'rgba(38,99,217,0.12)' as const, // sharedColors.primary at 12 %
};

// Tab-bar inactive icon tint (foreground color with reduced opacity).
export const tabInactiveTint = {
  dark: 'rgba(255,255,255,0.4)' as const, // white at 40 %
  light: 'rgba(0,0,0,0.3)' as const, // black at 30 %
};

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
    shadowColor: colors.foreground,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 8,
  },
});
