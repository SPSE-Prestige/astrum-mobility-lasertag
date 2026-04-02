"use client";

import { Square } from "lucide-react";
import type { Language } from "@/types/i18n";

interface GameControlsProps {
  raceTime: string;
  raceStatus: "idle" | "running" | "finished";
  language: Language;
  onStop: () => void;
}

export const GameControls = ({ raceTime, raceStatus, language, onStop }: GameControlsProps) => {
  const raceStatusLabel = {
    cs: { idle: "připraveno", running: "běží", finished: "dokončeno" },
    en: { idle: "idle", running: "running", finished: "finished" },
  };

  return (
    <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-6" role="timer" aria-label="Game clock">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-zinc-400">{language === "cs" ? "Časomíra" : "Race Clock"}</p>
          <div className="mt-1 text-5xl font-semibold tracking-wider text-zinc-100 md:text-6xl">{raceTime}</div>
          <p className="mt-2 text-xs uppercase tracking-[0.2em] text-zinc-500">
            {language === "cs" ? "Stav" : "Status"}: {raceStatusLabel[language][raceStatus]}
          </p>
        </div>

        <div className="flex gap-3">
          <button
            type="button"
            onClick={onStop}
            disabled={raceStatus !== "running"}
            className="inline-flex items-center gap-2 rounded-lg border border-[#ff0000]/80 bg-[#ff0000]/10 px-4 py-2 text-sm font-semibold text-[#ff0000] transition hover:bg-[#ff0000]/20 disabled:opacity-50"
          >
            <Square className="h-4 w-4" />
            STOP
          </button>
        </div>
      </div>
    </section>
  );
};