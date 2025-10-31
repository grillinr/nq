import React from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import Card from './ui/card';
import Badge from './ui/badge';
import ImageWithFallback from './figma/ImageWithFallback';
import { colors, fontSize, spacing } from './ui/tokens';

interface MediaCardProps {
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  onPress?: () => void;
}

function MediaCard({
  title,
  image,
  rating,
  genre,
  year,
  duration,
  description,
  onPress,
}: MediaCardProps) {
  return (
    <Pressable onPress={onPress} style={styles.pressable}>
      <Card style={styles.card}>
        <View style={styles.imageContainer}>
          <ImageWithFallback
            src={image}
            alt={title}
            style={styles.image}
          />
          <View style={styles.ratingBadge}>
            <Ionicons name="star" size={12} color="#fbbf24" />
            <Text style={styles.ratingText}>{rating.toFixed(1)}</Text>
          </View>
        </View>
        <View style={styles.content}>
          <Text style={styles.title} numberOfLines={1}>
            {title}
          </Text>
          <View style={styles.meta}>
            <View style={styles.metaItem}>
              <Ionicons name="calendar-outline" size={10} color={colors['muted-foreground']} />
              <Text style={styles.metaText}>{year}</Text>
            </View>
            {duration && (
              <>
                <Text style={styles.dot}>•</Text>
                <View style={styles.metaItem}>
                  <Ionicons name="time-outline" size={10} color={colors['muted-foreground']} />
                  <Text style={styles.metaText}>{duration}</Text>
                </View>
              </>
            )}
          </View>
          <View style={styles.genres}>
            {genre.slice(0, 3).map((g) => (
              <Badge key={g} variant="secondary" style={styles.genreBadge}>
                {g}
              </Badge>
            ))}
          </View>
          <Text style={styles.description} numberOfLines={2}>
            {description}
          </Text>
        </View>
      </Card>
    </Pressable>
  );
}

export default MediaCard;

const styles = StyleSheet.create({
  pressable: {
    marginBottom: spacing[4],
  },
  card: {
    overflow: 'hidden',
  },
  imageContainer: {
    aspectRatio: 2 / 3,
    backgroundColor: '#f3f4f6',
    overflow: 'hidden',
  },
  image: {
    width: '100%',
    height: '100%',
  },
  ratingBadge: {
    position: 'absolute',
    top: 8,
    right: 8,
    backgroundColor: 'rgba(0,0,0,0.8)',
    borderRadius: 4,
    paddingHorizontal: 6,
    paddingVertical: 2,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  ratingText: {
    color: 'white',
    fontSize: 12,
  },
  content: {
    padding: spacing[4],
  },
  title: {
    fontSize: fontSize.base,
    fontWeight: '500',
    marginBottom: spacing[2],
  },
  meta: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing[3],
    gap: 4,
  },
  metaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 2,
  },
  metaText: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
  },
  dot: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
  },
  genres: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 4,
    marginBottom: spacing[3],
  },
  genreBadge: {
    // fontSize handled by Badge component
  },
  description: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
    lineHeight: 18,
  },
});