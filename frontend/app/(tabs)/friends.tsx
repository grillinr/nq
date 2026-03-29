import React, { useMemo, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TextInput,
  TouchableOpacity,
  RefreshControl,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { useMutation, useQuery } from '@apollo/client/react';
import { router } from 'expo-router';
import PageHeader, { useHeaderHeight } from '../../src/components/PageHeader';
import { UserAutocomplete, UserSuggestion } from '../../src/components/UserAutocomplete';
import { useTheme } from '../../src/components/ui/theme-provider';
import {
  spacing,
  fontSize,
  radii,
  layout,
  fontWeights,
  ColorPalette,
} from '../../src/components/ui/tokens';
import {
  ME_FRIENDS_QUERY,
  FRIENDS_ACTIVITY_QUERY,
  SEND_FRIEND_REQUEST_MUTATION,
  ACCEPT_FRIEND_REQUEST_MUTATION,
  DECLINE_FRIEND_REQUEST_MUTATION,
  REMOVE_FRIEND_MUTATION,
} from '../../src/lib/graphql';
import { MeFriendsQuery, FriendsActivityQuery } from '../../src/__generated__/graphql';

// ── styles ────────────────────────────────────────────────────────────────────

function createStyles(colors: ColorPalette, headerHeight: number) {
  return StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      paddingHorizontal: spacing[4],
      paddingTop: headerHeight + spacing[4],
      paddingBottom: layout.tabBarHeight + spacing[4],
    },
    // search bar
    searchWrapper: {
      position: 'relative',
      marginBottom: spacing[6],
    },
    searchInput: {
      height: 40,
      borderWidth: 1,
      borderColor: colors.border,
      borderRadius: radii.lg,
      paddingHorizontal: spacing[3],
      fontSize: fontSize.base,
      backgroundColor: colors.inputBackground,
      color: colors.foreground,
    },
    // section
    sectionTitle: {
      fontSize: fontSize.sm,
      fontWeight: fontWeights.semibold,
      color: colors.mutedForeground,
      textTransform: 'uppercase',
      letterSpacing: 0.8,
      marginBottom: spacing[2],
      marginTop: spacing[4],
    },
    emptyText: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      marginBottom: spacing[3],
    },
    // user row
    userRow: {
      flexDirection: 'row',
      alignItems: 'center',
      paddingVertical: spacing[3],
      borderBottomWidth: 1,
      borderBottomColor: colors.border,
      gap: spacing[3],
    },
    avatar: {
      width: 40,
      height: 40,
      borderRadius: radii.full,
      backgroundColor: colors.inputBackground,
      alignItems: 'center',
      justifyContent: 'center',
      flexShrink: 0,
    },
    avatarText: {
      fontSize: fontSize.base,
      fontWeight: fontWeights.semibold,
      color: colors.mutedForeground,
    },
    userName: {
      flex: 1,
      fontSize: fontSize.base,
      color: colors.foreground,
    },
    // action buttons
    actionBtn: {
      paddingHorizontal: spacing[3],
      paddingVertical: spacing[1],
      borderRadius: radii.md,
      backgroundColor: colors.primary,
    },
    actionBtnDestructive: {
      backgroundColor: colors.destructiveBackground,
    },
    actionBtnSecondary: {
      backgroundColor: colors.muted,
    },
    actionBtnText: {
      fontSize: fontSize.sm,
      color: colors.primaryForeground,
      fontWeight: fontWeights.medium,
    },
    actionBtnDestructiveText: {
      fontSize: fontSize.sm,
      color: colors.destructive,
      fontWeight: fontWeights.medium,
    },
    actionBtnSecondaryText: {
      fontSize: fontSize.sm,
      color: colors.foreground,
      fontWeight: fontWeights.medium,
    },
    actionBtnRow: {
      flexDirection: 'row',
      gap: spacing[2],
    },
    // activity item
    activityRow: {
      paddingVertical: spacing[3],
      borderBottomWidth: 1,
      borderBottomColor: colors.border,
      gap: spacing[1],
    },
    activityPrimary: {
      fontSize: fontSize.base,
      color: colors.foreground,
      flexShrink: 1,
    },
    activitySecondary: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
    },
    activityMediaTitle: {
      fontSize: fontSize.base,
      color: colors.primary,
      fontWeight: fontWeights.medium,
    },
    activityRow2: {
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[2],
    },
    activityTextBlock: {
      flex: 1,
    },
    sendingIndicator: {
      position: 'absolute',
      right: spacing[3],
      top: spacing[2],
    },
  });
}

