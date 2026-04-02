import { httpClient } from "./client";
import type {
  EventResponse,
  GameFullResponse,
  GameResponse,
  GameSettingsDTO,
  PlayerResponse,
  TeamResponse,
} from "./types";

const enc = encodeURIComponent;

export const gameApi = {
  create: (settings?: GameSettingsDTO) =>
    httpClient.post<GameResponse>("/games", { settings }),

  get: (id: string) =>
    httpClient.get<GameResponse>(`/games/${enc(id)}`),

  getFull: (id: string) =>
    httpClient.get<GameFullResponse>(`/games/${enc(id)}/full`),

  list: () =>
    httpClient.get<GameResponse[]>("/games"),

  start: (id: string) =>
    httpClient.post<GameResponse>(`/games/${enc(id)}/start`),

  end: (id: string) =>
    httpClient.post<GameResponse>(`/games/${enc(id)}/end`),

  updateSettings: (id: string, settings: GameSettingsDTO) =>
    httpClient.patch<GameResponse>(`/games/${enc(id)}/settings`, { settings }),

  // ── Teams ──

  addTeam: (gameId: string, name: string, color: string) =>
    httpClient.post<TeamResponse>(`/games/${enc(gameId)}/teams`, { name, color }),

  listTeams: (gameId: string) =>
    httpClient.get<TeamResponse[]>(`/games/${enc(gameId)}/teams`),

  removeTeam: (gameId: string, teamId: string) =>
    httpClient.del<void>(`/games/${enc(gameId)}/teams/${enc(teamId)}`),

  // ── Players ──

  addPlayer: (gameId: string, deviceId: string, nickname: string, teamId?: string) =>
    httpClient.post<PlayerResponse>(`/games/${enc(gameId)}/players`, {
      device_id: deviceId,
      nickname,
      team_id: teamId ?? null,
    }),

  listPlayers: (gameId: string) =>
    httpClient.get<PlayerResponse[]>(`/games/${enc(gameId)}/players`),

  removePlayer: (gameId: string, playerId: string) =>
    httpClient.del<void>(`/games/${enc(gameId)}/players/${enc(playerId)}`),

  updatePlayerTeam: (gameId: string, playerId: string, teamId: string | null) =>
    httpClient.patch<void>(`/games/${enc(gameId)}/players/${enc(playerId)}/team`, {
      team_id: teamId,
    }),

  // ── Leaderboard & Events ──

  getLeaderboard: (gameId: string) =>
    httpClient.get<PlayerResponse[]>(`/games/${enc(gameId)}/leaderboard`),

  getEvents: (gameId: string) =>
    httpClient.get<EventResponse[]>(`/games/${enc(gameId)}/events`),
};
