import { ApolloClient, ApolloLink, HttpLink, InMemoryCache } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { persistCache, AsyncStorageWrapper } from 'apollo3-cache-persist';
import { getAccessToken } from './auth';
import { getAccountTokens, getCurrentAccountId, removeAccount } from './accountStorage';
import { logError, logInfo } from './logger';

// Use environment variable for API URL, with a fallback to localhost
const API_URL = process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/graphql';

let accountOperationInProgress = false;

export function setAccountOperationInProgress(inProgress: boolean): void {
  accountOperationInProgress = inProgress;
}

const authLink = setContext(async (_, { headers }) => {
  const hasAuthorizationHeader =
    typeof headers?.Authorization === 'string' || typeof headers?.authorization === 'string';

  if (hasAuthorizationHeader) {
    return { headers };
  }

  const token = await getAccessToken();
  return {
    headers: {
      ...headers,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  };
});

// Handle authentication errors
const errorLink = onError((errorOptions: any) => {
  if (
    errorOptions.networkError &&
    'statusCode' in errorOptions.networkError &&
    errorOptions.networkError.statusCode === 401
  ) {
    if (accountOperationInProgress) {
      logInfo('Skipping 401 account cleanup during account operation');
      return undefined;
    }

    logInfo('Authentication error - removing invalid account');

    const failedAuthHeader =
      errorOptions.operation.getContext()?.headers?.Authorization ??
      errorOptions.operation.getContext()?.headers?.authorization;
    const failedToken =
      typeof failedAuthHeader === 'string' && failedAuthHeader.startsWith('Bearer ')
        ? failedAuthHeader.slice('Bearer '.length)
        : null;

    // Remove the account only if the 401 came from the current account token.
    // This avoids deleting the newly selected account when older in-flight
    // requests fail during account switching.
    getCurrentAccountId()
      .then(async accountId => {
        if (!accountId || !failedToken) {
          return null;
        }

        const currentTokens = await getAccountTokens(accountId);
        if (currentTokens?.accessToken !== failedToken) {
          logInfo('Skipping account removal: 401 came from a non-current token');
          return null;
        }

        return removeAccount(accountId);
      })
      .catch(logError);

    // Retry the request without the token
    errorOptions.operation.setContext({
      headers: {
        ...errorOptions.operation.getContext().headers,
        Authorization: undefined,
      },
    });

    return errorOptions.forward(errorOptions.operation);
  }
  return undefined;
});

export const cache = new InMemoryCache({
  // Normalise all five media union types by their `id` field so that the
  // same item fetched in different queries (e.g. allMedia and me.activities)
  // shares a single cache entry. Without this Apollo falls back to
  // root-level keying, storing duplicates and missing cross-query updates.
  typePolicies: {
    Movie: { keyFields: ['id'] },
    TVShow: { keyFields: ['id'] },
    Book: { keyFields: ['id'] },
    Game: { keyFields: ['id'] },
    MusicAlbum: { keyFields: ['id'] },
    User: { keyFields: ['id'] },
  },
});

export const apolloClient = new ApolloClient({
  link: ApolloLink.from([errorLink, authLink, new HttpLink({ uri: API_URL })]),
  cache,
});

/**
 * Restores the Apollo cache from AsyncStorage before the app renders its first
 * frame. Call this once at the root of the app and wait for it to resolve
 * before mounting ApolloProvider, so the first render already has cached data.
 */
export async function initApolloCache(): Promise<void> {
  await persistCache({
    cache,
    storage: new AsyncStorageWrapper(AsyncStorage),
  });
}

// Prevent multiple cache clearing operations from running simultaneously
let clearingCache = false;

/**
 * Clear the Apollo cache when switching accounts to prevent data leakage
 * Uses a more robust approach to avoid clearStore() invariant errors
 */
export async function clearCacheForAccountSwitch(): Promise<void> {
  if (clearingCache) {
    logInfo('Cache clearing already in progress, skipping...');
    return;
  }

  clearingCache = true;

  try {
    logInfo('Clearing Apollo cache for account switch');

    cache.reset();
    await AsyncStorage.removeItem('apollo-cache-persist');

    logInfo('Apollo cache cleared for account switch');
  } catch (error) {
    logError('Failed to clear Apollo cache for account switch:', error);
  } finally {
    clearingCache = false;
  }
}
