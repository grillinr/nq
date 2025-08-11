import React, { useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ActivityIndicator,
} from "react-native";
import { ServiceConfig, AccountConnection } from "../types/accounts";
import AccountService from "../services/accountService";

interface ServiceCardProps {
  config: ServiceConfig;
  connection?: AccountConnection;
  onConnectionChange: () => void;
}

export default function ServiceCard({
  config,
  connection,
  onConnectionChange,
}: ServiceCardProps) {
  const [isConnecting, setIsConnecting] = useState(false);
  const [isDisconnecting, setIsDisconnecting] = useState(false);

  const handleConnect = async () => {
    setIsConnecting(true);
    try {
      let newConnection: AccountConnection | null = null;

      if (config.requiresOAuth) {
        switch (config.name) {
          case "spotify":
            newConnection = await AccountService.connectSpotify();
            break;
          case "googlebooks":
            newConnection = await AccountService.connectGoogleBooks();
            break;
          case "youtube":
            newConnection = await AccountService.connectYouTube();
            break;
          default:
            throw new Error("OAuth not implemented for this service");
        }
      } else {
        // For non-OAuth services, show manual connection dialog
        Alert.prompt(
          `Connect to ${config.displayName}`,
          "Enter your username or email:",
          [
            { text: "Cancel", style: "cancel" },
            {
              text: "Connect",
              onPress: async (username) => {
                if (username) {
                  newConnection = await AccountService.connectManualService(
                    config.name as any,
                    { username }
                  );
                  if (newConnection) {
                    onConnectionChange();
                  }
                }
              },
            },
          ],
          "plain-text"
        );
        setIsConnecting(false);
        return;
      }

      if (newConnection) {
        onConnectionChange();
        Alert.alert(
          "Success",
          `Successfully connected to ${config.displayName}!`
        );
      } else {
        Alert.alert(
          "Error",
          `Failed to connect to ${config.displayName}. Please try again.`
        );
      }
    } catch (error) {
      console.error("Connection error:", error);
      Alert.alert(
        "Error",
        `Failed to connect to ${config.displayName}. Please try again.`
      );
    } finally {
      setIsConnecting(false);
    }
  };

  const handleDisconnect = async () => {
    Alert.alert(
      "Disconnect",
      `Are you sure you want to disconnect from ${config.displayName}?`,
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Disconnect",
          style: "destructive",
          onPress: async () => {
            setIsDisconnecting(true);
            try {
              const success = await AccountService.disconnectService(
                config.name as any
              );
              if (success) {
                onConnectionChange();
                Alert.alert(
                  "Success",
                  `Disconnected from ${config.displayName}`
                );
              } else {
                Alert.alert(
                  "Error",
                  `Failed to disconnect from ${config.displayName}`
                );
              }
            } catch (error) {
              console.error("Disconnection error:", error);
              Alert.alert(
                "Error",
                `Failed to disconnect from ${config.displayName}`
              );
            } finally {
              setIsDisconnecting(false);
            }
          },
        },
      ]
    );
  };

  const handleSync = async () => {
    if (!connection) return;

    try {
      const success = await AccountService.syncService(config.name as any);
      if (success) {
        onConnectionChange();
        Alert.alert("Success", `Synced ${config.displayName} data`);
      } else {
        Alert.alert("Error", `Failed to sync ${config.displayName} data`);
      }
    } catch (error) {
      console.error("Sync error:", error);
      Alert.alert("Error", `Failed to sync ${config.displayName} data`);
    }
  };

  const isConnected = connection?.isConnected || false;

  return (
    <View style={[styles.card, { borderLeftColor: config.color }]}>
      <View style={styles.header}>
        <Text style={styles.icon}>{config.icon}</Text>
        <View style={styles.titleContainer}>
          <Text style={styles.title}>{config.displayName}</Text>
          <Text style={styles.description}>{config.description}</Text>
        </View>
        <View
          style={[
            styles.statusBadge,
            { backgroundColor: isConnected ? "#4CAF50" : "#9E9E9E" },
          ]}
        >
          <Text style={styles.statusText}>
            {isConnected ? "Connected" : "Not Connected"}
          </Text>
        </View>
      </View>

      {isConnected && connection?.lastSync && (
        <View style={styles.syncInfo}>
          <Text style={styles.syncText}>
            Last synced: {connection.lastSync.toLocaleDateString()}
          </Text>
        </View>
      )}

      <View style={styles.actions}>
        {!isConnected ? (
          <TouchableOpacity
            style={[
              styles.button,
              styles.connectButton,
              { backgroundColor: config.color },
            ]}
            onPress={handleConnect}
            disabled={isConnecting}
          >
            {isConnecting ? (
              <ActivityIndicator color="white" size="small" />
            ) : (
              <Text style={styles.buttonText}>Connect</Text>
            )}
          </TouchableOpacity>
        ) : (
          <View style={styles.connectedActions}>
            <TouchableOpacity
              style={[styles.button, styles.syncButton]}
              onPress={handleSync}
            >
              <Text style={styles.buttonText}>Sync</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={[styles.button, styles.disconnectButton]}
              onPress={handleDisconnect}
              disabled={isDisconnecting}
            >
              {isDisconnecting ? (
                <ActivityIndicator color="white" size="small" />
              ) : (
                <Text style={styles.buttonText}>Disconnect</Text>
              )}
            </TouchableOpacity>
          </View>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
    borderLeftWidth: 4,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 12,
  },
  icon: {
    fontSize: 32,
    marginRight: 12,
  },
  titleContainer: {
    flex: 1,
  },
  title: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 4,
  },
  description: {
    fontSize: 14,
    color: "#666",
    lineHeight: 18,
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 12,
  },
  statusText: {
    color: "white",
    fontSize: 12,
    fontWeight: "600",
  },
  syncInfo: {
    marginBottom: 16,
    paddingVertical: 8,
    paddingHorizontal: 12,
    backgroundColor: "#f5f5f5",
    borderRadius: 8,
  },
  syncText: {
    fontSize: 12,
    color: "#666",
  },
  actions: {
    flexDirection: "row",
    justifyContent: "flex-end",
  },
  button: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
    minWidth: 80,
    alignItems: "center",
  },
  connectButton: {
    // backgroundColor set dynamically
  },
  connectedActions: {
    flexDirection: "row",
    gap: 8,
  },
  syncButton: {
    backgroundColor: "#2196F3",
  },
  disconnectButton: {
    backgroundColor: "#F44336",
  },
  buttonText: {
    color: "white",
    fontSize: 14,
    fontWeight: "600",
  },
});
