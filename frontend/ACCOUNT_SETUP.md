# Account Connections Setup Guide

This guide will help you set up connections to various services in the NQ app.

## Prerequisites

- Node.js and npm/yarn installed
- Expo CLI installed (`npm install -g @expo/cli`)
- Basic understanding of OAuth applications

## Quick Start

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Configure OAuth applications (see sections below)
3. Update the OAuth configuration
4. Run the app:
   ```bash
   npm start
   ```

## OAuth Service Setup

### 1. Spotify

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Click "Create App"
3. Fill in the app details:
   - App name: `NQ App` (or your preferred name)
   - App description: `Personal dashboard app`
   - Redirect URI: `nq://auth/spotify`
   - Website: `http://localhost:3000` (for development)
4. Click "Save"
5. Copy the Client ID and Client Secret
6. Update `src/config/oauth.ts` with your credentials

### 2. Google (for Google Books and YouTube)

1. Go to [Google Cloud Console](https://console.developers.google.com/)
2. Create a new project or select existing one
3. Enable the following APIs:
   - Google Books API
   - YouTube Data API v3
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
5. Configure the OAuth consent screen if prompted
6. Set application type to "Mobile application"
7. Add redirect URI: `nq://auth/google`
8. Copy the Client ID and Client Secret
9. Update `src/config/oauth.ts` with your credentials

## Manual Service Setup

### Goodreads
- No OAuth required
- Enter your Goodreads username when connecting
- Note: Goodreads API access is limited

### Letterboxd
- No OAuth required
- Enter your Letterboxd username when connecting
- Note: Letterboxd doesn't have a public API

### Steam
- No OAuth required
- Enter your Steam username when connecting
- Note: Steam API requires an API key for full access

### Libby
- No OAuth required
- Enter your library card number or username
- Note: Libby doesn't have a public API

### Pocket Casts
- No OAuth required
- Enter your Pocket Casts username or email
- Note: Pocket Casts API access is limited

## Configuration Files

### Update OAuth Configuration

Edit `src/config/oauth.ts`:

```typescript
export const OAUTH_CONFIG = {
  spotify: {
    clientId: 'your_actual_spotify_client_id',
    clientSecret: 'your_actual_spotify_client_secret',
    // ... rest of config
  },
  google: {
    clientId: 'your_actual_google_client_id',
    clientSecret: 'your_actual_google_client_secret',
    // ... rest of config
  }
};
```

### Environment Variables (Recommended)

For production, use environment variables:

1. Create `.env` file:
   ```
   SPOTIFY_CLIENT_ID=your_spotify_client_id
   SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
   GOOGLE_CLIENT_ID=your_google_client_id
   GOOGLE_CLIENT_SECRET=your_google_client_secret
   ```

2. Install `react-native-dotenv`:
   ```bash
   npm install react-native-dotenv
   ```

3. Update the config to use environment variables

## Testing

1. Start the app: `npm start`
2. Navigate to "Account Connections"
3. Try connecting to different services
4. Check the console for any errors

## Troubleshooting

### Common Issues

1. **OAuth redirect errors**: Ensure redirect URIs match exactly
2. **Client ID not found**: Verify credentials are correctly copied
3. **API quota exceeded**: Check Google Cloud Console quotas
4. **Network errors**: Ensure internet connection and API endpoints are accessible

### Debug Mode

Enable debug logging in the account service:

```typescript
// In src/services/accountService.ts
console.log('OAuth URL:', authUrl);
console.log('OAuth result:', result);
```

## Security Notes

- Never commit OAuth credentials to version control
- Use environment variables for production
- Implement proper token refresh logic
- Store tokens securely using expo-secure-store
- Consider implementing PKCE for additional security

## Next Steps

After setting up account connections:

1. Implement data synchronization logic
2. Add error handling and retry mechanisms
3. Implement token refresh for OAuth services
4. Add data caching and offline support
5. Implement user preferences for sync frequency

## Support

For issues with specific services:
- Spotify: [Spotify Developer Support](https://developer.spotify.com/support/)
- Google: [Google Cloud Support](https://cloud.google.com/support/)
- General: Check the service's official documentation
