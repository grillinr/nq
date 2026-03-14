import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { Image } from 'expo-image';
import { useApolloClient } from '@apollo/client/react';
import { AUTOCOMPLETE_MEDIA_QUERY } from '../lib/graphql';
import { useTheme } from './ui/theme-provider';
import { fontSize, spacing, radii, zIndex } from './ui/tokens';

export type MediaSuggestion = {
  title: string;
  year?: number | null;
  externalId?: string | null;
  imageUrl?: string | null;
  subtitle?: string | null;
};

interface MediaAutocompleteProps {
  type: string;
  query: string;
  suppress: boolean;
  onSelect: (item: MediaSuggestion) => void;
  inputHeight: number;
}

const THUMB_SIZE = 36;

const createStyles = (colors: ReturnType<typeof useTheme>['colors'], inputBottom: number) =>
  StyleSheet.create({
    searching: {
      color: colors.mutedForeground,
      marginTop: spacing[2],
      fontSize: fontSize.sm,
    },
    suggestions: {
      position: 'absolute',
      top: inputBottom,
      left: 0,
      right: 0,
      zIndex: zIndex.modal,
      borderWidth: 1,
      borderColor: colors.border,
      borderRadius: radii.lg,
      backgroundColor: colors.background,
      overflow: 'hidden',
    },
    suggestionItem: {
      paddingVertical: spacing[3],
      paddingHorizontal: spacing[3],
      borderBottomWidth: 1,
      borderBottomColor: colors.border,
      flexDirection: 'row',
      alignItems: 'center',
      gap: spacing[3],
    },
    searchingIndicator: {
      position: 'absolute',
      top: inputBottom,
      left: 0,
    },
    suggestionItemLast: {
      borderBottomWidth: 0,
    },
    thumb: {
      width: THUMB_SIZE,
      height: THUMB_SIZE * 1.5,
      borderRadius: radii.sm,
      backgroundColor: colors.inputBackground,
      flexShrink: 0,
    },
    textBlock: {
      flex: 1,
    },
    titleRow: {
      flexDirection: 'row',
      alignItems: 'center',
      justifyContent: 'space-between',
      gap: spacing[2],
    },
    suggestionTitle: {
      fontSize: fontSize.base,
      color: colors.foreground,
      flexShrink: 1,
    },
    suggestionYear: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      flexShrink: 0,
    },
    suggestionSubtitle: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      marginTop: spacing[1],
    },
  });

export function MediaAutocomplete({
  type,
  query,
  suppress,
  onSelect,
  inputHeight,
}: MediaAutocompleteProps) {
  const { colors } = useTheme();
  const apolloClient = useApolloClient();

  const styles = React.useMemo(
    () => createStyles(colors, inputHeight + spacing[1]),
    [colors, inputHeight]
  );

  const [suggestions, setSuggestions] = useState<MediaSuggestion[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);

  useEffect(() => {
    if (suppress || !type || !query.trim()) {
      setSuggestions([]);
      setShowSuggestions(false);
      return undefined;
    }

    setIsSearching(true);
    const handle = setTimeout(async () => {
      try {
        const { data } = await apolloClient.query({
          query: AUTOCOMPLETE_MEDIA_QUERY,
          variables: {
            type: type.toUpperCase(),
            query: query.trim(),
          },
          fetchPolicy: 'no-cache',
        });
        setSuggestions((data as any)?.autocompleteMedia ?? []);
        setShowSuggestions(true);
      } catch {
        setSuggestions([]);
        setShowSuggestions(false);
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(handle);
  }, [apolloClient, suppress, type, query]);

  if (isSearching && query.trim()) {
    return (
      <ActivityIndicator
        size="small"
        color={colors.mutedForeground}
        style={styles.searchingIndicator}
      />
    );
  }

  if (!showSuggestions || suggestions.length === 0) {
    return null;
  }

  return (
    <View style={styles.suggestions}>
      {suggestions.map((item, index) => (
        <TouchableOpacity
          key={item.externalId ?? `${item.title}-${index}`}
          style={[
            styles.suggestionItem,
            index === suggestions.length - 1 ? styles.suggestionItemLast : null,
          ]}
          onPress={() => onSelect(item)}
        >
          {item.imageUrl ? (
            <Image
              source={{ uri: item.imageUrl }}
              style={styles.thumb}
              contentFit="cover"
              transition={120}
            />
          ) : (
            <View style={styles.thumb} />
          )}
          <View style={styles.textBlock}>
            <View style={styles.titleRow}>
              <Text style={styles.suggestionTitle} numberOfLines={1}>
                {item.title}
              </Text>
              {item.year ? <Text style={styles.suggestionYear}>{item.year}</Text> : null}
            </View>
            {item.subtitle ? (
              <Text style={styles.suggestionSubtitle} numberOfLines={1}>
                {item.subtitle}
              </Text>
            ) : null}
          </View>
        </TouchableOpacity>
      ))}
    </View>
  );
}
