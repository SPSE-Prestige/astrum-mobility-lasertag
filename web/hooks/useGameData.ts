"use client";

import { useEffect, useMemo, useState } from "react";
import type { AuthState, GameConfig, GameMode, GamePhase, GameState, MatchHistoryItem, Player, RegisteredPlayer } from "@/types/game";

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
    cartSpeed: 28,
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
      message: "event.systemReady",
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
const PLAYERS_STORAGE_KEY = "race-control-registered-players";

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

  const [registeredPlayers, setRegisteredPlayers] = useState<RegisteredPlayer[]>(() => {
    if (typeof window === "undefined") {
      return [];
    }

    const stored = window.localStorage.getItem(PLAYERS_STORAGE_KEY);
    if (!stored) {
      return [];
    }

    try {
      return JSON.parse(stored) as RegisteredPlayer[];
    } catch {
      return [];
    }
  });

  const [activeRoster, setActiveRoster] = useState<string[]>([]);
  // We want to persist this roster or load from current players on init maybe?
  // For now, let's keep it simple.

  useEffect(() => {
    localStorage.setItem(PLAYERS_STORAGE_KEY, JSON.stringify(registeredPlayers));
  }, [registeredPlayers]);

  const toggleRosterPlayer = (playerId: string) => {
    setActiveRoster((prev) => {
      const next = prev.includes(playerId) ? prev.filter((id) => id !== playerId) : [...prev, playerId];
      return next;
    });
  };

  // Sync state.players with activeRoster
  useEffect(() => {
    // Only update players list if we are in setup or players management phase
    if (state.phase !== "setup" && state.phase !== "players") return;

    if (activeRoster.length === 0) {
      // If roster cleared, maybe we don't want to show empty list? Or maybe we do.
      // Let's allow empty list for now.
      if (state.players.length > 0 && state.players.every((p) => registeredPlayers.some((rp) => rp.id === p.id))) {
         setState((prev) => ({ ...prev, players: [] }));
      }
      return;
    }

    const currentPlayersMap = new Map(state.players.map((p) => [p.id, p]));

    const newPlayersList = activeRoster
      .map((rosterId) => {
        const registered = registeredPlayers.find((rp) => rp.id === rosterId);
        if (!registered) return null;

        const existing = currentPlayersMap.get(rosterId);
        if (existing) return existing;

        return {
          id: registered.id,
          name: registered.name,
          team: "Unassigned",
          hp: config.playerTuning.hp,
          ammo: 60,
          status: "active",
          kills: 0,
          hits: 0,
          accuracy: 0,
          damageDealt: 0,
          cartConnected: true,
        } as Player;
      })
      .filter((p): p is Player => p !== null);

    setState((prev) => ({
      ...prev,
      players: newPlayersList,
    }));
  }, [activeRoster, registeredPlayers, config.playerTuning.hp]); // Removed state.phase to avoid loop, handled by check inside

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
                  message: `event.tag:${randomKiller.name}:${randomVictim.name}`,
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
    setState((prev) => {
      if (phase === prev.phase) {
        return prev;
      }

      // Manual switch to live phase is not allowed (only via start button).
      if (phase === "live") {
        return prev;
      }

      if (prev.raceStatus === "running") {
        return prev;
      }

      if (phase === "setup" && prev.phase === "results") {
        // Returning from results to setup prepares a fresh game state.
        return {
          ...initialState,
          phase: "setup",
          raceTimeSeconds: config.durationMinutes * 60,
          teamResults: calculateTeamResults(initialPlayers, config.gameMode),
        };
      }

      const allowedTransitions: Record<GamePhase, GamePhase[]> = {
        setup: ["players", "results"],
        players: ["setup", "results"],
        results: ["setup", "players"],
        live: [],
      };

      if (allowedTransitions[prev.phase].includes(phase)) {
        return {
          ...prev,
          phase,
        };
      }

      return prev;
    });
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
    setState((prev) => ({
      ...prev,
      raceStatus: prev.raceStatus === "paused" ? "running" : "paused",
    }));
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
        error: "invalid_credentials",
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

  const registerPlayer = (name: string, type: "guest" | "registered" = "registered") => {
    const SUPERHERO_NAMES = [
      "Hulk",
      "Thor",
      "Iron Man",
      "Spiderman",
      "Batman",
      "Superman",
      "Wonder Woman",
      "Flash",
      "Black Widow",
      "Captain America",
      "Deadpool",
      "Doctor Strange",
      "Black Panther",
      "Vision",
      "Scarlet Witch",
    ];

    let finalName = name.trim();
    if (type === "guest") {
      finalName = SUPERHERO_NAMES[Math.floor(Math.random() * SUPERHERO_NAMES.length)];
    }

    if (!finalName) return;

    const code = Math.floor(1000 + Math.random() * 9000).toString();
    const newPlayer: RegisteredPlayer = {
      id: crypto.randomUUID(),
      name: finalName,
      code,
      type,
      createdAt: new Date().toISOString(),
    };
    setRegisteredPlayers((prev) => [...prev, newPlayer]);
  };

  const deleteRegisteredPlayer = (id: string) => {
    setRegisteredPlayers((prev) => prev.filter((p) => p.id !== id));
  };

  return {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
    registeredPlayers,
    activeRoster,
    registerPlayer,
    deleteRegisteredPlayer,
    toggleRosterPlayer,
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