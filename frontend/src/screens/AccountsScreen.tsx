import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
  SafeAreaView,
} from "react-native";
import ServiceCard from "../components/ServiceCard";
import {
  SERVICE_CONFIGS,
  ServiceType,
  AccountConnection,
} from "../types/accounts";
import AccountService from "../services/accountService";

export default function AccountsScreen() {
  const [connections, setConnections] = useState<AccountConnection[]>([]);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    loadConnections();
  }, []);

  const loadConnections = async () => {
    try {
      await AccountService.initialize();
      const loadedConnections = await AccountService.getConnections();
      setConnections(loadedConnections);
    } catch (error) {
      console.error("Failed to load connections:", error);
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await loadConnections();
    setRefreshing(false);
  };

  const handleConnectionChange = () => {
    loadConnections();
  };

  const getConnectionForService = (
    service: ServiceType
  ): AccountConnection | undefined => {
    return connections.find((conn) => conn.service === service);
  };

  const getConnectedCount = (): number => {
    return connections.filter((conn) => conn.isConnected).length;
  };

  const getTotalCount = (): number => {
    return Object.keys(SERVICE_CONFIGS).length;
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContainer}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        <View style={styles.header}>
          <Text style={styles.title}>Account Connections</Text>
          <Text style={styles.subtitle}>
            Connect your accounts to sync data across platforms
          </Text>
          <View style={styles.statsContainer}>
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{getConnectedCount()}</Text>
              <Text style={styles.statLabel}>Connected</Text>
            </View>
            <View style={styles.statDivider} />
            <View style={styles.statItem}>
              <Text style={styles.statNumber}>{getTotalCount()}</Text>
              <Text style={styles.statLabel}>Total</Text>
            </View>
          </View>
        </View>

        <View style={styles.servicesSection}>
          <Text style={styles.sectionTitle}>Available Services</Text>

          {/* Music & Audio */}
          <View style={styles.categorySection}>
            <Text style={styles.categoryTitle}>🎵 Music & Audio</Text>
            <ServiceCard
              config={SERVICE_CONFIGS.spotify}
              connection={getConnectionForService("spotify")}
              onConnectionChange={handleConnectionChange}
            />
          </View>

          {/* Books & Reading */}
          <View style={styles.categorySection}>
            <Text style={styles.categoryTitle}>📚 Books & Reading</Text>
            <ServiceCard
              config={SERVICE_CONFIGS.goodreads}
              connection={getConnectionForService("goodreads")}
              onConnectionChange={handleConnectionChange}
            />
            <ServiceCard
              config={SERVICE_CONFIGS.googlebooks}
              connection={getConnectionForService("googlebooks")}
              onConnectionChange={handleConnectionChange}
            />
            <ServiceCard
              config={SERVICE_CONFIGS.libby}
              connection={getConnectionForService("libby")}
              onConnectionChange={handleConnectionChange}
            />
          </View>

          {/* Video & Entertainment */}
          <View style={styles.categorySection}>
            <Text style={styles.categoryTitle}>🎬 Video & Entertainment</Text>
            <ServiceCard
              config={SERVICE_CONFIGS.letterboxd}
              connection={getConnectionForService("letterboxd")}
              onConnectionChange={handleConnectionChange}
            />
            <ServiceCard
              config={SERVICE_CONFIGS.youtube}
              connection={getConnectionForService("youtube")}
              onConnectionChange={handleConnectionChange}
            />
          </View>

          {/* Gaming */}
          <View style={styles.categorySection}>
            <Text style={styles.categoryTitle}>🎮 Gaming</Text>
            <ServiceCard
              config={SERVICE_CONFIGS.steam}
              connection={getConnectionForService("steam")}
              onConnectionChange={handleConnectionChange}
            />
          </View>

          {/* Podcasts */}
          <View style={styles.categorySection}>
            <Text style={styles.categoryTitle}>🎧 Podcasts</Text>
            <ServiceCard
              config={SERVICE_CONFIGS.pocketcasts}
              connection={getConnectionForService("pocketcasts")}
              onConnectionChange={handleConnectionChange}
            />
          </View>
        </View>

        <View style={styles.infoSection}>
          <Text style={styles.infoTitle}>How it works</Text>
          <Text style={styles.infoText}>
            • Connect your accounts to automatically sync your data{"\n"}• OAuth
            services (Spotify, Google Books, YouTube) use secure authentication
            {"\n"}• Manual services require your username or API key{"\n"}• Data
            is synced locally and securely stored on your device{"\n"}• You can
            disconnect any service at any time
          </Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#f5f5f5",
  },
  scrollView: {
    flex: 1,
  },
  scrollContainer: {
    padding: 20,
  },
  header: {
    alignItems: "center",
    marginBottom: 30,
    marginTop: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 8,
    textAlign: "center",
  },
  subtitle: {
    fontSize: 16,
    color: "#666",
    textAlign: "center",
    marginBottom: 20,
    lineHeight: 22,
  },
  statsContainer: {
    flexDirection: "row",
    backgroundColor: "#fff",
    borderRadius: 16,
    padding: 20,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  statItem: {
    alignItems: "center",
    flex: 1,
  },
  statNumber: {
    fontSize: 28,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 14,
    color: "#666",
    textTransform: "uppercase",
    letterSpacing: 0.5,
  },
  statDivider: {
    width: 1,
    backgroundColor: "#e0e0e0",
    marginHorizontal: 20,
  },
  servicesSection: {
    marginBottom: 30,
  },
  sectionTitle: {
    fontSize: 24,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 20,
  },
  categorySection: {
    marginBottom: 24,
  },
  categoryTitle: {
    fontSize: 18,
    fontWeight: "600",
    color: "#1a1a1a",
    marginBottom: 12,
    marginLeft: 4,
  },
  infoSection: {
    backgroundColor: "#fff",
    padding: 20,
    borderRadius: 12,
    marginBottom: 20,
  },
  infoTitle: {
    fontSize: 18,
    fontWeight: "bold",
    color: "#1a1a1a",
    marginBottom: 12,
  },
  infoText: {
    fontSize: 14,
    color: "#666",
    lineHeight: 20,
  },
});
