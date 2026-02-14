import React, { useState } from 'react';
import { View, Text, Modal, StyleSheet, TouchableOpacity } from 'react-native';
import { StatusPicker, ActivityStatusId } from './ui/StatusPicker';
import { Button } from './ui/button';
import { useTheme } from './ui/ThemeProvider';

interface TrackItemModalProps {
  visible: boolean;
  onClose: () => void;
  onConfirm: (statusId: ActivityStatusId) => void;
  mediaTitle: string;
  loading?: boolean;
}

export function TrackItemModal({ 
  visible, 
  onClose, 
  onConfirm, 
  mediaTitle,
  loading = false 
}: TrackItemModalProps) {
  const { colors } = useTheme();
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
          style={[styles.modal, { backgroundColor: colors.card }]}
          onStartShouldSetResponder={() => true}
        >
          <Text style={[styles.title, { color: colors.foreground }]}>
            Track &quot;{mediaTitle}&quot;
          </Text>
          <Text style={[styles.subtitle, { color: colors.mutedForeground }]}>
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

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  modal: {
    width: '100%',
    maxWidth: 400,
    borderRadius: 12,
    padding: 24,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
    elevation: 8,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
    marginBottom: 20,
  },
  buttons: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 24,
  },
  button: {
    flex: 1,
  },
});
