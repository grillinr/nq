# Porting Figma UI to React Native

This document describes the process of porting a web-based Figma UI to React Native and the decisions made during the port.

## Overview

The original `figma_ui` is a React web app using Tailwind CSS, Radix UI, and other web-specific libraries. We ported it to React Native by:

- Converting CSS classes to `StyleSheet` objects
- Replacing web-only components with RN equivalents
- Adapting browser APIs to RN APIs
- Using Expo and React Navigation for navigation

## Key Changes

### Styling
- Extracted design tokens from `globals.css` to `app/components/ui/tokens.ts`
- Created `cn` utility in `utils.ts` for merging styles (replaces `clsx` + `tailwind-merge`)
- Converted Tailwind classes to inline `StyleSheet` objects

### Components
- **Button**: Uses `Pressable` with variants
- **Input**: Uses `TextInput` with custom styling
- **Card**: Simple `View` with shadow
- **Avatar**: `Image` or fallback `View`
- **Badge**: `Text` with background
- **Modal**: Uses `react-native-modal`
- **ImageWithFallback**: RN `Image` with error handling

### Navigation
- Replaced web routing with React Navigation bottom tabs
- Home, Add, History, Friends, Account screens

### Assets & Icons
- Replaced Lucide icons with `@expo/vector-icons` (Ionicons)
- Images use RN `Image` component

### Libraries Replaced
- Radix UI → Custom RN components or `react-native-modal`
- Tailwind → `StyleSheet`
- `lucide-react` → `@expo/vector-icons`
- Browser globals (window, document) → RN equivalents or removed

## Preserved
- `lib/createMedia.ts`: Kept as-is, already RN-compatible
- Core functionality: Adding media, displaying lists

## Files Added
- `app/components/ui/` - UI primitives (button, input, card, avatar, badge, modal, label, separator, switch, tabs, select, slider)
- `app/components/` - Screens (HomePage, AddMediaPage, AccountPage, HistoryPage, FriendsPage) and FilterPanel
- `app/App.tsx` - Main app with navigation

## Files Removed
- `app/components/MediaInput.tsx` - Replaced by full UI

## Dependencies Added
- `@react-navigation/native`, `@react-navigation/bottom-tabs`
- `react-native-screens`, `react-native-safe-area-context`
- `react-native-gesture-handler`, `react-native-reanimated`
- `react-native-svg`, `react-native-modal`
- `react-native-snap-carousel`, `react-native-chart-kit`
- `@react-native-async-storage/async-storage`
- `@expo/vector-icons`

## Testing
- Built and ran on Expo web for QA
- Basic functionality verified: navigation, adding media, displaying cards
- ESLint passes with no errors

## Completed
- All screens ported from figma_ui
- Image sizing fixed with resizeMode='cover'
- Custom UI components implemented
- Navigation wired for all screens

## Future Improvements
- Implement carousel for recommendations (if needed)
- Add dark mode support
- Enhance form validation and error handling
- Test on iOS/Android devices