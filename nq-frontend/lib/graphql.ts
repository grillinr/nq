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
          ... on Book {
            subjects {
              name
            }
            pages
          }
          ... on Game {
            genre
            themes
            keywords
            gameModes
            perspectives
            franchises
            platformsList
          }
        }
      }
    }
  }
`;

export const GET_MOVIES_QUERY = gql`
  query GetMovies($limit: Int, $offset: Int) {
    movies(limit: $limit, offset: $offset) {
      __typename
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

export const GET_HOME_MEDIA_QUERY = gql`
  query GetHomeMedia {
    allMedia {
      __typename
      id
      title
      coverUrl
      averageRating
      description
      releaseDate
      creators {
        id
        name
      }
      ... on Movie {
        genres {
          name
        }
        cast {
          id
          name
        }
      }
      ... on TVShow {
        genres {
          name
        }
        cast {
          id
          name
        }
      }
      ... on Book {
        subjects {
          name
        }
        authors {
          id
          name
        }
      }
      ... on Game {
        genre
        themes
        keywords
        gameModes
        perspectives
        franchises
        platformsList
      }
    }
    me {
      id
      activities {
        id
        media {
          __typename
          id
          title
          coverUrl
          averageRating
          description
          releaseDate
          creators {
            id
            name
          }
          ... on Movie {
            genres {
              name
            }
            cast {
              id
              name
            }
          }
          ... on TVShow {
            genres {
              name
            }
            cast {
              id
              name
            }
          }
          ... on Book {
            subjects {
              name
            }
            authors {
              id
              name
            }
          }
          ... on Game {
            genre
            themes
            keywords
            gameModes
            perspectives
            franchises
            platformsList
          }
        }
      }
    }
  }
`;

export const GET_MEDIA_DETAILS_QUERY = gql`
  query GetMediaDetails($id: UUID!) {
    media(id: $id) {
      __typename
      id
      title
      releaseDate
      description
      coverUrl
      averageRating
      creators {
        id
        name
      }
      tags {
        id
        name
      }
      ... on Movie {
        runtime
        genres {
          name
        }
        cast {
          id
          name
        }
        crew {
          id
          name
        }
      }
      ... on TVShow {
        seasons
        episodes
        genres {
          name
        }
        cast {
          id
          name
        }
        crew {
          id
          name
        }
      }
      ... on Book {
        authors {
          id
          name
        }
      }
      ... on Game {
        genre
        themes
        keywords
        gameModes
        perspectives
        franchises
        platformsList
      }
      ... on MusicAlbum {
        label
      }
    }
    allMedia {
      __typename
      id
      title
      releaseDate
      coverUrl
      averageRating
      description
      creators {
        id
        name
      }
      ... on Movie {
        genres {
          name
        }
        cast {
          id
          name
        }
      }
      ... on TVShow {
        genres {
          name
        }
        cast {
          id
          name
        }
      }
      ... on Book {
        subjects {
          name
        }
        authors {
          id
          name
        }
      }
      ... on Game {
        genre
        themes
        keywords
        gameModes
        perspectives
        franchises
        platformsList
      }
    }
  }
`;

export const AUTOCOMPLETE_MEDIA_QUERY = gql`
  query AutocompleteMedia($type: MediaType!, $query: String!) {
    autocompleteMedia(type: $type, query: $query) {
      title
      year
      externalId
      imageUrl
      subtitle
    }
  }
`;

export const RECURSIVE_SEARCH_STATUS_QUERY = gql`
  query RecursiveSearchStatus($mediaId: UUID!) {
    recursiveSearchStatus(mediaId: $mediaId) {
      state
      completedAt
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
