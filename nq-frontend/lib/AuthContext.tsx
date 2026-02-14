import React, { createContext, useContext, useCallback, useEffect, useState, ReactNode } from "react";
import { getAccessToken, loginWithAuth0, logout as logoutFromAuth } from "./auth";

interface AuthContextType {
  hasToken: boolean;
  isChecking: boolean;
  refreshAuth: () => Promise<void>;
  login: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [hasToken, setHasToken] = useState(false);
  const [isChecking, setIsChecking] = useState(true);

  const refreshAuth = useCallback(async () => {
    const token = await getAccessToken();
    setHasToken(!!token);
    setIsChecking(false);
  }, []);

  const login = useCallback(async () => {
    await loginWithAuth0();
    await refreshAuth();
  }, [refreshAuth]);

  const logout = useCallback(async () => {
    await logoutFromAuth();
    setHasToken(false);
  }, []);

  useEffect(() => {
    refreshAuth();
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
