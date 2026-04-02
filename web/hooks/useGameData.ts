"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { GameConfig, GamePhase, GameState } from "@/types/game";
import { deviceApi, gameApi, type TeamResponse } from "@/lib/api/index";
import {
  defaultConfig,
  initialState,
  loadPersistedState,
  loadPersistedConfig,
  persistGameState,
  persistConfig,
  clearPersistedGame,
  clearAllPersisted,
  toPlayer,
  toTeam,
  toDevice,
  buildKillFeed,
  calcTeamResults,
  formatRaceTime,
} from "@/lib/game";
import { useAuth } from "@/hooks/useAuth";
import { useRaceTimer } from "@/hooks/useRaceTimer";
import { useMatchHistory } from "@/hooks/useMatchHistory";
import { useGamePolling } from "@/hooks/useGamePolling";

// ── Loading / error helpers ──

function useActionState() {
  const [loading, setLoading] = useState<Record<string, boolean>>({});
  const [error, setError] = useState<string | null>(null);

  const setActionLoading = useCallback((key: string, value: boolean) => {
    setLoading((prev) => ({ ...prev, [key]: value }));
  }, []);

  const clearError = useCallback(() => setError(null), []);

  return { loading, error, setActionLoading, setError, clearError };
}

// ── Main hook (composition root) ──

