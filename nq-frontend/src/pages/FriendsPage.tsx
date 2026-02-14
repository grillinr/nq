import React from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { useTheme } from '../components/ui/ThemeProvider';
import PageHeader from '../components/PageHeader';
import { useScrollHeader } from '../hooks/useScrollHeader';
import { spacing } from '../components/ui/tokens';

function FriendsPage() {
  const { colors } = useTheme();
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
      paddingTop: spacing[4],
    },
    content: {
      padding: spacing[4],
      paddingTop: 140,
    },
    placeholderText: {
      color: colors.mutedForeground,
      textAlign: 'center',
      marginTop: spacing[6],
    },
  });

  return (
    <View style={{ flex: 1 }}>
      <PageHeader
        title="Friends"
        subtitle="See what your friends are rating and sharing"
        visible={isHeaderVisible}
      />
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.content}
        onScroll={handleHeaderScroll}
        scrollEventThrottle={16}
      >
        <Text style={styles.placeholderText}>Friend activity feed is coming soon.</Text>
      </ScrollView>
    </View>
  );
}

export default FriendsPage;
