import { Battery, Crosshair, Shield } from "lucide-react";
import type { Player } from "@/types/game";
import type { Language } from "@/types/i18n";

interface PlayerCardProps {
  player: Player;
  language: Language;
}

export const PlayerCard = ({ player, language }: PlayerCardProps) => {
  const hpWidth = `${player.hp}%`;
  const statusLabel =
    player.status === "active" ? (language === "cs" ? "aktivní" : "active") : language === "cs" ? "omráčen" : "stunned";

  return (
    <article className="rounded-xl border border-zinc-800 bg-zinc-950/70 p-4">
      <div className="mb-3 flex items-center justify-between">
        <div>
          <h4 className="text-sm font-semibold text-zinc-100">{player.name}</h4>
          <p className="text-xs text-zinc-500">{player.team}</p>
        </div>
        <span
          className={`rounded-full border px-2 py-1 text-xs uppercase tracking-[0.15em] ${
            player.status === "active"
              ? "border-[#00ff00]/50 bg-[#00ff00]/10 text-[#00ff00]"
              : "border-[#ff0000]/50 bg-[#ff0000]/10 text-[#ff0000]"
          }`}
        >
          {statusLabel}
        </span>
      </div>

      <div>
        <div className="mb-1 flex items-center justify-between text-xs text-zinc-400">
          <span className="inline-flex items-center gap-1">
            <Shield className="h-3.5 w-3.5" /> HP
          </span>
          <span>{player.hp}%</span>
        </div>
        <div className="h-2 overflow-hidden rounded-full bg-zinc-800">
          <div className="h-full rounded-full bg-[#00ff00] transition-all duration-500" style={{ width: hpWidth }} />
        </div>
      </div>

      <div className="mt-4 grid grid-cols-2 gap-2 text-xs text-zinc-300">
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Battery className="h-3.5 w-3.5" /> {language === "cs" ? "Munice" : "Ammo"}
          </p>
          <p className="mt-1 text-sm font-semibold">{player.ammo}</p>
        </div>
        <div className="rounded-md border border-zinc-800 bg-black/40 p-2">
          <p className="inline-flex items-center gap-1 text-zinc-500">
            <Crosshair className="h-3.5 w-3.5" /> {language === "cs" ? "Trefy" : "Hits"}
          </p>
          <p className="mt-1 text-sm font-semibold">{player.hits}</p>
        </div>
      </div>
    </article>
  );
};