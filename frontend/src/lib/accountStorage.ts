import * as SecureStore from 'expo-secure-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { logError, logInfo } from './logger';

const STORED_ACCOUNTS_KEY = 'stored_accounts';
const ACTIVE_ACCOUNT_ID_KEY = 'active_account_id';
const TOKEN_PREFIX = 'account_';
const TOKEN_SUFFIX_ACCESS = '_access_token';
const TOKEN_SUFFIX_REFRESH = '_refresh_token';
const TOKEN_SUFFIX_EXPIRY = '_expiry';

export interface AccountProfile {
  id: string; // user.id from backend
  name: string;
  email: string;
  avatarUrl: string | null;
  authProvider: string;
  authSubject: string;
}

export interface AccountTokens {
  accessToken: string;
  refreshToken?: string;
  expiryTime?: number;
}

/**
 * Get all stored account profiles
 */
export async function getStoredAccounts(): Promise<AccountProfile[]> {
  try {
    const stored = await AsyncStorage.getItem(STORED_ACCOUNTS_KEY);
    if (!stored) return [];
    return JSON.parse(stored);
  } catch (error) {
    logError('Failed to get stored accounts:', error);
    return [];
  }
}

/**
 * Get the currently active account ID
 */
export async function getCurrentAccountId(): Promise<string | null> {
  try {
    return await AsyncStorage.getItem(ACTIVE_ACCOUNT_ID_KEY);
  } catch (error) {
    logError('Failed to get current account ID:', error);
    return null;
  }
}

/**
 * Set the currently active account ID
 */
export async function setCurrentAccountId(accountId: string): Promise<void> {
  try {
    await AsyncStorage.setItem(ACTIVE_ACCOUNT_ID_KEY, accountId);
  } catch (error) {
    logError('Failed to set current account ID:', error);
  }
}

/**
 * Get tokens for a specific account
 */
export async function getAccountTokens(accountId: string): Promise<AccountTokens | null> {
  try {
    const accessToken = await SecureStore.getItemAsync(
      `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_ACCESS}`
    );
    if (!accessToken) return null;

    const refreshToken = await SecureStore.getItemAsync(
      `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_REFRESH}`
    );

    const expiryStr = await SecureStore.getItemAsync(
      `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_EXPIRY}`
    );

    return {
      accessToken,
      refreshToken: refreshToken ?? undefined,
      expiryTime: expiryStr ? parseInt(expiryStr, 10) : undefined,
    };
  } catch (error) {
    logError(`Failed to get tokens for account ${accountId}:`, error);
    return null;
  }
}

/**
 * Save tokens for a specific account
 */
export async function saveAccountTokens(accountId: string, tokens: AccountTokens): Promise<void> {
  try {
    await SecureStore.setItemAsync(
      `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_ACCESS}`,
      tokens.accessToken
    );

    if (tokens.refreshToken) {
      await SecureStore.setItemAsync(
        `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_REFRESH}`,
        tokens.refreshToken
      );
    } else {
      await SecureStore.deleteItemAsync(`${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_REFRESH}`);
    }

    if (tokens.expiryTime) {
      await SecureStore.setItemAsync(
        `${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_EXPIRY}`,
        tokens.expiryTime.toString()
      );
    } else {
      await SecureStore.deleteItemAsync(`${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_EXPIRY}`);
    }
  } catch (error) {
    logError(`Failed to save tokens for account ${accountId}:`, error);
  }
}

/**
 * Add a new account to storage
 */
export async function addAccount(profile: AccountProfile, tokens: AccountTokens): Promise<void> {
  try {
    const accounts = await getStoredAccounts();

    // Check if account already exists
    const existingIndex = accounts.findIndex(a => a.id === profile.id);
    if (existingIndex >= 0) {
      // Update existing account profile
      accounts[existingIndex] = profile;
    } else {
      // Add new account
      accounts.push(profile);
    }

    await AsyncStorage.setItem(STORED_ACCOUNTS_KEY, JSON.stringify(accounts));
    await saveAccountTokens(profile.id, tokens);

    logInfo(`Added account: ${profile.email}`);
  } catch (error) {
    logError('Failed to add account:', error);
  }
}

