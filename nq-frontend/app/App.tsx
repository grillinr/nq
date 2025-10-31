import React, { useState } from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { Ionicons } from "@expo/vector-icons";
import HomePage from "../src/pages/HomePage";
import AddMediaPage from "../src/pages/AddMediaPage";
import AccountPage from "../src/pages/AccountPage";
import HistoryPage from "../src/pages/HistoryPage";
import FriendsPage from "../src/pages/FriendsPage";
import { createMedia } from "../lib/createMedia";
import { useTheme } from "../src/components/ui/ThemeProvider";

interface Media {
  id: number;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: "movie" | "tv" | "book" | "music" | "game";
}

const mockData: Media[] = [
  {
    id: 1,
    title: "Inception",
    image:
      "https://m.media-amazon.com/images/M/MV5BZjhkNjM0ZTMtNGM5MC00ZTQ3LTk3YmYtZTkzYzdiNWE0ZTA2XkEyXkFqcGc@._V1_.jpg",
    rating: 8.8,
    genre: ["Sci-Fi", "Thriller", "Action"],
    year: 2010,
    duration: "2h 28m",
    description:
      "A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea.",
    type: "movie",
  },
];

const Tab = createBottomTabNavigator();

function AppContent({
  mediaList,
  handleAddMedia,
}: {
  mediaList: Media[];
  handleAddMedia: (m: Omit<Media, "id">) => Promise<void>;
}) {
  const { colors } = useTheme();

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
      <Tab.Screen name="Home">
        {() => <HomePage mediaList={mediaList} />}
      </Tab.Screen>
      <Tab.Screen name="History">
        {() => <HistoryPage mediaList={mediaList} />}
      </Tab.Screen>
      <Tab.Screen name="Add">
        {() => <AddMediaPage onBack={() => {}} onAddMedia={handleAddMedia} />}
      </Tab.Screen>
      <Tab.Screen name="Friends">{() => <FriendsPage />}</Tab.Screen>
      <Tab.Screen name="Account">{() => <AccountPage />}</Tab.Screen>
    </Tab.Navigator>
  );
}

export default function App() {
  const [mediaList, setMediaList] = useState<Media[]>(mockData);

  const handleAddMedia = async (newMedia: Omit<Media, "id">) => {
    try {
      const result = await createMedia(newMedia.type, newMedia.title);
      if (result) {
        const newId = Math.max(...mediaList.map((m) => m.id), 0) + 1;
        setMediaList([...mediaList, { ...newMedia, id: newId }]);
      }
    } catch (error) {
      console.error("Failed to add media:", error);
    }
  };

  return <AppContent mediaList={mediaList} handleAddMedia={handleAddMedia} />;
}
