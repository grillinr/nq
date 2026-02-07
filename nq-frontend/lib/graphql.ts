import { gql } from "@apollo/client";

export const ME_QUERY = gql`
  query Me {
    me {
      id
      name
      email
      authProvider
      authSubject
      avatarUrl
    }
  }
`;

export const ME_ACTIVITIES_QUERY = gql`
  query MeActivities {
    me {
      id
      activities {
        id
        media {
          __typename
          id
          title
          description
          releaseDate
          coverUrl
          averageRating
          ... on Movie {
            genres {
              name
            }
            runtime
          }
          ... on TVShow {
            genres {
              name
            }
          }
          ... on Game {
            genre
          }
        }
      }
    }
  }
`;

export const GET_MOVIES_QUERY = gql`
  query GetMovies($limit: Int, $offset: Int) {
    movies(limit: $limit, offset: $offset) {
      id
      title
      coverUrl
      averageRating
      genres {
        name
      }
      description
      releaseDate
      runtime
    }
  }
`;

export const UPDATE_USER_MUTATION = gql`
  mutation UpdateUser($id: UUID!, $input: UpdateUserInput!) {
    updateUser(id: $id, input: $input) {
      id
      name
      email
      avatarUrl
    }
  }
`;
