import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useTheme } from '../components/ui/ThemeProvider';
import { spacing, fontSize } from '../components/ui/tokens';

function FriendsPage() {
  const { colors } = useTheme();

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
    },
    header: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[3],
      marginBottom: spacing[4],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: '600',
      color: colors.foreground,
    },
    subtitle: {
      fontSize: fontSize.sm,
      color: colors['muted-foreground'],
    },
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <Ionicons name="people" size={24} color={colors.primary} />
        <View>
          <Text style={styles.title}>Friends</Text>
          <Text style={styles.subtitle}>See what your friends are rating and sharing</Text>
        </View>
      </View>

      <Text style={{ color: colors['muted-foreground'] }}>Friend activity feed is coming soon.</Text>
    </ScrollView>
  );
}

export default FriendsPage;
