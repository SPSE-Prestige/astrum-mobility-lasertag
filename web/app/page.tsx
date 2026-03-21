"use client";

import { useEffect, useMemo, useState } from "react";
import { LogOut, Settings, SlidersHorizontal, Users } from "lucide-react";
import { GameControls } from "@/components/race-control/GameControls";
import { Leaderboard } from "@/components/race-control/Leaderboard";
import { LoginPanel } from "@/components/race-control/LoginPanel";
import { MatchHistoryManager } from "@/components/race-control/MatchHistoryManager";
import { PhaseSidebar } from "@/components/race-control/PhaseSidebar";
import { PlayerCard } from "@/components/race-control/PlayerCard";
import { PlayerRegistration } from "@/components/race-control/PlayerRegistration";
import { useGameData } from "@/hooks/useGameData";
import type { Language } from "@/types/i18n";

const LANGUAGE_STORAGE_KEY = "race-control-language";

const formatRaceClock = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
    .toString()
    .padStart(2, "0");
  const secs = Math.max(seconds % 60, 0)
    .toString()
    .padStart(2, "0");
  return `${mins}:${secs}`;
};

const translateKillFeedMessage = (message: string, language: Language) => {
  if (message === "event.systemReady") {
    return language === "cs" ? "Systém je připraven. Čeká na start." : "System is ready. Waiting for start signal.";
  }

  if (message.startsWith("event.tag:")) {
    const [, killer, victim] = message.split(":");
    return language === "cs" ? `${killer} zasáhl ${victim}` : `${killer} tagged ${victim}`;
  }

  return message;
};

