import React, { useMemo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useTheme } from '../src/components/ui/theme-provider';
import { Button } from '../src/components/ui/button';
import { spacing, fontSize, fontWeights, ColorPalette } from '../src/components/ui/tokens';
import { useAuth } from '../src/lib/AuthContext';

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
      justifyContent: 'center',
      alignItems: 'center',
      padding: spacing[6],
      gap: spacing[4],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: fontWeights.semibold,
      color: colors.primary,
      textAlign: 'center',
    },
    subtitle: {
      fontSize: fontSize.base,
      color: colors.mutedForeground,
      textAlign: 'center',
    },
    actions: {
      width: '100%',
      gap: spacing[3],
    },
  });
}

export default function AuthPage() {
  const { login } = useAuth();
  const { colors } = useTheme();
  const router = useRouter();
  const styles = useMemo(() => createStyles(colors), [colors]);

  const handleLogin = async () => {
    const success = await login();
    if (success) {
      router.replace('/(tabs)');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Welcome to NQ</Text>
      <Text style={styles.subtitle}>Create an account or log in to start adding media.</Text>
      <View style={styles.actions}>
        <Button onPress={handleLogin}>Create account</Button>
        <Button variant="outline" onPress={handleLogin}>
          Log in
        </Button>
      </View>
    </View>
  );
}
