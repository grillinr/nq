import {
  ApolloClient,
  ApolloLink,
  HttpLink,
  InMemoryCache,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { getAccessToken } from "./auth";

const API_URL = process.env.EXPO_PUBLIC_API_URL || "http://localhost:8080/graphql";

const authLink = setContext(async (_, { headers }) => {
  const token = await getAccessToken();
  return {
    headers: {
      ...headers,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  };
});

export const apolloClient = new ApolloClient({
  link: ApolloLink.from([authLink, new HttpLink({ uri: API_URL })]),
  cache: new InMemoryCache(),
});
