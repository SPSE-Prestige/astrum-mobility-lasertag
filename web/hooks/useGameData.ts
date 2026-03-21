"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type {
  AuthState,
  Device,
  GameConfig,
  GameMode,
  GamePhase,
  GameState,
  KillFeedEvent,
  MatchHistoryItem,
  Player,
  Team,
  TeamResult,
} from "@/types/game";
import {
  api,
  type EventResponse,
  type PlayerResponse,
  type TeamResponse,
} from "@/lib/api";
import { useWebSocket } from "@/hooks/useWebSocket";

const TOKEN_KEY = "lasertag-token";
const USER_KEY = "lasertag-user";

const defaultConfig: GameConfig = {
  gameName: "Laser Tag Game",
  durationMinutes: 5,
  teamsCount: 2,
  gameMode: "team",
  friendlyFire: false,
  respawnDelay: 5,
  maxPlayers: 20,
};

const initialState: GameState = {
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

const initialAuth: AuthState = {
  isAuthenticated: false,
  username: null,
  token: null,
  error: null,
};

// ── Helpers ──

function toPlayer(p: PlayerResponse, teams: TeamResponse[]): Player {
  const team = teams.find((t) => t.id === p.team_id);
  return {
    id: p.id,
    name: p.nickname,
    team: team?.name ?? "Unassigned",
    teamId: p.team_id ?? null,
    deviceId: p.device_id,
    status: p.is_alive ? "alive" : "dead",
    kills: p.kills,
    deaths: p.deaths,
    score: p.score,
  };
}

function toTeam(t: TeamResponse): Team {
  return { id: t.id, name: t.name, color: t.color };
}

function toDevice(d: { id: string; device_id: string; status: string; last_seen: string }): Device {
  return { id: d.id, deviceId: d.device_id, status: d.status, lastSeen: d.last_seen };
}

function buildKillFeed(events: EventResponse[]): KillFeedEvent[] {
  return events
    .filter((e) => e.type === "kill")
    .reverse()
    .slice(0, 20)
    .map((e) => ({
      id: e.id,
      timestamp: new Date(e.timestamp).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      }),
      message: `${(e.payload?.attacker_nickname as string) ?? "?"} → ${(e.payload?.victim_nickname as string) ?? "?"}`,
    }));
}

function calcTeamResults(players: Player[], mode: GameMode): TeamResult[] {
  if (mode === "ffa") {
    return players
      .map((p) => ({ team: p.name, score: p.score, kills: p.kills, deaths: p.deaths }))
      .sort((a, b) => b.score - a.score);
  }
  const grouped = players
    .filter((p) => p.team !== "Unassigned")
    .reduce<Record<string, { score: number; kills: number; deaths: number }>>((acc, p) => {
      if (!acc[p.team]) acc[p.team] = { score: 0, kills: 0, deaths: 0 };
      acc[p.team].score += p.score;
      acc[p.team].kills += p.kills;
      acc[p.team].deaths += p.deaths;
      return acc;
    }, {});
  return Object.entries(grouped)
    .map(([team, s]) => ({ team, ...s }))
    .sort((a, b) => b.score - a.score);
}

const formatRaceTime = (seconds: number): string => {
  const m = Math.floor(seconds / 60).toString().padStart(2, "0");
  const s = Math.max(seconds % 60, 0).toString().padStart(2, "0");
  return `${m}:${s}`;
};

// ── Hook ──

