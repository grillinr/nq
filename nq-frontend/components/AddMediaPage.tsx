import React, { useState } from 'react';
import { View, Text, ScrollView, StyleSheet, Alert, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Button } from './ui/button';
import Input from './ui/input';
import Card from './ui/card';
import Badge from './ui/badge';
import { colors, fontSize, spacing } from './ui/tokens';

interface Media {
  id: number;
  title: string;
  image: string;
  rating: number;
  genre: string[];
  year: number;
  duration?: string;
  description: string;
  type: 'movie' | 'tv' | 'book' | 'music' | 'game';
}

interface AddMediaPageProps {
  onBack: () => void;
  onAddMedia: (media: Omit<Media, 'id'>) => void;
}

const genreOptions = {
  movie: ['Action', 'Sci-Fi', 'Drama', 'Thriller', 'Crime', 'Comedy', 'Horror', 'Romance'],
  tv: ['Crime', 'Drama', 'Sci-Fi', 'Horror', 'Historical', 'Action', 'Comedy', 'Reality'],
  book: ['Fiction', 'Classic', 'Dystopian', 'Sci-Fi', 'Drama', 'Mystery', 'Fantasy', 'Biography'],
  music: ['Rock', 'Pop', 'Jazz', 'R&B', 'Classic', 'Progressive', 'Hip-Hop', 'Electronic'],
  game: ['Action', 'RPG', 'Strategy', 'Puzzle', 'Adventure', 'Sports', 'Simulation', 'Horror'],
};

const typeOptions = [
  { label: 'Movie', value: 'movie', icon: 'film-outline' },
  { label: 'TV Show', value: 'tv', icon: 'tv-outline' },
  { label: 'Book', value: 'book', icon: 'book-outline' },
  { label: 'Music', value: 'music', icon: 'musical-notes-outline' },
  { label: 'Game', value: 'game', icon: 'game-controller-outline' },
];

