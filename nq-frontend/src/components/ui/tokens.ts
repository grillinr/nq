// Design tokens (originally from web styles, converted to React Native values)

export const lightColors = {
  // Light theme with improved contrast
  background: "#f8f9fa",
  foreground: "#1e293b",
  card: "#ffffff",
  "card-foreground": "#1e293b",
  popover: "#ffffff",
  "popover-foreground": "#1e293b",
  primary: "#2563eb",
  "primary-foreground": "#ffffff",
  secondary: "#f1f5f9",
  "secondary-foreground": "#1e293b",
  muted: "#e2e8f0",
  "muted-foreground": "#64748b",
  accent: "#e2e8f0",
  "accent-foreground": "#1e293b",
  destructive: "#dc2626",
  "destructive-foreground": "#ffffff",
  "destructive-background": "#fef2f2",
  border: "#e2e8f0",
  "primary-overlay": "rgba(37, 99, 235, 0.1)",
  "overlay-1": "rgba(0, 0, 0, 0.8)",
  input: "transparent",
  "input-background": "#f8f9fa",
  "switch-background": "#cbd5e1",
  ring: "#94a3b8",
  "chart-1": "#3b82f6",
  "chart-2": "#8b5cf6",
  "chart-3": "#06b6d4",
  "chart-4": "#f59e0b",
  "chart-5": "#ef4444",
  sidebar: "#f8f9fa",
  "sidebar-foreground": "#1e293b",
  "sidebar-primary": "#2563eb",
  "sidebar-primary-foreground": "#f8f9fa",
  "sidebar-accent": "#f1f5f9",
  "sidebar-accent-foreground": "#374151",
  "sidebar-border": "#e2e8f0",
  "sidebar-ring": "#94a3b8",
};

export const darkColors = {
  // Dark theme values
  background: "#071328",
  foreground: "#e6eef8",
  card: "#102441",
  "card-foreground": "#e6eef8",
  popover: "#071328",
  "popover-foreground": "#e6eef8",
  primary: "#60a5fa",
  "primary-foreground": "#ffffff",
  secondary: "#0f172a",
  "secondary-foreground": "#cbd5e1",
  muted: "#0f172a",
  "muted-foreground": "#9aa4b2",
  accent: "#0f172a",
  "accent-foreground": "#e6eef8",
  destructive: "#fb7185",
  "destructive-foreground": "#071328",
  "destructive-background": "#0b1220",
  border: "rgba(255, 255, 255, 0.06)",
  "primary-overlay": "rgba(7, 19, 40, 0.12)",
  "overlay-1": "rgba(255, 255, 255, 0.08)",
  input: "transparent",
  "input-background": "#0b1220",
  "switch-background": "#374151",
  ring: "#1f2937",
  "chart-1": "#60a5fa",
  "chart-2": "#a78bfa",
  "chart-3": "#06b6d4",
  "chart-4": "#f59e0b",
  "chart-5": "#fb7185",
  sidebar: "#050816",
  "sidebar-foreground": "#e6eef8",
  "sidebar-primary": "#60a5fa",
  "sidebar-primary-foreground": "#071328",
  "sidebar-accent": "#0b1220",
  "sidebar-accent-foreground": "#e6eef8",
  "sidebar-border": "#111827",
  "sidebar-ring": "#1f2937",
};

export type ColorPalette = typeof lightColors;

export const fontWeights = {
  medium: "500" as const,
  normal: "400" as const,
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
  "2xl": 24,
};
