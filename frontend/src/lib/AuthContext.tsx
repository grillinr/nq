import React, {
  createContext,
  useContext,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  ReactNode,
} from 'react';
import { AppState, AppStateStatus } from 'react-native';
import { useApolloClient } from '@apollo/client/react';
import { getAccessToken, loginWithAuth0, logout as logoutFromAuth } from './auth';
import {
  getStoredAccounts,
  getCurrentAccountId,
  setCurrentAccountId,
  switchToAccount,
  addAccount as addAccountToStorage,
  migrateOldStorage,
  type AccountProfile,
} from './accountStorage';
import { clearCacheForAccountSwitch, setAccountOperationInProgress } from './apolloClient';
import { ME_QUERY } from './graphql';
import { logError, logInfo } from './logger';

interface AuthContextType {
  hasToken: boolean;
  isChecking: boolean;
  currentAccountId: string | null;
  storedAccounts: AccountProfile[];
  refreshAuth: () => Promise<void>;
  login: () => Promise<boolean>;
  logout: (accountId?: string) => Promise<void>;
  switchAccount: (accountId: string) => Promise<void>;
  addAccount: () => Promise<boolean>;
  reloadAccounts: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [hasToken, setHasToken] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [currentAccountId, setCurrentAccountIdState] = useState<string | null>(null);
  const [storedAccounts, setStoredAccounts] = useState<AccountProfile[]>([]);
  const isPerformingAccountOperationRef = useRef(false);
  const appState = useRef(AppState.currentState);
  const apolloClient = useApolloClient();

  const beginAccountOperation = useCallback(() => {
    isPerformingAccountOperationRef.current = true;
    setAccountOperationInProgress(true);
  }, []);

  const endAccountOperation = useCallback(() => {
    isPerformingAccountOperationRef.current = false;
    setAccountOperationInProgress(false);
  }, []);

  const reloadAccounts = useCallback(async () => {
    const accounts = await getStoredAccounts();
    const accountId = await getCurrentAccountId();
    setStoredAccounts(accounts);
    setCurrentAccountIdState(accountId);
  }, []);

  const refreshAuth = useCallback(async () => {
    const token = await getAccessToken();
    setHasToken(!!token);
    setIsChecking(false);
    await reloadAccounts();
  }, [reloadAccounts]);

  const login = useCallback(async (): Promise<boolean> => {
    if (isPerformingAccountOperationRef.current) {
      logInfo('[AuthContext] Account operation already in progress, skipping login');
      return false;
    }

    beginAccountOperation();

    try {
      logInfo('[AuthContext] Starting login process');
      const loginResult = await loginWithAuth0(false); // Don't force login for initial login
      if (!loginResult) {
        logError('[AuthContext] Login failed - no login result');
        return false;
      }

      logInfo('[AuthContext] Login successful, fetching user data');

      // Fetch user data using the token directly in the query options
      const userResult = await apolloClient.query({
        query: ME_QUERY,
        fetchPolicy: 'network-only',
        context: {
          headers: {
            Authorization: `Bearer ${loginResult.accessToken}`,
          },
        },
      });

      const userData = (userResult.data as any)?.me;
      if (!userData) {
        logError('[AuthContext] Failed to fetch user data');
        return false;
      }

      logInfo('[AuthContext] User data fetched, storing account');

      // Migrate old storage if needed
      await migrateOldStorage(userData).catch(err =>
        logError('[AuthContext] Migration failed:', err)
      );

      // Store the account
      await addAccountToStorage(
        {
          id: userData.id,
          name: userData.name,
          email: userData.email,
          avatarUrl: userData.avatarUrl,
          authProvider: userData.authProvider,
          authSubject: userData.authSubject,
        },
        {
          accessToken: loginResult.accessToken,
          refreshToken: loginResult.refreshToken,
          expiryTime: loginResult.expiryTime,
        }
      );

      // Set as current account
      await setCurrentAccountId(userData.id);

      logInfo('[AuthContext] Account stored successfully, refreshing auth');

      // Refresh auth state
      await refreshAuth();

      return true;
    } catch (error) {
      logError('[AuthContext] Login error:', error);
      return false;
    } finally {
      endAccountOperation();
    }
  }, [apolloClient, beginAccountOperation, endAccountOperation, refreshAuth]);

