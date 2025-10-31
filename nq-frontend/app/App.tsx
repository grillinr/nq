import React, { useState } from "react";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { Ionicons } from "@expo/vector-icons";
import HomePage from "../components/HomePage";
import AddMediaPage from "../components/AddMediaPage";
import AccountPage from "../components/AccountPage";
import HistoryPage from "../components/HistoryPage";
import FriendsPage from "../components/FriendsPage";
import { colors } from "../components/ui/tokens";
import { createMedia } from "../lib/createMedia";

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
  // Same as in figma_ui App.tsx
  {
    id: 1,
    title: "Inception",
    image: "https://images.unsplash.com/photo-1524712245354-2c4e5e7121c0?w=400",
    rating: 8.8,
    genre: ["Sci-Fi", "Thriller", "Action"],
    year: 2010,
    duration: "2h 28m",
    description:
      "A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea.",
    type: "movie",
  },
  // Add more as needed
];

const Tab = createBottomTabNavigator();

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

