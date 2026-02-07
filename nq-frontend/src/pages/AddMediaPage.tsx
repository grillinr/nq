import React, { useEffect, useMemo, useState } from "react";
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Alert,
  TouchableOpacity,
  ActivityIndicator,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { Button } from "../components/ui/button";
import Input from "../components/ui/input";
import Card from "../components/ui/card";
import { fontSize, spacing } from "../components/ui/tokens";
import { useTheme } from "../components/ui/ThemeProvider";
import { Media } from "../types";
import { useApolloClient } from "@apollo/client/react";
import { AUTOCOMPLETE_MEDIA_QUERY } from "../../lib/graphql";

interface AddMediaPageProps {
  onBack: () => void;
  onAddMedia: (media: Omit<Media, "id">) => void;
  isLoading?: boolean;
}

const typeOptions = [
  { label: "Movie", value: "movie" as const, icon: "film-outline" as const },
  { label: "TV Show", value: "tv" as const, icon: "tv-outline" as const },
  { label: "Book", value: "book" as const, icon: "book-outline" as const },
  {
    label: "Music",
    value: "music" as const,
    icon: "musical-notes-outline" as const,
  },
  {
    label: "Game",
    value: "game" as const,
    icon: "game-controller-outline" as const,
  },
];

type MediaType = "movie" | "tv" | "book" | "music" | "game";

type MediaSuggestion = {
  title: string;
  year?: number | null;
  externalId?: string | null;
  imageUrl?: string | null;
  subtitle?: string | null;
};

