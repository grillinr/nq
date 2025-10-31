import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import Card from './ui/card';
import { Button } from './ui/button';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, radii, fontSize } from './ui/tokens';

interface Friend {
  id: number;
  name: string;
  avatar: string;
  recentActivity: {
    title: string;
    type: "movie" | "tv" | "book" | "music" | "game";
    rating: number;
    timestamp: string;
  }[];
}

const mockFriends: Friend[] = [
  {
    id: 1,
    name: "Sarah Johnson",
    avatar: "https://i.pravatar.cc/150?img=1",
    recentActivity: [
      { title: "Inception", type: "movie", rating: 9.0, timestamp: "2 hours ago" },
      { title: "Breaking Bad", type: "tv", rating: 9.5, timestamp: "1 day ago" },
      { title: "1984", type: "book", rating: 8.5, timestamp: "3 days ago" },
    ],
  },
  {
    id: 2,
    name: "Mike Chen",
    avatar: "https://i.pravatar.cc/150?img=2",
    recentActivity: [
      { title: "The Dark Knight", type: "movie", rating: 9.0, timestamp: "5 hours ago" },
      { title: "Abbey Road", type: "music", rating: 9.2, timestamp: "2 days ago" },
    ],
  },
  {
    id: 3,
    name: "Emily Davis",
    avatar: "https://i.pravatar.cc/150?img=3",
    recentActivity: [
      { title: "Elden Ring", type: "game", rating: 9.5, timestamp: "1 hour ago" },
      { title: "Stranger Things", type: "tv", rating: 8.7, timestamp: "6 hours ago" },
      { title: "Project Hail Mary", type: "book", rating: 8.9, timestamp: "1 day ago" },
    ],
  },
  {
    id: 4,
    name: "Alex Rodriguez",
    avatar: "https://i.pravatar.cc/150?img=4",
    recentActivity: [
      { title: "The Great Gatsby", type: "book", rating: 7.8, timestamp: "3 hours ago" },
      { title: "Parasite", type: "movie", rating: 8.6, timestamp: "1 day ago" },
    ],
  },
  {
    id: 5,
    name: "Jessica Lee",
    avatar: "https://i.pravatar.cc/150?img=5",
    recentActivity: [
      { title: "The Crown", type: "tv", rating: 8.6, timestamp: "4 hours ago" },
      { title: "Kind of Blue", type: "music", rating: 8.8, timestamp: "2 days ago" },
      { title: "Zelda: Breath of the Wild", type: "game", rating: 9.6, timestamp: "5 days ago" },
    ],
  },
];

const typeIcons = {
  movie: 'film',
  tv: 'tv',
  book: 'book',
  music: 'musical-notes',
  game: 'game-controller',
};

function FriendsPage() {
  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <View style={styles.headerContent}>
          <Ionicons name="people" size={24} color={colors.primary} />
          <View>
            <Text style={styles.title}>Friends Activity</Text>
            <Text style={styles.subtitle}>See what your friends are watching, reading, and playing</Text>
          </View>
        </View>
        <Button variant="outline">
          <Ionicons name="people" size={16} />
          <Text style={styles.addFriendsText}>Add Friends</Text>
        </Button>
      </View>

      <View style={styles.friendsGrid}>
        {mockFriends.map((friend) => (
          <Card key={friend.id} style={styles.friendCard}>
            <View style={styles.friendHeader}>
              <Avatar style={styles.friendAvatar}>
                <AvatarImage src={friend.avatar} alt={friend.name} />
                <AvatarFallback>{friend.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
              </Avatar>
              <View style={styles.friendInfo}>
                <Text style={styles.friendName}>{friend.name}</Text>
                <Text style={styles.activityCount}>{friend.recentActivity.length} recent activities</Text>
              </View>
            </View>

            <View style={styles.activities}>
              {friend.recentActivity.map((activity, index) => {
                const iconName = typeIcons[activity.type] as keyof typeof Ionicons.glyphMap;
                return (
                  <View key={index} style={styles.activityItem}>
                    <View style={styles.activityIcon}>
                      <Ionicons name={iconName} size={16} color={colors.primary} />
                    </View>
                    <View style={styles.activityContent}>
                      <Text style={styles.activityTitle} numberOfLines={1}>{activity.title}</Text>
                      <Text style={styles.activityTimestamp}>{activity.timestamp}</Text>
                    </View>
                    <View style={styles.activityRating}>
                      <Ionicons name="star" size={12} color="#fbbf24" />
                      <Text style={styles.ratingText}>{activity.rating.toFixed(1)}</Text>
                    </View>
                  </View>
                );
              })}
            </View>
          </Card>
        ))}
      </View>

      <Card style={styles.discoverCard}>
        <Ionicons name="people" size={48} color={colors.primary} />
        <Text style={styles.discoverTitle}>Discover More Friends</Text>
        <Text style={styles.discoverDescription}>
          Connect with friends to share recommendations and see what they&apos;re enjoying
        </Text>
        <Button>Find Friends</Button>
      </Card>
    </ScrollView>
  );
}

export default FriendsPage;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing[6],
    gap: spacing[6],
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  headerContent: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
    flex: 1,
  },
  title: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.foreground,
  },
  subtitle: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
  },
  addFriendsText: {
    marginLeft: spacing[2],
    fontSize: fontSize.sm,
  },
  friendsGrid: {
    gap: spacing[6],
  },
  friendCard: {
    padding: spacing[6],
  },
  friendHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[4],
    marginBottom: spacing[4],
  },
  friendAvatar: {
    width: 48,
    height: 48,
  },
  friendInfo: {
    flex: 1,
  },
  friendName: {
    fontSize: fontSize.base,
    fontWeight: '500',
    marginBottom: spacing[1],
  },
  activityCount: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
  },
  activities: {
    gap: spacing[3],
  },
  activityItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
    padding: spacing[3],
    backgroundColor: colors.muted,
    borderRadius: radii.lg,
  },
  activityIcon: {
    width: 32,
    height: 32,
    backgroundColor: colors.background,
    borderRadius: radii.md,
    alignItems: 'center',
    justifyContent: 'center',
  },
  activityContent: {
    flex: 1,
    minWidth: 0,
  },
  activityTitle: {
    fontSize: fontSize.sm,
    marginBottom: spacing[1],
  },
  activityTimestamp: {
    fontSize: fontSize.xs,
    color: colors['muted-foreground'],
  },
  activityRating: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[1],
  },
  ratingText: {
    fontSize: fontSize.sm,
    color: colors.foreground,
  },
  discoverCard: {
    padding: spacing[8],
    alignItems: 'center',
    backgroundColor: colors.muted,
    gap: spacing[3],
  },
  discoverTitle: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    textAlign: 'center',
  },
  discoverDescription: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
    textAlign: 'center',
    maxWidth: 300,
  },
});