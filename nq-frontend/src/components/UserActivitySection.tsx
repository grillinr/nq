import React, { useState } from 'react';
import { View, Text, TextInput, StyleSheet, Alert } from 'react-native';
import { StarRating } from './ui/StarRating';
import { CharacterCounter } from './ui/CharacterCounter';
import { StatusPicker, ActivityStatusId } from './ui/StatusPicker';
import { Button } from './ui/button';
import { useTheme } from './ui/ThemeProvider';
import { updateActivity } from '../../lib/updateActivity';

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

export function UserActivitySection({ 
  activity, 
  mediaId, 
  mediaTitle,
  onUpdate 
}: UserActivitySectionProps) {
  const { colors } = useTheme();
  const [isEditing, setIsEditing] = useState(false);
  const [rating, setRating] = useState(activity?.rating || 0);
  const [review, setReview] = useState(activity?.review || '');
  const [status, setStatus] = useState<ActivityStatusId>((activity?.status.id || 1) as ActivityStatusId);
  const [loading, setLoading] = useState(false);

  if (!activity) {
    return (
      <View style={[styles.container, { backgroundColor: colors.card }]}>
        <Text style={[styles.emptyText, { color: colors['muted-foreground'] }]}>
          Track this item to add a rating and review
        </Text>
      </View>
    );
  }

  const handleSave = async () => {
    // Check for illogical status
    if (review.trim() && status !== 3) {
      Alert.alert(
        'Warning',
        'You have a review but the status is not "Completed". Are you sure?',
        [
          { text: 'Cancel', style: 'cancel' },
          { text: 'Save Anyway', onPress: () => performSave() }
        ]
      );
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
        <Text style={[styles.title, { color: colors.foreground }]}>
          Rate & Review
        </Text>

        <View style={styles.section}>
          <Text style={[styles.label, { color: colors.foreground }]}>Rating</Text>
          <StarRating 
            value={rating} 
            onChange={setRating}
            showValue
            size="lg"
          />
        </View>

        <View style={styles.section}>
          <View style={styles.labelRow}>
            <Text style={[styles.label, { color: colors.foreground }]}>Review</Text>
            <CharacterCounter current={review.length} max={140} />
          </View>
          <TextInput
            style={[
              styles.textInput,
              { 
                color: colors.foreground,
                backgroundColor: colors['input-background'],
                borderColor: colors.border,
              }
            ]}
            value={review}
            onChangeText={setReview}
            placeholder="Share your thoughts (optional)"
            placeholderTextColor={colors['muted-foreground']}
            multiline
            maxLength={140}
            numberOfLines={4}
            textAlignVertical="top"
          />
        </View>

        <View style={styles.section}>
          <Text style={[styles.label, { color: colors.foreground }]}>Status</Text>
          <StatusPicker value={status} onChange={setStatus} />
        </View>

        <View style={styles.buttonRow}>
          <Button 
            variant="outline" 
            onPress={handleCancel}
            style={styles.button}
            disabled={loading}
          >
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
    <View style={[styles.container, { backgroundColor: colors.card }]}>
      <View style={styles.header}>
        <Text style={[styles.title, { color: colors.foreground }]}>
          Your Rating & Review
        </Text>
        <Button variant="ghost" onPress={handleEdit} size="sm">
          Edit
        </Button>
      </View>

      {hasRatingOrReview ? (
        <>
          {activity.rating && (
            <View style={styles.section}>
              <StarRating 
                value={activity.rating} 
                readonly
                showValue
                size="lg"
              />
            </View>
          )}

          {activity.review && (
            <View style={styles.section}>
              <Text style={[styles.reviewText, { color: colors.foreground }]}>
                {activity.review}
              </Text>
            </View>
          )}

          <View style={styles.statusBadge}>
            <Text style={[styles.statusText, { color: colors['muted-foreground'] }]}>
              Status: {activity.status.name}
            </Text>
          </View>
        </>
      ) : (
        <View style={styles.section}>
          <Text style={[styles.emptyText, { color: colors['muted-foreground'] }]}>
            No rating or review yet
          </Text>
          <Button onPress={handleEdit} style={{ marginTop: 12 }}>
            Add Rating & Review
          </Button>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: 16,
    borderRadius: 12,
    marginVertical: 16,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  section: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 8,
  },
  labelRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  textInput: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
    minHeight: 100,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 8,
  },
  button: {
    flex: 1,
  },
  reviewText: {
    fontSize: 16,
    lineHeight: 24,
  },
  statusBadge: {
    marginTop: 8,
  },
  statusText: {
    fontSize: 14,
    fontStyle: 'italic',
  },
  emptyText: {
    fontSize: 14,
    textAlign: 'center',
  },
});
