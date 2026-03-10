import React, { createContext, useContext, useCallback, useEffect, useRef, useState, ReactNode } from "react";
import { AppState, AppStateStatus } from "react-native";
import { getAccessToken, loginWithAuth0, logout as logoutFromAuth } from "./auth";

interface AuthContextType {
  hasToken: boolean;
  isChecking: boolean;
  refreshAuth: () => Promise<void>;
  login: () => Promise<boolean>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [hasToken, setHasToken] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const appState = useRef(AppState.currentState);

  const refreshAuth = useCallback(async () => {
    const token = await getAccessToken();
    setHasToken(!!token);
    setIsChecking(false);
  }, []);

  const login = useCallback(async (): Promise<boolean> => {
    const token = await loginWithAuth0();
    if (!token) return false;
    await refreshAuth();
    return true;
  }, [refreshAuth]);

  const logout = useCallback(async () => {
    await logoutFromAuth();
    setHasToken(false);
  }, []);

  // Initial token check on mount
  useEffect(() => {
    refreshAuth();
  }, [refreshAuth]);

  // Re-check auth whenever the app comes back to the foreground.
  // This handles the Expo Go deep-link case where the browser redirects back
  // and the component tree may have remounted before loginWithAuth0() returned.
  useEffect(() => {
    const subscription = AppState.addEventListener("change", (nextState: AppStateStatus) => {
      if (appState.current.match(/inactive|background/) && nextState === "active") {
        refreshAuth();
      }
      appState.current = nextState;
    });
    return () => subscription.remove();
  }, [refreshAuth]);

  return (
    <AuthContext.Provider value={{ hasToken, isChecking, refreshAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
