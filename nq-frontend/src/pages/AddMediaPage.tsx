import React, { useState } from "react";
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Alert,
  TouchableOpacity,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { Button } from "../components/ui/button";
import Input from "../components/ui/input";
import Card from "../components/ui/card";
import { fontSize, spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { Media } from "../types";

interface AddMediaPageProps {
  onBack: () => void;
  onAddMedia: (media: Omit<Media, "id">) => void;
}

const typeOptions = [
  { label: "Movie", value: "movie" as const, icon: "film-outline" as const },
  { label: "TV Show", value: "tv" as const, icon: "tv-outline" as const },
  { label: "Book", value: "book" as const, icon: "book-outline" as const },
  { label: "Music", value: "music" as const, icon: "musical-notes-outline" as const },
  { label: "Game", value: "game" as const, icon: "game-controller-outline" as const },
];

function AddMediaPage({ onBack, onAddMedia }: AddMediaPageProps) {
  const { colors } = useTheme();

  const [title, setTitle] = useState("");
  const [type, setType] = useState<"movie" | "tv" | "book" | "music" | "game">(
    "movie",
  );
  const [year, setYear] = useState(new Date().getFullYear().toString());

  const handleSubmit = () => {
    if (!title.trim()) {
      Alert.alert("Error", "Please enter a title");
      return;
    }

    const newMedia: Omit<Media, "id"> = {
      title: title.trim(),
      type,
      description: "", // Backend will enrich
      year: parseInt(year) || new Date().getFullYear(),
      rating: 0, // Backend will enrich
      duration: undefined, // Backend will enrich
      image: "", // Backend will enrich
      genre: [], // Backend will enrich
    };

    onAddMedia(newMedia);

    // Reset form
    setTitle("");
    setYear(new Date().getFullYear().toString());
  };

  const handleTypeChange = (value: string) => {
    setType(value as "movie" | "tv" | "book" | "music" | "game");
  };

  const styles = StyleSheet.create({
    container: {
      flex: 1,
      backgroundColor: colors.background,
    },
    content: {
      padding: spacing[4],
    },
    backButton: {
      marginBottom: spacing[6],
      alignSelf: "flex-start",
    },
    backText: {
      marginLeft: spacing[2],
      color: colors.foreground,
    },
    header: {
      marginBottom: spacing[6],
    },
    title: {
      fontSize: fontSize.xl,
      fontWeight: "600",
      color: colors.primary,
      marginBottom: spacing[2],
    },
    subtitle: {
      fontSize: fontSize.base,
      color: colors["muted-foreground"],
    },
    card: {
      padding: spacing[6],
    },
    form: {},
    field: {
      marginBottom: spacing[6],
    },
    label: {
      fontSize: fontSize.base,
      fontWeight: "500",
      color: colors.foreground,
      marginBottom: spacing[2],
    },
    typeOptions: {
      flexDirection: "row",
      flexWrap: "wrap",
      gap: spacing[2],
    },
    typeOption: {
      flexDirection: "row",
      alignItems: "center",
      gap: spacing[2],
      padding: spacing[3],
      borderRadius: 8,
      borderWidth: 1,
      borderColor: colors.border,
      backgroundColor: colors.background,
    },
    typeOptionSelected: {
      backgroundColor: colors.primary,
      borderColor: colors.primary,
    },
    typeText: {
      fontSize: fontSize.sm,
      color: colors.foreground,
    },
    typeTextSelected: {
      color: colors["primary-foreground"],
    },
    submitButton: {
      marginTop: spacing[6],
    },
    submitText: {
      color: colors["primary-foreground"],
      marginLeft: spacing[2],
    },
  });

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <View style={styles.header}>
        <Text style={styles.title}>Add New Media</Text>
        <Text style={styles.subtitle}>
          Enter the title and year. We'll fetch the rest of the details for you.
        </Text>
      </View>

      <Card style={styles.card}>
        <View style={styles.form}>
          <View style={styles.field}>
            <Text style={styles.label}>Title *</Text>
            <Input
              value={title}
              onChangeText={setTitle}
              placeholder="Enter media title"
            />
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Type *</Text>
            <View style={styles.typeOptions}>
              {typeOptions.map((option) => (
                <TouchableOpacity
                  key={option.value}
                  style={[
                    styles.typeOption,
                    type === option.value && styles.typeOptionSelected,
                  ]}
                  onPress={() => handleTypeChange(option.value)}
                >
                   <Ionicons
                     name={option.icon}
                     size={20}
                     color={
                       type === option.value
                         ? colors["primary-foreground"]
                         : colors.foreground
                     }
                   />
                  <Text
                    style={[
                      styles.typeText,
                      type === option.value && styles.typeTextSelected,
                    ]}
                  >
                    {option.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Year</Text>
            <Input
              value={year}
              onChangeText={setYear}
              placeholder="2024"
              keyboardType="numeric"
            />
          </View>

          <Button onPress={handleSubmit} style={styles.submitButton}>
            <Ionicons
              name="add"
              size={20}
              color={colors["primary-foreground"]}
            />
            <Text style={styles.submitText}>Add Media</Text>
          </Button>
        </View>
      </Card>
    </ScrollView>
  );
}

export default AddMediaPage;
