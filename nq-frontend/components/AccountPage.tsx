import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import Card from './ui/card';
import { Button } from './ui/button';
import Input from './ui/input';
import Label from './ui/label';
import Separator from './ui/separator';
import Switch from './ui/switch';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, fontSize } from './ui/tokens';

function AccountPage() {
  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <Ionicons name="person" size={24} color={colors.primary} />
        <View>
          <Text style={styles.title}>Account Settings</Text>
          <Text style={styles.subtitle}>Manage your profile and preferences</Text>
        </View>
      </View>

      <Card style={styles.card}>
        <Text style={styles.sectionTitle}>Profile Information</Text>
        <View style={styles.profileSection}>
          <Avatar style={styles.avatar}>
            <AvatarImage src="https://i.pravatar.cc/150?img=10" />
            <AvatarFallback>JD</AvatarFallback>
          </Avatar>
          <View style={styles.avatarActions}>
            <Button variant="outline" size="sm">
              Change Photo
            </Button>
            <Text style={styles.avatarHelp}>JPG, PNG or GIF. Max size 2MB.</Text>
          </View>
        </View>

        <View style={styles.form}>
          <View style={styles.row}>
            <View style={styles.field}>
              <Label>First Name</Label>
              <Input defaultValue="John" />
            </View>
            <View style={styles.field}>
              <Label>Last Name</Label>
              <Input defaultValue="Doe" />
            </View>
          </View>

          <View style={styles.field}>
            <Label>Email</Label>
            <Input defaultValue="john.doe@example.com" keyboardType="email-address" />
          </View>

          <View style={styles.field}>
            <Label>Bio</Label>
            <Input defaultValue="Movie enthusiast and avid reader" multiline />
          </View>

          <Button>Save Changes</Button>
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
              <Text style={styles.notificationDescription}>Get notified when friends rate new media</Text>
            </View>
            <Switch value={true} />
          </View>
          <Separator />
          <View style={styles.notificationItem}>
            <View style={styles.notificationContent}>
              <Label>Recommendations</Label>
              <Text style={styles.notificationDescription}>Receive personalized recommendations</Text>
            </View>
            <Switch value={true} />
          </View>
          <Separator />
          <View style={styles.notificationItem}>
            <View style={styles.notificationContent}>
              <Label>Email Digest</Label>
              <Text style={styles.notificationDescription}>Weekly summary of new additions</Text>
            </View>
            <Switch value={false} />
          </View>
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
            <Button variant="outline" style={styles.themeButton}>Light</Button>
            <Button variant="outline" style={styles.themeButton}>Dark</Button>
            <Button variant="outline" style={styles.themeButton}>Auto</Button>
          </View>
        </View>
      </Card>

      <Card style={styles.card}>
        <View style={styles.sectionHeader}>
          <Ionicons name="shield" size={20} color={colors.primary} />
          <Text style={styles.sectionTitle}>Privacy & Security</Text>
        </View>
        <View style={styles.privacy}>
          <View style={styles.privacyItem}>
            <View style={styles.privacyContent}>
              <Label>Profile Visibility</Label>
              <Text style={styles.privacyDescription}>Make your profile visible to other users</Text>
            </View>
            <Switch value={true} />
          </View>
          <Separator />
          <View style={styles.privacyItem}>
            <View style={styles.privacyContent}>
              <Label>Show Activity</Label>
              <Text style={styles.privacyDescription}>Let friends see your recent activity</Text>
            </View>
            <Switch value={true} />
          </View>
          <Separator />
          <Button variant="outline">Change Password</Button>
        </View>
      </Card>

      <Card style={{ ...styles.card, ...styles.dangerCard }}>
        <View style={styles.sectionHeader}>
          <Ionicons name="log-out" size={20} color="#d4183d" />
          <Text style={styles.dangerTitle}>Danger Zone</Text>
        </View>
        <View style={styles.dangerActions}>
          <Button variant="outline" style={styles.dangerButton}>
            Log Out
          </Button>
          <Button variant="outline" style={styles.dangerButton}>
            Delete Account
          </Button>
        </View>
      </Card>
    </ScrollView>
  );
}

export default AccountPage;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    maxWidth: 1024,
    alignSelf: 'center',
    width: '100%',
    padding: spacing[6],
    gap: spacing[6],
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[3],
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
  card: {
    padding: spacing[6],
  },
  sectionTitle: {
    fontSize: fontSize.lg,
    fontWeight: '600',
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
    width: 80,
    height: 80,
  },
  avatarActions: {
    flex: 1,
    gap: spacing[2],
  },
  avatarHelp: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
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
    color: colors['muted-foreground'],
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
  privacy: {
    gap: spacing[4],
  },
  privacyItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  privacyContent: {
    flex: 1,
    gap: spacing[1],
  },
  privacyDescription: {
    fontSize: fontSize.sm,
    color: colors['muted-foreground'],
  },
  dangerCard: {
    borderColor: '#d4183d',
    backgroundColor: '#fef2f2',
  },
  dangerTitle: {
    color: '#d4183d',
    fontSize: fontSize.lg,
    fontWeight: '600',
  },
  dangerActions: {
    gap: spacing[3],
  },
  dangerButton: {
    borderColor: '#d4183d',
    color: '#d4183d',
  },
});