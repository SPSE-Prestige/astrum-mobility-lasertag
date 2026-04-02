import { httpClient } from "./client";
import type { LoginResponse } from "./types";

export const authApi = {
  login: (username: string, password: string) =>
    httpClient.post<LoginResponse>("/auth/login", { username, password }),

  logout: () => httpClient.post<void>("/auth/logout"),
};