function AddMediaPage({
  onBack,
  onAddMedia,
  isLoading = false,
}: AddMediaPageProps) {
  const { colors } = useTheme();
  const apolloClient = useApolloClient();

  const [title, setTitle] = useState("");
  const [type, setType] = useState<MediaType | null>(null);
  const [year, setYear] = useState("");
  const [suggestions, setSuggestions] = useState<MediaSuggestion[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedExternalId, setSelectedExternalId] = useState<
    string | undefined
  >();
  const [selectedIsbn, setSelectedIsbn] = useState<string | undefined>();
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [suppressAutocomplete, setSuppressAutocomplete] = useState(false);

  const handleSubmit = () => {
    if (isLoading) return;
    if (!title.trim()) {
      Alert.alert("Error", "Please enter a title");
      return;
    }
    if (!type) {
      Alert.alert("Error", "Please select a media type");
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
      externalId: selectedExternalId,
      isbn: selectedIsbn,
    };

    onAddMedia(newMedia);

    // Reset form
    setTitle("");
    setYear("");
    setSelectedExternalId(undefined);
    setSelectedIsbn(undefined);
    setSuggestions([]);
    setShowSuggestions(false);
    setSuppressAutocomplete(false);
  };

  const handleTypeChange = (value: string) => {
    setType(value as MediaType);
    setSelectedExternalId(undefined);
    setSelectedIsbn(undefined);
    setSuggestions([]);
    setShowSuggestions(false);
    setSuppressAutocomplete(false);
  };

  const typeLabel = useMemo(
    () => typeOptions.find((option) => option.value === type)?.label,
    [type],
  );
  const canType = Boolean(type);
  const showMusicNotice = type === "music";

  useEffect(() => {
    if (suppressAutocomplete) {
      setShowSuggestions(false);
      return;
    }
    if (!type || !title.trim() || showMusicNotice) {
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }

    setIsSearching(true);
    const handle = setTimeout(async () => {
      try {
        const { data } = await apolloClient.query({
          query: AUTOCOMPLETE_MEDIA_QUERY,
          variables: {
            type: type.toUpperCase(),
            query: title.trim(),
          },
          fetchPolicy: "no-cache",
        });
        setSuggestions(data?.autocompleteMedia ?? []);
        setShowSuggestions(true);
      } catch {
        setSuggestions([]);
        setShowSuggestions(false);
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(handle);
  }, [apolloClient, showMusicNotice, suppressAutocomplete, title, type]);

  const handleSuggestionPress = (item: MediaSuggestion) => {
    setTitle(item.title);
    if (item.year) {
      setYear(String(item.year));
    }
    if (type === "book") {
      setSelectedIsbn(item.externalId || undefined);
      setSelectedExternalId(undefined);
    } else {
      setSelectedExternalId(item.externalId || undefined);
      setSelectedIsbn(undefined);
    }
    setShowSuggestions(false);
    setSuppressAutocomplete(true);
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
    helperText: {
      color: colors["muted-foreground"],
      marginTop: spacing[2],
      fontSize: fontSize.sm,
    },
    suggestions: {
      marginTop: spacing[2],
      borderWidth: 1,
      borderColor: colors.border,
      borderRadius: 8,
      backgroundColor: colors.background,
    },
    suggestionItem: {
      paddingVertical: spacing[3],
      paddingHorizontal: spacing[3],
      borderBottomWidth: 1,
      borderBottomColor: colors.border,
    },
    suggestionTitle: {
      fontSize: fontSize.base,
      color: colors.foreground,
    },
    suggestionSubtitle: {
      fontSize: fontSize.sm,
      color: colors["muted-foreground"],
      marginTop: spacing[1],
    },
    suggestionRow: {
      flexDirection: "row",
      alignItems: "center",
      justifyContent: "space-between",
      gap: spacing[2],
    },
    suggestionYear: {
      fontSize: fontSize.sm,
      color: colors["muted-foreground"],
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
          Enter the title and year. We&apos;ll fetch the rest of the details for
          you.
        </Text>
      </View>

      <Card style={styles.card}>
        <View style={styles.form}>
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
            {!type ? (
              <Text style={styles.helperText}>
                Select a type to enable title input.
              </Text>
            ) : null}
          </View>

          {canType ? (
            <View style={styles.field}>
              <Text style={styles.label}>Title *</Text>
              <Input
                value={title}
                onChangeText={(value) => {
                  setTitle(value);
                  setSelectedExternalId(undefined);
                  setSelectedIsbn(undefined);
                  setSuppressAutocomplete(false);
                }}
                placeholder={`Enter ${typeLabel ?? "media"} title`}
              />
              {showMusicNotice ? (
                <Text style={styles.helperText}>
                  Music autocomplete coming soon.
                </Text>
              ) : null}
              {isSearching && title.trim() && !showMusicNotice ? (
                <Text style={styles.helperText}>Searching…</Text>
              ) : null}
              {showSuggestions && suggestions.length > 0 ? (
                <View style={styles.suggestions}>
                  {suggestions.map((item, index) => (
                    <TouchableOpacity
                      key={`${item.title}-${index}`}
                      style={[
                        styles.suggestionItem,
                        index === suggestions.length - 1
                          ? { borderBottomWidth: 0 }
                          : null,
                      ]}
                      onPress={() => handleSuggestionPress(item)}
                    >
                      <View style={styles.suggestionRow}>
                        <Text style={styles.suggestionTitle}>{item.title}</Text>
                        {item.year ? (
                          <Text style={styles.suggestionYear}>{item.year}</Text>
                        ) : null}
                      </View>
                      {item.subtitle ? (
                        <Text style={styles.suggestionSubtitle}>
                          {item.subtitle}
                        </Text>
                      ) : null}
                    </TouchableOpacity>
                  ))}
                </View>
              ) : null}
            </View>
          ) : null}

          <View style={styles.field}>
            <Text style={styles.label}>Year</Text>
            <Input
              value={year}
              onChangeText={setYear}
              placeholder="Release year"
              keyboardType="numeric"
            />
          </View>

          <Button
            onPress={handleSubmit}
            style={styles.submitButton}
            disabled={isLoading}
          >
            {isLoading ? (
              <ActivityIndicator color={colors["primary-foreground"]} />
            ) : (
              <>
                <Ionicons
                  name="add"
                  size={20}
                  color={colors["primary-foreground"]}
                />
                <Text style={styles.submitText}>Add Media</Text>
              </>
            )}
          </Button>
        </View>
      </Card>
    </ScrollView>
  );
}

export default AddMediaPage;
