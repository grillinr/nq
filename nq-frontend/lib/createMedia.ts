
export type MediaType = "movie" | "tv" | "book" | "game" | "music";

async function getDefaultUrl(): Promise<string> {
  // Prefer explicit env override (e.g. via Expo constants or process.env)
  // If none provided, fall back to localhost for predictable local development.
  const envUrl = (global as any)?.GRAPHQL_URL ?? process.env.GRAPHQL_URL;
  if (envUrl) return envUrl;

  // For web or native development default to localhost
  return "http://localhost:8080/query";
}

const mutationMap: Record<
  MediaType,
  { mutation: string; responseField: string }
> = {
  movie: {
    mutation: `mutation CreateMovie($input: CreateMovieInput!) {
      createMovie(input: $input) { id title }
    }`,
    responseField: "createMovie",
  },
  tv: {
    mutation: `mutation CreateTVShow($input: CreateTVShowInput!) {
      createTVShow(input: $input) { id title }
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

export type CreateMediaResult = { id: string; title: string } | null;

export async function createMedia(
  mediaType: MediaType,
  title: string,
  opts?: { graphqlUrl?: string; signal?: AbortSignal },
): Promise<CreateMediaResult> {
  if (!title || !title.trim()) {
    throw new Error("title is required");
  }
  const entry = mutationMap[mediaType];
  if (!entry) {
    throw new Error(`unsupported media type: ${mediaType}`);
  }

  const graphqlUrl = opts?.graphqlUrl ?? (await getDefaultUrl());

  const body = {
    query: entry.mutation,
    variables: { input: { title: title.trim() } },
  };

  const res = await fetch(graphqlUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
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
