import React, { useEffect, useRef } from 'react';
import { Animated, Pressable, StyleSheet, View, ViewStyle } from 'react-native';
import ImageWithFallback from './ui/image-with-fallback';
import { useTheme } from './ui/theme-provider';
import { radii } from './ui/tokens';

interface MediaCoverCardProps {
  title: string;
  image: string;
  onPress?: () => void;
  aspectRatio?: number;
  style?: ViewStyle;
  isEnriching?: boolean;
}

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    pressable: {
      borderRadius: radii.lg,
      overflow: 'hidden',
      backgroundColor: colors.inputBackground,
    },
    imageContainer: {
      width: '100%',
      backgroundColor: colors.inputBackground,
    },
    image: {
      width: '100%',
      height: '100%',
    },
    shimmerOverlay: {
      ...StyleSheet.absoluteFillObject,
      borderRadius: radii.lg,
    },
  });

function MediaCoverCard({
  title,
  image,
  onPress,
  aspectRatio = 2 / 3,
  style,
  isEnriching = false,
}: MediaCoverCardProps) {
  const { colors } = useTheme();
  const styles = React.useMemo(() => createStyles(colors), [colors]);

  const shimmerOpacity = useRef(new Animated.Value(0.3)).current;
  useEffect(() => {
    if (!isEnriching) {
      shimmerOpacity.setValue(0);
      return undefined;
    }
    const loop = Animated.loop(
      Animated.sequence([
        Animated.timing(shimmerOpacity, { toValue: 0.6, duration: 700, useNativeDriver: true }),
        Animated.timing(shimmerOpacity, { toValue: 0.2, duration: 700, useNativeDriver: true }),
      ])
    );
    loop.start();
    return () => loop.stop();
  }, [isEnriching, shimmerOpacity]);

  return (
    <Pressable
      onPress={onPress}
      style={[styles.pressable, style]}
      accessibilityRole="button"
      accessibilityLabel={title}
    >
      <View style={[styles.imageContainer, { aspectRatio }]}>
        <ImageWithFallback src={image} alt={title} style={styles.image} />
        {isEnriching ? (
          <Animated.View
            style={[
              styles.shimmerOverlay,
              { backgroundColor: colors.primary, opacity: shimmerOpacity },
            ]}
          />
        ) : null}
      </View>
    </Pressable>
  );
}

export default React.memo(MediaCoverCard);
