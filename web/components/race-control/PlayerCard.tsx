import { Crosshair, Skull, Swords, Target, Zap } from "lucide-react";
import type { Player } from "@/types/game";
import type { Language } from "@/types/i18n";

interface PlayerCardProps {
  player: Player;
  language: Language;
}

export const PlayerCard = ({ player, language }: PlayerCardProps) => {
  const isAlive = player.status === "alive";
  const statusLabel = isAlive
    ? language === "cs" ? "naživu" : "alive"
    : language === "cs" ? "vyřazen" : "eliminated";

  return (
    <article
      role="region"
      aria-label={`${player.name} – ${statusLabel}`}
      className={`rounded-xl border p-4 transition-all ${
        isAlive
          ? "border-zinc-800 bg-zinc-950/70"
          : "border-red-900/50 bg-red-950/20 opacity-75"
      }`}
    >
      <div className="mb-3 flex items-center justify-between">
        <div>
          <h4 className="text-sm font-semibold text-zinc-100">{player.name}</h4>
          <p className="text-xs text-zinc-500">{player.team}</p>
        </div>
        <div className="flex items-center gap-2">
          {player.weaponLevel > 0 && (
            <span className="inline-flex items-center gap-1 rounded-full border border-amber-500/50 bg-amber-500/10 px-2 py-0.5 text-xs font-semibold text-amber-400">
              <Zap className="h-3 w-3" />
              LVL {player.weaponLevel}
            </span>
          )}
          <span
            className={`rounded-full border px-2 py-1 text-xs uppercase tracking-[0.15em] ${
              isAlive
                ? "border-[#ff0a0a]/50 bg-[#ff0a0a]/10 text-[#ff0a0a]"
                : "border-red-300/50 bg-red-300/10 text-red-300"
            }`}
          >
            {statusLabel}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-4 gap-2 text-xs text-zinc-300">
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Crosshair className="h-3.5 w-3.5" /> Kills
          </p>
          <p className="mt-1 text-sm font-semibold text-[#ff0a0a]">{player.kills}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Skull className="h-3.5 w-3.5" /> Deaths
          </p>
          <p className="mt-1 text-sm font-semibold text-red-300">{player.deaths}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Target className="h-3.5 w-3.5" /> {language === "cs" ? "Skóre" : "Score"}
          </p>
          <p className="mt-1 text-sm font-semibold">{player.score}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Swords className="h-3.5 w-3.5" /> Streak
          </p>
          <p className="mt-1 text-sm font-semibold text-amber-400">{player.killStreak}</p>
        </div>
      </div>
    </article>
  );
};