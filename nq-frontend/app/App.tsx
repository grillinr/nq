import React from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { Ionicons } from "@expo/vector-icons";
import HomePage from "../src/pages/HomePage";
import AddMediaPage from "../src/pages/AddMediaPage";
import AccountPage from "../src/pages/AccountPage";
import HistoryPage from "../src/pages/HistoryPage";
import FriendsPage from "../src/pages/FriendsPage";
import { createMedia } from "../lib/createMedia";
import { useTheme } from "../src/components/ui/ThemeProvider";
import {
  ApolloClient,
  HttpLink,
  InMemoryCache,
  ApolloLink,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { ApolloProvider, useApolloClient } from "@apollo/client/react";
import { Media } from "../src/types";
import { useAuth } from "../lib/AuthContext";
import { createActivity } from "../lib/createActivity";
import AuthPromptPage from "../src/pages/AuthPromptPage";
import { getAccessToken } from "../lib/auth";
import { GET_MOVIES_QUERY, ME_ACTIVITIES_QUERY } from "../lib/graphql";

const Tab = createBottomTabNavigator();

// Initialize Apollo Client
const API_URL = process.env.EXPO_PUBLIC_API_URL || "http://localhost:8080/graphql";

const authLink = setContext(async (_, { headers }) => {
  const token = await getAccessToken();
  return {
    headers: {
      ...headers,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  };
});

const client = new ApolloClient({
  link: ApolloLink.from([authLink, new HttpLink({ uri: API_URL })]),
  cache: new InMemoryCache(),
});

function AppContent() {
  const { colors } = useTheme();
  const apolloClient = useApolloClient();
  const [isAddingMedia, setIsAddingMedia] = React.useState(false);
  const { hasToken, isChecking, login } = useAuth();

  const handleAddMedia = async (newMedia: Omit<Media, "id">) => {
    setIsAddingMedia(true);
    try {
      const result = await createMedia(newMedia);
      if (result?.id) {
        await createActivity({
          mediaId: result.id,
          statusId: 1,
        });
      }
      await apolloClient.refetchQueries({ include: [GET_MOVIES_QUERY, ME_ACTIVITIES_QUERY] });
    } catch (error) {
      console.error("Failed to add media:", error);
    } finally {
      setIsAddingMedia(false);
    }
  };

  const handleLogin = async () => {
    await login();
  };

  if (isChecking) {
    return null;
  }

  if (!hasToken) {
    return <AuthPromptPage onLogin={handleLogin} onSignup={handleLogin} />;
  }

  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: any;

          if (route.name === "Home") {
            iconName = focused ? "home" : "home-outline";
          } else if (route.name === "Add") {
            iconName = focused ? "add-circle" : "add-circle-outline";
          } else if (route.name === "History") {
            iconName = focused ? "time" : "time-outline";
          } else if (route.name === "Friends") {
            iconName = focused ? "people" : "people-outline";
          } else if (route.name === "Account") {
            iconName = focused ? "person" : "person-outline";
          }

          return <Ionicons name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: colors.primary,
        tabBarInactiveTintColor: colors["muted-foreground"],
        tabBarStyle: {
          backgroundColor: colors.background,
          borderTopColor: colors.border,
          borderTopWidth: 1,
        },
        headerShown: false,
      })}
    >
      <Tab.Screen name="Home" component={HomePage} />
      <Tab.Screen name="History" component={HistoryPage} />
      <Tab.Screen name="Add">
        {() => <AddMediaPage onBack={() => {}} onAddMedia={handleAddMedia} isLoading={isAddingMedia} />}
      </Tab.Screen>
      <Tab.Screen name="Friends" component={FriendsPage} />
      <Tab.Screen name="Account" component={AccountPage} />
    </Tab.Navigator>
  );
}

export default function App() {
  return (
    <ApolloProvider client={client}>
      <AppContent />
    </ApolloProvider>
  );
}
