import React, { useState, useEffect } from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createStackNavigator } from "@react-navigation/stack";
import { StatusBar } from "expo-status-bar";
import {
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  ScrollView,
} from "react-native";
import { SafeAreaProvider, SafeAreaView } from "react-native-safe-area-context";
import AccountsScreen from "./src/screens/AccountsScreen";
import ConnectionDemo from "./src/components/ConnectionDemo";
import AccountService from "./src/services/accountService";
import { AccountConnection } from "./src/types/accounts";

// Define the stack navigator types
type RootStackParamList = {
  Home: undefined;
  Profile: undefined;
  Settings: undefined;
  Accounts: undefined;
};

const Stack = createStackNavigator<RootStackParamList>();

// Home Screen Component
function HomeScreen({ navigation }: any) {
  const [connections, setConnections] = useState<AccountConnection[]>([]);

  useEffect(() => {
    loadConnections();
  }, []);

  const loadConnections = async () => {
    try {
      await AccountService.initialize();
      const loadedConnections = await AccountService.getConnections();
      setConnections(loadedConnections);
    } catch (error) {
      console.error("Failed to load connections:", error);
    }
  };

  const connectedCount = connections.filter((conn) => conn.isConnected).length;

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <View style={styles.header}>
          <Text style={styles.title}>Welcome to NQ</Text>
          <Text style={styles.subtitle}>Your Personal Dashboard</Text>
        </View>

        <View style={styles.cardContainer}>
          <TouchableOpacity
            style={styles.card}
            onPress={() => navigation.navigate("Accounts")}
          >
            <Text style={styles.cardTitle}>Account Connections</Text>
            <Text style={styles.cardDescription}>
              Connect your Spotify, Goodreads, Letterboxd, and other accounts
            </Text>
            <View style={styles.connectionStatus}>
              <Text style={styles.connectionText}>
                {connectedCount} of{" "}
                {
                  Object.keys(require("./src/types/accounts").SERVICE_CONFIGS)
                    .length
                }{" "}
                services connected
              </Text>
            </View>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.card}
            onPress={() => navigation.navigate("Profile")}
          >
            <Text style={styles.cardTitle}>Profile</Text>
            <Text style={styles.cardDescription}>
              View and edit your profile information
            </Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={styles.card}
            onPress={() => navigation.navigate("Settings")}
          >
            <Text style={styles.cardTitle}>Settings</Text>
            <Text style={styles.cardDescription}>
              Configure app preferences
            </Text>
          </TouchableOpacity>
        </View>

        {connectedCount > 0 && (
          <View style={styles.demoSection}>
            <Text style={styles.sectionTitle}>Your Connected Services</Text>
            <ConnectionDemo connections={connections} />
          </View>
        )}

        <View style={styles.infoSection}>
          <Text style={styles.infoTitle}>Quick Actions</Text>
          <Text style={styles.infoText}>
            Tap on any card above to navigate to different sections of the app.
          </Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

// Profile Screen Component
function ProfileScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <View style={styles.header}>
          <Text style={styles.title}>Profile</Text>
          <Text style={styles.subtitle}>Your Information</Text>
        </View>

        <View style={styles.profileSection}>
          <View style={styles.profileItem}>
            <Text style={styles.profileLabel}>Name:</Text>
            <Text style={styles.profileValue}>John Doe</Text>
          </View>
          <View style={styles.profileItem}>
            <Text style={styles.profileLabel}>Email:</Text>
            <Text style={styles.profileValue}>john.doe@example.com</Text>
          </View>
          <View style={styles.profileItem}>
            <Text style={styles.profileLabel}>Member Since:</Text>
            <Text style={styles.profileValue}>January 2024</Text>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

// Settings Screen Component
function SettingsScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContainer}>
        <View style={styles.header}>
          <Text style={styles.title}>Settings</Text>
          <Text style={styles.subtitle}>App Configuration</Text>
        </View>

        <View style={styles.settingsSection}>
          <View style={styles.settingItem}>
            <Text style={styles.settingLabel}>Notifications</Text>
            <Text style={styles.settingValue}>Enabled</Text>
          </View>
          <View style={styles.settingItem}>
            <Text style={styles.settingLabel}>Dark Mode</Text>
            <Text style={styles.settingValue}>Light</Text>
          </View>
          <View style={styles.settingItem}>
            <Text style={styles.settingLabel}>Language</Text>
            <Text style={styles.settingValue}>English</Text>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

// Main App Component
export default function App() {
  return (
    <SafeAreaProvider>
      <NavigationContainer>
        <Stack.Navigator
          initialRouteName="Home"
          screenOptions={{
            headerStyle: {
              backgroundColor: "#1a1a1a",
            },
            headerTintColor: "#fff",
            headerTitleStyle: {
              fontWeight: "bold",
            },
          }}
        >
          <Stack.Screen
            name="Home"
            component={HomeScreen}
            options={{ title: "NQ" }}
          />
          <Stack.Screen
            name="Accounts"
            component={AccountsScreen}
            options={{ title: "Account Connections" }}
          />
          <Stack.Screen
            name="Profile"
            component={ProfileScreen}
            options={{ title: "Profile" }}
          />
          <Stack.Screen
            name="Settings"
            component={SettingsScreen}
            options={{ title: "Settings" }}
          />
        </Stack.Navigator>
        <StatusBar style="light" />
      </NavigationContainer>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#f5f5f5",
  },
  scrollContainer: {
    flexGrow: 1,
    padding: 20,
  },
  header: {
    alignItems: "center",
    marginBottom: 30,
    marginTop: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: "#666",
    textAlign: "center",
  },
  cardContainer: {
    marginBottom: 30,
  },
  card: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 16,
    shadowColor: "#000",
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  cardTitle: {
    fontSize: 20,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 8,
  },
  cardDescription: {
    fontSize: 14,
    color: "#666",
    lineHeight: 20,
  },
  connectionStatus: {
    backgroundColor: "#e0f7fa",
    padding: 10,
    borderRadius: 8,
    marginTop: 10,
    alignItems: "center",
  },
  connectionText: {
    fontSize: 14,
    color: "#00796b",
    fontWeight: "bold",
  },
  demoSection: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 15,
  },
  infoSection: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
  },
  infoTitle: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 12,
  },
  infoText: {
    fontSize: 14,
    color: "#666",
    lineHeight: 20,
  },
  profileSection: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
  },
  profileItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#f0f0f0",
  },
  profileLabel: {
    fontSize: 16,
    fontWeight: "600",
    color: "#1a1a1a",
  },
  profileValue: {
    fontSize: 16,
    color: "#666",
  },
  settingsSection: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
  },
  settingItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#f0f0f0",
  },
  settingLabel: {
    fontSize: 16,
    fontWeight: "600",
    color: "#1a1a1a",
  },
  settingValue: {
    fontSize: 16,
    color: "#666",
  },
});
