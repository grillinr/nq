import React from "react";
import { View, Text, StyleSheet } from "react-native";
import { useTheme } from "../components/ui/ThemeProvider";
import { Button } from "../components/ui/button";
import { spacing, fontSize } from "../components/ui/tokens";

interface AuthPromptPageProps {
  onLogin: () => Promise<void> | void;
  onSignup: () => Promise<void> | void;
}

function AuthPromptPage({ onLogin, onSignup }: AuthPromptPageProps) {
  const { colors } = useTheme();

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
      justifyContent: "center",
      alignItems: "center",
      padding: spacing[6],
      gap: spacing[4],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: "600",
      color: colors.primary,
      textAlign: "center",
    },
    subtitle: {
      fontSize: fontSize.base,
      color: colors["muted-foreground"],
      textAlign: "center",
    },
    actions: {
      width: "100%",
      gap: spacing[3],
    },
  });

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Welcome to NQ</Text>
      <Text style={styles.subtitle}>
        Create an account or log in to start adding media.
      </Text>
      <View style={styles.actions}>
        <Button onPress={onSignup}>Create account</Button>
        <Button variant="outline" onPress={onLogin}>
          Log in
        </Button>
      </View>
    </View>
  );
}

export default AuthPromptPage;
