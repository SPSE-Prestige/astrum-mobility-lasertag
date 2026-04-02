import { useMemo } from "react";
import { Leaderboard } from "@/components/race-control/Leaderboard";
import { MatchHistoryManager } from "@/components/race-control/MatchHistoryManager";
import { t } from "@/lib/i18n";
import type { Language } from "@/types/i18n";
import type { Player, TeamResult, MatchHistoryItem } from "@/types/game";

interface ResultsPhaseProps {
  teamResults: TeamResult[];
  leaderboard: Player[];
  matchHistory: MatchHistoryItem[];
  gameMode: "team" | "ffa";
  language: Language;
  onNewGame: () => void;
}

export function ResultsPhase({ teamResults, leaderboard, matchHistory, gameMode, language, onNewGame }: ResultsPhaseProps) {
  const winner = useMemo(
    () => [...teamResults].sort((a, b) => b.score - a.score)[0],
    [teamResults],
  );

  return (
    <section className="space-y-5">
      <header>
        <h2 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">{t("phase.results", language)}</h2>
      </header>

      <article className="rounded-2xl border border-[#00ff00]/40 bg-[#00ff00]/5 p-5">
        <p className="text-xs uppercase tracking-[0.18em] text-zinc-400">
          {gameMode === "ffa" ? t("results.winningPlayer", language) : t("results.winningTeam", language)}
        </p>
        <h3 className="mt-2 text-2xl font-semibold text-[#00ff00]">
          {winner?.team ?? t("results.tbd", language)}
        </h3>
      </article>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {teamResults.map((r) => (
          <article key={r.team} className="rounded-xl border border-zinc-800 bg-black/40 p-4 text-zinc-200">
            <h4 className="text-lg font-semibold">{r.team}</h4>
            <div className="mt-3 space-y-2 text-sm">
              <p>
                {t("results.score", language)}: <span className="text-[#ff0000]">{r.score}</span>
              </p>
              <p>
                Kills: <span className="text-[#00ff00]">{r.kills}</span>
              </p>
              <p>
                Deaths: <span className="text-[#ff0000]">{r.deaths}</span>
              </p>
              <p>
                {language === "cs" ? "Přesnost" : "Accuracy"}:{" "}
                <span className="text-cyan-400">
                  {r.shotsFired > 0 ? `${((r.kills / r.shotsFired) * 100).toFixed(1)}%` : "—"}
                </span>
              </p>
            </div>
          </article>
        ))}
      </div>

      <Leaderboard players={leaderboard} gameMode={gameMode} language={language} />

      <button
        type="button"
        onClick={onNewGame}
        className="w-full rounded-xl border border-zinc-700 bg-zinc-900 px-6 py-3 text-lg font-semibold uppercase tracking-[0.16em] text-zinc-100 transition hover:border-zinc-500 hover:bg-zinc-800"
      >
        {t("results.newGame", language)}
      </button>

      <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
        <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
          {t("results.history", language)}
        </h3>
        <MatchHistoryManager items={matchHistory} language={language} />
      </article>
    </section>
  );
}
