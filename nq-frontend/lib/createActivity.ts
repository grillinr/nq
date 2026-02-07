import { getAccessToken } from "./auth";

interface CreateActivityInput {
  mediaId: string;
  statusId: number;
  rating?: number;
  review?: string;
  startedAt?: string;
  finishedAt?: string;
}

interface CreateActivityResult {
  id: string;
}

async function getDefaultUrl(): Promise<string> {
  const envUrl = process.env.EXPO_PUBLIC_API_URL;
  if (envUrl) return envUrl;
  return "http://localhost:8080/graphql";
}

const mutation = `mutation CreateActivity($input: CreateActivityInput!) {
  createActivity(input: $input) { id }
}`;

export async function createActivity(
  input: CreateActivityInput,
  opts?: { graphqlUrl?: string; signal?: AbortSignal },
): Promise<CreateActivityResult | null> {
  const graphqlUrl = opts?.graphqlUrl ?? (await getDefaultUrl());
  const accessToken = await getAccessToken();
  const res = await fetch(graphqlUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
    },
    body: JSON.stringify({ query: mutation, variables: { input } }),
    signal: opts?.signal,
  });

  const json = await res.json();
  if (json.errors) {
    throw new Error(JSON.stringify(json.errors));
  }

  return json.data ? json.data.createActivity : null;
}
