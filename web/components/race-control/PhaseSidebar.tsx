"use client";

import { Flag, Gauge, Settings2 } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import type { GamePhase } from "@/types/game";

interface PhaseSidebarProps {
  phase: GamePhase;
  onPhaseChange: (phase: GamePhase) => void;
}

const phases: Array<{ key: GamePhase; label: string; icon: LucideIcon }> = [
  { key: "setup", label: "Setup", icon: Settings2 },
  { key: "live", label: "Live Race", icon: Gauge },
  { key: "results", label: "Results", icon: Flag },
];

export const PhaseSidebar = ({ phase, onPhaseChange }: PhaseSidebarProps) => {
  return (
    <aside className="w-full rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4 shadow-[0_0_40px_rgba(0,0,0,0.35)] backdrop-blur md:w-72">
      <div className="mb-6">
        <p className="text-xs uppercase tracking-[0.2em] text-zinc-400">Race Control</p>
        <h2 className="mt-2 text-2xl font-semibold text-zinc-100">Telemetry Hub</h2>
      </div>

      <nav className="space-y-2">
        {phases.map(({ key, label, icon: Icon }) => {
          const active = phase === key;

          return (
            <button
              key={key}
              type="button"
              onClick={() => onPhaseChange(key)}
              className={`flex w-full items-center justify-between rounded-xl border px-3 py-2 text-left transition ${
                active
                  ? "border-[#00ff00]/70 bg-[#00ff00]/10 text-zinc-50"
                  : "border-zinc-800 bg-zinc-900/60 text-zinc-300 hover:border-zinc-600"
              }`}
            >
              <span className="font-medium">{label}</span>
              <Icon className={`h-4 w-4 ${active ? "text-[#00ff00]" : "text-zinc-500"}`} />
            </button>
          );
        })}
      </nav>

      <div className="mt-6 rounded-xl border border-zinc-800 bg-black/50 p-3 text-xs text-zinc-400">
        Prepared for mock real-time stream.
      </div>
    </aside>
  );
};