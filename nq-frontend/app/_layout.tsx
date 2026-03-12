import { Stack } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { ApolloProvider } from "@apollo/client/react";
import { ThemeProvider, useTheme } from "../src/components/ui/theme-provider";
import { AuthProvider } from "../src/lib/AuthContext";
import { apolloClient, initApolloCache } from "../src/lib/apolloClient";
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
            <Stack screenOptions={{ headerShown: false }} />
          </AuthProvider>
        </ThemedSafeAreaView>
      </ThemeProvider>
    </ApolloProvider>
  );
}
