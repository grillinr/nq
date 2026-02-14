import React from "react";
import { View, Text, StyleSheet, ScrollView } from "react-native";
import { Avatar, AvatarImage } from "../components/ui/avatar";
import Card from "../components/ui/card";
import { Button } from "../components/ui/button";
import Input from "../components/ui/input";
import Label from "../components/ui/label";
import Separator from "../components/ui/separator";
import Switch from "../components/ui/switch";
import PageHeader from "../components/PageHeader";
import { Ionicons } from "@expo/vector-icons";
import { spacing, fontSize } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { useScrollHeader } from "../hooks/useScrollHeader";
import { getAccessToken } from "../../lib/auth";
import { useAuth } from "../../lib/AuthContext";
import { useApolloClient, useMutation, useQuery } from "@apollo/client/react";
import { ME_QUERY, UPDATE_USER_MUTATION } from "../../lib/graphql";

function AccountPage() {
  const { colors, theme, setTheme } = useTheme();
  const apolloClient = useApolloClient();
  const { login, logout: logoutFromAuth, refreshAuth } = useAuth();
  const { isHeaderVisible, handleScroll: handleHeaderScroll } = useScrollHeader(50);
  const [hasToken, setHasToken] = React.useState(false);
  const { data, loading, error, refetch } = useQuery(ME_QUERY, {
    fetchPolicy: "cache-first",
    errorPolicy: "all",
    skip: !hasToken,
  });
  const [updateUser, { loading: saving }] = useMutation(UPDATE_USER_MUTATION);
  const currentUser = (data as any)?.me;
  const [firstName, setFirstName] = React.useState("");
  const [lastName, setLastName] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [avatarUrlInput, setAvatarUrlInput] = React.useState("");
  const [statusMessage, setStatusMessage] = React.useState<string | null>(null);
  const [avatarError, setAvatarError] = React.useState<string | null>(null);
  const avatarUrl =
    currentUser?.avatarUrl ??
    "https://toppng.com/uploads/preview/avatar-png-115540218987bthtxfhls.png";

  React.useEffect(() => {
    let active = true;
    getAccessToken().then((token) => {
      if (active) setHasToken(!!token);
    });
    return () => {
      active = false;
    };
  }, []);

  React.useEffect(() => {
    if (!currentUser) {
      setFirstName("");
      setLastName("");
      setEmail("");
      setAvatarUrlInput("");
      return;
    }
    const [first, ...rest] = currentUser.name.split(" ");
    setFirstName(first ?? "");
    setLastName(rest.join(" "));
    setEmail(currentUser.email ?? "");
    setAvatarUrlInput(currentUser.avatarUrl ?? "");
  }, [currentUser]);

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
      paddingTop: 140,
      gap: spacing[6],
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
      color: colors.mutedForeground,
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
      color: colors.mutedForeground,
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
      color: colors.mutedForeground,
    },
    dangerCard: {
      borderColor: colors.destructive,
      backgroundColor: colors.destructiveBackground,
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
    subtitle: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      marginTop: spacing[2],
    },
  });

  const handleLogin = async () => {
    await login();
    await refreshAuth();
    setHasToken(true);
    await refetch();
  };

  const handleLogout = async () => {
    await logoutFromAuth();
    setStatusMessage(null);
    setHasToken(false);
    setFirstName("");
    setLastName("");
    setEmail("");
    setAvatarUrlInput("");
    setAvatarError(null);
    apolloClient.clearStore().catch(() => undefined);
  };

  const handleSave = async () => {
    if (!currentUser) return;
    setStatusMessage(null);
    const name = [firstName, lastName].filter(Boolean).join(" ").trim();
    const avatarCandidate = avatarUrlInput.trim();
    if (avatarCandidate && !isValidUrl(avatarCandidate)) {
      setAvatarError("Avatar URL must be a valid http(s) URL");
      return;
    }
    setAvatarError(null);
    try {
      await updateUser({
        variables: {
          id: currentUser.id,
          input: {
            name: name || currentUser.name,
            email: email || currentUser.email,
            avatarUrl: avatarCandidate || currentUser.avatarUrl,
          },
        },
      });
      setStatusMessage("Profile updated");
      await refetch();
    } catch (saveError) {
      setStatusMessage("Failed to update profile");
      console.error(saveError);
    }
  };

  return (
    <View style={{ flex: 1 }}>
      <PageHeader 
        title="Account Settings" 
        subtitle="Manage your profile and preferences"
        visible={isHeaderVisible}
      />
      <ScrollView 
        style={styles.container} 
        contentContainerStyle={styles.content}
        onScroll={handleHeaderScroll}
        scrollEventThrottle={16}
      >
      <Card style={styles.card}>
        <Text style={styles.sectionTitle}>Profile Information</Text>
        {loading && <Text style={{ fontSize: fontSize.sm, color: colors.mutedForeground, textAlign: "center" }}>Loading profile…</Text>}
        {error && (
          <Text style={{ fontSize: fontSize.sm, color: colors.mutedForeground, textAlign: "center" }}>
            Failed to load profile. Try again.
          </Text>
        )}
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

          <View style={styles.field}>
            <Label>Avatar URL</Label>
            <Input
              placeholder="https://"
              value={avatarUrlInput}
              onChangeText={(value) => {
                setAvatarUrlInput(value);
                if (avatarError) setAvatarError(null);
              }}
              disabled={!currentUser}
            />
            {avatarError && <Text style={styles.subtitle}>{avatarError}</Text>}
          </View>

          <Button disabled={!currentUser || saving} onPress={handleSave}>
            {saving ? "Saving…" : "Save Changes"}
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
              variant={theme === "light" ? "default" : "outline"}
              style={styles.themeButton}
              onPress={() => setTheme("light")}
            >
              <Ionicons
                name="sunny"
                size={16}
                color={
                  theme === "light"
                    ? colors.primaryForeground
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
                    ? colors.primaryForeground
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
                    ? colors.primaryForeground
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
          {hasToken ? (
            <Button variant="outline" style={styles.dangerButton} onPress={handleLogout}>
              Log Out
            </Button>
          ) : (
            <Button onPress={handleLogin}>Sign In with Auth0</Button>
          )}
          <Button variant="outline" style={styles.dangerButton}>
            Delete Account
          </Button>
        </View>
      </Card>
    </ScrollView>
    </View>
  );
}

function isValidUrl(value: string): boolean {
  try {
    const url = new URL(value);
    return url.protocol === "http:" || url.protocol === "https:";
  } catch {
    return false;
  }
}

export default AccountPage;
