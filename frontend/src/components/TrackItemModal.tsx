import React, { useMemo, useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import Modal from 'react-native-modal';
import { Ionicons } from '@expo/vector-icons';
import * as Haptics from 'expo-haptics';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Button } from './ui/button';
import { useTheme } from './ui/theme-provider';
import { createShadows, fontSize, fontWeights, radii, spacing, ColorPalette } from './ui/tokens';

export type ActivityStatusId = 1 | 2 | 3;

interface TrackItemModalProps {
  visible: boolean;
  onClose: () => void;
  onConfirm: (statusId: ActivityStatusId) => void;
  mediaTitle: string;
  loading?: boolean;
}

const STATUS_OPTIONS: {
  id: ActivityStatusId;
  name: string;
  icon: keyof typeof Ionicons.glyphMap;
}[] = [
  { id: 1, name: 'Planned', icon: 'bookmark-outline' },
  { id: 2, name: 'In Progress', icon: 'play-outline' },
  { id: 3, name: 'Completed', icon: 'checkmark-circle-outline' },
];

function createStyles(colors: ColorPalette) {
  const shadows = createShadows(colors);
  return StyleSheet.create({
    sheet: {
      margin: 0,
      justifyContent: 'flex-end',
    },
    container: {
      backgroundColor: colors.card,
      borderTopLeftRadius: 24,
      borderTopRightRadius: 24,
      paddingHorizontal: spacing[5],
      paddingTop: spacing[3],
      ...shadows.modal,
    },
    handle: {
      width: 36,
      height: 4,
      borderRadius: 2,
      backgroundColor: colors.border,
      alignSelf: 'center',
      marginBottom: spacing[5],
    },
    title: {
      fontSize: fontSize.lg,
      fontWeight: fontWeights.semibold,
      color: colors.foreground,
      marginBottom: spacing[1],
    },
    subtitle: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      marginBottom: spacing[5],
    },
    statusRow: {
      flexDirection: 'row',
      gap: spacing[2],
      marginBottom: spacing[5],
    },
    statusOption: {
      flex: 1,
      alignItems: 'center',
      paddingVertical: spacing[3],
      paddingHorizontal: spacing[2],
      borderRadius: radii.lg,
      borderWidth: 1.5,
      gap: spacing[1],
    },
    statusOptionSelected: {
      backgroundColor: colors.primary,
      borderColor: colors.primary,
    },
    statusOptionUnselected: {
      backgroundColor: colors.muted,
      borderColor: colors.input,
    },
    statusLabel: {
      fontSize: fontSize.xs,
      fontWeight: fontWeights.medium,
      textAlign: 'center',
    },
    confirmButton: {
      marginBottom: spacing[2],
    },
    cancelButton: {
      marginBottom: spacing[2],
    },
  });
}

export function TrackItemModal({
  visible,
  onClose,
  onConfirm,
  mediaTitle,
  loading = false,
}: TrackItemModalProps) {
  const { colors } = useTheme();
  const insets = useSafeAreaInsets();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const [selectedStatus, setSelectedStatus] = useState<ActivityStatusId>(1);

  const handleStatusSelect = (id: ActivityStatusId) => {
    Haptics.selectionAsync();
    setSelectedStatus(id);
  };

  const handleConfirm = () => {
    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
    onConfirm(selectedStatus);
  };

  const handleClose = () => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    onClose();
  };

  return (
    <Modal
      isVisible={visible}
      onBackdropPress={handleClose}
      onSwipeComplete={handleClose}
      swipeDirection={['down']}
      style={styles.sheet}
      backdropOpacity={0.4}
      propagateSwipe
      useNativeDriverForBackdrop
    >
      <View style={[styles.container, { paddingBottom: Math.max(insets.bottom, spacing[5]) }]}>
        <View style={styles.handle} />

        <Text style={styles.title} numberOfLines={1}>
          Track &quot;{mediaTitle}&quot;
        </Text>
        <Text style={styles.subtitle}>What&apos;s your status?</Text>

        <View style={styles.statusRow}>
          {STATUS_OPTIONS.map(option => {
            const isSelected = selectedStatus === option.id;
            return (
              <TouchableOpacity
                key={option.id}
                style={[
                  styles.statusOption,
                  isSelected ? styles.statusOptionSelected : styles.statusOptionUnselected,
                ]}
                onPress={() => handleStatusSelect(option.id)}
                activeOpacity={0.7}
              >
                <Ionicons
                  name={option.icon}
                  size={22}
                  color={isSelected ? colors.primaryForeground : colors.mutedForeground}
                />
                <Text
                  style={[
                    styles.statusLabel,
                    { color: isSelected ? colors.primaryForeground : colors.mutedForeground },
                  ]}
                >
                  {option.name}
                </Text>
              </TouchableOpacity>
            );
          })}
        </View>

        <Button onPress={handleConfirm} style={styles.confirmButton} disabled={loading} size="lg">
          {loading ? 'Tracking...' : 'Track Item'}
        </Button>
        <Button
          variant="ghost"
          onPress={handleClose}
          style={styles.cancelButton}
          disabled={loading}
        >
          Cancel
        </Button>
      </View>
    </Modal>
  );
}
