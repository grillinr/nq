import React, { createContext, useContext, useState, ReactNode } from 'react';
import { createMedia } from '../lib/createMedia';
import { logError } from './logger';

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

const mockData: Media[] = [
  {
    id: 1,
    title: 'Inception',
    image: 'https://images.unsplash.com/photo-1524712245354-2c4e5e7121c0?w=400',
    rating: 8.8,
    genre: ['Sci-Fi', 'Thriller', 'Action'],
    year: 2010,
    duration: '2h 28m',
    description: 'A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea.',
    type: 'movie',
  },
  // Add more as needed
];

interface MediaContextType {
  mediaList: Media[];
  addMedia: (newMedia: Omit<Media, 'id'>) => void;
}

const MediaContext = createContext<MediaContextType | undefined>(undefined);

export function MediaProvider({ children }: { children: ReactNode }) {
  const [mediaList, setMediaList] = useState<Media[]>(mockData);

  const addMedia = async (newMedia: Omit<Media, 'id'>) => {
    try {
      const result = await createMedia(newMedia.type, newMedia.title);
      if (result) {
        const newId = Math.max(...mediaList.map((m) => m.id), 0) + 1;
        setMediaList([...mediaList, { ...newMedia, id: newId }]);
      }
    } catch (error) {
      logError('Failed to add media:', error);
    }
  };

  return (
    <MediaContext.Provider value={{ mediaList, addMedia }}>
      {children}
    </MediaContext.Provider>
  );
}

export function useMedia() {
  const context = useContext(MediaContext);
  if (!context) {
    throw new Error('useMedia must be used within a MediaProvider');
  }
  return context;
}