export const useGameData = () => {
  const [config, setConfig] = useState<GameConfig>(defaultConfig);
  const [state, setState] = useState<GameState>(initialState);
  const [matchHistory, setMatchHistory] = useState<MatchHistoryItem[]>([]);
  const [auth, setAuth] = useState<AuthState>(() => {
    if (typeof window === "undefined") return initialAuth;
    const token = localStorage.getItem(TOKEN_KEY);
    const username = localStorage.getItem(USER_KEY);
    if (token && username) {
      api.setToken(token);
      return { isAuthenticated: true, username, token, error: null };
    }
    return initialAuth;
  });

  const pollRef = useRef<ReturnType<typeof setInterval> | undefined>(undefined);

  // Persist auth token
  useEffect(() => {
    if (auth.token) {
      localStorage.setItem(TOKEN_KEY, auth.token);
      localStorage.setItem(USER_KEY, auth.username ?? "");
      api.setToken(auth.token);
    } else {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      api.setToken(null);
    }
  }, [auth.token, auth.username]);

  // ── Polling during live phase ──

  const pollGameState = useCallback(async () => {
    const gameId = state.gameId;
    if (!gameId) return;
    try {
      const full = await api.getGameFull(gameId);
      const players = full.players.map((p) => toPlayer(p, full.teams));
      const teams = full.teams.map(toTeam);
      const killFeed = buildKillFeed(full.events);
      const teamResults = calcTeamResults(players, config.gameMode);

      if (full.game.status === "finished") {
        clearInterval(pollRef.current);
        setState((prev) => ({
          ...prev,
          phase: "results",
          raceStatus: "finished",
          players,
          teams,
          killFeed,
          teamResults,
        }));
        return;
      }

      let raceTimeSeconds = config.durationMinutes * 60;
      if (full.game.started_at && full.game.settings.game_duration > 0) {
        const elapsed = (Date.now() - new Date(full.game.started_at).getTime()) / 1000;
        raceTimeSeconds = Math.max(0, full.game.settings.game_duration - Math.floor(elapsed));
      }

      setState((prev) => ({
        ...prev,
        raceTimeSeconds,
        players,
        teams,
        killFeed,
        teamResults,
      }));
    } catch (err) {
      console.error("[Poll]", err);
    }
  }, [state.gameId, config.gameMode, config.durationMinutes]);

  useEffect(() => {
    if (state.phase !== "live" || !state.gameId) {
      clearInterval(pollRef.current);
      return;
    }
    pollGameState();
    pollRef.current = setInterval(pollGameState, 2000);
    return () => clearInterval(pollRef.current);
  }, [state.phase, state.gameId, pollGameState]);

  // WebSocket — triggers immediate poll on events
  useWebSocket(
    state.phase === "live" && state.raceStatus === "running",
    state.gameId,
    useCallback(() => {
      pollGameState();
    }, [pollGameState]),
  );

  // Timer countdown (visual only, server is authoritative)
  useEffect(() => {
    if (state.phase !== "live" || state.raceStatus !== "running") return;
    const timer = setInterval(() => {
      setState((prev) => {
        if (prev.raceTimeSeconds <= 0) return prev;
        return { ...prev, raceTimeSeconds: prev.raceTimeSeconds - 1 };
      });
    }, 1000);
    return () => clearInterval(timer);
  }, [state.phase, state.raceStatus]);

  // Load match history when entering results
  useEffect(() => {
    if (state.phase === "results" && auth.isAuthenticated) {
      api
        .listGames()
        .then((games) =>
          setMatchHistory(
            games
              .filter((g) => g.status === "finished")
              .map((g) => ({
                id: g.id,
                gameName: g.code,
                gameMode: "team" as GameMode,
                playedAt: g.created_at,
                durationMinutes: Math.round((g.settings?.game_duration ?? 300) / 60),
                status: g.status,
              }))
              .sort((a, b) => new Date(b.playedAt).getTime() - new Date(a.playedAt).getTime()),
          ),
        )
        .catch(() => {});
    }
  }, [state.phase, auth.isAuthenticated]);

  const leaderboard = useMemo(
    () => [...state.players].sort((a, b) => b.score - a.score || b.kills - a.kills),
    [state.players],
  );

  // ─── Actions ───

  const login = async (username: string, password: string): Promise<boolean> => {
    try {
      const res = await api.login(username, password);
      setAuth({ isAuthenticated: true, username, token: res.token, error: null });
      return true;
    } catch {
      setAuth({ isAuthenticated: false, username: null, token: null, error: "invalid_credentials" });
      return false;
    }
  };

  const logout = async () => {
    try {
      await api.logout();
    } catch {}
    setAuth(initialAuth);
    setState(initialState);
  };

  const updateConfig = (next: GameConfig) => {
    setConfig(next);
    setState((prev) => ({ ...prev, raceTimeSeconds: next.durationMinutes * 60 }));
  };

  const createGame = async () => {
    try {
      const game = await api.createGame({
        respawn_delay: config.respawnDelay,
        game_duration: config.durationMinutes * 60,
        friendly_fire: config.friendlyFire,
        max_players: config.maxPlayers,
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
          await api.addTeam(game.id, defs[i].name, defs[i].color);
        }
        teams = await api.listTeams(game.id);
      }

      const devices = await api.listAvailableDevices();

      setState((prev) => ({
        ...prev,
        gameId: game.id,
        phase: "players",
        teams: teams.map(toTeam),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      console.error("[CreateGame]", err);
    }
  };

  const refreshDevices = async () => {
    try {
      const devices = await api.listAvailableDevices();
      setState((prev) => ({ ...prev, devices: devices.map(toDevice) }));
    } catch {}
  };

  const addPlayer = async (deviceId: string, nickname: string, teamId?: string) => {
    if (!state.gameId) return;
    try {
      await api.addPlayer(state.gameId, deviceId, nickname, teamId);
      const [players, devices] = await Promise.all([
        api.listPlayers(state.gameId),
        api.listAvailableDevices(),
      ]);
      const teamResponses = state.teams.length > 0 ? await api.listTeams(state.gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      console.error("[AddPlayer]", err);
    }
  };

  const removePlayer = async (playerId: string) => {
    if (!state.gameId) return;
    try {
      await api.removePlayer(state.gameId, playerId);
      const [players, devices] = await Promise.all([
        api.listPlayers(state.gameId),
        api.listAvailableDevices(),
      ]);
      const teamResponses = state.teams.length > 0 ? await api.listTeams(state.gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
        devices: devices.map(toDevice),
      }));
    } catch (err) {
      console.error("[RemovePlayer]", err);
    }
  };

  const assignPlayerTeam = async (playerId: string, teamId: string | null) => {
    if (!state.gameId) return;
    try {
      await api.updatePlayerTeam(state.gameId, playerId, teamId);
      const players = await api.listPlayers(state.gameId);
      const teamResponses = state.teams.length > 0 ? await api.listTeams(state.gameId) : [];
      setState((prev) => ({
        ...prev,
        players: players.map((p) => toPlayer(p, teamResponses)),
      }));
    } catch (err) {
      console.error("[AssignTeam]", err);
    }
  };

  const startRace = async () => {
    if (!state.gameId) return;
    try {
      await api.startGame(state.gameId);
      setState((prev) => ({
        ...prev,
        phase: "live",
        raceStatus: "running",
        raceTimeSeconds: config.durationMinutes * 60,
      }));
    } catch (err) {
      console.error("[StartRace]", err);
    }
  };

  const stopRace = async () => {
    if (!state.gameId) return;
    try {
      await api.endGame(state.gameId);
      const full = await api.getGameFull(state.gameId);
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
      console.error("[StopRace]", err);
    }
  };

  const updatePhase = (phase: GamePhase) => {
    setState((prev) => {
      if (phase === prev.phase) return prev;
      if (phase === "live") return prev;
      if (prev.raceStatus === "running") return prev;
      if (phase === "players" && !prev.gameId) return prev;

      if (phase === "setup" && prev.phase === "results") {
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
  };

  return {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
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