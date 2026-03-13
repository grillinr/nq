import React, { useMemo, useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Alert,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useApolloClient } from '@apollo/client/react';
import { router } from 'expo-router';
import { Button } from '../../src/components/ui/button';
import Input from '../../src/components/ui/input';
import Card from '../../src/components/ui/card';
import PageHeader from '../../src/components/PageHeader';
import {
  fontSize,
  spacing,
  radii,
  fontWeights,
  layout,
  zIndex,
} from '../../src/components/ui/tokens';
import { useTheme } from '../../src/components/ui/theme-provider';
import { useScrollHeader } from '../../src/hooks/useScrollHeader';
import { Media } from '../../src/types';
import { GET_HOME_MEDIA_QUERY, ME_ACTIVITIES_QUERY } from '../../src/lib/graphql';
import { StarRating } from '../../src/components/ui/star-rating';
import { CharacterCounter } from '../../src/components/ui/character-counter';
import { StatusPicker, ActivityStatusId } from '../../src/components/ui/status-picker';
import { createMedia } from '../../src/lib/createMedia';
import { createActivity } from '../../src/lib/createActivity';
import { useAuth } from '../../src/lib/AuthContext';
import { MediaAutocomplete, MediaSuggestion } from '../../src/components/MediaAutocomplete';

const typeOptions = [
  { label: 'Movie', value: 'movie' as const, icon: 'film-outline' as const },
  { label: 'TV Show', value: 'tv' as const, icon: 'tv-outline' as const },
  { label: 'Book', value: 'book' as const, icon: 'book-outline' as const },
  {
    label: 'Music',
    value: 'music' as const,
    icon: 'musical-notes-outline' as const,
  },
  {
    label: 'Game',
    value: 'game' as const,
    icon: 'game-controller-outline' as const,
  },
];

type MediaType = 'movie' | 'tv' | 'book' | 'music' | 'game';

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      maxWidth: 600,
      alignSelf: 'center',
      width: '100%',
      paddingHorizontal: spacing[4],
      paddingBottom: spacing[6],
      paddingTop: spacing[4],
    },
    card: {
      padding: spacing[6],
    },
    form: {},
    field: {
      marginBottom: spacing[6],
    },
    helperText: {
      color: colors.mutedForeground,
      marginTop: spacing[2],
      fontSize: fontSize.sm,
    },
    label: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.medium,
      color: colors.foreground,
      marginBottom: spacing[2],
    },
    typeOptions: {
      flexDirection: 'row',
      flexWrap: 'wrap',
      gap: spacing[2],
    },
    typeOption: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[2],
      padding: spacing[3],
      borderRadius: radii.lg,
      borderWidth: 1,
      borderColor: colors.border,
      backgroundColor: colors.background,
    },
    typeOptionSelected: {
      backgroundColor: colors.primary,
      borderColor: colors.primary,
    },
    typeText: {
      fontSize: fontSize.sm,
      color: colors.foreground,
    },
    typeTextSelected: {
      color: colors.primaryForeground,
    },
    submitButton: {
      marginTop: spacing[6],
    },
    submitText: {
      color: colors.primaryForeground,
      marginLeft: spacing[2],
    },
    sectionToggle: {
      flexDirection: 'row',
      justifyContent: 'space-between',
      alignItems: 'center',
    },
    ratingSection: {
      marginTop: spacing[4],
      paddingTop: spacing[4],
      borderTopWidth: 1,
      borderTopColor: colors.border,
    },
    labelRow: {
      flexDirection: 'row',
      justifyContent: 'space-between',
      alignItems: 'center',
      marginBottom: spacing[2],
    },
    headerSpacer: {
      paddingTop: layout.headerHeight,
    },
    titleInputContainer: {
      position: 'relative',
      zIndex: zIndex.modal,
    },
  });

