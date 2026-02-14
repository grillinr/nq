import React, { createContext, useContext, useState, ReactNode } from 'react';
import { Media } from '../src/types';

interface MediaContextType {
  mediaList: Media[];
  addMedia: (newMedia: Omit<Media, 'id'>) => void;
}

const MediaContext = createContext<MediaContextType | undefined>(undefined);

export function MediaProvider({ children }: { children: ReactNode }) {
  const [mediaList, setMediaList] = useState<Media[]>([]);

  const addMedia = (newMedia: Omit<Media, 'id'>) => {
    const numericIds = mediaList
      .map((m) => {
        const idNum = typeof m.id === 'number' ? m.id : parseInt(m.id, 10);
        return Number.isFinite(idNum) ? idNum : 0;
      })
      .filter((id) => Number.isFinite(id));
    const newId = (numericIds.length > 0 ? Math.max(...numericIds) : 0) + 1;
    const nextList = [...mediaList, { ...newMedia, id: String(newId) }];
    setMediaList(nextList);
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
