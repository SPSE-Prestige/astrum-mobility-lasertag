"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { GameMode, GameState } from "@/types/game";
import { gameApi } from "@/lib/api";
import { toPlayer, toTeam, buildKillFeed, calcTeamResults } from "@/lib/game";
import { useWebSocket } from "@/hooks/useWebSocketV2";

const POLL_INTERVAL_MS = 2000;

interface UseGamePollingOptions {
  phase: GameState["phase"];
  gameId: string | null;
  gameMode: GameMode;
  durationMinutes: number;
  raceStatus: GameState["raceStatus"];
  setState: React.Dispatch<React.SetStateAction<GameState>>;
}

/**
 * Handles game state polling and WebSocket triggers during live gameplay.
 * Polls every 2s during live phase and triggers immediate poll on WS events.
 */
export function useGamePolling({
  phase,
  gameId,
  gameMode,
  durationMinutes,
  raceStatus,
  setState,
}: UseGamePollingOptions) {
  const [wsStatus, setWsStatus] = useState<"connecting" | "connected" | "disconnected">("disconnected");
  const pollRef = useRef<ReturnType<typeof setInterval> | undefined>(undefined);
  const pollGameStateRef = useRef<(() => Promise<void>) | undefined>(undefined);

  const pollGameState = useCallback(async () => {
    if (!gameId) return;
    try {
      const full = await gameApi.getFull(gameId);
      const players = full.players.map((p) => toPlayer(p, full.teams));
      const teams = full.teams.map(toTeam);
      const killFeed = buildKillFeed(full.events);
      const teamResults = calcTeamResults(players, gameMode);

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

      let raceTimeSeconds = durationMinutes * 60;
      if (full.game.started_at && full.game.settings.game_duration > 0) {
        const elapsed = (Date.now() - new Date(full.game.started_at).getTime()) / 1000;
        raceTimeSeconds = Math.max(0, full.game.settings.game_duration - Math.floor(elapsed));
      }

      setState((prev) => ({ ...prev, raceTimeSeconds, players, teams, killFeed, teamResults }));
    } catch (err) {
      console.error("[Poll]", err);
    }
  }, [gameId, gameMode, durationMinutes, setState]);

  useEffect(() => {
    pollGameStateRef.current = pollGameState;
  }, [pollGameState]);

  // Stable polling interval
  useEffect(() => {
    if (phase !== "live" || !gameId) {
      clearInterval(pollRef.current);
      return;
    }
    pollGameStateRef.current?.();
    pollRef.current = setInterval(() => pollGameStateRef.current?.(), POLL_INTERVAL_MS);
    return () => clearInterval(pollRef.current);
  }, [phase, gameId]);

  // WebSocket triggers immediate poll
  useWebSocket({
    enabled: phase === "live" && raceStatus === "running",
    gameId,
    onEvent: useCallback(() => {
      pollGameStateRef.current?.();
    }, []),
    onStatusChange: setWsStatus,
  });

  return { wsStatus, triggerPoll: () => pollGameStateRef.current?.() };
}
