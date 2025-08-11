# NQ Frontend

A React Native mobile application built with Expo for the NQ project.

## Features

- **Navigation**: Stack-based navigation between screens
- **Home Screen**: Dashboard with quick action cards
- **Profile Screen**: User profile information display
- **Settings Screen**: App configuration options
- **Modern UI**: Clean, card-based design with shadows and proper spacing
- **TypeScript**: Full TypeScript support for better development experience

## Prerequisites

- Node.js (v16 or higher)
- npm or yarn
- Expo CLI (optional, but recommended)
- iOS Simulator (macOS) or Android Emulator

## Installation

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

## Running the App

### Start the development server:
```bash
npm start
```

### Run on specific platforms:

**Android:**
```bash
npm run android
```

**iOS (macOS only):**
```bash
npm run ios
```

**Web:**
```bash
npm run web
```

## Project Structure

```
frontend/
├── App.tsx              # Main app component with navigation
├── index.ts             # Entry point
├── package.json         # Dependencies and scripts
├── app.json            # Expo configuration
├── tsconfig.json       # TypeScript configuration
├── assets/             # Images and static assets
└── node_modules/       # Dependencies
```

## Dependencies

### Core Dependencies
- `expo`: Expo framework
- `react`: React library
- `react-native`: React Native framework
- `expo-status-bar`: Status bar management

### Navigation Dependencies
- `@react-navigation/native`: Core navigation library
- `@react-navigation/stack`: Stack navigator
- `react-native-screens`: Native screen components
- `react-native-safe-area-context`: Safe area handling
- `react-native-gesture-handler`: Gesture handling

### Development Dependencies
- `typescript`: TypeScript compiler
- `@types/react`: React type definitions
- `@types/react-native-vector-icons`: Vector icons type definitions

## Development

### Adding New Screens

1. Create a new screen component in `App.tsx` or in a separate file
2. Add the screen to the `RootStackParamList` type
3. Add a new `Stack.Screen` in the navigator

Example:
```typescript
type RootStackParamList = {
  Home: undefined;
  Profile: undefined;
  Settings: undefined;
  NewScreen: undefined; // Add new screen
};

// Add the screen component
<Stack.Screen 
  name="NewScreen" 
  component={NewScreenComponent} 
  options={{ title: 'New Screen' }}
/>
```

### Styling

The app uses React Native's StyleSheet API with a consistent design system:
- Primary color: `#1a1a1a` (dark gray)
- Background: `#f5f5f5` (light gray)
- Cards: `#ffffff` (white)
- Text: `#1a1a1a` (dark) and `#666` (medium gray)

### Navigation

The app uses React Navigation v6 with stack navigation. Each screen can access the navigation object to move between screens:

```typescript
navigation.navigate('ScreenName');
```

## Building for Production

### Android APK:
```bash
expo build:android
```

### iOS IPA:
```bash
expo build:ios
```

## Troubleshooting

### Common Issues

1. **Metro bundler issues**: Try clearing the cache:
   ```bash
   npx expo start --clear
   ```

2. **Dependencies issues**: Delete node_modules and reinstall:
   ```bash
   rm -rf node_modules
   npm install
   ```

3. **iOS build issues**: Make sure you're on macOS and have Xcode installed

### Getting Help

- Check the [Expo documentation](https://docs.expo.dev/)
- Review [React Native documentation](https://reactnative.dev/)
- Check the [React Navigation documentation](https://reactnavigation.org/)

## Contributing

1. Follow the existing code style and patterns
2. Use TypeScript for all new code
3. Test your changes on both iOS and Android
4. Update this README if you add new features or change the setup

## License

This project is part of the NQ project. See the main project LICENSE file for details.
