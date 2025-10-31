import { Slot } from 'expo-router';
import { SafeAreaView } from 'react-native-safe-area-context';
import { MediaProvider } from '../lib/MediaContext';

export default function RootLayout() {
  return (
    <SafeAreaView style={{ flex: 1 }}>
      <MediaProvider>
        <Slot />
      </MediaProvider>
    </SafeAreaView>
  );
}
