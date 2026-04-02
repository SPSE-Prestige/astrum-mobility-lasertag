"use client";

import { useEffect, useState } from "react";
import type { GameMode, GamePhase, MatchHistoryItem } from "@/types/game";
import { gameApi } from "@/lib/api";

/**
 * Fetches match history when the results phase is active.
 */
export function useMatchHistory(phase: GamePhase, isAuthenticated: boolean) {
  const [matchHistory, setMatchHistory] = useState<MatchHistoryItem[]>([]);

  useEffect(() => {
    if (phase !== "results" || !isAuthenticated) return;

    gameApi
      .list()
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
  }, [phase, isAuthenticated]);

  return matchHistory;
}
