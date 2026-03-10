import { Stack } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { ApolloProvider } from "@apollo/client/react";
import { MediaProvider } from "../lib/MediaContext";
import { ThemeProvider, useTheme } from "../src/components/ui/ThemeProvider";
import { AuthProvider } from "../lib/AuthContext";
import { apolloClient, initApolloCache } from "../lib/apolloClient";
import { useEffect } from "react";

function ThemedSafeAreaView({ children }: { children: React.ReactNode }) {
  const { colors } = useTheme();
  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }} edges={['top', 'left', 'right']}>
      {children}
    </SafeAreaView>
  );
}

export default function RootLayout() {
  useEffect(() => {
    // Restore the Apollo cache from AsyncStorage in the background.
    // We don't block rendering — Apollo will re-render when the cache is hydrated.
    initApolloCache().catch(console.error);
  }, []);

  return (
    <ApolloProvider client={apolloClient}>
      <ThemeProvider>
        <ThemedSafeAreaView>
          <AuthProvider>
            <MediaProvider>
              <Stack screenOptions={{ headerShown: false }} />
            </MediaProvider>
          </AuthProvider>
        </ThemedSafeAreaView>
      </ThemeProvider>
    </ApolloProvider>
  );
}
