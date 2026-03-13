import React, { useState } from 'react';
import { router, useLocalSearchParams } from 'expo-router';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Pressable,
  ActivityIndicator,
  FlatList,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import MediaCoverCard from '../../src/components/MediaCoverCard';
import ImageWithFallback from '../../src/components/ui/image-with-fallback';
import Badge from '../../src/components/ui/badge';
import { useTheme } from '../../src/components/ui/theme-provider';
import { fontSize, fontWeights, lineHeight, radii, spacing } from '../../src/components/ui/tokens';
import { useMediaDetails } from '../../src/hooks/useMediaDetails';
import { UserActivitySection } from '../../src/components/UserActivitySection';
import { TrackItemModal } from '../../src/components/TrackItemModal';
import { ActivityStatusId } from '../../src/components/ui/status-picker';
import { createActivity } from '../../src/lib/createActivity';

const COVER_RATIO = 2 / 3;

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
      paddingBottom: spacing[8],
    },
    loadingContainer: {
      flex: 1,
      backgroundColor: colors.background,
      padding: spacing[4],
    },
    loadingCenter: {
      flex: 1,
      alignItems: 'center',
      justifyContent: 'center',
    },
    backButton: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[2],
      marginBottom: spacing[4],
    },
    backText: {
      color: colors.foreground,
      fontSize: fontSize.base,
    },
    hero: {
      flexDirection: 'row',
      gap: spacing[4],
      marginBottom: spacing[6],
    },
    coverWrap: {
      width: 140,
      aspectRatio: COVER_RATIO,
      borderRadius: radii.lg,
      overflow: 'hidden',
      backgroundColor: colors.inputBackground,
    },
    cover: {
      width: '100%',
      height: '100%',
    },
    heroText: {
      flex: 1,
    },
    title: {
      fontSize: fontSize['2xl'],
      fontWeight: fontWeights.semibold,
      color: colors.foreground,
      marginBottom: spacing[2],
    },
    metaRow: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
      marginBottom: spacing[2],
    },
    metaText: {
      color: colors.mutedForeground,
      fontSize: fontSize.sm,
    },
    ratingRow: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[1],
      marginBottom: spacing[2],
    },
    ratingText: {
      color: colors.foreground,
      fontSize: fontSize.base,
      fontWeight: fontWeights.medium,
    },
    metaDetail: {
      color: colors.mutedForeground,
      fontSize: fontSize.sm,
      marginBottom: spacing[2],
    },
    genreRow: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
    },
    section: {
      marginBottom: spacing[6],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: fontWeights.semibold,
      color: colors.foreground,
      marginBottom: spacing[3],
    },
    bodyText: {
      fontSize: fontSize.base,
      color: colors.mutedForeground,
      lineHeight: lineHeight.md,
    },
    chipRow: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
    },
    chip: {
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[1],
      borderRadius: radii.md,
      backgroundColor: colors.secondary,
    },
    chipText: {
      fontSize: fontSize.sm,
      color: colors.secondaryForeground,
    },
    relatedList: {
      gap: spacing[3],
    },
    relatedItem: {
      width: 120,
    },
    relatedTitle: {
      marginTop: spacing[2],
      color: colors.foreground,
      fontSize: fontSize.sm,
      fontWeight: fontWeights.medium,
    },
    relatedMeta: {
      color: colors.mutedForeground,
      fontSize: fontSize.xs,
      marginTop: 2,
    },
    emptyText: {
      color: colors.mutedForeground,
      fontSize: fontSize.base,
    },
    trackButton: {
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'center',
      gap: spacing[2],
      paddingVertical: spacing[3],
      paddingHorizontal: spacing[4],
      borderRadius: radii.md,
    },
    trackButtonText: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.semibold,
    },
  });

