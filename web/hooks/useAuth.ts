"use client";

import { useCallback, useEffect, useState } from "react";
import type { AuthState } from "@/types/game";
import { authApi, httpClient } from "@/lib/api/index";

const TOKEN_KEY = "lasertag-token";
const USER_KEY = "lasertag-user";

const initialAuth: AuthState = {
  isAuthenticated: false,
  username: null,
  token: null,
  error: null,
};

export function useAuth() {
  const [auth, setAuth] = useState<AuthState>(() => {
    if (typeof window === "undefined") return initialAuth;
    const token = localStorage.getItem(TOKEN_KEY);
    const username = localStorage.getItem(USER_KEY);
    if (token && username) {
      httpClient.setToken(token);
      return { isAuthenticated: true, username, token, error: null };
    }
    return initialAuth;
  });

  // Persist auth token to localStorage and sync with httpClient
  useEffect(() => {
    if (auth.token) {
      localStorage.setItem(TOKEN_KEY, auth.token);
      localStorage.setItem(USER_KEY, auth.username ?? "");
      httpClient.setToken(auth.token);
    } else {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      httpClient.setToken(null);
    }
  }, [auth.token, auth.username]);

  const login = useCallback(async (username: string, password: string): Promise<boolean> => {
    try {
      const res = await authApi.login(username, password);
      setAuth({ isAuthenticated: true, username, token: res.token, error: null });
      return true;
    } catch {
      setAuth({ isAuthenticated: false, username: null, token: null, error: "invalid_credentials" });
      return false;
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
    } catch {}
    setAuth(initialAuth);
  }, []);

  return { auth, login, logout };
}
