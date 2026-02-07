export const MAX_DEPTH = 2;
export const LAYER_DECAY = 0.5;
export const CANDIDATE_CAP = 300;

type NamedEntity = { id?: string | null; name?: string | null };

export type GraphMediaNode = {
  id: string | number;
  title?: string | null;
  averageRating?: number | null;
  genres?: { name?: string | null }[] | null;
  subjects?: { name?: string | null }[] | null;
  genre?: string[] | null;
  themes?: string[] | null;
  keywords?: string[] | null;
  gameModes?: string[] | null;
  perspectives?: string[] | null;
  franchises?: string[] | null;
  platformsList?: string[] | null;
  creators?: NamedEntity[] | null;
  cast?: NamedEntity[] | null;
  authors?: NamedEntity[] | null;
  tags?: { name?: string | null }[] | null;
};

type MediaFeatures = {
  id: string;
  tags: string[];
  people: string[];
};

type FeatureIndexes = {
  featuresById: Map<string, MediaFeatures>;
  adjacency: Map<string, Set<string>>;
};

export function capMediaCandidates(media: GraphMediaNode[], cap: number = CANDIDATE_CAP) {
  if (media.length <= cap) return media;
  return [...media]
    .sort((a, b) => {
      const ratingDiff = (b.averageRating ?? 0) - (a.averageRating ?? 0);
      if (ratingDiff !== 0) return ratingDiff;
      return (a.title ?? "").localeCompare(b.title ?? "");
    })
    .slice(0, cap);
}

export function scoreMediaFromUser(options: {
  candidates: GraphMediaNode[];
  activityMedia: GraphMediaNode[];
  maxDepth?: number;
  layerDecay?: number;
}) {
  const { candidates, activityMedia, maxDepth = MAX_DEPTH, layerDecay = LAYER_DECAY } = options;
  const allNodes = dedupeMedia([...candidates, ...activityMedia]);
  const { adjacency } = buildIndexes(allNodes);
  const startIds = activityMedia
    .map((item) => String(item.id))
    .filter((id) => adjacency.has(id));
  return computeScores({ adjacency, startIds, maxDepth, layerDecay });
}

export function scoreMediaFromRootMedia(options: {
  candidates: GraphMediaNode[];
  rootMedia: GraphMediaNode;
  maxDepth?: number;
  layerDecay?: number;
}) {
  const { candidates, rootMedia, maxDepth = MAX_DEPTH, layerDecay = LAYER_DECAY } = options;
  const allNodes = dedupeMedia([...candidates, rootMedia]);
  const { adjacency } = buildIndexes(allNodes);
  const rootId = String(rootMedia.id);
  const startIds = adjacency.get(rootId) ? Array.from(adjacency.get(rootId) ?? []) : [];
  return computeScores({ adjacency, startIds, maxDepth, layerDecay });
}

function computeScores(options: {
  adjacency: Map<string, Set<string>>;
  startIds: string[];
  maxDepth: number;
  layerDecay: number;
}) {
  const { adjacency, startIds, maxDepth, layerDecay } = options;
  const scores = new Map<string, number>();
  const seenDepth = new Map<string, number>();

  let frontierCounts = new Map<string, number>();
  for (const id of startIds) {
    frontierCounts.set(id, (frontierCounts.get(id) ?? 0) + 1);
  }

  for (let layer = 1; layer <= maxDepth && frontierCounts.size > 0; layer += 1) {
    const layerWeight = Math.pow(layerDecay, layer);
    const nextCounts = new Map<string, number>();
    const newlySeen = new Set<string>();

    for (const [id, count] of frontierCounts.entries()) {
      if (seenDepth.has(id)) continue;
      newlySeen.add(id);
      scores.set(id, (scores.get(id) ?? 0) + count * layerWeight);
    }

    for (const id of newlySeen) {
      seenDepth.set(id, layer);
      const neighbors = adjacency.get(id);
      if (!neighbors) continue;
      for (const neighborId of neighbors) {
        if (seenDepth.has(neighborId)) continue;
        nextCounts.set(neighborId, (nextCounts.get(neighborId) ?? 0) + 1);
      }
    }

    frontierCounts = nextCounts;
  }

  return scores;
}

function buildIndexes(media: GraphMediaNode[]): FeatureIndexes {
  const featuresById = new Map<string, MediaFeatures>();
  const tagIndex = new Map<string, Set<string>>();
  const personIndex = new Map<string, Set<string>>();

  for (const item of media) {
    const id = String(item.id);
    const features = extractMediaFeatures(item);
    featuresById.set(id, features);
    for (const tag of features.tags) {
      if (!tagIndex.has(tag)) tagIndex.set(tag, new Set());
      tagIndex.get(tag)!.add(id);
    }
    for (const person of features.people) {
      if (!personIndex.has(person)) personIndex.set(person, new Set());
      personIndex.get(person)!.add(id);
    }
  }

  const adjacency = new Map<string, Set<string>>();
  for (const [id, features] of featuresById.entries()) {
    const neighbors = new Set<string>();
    for (const tag of features.tags) {
      const ids = tagIndex.get(tag);
      if (!ids) continue;
      for (const neighborId of ids) {
        if (neighborId !== id) neighbors.add(neighborId);
      }
    }
    for (const person of features.people) {
      const ids = personIndex.get(person);
      if (!ids) continue;
      for (const neighborId of ids) {
        if (neighborId !== id) neighbors.add(neighborId);
      }
    }
    adjacency.set(id, neighbors);
  }

  return { featuresById, adjacency };
}

function extractMediaFeatures(media: GraphMediaNode): MediaFeatures {
  const tags = new Set<string>();
  const people = new Set<string>();

  const addTag = (value?: string | null) => {
    const normalized = normalizeValue(value);
    if (!normalized) return;
    tags.add(`tag:${normalized}`);
  };

  const addPlatform = (value?: string | null) => {
    const normalized = normalizeValue(value);
    if (!normalized) return;
    tags.add(`tag:platform:${normalized}`);
  };

  const addPerson = (value?: string | null, prefix: string = "person") => {
    const normalized = normalizeValue(value);
    if (!normalized) return;
    people.add(`${prefix}:${normalized}`);
  };

  media.genres?.forEach((g) => addTag(g?.name ?? undefined));
  media.subjects?.forEach((s) => addTag(s?.name ?? undefined));
  media.genre?.forEach((g) => addTag(g));
  media.themes?.forEach((t) => addTag(t));
  media.keywords?.forEach((k) => addTag(k));
  media.gameModes?.forEach((g) => addTag(g));
  media.perspectives?.forEach((p) => addTag(p));
  media.franchises?.forEach((f) => addTag(f));
  media.platformsList?.forEach((p) => addPlatform(p));
  media.tags?.forEach((t) => addTag(t?.name ?? undefined));

  media.creators?.forEach((c) => addPerson(c?.id ?? c?.name ?? undefined));
  media.cast?.forEach((c) => addPerson(c?.id ?? c?.name ?? undefined));
  media.authors?.forEach((a) => addPerson(a?.id ?? a?.name ?? undefined));

  return { id: String(media.id), tags: Array.from(tags), people: Array.from(people) };
}

function normalizeValue(value?: string | null) {
  if (!value) return "";
  return value.trim().toLowerCase();
}

function dedupeMedia(media: GraphMediaNode[]) {
  const byId = new Map<string, GraphMediaNode>();
  for (const item of media) {
    if (!item?.id) continue;
    byId.set(String(item.id), item);
  }
  return Array.from(byId.values());
}
