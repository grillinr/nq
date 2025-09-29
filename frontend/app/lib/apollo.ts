import { ApolloClient, HttpLink, InMemoryCache } from '@apollo/client';

const httpLink = new HttpLink({
  uri: 'http://localhost:8080/query', // Update this with your GraphQL endpoint
});

export const client = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
});
