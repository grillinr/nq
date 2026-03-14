import { useMemo } from "react";
import { useQuery } from "@apollo/client/react";
import { GET_MEDIA_DETAILS_QUERY } from "../../lib/graphql";
import { Media, MediaType, UserActivity } from "../types";
import { useAppStateRefetch } from "./useAppStateRefetch";

type MediaDetails = Media & {
  actors: { id: string; name: string }[];
  creators: { id: string; name: string }[];
  related: Media[];
  metaLabel?: string;
  metaValue?: string;
  tags: { id: string; name: string }[];
  myActivity?: UserActivity | null;
};

interface MediaDetailsData {
  media: any;
}

interface MediaDetailsVars {
  id: string;
}


export function useMediaDetails(id?: string) {
  const { data, loading, error, refetch } = useQuery<MediaDetailsData, MediaDetailsVars>(
    GET_MEDIA_DETAILS_QUERY,
    {
      variables: { id: id ?? "" },
      skip: !id,
      fetchPolicy: "cache-and-network",
      nextFetchPolicy: "cache-first",
    }
  );

  useAppStateRefetch(refetch);

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
      myActivity: media.myActivity ?? null,
      related: mapRelatedMedia(media.relatedMedia ?? []),
      ...buildMeta(media),
    } as MediaDetails;
  }, [data]);

  return { details, loading, error, refetch };
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

function mapRelatedMedia(items: any[]): Media[] {
  return items.map((item) => ({
    id: item.id,
    title: item.title ?? "Untitled",
    image:
      item.coverUrl ||
      `https://placehold.co/400x600?text=${encodeURIComponent(item.title ?? "Untitled")}`,
    rating: item.averageRating || 0,
    genre: [],
    year: item.releaseDate
      ? parseInt(item.releaseDate.substring(0, 4))
      : new Date().getFullYear(),
    duration: undefined,
    description: item.description || "",
    type: mapMediaType(item.__typename),
  }));
}

