"use client";

import { useEffect, useState } from "react";
import { LoginPanel } from "@/components/race-control/LoginPanel";
import { PhaseSidebar } from "@/components/race-control/PhaseSidebar";
import { PlayerRegistration } from "@/components/race-control/PlayerRegistration";
import { TopBar } from "@/components/layout/TopBar";
import { SetupPhase } from "@/components/phases/SetupPhase";
import { LivePhase } from "@/components/phases/LivePhase";
import { ResultsPhase } from "@/components/phases/ResultsPhase";
import { ToastContainer, showToast } from "@/components/ui/Toast";
import { ErrorBoundary } from "@/components/ui/ErrorBoundary";
import { useGameData } from "@/hooks/useGameData";
import { t } from "@/lib/i18n";
import type { Language } from "@/types/i18n";

const LANGUAGE_STORAGE_KEY = "race-control-language";
const BG_GRADIENT = "min-h-screen bg-[radial-gradient(circle_at_20%_20%,rgba(0,255,0,0.14),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(255,0,0,0.12),transparent_30%),#020303] p-4 md:p-8";

export default function Home() {
  const {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
    wsStatus,
    loading,
    error,
    clearError,
    formatRaceTime,
    updateConfig,
    updatePhase,
    createGame,
    addPlayer,
    removePlayer,
    assignPlayerTeam,
    refreshDevices,
    startRace,
    stopRace,
    login,
    logout,
  } = useGameData();

  const [language, setLanguage] = useState<Language>(() => {
    if (typeof window === "undefined") return "cs";
    const stored = localStorage.getItem(LANGUAGE_STORAGE_KEY);
    return stored === "en" ? "en" : "cs";
  });

  useEffect(() => {
    document.documentElement.lang = language;
    localStorage.setItem(LANGUAGE_STORAGE_KEY, language);
  }, [language]);

  // Show toast on errors
  useEffect(() => {
    if (error) {
      showToast("error", error);
      clearError();
    }
  }, [error, clearError]);

  // ── Login screen ──
  if (!auth.isAuthenticated) {
    return (
      <div className={BG_GRADIENT}>
        <div className="mx-auto flex min-h-[80vh] w-full max-w-xl items-center justify-center">
          <div className="w-full rounded-2xl border border-zinc-800 bg-zinc-950/70 p-5 shadow-[0_0_40px_rgba(0,0,0,0.35)] backdrop-blur">
            <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">{t("auth.access", language)}</p>
            <h1 className="mt-2 text-3xl font-semibold text-zinc-100">{t("auth.login", language)}</h1>
            <p className="mt-2 text-sm text-zinc-400">{t("auth.loginHint", language)}</p>
            <div className="mt-4 flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setLanguage("cs")}
                className={`rounded-md border px-2 py-1 text-xs ${language === "cs" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"}`}
              >
                CZ
              </button>
              <button
                type="button"
                onClick={() => setLanguage("en")}
                className={`rounded-md border px-2 py-1 text-xs ${language === "en" ? "border-[#00ff00] text-[#00ff00]" : "border-zinc-700 text-zinc-400"}`}
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
        <ToastContainer />
      </div>
    );
  }

  // ── Dashboard ──
  return (
    <ErrorBoundary>
      <div className={BG_GRADIENT}>
        {/* Mobile top bar */}
        <div className="mx-auto mb-4 w-full max-w-350 md:hidden">
          <TopBar
            language={language}
            username={auth.username}
            wsStatus={wsStatus}
            showConnectionStatus={state.phase === "live"}
            onLanguageChange={setLanguage}
            onLogout={logout}
          />
        </div>

        <div className="mx-auto grid w-full max-w-350 gap-4 md:grid-cols-[288px_1fr]">
          <PhaseSidebar phase={state.phase} language={language} raceStatus={state.raceStatus} onPhaseChange={updatePhase} />

          <main className="rounded-2xl border border-zinc-800 bg-zinc-950/60 p-4 md:p-6">
            {/* Desktop top bar */}
            <div className="mb-4 hidden md:block">
              <TopBar
                language={language}
                username={auth.username}
                wsStatus={wsStatus}
                showConnectionStatus={state.phase === "live"}
                onLanguageChange={setLanguage}
                onLogout={logout}
              />
            </div>

            {state.phase === "setup" && (
              <SetupPhase
                config={config}
                language={language}
                loading={!!loading.createGame}
                onCreateGame={createGame}
                onUpdateConfig={updateConfig}
              />
            )}

            {state.phase === "players" && (
              <PlayerRegistration
                players={state.players}
                devices={state.devices}
                teams={state.teams}
                gameMode={config.gameMode}
                language={language}
                actionLoading={loading}
                onAddPlayer={addPlayer}
                onRemovePlayer={removePlayer}
                onAssignTeam={assignPlayerTeam}
                onRefreshDevices={refreshDevices}
                onStartGame={startRace}
              />
            )}

            {state.phase === "live" && (
              <LivePhase
                players={state.players}
                leaderboard={leaderboard}
                killFeed={state.killFeed}
                raceTime={formatRaceTime(state.raceTimeSeconds)}
                raceStatus={state.raceStatus}
                gameMode={config.gameMode}
                language={language}
                onStop={stopRace}
              />
            )}

            {state.phase === "results" && (
              <ResultsPhase
                teamResults={state.teamResults}
                leaderboard={leaderboard}
                matchHistory={matchHistory}
                gameMode={config.gameMode}
                language={language}
                onNewGame={() => updatePhase("setup")}
              />
            )}
          </main>
        </div>
        <ToastContainer />
      </div>
    </ErrorBoundary>
  );
}
