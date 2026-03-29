import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Haptics from 'expo-haptics';
import { useApolloClient, useMutation, useQuery } from '@apollo/client/react';
import { Avatar, AvatarImage } from '../../src/components/ui/avatar';
import Card from '../../src/components/ui/card';
import { Button } from '../../src/components/ui/button';
import Input from '../../src/components/ui/input';
import Label from '../../src/components/ui/label';
import Separator from '../../src/components/ui/separator';
import Switch from '../../src/components/ui/switch';
import PageHeader, { useHeaderHeight } from '../../src/components/PageHeader';
import { spacing, fontSize, fontWeights, layout } from '../../src/components/ui/tokens';
import { useTheme } from '../../src/components/ui/theme-provider';
import { getAccessToken } from '../../src/lib/auth';
import { logError, logInfo } from '../../src/lib/logger';
import { AccountSwitcherModal } from '../../src/components/AccountSwitcherModal';
import { clearCacheForAccountSwitch } from '../../src/lib/apolloClient';
import { useAuth } from '../../src/lib/AuthContext';
import { ME_QUERY, UPDATE_USER_MUTATION } from '../../src/lib/graphql';

const createStyles = (colors: ReturnType<typeof useTheme>['colors'], headerHeight: number) =>
  StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      maxWidth: 1024,
      alignSelf: 'center',
      width: '100%',
      padding: spacing[6],
      paddingTop: headerHeight,
      paddingBottom: layout.tabBarHeight,
      gap: spacing[6],
    },
    card: {
      padding: spacing[6],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: fontWeights.semibold,
      color: colors.foreground,
      marginBottom: spacing[4],
    },
    sectionHeader: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[3],
      marginBottom: spacing[4],
    },
    profileSection: {
      flexDirection: 'row',
      alignItems: 'flex-start',
      gap: spacing[6],
      marginBottom: spacing[6],
    },
    avatar: {
      width: 100,
      height: 100,
    },
    avatarActions: {
      flex: 1,
      gap: spacing[2],
    },
    avatarHelp: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
    },
    form: {
      gap: spacing[4],
    },
    row: {
      flexDirection: 'row',
      gap: spacing[4],
    },
    field: {
      flex: 1,
      gap: spacing[2],
    },
    notifications: {
      gap: spacing[4],
    },
    notificationItem: {
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'space-between',
    },
    notificationContent: {
      flex: 1,
      gap: spacing[1],
    },
    notificationDescription: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
    },
    appearance: {
      gap: spacing[4],
    },
    themeButtons: {
      flexDirection: 'row',
      gap: spacing[2],
    },
    themeButton: {
      flex: 1,
    },
    logoutCard: {
      padding: spacing[6],
    },
    accountActions: {
      gap: spacing[3],
    },
    accountActionRow: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[3],
    },
    dangerCard: {
      borderColor: colors.destructive,
      backgroundColor: colors.destructiveBackground,
      padding: spacing[4],
    },
    dangerLinkButton: {
      alignSelf: 'flex-start',
    },
    dangerLinkText: {
      fontSize: fontSize.sm,
      color: colors.destructive,
      textDecorationLine: 'underline',
    },
    subtitle: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      marginTop: spacing[2],
    },
    root: {
      flex: 1,
    },
    statusText: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      textAlign: 'center',
    },
  });

