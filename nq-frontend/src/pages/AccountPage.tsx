import React from "react";
import { View, Text, StyleSheet, ScrollView } from "react-native";
import { Avatar, AvatarImage } from "../components/ui/avatar";
import Card from "../components/ui/card";
import { Button } from "../components/ui/button";
import Input from "../components/ui/input";
import Label from "../components/ui/label";
import Separator from "../components/ui/separator";
import Switch from "../components/ui/switch";
import { Ionicons } from "@expo/vector-icons";
import { spacing, fontSize } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";

function AccountPage() {
  const { colors, theme, setTheme } = useTheme();

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      maxWidth: 1024,
      alignSelf: "center",
      width: "100%",
      padding: spacing[6],
      gap: spacing[6],
    },
    header: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[3],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: "600",
      color: colors.primary,
    },
    subtitle: {
      fontSize: fontSize.sm,
      color: colors["muted-foreground"],
    },
    card: {
      padding: spacing[6],
    },
    sectionTitle: {
      fontSize: fontSize.lg,
      fontWeight: "600",
      color: colors.foreground,
      marginBottom: spacing[4],
    },
    sectionHeader: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[3],
      marginBottom: spacing[4],
    },
    profileSection: {
      flexDirection: "row",
      alignItems: "flex-start",
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
      color: colors["muted-foreground"],
    },
    form: {
      gap: spacing[4],
    },
    row: {
      flexDirection: "row",
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
      flexDirection: "row",
      alignItems: "center",
      justifyContent: "space-between",
    },
    notificationContent: {
      flex: 1,
      gap: spacing[1],
    },
    notificationDescription: {
      fontSize: fontSize.sm,
      color: colors["muted-foreground"],
    },
    appearance: {
      gap: spacing[4],
    },
    themeButtons: {
      flexDirection: "row",
      gap: spacing[2],
    },
    themeButton: {
      flex: 1,
    },
    privacy: {
      gap: spacing[4],
    },
    privacyItem: {
      flexDirection: "row",
      alignItems: "center",
      justifyContent: "space-between",
    },
    privacyContent: {
      flex: 1,
      gap: spacing[1],
    },
    privacyDescription: {
      fontSize: fontSize.sm,
      color: colors["muted-foreground"],
    },
    dangerCard: {
      borderColor: colors.destructive,
      backgroundColor: colors["destructive-background"],
    },
    dangerTitle: {
      color: colors.destructive,
      fontSize: fontSize.lg,
      fontWeight: "600",
    },
    dangerActions: {
      gap: spacing[3],
    },
    dangerButton: {
      borderColor: colors.destructive,
      color: colors.destructive,
    },
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <Ionicons name="person" size={24} color={colors.primary} />
        <View>
          <Text style={styles.title}>Account Settings</Text>
          <Text style={styles.subtitle}>
            Manage your profile and preferences
          </Text>
        </View>
      </View>

      <Card style={styles.card}>
        <Text style={styles.sectionTitle}>Profile Information</Text>
        <View style={styles.profileSection}>
          <Avatar style={styles.avatar}>
            <AvatarImage src="https://toppng.com/uploads/preview/avatar-png-115540218987bthtxfhls.png" />
          </Avatar>
          <View style={styles.avatarActions}>
            <Button variant="outline" size="sm">
              Change Photo
            </Button>
          </View>
        </View>

        <View style={styles.form}>
          <View style={styles.row}>
            <View style={styles.field}>
              <Label>First Name</Label>
              <Input placeholder="First Name" />
            </View>
            <View style={styles.field}>
              <Label>Last Name</Label>
              <Input placeholder="Last Name" />
            </View>
          </View>

          <View style={styles.field}>
            <Label>Email</Label>
            <Input
              keyboardType="email-address"
              placeholder="youremail@example.com"
            />
          </View>

          <View style={styles.field}>
            <Label>Bio</Label>
            <Input placeholder="Short bio about yourself" multiline />
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
              variant={theme === "light" ? "default" : "outline"}
              style={styles.themeButton}
              onPress={() => setTheme("light")}
            >
              <Ionicons
                name="sunny"
                size={16}
                color={
                  theme === "light"
                    ? colors["primary-foreground"]
                    : colors.foreground
                }
              />
            </Button>
            <Button
              variant={theme === "dark" ? "default" : "outline"}
              style={styles.themeButton}
              onPress={() => setTheme("dark")}
            >
              <Ionicons
                name="moon"
                size={16}
                color={
                  theme === "dark"
                    ? colors["primary-foreground"]
                    : colors.foreground
                }
              />
            </Button>
            <Button
              variant={theme === "auto" ? "default" : "outline"}
              style={styles.themeButton}
              onPress={() => setTheme("auto")}
            >
              <Ionicons
                name="sync"
                size={16}
                color={
                  theme === "dark"
                    ? colors["primary-foreground"]
                    : colors.foreground
                }
              />
            </Button>
          </View>
        </View>
      </Card>

      <Card style={{ ...styles.card, ...styles.dangerCard }}>
        <View style={styles.sectionHeader}>
          <Ionicons name="shield" size={20} color={colors.destructive} />
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