export default function AddTabPage() {
  const { colors } = useTheme();
  const apolloClient = useApolloClient();
  const { hasToken } = useAuth();
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);

  const [isAddingMedia, setIsAddingMedia] = useState(false);

  const [title, setTitle] = useState('');
  const [type, setType] = useState<MediaType | null>(null);
  const [year, setYear] = useState('');
  const [selectedExternalId, setSelectedExternalId] = useState<string | undefined>();
  const [selectedIsbn, setSelectedIsbn] = useState<string | undefined>();
  const [suppressAutocomplete, setSuppressAutocomplete] = useState(false);
  const [titleInputHeight, setTitleInputHeight] = useState(0);

  // Rating, Review, and Status fields
  const [rating, setRating] = useState(0);
  const [review, setReview] = useState('');
  const [status, setStatus] = useState<ActivityStatusId>(1);
  const [showRatingSection, setShowRatingSection] = useState(false);

  const handleAddMedia = async (
    newMedia: Omit<Media, 'id'>,
    activityData?: { rating?: number; review?: string; statusId: number }
  ) => {
    setIsAddingMedia(true);
    try {
      const result = await createMedia(newMedia);
      if (result?.id) {
        await createActivity({
          mediaId: result.id,
          statusId: activityData?.statusId || 1,
          rating: activityData?.rating,
          review: activityData?.review,
        });
        router.replace({
          pathname: '/history',
          params: { addedMediaId: result.id },
        });
      }
      const queries: Promise<unknown>[] = [
        apolloClient.query({ query: GET_HOME_MEDIA_QUERY, fetchPolicy: 'network-only' }),
      ];
      if (hasToken) {
        queries.push(
          apolloClient.query({ query: ME_ACTIVITIES_QUERY, fetchPolicy: 'network-only' })
        );
      }
      Promise.all(queries).catch(() => undefined);
    } catch (error) {
      console.error('Failed to add media:', error);
    } finally {
      setIsAddingMedia(false);
    }
  };

  const handleSubmit = () => {
    if (isAddingMedia) return;
    if (!title.trim()) {
      Alert.alert('Error', 'Please enter a title');
      return;
    }
    if (!type) {
      Alert.alert('Error', 'Please select a media type');
      return;
    }

    // Validate review length
    if (review.length > 140) {
      Alert.alert('Error', 'Review must be 140 characters or less');
      return;
    }

    const newMedia: Omit<Media, 'id'> = {
      title: title.trim(),
      type,
      description: '', // Backend will enrich
      year: parseInt(year, 10) || new Date().getFullYear(),
      rating: 0, // Backend will enrich
      duration: undefined, // Backend will enrich
      image: '', // Backend will enrich
      genre: [], // Backend will enrich
      externalId: selectedExternalId,
      isbn: selectedIsbn,
    };

    // Prepare activity data if rating or review is provided
    let activityData: { rating?: number; review?: string; statusId: number } | undefined;
    if (rating > 0 || review.trim()) {
      // Auto-set to Completed if review exists
      const finalStatus = review.trim() ? 3 : status;
      activityData = {
        rating: rating > 0 ? rating : undefined,
        review: review.trim() || undefined,
        statusId: finalStatus,
      };
    } else {
      // No rating or review, just use selected status
      activityData = {
        statusId: status,
      };
    }

    handleAddMedia(newMedia, activityData);

    // Reset form
    setTitle('');
    setYear('');
    setSelectedExternalId(undefined);
    setSelectedIsbn(undefined);
    setSuppressAutocomplete(false);
    setRating(0);
    setReview('');
    setStatus(1);
    setShowRatingSection(false);
  };

  const handleTypeChange = (value: string) => {
    setType(value as MediaType);
    setSelectedExternalId(undefined);
    setSelectedIsbn(undefined);
    setSuppressAutocomplete(false);
  };

  const typeLabel = useMemo(() => typeOptions.find(option => option.value === type)?.label, [type]);
  const canType = Boolean(type);
  const showMusicNotice = type === 'music';

  const handleSuggestionSelect = (item: MediaSuggestion) => {
    setTitle(item.title);
    if (item.year) {
      setYear(String(item.year));
    }
    if (type === 'book') {
      setSelectedIsbn(item.externalId || undefined);
      setSelectedExternalId(undefined);
    } else {
      setSelectedExternalId(item.externalId || undefined);
      setSelectedIsbn(undefined);
    }
    setSuppressAutocomplete(true);
  };

  const styles = React.useMemo(() => createStyles(colors), [colors]);

  return (
    <View style={{ flex: 1 }}>
      <PageHeader
        title="Add New Media"
        subtitle="Enter the title and year. We'll fetch the rest of the details for you."
        visible={isHeaderVisible}
      />
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.content}
        onScroll={handleHeaderScroll}
        scrollEventThrottle={16}
      >
        <View style={styles.headerSpacer} />

        <Card style={styles.card}>
          <View style={styles.form}>
            <View style={styles.field}>
              <Text style={styles.label}>Type *</Text>
              <View style={styles.typeOptions}>
                {typeOptions.map(option => (
                  <TouchableOpacity
                    key={option.value}
                    style={[styles.typeOption, type === option.value && styles.typeOptionSelected]}
                    onPress={() => handleTypeChange(option.value)}
                  >
                    <Ionicons
                      name={option.icon}
                      size={20}
                      color={type === option.value ? colors.primaryForeground : colors.foreground}
                    />
                    <Text
                      style={[styles.typeText, type === option.value && styles.typeTextSelected]}
                    >
                      {option.label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
              {!type ? <Text style={styles.helperText}>Select media type...</Text> : null}
            </View>

            {canType ? (
              <View style={styles.field}>
                <Text style={styles.label}>Title *</Text>
                <View style={styles.titleInputContainer}>
                  <View onLayout={e => setTitleInputHeight(e.nativeEvent.layout.height)}>
                    <Input
                      value={title}
                      onChangeText={value => {
                        setTitle(value);
                        setSelectedExternalId(undefined);
                        setSelectedIsbn(undefined);
                        setSuppressAutocomplete(false);
                      }}
                      placeholder={`Enter ${typeLabel ?? 'media'} title`}
                    />
                  </View>
                  {showMusicNotice ? (
                    <Text style={styles.helperText}>Music autocomplete coming soon.</Text>
                  ) : null}
                  {!showMusicNotice && type ? (
                    <MediaAutocomplete
                      type={type}
                      query={title}
                      suppress={suppressAutocomplete}
                      onSelect={handleSuggestionSelect}
                      inputHeight={titleInputHeight}
                    />
                  ) : null}
                </View>
              </View>
            ) : null}

            {type ? (
              <View style={styles.field}>
                <Text style={styles.label}>Year</Text>
                <Input
                  value={year}
                  onChangeText={setYear}
                  placeholder="Release year"
                  keyboardType="numeric"
                />
              </View>
            ) : null}

            {/* Rating, Review, and Status Section */}
            {type ? (
              <View style={styles.field}>
                <TouchableOpacity
                  onPress={() => setShowRatingSection(!showRatingSection)}
                  style={styles.sectionToggle}
                >
                  <Text style={styles.label}>Rate & Review (Optional)</Text>
                  <Ionicons
                    name={showRatingSection ? 'chevron-up' : 'chevron-down'}
                    size={20}
                    color={colors.foreground}
                  />
                </TouchableOpacity>

                {showRatingSection && (
                  <View style={styles.ratingSection}>
                    <View style={styles.field}>
                      <Text style={styles.label}>Your Rating</Text>
                      <StarRating value={rating} onChange={setRating} showValue size="lg" />
                    </View>

                    <View style={styles.field}>
                      <View style={styles.labelRow}>
                        <Text style={styles.label}>Review</Text>
                        <CharacterCounter current={review.length} max={140} />
                      </View>
                      <Input
                        value={review}
                        onChangeText={setReview}
                        placeholder="Share your thoughts (optional)"
                        multiline
                        numberOfLines={4}
                        maxLength={140}
                        style={{ minHeight: 100 }}
                      />
                      <Text style={styles.helperText}>
                        {review.trim()
                          ? 'Adding a review will set status to Completed'
                          : 'Max 140 characters'}
                      </Text>
                    </View>

                    {!review.trim() && (
                      <View style={styles.field}>
                        <Text style={styles.label}>Status</Text>
                        <StatusPicker value={status} onChange={setStatus} />
                      </View>
                    )}
                  </View>
                )}
              </View>
            ) : null}

            {type ? (
              <Button onPress={handleSubmit} style={styles.submitButton} disabled={isAddingMedia}>
                {isAddingMedia ? (
                  <ActivityIndicator color={colors.primaryForeground} />
                ) : (
                  <>
                    <Ionicons name="add" size={20} color={colors.primaryForeground} />
                    <Text style={styles.submitText}>Add Media</Text>
                  </>
                )}
              </Button>
            ) : null}
          </View>
        </Card>
      </ScrollView>
    </View>
  );
}
