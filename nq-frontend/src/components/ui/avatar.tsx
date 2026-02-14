import React from 'react';
import { View, Image, Text, StyleSheet, ViewStyle } from 'react-native';
import { cn } from './utils';
import { radii } from './tokens';
import { useTheme } from './ThemeProvider';

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

export function Avatar({ src, alt, fallback = 'U', size = 40, style, children }: AvatarProps) {
  const { colors } = useTheme();
  const styles = StyleSheet.create({
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
      fontSize: 16,
      fontWeight: 'bold',
    },
  });

  const avatarStyle = cn([styles.base, { width: size, height: size }, style]);

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
  const styles = StyleSheet.create({ image: { width: '100%', height: '100%' } });
  return <Image source={{ uri: src }} style={styles.image} accessibilityLabel={alt} />;
}

export function AvatarFallback({ children }: AvatarFallbackProps) {
  const styles = StyleSheet.create({ fallback: { color: '#666', fontSize: 16, fontWeight: 'bold' } as any });
  return <Text style={styles.fallback}>{children}</Text>;
}