export default function MediaDetailsPage() {
  const { id } = useLocalSearchParams<{ id?: string | string[] }>();
  const { colors } = useTheme();
  const mediaId = Array.isArray(id) ? id[0] : id;
  const { details, loading, error, refetch } = useMediaDetails(mediaId);
  const styles = React.useMemo(() => createStyles(colors), [colors]);
  const [trackModalVisible, setTrackModalVisible] = useState(false);
  const [trackingItem, setTrackingItem] = useState(false);
  const isMountedRef = React.useRef(true);

  React.useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  React.useEffect(() => {
    if (error) {
      console.error('Failed to load media details:', error);
    }
  }, [error]);

  const handleTrackItem = async (statusId: ActivityStatusId) => {
    if (!mediaId) return;

    try {
      setTrackingItem(true);
      await createActivity({
        mediaId,
        statusId,
      });

      if (!isMountedRef.current) return;

      setTrackModalVisible(false);

      // Refetch to get updated myActivity
      if (refetch) {
        try {
          await refetch();
        } catch (refetchError: any) {
          // Ignore abort errors when component unmounts
          if (refetchError.name !== 'AbortError') {
            console.error('Failed to refetch:', refetchError);
          }
        }
      }
    } catch (err: any) {
      if (!isMountedRef.current) return;
      console.error('Failed to track item:', err);
      // Only show error if not aborted
      if (err.name !== 'AbortError') {
        Alert.alert('Error', 'Failed to track this item');
      }
    } finally {
      if (isMountedRef.current) {
        setTrackingItem(false);
      }
    }
  };

  const handleActivityUpdate = async () => {
    // Refetch to get updated myActivity
    if (refetch) {
      try {
        await refetch();
      } catch (refetchError: any) {
        // Ignore abort errors when component unmounts
        if (refetchError.name !== 'AbortError') {
          console.error('Failed to refetch:', refetchError);
        }
      }
    }
  };

  const backButton = (
    <Pressable onPress={() => router.back()} style={styles.backButton} accessibilityRole="button">
      <Ionicons name="arrow-back" size={18} color={colors.foreground} />
      <Text style={styles.backText}>Back</Text>
    </Pressable>
  );

  if (loading && !details) {
    return (
      <View style={styles.loadingContainer}>
        {backButton}
        <View style={styles.loadingCenter}>
          <ActivityIndicator color={colors.primary} />
        </View>
      </View>
    );
  }

  if (!details) {
    return (
      <View style={styles.loadingContainer}>
        {backButton}
        <View style={styles.loadingCenter}>
          <Text style={styles.emptyText}>
            {error ? 'We hit an error loading this title.' : "We couldn't find that title."}
          </Text>
        </View>
      </View>
    );
  }

  const actors = details.actors.slice(0, 10);
  const { related } = details;

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      {backButton}

      <View style={styles.hero}>
        <View style={styles.coverWrap}>
          <ImageWithFallback src={details.image} alt={details.title} style={styles.cover} />
        </View>
        <View style={styles.heroText}>
          <Text style={styles.title}>{details.title}</Text>
          <View style={styles.metaRow}>
            <Text style={styles.metaText}>{details.year}</Text>
            {details.duration ? <Text style={styles.metaText}>· {details.duration}</Text> : null}
            <Text style={styles.metaText}>· {details.type.toUpperCase()}</Text>
          </View>
          <View style={styles.ratingRow}>
            <Ionicons name="star" size={14} color={colors.chart4} />
            <Text style={styles.ratingText}>{details.rating.toFixed(1)}</Text>
          </View>
          {details.metaLabel && details.metaValue ? (
            <Text style={styles.metaDetail}>
              {details.metaLabel}: {details.metaValue}
            </Text>
          ) : null}
          {details.genre.length > 0 ? (
            <View style={styles.genreRow}>
              {details.genre.slice(0, 4).map(item => (
                <Badge key={item} variant="secondary">
                  {item}
                </Badge>
              ))}
            </View>
          ) : null}
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Overview</Text>
        <Text style={styles.bodyText}>
          {details.description || 'No description available yet.'}
        </Text>
      </View>

      {!details.myActivity ? (
        <View style={styles.section}>
          <Pressable
            style={[styles.trackButton, { backgroundColor: colors.primary }]}
            onPress={() => setTrackModalVisible(true)}
          >
            <Ionicons name="add-circle-outline" size={20} color={colors.primaryForeground} />
            <Text style={[styles.trackButtonText, { color: colors.primaryForeground }]}>
              Track this item
            </Text>
          </Pressable>
        </View>
      ) : (
        <UserActivitySection activity={details.myActivity} onUpdate={handleActivityUpdate} />
      )}

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Actors & Creators</Text>
        {actors.length === 0 && details.creators.length === 0 ? (
          <Text style={styles.bodyText}>No cast or creator info yet.</Text>
        ) : (
          <View style={styles.chipRow}>
            {actors.map(actor => (
              <View key={actor.id} style={styles.chip}>
                <Text style={styles.chipText}>{actor.name}</Text>
              </View>
            ))}
            {details.creators.slice(0, 6).map(creator => (
              <View key={creator.id} style={styles.chip}>
                <Text style={styles.chipText}>{creator.name}</Text>
              </View>
            ))}
          </View>
        )}
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Related</Text>
        {related.length === 0 ? (
          <Text style={styles.bodyText}>No related titles yet.</Text>
        ) : (
          <FlatList
            data={related}
            horizontal
            showsHorizontalScrollIndicator={false}
            keyExtractor={item => String(item.id)}
            contentContainerStyle={styles.relatedList}
            renderItem={({ item }) => (
              <View style={styles.relatedItem}>
                <MediaCoverCard
                  title={item.title}
                  image={item.image}
                  aspectRatio={COVER_RATIO}
                  onPress={() => router.push({ pathname: '/media/[id]', params: { id: item.id } })}
                />
                <Text style={styles.relatedTitle} numberOfLines={1}>
                  {item.title}
                </Text>
                <Text style={styles.relatedMeta} numberOfLines={1}>
                  {item.year} · {item.type.toUpperCase()}
                </Text>
              </View>
            )}
          />
        )}
      </View>

      <TrackItemModal
        visible={trackModalVisible}
        onClose={() => setTrackModalVisible(false)}
        onConfirm={handleTrackItem}
        mediaTitle={details.title}
        loading={trackingItem}
      />
    </ScrollView>
  );
}
