import React, { useMemo } from 'react';
import { View, Image, Text, StyleSheet, ViewStyle } from 'react-native';
import { flattenStyles } from './utils';
import { fontSize, fontWeights, radii, ColorPalette } from './tokens';
import { useTheme } from './theme-provider';

interface AvatarProps {
  src?: string;
  alt?: string;
  fallback?: string;
  size?: number;
  style?: ViewStyle;
  children?: React.ReactNode;
}

interface AvatarImageProps {
  src: string;
  alt?: string;
}

interface AvatarFallbackProps {
  children: string;
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    base: {
      backgroundColor: colors.muted,
      alignItems: 'center',
      justifyContent: 'center',
      borderRadius: radii.md,
      overflow: 'hidden',
    },
    image: {
      width: '100%',
      height: '100%',
    },
    fallback: {
      color: colors.mutedForeground,
      fontSize: fontSize.base,
      fontWeight: fontWeights.bold,
    },
  });
}

export function Avatar({ src, alt, fallback = 'U', size = 40, style, children }: AvatarProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);

  const avatarStyle = flattenStyles([styles.base, { width: size, height: size }, style]);

  if (children) {
    return <View style={avatarStyle}>{children}</View>;
  }

  return (
    <View style={avatarStyle}>
      {src ? (
        <Image source={{ uri: src }} style={styles.image} accessibilityLabel={alt} />
      ) : (
        <Text style={styles.fallback}>{fallback}</Text>
      )}
    </View>
  );
}

export function AvatarImage({ src, alt }: AvatarImageProps) {
  return <Image source={{ uri: src }} style={avatarImageStyle} accessibilityLabel={alt} />;
}

const avatarImageStyle = StyleSheet.create({ image: { width: '100%', height: '100%' } }).image;

export function AvatarFallback({ children }: AvatarFallbackProps) {
  const { colors } = useTheme();
  const styles = useMemo(
    () =>
      StyleSheet.create({
        fallback: {
          color: colors.mutedForeground,
          fontSize: fontSize.base,
          fontWeight: fontWeights.bold,
        } as any,
      }),
    [colors]
  );
  return <Text style={styles.fallback}>{children}</Text>;
}
