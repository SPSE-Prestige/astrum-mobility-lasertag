export type GamePhase = "setup" | "players" | "live" | "results";
export type GameMode = "team" | "ffa";

export interface GameConfig {
  gameName: string;
  durationMinutes: number;
  teamsCount: number;
  gameMode: GameMode;
  friendlyFire: boolean;
  respawnDelay: number;
  maxPlayers: number;
}

export interface Player {
  id: string;
  loginCode?: string; // Unique admin-assigned code for mobile login
  name: string;
  team: string;
  teamId: string | null;
  deviceId: string;
  status: "alive" | "dead";
  kills: number;
  deaths: number;
  score: number;
}

export interface Team {
  id: string;
  name: string;
  color: string;
}

export interface Device {
  id: string;
  deviceId: string;
  status: string;
  lastSeen: string;
}

export interface KillFeedEvent {
  id: string;
  timestamp: string;
  message: string;
}

export interface TeamResult {
  team: string;
  score: number;
  kills: number;
  deaths: number;
}

export interface MatchHistoryItem {
  id: string;
  gameName: string;
  gameMode: GameMode;
  playedAt: string;
  durationMinutes: number;
  status: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  username: string | null;
  token: string | null;
  error: string | null;
}

export interface GameState {
  phase: GamePhase;
  gameId: string | null;
  raceTimeSeconds: number;
  raceStatus: "idle" | "running" | "finished";
  players: Player[];
  teams: Team[];
  devices: Device[];
  killFeed: KillFeedEvent[];
  teamResults: TeamResult[];
}