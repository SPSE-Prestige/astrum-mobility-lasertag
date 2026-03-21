"use client";

import { useEffect, useMemo, useState } from "react";
import type { AuthState, GameConfig, GameMode, GamePhase, GameState, MatchHistoryItem, Player } from "@/types/game";

const initialConfig: GameConfig = {
  gameName: "Aimtec Night Grand Prix",
  durationMinutes: 12,
  teamsCount: 2,
  gameMode: "team",
  friendlyFire: false,
  weaponTuning: {
    damage: 24,
    fireRate: 6,
    reloadTime: 2,
  },
  playerTuning: {
    hp: 100,
    respawnDelay: 4,
  },
};

const initialPlayers: Player[] = [
  {
    id: "p1",
    name: "Rider 01",
    team: "Neon Green",
    hp: 100,
    ammo: 36,
    status: "active",
    kills: 2,
    hits: 9,
    accuracy: 56,
    damageDealt: 220,
    cartConnected: true,
  },
  {
    id: "p2",
    name: "Rider 02",
    team: "Neon Red",
    hp: 90,
    ammo: 29,
    status: "active",
    kills: 3,
    hits: 11,
    accuracy: 62,
    damageDealt: 265,
    cartConnected: true,
  },
  {
    id: "p3",
    name: "Rider 03",
    team: "Neon Green",
    hp: 70,
    ammo: 16,
    status: "stunned",
    kills: 1,
    hits: 7,
    accuracy: 40,
    damageDealt: 155,
    cartConnected: true,
  },
  {
    id: "p4",
    name: "Rider 04",
    team: "Unassigned",
    hp: 100,
    ammo: 40,
    status: "active",
    kills: 0,
    hits: 2,
    accuracy: 20,
    damageDealt: 40,
    cartConnected: false,
  },
];

const initialState: GameState = {
  phase: "setup",
  raceTimeSeconds: initialConfig.durationMinutes * 60,
  raceStatus: "idle",
  players: initialPlayers,
  killFeed: [
    {
      id: "k1",
      timestamp: "00:00",
      message: "System ready. Waiting for start signal.",
    },
  ],
  teamResults: [
    {
      team: "Neon Green",
      score: 0,
      accuracy: 0,
      damageDealt: 0,
    },
    {
      team: "Neon Red",
      score: 0,
      accuracy: 0,
      damageDealt: 0,
    },
  ],
};

const HISTORY_STORAGE_KEY = "race-control-history";

const initialAuthState: AuthState = {
  isAuthenticated: false,
  username: null,
  error: null,
};

const seedMatchHistory: MatchHistoryItem[] = [
  {
    id: "m1",
    gameName: "Friday Neon Sprint",
    gameMode: "team",
    playedAt: "2026-03-20T18:30:00.000Z",
    durationMinutes: 10,
    winner: "Neon Green",
    totalKills: 18,
  },
  {
    id: "m2",
    gameName: "Solo Chaos Run",
    gameMode: "ffa",
    playedAt: "2026-03-20T20:00:00.000Z",
    durationMinutes: 8,
    winner: "Rider 02",
    totalKills: 14,
  },
];

const calculateTeamResults = (players: Player[], mode: GameMode) => {
  if (mode === "ffa") {
    return players
      .filter((player) => player.cartConnected)
      .map((player) => ({
        team: player.name,
        score: player.kills * 10 + player.hits,
        accuracy: player.accuracy,
        damageDealt: player.damageDealt,
      }))
      .sort((a, b) => b.score - a.score);
  }

  const grouped = players
    .filter((player) => player.team !== "Unassigned")
    .reduce<Record<string, { kills: number; hits: number; damage: number; accuracySum: number; count: number }>>(
      (acc, player) => {
        if (!acc[player.team]) {
          acc[player.team] = { kills: 0, hits: 0, damage: 0, accuracySum: 0, count: 0 };
        }

        acc[player.team].kills += player.kills;
        acc[player.team].hits += player.hits;
        acc[player.team].damage += player.damageDealt;
        acc[player.team].accuracySum += player.accuracy;
        acc[player.team].count += 1;
        return acc;
      },
      {},
    );

  return Object.entries(grouped)
    .map(([team, metrics]) => ({
      team,
      score: metrics.kills * 10 + metrics.hits,
      accuracy: metrics.count > 0 ? Math.round(metrics.accuracySum / metrics.count) : 0,
      damageDealt: metrics.damage,
    }))
    .sort((a, b) => b.score - a.score);
};

