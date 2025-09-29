import { Stack } from "expo-router";
import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";
import { ApolloProvider } from "@apollo/client/react";

// Create an HTTP link to your GraphQL server
const httpLink = new HttpLink({
  uri: "http://localhost:8080/query",
});

// Initialize Apollo Client
const client = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          // Define any custom field policies here
        },
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: "cache-and-network",
    },
    query: {
      fetchPolicy: "network-only",
      errorPolicy: "all",
    },
    mutate: {
      errorPolicy: "all",
    },
  },
});

export default function RootLayout() {
  return (
    <ApolloProvider client={client}>
      <Stack />
    </ApolloProvider>
  );
}
