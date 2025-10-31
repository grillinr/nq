// Design tokens (originally from web styles, converted to React Native values)

export const colors = {
  // Light theme
  background: '#ffffff',
  foreground: '#0f172a', // oklch(0.145 0 0)
  card: '#ffffff',
  'card-foreground': '#0f172a',
  popover: '#ffffff',
  'popover-foreground': '#0f172a',
  primary: '#030213',
  'primary-foreground': '#ffffff',
  secondary: '#f1f5f9', // oklch(0.95 0.0058 264.53)
  'secondary-foreground': '#030213',
  muted: '#ececf0',
  'muted-foreground': '#717182',
  accent: '#e9ebef',
  'accent-foreground': '#030213',
  destructive: '#d4183d',
  'destructive-foreground': '#ffffff',
  border: 'rgba(0, 0, 0, 0.1)',
  input: 'transparent',
  'input-background': '#f3f3f5',
  'switch-background': '#cbced4',
  ring: '#cbd5e1', // oklch(0.708 0 0)
  'chart-1': '#3b82f6', // oklch(0.646 0.222 41.116)
  'chart-2': '#8b5cf6', // oklch(0.6 0.118 184.704)
  'chart-3': '#06b6d4', // oklch(0.398 0.07 227.392)
  'chart-4': '#f59e0b', // oklch(0.828 0.189 84.429)
  'chart-5': '#ef4444', // oklch(0.769 0.188 70.08)
  sidebar: '#fafafa', // oklch(0.985 0 0)
  'sidebar-foreground': '#0f172a',
  'sidebar-primary': '#030213',
  'sidebar-primary-foreground': '#fafafa',
  'sidebar-accent': '#f5f5f5', // oklch(0.97 0 0)
  'sidebar-accent-foreground': '#1e293b', // oklch(0.205 0 0)
  'sidebar-border': '#e2e8f0', // oklch(0.922 0 0)
  'sidebar-ring': '#cbd5e1',

  // Dark theme (can add later if needed)
  // ...
};

export const fontWeights = {
  medium: '500' as const,
  normal: '400' as const,
};

export const radii = {
  sm: 4,
  md: 6,
  lg: 8,
  xl: 12,
};

export const spacing = {
  1: 4,
  2: 8,
  3: 12,
  4: 16,
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