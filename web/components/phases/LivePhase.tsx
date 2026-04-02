import { GameControls } from "@/components/race-control/GameControls";
import { Leaderboard } from "@/components/race-control/Leaderboard";
import { PlayerCard } from "@/components/race-control/PlayerCard";
import { t } from "@/lib/i18n";
import type { Language } from "@/types/i18n";
import type { Player, KillFeedEvent } from "@/types/game";

interface LivePhaseProps {
  players: Player[];
  leaderboard: Player[];
  killFeed: KillFeedEvent[];
  raceTime: string;
  raceStatus: "idle" | "running" | "finished";
  gameMode: "team" | "ffa";
  language: Language;
  onStop: () => void;
}

export function LivePhase({ players, leaderboard, killFeed, raceTime, raceStatus, gameMode, language, onStop }: LivePhaseProps) {
  return (
    <section className="space-y-4">
      <GameControls
        raceTime={raceTime}
        raceStatus={raceStatus}
        language={language}
        onStop={onStop}
      />

      <div className="grid gap-4 xl:grid-cols-[2fr_1fr]">
        <div className="grid gap-3 md:grid-cols-2">
          {players.map((player) => (
            <PlayerCard key={player.id} player={player} language={language} />
          ))}
        </div>

        <div className="space-y-4">
          <Leaderboard players={leaderboard} gameMode={gameMode} language={language} />
          <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4">
            <p className="mb-3 text-xs uppercase tracking-[0.2em] text-zinc-500">
              {t("live.killFeed", language)}
            </p>
            <div className="space-y-2">
              {killFeed.map((entry) => (
                <div key={entry.id} className="rounded-md border border-zinc-800 bg-black/40 px-3 py-2 text-xs text-zinc-300">
                  <span className="mr-2 text-[#ff0a0a]">[{entry.timestamp}]</span>
                  {entry.message}
                </div>
              ))}
              {killFeed.length === 0 && (
                <p className="text-xs text-zinc-600 italic">
                  {t("live.noEvents", language)}
                </p>
              )}
            </div>
          </section>
        </div>
      </div>
    </section>
  );
}
