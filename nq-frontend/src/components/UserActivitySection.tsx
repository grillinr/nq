import React, { useMemo, useState } from 'react';
import { View, Text, TextInput, StyleSheet, Alert } from 'react-native';
import { StarRating } from './ui/star-rating';
import { CharacterCounter } from './ui/character-counter';
import { StatusPicker, ActivityStatusId } from './ui/status-picker';
import { Button } from './ui/button';
import { useTheme } from './ui/theme-provider';
import { updateActivity } from '../lib/updateActivity';
import {
  fontSize,
  fontWeights,
  lineHeight,
  radii,
  spacing,
  ColorPalette,
} from './ui/tokens';

interface UserActivity {
  id: string;
  rating?: number | null;
  review?: string | null;
  status: {
    id: number;
    name: string;
  };
}

interface UserActivitySectionProps {
  activity?: UserActivity | null;
  mediaId: string;
  mediaTitle: string;
  onUpdate: () => void;
}

function createStyles(colors: ColorPalette) {
  return StyleSheet.create({
    container: {
      padding: spacing[4],
      borderRadius: radii.xl,
      marginVertical: spacing[4],
    },
    header: {
      flexDirection: 'row',
      justifyContent: 'space-between',
      alignItems: 'center',
      marginBottom: spacing[4],
    },
    title: {
      fontSize: fontSize.lg,
      fontWeight: fontWeights.bold,
      color: colors.foreground,
    },
    section: {
      marginBottom: spacing[4],
    },
    label: {
      fontSize: fontSize.sm,
      fontWeight: fontWeights.semibold,
      marginBottom: spacing[2],
      color: colors.foreground,
    },
    labelRow: {
      flexDirection: 'row',
      justifyContent: 'space-between',
      alignItems: 'center',
      marginBottom: spacing[2],
    },
    textInput: {
      borderWidth: 1,
      borderRadius: radii.lg,
      padding: spacing[3],
      fontSize: fontSize.base,
      minHeight: 100,
      color: colors.foreground,
      backgroundColor: colors.inputBackground,
      borderColor: colors.border,
    },
    buttonRow: {
      flexDirection: 'row',
      gap: spacing[3],
      marginTop: spacing[2],
    },
    button: {
      flex: 1,
    },
    addButton: {
      marginTop: spacing[3],
    },
    reviewText: {
      fontSize: fontSize.base,
      lineHeight: lineHeight.lg,
      color: colors.foreground,
    },
    statusBadge: {
      marginTop: spacing[2],
    },
    statusText: {
      fontSize: fontSize.sm,
      fontStyle: 'italic',
      color: colors.mutedForeground,
    },
    emptyText: {
      fontSize: fontSize.sm,
      textAlign: 'center',
      color: colors.mutedForeground,
    },
  });
}

export function UserActivitySection({
  activity,
  mediaId,
  mediaTitle,
  onUpdate,
}: UserActivitySectionProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const [isEditing, setIsEditing] = useState(false);
  const [rating, setRating] = useState(activity?.rating || 0);
  const [review, setReview] = useState(activity?.review || '');
  const [status, setStatus] = useState<ActivityStatusId>(
    (activity?.status.id || 1) as ActivityStatusId
  );
  const [loading, setLoading] = useState(false);

  if (!activity) {
    return (
      <View style={[styles.container, { backgroundColor: colors.background }]}>
        <Text style={styles.emptyText}>
          Track this item to add a rating and review
        </Text>
      </View>
    );
  }

  const handleSave = async () => {
    // Check for illogical status
    if (review.trim() && status !== 3) {
      Alert.alert('Warning', 'You have a review but the status is not "Completed". Are you sure?', [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Save Anyway', onPress: () => performSave() },
      ]);
      return;
    }

    await performSave();
  };

  const performSave = async () => {
    try {
      setLoading(true);

      await updateActivity(activity.id, {
        rating: rating || undefined,
        review: review.trim() || undefined,
        statusId: status,
      });

      setIsEditing(false);
      onUpdate();
    } catch (error) {
      console.error('Failed to update activity:', error);
      Alert.alert('Error', 'Failed to update rating and review');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setRating(activity?.rating || 0);
    setReview(activity?.review || '');
    setStatus((activity?.status.id || 1) as ActivityStatusId);
    setIsEditing(false);
  };

  const handleEdit = () => {
    setRating(activity?.rating || 0);
    setReview(activity?.review || '');
    setStatus((activity?.status.id || 1) as ActivityStatusId);
    setIsEditing(true);
  };

  if (isEditing) {
    return (
      <View style={[styles.container, { backgroundColor: colors.card }]}>
        <Text style={styles.title}>Rate & Review</Text>

        <View style={styles.section}>
          <Text style={styles.label}>Rating</Text>
          <StarRating value={rating} onChange={setRating} showValue size="lg" />
        </View>

        <View style={styles.section}>
          <View style={styles.labelRow}>
            <Text style={styles.label}>Review</Text>
            <CharacterCounter current={review.length} max={140} />
          </View>
          <TextInput
            style={styles.textInput}
            value={review}
            onChangeText={setReview}
            placeholder="Share your thoughts (optional)"
            placeholderTextColor={colors.mutedForeground}
            multiline
            maxLength={140}
            numberOfLines={4}
            textAlignVertical="top"
          />
        </View>

        <View style={styles.section}>
          <Text style={styles.label}>Status</Text>
          <StatusPicker value={status} onChange={setStatus} />
        </View>

        <View style={styles.buttonRow}>
          <Button variant="outline" onPress={handleCancel} style={styles.button} disabled={loading}>
            Cancel
          </Button>
          <Button
            onPress={handleSave}
            style={styles.button}
            disabled={loading || review.length > 140}
          >
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </View>
      </View>
    );
  }

  // Display mode
  const hasRatingOrReview = activity.rating || activity.review;

  return (
    <View
      style={[styles.container, { backgroundColor: colors.background, borderColor: colors.border }]}
    >
      <View style={styles.header}>
        <Text style={styles.title}>Your Rating & Review</Text>
        <Button variant="ghost" onPress={handleEdit} size="sm">
          Edit
        </Button>
      </View>

      {hasRatingOrReview ? (
        <>
          {activity.rating && (
            <View style={styles.section}>
              <StarRating value={activity.rating} readonly showValue size="lg" />
            </View>
          )}

          {activity.review && (
            <View style={styles.section}>
              <Text style={styles.reviewText}>
                {activity.review}
              </Text>
            </View>
          )}

          <View style={styles.statusBadge}>
            <Text style={styles.statusText}>
              Status: {activity.status.name}
            </Text>
          </View>
        </>
      ) : (
        <View style={styles.section}>
          <Text style={styles.emptyText}>
            No rating or review yet
          </Text>
          <Button onPress={handleEdit} style={styles.addButton}>
            Add Rating & Review
          </Button>
        </View>
      )}
    </View>
  );
}
