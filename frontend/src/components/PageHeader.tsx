import React, { useEffect } from 'react';
import { Text, StyleSheet } from 'react-native';
import { Image } from 'expo-image';
import { Ionicons } from '@expo/vector-icons';
import Animated, { useSharedValue, useAnimatedStyle, withTiming } from 'react-native-reanimated';
import { useTheme } from './ui/theme-provider';
import { spacing, fontSize, fontWeights, zIndex } from './ui/tokens';

interface PageHeaderProps {
  title: string;
  icon?: keyof typeof Ionicons.glyphMap;
  visible?: boolean;
}

// Hoist module-level so the require is not called on every render.
const LOGO_SOURCE = require('../../assets/images/nq-logo.svg');

export const BAR_CONTENT_HEIGHT = 44; // icon + text row

/** Returns the total header height (safe-area top + bar content row). */
export function useHeaderHeight() {
  return BAR_CONTENT_HEIGHT;
}

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      height: BAR_CONTENT_HEIGHT,
      backgroundColor: colors.background,
      paddingHorizontal: spacing[4],
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[2],
      zIndex: zIndex.header,
      borderBottomWidth: StyleSheet.hairlineWidth,
      borderBottomColor: colors.border,
    },
    logo: {
      width: 24,
      height: 24,
    },
    title: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.semibold,
      color: colors.primary,
    },
  });

function PageHeader({ title, icon, visible = true }: PageHeaderProps) {
  const { colors } = useTheme();
  const styles = React.useMemo(() => createStyles(colors), [colors]);
  const translateY = useSharedValue(0);

  useEffect(() => {
    translateY.value = withTiming(visible ? 0 : -BAR_CONTENT_HEIGHT, { duration: 250 });
  }, [visible, translateY]);

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateY: translateY.value }],
  }));

  return (
    <Animated.View
      style={[styles.container, animatedStyle]}
      pointerEvents={visible ? 'auto' : 'none'}
    >
      <Image source={LOGO_SOURCE} style={styles.logo} contentFit="contain" />
      {icon && <Ionicons name={icon} size={22} color={colors.primary} />}
      <Text style={styles.title}>{title}</Text>
    </Animated.View>
  );
}

export default PageHeader;
