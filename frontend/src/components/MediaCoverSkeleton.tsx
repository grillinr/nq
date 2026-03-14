import React from 'react';
import { Animated, Easing, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { radii, shimmerColors } from './ui/tokens';
import { useTheme } from './ui/theme-provider';

const SHIMMER_DURATION = 1200;

const shimmerStyles = StyleSheet.create({
  shimmer: {
    position: 'absolute',
    top: 0,
    bottom: 0,
    width: 160,
  },
});

function createStyles(cardBg: string, aspectRatio: number) {
  return StyleSheet.create({
    container: {
      width: '100%',
      borderRadius: radii.lg,
      overflow: 'hidden',
      backgroundColor: cardBg,
      aspectRatio,
    },
  });
}

interface MediaCoverSkeletonProps {
  style?: StyleProp<ViewStyle>;
  aspectRatio?: number;
}

function MediaCoverSkeleton({ style, aspectRatio = 2 / 3 }: MediaCoverSkeletonProps) {
  const { colors, resolved } = useTheme();
  const translateX = React.useRef(new Animated.Value(0)).current;

  React.useEffect(() => {
    translateX.setValue(-160);
    const animation = Animated.loop(
      Animated.timing(translateX, {
        toValue: 160,
        duration: SHIMMER_DURATION,
        easing: Easing.linear,
        useNativeDriver: true,
      })
    );
    animation.start();
    return () => animation.stop();
  }, [translateX]);

  // Theme-aware shimmer: white highlight in dark, dark highlight in light
  const shimmerGradient = resolved === 'dark' ? shimmerColors.dark : shimmerColors.light;

  const cardBg = resolved === 'dark' ? colors.inputBackground : colors.muted;
  const styles = React.useMemo(() => createStyles(cardBg, aspectRatio), [cardBg, aspectRatio]);

  return (
    <View style={[styles.container, style]}>
      <Animated.View style={[shimmerStyles.shimmer, { transform: [{ translateX }] }]}>
        <LinearGradient
          colors={shimmerGradient}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 0 }}
          style={StyleSheet.absoluteFill}
        />
      </Animated.View>
    </View>
  );
}

export default React.memo(MediaCoverSkeleton);
