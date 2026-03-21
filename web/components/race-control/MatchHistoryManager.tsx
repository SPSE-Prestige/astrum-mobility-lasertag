import type { MatchHistoryItem } from "@/types/game";
import type { Language } from "@/types/i18n";

interface MatchHistoryManagerProps {
  items: MatchHistoryItem[];
  language: Language;
}

export const MatchHistoryManager = ({ items, language }: MatchHistoryManagerProps) => {
  if (items.length === 0) {
    return (
      <p className="text-sm text-zinc-500 italic">
        {language === "cs" ? "Zatím žádné dokončené hry." : "No finished games yet."}
      </p>
    );
  }

  return (
    <div className="space-y-3">
      {items.map((item) => (
        <article key={item.id} className="rounded-lg border border-zinc-800 bg-zinc-900/60 p-3">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-sm font-semibold text-zinc-100">{item.gameName}</p>
              <p className="text-xs text-zinc-500">
                {new Date(item.playedAt).toLocaleString(language === "cs" ? "cs-CZ" : "en-US")} |{" "}
                {item.gameMode === "ffa"
                  ? language === "cs" ? "Každý proti každému" : "Free For All"
                  : language === "cs" ? "Týmová" : "Team"}
              </p>
            </div>
            <span className="rounded-full border border-zinc-700 bg-zinc-800 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider text-zinc-400">
              {item.status}
            </span>
          </div>
          <div className="mt-2 text-xs text-zinc-400">
            {language === "cs" ? "Délka" : "Duration"}: {item.durationMinutes} min
          </div>
        </article>
      ))}
    </div>
  );
};