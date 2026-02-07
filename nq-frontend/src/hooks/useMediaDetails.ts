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

const PLACEHOLDER_IMAGE = "https://placehold.co/400x600?text=No+Image";

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
      image: media.coverUrl || PLACEHOLDER_IMAGE,
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
  if (Array.isArray(media.genres)) {
    return media.genres.map((g: any) => g.name).filter(Boolean);
  }
  if (Array.isArray(media.genre)) {
    return media.genre.filter(Boolean);
  }
  return [];
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
function buildRelatedMedia(media: any, allMedia: any[]): Media[] {
  const currentId = media.id;
  const currentGenres = new Set(extractGenres(media));
  const related = allMedia
    .filter(
      (item) =>
        item?.id &&
        item.id !== currentId &&
        (item.__typename === "Movie" || item.__typename === "TVShow")
    )
    .map((item) => ({
      id: item.id,
      title: item.title ?? "Untitled",
      image: item.coverUrl || PLACEHOLDER_IMAGE,
      rating: item.averageRating || 0,
      genre: extractGenres(item),
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
  return genres.reduce((score, genre) => (currentGenres.has(genre) ? score + 1 : score), 0);
}
