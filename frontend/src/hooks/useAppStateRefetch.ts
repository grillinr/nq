import { useEffect, useRef } from 'react';
import { AppState, AppStateStatus } from 'react-native';

/**
 * Calls `refetch` whenever the app transitions from background/inactive to active.
 * This keeps Apollo queries fresh after the user returns to the app without
 * coupling the Apollo client directly to the AppState listener.
 */
export function useAppStateRefetch(refetch: () => void) {
  const appState = useRef<AppStateStatus>(AppState.currentState);

  useEffect(() => {
    const subscription = AppState.addEventListener('change', (nextState: AppStateStatus) => {
      if (appState.current.match(/inactive|background/) && nextState === 'active') {
        refetch();
      }
      appState.current = nextState;
    });
    return () => subscription.remove();
  }, [refetch]);
}
