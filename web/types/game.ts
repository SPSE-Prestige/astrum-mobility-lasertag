export type GamePhase = "setup" | "live" | "results";
export type GameMode = "team" | "ffa";

export interface GameConfig {
  gameName: string;
  durationMinutes: number;
  teamsCount: number;
  gameMode: GameMode;
  friendlyFire: boolean;
  weaponTuning: {
    damage: number;
    fireRate: number;
    reloadTime: number;
  };
  playerTuning: {
    hp: number;
    respawnDelay: number;
  };
}

export interface Player {
  id: string;
  name: string;
  team: string;
  hp: number;
  ammo: number;
  status: "active" | "stunned";
  kills: number;
  hits: number;
  accuracy: number;
  damageDealt: number;
  cartConnected: boolean;
}

export interface KillFeedEvent {
  id: string;
  timestamp: string;
  message: string;
}

export interface TeamResult {
  team: string;
  score: number;
  accuracy: number;
  damageDealt: number;
}

export interface MatchHistoryItem {
  id: string;
  gameName: string;
  gameMode: GameMode;
  playedAt: string;
  durationMinutes: number;
  winner: string;
  totalKills: number;
}

export interface AuthState {
  isAuthenticated: boolean;
  username: string | null;
  error: string | null;
}

export interface GameState {
  phase: GamePhase;
  raceTimeSeconds: number;
  raceStatus: "idle" | "running" | "paused" | "finished";
  players: Player[];
  killFeed: KillFeedEvent[];
  teamResults: TeamResult[];
}