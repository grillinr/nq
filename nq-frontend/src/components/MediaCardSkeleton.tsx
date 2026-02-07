import React from 'react';
import { View, StyleSheet, Animated, Easing, useWindowDimensions, StyleProp, ViewStyle } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { spacing, radii } from './ui/tokens';
import { useTheme } from './ui/ThemeProvider';

const SHIMMER_DURATION = 1200;

function MediaCardSkeleton() {
  const { colors } = useTheme();
  const { width } = useWindowDimensions();
  const translateX = React.useRef(new Animated.Value(0)).current;

  React.useEffect(() => {
    translateX.setValue(-width);
    const animation = Animated.loop(
      Animated.timing(translateX, {
        toValue: width,
        duration: SHIMMER_DURATION,
        easing: Easing.linear,
        useNativeDriver: true,
      })
    );
    animation.start();
    return () => animation.stop();
  }, [translateX, width]);

  const styles = React.useMemo(() => createStyles(colors), [colors]);

  return (
    <View style={styles.container}>
      <View style={styles.card}>
        <View style={styles.imageBlock}>
          <Shimmer translateX={translateX} style={styles.shimmer} />
        </View>
        <View style={styles.content}>
          <View style={[styles.line, styles.titleLine]}>
            <Shimmer translateX={translateX} style={styles.shimmer} />
          </View>
          <View style={[styles.line, styles.metaLine]}>
            <Shimmer translateX={translateX} style={styles.shimmer} />
          </View>
          <View style={styles.badgesRow}>
            <View style={[styles.badge, styles.badgeShort]}>
              <Shimmer translateX={translateX} style={styles.shimmer} />
            </View>
            <View style={[styles.badge, styles.badgeMedium]}>
              <Shimmer translateX={translateX} style={styles.shimmer} />
            </View>
            <View style={[styles.badge, styles.badgeShort]}>
              <Shimmer translateX={translateX} style={styles.shimmer} />
            </View>
          </View>
          <View style={[styles.line, styles.descLine]}>
            <Shimmer translateX={translateX} style={styles.shimmer} />
          </View>
          <View style={[styles.line, styles.descLineShort]}>
            <Shimmer translateX={translateX} style={styles.shimmer} />
          </View>
        </View>
      </View>
    </View>
  );
}

function Shimmer({ translateX, style }: { translateX: Animated.Value; style: StyleProp<ViewStyle> }) {
  return (
    <Animated.View style={[style, { transform: [{ translateX }] }]}> 
      <LinearGradient
        colors={['rgba(255,255,255,0)', 'rgba(255,255,255,0.25)', 'rgba(255,255,255,0)']}
        start={{ x: 0, y: 0 }}
        end={{ x: 1, y: 0 }}
        style={StyleSheet.absoluteFill}
      />
    </Animated.View>
  );
}

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      width: '100%',
    },
    card: {
      borderRadius: radii.lg,
      backgroundColor: colors.card,
      overflow: 'hidden',
      padding: spacing[4],
    },
    imageBlock: {
      aspectRatio: 2 / 3,
      backgroundColor: colors['input-background'],
      overflow: 'hidden',
    },
    content: {
      padding: spacing[4],
      gap: spacing[2],
    },
    line: {
      backgroundColor: colors['input-background'],
      borderRadius: radii.sm,
      overflow: 'hidden',
    },
    titleLine: {
      height: 18,
      width: '65%',
    },
    metaLine: {
      height: 12,
      width: '45%',
      marginTop: spacing[1],
    },
    badgesRow: {
      flexDirection: 'row',
      gap: spacing[1],
      marginTop: spacing[1],
      marginBottom: spacing[1],
    },
    badge: {
      height: 16,
      borderRadius: radii.sm,
      backgroundColor: colors['input-background'],
      overflow: 'hidden',
    },
    badgeShort: {
      width: 48,
    },
    badgeMedium: {
      width: 72,
    },
    descLine: {
      height: 12,
      width: '90%',
      marginTop: spacing[1],
    },
    descLineShort: {
      height: 12,
      width: '70%',
    },
    shimmer: {
      position: 'absolute',
      top: 0,
      bottom: 0,
      width: 160,
    },
  });

const MemoizedMediaCardSkeleton = React.memo(MediaCardSkeleton);
MemoizedMediaCardSkeleton.displayName = 'MediaCardSkeleton';

export default MemoizedMediaCardSkeleton;
