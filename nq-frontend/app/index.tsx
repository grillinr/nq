import React from "react";
import { View } from "react-native";
import MediaInput from "./components/MediaInput";

export default function Index() {
  return (
    <View
      style={{
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
        padding: 16,
      }}
    >
      <MediaInput
        onCreated={(data) => {
          console.log("Created media:", data);
        }}
      />
    </View>
  );
}