export const useGameData = () => {
  const { auth, login, logout: rawLogout } = useAuth();

  const [config, setConfig] = useState<GameConfig>(() => loadPersistedConfig() ?? defaultConfig);
  const [state, setState] = useState<GameState>(() => {
    const persisted = loadPersistedState();
    if (persisted) {
      const cfg = loadPersistedConfig() ?? defaultConfig;
      return { ...initialState, ...persisted, raceTimeSeconds: cfg.durationMinutes * 60 };
    }
    return initialState;
  });

  const { loading, error, setActionLoading, setError, clearError } = useActionState();

  // ── Composed hooks ──

  useRaceTimer(state.phase, state.raceStatus, setState);

  const matchHistory = useMatchHistory(state.phase, auth.isAuthenticated);

  const { wsStatus } = useGamePolling({
    phase: state.phase,
    gameId: state.gameId,
    gameMode: config.gameMode,
    durationMinutes: config.durationMinutes,
    raceStatus: state.raceStatus,
    setState,
  });

  // ── Persistence effects ──

  useEffect(() => {
    persistGameState(state.gameId, state.phase);
  }, [state.gameId, state.phase]);

  useEffect(() => {
    persistConfig(config);
  }, [config]);

  // ── Recover game state on mount ──

  const hasRecoveredRef = useRef(false);
  useEffect(() => {
    if (hasRecoveredRef.current || !state.gameId || !auth.isAuthenticated) return;
    hasRecoveredRef.current = true;

    (async () => {
      try {
        const full = await gameApi.getFull(state.gameId!);
        const players = full.players.map((p) => toPlayer(p, full.teams));
        const teams = full.teams.map(toTeam);
        const killFeed = buildKillFeed(full.events);
        const teamResults = calcTeamResults(players, config.gameMode);

        if (full.game.status === "finished") {
          setState((prev) => ({ ...prev, phase: "results", raceStatus: "finished", players, teams, killFeed, teamResults }));
          return;
        }

        const isRunning = full.game.status === "running" || full.game.status === "started";
        let raceTimeSeconds = config.durationMinutes * 60;
        if (isRunning && full.game.started_at && full.game.settings.game_duration > 0) {
          const elapsed = (Date.now() - new Date(full.game.started_at).getTime()) / 1000;
          raceTimeSeconds = Math.max(0, full.game.settings.game_duration - Math.floor(elapsed));
        }

        setState((prev) => ({
          ...prev,
          phase: isRunning ? "live" : prev.phase === "live" ? "players" : prev.phase,
          raceStatus: isRunning ? "running" : "idle",
          raceTimeSeconds,
          players,
          teams,
          killFeed,
          teamResults,
        }));
      } catch {
        clearPersistedGame();
        setState(initialState);
      }
    })();
  }, [state.gameId, auth.isAuthenticated, config.gameMode, config.durationMinutes]);

  // ── Computed ──

  const leaderboard = useMemo(
    () => [...state.players].sort((a, b) => b.score - a.score || b.kills - a.kills),
    [state.players],
  );

  // ── Actions ──

  const updateConfig = useCallback((next: GameConfig) => {
    setConfig(next);
    setState((prev) => ({ ...prev, raceTimeSeconds: next.durationMinutes * 60 }));
  }, []);

  const logout = useCallback(async () => {
    await rawLogout();
    setState(initialState);
    clearAllPersisted();
  }, [rawLogout]);

  const createGame = useCallback(async () => {
    setActionLoading("createGame", true);
    setError(null);
    try {
      const game = await gameApi.create({
        respawn_delay: config.respawnDelay,
        game_duration: config.durationMinutes * 60,
        friendly_fire: config.friendlyFire,
        max_players: config.maxPlayers,
        score_per_kill: 100,
        kills_per_upgrade: config.killsPerUpgrade,
      });

      let teams: TeamResponse[] = [];
      if (config.gameMode === "team") {
        const defs = [
          { name: "Neon Green", color: "#00ff00" },
          { name: "Neon Red", color: "#ff0000" },
          { name: "Neon Blue", color: "#0088ff" },
          { name: "Neon Yellow", color: "#ffff00" },
        ];
        for (let i = 0; i < config.teamsCount; i++) {
          await gameApi.addTeam(game.id, defs[i].name, defs[i].color);
        }
        teams = await gameApi.listTeams(game.id);
      }

      const devices = await deviceApi.listAvailable();

      setState((prev) => ({
        ...prev,
        gameId: game.id,
        phase: "players",
        teams: teams.map(toTeam),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create game");
    } finally {
      setActionLoading("createGame", false);
    }
  }, [config, setActionLoading, setError]);

  const refreshDevices = useCallback(async () => {
    setActionLoading("refreshDevices", true);
    try {
      const devices = await deviceApi.listAvailable();
      setState((prev) => ({ ...prev, devices: devices.map(toDevice) }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to refresh devices");
    } finally {
      setActionLoading("refreshDevices", false);
    }
  }, [setActionLoading, setError]);

  const addPlayer = useCallback(async (deviceId: string, nickname: string, teamId?: string) => {
    const gameId = state.gameId;
    if (!gameId) return;
    setActionLoading("addPlayer", true);
    setError(null);
    try {
      await gameApi.addPlayer(gameId, deviceId, nickname, teamId);
      const [players, devices] = await Promise.all([
        gameApi.listPlayers(gameId),
        deviceApi.listAvailable(),
      ]);
      const teamResponses = state.teams.length > 0 ? await gameApi.listTeams(gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add player");
    } finally {
      setActionLoading("addPlayer", false);
    }
  }, [state.gameId, state.teams.length, setActionLoading, setError]);

  const removePlayer = useCallback(async (playerId: string) => {
    const gameId = state.gameId;
    if (!gameId) return;
    setActionLoading("removePlayer", true);
    try {
      await gameApi.removePlayer(gameId, playerId);
      const [players, devices] = await Promise.all([
        gameApi.listPlayers(gameId),
        deviceApi.listAvailable(),
      ]);
      const teamResponses = state.teams.length > 0 ? await gameApi.listTeams(gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to remove player");
    } finally {
      setActionLoading("removePlayer", false);
    }
  }, [state.gameId, state.teams.length, setActionLoading, setError]);

  const assignPlayerTeam = useCallback(async (playerId: string, teamId: string | null) => {
    const gameId = state.gameId;
    if (!gameId) return;
    try {
      await gameApi.updatePlayerTeam(gameId, playerId, teamId);
      const players = await gameApi.listPlayers(gameId);
      const teamResponses = state.teams.length > 0 ? await gameApi.listTeams(gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to assign team");
    }
  }, [state.gameId, state.teams.length, setError]);

  const startRace = useCallback(async () => {
    const gameId = state.gameId;
    if (!gameId) return;
    setActionLoading("startRace", true);
    setError(null);
    try {
      await gameApi.start(gameId);
      setState((prev) => ({
        ...prev,
        phase: "live",
        raceStatus: "running",
        raceTimeSeconds: config.durationMinutes * 60,
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to start game");
    } finally {
      setActionLoading("startRace", false);
    }
  }, [state.gameId, config.durationMinutes, setActionLoading, setError]);

  const stopRace = useCallback(async () => {
    const gameId = state.gameId;
    if (!gameId) return;
    setActionLoading("stopRace", true);
    try {
      await gameApi.end(gameId);
      const full = await gameApi.getFull(gameId);
      const players = full.players.map((p) => toPlayer(p, full.teams));
      setState((prev) => ({
        ...prev,
        phase: "results",
        raceStatus: "finished",
        players,
        teamResults: calcTeamResults(players, config.gameMode),
        killFeed: buildKillFeed(full.events),
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to stop game");
    } finally {
      setActionLoading("stopRace", false);
    }
  }, [state.gameId, config.gameMode, setActionLoading, setError]);

  const updatePhase = useCallback((phase: GamePhase) => {
    setState((prev) => {
      if (phase === prev.phase) return prev;
      if (phase === "live") return prev;
      if (prev.raceStatus === "running") return prev;
      if (phase === "players" && !prev.gameId) return prev;

      if (phase === "setup" && prev.phase === "results") {
        clearPersistedGame();
        return { ...initialState, phase: "setup", raceTimeSeconds: config.durationMinutes * 60 };
      }

      const allowed: Record<GamePhase, GamePhase[]> = {
        setup: ["players", "results"],
        players: ["setup", "results"],
        results: ["setup", "players"],
        live: [],
      };
      if (allowed[prev.phase].includes(phase)) {
        return { ...prev, phase };
      }
      return prev;
    });
  }, [config.durationMinutes]);

  return {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
    wsStatus,
    loading,
    error,
    clearError,
    formatRaceTime,
    updateConfig,
    updatePhase,
    createGame,
    addPlayer,
    removePlayer,
    assignPlayerTeam,
    refreshDevices,
    startRace,
    stopRace,
    login,
    logout,
  };
};
