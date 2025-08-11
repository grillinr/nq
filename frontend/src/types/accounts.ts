export interface AccountConnection {
  id: string;
  service: ServiceType;
  username?: string;
  email?: string;
  isConnected: boolean;
  lastSync?: Date;
  profileImage?: string;
  accessToken?: string;
  refreshToken?: string;
  expiresAt?: Date;
}

export type ServiceType = 
  | 'spotify'
  | 'goodreads'
  | 'letterboxd'
  | 'googlebooks'
  | 'youtube'
  | 'steam'
  | 'libby'
  | 'pocketcasts';

export interface ServiceConfig {
  name: string;
  displayName: string;
  description: string;
  icon: string;
  color: string;
  authUrl?: string;
  clientId?: string;
  scopes?: string[];
  requiresOAuth: boolean;
  apiBaseUrl?: string;
}

export const SERVICE_CONFIGS: Record<ServiceType, ServiceConfig> = {
  spotify: {
    name: 'spotify',
    displayName: 'Spotify',
    description: 'Connect your Spotify account to sync your music library and playlists',
    icon: '🎵',
    color: '#1DB954',
    authUrl: 'https://accounts.spotify.com/authorize',
    clientId: 'YOUR_SPOTIFY_CLIENT_ID',
    scopes: ['user-read-private', 'user-read-email', 'playlist-read-private', 'user-library-read'],
    requiresOAuth: true,
    apiBaseUrl: 'https://api.spotify.com/v1'
  },
  goodreads: {
    name: 'goodreads',
    displayName: 'Goodreads',
    description: 'Connect your Goodreads account to sync your reading lists and reviews',
    icon: '📚',
    color: '#663300',
    requiresOAuth: false,
    apiBaseUrl: 'https://www.goodreads.com'
  },
  letterboxd: {
    name: 'letterboxd',
    displayName: 'Letterboxd',
    description: 'Connect your Letterboxd account to sync your movie watchlist and reviews',
    icon: '🎬',
    color: '#00C030',
    requiresOAuth: false,
    apiBaseUrl: 'https://api.letterboxd.com'
  },
  googlebooks: {
    name: 'googlebooks',
    displayName: 'Google Books',
    description: 'Connect your Google Books account to sync your reading progress',
    icon: '📖',
    color: '#4285F4',
    authUrl: 'https://accounts.google.com/oauth/authorize',
    clientId: 'YOUR_GOOGLE_CLIENT_ID',
    scopes: ['https://www.googleapis.com/auth/books'],
    requiresOAuth: true,
    apiBaseUrl: 'https://www.googleapis.com/books/v1'
  },
  youtube: {
    name: 'youtube',
    displayName: 'YouTube',
    description: 'Connect your YouTube account to sync your subscriptions and watch history',
    icon: '📺',
    color: '#FF0000',
    authUrl: 'https://accounts.google.com/oauth/authorize',
    clientId: 'YOUR_GOOGLE_CLIENT_ID',
    scopes: ['https://www.googleapis.com/auth/youtube.readonly'],
    requiresOAuth: true,
    apiBaseUrl: 'https://www.googleapis.com/youtube/v3'
  },
  steam: {
    name: 'steam',
    displayName: 'Steam',
    description: 'Connect your Steam account to sync your game library and achievements',
    icon: '🎮',
    color: '#171a21',
    requiresOAuth: false,
    apiBaseUrl: 'https://api.steampowered.com'
  },
  libby: {
    name: 'libby',
    displayName: 'Libby',
    description: 'Connect your Libby account to sync your library books and reading progress',
    icon: '🏛️',
    color: '#2E5BBA',
    requiresOAuth: false,
    apiBaseUrl: 'https://libbyapp.com'
  },
  pocketcasts: {
    name: 'pocketcasts',
    displayName: 'Pocket Casts',
    description: 'Connect your podcast app to sync your subscriptions and listening history',
    icon: '🎧',
    color: '#F43E37',
    requiresOAuth: false,
    apiBaseUrl: 'https://api.pocketcasts.com'
  }
};
