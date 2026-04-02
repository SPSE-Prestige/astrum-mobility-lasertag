import { Crosshair, Zap } from "lucide-react";
import type { Player } from "@/types/game";
import type { GameMode } from "@/types/game";
import type { Language } from "@/types/i18n";

function formatAccuracy(kills: number, shotsFired: number): string {
  if (shotsFired === 0) return "—";
  return `${((kills / shotsFired) * 100).toFixed(1)}%`;
}

interface LeaderboardProps {
  players: Player[];
  gameMode: GameMode;
  language: Language;
}

export const Leaderboard = ({ players, gameMode, language }: LeaderboardProps) => {
  return (
    <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4">
      <p className="mb-3 text-xs uppercase tracking-[0.2em] text-zinc-500">{language === "cs" ? "Průběžné pořadí" : "Live Leaderboard"}</p>
      <div className="overflow-x-auto rounded-lg border border-zinc-800">
        <table className="w-full text-left text-sm">
          <thead className="bg-zinc-900 text-xs uppercase tracking-[0.14em] text-zinc-500">
            <tr>
              <th className="px-3 py-2">{language === "cs" ? "Hráč" : "Player"}</th>
              <th className="px-3 py-2">{language === "cs" ? "Tým" : "Team"}</th>
              <th className="px-3 py-2">Kills</th>
              <th className="px-3 py-2">Deaths</th>
              <th className="px-3 py-2">{language === "cs" ? "Skóre" : "Score"}</th>
              <th className="px-3 py-2"><Crosshair className="inline h-3.5 w-3.5 text-cyan-400" /></th>
              <th className="px-3 py-2"><Zap className="inline h-3.5 w-3.5 text-amber-400" /></th>
            </tr>
          </thead>
          <tbody>
            {players.map((player) => (
              <tr key={player.id} className="border-t border-zinc-800 text-zinc-200">
                <td className="px-3 py-2">{player.name}</td>
                <td className="px-3 py-2">{gameMode === "ffa" ? "Solo" : player.team}</td>
                <td className="px-3 py-2 text-[#ff0a0a]">{player.kills}</td>
                <td className="px-3 py-2 text-red-300">{player.deaths}</td>
                <td className="px-3 py-2 font-semibold">{player.score}</td>
                <td className="px-3 py-2 text-cyan-400">{formatAccuracy(player.kills, player.shotsFired)}</td>
                <td className="px-3 py-2 text-amber-400">{player.weaponLevel > 0 ? `LVL ${player.weaponLevel}` : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
};