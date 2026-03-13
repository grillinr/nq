import React from 'react';
import { Animated, Easing, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { radii } from './ui/tokens';
import { useTheme } from './ui/theme-provider';

const SHIMMER_DURATION = 1200;

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
  const shimmerColors =
    resolved === 'dark'
      ? (['rgba(255,255,255,0)', 'rgba(255,255,255,0.12)', 'rgba(255,255,255,0)'] as const)
      : (['rgba(0,0,0,0)', 'rgba(0,0,0,0.06)', 'rgba(0,0,0,0)'] as const);

  const cardBg = resolved === 'dark' ? colors.inputBackground : colors.muted;

  return (
    <View
      style={[
        {
          width: '100%',
          borderRadius: radii.lg,
          overflow: 'hidden',
          backgroundColor: cardBg,
          aspectRatio,
        },
        style,
      ]}
    >
      <Animated.View
        style={{
          position: 'absolute',
          top: 0,
          bottom: 0,
          width: 160,
          transform: [{ translateX }],
        }}
      >
        <LinearGradient
          colors={shimmerColors}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 0 }}
          style={StyleSheet.absoluteFill}
        />
      </Animated.View>
    </View>
  );
}

export default React.memo(MediaCoverSkeleton);
