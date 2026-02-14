import React from "react";
import { View, StyleSheet } from "react-native";
import AuthPromptPage from "../src/pages/AuthPromptPage";
import { useAuth } from "../lib/AuthContext";
import { useTheme } from "../src/components/ui/ThemeProvider";

export default function AuthPage() {
  const { login } = useAuth();
  const { colors } = useTheme();

  const handleLogin = async () => {
    await login();
  };

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <AuthPromptPage onLogin={handleLogin} onSignup={handleLogin} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
});