// ── sub-components ────────────────────────────────────────────────────────────

type StylesType = ReturnType<typeof createStyles>;

function AvatarCircle({ name, styles }: { name: string; styles: StylesType }) {
  return (
    <View style={styles.avatar}>
      <Text style={styles.avatarText}>{name.charAt(0).toUpperCase()}</Text>
    </View>
  );
}

// ── main component ────────────────────────────────────────────────────────────

export default function FriendsPage() {
  const { colors } = useTheme();
  const headerHeight = useHeaderHeight();
  const styles = useMemo(() => createStyles(colors, headerHeight), [colors, headerHeight]);

  const [searchText, setSearchText] = useState('');
  const [suppressAutocomplete, setSuppressAutocomplete] = useState(false);
  const [searchBarHeight, setSearchBarHeight] = useState(40);

  // ── queries ──
  const {
    data: friendsData,
    loading: friendsLoading,
    refetch: refetchFriends,
  } = useQuery<MeFriendsQuery>(ME_FRIENDS_QUERY, { fetchPolicy: 'cache-and-network' });

  const {
    data: activityData,
    loading: activityLoading,
    refetch: refetchActivity,
  } = useQuery<FriendsActivityQuery>(FRIENDS_ACTIVITY_QUERY, {
    variables: { limit: 30 },
    fetchPolicy: 'cache-and-network',
  });

  const isRefreshing = friendsLoading || activityLoading;

  // ── mutations ──
  const [sendRequest, { loading: sendingRequest }] = useMutation(SEND_FRIEND_REQUEST_MUTATION, {
    onCompleted: () => {
      setSuppressAutocomplete(true);
      setSearchText('');
      refetchFriends();
    },
    onError: err => Alert.alert('Error', err.message),
  });

  const [acceptRequest] = useMutation(ACCEPT_FRIEND_REQUEST_MUTATION, {
    onCompleted: () => refetchFriends(),
    onError: err => Alert.alert('Error', err.message),
  });

  const [declineRequest] = useMutation(DECLINE_FRIEND_REQUEST_MUTATION, {
    onCompleted: () => refetchFriends(),
    onError: err => Alert.alert('Error', err.message),
  });

  const [removeFriend] = useMutation(REMOVE_FRIEND_MUTATION, {
    onCompleted: () => refetchFriends(),
    onError: err => Alert.alert('Error', err.message),
  });

  // ── derived data ──
  const me = friendsData?.me;
  const friends = me?.friends ?? [];
  const pendingRequests = me?.pendingFriendRequests ?? [];
  const sentRequests = me?.sentFriendRequests ?? [];
  const activities = activityData?.friendsActivity ?? [];

  // ── handlers ──
  const handleSelectUser = (user: UserSuggestion) => {
    setSuppressAutocomplete(true);
    setSearchText('');
    sendRequest({ variables: { toUserID: user.id } });
  };

  const handleRemoveFriend = (friendId: string, friendName: string) => {
    Alert.alert('Remove Friend', `Remove ${friendName} from your friends?`, [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Remove',
        style: 'destructive',
        onPress: () => removeFriend({ variables: { friendID: friendId } }),
      },
    ]);
  };

  const handleRefresh = async () => {
    await Promise.all([refetchFriends(), refetchActivity()]);
  };

  return (
    <View style={styles.container}>
      <PageHeader title="Friends" />
      <ScrollView
        style={styles.container}
        contentContainerStyle={styles.content}
        keyboardShouldPersistTaps="handled"
        refreshControl={<RefreshControl refreshing={isRefreshing} onRefresh={handleRefresh} />}
      >
        {/* ── Add Friend Search ── */}
        <Text style={styles.sectionTitle}>Add Friend</Text>
        <View style={styles.searchWrapper}>
          <TextInput
            style={styles.searchInput}
            placeholder="Search by name…"
            placeholderTextColor={colors.mutedForeground}
            value={searchText}
            onChangeText={text => {
              setSearchText(text);
              setSuppressAutocomplete(false);
            }}
            onLayout={e => setSearchBarHeight(e.nativeEvent.layout.height)}
            returnKeyType="search"
          />
          {sendingRequest && (
            <ActivityIndicator
              size="small"
              color={colors.mutedForeground}
              style={styles.sendingIndicator}
            />
          )}
          <UserAutocomplete
            query={searchText}
            suppress={suppressAutocomplete}
            onSelect={handleSelectUser}
            inputHeight={searchBarHeight}
          />
        </View>

        {/* ── Pending Requests ── */}
        {pendingRequests.length > 0 && (
          <>
            <Text style={styles.sectionTitle}>Pending Requests</Text>
            {pendingRequests.map(req => (
              <View key={req.id} style={styles.userRow}>
                <AvatarCircle name={req.from?.name ?? '?'} styles={styles} />
                <Text style={styles.userName} numberOfLines={1}>
                  {req.from?.name}
                </Text>
                <View style={styles.actionBtnRow}>
                  <TouchableOpacity
                    style={styles.actionBtn}
                    onPress={() => acceptRequest({ variables: { requestID: req.id } })}
                  >
                    <Text style={styles.actionBtnText}>Accept</Text>
                  </TouchableOpacity>
                  <TouchableOpacity
                    style={[styles.actionBtn, styles.actionBtnSecondary]}
                    onPress={() => declineRequest({ variables: { requestID: req.id } })}
                  >
                    <Text style={styles.actionBtnSecondaryText}>Decline</Text>
                  </TouchableOpacity>
                </View>
              </View>
            ))}
          </>
        )}

        {/* ── Sent Requests ── */}
        {sentRequests.length > 0 && (
          <>
            <Text style={styles.sectionTitle}>Sent Requests</Text>
            {sentRequests.map(req => (
              <View key={req.id} style={styles.userRow}>
                <AvatarCircle name={req.to?.name ?? '?'} styles={styles} />
                <Text style={styles.userName} numberOfLines={1}>
                  {req.to?.name}
                </Text>
                <Text style={styles.activitySecondary}>Pending</Text>
              </View>
            ))}
          </>
        )}

        {/* ── Friends List ── */}
        <Text style={styles.sectionTitle}>
          Friends {friends.length > 0 ? `(${friends.length})` : ''}
        </Text>
        {friends.length === 0 ? (
          <Text style={styles.emptyText}>No friends yet. Search above to add someone.</Text>
        ) : (
          friends.map(friend => (
            <View key={friend.id} style={styles.userRow}>
              <AvatarCircle name={friend.name} styles={styles} />
              <Text style={styles.userName} numberOfLines={1}>
                {friend.name}
              </Text>
              <TouchableOpacity
                style={[styles.actionBtn, styles.actionBtnDestructive]}
                onPress={() => handleRemoveFriend(friend.id, friend.name)}
              >
                <Text style={styles.actionBtnDestructiveText}>Remove</Text>
              </TouchableOpacity>
            </View>
          ))
        )}

        {/* ── Activity Feed ── */}
        <Text style={styles.sectionTitle}>Recent Activity</Text>
        {renderActivityFeed(activityLoading, activities, styles, colors, router)}
      </ScrollView>
    </View>
  );
}

