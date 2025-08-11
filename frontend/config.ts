// App Configuration
export const CONFIG = {
  // API Configuration
  API_BASE_URL: process.env.EXPO_PUBLIC_API_BASE_URL || 'http://localhost:8000',
  API_TIMEOUT: 10000,
  
  // App Configuration
  APP_NAME: 'NQ',
  APP_VERSION: '1.0.0',
  
  // Feature Flags
  ENABLE_ANALYTICS: false,
  ENABLE_CRASH_REPORTING: false,
  
  // UI Configuration
  THEME: {
    PRIMARY_COLOR: '#1a1a1a',
    BACKGROUND_COLOR: '#f5f5f5',
    CARD_BACKGROUND: '#ffffff',
    TEXT_PRIMARY: '#1a1a1a',
    TEXT_SECONDARY: '#666',
  },
  
  // Navigation
  INITIAL_ROUTE: 'Home',
} as const;

export default CONFIG;
