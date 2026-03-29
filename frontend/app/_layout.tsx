import { Stack } from 'expo-router';
import { StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { ApolloProvider } from '@apollo/client/react';
import { useEffect } from 'react';
import { ThemeProvider, useTheme } from '../src/components/ui/theme-provider';
import { AuthProvider } from '../src/lib/AuthContext';
import { apolloClient, initApolloCache } from '../src/lib/apolloClient';
import { logError } from '../src/lib/logger';

// Enable Apollo Client error messages in development
if (__DEV__) {
  import('@apollo/client/dev').then(({ loadErrorMessages, loadDevMessages }) => {
    loadDevMessages();
    loadErrorMessages();
  });
}

const styles = StyleSheet.create({
  fill: {
    flex: 1,
  },
});

function ThemedSafeAreaView({ children }: { children: React.ReactNode }) {
  const { colors } = useTheme();
  return (
    <SafeAreaView
      style={[styles.fill, { backgroundColor: colors.background }]}
      edges={['top', 'left', 'right']}
    >
      {children}
    </SafeAreaView>
  );
}

export default function RootLayout() {
  useEffect(() => {
    // Restore the Apollo cache from AsyncStorage in the background.
    // We don't block rendering — Apollo will re-render when the cache is hydrated.
    initApolloCache().catch(logError);
  }, []);

  return (
    <ApolloProvider client={apolloClient}>
      <ThemeProvider>
        <ThemedSafeAreaView>
          <AuthProvider>
            <Stack screenOptions={{ headerShown: false }} />
          </AuthProvider>
        </ThemedSafeAreaView>
      </ThemeProvider>
    </ApolloProvider>
  );
}
