import React, { useMemo } from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { useTheme } from '../../src/components/ui/theme-provider';
import PageHeader, { useHeaderHeight } from '../../src/components/PageHeader';
import { layout, spacing, ColorPalette } from '../../src/components/ui/tokens';

function createStyles(colors: ColorPalette, headerHeight: number) {
  return StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
      paddingTop: headerHeight,
      paddingBottom: layout.tabBarHeight,
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
  const headerHeight = useHeaderHeight();
  const styles = useMemo(() => createStyles(colors, headerHeight), [colors, headerHeight]);

  return (
    <View style={flexOne}>
      <PageHeader title="Friends" />
      <ScrollView style={styles.container} contentContainerStyle={styles.content}>
        <Text style={styles.placeholderText}>Friend activity feed is coming soon.</Text>
      </ScrollView>
    </View>
  );
}
