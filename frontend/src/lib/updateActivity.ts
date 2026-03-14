import { getAccessToken } from './auth';

interface UpdateActivityInput {
  statusId?: number;
  rating?: number;
  review?: string;
  startedAt?: string;
  finishedAt?: string;
}

interface UpdateActivityResult {
  id: string;
  rating?: number;
  review?: string;
  status: {
    id: number;
    name: string;
  };
}

async function getDefaultUrl(): Promise<string> {
  const envUrl = process.env.EXPO_PUBLIC_API_URL;
  if (envUrl) return envUrl;
  return 'http://localhost:8080/graphql';
}

const mutation = `mutation UpdateActivity($id: UUID!, $input: UpdateActivityInput!) {
  updateActivity(id: $id, input: $input) {
    id
    rating
    review
    status {
      id
      name
    }
    startedAt
    finishedAt
  }
}`;

export async function updateActivity(
  id: string,
  input: UpdateActivityInput,
  opts?: { graphqlUrl?: string; signal?: AbortSignal }
): Promise<UpdateActivityResult | null> {
  const graphqlUrl = opts?.graphqlUrl ?? (await getDefaultUrl());
  const accessToken = await getAccessToken();
  const res = await fetch(graphqlUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
    },
    body: JSON.stringify({ query: mutation, variables: { id, input } }),
    signal: opts?.signal,
  });

  const json = await res.json();
  if (json.errors) {
    throw new Error(JSON.stringify(json.errors));
  }

  return json.data ? json.data.updateActivity : null;
}