const resolveWinner = (players: Player[], mode: GameMode): string => {
  if (mode === "ffa") {
    const bestPlayer = [...players].sort((a, b) => b.kills - a.kills || b.hits - a.hits)[0];
    return bestPlayer?.name ?? "Unknown";
  }

  const teamResults = calculateTeamResults(players, mode);
  return teamResults[0]?.team ?? "Unknown";
};

const totalKills = (players: Player[]): number => players.reduce((sum, player) => sum + player.kills, 0);

const formatRaceTime = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
    .toString()
    .padStart(2, "0");
  const secs = Math.max(seconds % 60, 0)
    .toString()
    .padStart(2, "0");
  return `${mins}:${secs}`;
};

export const useGameData = () => {
  const [config, setConfig] = useState<GameConfig>(initialConfig);
  const [state, setState] = useState<GameState>(initialState);
  const [auth, setAuth] = useState<AuthState>(initialAuthState);
  const [matchHistory, setMatchHistory] = useState<MatchHistoryItem[]>(() => {
    if (typeof window === "undefined") {
      return seedMatchHistory;
    }

    const stored = window.localStorage.getItem(HISTORY_STORAGE_KEY);
    if (!stored) {
      return seedMatchHistory;
    }

    try {
      const parsed = JSON.parse(stored) as MatchHistoryItem[];
      return parsed.length > 0 ? parsed : seedMatchHistory;
    } catch {
      return seedMatchHistory;
    }
  });

  const appendHistoryEntry = (players: Player[], mode: GameMode, currentConfig: GameConfig) => {
    setMatchHistory((prev) => [
      {
        id: crypto.randomUUID(),
        gameName: currentConfig.gameName,
        gameMode: mode,
        playedAt: new Date().toISOString(),
        durationMinutes: currentConfig.durationMinutes,
        winner: resolveWinner(players, mode),
        totalKills: totalKills(players),
      },
      ...prev,
    ]);
  };

  useEffect(() => {
    localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(matchHistory));
  }, [matchHistory]);

  useEffect(() => {
    const tick = setInterval(() => {
      setState((prev) => {
        if (prev.phase !== "live" || prev.raceStatus !== "running") {
          return prev;
        }

        const nextTime = prev.raceTimeSeconds > 0 ? prev.raceTimeSeconds - 1 : 0;

        const updatedPlayers: Player[] = prev.players.map((player): Player => {
          if (!player.cartConnected) {
            return player;
          }

          const hpShift = Math.random() > 0.7 ? Math.floor(Math.random() * 12) : 0;
          const ammoShift = Math.random() > 0.45 ? Math.floor(Math.random() * 3) : 0;
          const nextHp = Math.max(player.hp - hpShift, 0);
          const stunned = nextHp <= 20 && Math.random() > 0.5;

          return {
            ...player,
            hp: nextHp,
            ammo: Math.max(player.ammo - ammoShift, 0),
            status: stunned ? "stunned" : "active",
            hits: player.hits + (Math.random() > 0.55 ? 1 : 0),
            kills: player.kills + (Math.random() > 0.9 ? 1 : 0),
            damageDealt: player.damageDealt + hpShift,
            accuracy: Math.min(99, Math.max(15, player.accuracy + (Math.random() > 0.5 ? 1 : -1))),
          };
        });

        const randomKiller = updatedPlayers[Math.floor(Math.random() * updatedPlayers.length)];
        const randomVictim = updatedPlayers[Math.floor(Math.random() * updatedPlayers.length)];

        const nextFeed =
          Math.random() > 0.7 && randomKiller && randomVictim && randomKiller.id !== randomVictim.id
            ? [
                {
                  id: crypto.randomUUID(),
                  timestamp: formatRaceTime(config.durationMinutes * 60 - nextTime),
                  message: `${randomKiller.name} tagged ${randomVictim.name}`,
                },
                ...prev.killFeed,
              ].slice(0, 8)
            : prev.killFeed;

        const nextStatus = nextTime === 0 ? "finished" : prev.raceStatus;
        const nextPhase: GamePhase = nextTime === 0 ? "results" : prev.phase;

        if (nextStatus === "finished") {
          appendHistoryEntry(updatedPlayers, config.gameMode, config);
        }

        return {
          ...prev,
          raceTimeSeconds: nextTime,
          raceStatus: nextStatus,
          phase: nextPhase,
          players: updatedPlayers,
          killFeed: nextFeed,
          teamResults: calculateTeamResults(updatedPlayers, config.gameMode),
        };
      });
    }, 1000);

    return () => clearInterval(tick);
  }, [config]);

  const leaderboard = useMemo(
    () =>
      [...state.players]
        .filter((player) => config.gameMode === "ffa" || player.team !== "Unassigned")
        .sort((a, b) => b.kills - a.kills || b.hits - a.hits),
    [config.gameMode, state.players],
  );

  const updatePhase = (phase: GamePhase) => {
    setState((prev) => ({ ...prev, phase }));
  };

  const updateConfig = (nextConfig: GameConfig) => {
    setConfig(nextConfig);
    setState((prev) => ({
      ...prev,
      raceTimeSeconds: nextConfig.durationMinutes * 60,
      teamResults: calculateTeamResults(prev.players, nextConfig.gameMode),
    }));
  };

  const assignPlayerTeam = (playerId: string, team: string) => {
    setState((prev) => {
      const updatedPlayers = prev.players.map((player) =>
        player.id === playerId
          ? {
              ...player,
              team,
            }
          : player,
      );

      return {
        ...prev,
        players: updatedPlayers,
        teamResults: calculateTeamResults(updatedPlayers, config.gameMode),
      };
    });
  };

  const startRace = () => {
    setState((prev) => ({ ...prev, phase: "live", raceStatus: "running" }));
  };

  const pauseRace = () => {
    setState((prev) => ({ ...prev, raceStatus: "paused" }));
  };

  const stopRace = () => {
    if (state.raceStatus === "finished") {
      return;
    }

    appendHistoryEntry(state.players, config.gameMode, config);

    setState((prev) => ({
      ...prev,
      raceStatus: "finished",
      phase: "results",
      teamResults: calculateTeamResults(prev.players, config.gameMode),
    }));
  };

  const login = (username: string, password: string) => {
    const valid = username.trim() === "admin" && password === "admin123";

    if (!valid) {
      setAuth({
        isAuthenticated: false,
        username: null,
        error: "Neplatné uživatelské jméno nebo heslo.",
      });
      return false;
    }

    setAuth({
      isAuthenticated: true,
      username: username.trim(),
      error: null,
    });
    return true;
  };

  const logout = () => {
    setAuth(initialAuthState);
  };

  const updateMatchHistoryItem = (
    matchId: string,
    patch: Partial<Pick<MatchHistoryItem, "gameName" | "winner" | "durationMinutes" | "totalKills">>,
  ) => {
    setMatchHistory((prev) => prev.map((item) => (item.id === matchId ? { ...item, ...patch } : item)));
  };

  const deleteMatchHistoryItem = (matchId: string) => {
    setMatchHistory((prev) => prev.filter((item) => item.id !== matchId));
  };

  return {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
    updateConfig,
    updatePhase,
    assignPlayerTeam,
    startRace,
    pauseRace,
    stopRace,
    login,
    logout,
    updateMatchHistoryItem,
    deleteMatchHistoryItem,
  };
};