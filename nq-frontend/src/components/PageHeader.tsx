import React, { useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Image } from 'expo-image';
import { Ionicons } from '@expo/vector-icons';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withTiming,
  SharedValue,
} from 'react-native-reanimated';
import { useTheme } from './ui/ThemeProvider';
import { spacing, fontSize } from './ui/tokens';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  icon?: keyof typeof Ionicons.glyphMap;
  visible?: boolean;
  onTranslateYChange?: (translateY: SharedValue<number>) => void;
}

function PageHeader({
  title,
  subtitle,
  icon,
  visible = true,
  onTranslateYChange,
}: PageHeaderProps) {
  const { colors } = useTheme();
  const translateY = useSharedValue(0);

  useEffect(() => {
    if (onTranslateYChange) {
      onTranslateYChange(translateY);
    }
  }, [onTranslateYChange, translateY]);

  useEffect(() => {
    translateY.value = withTiming(visible ? 0 : -120, { duration: 300 });
  }, [visible, translateY]);

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateY: translateY.value }],
  }));

  const styles = StyleSheet.create({
    container: {
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      backgroundColor: colors.background,
      paddingTop: spacing[1],
      paddingBottom: spacing[1],
      alignItems: 'center',
      zIndex: 1000,
      borderBottomWidth: 1,
      borderBottomColor: colors.border,
    },
    logo: {
      width: 48,
      height: 48,
      marginBottom: 0,
    },
    headerContent: {
      alignItems: 'center',
      gap: spacing[1],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: '600',
      color: colors.primary,
      textAlign: 'center',
    },
    subtitle: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      textAlign: 'center',
      paddingHorizontal: spacing[4],
    },
  });

  return (
    <Animated.View
      style={[styles.container, animatedStyle]}
      pointerEvents={visible ? 'auto' : 'none'}
    >
      <Image
        source={require('../../assets/images/nq-logo.svg')}
        style={styles.logo}
        contentFit="contain"
      />
      <View style={styles.headerContent}>
        {icon && <Ionicons name={icon} size={32} color={colors.primary} />}
        <Text style={styles.title}>{title}</Text>
        {subtitle && <Text style={styles.subtitle}>{subtitle}</Text>}
      </View>
    </Animated.View>
  );
}

export default PageHeader;