/**
 * Remove an account and its tokens
 */
export async function removeAccount(accountId: string): Promise<void> {
  try {
    // Remove from account list
    const accounts = await getStoredAccounts();
    const filtered = accounts.filter(a => a.id !== accountId);
    await AsyncStorage.setItem(STORED_ACCOUNTS_KEY, JSON.stringify(filtered));

    // Remove tokens
    await SecureStore.deleteItemAsync(`${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_ACCESS}`);
    await SecureStore.deleteItemAsync(`${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_REFRESH}`);
    await SecureStore.deleteItemAsync(`${TOKEN_PREFIX}${accountId}${TOKEN_SUFFIX_EXPIRY}`);

    // If this was the active account, clear it
    const currentId = await getCurrentAccountId();
    if (currentId === accountId) {
      await AsyncStorage.removeItem(ACTIVE_ACCOUNT_ID_KEY);
    }

    logInfo(`Removed account: ${accountId}`);
  } catch (error) {
    logError('Failed to remove account:', error);
  }
}

/**
 * Get the current active account profile
 */
export async function getCurrentAccount(): Promise<AccountProfile | null> {
  try {
    const accountId = await getCurrentAccountId();
    if (!accountId) return null;

    const accounts = await getStoredAccounts();
    return accounts.find(a => a.id === accountId) ?? null;
  } catch (error) {
    logError('Failed to get current account:', error);
    return null;
  }
}

/**
 * Switch to a different account
 */
export async function switchToAccount(accountId: string): Promise<boolean> {
  try {
    const accounts = await getStoredAccounts();
    const account = accounts.find(a => a.id === accountId);
    if (!account) {
      logError(`Account ${accountId} not found`);
      return false;
    }

    await setCurrentAccountId(accountId);
    logInfo(`Switched to account: ${account.email}`);
    return true;
  } catch (error) {
    logError('Failed to switch account:', error);
    return false;
  }
}

/**
 * Migrate old single-account storage to new multi-account format
 * This should be called once on app startup
 */
export async function migrateOldStorage(userData?: {
  id: string;
  name: string;
  email: string;
  avatarUrl: string | null;
  authProvider: string;
  authSubject: string;
}): Promise<void> {
  try {
    // Check if migration already done
    const accounts = await getStoredAccounts();
    if (accounts.length > 0) return;

    // Check for old token keys
    const oldAccessToken = await SecureStore.getItemAsync('auth0_access_token');
    if (!oldAccessToken || !userData) return;

    logInfo('Migrating old account storage to new format...');

    // Get old tokens
    const oldRefreshToken = await SecureStore.getItemAsync('auth0_refresh_token');
    const oldExpiryStr = await SecureStore.getItemAsync('auth0_token_expiry');

    // Create account profile
    const profile: AccountProfile = {
      id: userData.id,
      name: userData.name,
      email: userData.email,
      avatarUrl: userData.avatarUrl,
      authProvider: userData.authProvider,
      authSubject: userData.authSubject,
    };

    // Save to new format
    const tokens: AccountTokens = {
      accessToken: oldAccessToken,
      refreshToken: oldRefreshToken ?? undefined,
      expiryTime: oldExpiryStr ? parseInt(oldExpiryStr, 10) : undefined,
    };

    await addAccount(profile, tokens);
    await setCurrentAccountId(userData.id);

    // Remove old keys
    await SecureStore.deleteItemAsync('auth0_access_token');
    await SecureStore.deleteItemAsync('auth0_refresh_token');
    await SecureStore.deleteItemAsync('auth0_token_expiry');

    logInfo('Migration completed successfully');
  } catch (error) {
    logError('Failed to migrate old storage:', error);
  }
}
