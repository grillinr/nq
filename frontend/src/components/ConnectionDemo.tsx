import React from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
} from "react-native";
import { AccountConnection } from "../types/accounts";

interface ConnectionDemoProps {
  connections: AccountConnection[];
}

export default function ConnectionDemo({ connections }: ConnectionDemoProps) {
  const connectedServices = connections.filter((conn) => conn.isConnected);

  if (connectedServices.length === 0) {
    return (
      <View style={styles.emptyState}>
        <Text style={styles.emptyTitle}>No Connected Services</Text>
        <Text style={styles.emptyText}>
          Connect your accounts to see your data here
        </Text>
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      horizontal
      showsHorizontalScrollIndicator={false}
    >
      {connectedServices.map((connection) => (
        <View key={connection.id} style={styles.demoCard}>
          <Text style={styles.serviceName}>
            {getServiceDisplayName(connection.service)}
          </Text>
          <View style={styles.demoContent}>
            {renderDemoContent(connection.service)}
          </View>
          <Text style={styles.lastSync}>
            Last synced: {connection.lastSync?.toLocaleDateString() || "Never"}
          </Text>
        </View>
      ))}
    </ScrollView>
  );
}

function getServiceDisplayName(service: string): string {
  const names: Record<string, string> = {
    spotify: "Spotify",
    goodreads: "Goodreads",
    letterboxd: "Letterboxd",
    googlebooks: "Google Books",
    youtube: "YouTube",
    steam: "Steam",
    libby: "Libby",
    pocketcasts: "Pocket Casts",
  };
  return names[service] || service;
}

function renderDemoContent(service: string): React.ReactNode {
  switch (service) {
    case "spotify":
      return (
        <>
          <Text style={styles.demoTitle}>Recently Played</Text>
          <Text style={styles.demoItem}>🎵 "Bohemian Rhapsody" - Queen</Text>
          <Text style={styles.demoItem}>🎵 "Hotel California" - Eagles</Text>
          <Text style={styles.demoItem}>
            🎵 "Stairway to Heaven" - Led Zeppelin
          </Text>
        </>
      );

    case "goodreads":
      return (
        <>
          <Text style={styles.demoTitle}>Currently Reading</Text>
          <Text style={styles.demoItem}>
            📚 "The Great Gatsby" - F. Scott Fitzgerald
          </Text>
          <Text style={styles.demoItem}>📚 "1984" - George Orwell</Text>
          <Text style={styles.demoItem}>
            📚 "Pride and Prejudice" - Jane Austen
          </Text>
        </>
      );

    case "letterboxd":
      return (
        <>
          <Text style={styles.demoTitle}>Recently Watched</Text>
          <Text style={styles.demoItem}>🎬 "Inception" (2010) ⭐⭐⭐⭐⭐</Text>
          <Text style={styles.demoItem}>
            🎬 "The Shawshank Redemption" (1994) ⭐⭐⭐⭐⭐
          </Text>
          <Text style={styles.demoItem}>
            🎬 "Pulp Fiction" (1994) ⭐⭐⭐⭐⭐
          </Text>
        </>
      );

    case "googlebooks":
      return (
        <>
          <Text style={styles.demoTitle}>Reading Progress</Text>
          <Text style={styles.demoItem}>📖 "Dune" - 45% complete</Text>
          <Text style={styles.demoItem}>📖 "The Hobbit" - 78% complete</Text>
          <Text style={styles.demoItem}>📖 "Neuromancer" - 12% complete</Text>
        </>
      );

    case "youtube":
      return (
        <>
          <Text style={styles.demoTitle}>Subscriptions</Text>
          <Text style={styles.demoItem}>📺 "Veritasium" - Science videos</Text>
          <Text style={styles.demoItem}>
            📺 "Kurzgesagt" - Educational content
          </Text>
          <Text style={styles.demoItem}>
            📺 "CGP Grey" - History & politics
          </Text>
        </>
      );

    case "steam":
      return (
        <>
          <Text style={styles.demoTitle}>Recent Games</Text>
          <Text style={styles.demoItem}>🎮 "Portal 2" - 2.5 hours</Text>
          <Text style={styles.demoItem}>🎮 "Half-Life 2" - 15.3 hours</Text>
          <Text style={styles.demoItem}>
            🎮 "Team Fortress 2" - 127.8 hours
          </Text>
        </>
      );

    case "libby":
      return (
        <>
          <Text style={styles.demoTitle}>Library Books</Text>
          <Text style={styles.demoItem}>🏛️ "The Martian" - Available</Text>
          <Text style={styles.demoItem}>
            🏛️ "Ready Player One" - Checked out
          </Text>
          <Text style={styles.demoItem}>
            🏛️ "The Name of the Wind" - On hold
          </Text>
        </>
      );

    case "pocketcasts":
      return (
        <>
          <Text style={styles.demoTitle}>Podcast Subscriptions</Text>
          <Text style={styles.demoItem}>
            🎧 "This American Life" - 3 new episodes
          </Text>
          <Text style={styles.demoItem}>🎧 "Radiolab" - 1 new episode</Text>
          <Text style={styles.demoItem}>🎧 "Serial" - 2 new episodes</Text>
        </>
      );

    default:
      return <Text style={styles.demoText}>Data from {service}</Text>;
  }
}

const styles = StyleSheet.create({
  container: {
    flexGrow: 0,
  },
  emptyState: {
    padding: 20,
    alignItems: "center",
  },
  emptyTitle: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#666",
    marginBottom: 8,
  },
  emptyText: {
    fontSize: 14,
    color: "#999",
    textAlign: "center",
  },
  demoCard: {
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 16,
    marginRight: 16,
    width: 280,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  serviceName: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 12,
    textAlign: "center",
  },
  demoContent: {
    marginBottom: 12,
  },
  demoTitle: {
    fontSize: 14,
    fontWeight: "600",
    color: "#1a1a1a",
    marginBottom: 8,
  },
  demoItem: {
    fontSize: 12,
    color: "#666",
    marginBottom: 4,
    lineHeight: 16,
  },
  demoText: {
    fontSize: 12,
    color: "#666",
  },
  lastSync: {
    fontSize: 10,
    color: "#999",
    textAlign: "center",
    fontStyle: "italic",
  },
});
