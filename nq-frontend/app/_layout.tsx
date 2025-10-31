import { Slot } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { MediaProvider } from "../lib/MediaContext";
import { ThemeProvider, useTheme } from "../src/components/ui/ThemeProvider";

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
    <ThemeProvider>
      <ThemedSafeAreaView>
        <MediaProvider>
          <Slot />
        </MediaProvider>
      </ThemedSafeAreaView>
    </ThemeProvider>
  );
}
