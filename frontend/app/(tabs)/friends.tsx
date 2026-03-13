import React, { useMemo } from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { useTheme } from '../../src/components/ui/theme-provider';
import PageHeader from '../../src/components/PageHeader';
import { useScrollHeader } from '../../src/hooks/useScrollHeader';
import { layout, spacing, ColorPalette } from '../../src/components/ui/tokens';

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
      paddingTop: spacing[4],
    },
    content: {
      padding: spacing[4],
      paddingTop: layout.headerHeight,
    },
    placeholderText: {
      color: colors.mutedForeground,
      textAlign: 'center',
      marginTop: spacing[6],
    },
  });
}

const flexOne = { flex: 1 };

export default function FriendsPage() {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);

  return (
    <View style={flexOne}>
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
