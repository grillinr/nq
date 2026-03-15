import React, { useEffect, useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';
import { useApolloClient } from '@apollo/client/react';
import { SEARCH_USERS_QUERY } from '../lib/graphql';
import { useTheme } from './ui/theme-provider';
import { fontSize, spacing, radii, zIndex } from './ui/tokens';

export type UserSuggestion = {
  id: string;
  name: string;
  avatarUrl?: string | null;
};

interface UserAutocompleteProps {
  query: string;
  suppress: boolean;
  onSelect: (user: UserSuggestion) => void;
  inputHeight: number;
}

const createStyles = (colors: ReturnType<typeof useTheme>['colors'], inputBottom: number) =>
  StyleSheet.create({
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
    searchingIndicator: {
      position: 'absolute',
      top: inputBottom,
      left: 0,
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
    suggestionItemLast: {
      borderBottomWidth: 0,
    },
    avatar: {
      width: 32,
      height: 32,
      borderRadius: radii.full,
      backgroundColor: colors.inputBackground,
      alignItems: 'center',
      justifyContent: 'center',
      flexShrink: 0,
    },
    avatarText: {
      fontSize: fontSize.sm,
      color: colors.mutedForeground,
      fontWeight: '600',
    },
    suggestionName: {
      flex: 1,
      fontSize: fontSize.base,
      color: colors.foreground,
    },
  });

export function UserAutocomplete({
  query,
  suppress,
  onSelect,
  inputHeight,
}: UserAutocompleteProps) {
  const { colors } = useTheme();
  const apolloClient = useApolloClient();
  const styles = React.useMemo(
    () => createStyles(colors, inputHeight + spacing[1]),
    [colors, inputHeight]
  );

  const [suggestions, setSuggestions] = useState<UserSuggestion[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);

  useEffect(() => {
    if (suppress || !query.trim()) {
      setSuggestions([]);
      setShowSuggestions(false);
      return undefined;
    }

    setIsSearching(true);
    const handle = setTimeout(async () => {
      try {
        const { data } = await apolloClient.query({
          query: SEARCH_USERS_QUERY,
          variables: { query: query.trim() },
          fetchPolicy: 'no-cache',
        });
        setSuggestions((data as { searchUsers: UserSuggestion[] })?.searchUsers ?? []);
        setShowSuggestions(true);
      } catch {
        setSuggestions([]);
        setShowSuggestions(false);
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(handle);
  }, [apolloClient, suppress, query]);

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
      {suggestions.map((user, index) => (
        <TouchableOpacity
          key={user.id}
          style={[
            styles.suggestionItem,
            index === suggestions.length - 1 ? styles.suggestionItemLast : null,
          ]}
          onPress={() => onSelect(user)}
        >
          <View style={styles.avatar}>
            <Text style={styles.avatarText}>{user.name.charAt(0).toUpperCase()}</Text>
          </View>
          <Text style={styles.suggestionName} numberOfLines={1}>
            {user.name}
          </Text>
        </TouchableOpacity>
      ))}
    </View>
  );
}
