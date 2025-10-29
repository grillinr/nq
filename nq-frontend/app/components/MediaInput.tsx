import React, { useState } from "react";
import { View, TextInput, Text, StyleSheet, Button, ActivityIndicator, Alert } from "react-native";
import { Picker } from "@react-native-picker/picker";
import { createMedia, MediaType } from "../lib/createMedia";

type Props = {
  onCreated?: (data: any) => void;
};

export default function MediaInput({ onCreated }: Props) {
  const [title, setTitle] = useState("");
  const [mediaType, setMediaType] = useState<MediaType>("movie");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!title.trim()) {
      Alert.alert("Validation", "Title cannot be empty");
      return;
    }
    setLoading(true);
    try {
      const data = await createMedia(mediaType, title);
      setTitle("");
      onCreated?.(data);
      Alert.alert("Success", `Created ${data?.title ?? "item"}`);
    } catch (err: any) {
      Alert.alert("Error", err.message || String(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.root}>
      <Text style={styles.label}>Title</Text>
      <TextInput
        value={title}
        onChangeText={setTitle}
        style={styles.input}
        placeholder="Enter title"
        returnKeyType="done"
      />

      <Text style={[styles.label, { marginTop: 12 }]}>Media type</Text>
      <View style={styles.pickerWrapper}>
        <Picker selectedValue={mediaType} onValueChange={(v) => setMediaType(v as MediaType)} style={styles.picker}>
          <Picker.Item label="Movie" value="movie" />
          <Picker.Item label="TV Show" value="tv" />
          <Picker.Item label="Book" value="book" />
          <Picker.Item label="Game" value="game" />
          <Picker.Item label="Music Album" value="music" />
        </Picker>
      </View>

      <View style={{ marginTop: 14 }}>
        {loading ? (
          <ActivityIndicator />
        ) : (
          <Button title="Create" onPress={handleSubmit} />
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  root: {
    width: "100%",
    maxWidth: 540,
  },
  label: {
    fontSize: 14,
    marginBottom: 6,
  },
  input: {
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 6,
    paddingHorizontal: 12,
    paddingVertical: 10,
    fontSize: 16,
    backgroundColor: "white",
  },
  pickerWrapper: {
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 6,
    overflow: "hidden",
    backgroundColor: "white",
  },
  picker: {
    height: 44,
    width: "100%",
  },
});
