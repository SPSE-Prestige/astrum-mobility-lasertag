import type { Player } from "@/types/game";
import type { GameMode } from "@/types/game";

interface LeaderboardProps {
  players: Player[];
  gameMode: GameMode;
}

export const Leaderboard = ({ players, gameMode }: LeaderboardProps) => {
  return (
    <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4">
      <p className="mb-3 text-xs uppercase tracking-[0.2em] text-zinc-500">Live Leaderboard</p>
      <div className="overflow-hidden rounded-lg border border-zinc-800">
        <table className="w-full text-left text-sm">
          <thead className="bg-zinc-900 text-xs uppercase tracking-[0.14em] text-zinc-500">
            <tr>
              <th className="px-3 py-2">Player</th>
              <th className="px-3 py-2">{gameMode === "ffa" ? "Mode" : "Team"}</th>
              <th className="px-3 py-2">Kills</th>
              <th className="px-3 py-2">Hits</th>
            </tr>
          </thead>
          <tbody>
            {players.map((player) => (
              <tr key={player.id} className="border-t border-zinc-800 text-zinc-200">
                <td className="px-3 py-2">{player.name}</td>
                <td className="px-3 py-2">{gameMode === "ffa" ? "Solo" : player.team}</td>
                <td className="px-3 py-2 text-[#ff0000]">{player.kills}</td>
                <td className="px-3 py-2 text-[#00ff00]">{player.hits}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
};