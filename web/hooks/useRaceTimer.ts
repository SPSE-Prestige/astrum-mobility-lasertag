"use client";

import { useEffect, useRef } from "react";
import type { GameState } from "@/types/game";

/**
 * Manages the visual countdown timer during a live game.
 * Decrements `raceTimeSeconds` every second when the race is running.
 */
export function useRaceTimer(
  phase: GameState["phase"],
  raceStatus: GameState["raceStatus"],
  setState: React.Dispatch<React.SetStateAction<GameState>>,
) {
  const timerRef = useRef<ReturnType<typeof setInterval> | undefined>(undefined);

  useEffect(() => {
    if (phase !== "live" || raceStatus !== "running") {
      clearInterval(timerRef.current);
      return;
    }

    timerRef.current = setInterval(() => {
      setState((prev) => {
        if (prev.raceTimeSeconds <= 0) return prev;
        return { ...prev, raceTimeSeconds: prev.raceTimeSeconds - 1 };
      });
    }, 1000);

    return () => clearInterval(timerRef.current);
  }, [phase, raceStatus, setState]);
}
