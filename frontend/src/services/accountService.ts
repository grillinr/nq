import * as AuthSession from 'expo-auth-session';
import * as SecureStore from 'expo-secure-store';
import axios from 'axios';
import { AccountConnection, ServiceType, SERVICE_CONFIGS } from '../types/accounts';

const ACCOUNT_STORAGE_KEY = 'nq_account_connections';

export class AccountService {
  private static instance: AccountService;
  private connections: Map<string, AccountConnection> = new Map();

  static getInstance(): AccountService {
    if (!AccountService.instance) {
      AccountService.instance = new AccountService();
    }
    return AccountService.instance;
  }

  async initialize(): Promise<void> {
    await this.loadConnections();
  }

  async getConnections(): Promise<AccountConnection[]> {
    return Array.from(this.connections.values());
  }

  async getConnection(service: ServiceType): Promise<AccountConnection | null> {
    return this.connections.get(service) || null;
  }

  async connectSpotify(): Promise<AccountConnection | null> {
    try {
      const config = SERVICE_CONFIGS.spotify;
      const redirectUri = AuthSession.makeRedirectUri({
        scheme: 'nq',
        path: 'auth/spotify'
      });

      const authUrl = `${config.authUrl}?client_id=${config.clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${encodeURIComponent(config.scopes!.join(' '))}`;

      const result = await AuthSession.startAsync({
        authUrl,
        returnUrl: redirectUri,
      });

      if (result.type === 'success' && result.params.code) {
        // In a real app, you'd exchange the code for tokens via your backend
        const connection: AccountConnection = {
          id: `spotify_${Date.now()}`,
          service: 'spotify',
          isConnected: true,
          accessToken: 'dummy_token', // Replace with actual token
          lastSync: new Date(),
        };

        await this.saveConnection(connection);
        return connection;
      }
    } catch (error) {
      console.error('Spotify connection failed:', error);
    }
    return null;
  }

  async connectGoogleBooks(): Promise<AccountConnection | null> {
    try {
      const config = SERVICE_CONFIGS.googlebooks;
      const redirectUri = AuthSession.makeRedirectUri({
        scheme: 'nq',
        path: 'auth/googlebooks'
      });

      const authUrl = `${config.authUrl}?client_id=${config.clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${encodeURIComponent(config.scopes!.join(' '))}`;

      const result = await AuthSession.startAsync({
        authUrl,
        returnUrl: redirectUri,
      });

      if (result.type === 'success' && result.params.code) {
        const connection: AccountConnection = {
          id: `googlebooks_${Date.now()}`,
          service: 'googlebooks',
          isConnected: true,
          accessToken: 'dummy_token',
          lastSync: new Date(),
        };

        await this.saveConnection(connection);
        return connection;
      }
    } catch (error) {
      console.error('Google Books connection failed:', error);
    }
    return null;
  }

  async connectYouTube(): Promise<AccountConnection | null> {
    try {
      const config = SERVICE_CONFIGS.youtube;
      const redirectUri = AuthSession.makeRedirectUri({
        scheme: 'nq',
        path: 'auth/youtube'
      });

      const authUrl = `${config.authUrl}?client_id=${config.clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${encodeURIComponent(config.scopes!.join(' '))}`;

      const result = await AuthSession.startAsync({
        authUrl,
        returnUrl: redirectUri,
      });

      if (result.type === 'success' && result.params.code) {
        const connection: AccountConnection = {
          id: `youtube_${Date.now()}`,
          service: 'youtube',
          isConnected: true,
          accessToken: 'dummy_token',
          lastSync: new Date(),
        };

        await this.saveConnection(connection);
        return connection;
      }
    } catch (error) {
      console.error('YouTube connection failed:', error);
    }
    return null;
  }

  async connectManualService(service: ServiceType, credentials: { username?: string; email?: string; apiKey?: string }): Promise<AccountConnection | null> {
    try {
      const connection: AccountConnection = {
        id: `${service}_${Date.now()}`,
        service,
        username: credentials.username,
        email: credentials.email,
        isConnected: true,
        lastSync: new Date(),
      };

      await this.saveConnection(connection);
      return connection;
    } catch (error) {
      console.error(`${service} connection failed:`, error);
      return null;
    }
  }

  async disconnectService(service: ServiceType): Promise<boolean> {
    try {
      const connection = this.connections.get(service);
      if (connection) {
        connection.isConnected = false;
        connection.accessToken = undefined;
        connection.refreshToken = undefined;
        connection.expiresAt = undefined;
        await this.saveConnection(connection);
        return true;
      }
      return false;
    } catch (error) {
      console.error(`Failed to disconnect ${service}:`, error);
      return false;
    }
  }

  async syncService(service: ServiceType): Promise<boolean> {
    try {
      const connection = this.connections.get(service);
      if (connection && connection.isConnected) {
        // In a real app, you'd make API calls to sync data
        connection.lastSync = new Date();
        await this.saveConnection(connection);
        return true;
      }
      return false;
    } catch (error) {
      console.error(`Failed to sync ${service}:`, error);
      return false;
    }
  }

  private async saveConnection(connection: AccountConnection): Promise<void> {
    this.connections.set(connection.service, connection);
    await this.persistConnections();
  }

  private async loadConnections(): Promise<void> {
    try {
      const stored = await SecureStore.getItemAsync(ACCOUNT_STORAGE_KEY);
      if (stored) {
        const connections = JSON.parse(stored);
        this.connections = new Map(Object.entries(connections));
      }
    } catch (error) {
      console.error('Failed to load connections:', error);
    }
  }

  private async persistConnections(): Promise<void> {
    try {
      const connectionsObj = Object.fromEntries(this.connections);
      await SecureStore.setItemAsync(ACCOUNT_STORAGE_KEY, JSON.stringify(connectionsObj));
    } catch (error) {
      console.error('Failed to persist connections:', error);
    }
  }
}

export default AccountService.getInstance();
