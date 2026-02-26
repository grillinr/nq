import {
  ApolloClient,
  ApolloLink,
  HttpLink,
  InMemoryCache,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { onError } from "@apollo/client/link/error";
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

export const apolloClient = new ApolloClient({
  link: ApolloLink.from([errorLink, authLink, new HttpLink({ uri: API_URL })]),
  cache: new InMemoryCache(),
});
