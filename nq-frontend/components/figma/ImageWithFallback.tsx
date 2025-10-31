import React, { useState } from 'react';
import { Image, View, StyleSheet, ImageStyle } from 'react-native';

const ERROR_IMG_SRC = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODgiIGhlaWdodD0iODgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgc3Ryb2tlPSIjMDAwIiBzdHJva2UtbGluZWpvaW49InJvdW5kIiBvcGFjaXR5PSIuMyIgZmlsbD0ibm9uZSIgc3Ryb2tlLXdpZHRoPSIzLjciPjxyZWN0IHg9IjE2IiB5PSIxNiIgd2lkdGg9IjU2IiBoZWlnaHQ9IjU2IiByeD0iNiIvPjxwYXRoIGQ9Im0xNiA1OCAxNi0xOCAzMiAzMiIvPjxjaXJjbGUgY3g9IjUzIiBjeT0iMzUiIHI9IjciLz48L3N2Zz4KCg==';

interface ImageWithFallbackProps {
  src: string;
  alt?: string;
  style?: ImageStyle;
  // Add other Image props as needed
}

function ImageWithFallback({ src, alt, style, ...props }: ImageWithFallbackProps) {
  const [didError, setDidError] = useState(false);

  const handleError = () => {
    setDidError(true);
  };

  if (didError) {
    return (
      <View style={[styles.fallback, style]}>
        <Image source={{ uri: ERROR_IMG_SRC }} style={styles.errorImage} />
      </View>
    );
  }

  return (
    <Image
      source={{ uri: src }}
      style={style}
      resizeMode="cover"
      onError={handleError}
      accessibilityLabel={alt}
      {...props}
    />
  );
}

export default ImageWithFallback;

const styles = StyleSheet.create({
  fallback: {
    backgroundColor: '#f3f4f6', // gray-100
    alignItems: 'center',
    justifyContent: 'center',
  },
  errorImage: {
    width: 88,
    height: 88,
  },
});