"use client";

import { useState } from "react";
import { Pencil, Save, Trash2, X } from "lucide-react";
import type { MatchHistoryItem } from "@/types/game";
import type { Language } from "@/types/i18n";

interface MatchHistoryManagerProps {
  items: MatchHistoryItem[];
  language: Language;
  onUpdate: (
    matchId: string,
    patch: Partial<Pick<MatchHistoryItem, "gameName" | "winner" | "durationMinutes" | "totalKills">>,
  ) => void;
  onDelete: (matchId: string) => void;
}

export const MatchHistoryManager = ({ items, language, onUpdate, onDelete }: MatchHistoryManagerProps) => {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [draft, setDraft] = useState({ gameName: "", winner: "", durationMinutes: 0, totalKills: 0 });

  const beginEdit = (item: MatchHistoryItem) => {
    setEditingId(item.id);
    setDraft({
      gameName: item.gameName,
      winner: item.winner,
      durationMinutes: item.durationMinutes,
      totalKills: item.totalKills,
    });
  };

  return (
    <section className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
      <p className="mb-3 text-xs uppercase tracking-[0.18em] text-zinc-500">{language === "cs" ? "Historie zápasů" : "Match History"}</p>
      <div className="space-y-3">
        {items.map((item) => {
          const isEditing = editingId === item.id;

          return (
            <article key={item.id} className="rounded-lg border border-zinc-800 bg-zinc-900/60 p-3">
              {!isEditing && (
                <>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-zinc-100">{item.gameName}</p>
                      <p className="text-xs text-zinc-500">
                        {new Date(item.playedAt).toLocaleString(language === "cs" ? "cs-CZ" : "en-US")} | {item.gameMode === "ffa" ? (language === "cs" ? "Každý proti každému" : "Free For All") : language === "cs" ? "Týmová" : "Team"}
                      </p>
                    </div>
                    <div className="flex gap-2">
                      <button
                        type="button"
                        onClick={() => beginEdit(item)}
                        className="rounded-md border border-zinc-700 bg-zinc-900 p-2 text-zinc-300 transition hover:border-zinc-500"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => onDelete(item.id)}
                        className="rounded-md border border-[#ff0000]/60 bg-[#ff0000]/10 p-2 text-[#ff0000] transition hover:bg-[#ff0000]/20"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                  <div className="mt-2 grid grid-cols-2 gap-2 text-xs text-zinc-300 md:grid-cols-4">
                    <p>{language === "cs" ? "Vítěz" : "Winner"}: {item.winner}</p>
                    <p>{language === "cs" ? "Délka" : "Duration"}: {item.durationMinutes} min</p>
                    <p>Kills: {item.totalKills}</p>
                  </div>
                </>
              )}

              {isEditing && (
                <div className="space-y-2">
                  <input
                    value={draft.gameName}
                    onChange={(event) => setDraft((prev) => ({ ...prev, gameName: event.target.value }))}
                    className="w-full rounded-md border border-zinc-700 bg-black px-3 py-2 text-sm text-zinc-100"
                  />
                  <div className="grid gap-2 md:grid-cols-3">
                    <input
                      value={draft.winner}
                      onChange={(event) => setDraft((prev) => ({ ...prev, winner: event.target.value }))}
                      placeholder={language === "cs" ? "Vítěz" : "Winner"}
                      className="rounded-md border border-zinc-700 bg-black px-3 py-2 text-sm text-zinc-100"
                    />
                    <input
                      type="number"
                      value={draft.durationMinutes}
                      onChange={(event) => setDraft((prev) => ({ ...prev, durationMinutes: Number(event.target.value) }))}
                      placeholder={language === "cs" ? "Délka" : "Duration"}
                      className="rounded-md border border-zinc-700 bg-black px-3 py-2 text-sm text-zinc-100"
                    />
                    <input
                      type="number"
                      value={draft.totalKills}
                      onChange={(event) => setDraft((prev) => ({ ...prev, totalKills: Number(event.target.value) }))}
                      placeholder="Kills"
                      className="rounded-md border border-zinc-700 bg-black px-3 py-2 text-sm text-zinc-100"
                    />
                  </div>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      onClick={() => {
                        onUpdate(item.id, draft);
                        setEditingId(null);
                      }}
                      className="inline-flex items-center gap-1 rounded-md border border-[#00ff00]/70 bg-[#00ff00]/10 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-[#00ff00]"
                    >
                      <Save className="h-3.5 w-3.5" />
                      {language === "cs" ? "Uložit" : "Save"}
                    </button>
                    <button
                      type="button"
                      onClick={() => setEditingId(null)}
                      className="inline-flex items-center gap-1 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-zinc-300"
                    >
                      <X className="h-3.5 w-3.5" />
                      {language === "cs" ? "Zrušit" : "Cancel"}
                    </button>
                  </div>
                </div>
              )}
            </article>
          );
        })}
      </div>
    </section>
  );
};