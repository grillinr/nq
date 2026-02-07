import { useMemo } from "react";
import { useQuery } from "@apollo/client/react";
import { GET_MEDIA_DETAILS_QUERY } from "../../lib/graphql";
import { Media, MediaType } from "../types";

type MediaDetails = Media & {
  actors: { id: string; name: string }[];
  creators: { id: string; name: string }[];
  related: Media[];
  metaLabel?: string;
  metaValue?: string;
  tags: { id: string; name: string }[];
};

interface MediaDetailsData {
  media: any;
  allMedia: any[];
}

interface MediaDetailsVars {
  id: string;
}


export function useMediaDetails(id?: string) {
  const { data, loading, error } = useQuery<MediaDetailsData, MediaDetailsVars>(
    GET_MEDIA_DETAILS_QUERY,
    {
      variables: { id: id ?? "" },
      skip: !id,
      fetchPolicy: "cache-and-network",
    }
  );

  const details = useMemo(() => {
    if (!data?.media) return null;
    const media = data.media;
    const base: Media = {
      id: media.id,
      title: media.title ?? "Untitled",
      image:
        media.coverUrl ||
        `https://placehold.co/400x600?text=${encodeURIComponent(media.title ?? "Untitled")}`,
      rating: media.averageRating || 0,
      genre: extractGenres(media),
      year: media.releaseDate
        ? parseInt(media.releaseDate.substring(0, 4))
        : new Date().getFullYear(),
      duration: formatDuration(media),
      description: media.description || "",
      type: mapMediaType(media.__typename),
    };

    return {
      ...base,
      actors: extractActors(media),
      creators: media.creators ?? [],
      tags: media.tags ?? [],
      related: buildRelatedMedia(media, data.allMedia ?? []),
      ...buildMeta(media),
    } as MediaDetails;
  }, [data]);

  return { details, loading, error };
}

function mapMediaType(typename?: string): MediaType {
  switch (typename) {
    case "TVShow":
      return "tv";
    case "Book":
      return "book";
    case "Game":
      return "game";
    case "MusicAlbum":
      return "music";
    default:
      return "movie";
  }
}

function extractGenres(media: any): string[] {
  const tags: string[] = [];
  if (Array.isArray(media.genres)) {
    tags.push(...media.genres.map((g: any) => g.name).filter(Boolean));
  }
  if (Array.isArray(media.subjects)) {
    tags.push(...media.subjects.map((s: any) => s.name).filter(Boolean));
  }
  if (Array.isArray(media.genre)) {
    tags.push(...media.genre.filter(Boolean));
  }
  if (Array.isArray(media.themes)) {
    tags.push(...media.themes.filter(Boolean));
  }
  if (Array.isArray(media.keywords)) {
    tags.push(...media.keywords.filter(Boolean));
  }
  if (Array.isArray(media.gameModes)) {
    tags.push(...media.gameModes.filter(Boolean));
  }
  if (Array.isArray(media.perspectives)) {
    tags.push(...media.perspectives.filter(Boolean));
  }
  if (Array.isArray(media.franchises)) {
    tags.push(...media.franchises.filter(Boolean));
  }
  if (Array.isArray(media.platformsList)) {
    tags.push(...media.platformsList.filter(Boolean).map((p: string) => `platform:${p}`));
  }
  return tags;
}

function extractActors(media: any): { id: string; name: string }[] {
  if (Array.isArray(media.cast)) return media.cast;
  if (Array.isArray(media.authors)) return media.authors;
  return [];
}

function formatDuration(media: any): string | undefined {
  if (media.runtime) {
    return `${Math.floor(media.runtime / 60)}h ${media.runtime % 60}m`;
  }
  if (media.duration) {
    const minutes = Math.round(media.duration / 60);
    return `${Math.floor(minutes / 60)}h ${minutes % 60}m`;
  }
  if (media.pages) {
    return `${media.pages} pages`;
  }
  return undefined;
}

function buildMeta(media: any) {
  if (media.__typename === "TVShow") {
    const seasons = media.seasons ? `${media.seasons} seasons` : undefined;
    const episodes = media.episodes ? `${media.episodes} episodes` : undefined;
    if (seasons || episodes) {
      return { metaLabel: "Run", metaValue: [seasons, episodes].filter(Boolean).join(" · ") };
    }
  }
  if (media.__typename === "Book" && media.publisher) {
    return { metaLabel: "Publisher", metaValue: media.publisher };
  }
  if (media.__typename === "MusicAlbum" && media.label) {
    return { metaLabel: "Label", metaValue: media.label };
  }
  return {};
}
const LOW_WEIGHT_TAGS = new Set(["platform:"]);

function normalizeTag(tag: string) {
  return tag.trim().toLowerCase();
}

function tagWeight(tag: string) {
  for (const prefix of LOW_WEIGHT_TAGS) {
    if (tag.startsWith(prefix)) return 0.25;
  }
  return 1;
}

function buildRelatedMedia(media: any, allMedia: any[]): Media[] {
  const currentId = media.id;
  const currentGenres = new Set(extractGenres(media).map((g) => normalizeTag(g)));
  const related = allMedia
    .filter(
      (item) =>
        item?.id &&
        item.id !== currentId &&
        (item.__typename === "Movie" || item.__typename === "TVShow" || item.__typename === "Book" || item.__typename === "Game")
    )
    .map((item) => ({
      id: item.id,
      title: item.title ?? "Untitled",
      image:
        item.coverUrl ||
        `https://placehold.co/400x600?text=${encodeURIComponent(item.title ?? "Untitled")}`,
      rating: item.averageRating || 0,
      genre: extractGenres(item).map((g) => normalizeTag(g)),
      year: item.releaseDate
        ? parseInt(item.releaseDate.substring(0, 4))
        : new Date().getFullYear(),
      duration: undefined,
      description: item.description || "",
      type: mapMediaType(item.__typename),
    }))
    .sort((a, b) => scoreRelated(b.genre, currentGenres) - scoreRelated(a.genre, currentGenres))
    .slice(0, 12);
  return related;
}

function scoreRelated(genres: string[], currentGenres: Set<string>) {
  if (!genres.length || currentGenres.size === 0) return 0;
  return genres.reduce((score, genre) => {
    if (!currentGenres.has(genre)) return score;
    return score + tagWeight(genre);
  }, 0);
}
