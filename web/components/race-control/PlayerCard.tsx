import { Crosshair, Skull, Target } from "lucide-react";
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
      className={`rounded-xl border p-4 ${
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
        <span
          className={`rounded-full border px-2 py-1 text-xs uppercase tracking-[0.15em] ${
            isAlive
              ? "border-[#00ff00]/50 bg-[#00ff00]/10 text-[#00ff00]"
              : "border-[#ff0000]/50 bg-[#ff0000]/10 text-[#ff0000]"
          }`}
        >
          {statusLabel}
        </span>
      </div>

      <div className="grid grid-cols-3 gap-2 text-xs text-zinc-300">
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Crosshair className="h-3.5 w-3.5" /> Kills
          </p>
          <p className="mt-1 text-sm font-semibold text-[#00ff00]">{player.kills}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Skull className="h-3.5 w-3.5" /> Deaths
          </p>
          <p className="mt-1 text-sm font-semibold text-[#ff0000]">{player.deaths}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Target className="h-3.5 w-3.5" /> {language === "cs" ? "Skóre" : "Score"}
          </p>
          <p className="mt-1 text-sm font-semibold">{player.score}</p>
        </div>
      </div>
    </article>
  );
};