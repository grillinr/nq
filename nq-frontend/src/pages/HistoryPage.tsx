import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { fontSize, spacing } from '../components/ui/tokens';
import { useTheme } from '../components/ui/ThemeProvider';

function HistoryPage() {
  const { colors } = useTheme();

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
    },
    section: {
      marginBottom: spacing[6],
    },
    list: {
      gap: spacing[4],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: '600',
      color: colors.foreground,
      marginBottom: spacing[3],
    },
    placeholderText: {
        color: colors.foreground,
        textAlign: 'center',
        marginTop: spacing[4],
    }
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Recently Viewed</Text>
        <View style={styles.list}>
          <Text style={styles.placeholderText}>History functionality coming soon...</Text>
        </View>
      </View>
    </ScrollView>
  );
}

export default HistoryPage;
