import React, { useState } from 'react';
import { View, ImageStyle } from 'react-native';
import { Image } from 'expo-image';
import { useTheme } from './ThemeProvider';

const ERROR_IMG_SRC =
  'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODgiIGhlaWdodD0iODgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgc3Ryb2tlPSIjMDAwIiBzdHJva2UtbGluZWpvaW49InJvdW5kIiBvcGFjaXR5PSIuMyIgZmlsbD0ibm9uZSIgc3Ryb2tlLXdpZHRoPSIzLjciPjxyZWN0IHg9IjE2IiB5PSIxNiIgd2lkdGg9IjU2IiBoZWlnaHQ9IjU2IiByeD0iNiIvPjxwYXRoIGQ9Im0xNiA1OCAxNi0xOCAzMiAzMiIvPjxjaXJjbGUgY3g9IjUzIiBjeT0iMzUiIHI9IjciLz48L3N2Zz4KCg==';

interface ImageWithFallbackProps {
  src: string;
  alt?: string;
  style?: ImageStyle;
  contentFit?: 'cover' | 'contain' | 'fill' | 'scale-down';
}

function ImageWithFallback({
  src,
  alt,
  style,
  contentFit = 'cover',
  ...props
}: ImageWithFallbackProps) {
  const [didError, setDidError] = useState(false);
  const { colors } = useTheme();

  const handleError = () => {
    setDidError(true);
  };

  if (didError) {
    return (
      <View
        style={[
          {
            backgroundColor: colors.inputBackground,
            alignItems: 'center',
            justifyContent: 'center',
          },
          style,
        ]}
      >
        <Image source={{ uri: ERROR_IMG_SRC }} style={{ width: 88, height: 88 }} />
      </View>
    );
  }

  return (
    <Image
      source={{ uri: src }}
      style={style}
      contentFit={contentFit}
      transition={120}
      onError={handleError}
      accessibilityLabel={alt}
      {...props}
    />
  );
}

export default ImageWithFallback;
