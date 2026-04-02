import { describe, it, expect, beforeEach } from "vitest";
import {
  defaultConfig,
  initialState,
  loadPersistedState,
  loadPersistedConfig,
  persistGameState,
  persistConfig,
  clearPersistedGame,
  clearAllPersisted,
} from "@/lib/game/persistence";

beforeEach(() => {
  localStorage.clear();
});

describe("defaultConfig", () => {
  it("has sensible defaults", () => {
    expect(defaultConfig.durationMinutes).toBe(5);
    expect(defaultConfig.gameMode).toBe("team");
    expect(defaultConfig.killsPerUpgrade).toBe(3);
    expect(defaultConfig.maxPlayers).toBe(20);
  });
});

describe("initialState", () => {
  it("starts in setup phase with no game", () => {
    expect(initialState.phase).toBe("setup");
    expect(initialState.gameId).toBeNull();
    expect(initialState.raceStatus).toBe("idle");
    expect(initialState.raceTimeSeconds).toBe(300);
  });
});

describe("persistGameState / loadPersistedState", () => {
  it("returns null when nothing is persisted", () => {
    expect(loadPersistedState()).toBeNull();
  });

  it("persists and loads game state", () => {
    persistGameState("game-123", "players");
    const state = loadPersistedState();
    expect(state).toEqual({
      gameId: "game-123",
      phase: "players",
      raceStatus: "idle",
    });
  });

  it("sets raceStatus to running when phase is live", () => {
    persistGameState("game-123", "live");
    const state = loadPersistedState();
    expect(state?.raceStatus).toBe("running");
  });

  it("falls back to setup for invalid phase", () => {
    localStorage.setItem("lasertag-gameId", "game-123");
    localStorage.setItem("lasertag-phase", "INVALID");
    const state = loadPersistedState();
    expect(state?.phase).toBe("setup");
  });

  it("clears persistence when gameId is null", () => {
    persistGameState("game-123", "players");
    persistGameState(null, "setup");
    expect(loadPersistedState()).toBeNull();
  });
});

describe("persistConfig / loadPersistedConfig", () => {
  it("returns null when nothing is persisted", () => {
    expect(loadPersistedConfig()).toBeNull();
  });

  it("persists and loads config", () => {
    persistConfig(defaultConfig);
    const config = loadPersistedConfig();
    expect(config).toEqual(defaultConfig);
  });

  it("returns null for corrupted JSON", () => {
    localStorage.setItem("lasertag-config", "NOT_JSON{{{");
    expect(loadPersistedConfig()).toBeNull();
  });
});

describe("clearPersistedGame", () => {
  it("removes game id and phase but keeps config", () => {
    persistGameState("game-123", "players");
    persistConfig(defaultConfig);
    clearPersistedGame();
    expect(loadPersistedState()).toBeNull();
    expect(loadPersistedConfig()).toEqual(defaultConfig);
  });
});

describe("clearAllPersisted", () => {
  it("removes everything", () => {
    persistGameState("game-123", "players");
    persistConfig(defaultConfig);
    clearAllPersisted();
    expect(loadPersistedState()).toBeNull();
    expect(loadPersistedConfig()).toBeNull();
  });
});
