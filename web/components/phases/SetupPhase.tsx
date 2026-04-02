import { Settings, SlidersHorizontal } from "lucide-react";
import { t } from "@/lib/i18n";
import type { Language } from "@/types/i18n";
import type { GameConfig } from "@/types/game";

interface SetupPhaseProps {
  config: GameConfig;
  language: Language;
  loading: boolean;
  onCreateGame: () => void;
  onUpdateConfig: (config: GameConfig) => void;
}

export function SetupPhase({ config, language, loading, onCreateGame, onUpdateConfig }: SetupPhaseProps) {
  return (
    <section className="space-y-6">
      <header>
        <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">{t("phase.1", language)}</p>
        <h1 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">{t("phase.setup", language)}</h1>
      </header>

      <button
        type="button"
        onClick={onCreateGame}
        disabled={loading}
        className="w-full rounded-xl border-2 border-[#ff0a0a] bg-[#ff0a0a]/20 px-6 py-3 text-lg font-semibold uppercase tracking-[0.16em] text-[#ff0a0a] shadow-[0_0_20px_rgba(255,10,10,0.3)] transition hover:bg-[#ff0a0a]/30 hover:shadow-[0_0_30px_rgba(255,10,10,0.4)] disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? t("setup.creating", language) : t("setup.createGame", language)}
      </button>

      <div className="grid gap-4 xl:grid-cols-2 xl:items-start">
        {/* Game Parameters */}
        <article className="self-start rounded-2xl border border-zinc-800 bg-black/40 px-3 py-2 pb-1">
          <h3 className="mb-2 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <Settings className="h-4 w-4 text-[#ff0a0a]" />
            {t("setup.params", language)}
          </h3>
          <div className="grid gap-0.5">
            <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
              {t("setup.name", language)}
              <input
                value={config.gameName}
                onChange={(e) => onUpdateConfig({ ...config, gameName: e.target.value })}
                className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-sm text-zinc-100 outline-none ring-[#ff0a0a] transition focus:ring-1"
              />
            </label>
            <div className="grid grid-cols-2 gap-2">
              <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                {t("setup.duration", language)}
                <input
                  type="number"
                  min={1}
                  value={config.durationMinutes}
                  onChange={(e) => onUpdateConfig({ ...config, durationMinutes: Number(e.target.value) })}
                  className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-sm text-zinc-100 outline-none ring-[#ff0a0a] transition focus:ring-1"
                />
              </label>
              <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                {t("setup.mode", language)}
                <select
                  value={config.gameMode}
                  onChange={(e) => onUpdateConfig({ ...config, gameMode: e.target.value as "team" | "ffa" })}
                  className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-sm text-zinc-100 outline-none ring-[#ff0a0a] transition focus:ring-1"
                >
                  <option value="team">{t("setup.modeTeam", language)}</option>
                  <option value="ffa">{t("setup.modeFFA", language)}</option>
                </select>
              </label>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                {t("setup.teams", language)}
                <select
                  value={config.teamsCount}
                  onChange={(e) => onUpdateConfig({ ...config, teamsCount: Number(e.target.value) })}
                  disabled={config.gameMode === "ffa"}
                  className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-sm text-zinc-100 outline-none ring-[#ff0a0a] transition focus:ring-1"
                >
                  {[2, 3, 4].map((n) => (
                    <option key={n} value={n}>{n}</option>
                  ))}
                </select>
              </label>
              <div className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                <div className="mt-5 inline-flex w-full items-center justify-between rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-sm text-zinc-200">
                  <span>Friendly Fire</span>
                  <button
                    type="button"
                    onClick={() => onUpdateConfig({ ...config, friendlyFire: !config.friendlyFire })}
                    className={`mt-1 h-3.5 w-12 rounded-full px-1 transition ${config.friendlyFire ? "bg-[#ff0000]/60" : "bg-[#ff0a0a]/30"}`}
                  >
                    <span className={`block h-2 w-4 rounded-full bg-white transition ${config.friendlyFire ? "translate-x-6" : "translate-x-0"}`} />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </article>

        {/* Additional Settings */}
        <article className="self-start rounded-2xl border border-zinc-800 bg-black/40 p-4">
          <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <SlidersHorizontal className="h-4 w-4 text-[#ff0000]" />
            {t("setup.additional", language)}
          </h3>
          <div className="grid gap-3">
            <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
              {t("setup.respawnDelay", language)}
              <input
                type="number"
                min={0}
                value={config.respawnDelay}
                onChange={(e) => onUpdateConfig({ ...config, respawnDelay: Number(e.target.value) })}
                className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
              />
            </label>
            <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
              {t("setup.maxPlayers", language)}
              <input
                type="number"
                min={2}
                value={config.maxPlayers}
                onChange={(e) => onUpdateConfig({ ...config, maxPlayers: Number(e.target.value) })}
                className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
              />
            </label>
            <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
              {t("setup.killsPerUpgrade", language)}
              <input
                type="number"
                min={0}
                max={50}
                value={config.killsPerUpgrade}
                onChange={(e) => onUpdateConfig({ ...config, killsPerUpgrade: Number(e.target.value) })}
                className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
              />
              <span className="mt-1 block text-[10px] text-zinc-600">
                {t("setup.killsPerUpgradeHint", language)}
              </span>
            </label>
          </div>
        </article>
      </div>
    </section>
  );
}
