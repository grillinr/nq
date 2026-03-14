import {
  ApolloClient,
  ApolloLink,
  HttpLink,
  InMemoryCache,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { onError } from "@apollo/client/link/error";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { persistCache, AsyncStorageWrapper } from "apollo3-cache-persist";
import { getAccessToken, logout } from "./auth";

// Use environment variable for API URL, with a fallback to localhost
const API_URL = process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080/graphql";

const authLink = setContext(async (_, { headers }) => {
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
  if (errorOptions.networkError && 'statusCode' in errorOptions.networkError && errorOptions.networkError.statusCode === 401) {
    console.warn("Authentication error - clearing invalid token");
    // Clear the invalid token
    logout().catch(console.error);
    
    // Retry the request without the token
    errorOptions.operation.setContext({
      headers: {
        ...errorOptions.operation.getContext().headers,
        Authorization: undefined,
      },
    });
    
    return errorOptions.forward(errorOptions.operation);
  }
});

export const cache = new InMemoryCache({
  // Normalise all five media union types by their `id` field so that the
  // same item fetched in different queries (e.g. allMedia and me.activities)
  // shares a single cache entry. Without this Apollo falls back to
  // root-level keying, storing duplicates and missing cross-query updates.
  typePolicies: {
    Movie:      { keyFields: ["id"] },
    TVShow:     { keyFields: ["id"] },
    Book:       { keyFields: ["id"] },
    Game:       { keyFields: ["id"] },
    MusicAlbum: { keyFields: ["id"] },
    User:       { keyFields: ["id"] },
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
