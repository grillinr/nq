# NQ - Personal Dashboard App

A comprehensive personal dashboard app that allows users to connect and sync data from multiple services including Spotify, Goodreads, Letterboxd, Google Books, YouTube, Steam, Libby, and Pocket Casts.

## Features

### 🔗 Account Connections
- **Spotify**: Connect to sync music library, playlists, and listening history
- **Goodreads**: Connect to sync reading lists and book reviews
- **Letterboxd**: Connect to sync movie watchlist and ratings
- **Google Books**: Connect to sync reading progress and library
- **YouTube**: Connect to sync subscriptions and watch history
- **Steam**: Connect to sync game library and achievements
- **Libby**: Connect to sync library books and reading progress
- **Pocket Casts**: Connect to sync podcast subscriptions

### 🎯 Key Capabilities
- Secure OAuth authentication for supported services
- Manual connection for services without OAuth
- Local data storage with secure token management
- Real-time sync status and connection management
- Beautiful, intuitive user interface
- Cross-platform support (iOS, Android, Web)

## Project Structure

```
nq/
├── frontend/                 # React Native/Expo frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── screens/         # App screens
│   │   ├── services/        # Business logic and API services
│   │   ├── types/          # TypeScript type definitions
│   │   └── config/         # Configuration files
│   ├── assets/             # Images and static assets
│   └── package.json        # Frontend dependencies
├── backend/                 # Rust backend (optional)
├── docker-compose.yml       # Development environment setup
└── README.md               # This file
```

## Quick Start

### Prerequisites
- Node.js 18+ and npm/yarn
- Expo CLI (`npm install -g @expo/cli`)
- For OAuth services: Developer accounts (see setup guide)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd nq
   ```

2. **Install frontend dependencies**
   ```bash
   cd frontend
   npm install
   ```

3. **Configure OAuth services** (see `frontend/ACCOUNT_SETUP.md`)
   - Set up Spotify Developer account
   - Set up Google Cloud Console project
   - Update configuration files

4. **Start the development server**
   ```bash
   npm start
   ```

5. **Run on your device**
   - Scan QR code with Expo Go app (mobile)
   - Press 'w' for web version
   - Press 'a' for Android emulator
   - Press 'i' for iOS simulator

## Configuration

### OAuth Setup
Detailed setup instructions are available in `frontend/ACCOUNT_SETUP.md`. You'll need to:

1. Create developer accounts for OAuth services
2. Configure redirect URIs
3. Update client IDs and secrets
4. Set up API permissions

### Environment Variables
For production, use environment variables for sensitive data:

```bash
# .env
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

## Development

### Adding New Services
1. Add service configuration to `src/types/accounts.ts`
2. Implement connection logic in `src/services/accountService.ts`
3. Add UI components as needed
4. Update the accounts screen

### Code Style
- TypeScript for type safety
- Functional components with hooks
- Consistent naming conventions
- Proper error handling
- Comprehensive documentation

## Security Features

- Secure token storage using `expo-secure-store`
- OAuth 2.0 implementation for supported services
- Local data encryption
- No sensitive data in app bundle
- Secure redirect URI handling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Documentation**: Check `frontend/ACCOUNT_SETUP.md` for detailed setup
- **Issues**: Create an issue in the repository
- **Questions**: Open a discussion or contact the maintainers

## Roadmap

- [ ] Data synchronization scheduling
- [ ] Offline data caching
- [ ] Data export functionality
- [ ] Advanced analytics dashboard
- [ ] API rate limiting and optimization
- [ ] Multi-user support
- [ ] Cloud data backup

## Acknowledgments

- Expo team for the excellent development platform
- React Native community for the robust ecosystem
- Service providers for their APIs and documentation
