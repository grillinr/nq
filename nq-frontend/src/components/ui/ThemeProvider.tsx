import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { Appearance } from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { lightColors, darkColors, ColorPalette } from './tokens';

type ThemeChoice = 'light' | 'dark' | 'auto';
type Resolved = 'light' | 'dark';

const THEME_KEY = 'nq.theme.choice';

type ThemeContextValue = {
  theme: ThemeChoice;
  resolved: Resolved;
  colors: ColorPalette;
  setTheme: (t: ThemeChoice) => void;
};

const defaultColors = lightColors;

const ThemeContext = createContext<ThemeContextValue>({
  theme: 'auto',
  resolved: 'light',
  colors: defaultColors,
  setTheme: () => {},
});

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setThemeState] = useState<ThemeChoice>('auto');
  const [systemScheme, setSystemScheme] = useState<Resolved>(
    Appearance.getColorScheme() === 'dark' ? 'dark' : 'light'
  );

  useEffect(() => {
    (async () => {
      try {
        const stored = await AsyncStorage.getItem(THEME_KEY);
        if (stored === 'light' || stored === 'dark' || stored === 'auto') {
          setThemeState(stored);
        }
      } catch {
        // ignore
      }
    })();

    const sub = Appearance.addChangeListener(({ colorScheme }) => {
      setSystemScheme(colorScheme === 'dark' ? 'dark' : 'light');
    });
    return () => sub.remove();
  }, []);

  const setTheme = async (t: ThemeChoice) => {
    try {
      await AsyncStorage.setItem(THEME_KEY, t);
    } catch {
      // ignore
    }
    setThemeState(t);
  };

  const resolved: Resolved = theme === 'auto' ? systemScheme : theme;
  const colors = useMemo(() => (resolved === 'dark' ? darkColors : lightColors), [resolved]);

  return (
    <ThemeContext.Provider value={{ theme, resolved, colors, setTheme }}>{children}</ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