function AddMediaPage({ onBack, onAddMedia }: AddMediaPageProps) {
  const [title, setTitle] = useState('');
  const [type, setType] = useState<'movie' | 'tv' | 'book' | 'music' | 'game'>('movie');
  const [description, setDescription] = useState('');
  const [year, setYear] = useState(new Date().getFullYear().toString());
  const [rating, setRating] = useState('7.5');
  const [duration, setDuration] = useState('');
  const [imageUrl, setImageUrl] = useState('');
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);

  const handleGenreToggle = (genre: string) => {
    setSelectedGenres((prev) =>
      prev.includes(genre) ? prev.filter((g) => g !== genre) : [...prev, genre]
    );
  };

  const handleSubmit = () => {
    if (!title.trim()) {
      Alert.alert('Error', 'Please enter a title');
      return;
    }

    const newMedia: Omit<Media, 'id'> = {
      title: title.trim(),
      type,
      description: description.trim() || 'No description provided.',
      year: parseInt(year) || new Date().getFullYear(),
      rating: parseFloat(rating) || 7.5,
      duration: duration.trim() || undefined,
      image: imageUrl.trim() || 'https://images.unsplash.com/photo-1485846234645-a62644f84728?w=400',
      genre: selectedGenres.length > 0 ? selectedGenres : ['Uncategorized'],
    };

    onAddMedia(newMedia);

    // Reset form
    setTitle('');
    setDescription('');
    setYear(new Date().getFullYear().toString());
    setRating('7.5');
    setDuration('');
    setImageUrl('');
    setSelectedGenres([]);
  };

  const handleTypeChange = (value: string) => {
    setType(value as 'movie' | 'tv' | 'book' | 'music' | 'game');
    setSelectedGenres([]);
  };

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <Button variant="ghost" onPress={onBack} style={styles.backButton}>
        <Ionicons name="arrow-back" size={16} color={colors.foreground} />
        <Text style={styles.backText}>Back to Media List</Text>
      </Button>

      <View style={styles.header}>
        <Text style={styles.title}>Add New Media</Text>
        <Text style={styles.subtitle}>Fill in the details to add a new item to your collection</Text>
      </View>

      <Card style={styles.card}>
        <View style={styles.form}>
          <View style={styles.field}>
            <Text style={styles.label}>Title *</Text>
            <Input
              value={title}
              onChangeText={setTitle}
              placeholder="Enter media title"
            />
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Type *</Text>
            <View style={styles.typeOptions}>
              {typeOptions.map((option) => (
                <TouchableOpacity
                  key={option.value}
                  style={[styles.typeOption, type === option.value && styles.typeOptionSelected]}
                  onPress={() => handleTypeChange(option.value)}
                >
                  <Ionicons name={option.icon as any} size={20} color={type === option.value ? 'white' : colors.foreground} />
                  <Text style={[styles.typeText, type === option.value && styles.typeTextSelected]}>
                    {option.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Description</Text>
            <Input
              value={description}
              onChangeText={setDescription}
              placeholder="Enter description"
              multiline
              numberOfLines={3}
            />
          </View>

          <View style={styles.row}>
            <View style={styles.field}>
              <Text style={styles.label}>Year</Text>
              <Input
                value={year}
                onChangeText={setYear}
                placeholder="2024"
                keyboardType="numeric"
              />
            </View>
            <View style={styles.field}>
              <Text style={styles.label}>Rating</Text>
              <Input
                value={rating}
                onChangeText={setRating}
                placeholder="7.5"
                keyboardType="numeric"
              />
            </View>
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Duration</Text>
            <Input
              value={duration}
              onChangeText={setDuration}
              placeholder="2h 30m or 300 pages"
            />
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Image URL</Text>
            <Input
              value={imageUrl}
              onChangeText={setImageUrl}
              placeholder="https://..."
            />
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Genres</Text>
            <View style={styles.genres}>
              {genreOptions[type].map((genre) => (
                <TouchableOpacity
                  key={genre}
                  onPress={() => handleGenreToggle(genre)}
                  style={styles.genreButton}
                >
                  <Badge variant={selectedGenres.includes(genre) ? 'default' : 'secondary'}>
                    {genre}
                  </Badge>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <Button onPress={handleSubmit} style={styles.submitButton}>
            <Ionicons name="add" size={20} color="white" />
            <Text style={styles.submitText}>Add Media</Text>
          </Button>
        </View>
      </Card>
    </ScrollView>
  );
}

export default AddMediaPage;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing[4],
  },
  backButton: {
    marginBottom: spacing[6],
    alignSelf: 'flex-start',
  },
  backText: {
    marginLeft: spacing[2],
    color: colors.foreground,
  },
  header: {
    marginBottom: spacing[6],
  },
  title: {
    fontSize: fontSize.xl,
    fontWeight: '600',
    color: colors.primary,
    marginBottom: spacing[2],
  },
  subtitle: {
    fontSize: fontSize.base,
    color: colors['muted-foreground'],
  },
  card: {
    padding: spacing[6],
  },
  form: {},
  field: {
    marginBottom: spacing[6],
  },
  label: {
    fontSize: fontSize.base,
    fontWeight: '500',
    color: colors.foreground,
    marginBottom: spacing[2],
  },
  typeOptions: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing[2],
  },
  typeOption: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing[2],
    padding: spacing[3],
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.background,
  },
  typeOptionSelected: {
    backgroundColor: colors.primary,
    borderColor: colors.primary,
  },
  typeText: {
    fontSize: fontSize.sm,
    color: colors.foreground,
  },
  typeTextSelected: {
    color: 'white',
  },
  row: {
    flexDirection: 'row',
    gap: spacing[4],
  },
  genres: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing[2],
  },
  genreButton: {},
  submitButton: {
    marginTop: spacing[6],
  },
  submitText: {
    color: 'white',
    marginLeft: spacing[2],
  },
});