function AccountPage() {
  const { colors, theme, setTheme } = useTheme();
  const apolloClient = useApolloClient();
  const { login, logout, refreshAuth, storedAccounts, addAccount, hasToken: authHasToken } = useAuth();
  const [showAccountSwitcher, setShowAccountSwitcher] = React.useState(false);
  
  // Debug logging
  React.useEffect(() => {
    logInfo('[Account] storedAccounts changed:', storedAccounts.length, 'accounts');
    storedAccounts.forEach((account, index) => {
      logInfo(`[Account] Account ${index}:`, account.name, account.email);
    });
  }, [storedAccounts]);

  // Debug current user data
  React.useEffect(() => {
    if (currentUser) {
      logInfo('[Account] Current user data changed:', currentUser.name, currentUser.email);
    } else {
      logInfo('[Account] No current user data');
    }
  }, [currentUser]);

  // Debug hasToken state
  React.useEffect(() => {
    logInfo('[Account] authHasToken state changed:', authHasToken);
  }, [authHasToken]);
  const { data, loading, error, refetch } = useQuery(ME_QUERY, {
    fetchPolicy: 'cache-and-network', // Changed from 'cache-first' to ensure fresh data after account switch
    errorPolicy: 'all',
    skip: !authHasToken,
  });
  const [updateUser, { loading: saving }] = useMutation(UPDATE_USER_MUTATION);
  const currentUser = (data as any)?.me;
  const [firstName, setFirstName] = React.useState('');
  const [lastName, setLastName] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [statusMessage, setStatusMessage] = React.useState<string | null>(null);
  const avatarUrl =
    currentUser?.avatarUrl ??
    'https://toppng.com/uploads/preview/avatar-png-115540218987bthtxfhls.png';

  React.useEffect(() => {
    if (!currentUser) {
      setFirstName('');
      setLastName('');
      setEmail('');
      return;
    }
    const [first, ...rest] = currentUser.name.split(' ');
    setFirstName(first ?? '');
    setLastName(rest.join(' '));
    setEmail(currentUser.email ?? '');
  }, [currentUser]);

  const headerHeight = useHeaderHeight();
  const styles = React.useMemo(() => createStyles(colors, headerHeight), [colors, headerHeight]);

  // Remove the complex login handling - AuthContext now handles everything
  const handleLogin = async () => {
    try {
      logInfo('[Account] Starting login process');
      const success = await login();
      if (success) {
        logInfo('[Account] Login successful');
        await refreshAuth();
        await refetch();
      } else {
        logError('[Account] Login failed');
      }
    } catch (err) {
      logError('[Account] Login error:', err);
    }
  };

  const handleAddAccount = async () => {
    try {
      logInfo('[Account] Starting add account process');
      logInfo('[Account] Current storedAccounts before add:', storedAccounts.length);
      
      const success = await addAccount();
      logInfo('[Account] Add account result:', success);
      
      if (success) {
        logInfo('[Account] Add account successful, refreshing token state');
        // Don't call refreshAuth() here since addAccount() already calls it
        await refetch();
        
        logInfo('[Account] storedAccounts after add:', storedAccounts.length);
      } else {
        logInfo('[Account] Add account failed');
      }
    } catch (err) {
      logError('[Account] Add account error:', err);
    }
  };

  const handleLogout = async () => {
    try {
      await Haptics.notificationAsync(Haptics.NotificationFeedbackType.Warning);
      await logout();
      setStatusMessage(null);
      await refreshAuth();

      // Check if we still have a token (another account)
      const token = await getAccessToken();

      if (token) {
        await refetch();
      } else {
        setFirstName('');
        setLastName('');
        setEmail('');
        // Use our robust cache clearing function instead of direct clearStore()
        await clearCacheForAccountSwitch();
      }
    } catch (error) {
      logError('[Account] Logout error:', error);
    }
  };

  const handleSwitchAccount = () => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    setShowAccountSwitcher(true);
  };

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'auto') => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
    setTheme(newTheme);
  };

  const handleSave = async () => {
    if (!currentUser) return;
    setStatusMessage(null);
    const name = [firstName, lastName].filter(Boolean).join(' ').trim();
    try {
      await updateUser({
        variables: {
          id: currentUser.id,
          input: {
            name: name || currentUser.name,
            email: email || currentUser.email,
            avatarUrl: currentUser.avatarUrl,
          },
        },
      });
      setStatusMessage('Profile updated');
      await refetch();
    } catch (saveError) {
      setStatusMessage('Failed to update profile');
      logError(saveError);
    }
  };

  return (
    <View style={styles.root}>
      <PageHeader title="Account Settings" />
      <ScrollView style={styles.container} contentContainerStyle={styles.content}>
        <Card style={styles.card}>
          <Text style={styles.sectionTitle}>Profile Information</Text>
          {loading && <Text style={styles.statusText}>Loading profile…</Text>}
          {error && <Text style={styles.statusText}>Failed to load profile. Try again.</Text>}
          <View style={styles.profileSection}>
            <Avatar style={styles.avatar}>
              <AvatarImage src={avatarUrl} />
            </Avatar>
            <View style={styles.avatarActions}>
              <Button variant="outline" size="sm">
                Change Photo
              </Button>
              {currentUser ? (
                <Text style={styles.avatarHelp}>Signed in as {currentUser.email}</Text>
              ) : (
                <Text style={styles.avatarHelp}>Sign in to see your profile</Text>
              )}
            </View>
          </View>

          <View style={styles.form}>
            <View style={styles.row}>
              <View style={styles.field}>
                <Label>First Name</Label>
                <Input
                  placeholder="First Name"
                  value={firstName}
                  onChangeText={setFirstName}
                  disabled={!currentUser}
                />
              </View>
              <View style={styles.field}>
                <Label>Last Name</Label>
                <Input
                  placeholder="Last Name"
                  value={lastName}
                  onChangeText={setLastName}
                  disabled={!currentUser}
                />
              </View>
            </View>

            <View style={styles.field}>
              <Label>Email</Label>
              <Input
                keyboardType="email-address"
                placeholder="youremail@example.com"
                value={email}
                onChangeText={setEmail}
                disabled={!currentUser}
              />
            </View>

            <Button disabled={!currentUser || saving} onPress={handleSave}>
              {saving ? 'Saving…' : 'Save Changes'}
            </Button>
            {statusMessage && <Text style={styles.subtitle}>{statusMessage}</Text>}
          </View>
        </Card>

        <Card style={styles.card}>
          <View style={styles.sectionHeader}>
            <Ionicons name="notifications" size={20} color={colors.primary} />
            <Text style={styles.sectionTitle}>Notifications</Text>
          </View>
          <View style={styles.notifications}>
            <View style={styles.notificationItem}>
              <View style={styles.notificationContent}>
                <Label>Friend Activity</Label>
                <Text style={styles.notificationDescription}>
                  Get notified when friends rate new media
                </Text>
              </View>
              <Switch />
            </View>
            <Separator />
            <View style={styles.notificationItem}>
              <View style={styles.notificationContent}>
                <Label>Recommendations</Label>
                <Text style={styles.notificationDescription}>
                  Receive personalized recommendations
                </Text>
              </View>
              <Switch />
            </View>
            <Separator />
          </View>
        </Card>

        <Card style={styles.card}>
          <View style={styles.sectionHeader}>
            <Ionicons name="color-palette" size={20} color={colors.primary} />
            <Text style={styles.sectionTitle}>Appearance</Text>
          </View>
          <View style={styles.appearance}>
            <Label>Theme</Label>
            <View style={styles.themeButtons}>
              <Button
                variant={theme === 'light' ? 'default' : 'outline'}
                style={styles.themeButton}
                onPress={() => handleThemeChange('light')}
              >
                <Ionicons
                  name="sunny"
                  size={16}
                  color={theme === 'light' ? colors.primaryForeground : colors.foreground}
                />
              </Button>
              <Button
                variant={theme === 'dark' ? 'default' : 'outline'}
                style={styles.themeButton}
                onPress={() => handleThemeChange('dark')}
              >
                <Ionicons
                  name="moon"
                  size={16}
                  color={theme === 'dark' ? colors.primaryForeground : colors.foreground}
                />
              </Button>
              <Button
                variant={theme === 'auto' ? 'default' : 'outline'}
                style={styles.themeButton}
                onPress={() => handleThemeChange('auto')}
              >
                <Ionicons
                  name="sync"
                  size={16}
                  color={theme === 'auto' ? colors.primaryForeground : colors.foreground}
                />
              </Button>
            </View>
          </View>
        </Card>

        {/* Account Management — Switch Account / Log Out / Sign In */}
        <Card style={styles.logoutCard}>
          {authHasToken ? (
            <View style={styles.accountActions}>
              {storedAccounts.length > 1 && (
                <View style={styles.accountActionRow}>
                  <Ionicons name="swap-horizontal-outline" size={20} color={colors.foreground} />
                  <Button variant="ghost" onPress={handleSwitchAccount}>
                    Switch Account
                  </Button>
                </View>
              )}
              {storedAccounts.length > 1 && <Separator />}
              <View style={styles.accountActionRow}>
                <Ionicons name="add-circle-outline" size={20} color={colors.foreground} />
                <Button variant="ghost" onPress={handleAddAccount}>
                  Add Account
                </Button>
              </View>
              <Separator />
              <View style={styles.accountActionRow}>
                <Ionicons name="log-out-outline" size={20} color={colors.foreground} />
                <Button variant="ghost" onPress={handleLogout}>
                  Log Out
                </Button>
              </View>
            </View>
          ) : (
            <Button onPress={handleLogin}>Sign In with Auth0</Button>
          )}
        </Card>

        {/* Danger Zone — delete account only */}
        <Card style={styles.dangerCard}>
          <Button variant="ghost" style={styles.dangerLinkButton} onPress={() => {}}>
            <Text style={styles.dangerLinkText}>Delete Account</Text>
          </Button>
        </Card>
      </ScrollView>

      <AccountSwitcherModal
        visible={showAccountSwitcher}
        onClose={() => setShowAccountSwitcher(false)}
      />
    </View>
  );
}

export default AccountPage;
