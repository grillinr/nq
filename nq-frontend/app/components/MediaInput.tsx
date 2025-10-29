import React, { useState } from "react";
import {
  View,
  TextInput,
  Text,
  StyleSheet,
  Button,
  ActivityIndicator,
  Alert,
  Platform,
  TouchableOpacity,
  Modal,
  FlatList,
} from "react-native";
import { ActionSheetIOS } from "react-native";
import { createMedia, MediaType } from "../lib/createMedia";

type Props = {
  onCreated?: (data: any) => void;
};

const OPTIONS: { label: string; value: MediaType }[] = [
  { label: "Movie", value: "movie" },
  { label: "TV Show", value: "tv" },
  { label: "Book", value: "book" },
  { label: "Game", value: "game" },
  { label: "Music Album", value: "music" },
];

export default function MediaInput({ onCreated }: Props) {
  const [title, setTitle] = useState("");
  const [mediaType, setMediaType] = useState<MediaType>("movie");
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);

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

  const openSelector = () => {
    if (Platform.OS === "ios") {
      ActionSheetIOS.showActionSheetWithOptions(
        {
          options: OPTIONS.map((o) => o.label).concat(["Cancel"]),
          cancelButtonIndex: OPTIONS.length,
        },
        (buttonIndex) => {
          if (buttonIndex >= 0 && buttonIndex < OPTIONS.length) {
            setMediaType(OPTIONS[buttonIndex].value);
          }
        }
      );
    } else {
      setModalVisible(true);
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

      <TouchableOpacity onPress={openSelector} style={styles.selectorButton}>
        <Text style={styles.selectorText}>{OPTIONS.find((o) => o.value === mediaType)?.label}</Text>
      </TouchableOpacity>

      <Modal visible={modalVisible} animationType="slide" transparent={true} onRequestClose={() => setModalVisible(false)}>
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <FlatList
              data={OPTIONS}
              keyExtractor={(item) => item.value}
              renderItem={({ item }) => (
                <TouchableOpacity
                  style={styles.modalItem}
                  onPress={() => {
                    setMediaType(item.value);
                    setModalVisible(false);
                  }}
                >
                  <Text style={styles.modalItemText}>{item.label}</Text>
                </TouchableOpacity>
              )}
            />
            <Button title="Cancel" onPress={() => setModalVisible(false)} />
          </View>
        </View>
      </Modal>

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
  selectorButton: {
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 6,
    paddingHorizontal: 12,
    paddingVertical: 10,
    backgroundColor: "white",
  },
  selectorText: {
    fontSize: 16,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.4)",
    justifyContent: "flex-end",
  },
  modalContent: {
    backgroundColor: "white",
    padding: 16,
    borderTopLeftRadius: 12,
    borderTopRightRadius: 12,
    maxHeight: "50%",
  },
  modalItem: {
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#eee",
  },
  modalItemText: {
    fontSize: 16,
  },
});