export default function Home() {
  const {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
    registeredPlayers,
    activeRoster,
    registerPlayer,
    deleteRegisteredPlayer,
    toggleRosterPlayer,
    updateConfig,
    updatePhase,
    assignPlayerTeam,
    startRace,
    pauseRace,
    stopRace,
    login,
    logout,
    updateMatchHistoryItem,
    deleteMatchHistoryItem,
  } = useGameData();
  const [language, setLanguage] = useState<Language>(() => {
    if (typeof window === "undefined") {
      return "cs";
    }
    const storedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY);
    return storedLanguage === "en" ? "en" : "cs";
  });

  useEffect(() => {
    document.documentElement.lang = language;
    localStorage.setItem(LANGUAGE_STORAGE_KEY, language);
  }, [language]);

  const winner = useMemo(() => {
    return [...state.teamResults].sort((a, b) => b.score - a.score)[0];
  }, [state.teamResults]);

  if (!auth.isAuthenticated) {
    return (
      <div className="min-h-screen bg-[radial-gradient(circle_at_20%_20%,rgba(0,255,0,0.14),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(255,0,0,0.12),transparent_30%),#020303] p-4 md:p-8">
        <div className="mx-auto flex min-h-[80vh] w-full max-w-xl items-center justify-center">
          <div className="w-full rounded-2xl border border-zinc-800 bg-zinc-950/70 p-5 shadow-[0_0_40px_rgba(0,0,0,0.35)] backdrop-blur">
            <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">{language === "cs" ? "Přístup do řízení závodu" : "Race Control Access"}</p>
            <h1 className="mt-2 text-3xl font-semibold text-zinc-100">
              {language === "cs" ? "Přihlášení do dashboardu" : "Dashboard Login"}
            </h1>
            <p className="mt-2 text-sm text-zinc-400">
              {language === "cs"
                ? "Bez přihlášení není přístup k nastavení ani živé obrazovce."
                : "Without login, settings and live screens are not accessible."}
            </p>
            <div className="mt-4 flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setLanguage("cs")}
                className={`rounded-md border px-2 py-1 text-xs ${
                  language === "cs" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
                }`}
              >
                CZ
              </button>
              <button
                type="button"
                onClick={() => setLanguage("en")}
                className={`rounded-md border px-2 py-1 text-xs ${
                  language === "en" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
                }`}
              >
                EN
              </button>
            </div>
            <div className="mt-4">
              <LoginPanel
                isAuthenticated={auth.isAuthenticated}
                username={auth.username}
                error={auth.error}
                language={language}
                onLogin={login}
                onLogout={logout}
              />
            </div>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_20%_20%,rgba(0,255,0,0.14),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(255,0,0,0.12),transparent_30%),#020303] p-4 md:p-8">
      <div className="mx-auto mb-4 flex w-full max-w-350 items-center justify-between gap-3 rounded-xl border border-zinc-800 bg-black/30 px-3 py-2 md:hidden">
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setLanguage("cs")}
            className={`rounded-md border px-2 py-1 text-xs ${
              language === "cs" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
            }`}
          >
            CZ
          </button>
          <button
            type="button"
            onClick={() => setLanguage("en")}
            className={`rounded-md border px-2 py-1 text-xs ${
              language === "en" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
            }`}
          >
            EN
          </button>
        </div>
        <span className="text-xs uppercase tracking-[0.16em] text-zinc-400">{auth.username}</span>
        <button
          type="button"
          onClick={logout}
          className="inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-zinc-200 transition hover:border-zinc-500"
        >
          <LogOut className="h-3.5 w-3.5" />
          {language === "cs" ? "Odhlásit" : "Logout"}
        </button>
      </div>

      <div className="mx-auto grid w-full max-w-350 gap-4 md:grid-cols-[288px_1fr]">
        <PhaseSidebar phase={state.phase} language={language} raceStatus={state.raceStatus} onPhaseChange={updatePhase} />

        <main className="rounded-2xl border border-zinc-800 bg-zinc-950/60 p-4 md:p-6">
          <div className="mb-4 hidden items-center justify-between gap-3 rounded-xl border border-zinc-800 bg-black/30 px-3 py-2 md:flex">
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => setLanguage("cs")}
                className={`rounded-md border px-2 py-1 text-xs ${
                  language === "cs" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
                }`}
              >
                CZ
              </button>
              <button
                type="button"
                onClick={() => setLanguage("en")}
                className={`rounded-md border px-2 py-1 text-xs ${
                  language === "en" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"
                }`}
              >
                EN
              </button>
            </div>
            <span className="text-xs uppercase tracking-[0.16em] text-zinc-400">{auth.username}</span>
            <button
              type="button"
              onClick={logout}
              className="inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-zinc-200 transition hover:border-zinc-500"
            >
              <LogOut className="h-3.5 w-3.5" />
              {language === "cs" ? "Odhlásit" : "Logout"}
            </button>
          </div>

          {state.phase === "setup" && (

            <section className="space-y-6">
              <header>
                <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">{language === "cs" ? "Fáze 1" : "Phase 1"}</p>
                <h1 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">{language === "cs" ? "Nastavení hry" : "Setup Control"}</h1>
              </header>

              <button
                type="button"
                onClick={startRace}
                className="w-full rounded-xl border-2 border-[#00ff00] bg-[#00ff00]/20 px-6 py-3 text-lg font-semibold uppercase tracking-[0.16em] text-[#00ff00] shadow-[0_0_20px_rgba(0,255,0,0.3)] transition hover:bg-[#00ff00]/30 hover:shadow-[0_0_30px_rgba(0,255,0,0.4)]"
              >
                {language === "cs" ? "Začít hru" : "Start Game"}
              </button>

              <div className="grid gap-4 xl:grid-cols-2">
                <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                  <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                    <Settings className="h-4 w-4 text-[#00ff00]" />
                    {language === "cs" ? "Parametry hry" : "Game Parameters"}
                  </h3>
                  <div className="grid gap-3">
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      {language === "cs" ? "Název" : "Name"}
                      <input
                        value={config.gameName}
                        onChange={(event) => updateConfig({ ...config, gameName: event.target.value })}
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                      />
                    </label>
                    <div className="grid grid-cols-2 gap-3">
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        {language === "cs" ? "Délka (min)" : "Duration (min)"}
                        <input
                          type="number"
                          min={5}
                          value={config.durationMinutes}
                          onChange={(event) => updateConfig({ ...config, durationMinutes: Number(event.target.value) })}
                          className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                        />
                      </label>
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        {language === "cs" ? "Herní režim" : "Game Mode"}
                        <select
                          value={config.gameMode}
                          onChange={(event) =>
                            updateConfig({
                              ...config,
                              gameMode: event.target.value as "team" | "ffa",
                            })
                          }
                          className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                        >
                          <option value="team">{language === "cs" ? "Týmová hra" : "Team Game"}</option>
                          <option value="ffa">{language === "cs" ? "Každý proti každému" : "Free For All"}</option>
                        </select>
                      </label>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        {language === "cs" ? "Týmy" : "Teams"}
                        <select
                          value={config.teamsCount}
                          onChange={(event) => updateConfig({ ...config, teamsCount: Number(event.target.value) })}
                          disabled={config.gameMode === "ffa"}
                          className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                        >
                          {[2, 3, 4].map((teamCount) => (
                            <option key={teamCount} value={teamCount}>
                              {teamCount}
                            </option>
                          ))}
                        </select>
                      </label>
                      <div className="rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-xs text-zinc-400">
                        {config.gameMode === "ffa"
                          ? language === "cs"
                            ? "V režimu FFA se týmy nepoužívají."
                            : "Teams are not used in FFA mode."
                          : language === "cs"
                            ? "V týmovém režimu lze přiřadit vozíky do týmu."
                            : "In team mode, karts can be assigned to teams."}
                      </div>
                    </div>
                    <label className="mt-1 inline-flex items-center justify-between rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-200">
                      {language === "cs" ? "Friendly Fire" : "Friendly Fire"}
                      <button
                        type="button"
                        onClick={() => updateConfig({ ...config, friendlyFire: !config.friendlyFire })}
                        className={`h-6 w-12 rounded-full p-1 transition ${
                          config.friendlyFire ? "bg-[#ff0000]/60" : "bg-[#00ff00]/30"
                        }`}
                      >
                        <span
                          className={`block h-4 w-4 rounded-full bg-white transition ${config.friendlyFire ? "translate-x-6" : "translate-x-0"}`}
                        />
                      </button>
                    </label>
                  </div>
                </article>

                <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                  <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                    <SlidersHorizontal className="h-4 w-4 text-[#ff0000]" />
                    {language === "cs" ? "Ladění hráčů a vozíku" : "Player & Cart Tuning"}
                  </h3>
                  <div className="grid grid-cols-2 gap-3">
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      {language === "cs" ? "Prodleva respawnu" : "Respawn Delay"}
                      <input
                        type="number"
                        value={config.playerTuning.respawnDelay}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            playerTuning: { ...config.playerTuning, respawnDelay: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      {language === "cs" ? "Rychlost vozíku" : "Cart Speed"}
                      <input
                        type="number"
                        value={config.playerTuning.cartSpeed}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            playerTuning: { ...config.playerTuning, cartSpeed: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                  </div>
                </article>

                <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                  <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                    <Users className="h-4 w-4 text-[#00ff00]" />
                    {language === "cs" ? "Přiřazení hráčů do týmů" : "Assign Players to Teams"}
                  </h3>
                  
                  {config.gameMode === "ffa" ? (
                    <div className="rounded-lg border border-zinc-700 bg-zinc-900/50 px-4 py-3 text-sm text-zinc-400">
                      {language === "cs"
                        ? "V režimu 'Každý proti každému' se týmy nepoužívají."
                        : "In 'Free For All' mode, teams are not used."}
                    </div>
                  ) : (
                    <div className="space-y-3">
                      {Array.from({ length: config.teamsCount }).map((_, teamIdx) => {
                        const teams = ["Neon Green", "Neon Red", "Neon Blue", "Neon Yellow"];
                        const teamName = teams[teamIdx] || `Team ${teamIdx + 1}`;
                        const playersInTeam = state.players.filter((p) => p.team === teamName);

                        return (
                          <div key={teamIdx} className="rounded-lg border border-zinc-700 bg-zinc-900/50 p-3">
                            <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-zinc-300">{teamName}</p>
                            <div className="space-y-2">
                              {state.players.map((player) => {
                                const isInThisTeam = player.team === teamName;
                                return (
                                  <button
                                    key={player.id}
                                    onClick={() => assignPlayerTeam(player.id, isInThisTeam ? "Unassigned" : teamName)}
                                    className={`w-full rounded-md border px-3 py-2 text-left text-sm transition ${
                                      isInThisTeam
                                        ? "border-[#00ff00]/50 bg-[#00ff00]/10 text-[#00ff00]"
                                        : "border-zinc-700 bg-zinc-800/30 text-zinc-400 hover:border-zinc-600 hover:bg-zinc-800/50"
                                    }`}
                                  >
                                    {player.name}
                                  </button>
                                );
                              })}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  )}
                </article>
              </div>
            </section>
          )}

          {state.phase === "players" && (
            <PlayerRegistration
              activePlayers={state.players}
              registeredPlayers={registeredPlayers}
              activeRoster={activeRoster}
              gameMode={config.gameMode}
              language={language}
              onRegister={registerPlayer}
              onDeleteRegistry={deleteRegisteredPlayer}
              onToggleRoster={toggleRosterPlayer}
            />
          )}

          {state.phase === "live" && (
            <section className="space-y-4">
              <GameControls
                raceTime={formatRaceClock(state.raceTimeSeconds)}
                raceStatus={state.raceStatus}
                language={language}
                onPause={pauseRace}
                onStop={stopRace}
              />

              <div className="grid gap-4 xl:grid-cols-[2fr_1fr]">
                <div className="grid gap-3 md:grid-cols-2">
                  {state.players.map((player) => (
                    <PlayerCard key={player.id} player={player} language={language} />
                  ))}
                </div>

                <div className="space-y-4">
                  <Leaderboard players={leaderboard} gameMode={config.gameMode} language={language} />
                  <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4">
                    <p className="mb-3 text-xs uppercase tracking-[0.2em] text-zinc-500">{language === "cs" ? "Feed událostí" : "Kill Feed"}</p>
                    <div className="space-y-2">
                      {state.killFeed.map((entry) => (
                        <div key={entry.id} className="rounded-md border border-zinc-800 bg-black/40 px-3 py-2 text-xs text-zinc-300">
                          <span className="mr-2 text-[#00ff00]">[{entry.timestamp}]</span>
                          {translateKillFeedMessage(entry.message, language)}
                        </div>
                      ))}
                    </div>
                  </section>
                </div>
              </div>
            </section>
          )}

          {state.phase === "results" && (
            <section className="space-y-5">
              <header>
                <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">{language === "cs" ? "Fáze 3" : "Phase 3"}</p>
                <h2 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">{language === "cs" ? "Finální výsledky" : "Final Results"}</h2>
              </header>

              <article className="rounded-2xl border border-[#00ff00]/40 bg-[#00ff00]/5 p-5">
                <p className="text-xs uppercase tracking-[0.18em] text-zinc-400">
                  {config.gameMode === "ffa"
                    ? language === "cs"
                      ? "Vítězný hráč"
                      : "Winning Player"
                    : language === "cs"
                      ? "Vítězný tým"
                      : "Winning Team"}
                </p>
                <h3 className="mt-2 text-2xl font-semibold text-[#00ff00]">{winner?.team ?? (language === "cs" ? "Neurčeno" : "TBD")}</h3>
              </article>

              <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {state.teamResults.map((teamResult) => (
                  <article key={teamResult.team} className="rounded-xl border border-zinc-800 bg-black/40 p-4 text-zinc-200">
                    <h4 className="text-lg font-semibold">{teamResult.team}</h4>
                    <div className="mt-3 space-y-2 text-sm">
                      <p>
                        {language === "cs" ? "Skóre" : "Score"}: <span className="text-[#ff0000]">{teamResult.score}</span>
                      </p>
                      <p>
                        {language === "cs" ? "Přesnost" : "Accuracy"}: <span className="text-[#00ff00]">{teamResult.accuracy}%</span>
                      </p>
                      <p>
                        {language === "cs" ? "Udělené poškození" : "Damage Dealt"}: <span className="text-[#ff0000]">{teamResult.damageDealt}</span>
                      </p>
                    </div>
                  </article>
                ))}
              </div>

              <button
                type="button"
                onClick={() => updatePhase("setup")}
                className="w-full rounded-xl border border-zinc-700 bg-zinc-900 px-6 py-3 text-lg font-semibold uppercase tracking-[0.16em] text-zinc-100 transition hover:border-zinc-500 hover:bg-zinc-800"
              >
                {language === "cs" ? "Zpátky na nastavení" : "Back to Setup"}
              </button>

              <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                  {language === "cs" ? "Historie her" : "Match History"}
                </h3>
                <MatchHistoryManager
                  items={matchHistory}
                  language={language}
                  onUpdate={updateMatchHistoryItem}
                  onDelete={deleteMatchHistoryItem}
                />
              </article>
            </section>
          )}
        </main>
      </div>
    </div>
  );
}
