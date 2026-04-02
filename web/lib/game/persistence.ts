import type { GameConfig, GamePhase, GameState } from "@/types/game";

const GAME_ID_KEY = "lasertag-gameId";
const PHASE_KEY = "lasertag-phase";
const CONFIG_KEY = "lasertag-config";

export const defaultConfig: GameConfig = {
  gameName: "Laser Tag Game",
  durationMinutes: 5,
  teamsCount: 2,
  gameMode: "team",
  friendlyFire: false,
  respawnDelay: 5,
  maxPlayers: 20,
  killsPerUpgrade: 3,
};

export const initialState: GameState = {
  phase: "setup",
  gameId: null,
  raceTimeSeconds: defaultConfig.durationMinutes * 60,
  raceStatus: "idle",
  players: [],
  teams: [],
  devices: [],
  killFeed: [],
  teamResults: [],
};

export function loadPersistedState(): Partial<GameState> | null {
  if (typeof window === "undefined") return null;
  const gameId = localStorage.getItem(GAME_ID_KEY);
  const phase = localStorage.getItem(PHASE_KEY) as GamePhase | null;
  if (!gameId) return null;
  const validPhases: GamePhase[] = ["setup", "players", "live", "results"];
  const resolvedPhase = phase && validPhases.includes(phase) ? phase : "setup";
  return {
    gameId,
    phase: resolvedPhase,
    raceStatus: resolvedPhase === "live" ? "running" : "idle",
  };
}

export function loadPersistedConfig(): GameConfig | null {
  if (typeof window === "undefined") return null;
  try {
    const raw = localStorage.getItem(CONFIG_KEY);
    return raw ? (JSON.parse(raw) as GameConfig) : null;
  } catch {
    return null;
  }
}

export function persistGameState(gameId: string | null, phase: GamePhase) {
  if (gameId) {
    localStorage.setItem(GAME_ID_KEY, gameId);
    localStorage.setItem(PHASE_KEY, phase);
  } else {
    clearPersistedGame();
  }
}

export function persistConfig(config: GameConfig) {
  localStorage.setItem(CONFIG_KEY, JSON.stringify(config));
}

export function clearPersistedGame() {
  localStorage.removeItem(GAME_ID_KEY);
  localStorage.removeItem(PHASE_KEY);
}

export function clearAllPersisted() {
  clearPersistedGame();
  localStorage.removeItem(CONFIG_KEY);
}
