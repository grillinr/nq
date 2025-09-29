import React, { useState } from "react";
import {
  View,
  Text,
  TextInput,
  Button,
  StyleSheet,
  ScrollView,
  ActivityIndicator,
} from "react-native";
import { Picker } from "@react-native-picker/picker";
import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react/hooks";

type MediaType = "BOOK" | "GAME" | "MOVIE" | "TV_SHOW";

const ADD_MEDIA_ITEM = gql`
  mutation AddMediaItem($input: MediaItemInput!) {
    addMediaItem(input: $input) {
      id
      title
      type
      description
      year
      imageUrl
      url
    }
  }
`;

const AddMediaForm = () => {
  const [mediaType, setMediaType] = useState<MediaType>("BOOK");
  const [title, setTitle] = useState("");
  const [identifier, setIdentifier] = useState(""); // ISBN for books, IGDB ID for games, TMDB ID for videos
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const [addMediaItem] = useMutation(ADD_MEDIA_ITEM, {
    onCompleted: () => {
      setSuccess("Media item added successfully!");
      setTitle("");
      setIdentifier("");
      setTimeout(() => setSuccess(""), 3000);
    },
    onError: (error: Error) => {
      setError(error.message);
      setTimeout(() => setError(""), 5000);
    },
  });

  const fetchMetadata = async () => {
    if (!title && !identifier) {
      setError("Please provide either a title or identifier");
      setTimeout(() => setError(""), 3000);
      return;
    }

    setIsLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await fetch(`http://localhost:8080/api/metadata`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          type: mediaType,
          title: title || undefined,
          id: identifier || undefined,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to fetch metadata");
      }

      const data = await response.json();

      // Add the media item to the database
      await addMediaItem({
        variables: {
          input: {
            ...data,
            type: mediaType,
          },
        },
      });
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to fetch metadata";
      setError(errorMessage);
      setTimeout(() => setError(""), 5000);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>Add New Media</Text>

      <View style={styles.formGroup}>
        <Text style={styles.label}>Media Type</Text>
        <View style={styles.pickerContainer}>
          <Picker
            selectedValue={mediaType}
            onValueChange={(itemValue: MediaType) => setMediaType(itemValue)}
            style={styles.picker}
          >
            <Picker.Item label="Book" value="BOOK" />
            <Picker.Item label="Game" value="GAME" />
            <Picker.Item label="Movie" value="MOVIE" />
            <Picker.Item label="TV Show" value="TV_SHOW" />
          </Picker>
        </View>
      </View>

      <View style={styles.formGroup}>
        <Text style={styles.label}>
          {mediaType === "BOOK"
            ? "Title (optional if ISBN is provided)"
            : "Title (optional if ID is provided)"}
        </Text>
        <TextInput
          style={styles.input}
          value={title}
          onChangeText={setTitle}
          placeholder="Enter title"
        />
      </View>

      <View style={styles.formGroup}>
        <Text style={styles.label}>
          {mediaType === "BOOK"
            ? "ISBN (optional if title is provided)"
            : mediaType === "GAME"
            ? "IGDB ID (optional if title is provided)"
            : "TMDB ID (optional if title is provided)"}
        </Text>
        <TextInput
          style={styles.input}
          value={identifier}
          onChangeText={setIdentifier}
          placeholder={
            mediaType === "BOOK"
              ? "Enter ISBN"
              : mediaType === "GAME"
              ? "Enter IGDB ID"
              : "Enter TMDB ID"
          }
          keyboardType="default"
        />
      </View>

      {error ? <Text style={styles.error}>{error}</Text> : null}
      {success ? <Text style={styles.success}>{success}</Text> : null}

      <View style={styles.buttonContainer}>
        <Button
          title={isLoading ? "Fetching..." : "Add Media"}
          onPress={fetchMetadata}
          disabled={isLoading || (!title && !identifier)}
        />
      </View>

      {isLoading && (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#0000ff" />
          <Text style={styles.loadingText}>Fetching metadata...</Text>
        </View>
      )}
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    backgroundColor: "#fff",
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 20,
    textAlign: "center",
  },
  formGroup: {
    marginBottom: 20,
  },
  label: {
    fontSize: 16,
    marginBottom: 8,
    fontWeight: "500",
  },
  input: {
    borderWidth: 1,
    borderColor: "#ddd",
    padding: 10,
    borderRadius: 4,
    fontSize: 16,
  },
  pickerContainer: {
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 4,
    overflow: "hidden",
  },
  picker: {
    height: 50,
    width: "100%",
  },
  buttonContainer: {
    marginTop: 20,
  },
  error: {
    color: "red",
    marginTop: 10,
    textAlign: "center",
  },
  success: {
    color: "green",
    marginTop: 10,
    textAlign: "center",
  },
  loadingContainer: {
    marginTop: 20,
    alignItems: "center",
  },
  loadingText: {
    marginTop: 10,
    color: "#666",
  },
});

export default AddMediaForm;
