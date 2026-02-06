import { Media, MediaType, CreateMediaResult } from '../src/types';
import { getAccessToken } from './auth';

function parseDuration(durationStr?: string): number | undefined {
  if (!durationStr) return undefined;
  // Simple parser: looks for "Xh Ym" or just number
  let minutes = 0;
  const hoursMatch = durationStr.match(/(\d+)\s*h/);
  const minsMatch = durationStr.match(/(\d+)\s*m/);
  
  if (hoursMatch) minutes += parseInt(hoursMatch[1]) * 60;
  if (minsMatch) minutes += parseInt(minsMatch[1]);
  
  // validation fallback: if just a number, treat as minutes
  if (!hoursMatch && !minsMatch) {
    const raw = parseInt(durationStr);
    if (!isNaN(raw)) return raw;
  }
  
  return minutes > 0 ? minutes : undefined;
}

function formatDate(year?: number): string | undefined {
  if (!year) return undefined;
  return `${year}-01-01`;
}

// Helper to convert empty strings to undefined/null for backend
function emptyToNull(str?: string): string | undefined {
  if (!str || str.trim() === "") return undefined;
  return str;
}

async function getDefaultUrl(): Promise<string> {
  // Prefer explicit env override (e.g. via Expo constants or process.env)
  // If none provided, fall back to localhost for predictable local development.
  const envUrl = process.env.EXPO_PUBLIC_API_URL;
  if (envUrl) return envUrl;

  // For web or native development default to localhost
  return "http://localhost:8080/graphql";
}

const mutationMap: Record<
  MediaType,
  { mutation: string; responseField: string }
> = {
  movie: {
    mutation: `mutation CreateMovie($input: CreateMovieInput!) {
      createMovie(input: $input) { 
        id 
        title 
        description
        releaseDate
        coverUrl
      }
    }`,
    responseField: "createMovie",
  },
  tv: {
    mutation: `mutation CreateTVShow($input: CreateTVShowInput!) {
      createTVShow(input: $input) { 
        id 
        title 
        description
        releaseDate
        coverUrl
      }
    }`,
    responseField: "createTVShow",
  },
  book: {
    mutation: `mutation CreateBook($input: CreateBookInput!) {
      createBook(input: $input) { id title }
    }`,
    responseField: "createBook",
  },
  game: {
    mutation: `mutation CreateGame($input: CreateGameInput!) {
      createGame(input: $input) { id title }
    }`,
    responseField: "createGame",
  },
  music: {
    mutation: `mutation CreateMusicAlbum($input: CreateMusicAlbumInput!) {
      createMusicAlbum(input: $input) { id title }
    }`,
    responseField: "createMusicAlbum",
  },
};

export async function createMedia(
	mediaData: Omit<Media, "id">,
	opts?: { graphqlUrl?: string; signal?: AbortSignal },
): Promise<CreateMediaResult> {
  const { title, type } = mediaData;
  
  if (!title || !title.trim()) {
    throw new Error("title is required");
  }
  const entry = mutationMap[type];
  if (!entry) {
    throw new Error(`unsupported media type: ${type}`);
  }

  const graphqlUrl = opts?.graphqlUrl ?? (await getDefaultUrl());

  // Base input construction
  let input: any = {
    title: title.trim(),
    description: emptyToNull(mediaData.description),
    releaseDate: formatDate(mediaData.year),
    coverUrl: emptyToNull(mediaData.image),
    genres: mediaData.genre && mediaData.genre.length > 0 ? mediaData.genre : undefined,
  };

  // Add specific fields based on type
  if (type === 'movie') {
    input = {
      ...input,
      runtime: parseDuration(mediaData.duration),
    };
  }
  // TODO: Add specific mappings for other types (tv seasons, book pages, etc)

  const body = {
    query: entry.mutation,
    variables: { input },
  };

	const accessToken = await getAccessToken();
	const res = await fetch(graphqlUrl, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
		},
		body: JSON.stringify(body),
		signal: opts?.signal,
	});

  const json = await res.json();
  if (json.errors) {
    // surface GraphQL errors
    throw new Error(JSON.stringify(json.errors));
  }

  return json.data ? json.data[entry.responseField] : null;
}
