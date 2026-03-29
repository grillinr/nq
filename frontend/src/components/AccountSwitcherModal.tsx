import React, { useMemo } from 'react';
import { View, Text, StyleSheet, Pressable, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Haptics from 'expo-haptics';
import { Modal } from './ui/modal';
import { Avatar, AvatarImage } from './ui/avatar';
import { Button } from './ui/button';
import Separator from './ui/separator';
import { spacing, fontSize, fontWeights, radii } from './ui/tokens';
import { useTheme } from './ui/theme-provider';
import { useAuth } from '../lib/AuthContext';
import { logInfo, logError } from '../lib/logger';

interface AccountSwitcherModalProps {
  visible: boolean;
  onClose: () => void;
}

const createStyles = (colors: ReturnType<typeof useTheme>['colors']) =>
  StyleSheet.create({
    accountList: {
      gap: spacing[2],
      maxHeight: 400,
    },
    accountItem: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[3],
      padding: spacing[4],
      borderRadius: radii.md,
      backgroundColor: colors.card,
      borderWidth: 1,
      borderColor: colors.border,
    },
    accountItemActive: {
      borderColor: colors.primary,
      backgroundColor: colors.secondary,
    },
    avatar: {
      width: 48,
      height: 48,
    },
    accountInfo: {
      flex: 1,
      gap: spacing[1],
    },
    accountName: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.medium,
      color: colors.foreground,
    },
    accountEmail: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
    },
    checkIcon: {
      marginLeft: spacing[2],
    },
    addAccountButton: {
      marginTop: spacing[4],
    },
    emptyState: {
      padding: spacing[6],
      alignItems: 'center',
      gap: spacing[3],
    },
    emptyText: {
      fontSize: fontSize.base,
      color: colors.mutedForeground,
      textAlign: 'center',
    },
  });

export function AccountSwitcherModal({ visible, onClose }: AccountSwitcherModalProps) {
  const { colors } = useTheme();
  const styles = useMemo(() => createStyles(colors), [colors]);
  const { storedAccounts, currentAccountId, switchAccount, addAccount } = useAuth();

  const handleSwitchAccount = async (accountId: string) => {
    if (accountId === currentAccountId) {
      logInfo('[AccountSwitcher] Tapped on already active account');
      onClose();
      return;
    }

    logInfo('[AccountSwitcher] Switching to account:', accountId);
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
    
    try {
      await switchAccount(accountId);
      logInfo('[AccountSwitcher] Account switch completed successfully');
    } catch (error) {
      logError('[AccountSwitcher] Account switch failed:', error);
    }
    
    onClose();
  };

  const handleAddAccount = async () => {
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    onClose();
    await addAccount();
  };

  return (
    <Modal visible={visible} onClose={onClose} title="Switch Account">
      {storedAccounts.length === 0 ? (
        <View style={styles.emptyState}>
          <Ionicons name="people-outline" size={48} color={colors.mutedForeground} />
          <Text style={styles.emptyText}>No accounts found. Please sign in to add an account.</Text>
        </View>
      ) : (
        <ScrollView style={styles.accountList} showsVerticalScrollIndicator={false}>
          {storedAccounts.map((account, index) => {
            const isActive = account.id === currentAccountId;
            const avatarUrl =
              account.avatarUrl ??
              'https://toppng.com/uploads/preview/avatar-png-115540218987bthtxfhls.png';

            return (
              <React.Fragment key={account.id}>
                {index > 0 && <Separator />}
                <Pressable
                  style={[styles.accountItem, isActive && styles.accountItemActive]}
                  onPress={() => handleSwitchAccount(account.id)}
                >
                  <Avatar style={styles.avatar}>
                    <AvatarImage src={avatarUrl} />
                  </Avatar>
                  <View style={styles.accountInfo}>
                    <Text style={styles.accountName}>{account.name}</Text>
                    <Text style={styles.accountEmail}>{account.email}</Text>
                  </View>
                  {isActive && (
                    <Ionicons
                      name="checkmark-circle"
                      size={24}
                      color={colors.primary}
                      style={styles.checkIcon}
                    />
                  )}
                </Pressable>
              </React.Fragment>
            );
          })}
        </ScrollView>
      )}

      <Button style={styles.addAccountButton} variant="outline" onPress={handleAddAccount}>
        <Ionicons name="add-circle-outline" size={20} color={colors.foreground} />
        <Text> Add Account</Text>
      </Button>
    </Modal>
  );
}
