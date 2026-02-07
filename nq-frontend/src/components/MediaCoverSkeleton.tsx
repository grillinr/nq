import React from "react";
import {
  Animated,
  Easing,
  StyleProp,
  StyleSheet,
  View,
  ViewStyle,
} from "react-native";
import { LinearGradient } from "expo-linear-gradient";
import { radii } from "./ui/tokens";
import { useTheme } from "./ui/ThemeProvider";

const SHIMMER_DURATION = 1200;

interface MediaCoverSkeletonProps {
  style?: StyleProp<ViewStyle>;
  aspectRatio?: number;
}

function MediaCoverSkeleton({ style, aspectRatio = 2 / 3 }: MediaCoverSkeletonProps) {
  const { colors } = useTheme();
  const translateX = React.useRef(new Animated.Value(0)).current;

  React.useEffect(() => {
    translateX.setValue(-120);
    const animation = Animated.loop(
      Animated.timing(translateX, {
        toValue: 120,
        duration: SHIMMER_DURATION,
        easing: Easing.linear,
        useNativeDriver: true,
      })
    );
    animation.start();
    return () => animation.stop();
  }, [translateX]);

  const styles = React.useMemo(() => createStyles(colors), [colors]);

  return (
    <View style={[styles.card, { aspectRatio }, style]}>
      <Animated.View style={[styles.shimmer, { transform: [{ translateX }] }]}>
        <LinearGradient
          colors={[
            "rgba(255,255,255,0)",
            "rgba(255,255,255,0.2)",
            "rgba(255,255,255,0)",
          ]}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 0 }}
          style={StyleSheet.absoluteFill}
        />
      </Animated.View>
    </View>
  );
}

export default React.memo(MediaCoverSkeleton);

const createStyles = (colors: ReturnType<typeof useTheme>["colors"]) =>
  StyleSheet.create({
    card: {
      width: "100%",
      borderRadius: radii.lg,
      overflow: "hidden",
      backgroundColor: colors["input-background"],
    },
    shimmer: {
      position: "absolute",
      top: 0,
      bottom: 0,
      width: 160,
    },
  });
