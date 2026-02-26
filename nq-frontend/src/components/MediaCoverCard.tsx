import React from 'react';
import { Pressable, StyleSheet, View, ViewStyle } from 'react-native';
import ImageWithFallback from './ui/ImageWithFallback';
import { useTheme } from './ui/ThemeProvider';
import { radii } from './ui/tokens';

interface MediaCoverCardProps {
  title: string;
  image: string;
  onPress?: () => void;
  aspectRatio?: number;
  style?: ViewStyle;
}

function MediaCoverCard({
  title,
  image,
  onPress,
  aspectRatio = 2 / 3,
  style,
}: MediaCoverCardProps) {
  const { colors } = useTheme();
  const styles = React.useMemo(() => createStyles(colors), [colors]);

  return (
    <Pressable
      onPress={onPress}
      style={[styles.pressable, style]}
      accessibilityRole="button"
      accessibilityLabel={title}
    >
      <View style={[styles.imageContainer, { aspectRatio }]}>
        <ImageWithFallback src={image} alt={title} style={styles.image} />
      </View>
    </Pressable>
  );
}

export default React.memo(MediaCoverCard);

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
  });