function statusVerb(statusName?: string | null): string {
  if (!statusName) return 'logged';
  const lower = statusName.toLowerCase();
  if (lower.includes('complet') || lower.includes('finish')) return 'finished';
  if (lower.includes('start') || lower.includes('watching') || lower.includes('playing'))
    return 'started';
  if (lower.includes('drop')) return 'dropped';
  if (lower.includes('want') || lower.includes('plan') || lower.includes('wishlist'))
    return 'wants to watch';
  if (lower.includes('pause') || lower.includes('hold')) return 'put on hold';
  return 'logged';
}

type ActivityItem = NonNullable<FriendsActivityQuery['friendsActivity']>[number];
type RouterType = typeof router;

function renderActivityFeed(
  loading: boolean,
  activities: ActivityItem[],
  styles: StylesType,
  colors: ColorPalette,
  nav: RouterType
): React.ReactNode {
  if (loading && activities.length === 0) {
    return <ActivityIndicator size="small" color={colors.mutedForeground} />;
  }
  if (activities.length === 0) {
    return <Text style={styles.emptyText}>No recent activity from friends.</Text>;
  }
  return activities.map(item => (
    <TouchableOpacity
      key={item.id}
      style={styles.activityRow}
      onPress={() =>
        item.media?.id ? nav.push(`/media/${item.media.id}` as `/media/${string}`) : null
      }
    >
      <View style={styles.activityRow2}>
        <AvatarCircle name={item.user?.name ?? '?'} styles={styles} />
        <View style={styles.activityTextBlock}>
          <Text style={styles.activityPrimary} numberOfLines={1}>
            {item.user?.name}{' '}
            <Text style={styles.activitySecondary}>{statusVerb(item.status?.name)}</Text>
          </Text>
          <Text style={styles.activityMediaTitle} numberOfLines={1}>
            {item.media?.title ?? 'Unknown'}
          </Text>
        </View>
      </View>
    </TouchableOpacity>
  ));
}
