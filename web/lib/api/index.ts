// Barrel export for the API layer
export { ApiError, httpClient } from "./client";
export type { TokenExpiredCallback } from "./client";
export { authApi } from "./auth";
export { deviceApi } from "./devices";
export { gameApi } from "./games";
export type {
  LoginResponse,
  DeviceResponse,
  GameSettingsDTO,
  GameResponse,
  TeamResponse,
  PlayerResponse,
  GameFullResponse,
  EventResponse,
} from "./types";
