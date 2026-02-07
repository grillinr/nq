import { Slot } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { ApolloProvider } from "@apollo/client/react";
import { MediaProvider } from "../lib/MediaContext";
import { ThemeProvider, useTheme } from "../src/components/ui/ThemeProvider";
import { AuthProvider } from "../lib/AuthContext";
import { apolloClient } from "../lib/apolloClient";

function ThemedSafeAreaView({ children }: { children: React.ReactNode }) {
  const { colors } = useTheme();
  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }} edges={['top', 'left', 'right']}>
      {children}
    </SafeAreaView>
  );
}

export default function RootLayout() {
  return (
    <ApolloProvider client={apolloClient}>
      <ThemeProvider>
        <ThemedSafeAreaView>
          <AuthProvider>
            <MediaProvider>
              <Slot />
            </MediaProvider>
          </AuthProvider>
        </ThemedSafeAreaView>
      </ThemeProvider>
    </ApolloProvider>
  );
}
