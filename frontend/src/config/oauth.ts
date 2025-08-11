// OAuth Configuration
// You'll need to set up OAuth applications for these services and update the client IDs

export const OAUTH_CONFIG = {
  spotify: {
    clientId: 'YOUR_SPOTIFY_CLIENT_ID', // Get from https://developer.spotify.com/dashboard
    clientSecret: 'YOUR_SPOTIFY_CLIENT_SECRET',
    redirectUri: 'nq://auth/spotify',
    scopes: [
      'user-read-private',
      'user-read-email',
      'playlist-read-private',
      'user-library-read',
      'user-top-read',
      'user-read-recently-played'
    ]
  },
  google: {
    clientId: 'YOUR_GOOGLE_CLIENT_ID', // Get from https://console.developers.google.com/
    clientSecret: 'YOUR_GOOGLE_CLIENT_SECRET',
    redirectUri: 'nq://auth/google',
    scopes: {
      books: ['https://www.googleapis.com/auth/books'],
      youtube: ['https://www.googleapis.com/auth/youtube.readonly']
    }
  }
};

// Instructions for setting up OAuth:
//
// 1. SPOTIFY:
//    - Go to https://developer.spotify.com/dashboard
//    - Create a new app
//    - Add redirect URI: nq://auth/spotify
//    - Copy Client ID and Client Secret
//
// 2. GOOGLE (for Google Books and YouTube):
//    - Go to https://console.developers.google.com/
//    - Create a new project
//    - Enable Google Books API and YouTube Data API v3
//    - Create OAuth 2.0 credentials
//    - Add redirect URI: nq://auth/google
//    - Copy Client ID and Client Secret
//
// 3. Update the clientId and clientSecret values above
// 4. For production, store these securely (not in the app bundle)
