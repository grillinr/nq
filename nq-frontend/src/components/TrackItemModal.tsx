import React, { useMemo, useState } from 'react';
import { View, Text, Modal, StyleSheet, TouchableOpacity } from 'react-native';
import { StatusPicker, ActivityStatusId } from './ui/status-picker';
import { Button } from './ui/button';
import { useTheme } from './ui/theme-provider';
import {
  createShadows,
  fontSize,
  fontWeights,
  radii,
  spacing,
  ColorPalette,
} from './ui/tokens';

interface TrackItemModalProps {
  visible: boolean;
  onClose: () => void;
  onConfirm: (statusId: ActivityStatusId) => void;
  mediaTitle: string;
  loading?: boolean;
}

function createStyles(colors: ColorPalette) {
  const shadows = createShadows(colors);
  return StyleSheet.create({
    overlay: {
      flex: 1,
      backgroundColor: 'rgba(0,0,0,0.5)',
      justifyContent: 'center',
      alignItems: 'center',
      padding: spacing[5],
    },
    modal: {
      width: '100%',
      maxWidth: 400,
      borderRadius: radii.xl,
      padding: spacing[6],
      backgroundColor: colors.card,
      ...shadows.modal,
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: fontWeights.bold,
      marginBottom: spacing[2],
      color: colors.foreground,
    },
    subtitle: {
      fontSize: fontSize.sm,
      marginBottom: spacing[5],
      color: colors.mutedForeground,
    },
    buttons: {
      flexDirection: 'row',
      gap: spacing[3],
      marginTop: spacing[6],
    },
    button: {
      flex: 1,
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
  const styles = useMemo(() => createStyles(colors), [colors]);
  const [selectedStatus, setSelectedStatus] = useState<ActivityStatusId>(1);

  const handleConfirm = () => {
    onConfirm(selectedStatus);
  };

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <TouchableOpacity
        style={styles.overlay}
        activeOpacity={1}
        onPress={onClose}
      >
        <View
          style={styles.modal}
          onStartShouldSetResponder={() => true}
        >
          <Text style={styles.title}>
            Track &quot;{mediaTitle}&quot;
          </Text>
          <Text style={styles.subtitle}>
            What&apos;s your status?
          </Text>

          <StatusPicker value={selectedStatus} onChange={setSelectedStatus} />

          <View style={styles.buttons}>
            <Button
              variant="outline"
              onPress={onClose}
              style={styles.button}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button
              onPress={handleConfirm}
              style={styles.button}
              disabled={loading}
            >
              {loading ? 'Tracking...' : 'Track Item'}
            </Button>
          </View>
        </View>
      </TouchableOpacity>
    </Modal>
  );
}
