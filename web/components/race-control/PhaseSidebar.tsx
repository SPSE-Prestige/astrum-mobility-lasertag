"use client";

import { Flag, Gauge, Settings2, Users } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import type { GamePhase } from "@/types/game";
import type { Language } from "@/types/i18n";

interface PhaseSidebarProps {
  phase: GamePhase;
  language: Language;
  raceStatus: "idle" | "running" | "paused" | "finished";
  onPhaseChange: (phase: GamePhase) => void;
}

const phases: Array<{ key: GamePhase; label: Record<Language, string>; icon: LucideIcon }> = [
  { key: "setup", label: { cs: "Nastavení", en: "Setup" }, icon: Settings2 },
  { key: "players", label: { cs: "Hráči", en: "Players" }, icon: Users },
  { key: "live", label: { cs: "Závod", en: "Race" }, icon: Gauge },
  { key: "results", label: { cs: "Výsledky", en: "Results" }, icon: Flag },
];

export const PhaseSidebar = ({ phase, language, raceStatus, onPhaseChange }: PhaseSidebarProps) => {
  const canNavigateTo = (targetPhase: GamePhase): boolean => {
    // Always allow clicking the active phase (mainly for UI consistency)
    if (targetPhase === phase) return true;

    // Rule: Players card is only accessible from Setup and Results
    // & From Players card, only Setup and Results are accessible
    if (phase === "players") {
      return targetPhase === "setup" || targetPhase === "results";
    }
    if (targetPhase === "players") {
      return phase === "setup" || phase === "results";
    }

    // If the race is running, prevent navigation elsewhere (Live is locked)
    // Note: Navigation to 'players' is already restricted above based on phase expectations
    if (raceStatus === "running") {
      return targetPhase === "live";
    }

    // Do závodu se jde jen z nastavení a jen když se spustí hra (tlačítkem Start, ne v menu)
    if (targetPhase === "live") {
      return false;
    }

    // Can navigate to anything else (Setup, Results)
    return true;
  };

  return (
    <aside className="w-full rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4 shadow-[0_0_40px_rgba(0,0,0,0.35)] backdrop-blur md:w-72">
      <div className="mb-6">
        <p className="text-xs uppercase tracking-[0.2em] text-zinc-400">{language === "cs" ? "Řízení závodu" : "Race Control"}</p>
        <h2 className="mt-2 text-2xl font-semibold text-zinc-100">{language === "cs" ? "Lasertag" : "Lasertag"}</h2>
      </div>

      <nav className="space-y-2">
        {phases.map(({ key, label, icon: Icon }) => {
          const active = phase === key;
          const isDisabled = !canNavigateTo(key);

          return (
            <button
              key={key}
              type="button"
              onClick={() => !isDisabled && onPhaseChange(key)}
              disabled={isDisabled}
              className={`flex w-full items-center justify-between rounded-xl border px-3 py-2 text-left transition ${
                isDisabled
                  ? "border-zinc-800 bg-zinc-900/30 text-zinc-500 cursor-not-allowed opacity-50"
                  : active
                    ? "border-[#00ff00]/70 bg-[#00ff00]/10 text-zinc-50"
                    : "border-zinc-800 bg-zinc-900/60 text-zinc-300 hover:border-zinc-600"
              }`}
            >
              <span className="font-medium">{label[language]}</span>
              <Icon className={`h-4 w-4 ${active ? "text-[#00ff00]" : isDisabled ? "text-zinc-700" : "text-zinc-500"}`} />
            </button>
          );
        })}
      </nav>

      <div className="mt-6 rounded-xl border border-zinc-800 bg-black/50 p-3 text-xs text-zinc-400">
        {phase === "live" && raceStatus === "running"
          ? language === "cs"
            ? "Závod probíhá. Nelze přejít na jinou sekci."
            : "Race in progress. Cannot switch sections."
          : language === "cs"
            ? "Připraveno na mock real-time stream."
            : "Prepared for mock real-time stream."}
      </div>
    </aside>
  );
};