  const addAccount = useCallback(async (): Promise<boolean> => {
    if (isPerformingAccountOperationRef.current) {
      logInfo('[AuthContext] Account operation already in progress, skipping add account');
      return false;
    }

    beginAccountOperation();

    try {
      logInfo('[AuthContext] Starting add account process');
      const loginResult = await loginWithAuth0(true); // Force fresh login for new account
      if (!loginResult) {
        logError('[AuthContext] Add account failed - no login result');
        return false;
      }

      logInfo('[AuthContext] Add account login successful, fetching user data');

      // Fetch user data using the token directly in the query options
      const userResult = await apolloClient.query({
        query: ME_QUERY,
        fetchPolicy: 'network-only',
        context: {
          headers: {
            Authorization: `Bearer ${loginResult.accessToken}`,
          },
        },
      });

      const userData = (userResult.data as any)?.me;
      if (!userData) {
        logError('[AuthContext] Failed to fetch user data for new account');
        return false;
      }

      logInfo('[AuthContext] New account user data fetched, storing account');

      // Check if this account already exists
      const existingAccounts = await getStoredAccounts();
      const existingAccount = existingAccounts.find(account => account.id === userData.id);

      if (existingAccount) {
        logInfo('[AuthContext] Account already exists, updating tokens and switching');
        // Update existing account's tokens
        await addAccountToStorage(
          {
            id: userData.id,
            name: userData.name,
            email: userData.email,
            avatarUrl: userData.avatarUrl,
            authProvider: userData.authProvider,
            authSubject: userData.authSubject,
          },
          {
            accessToken: loginResult.accessToken,
            refreshToken: loginResult.refreshToken,
            expiryTime: loginResult.expiryTime,
          }
        );
      } else {
        logInfo('[AuthContext] New account, storing');
        // Store the new account
        await addAccountToStorage(
          {
            id: userData.id,
            name: userData.name,
            email: userData.email,
            avatarUrl: userData.avatarUrl,
            authProvider: userData.authProvider,
            authSubject: userData.authSubject,
          },
          {
            accessToken: loginResult.accessToken,
            refreshToken: loginResult.refreshToken,
            expiryTime: loginResult.expiryTime,
          }
        );
      }

      // Set as current account
      await setCurrentAccountId(userData.id);

      logInfo('[AuthContext] New account stored successfully, refreshing auth');

      // Clear cache to prevent data leakage between accounts
      await clearCacheForAccountSwitch().catch(cacheError => {
        logError('[AuthContext] Cache clear failed during add account:', cacheError);
      });

      // Refresh auth state
      await refreshAuth();

      return true;
    } catch (error) {
      logError('[AuthContext] Add account error:', error);
      return false;
    } finally {
      endAccountOperation();
    }
  }, [apolloClient, beginAccountOperation, endAccountOperation, refreshAuth]);

  const logout = useCallback(
    async (accountId?: string) => {
      await logoutFromAuth(accountId);

      // If we removed the current account, switch to another one if available
      const currentId = await getCurrentAccountId();
      if (!currentId) {
        // No current account set, try to switch to first available
        const accounts = await getStoredAccounts();
        if (accounts.length > 0) {
          await setCurrentAccountId(accounts[0].id);
        }
      }

      await refreshAuth();
    },
    [refreshAuth]
  );

  const switchAccount = useCallback(
    async (accountId: string) => {
      if (isPerformingAccountOperationRef.current) {
        logInfo('[AuthContext] Account operation already in progress, skipping switch');
        return;
      }

      beginAccountOperation();

      try {
        const success = await switchToAccount(accountId);
        if (success) {
          await clearCacheForAccountSwitch().catch(cacheError => {
            logError('[AuthContext] Cache clear failed during switch account:', cacheError);
          });
          await refreshAuth();
        }
      } finally {
        endAccountOperation();
      }
    },
    [beginAccountOperation, endAccountOperation, refreshAuth]
  );

  // Initial token check on mount
  useEffect(() => {
    refreshAuth();
  }, [refreshAuth]);

  // Re-check auth whenever the app comes back to the foreground.
  // This handles the Expo Go deep-link case where the browser redirects back
  // and the component tree may have remounted before loginWithAuth0() returned.
  useEffect(() => {
    const subscription = AppState.addEventListener('change', (nextState: AppStateStatus) => {
      if (appState.current.match(/inactive|background/) && nextState === 'active') {
        // Skip refresh if we're in the middle of an account operation
        if (!isPerformingAccountOperationRef.current) {
          refreshAuth();
        } else {
          logInfo('[AuthContext] Skipping app state refresh during account operation');
        }
      }
      appState.current = nextState;
    });
    return () => subscription.remove();
  }, [refreshAuth]);

  const contextValue = useMemo(
    () => ({
      hasToken,
      isChecking,
      currentAccountId,
      storedAccounts,
      refreshAuth,
      login,
      logout,
      switchAccount,
      addAccount,
      reloadAccounts,
    }),
    [
      hasToken,
      isChecking,
      currentAccountId,
      storedAccounts,
      refreshAuth,
      login,
      logout,
      switchAccount,
      addAccount,
      reloadAccounts,
    ]
  );

  return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

// Export types for use in other files
export type { AccountProfile };
