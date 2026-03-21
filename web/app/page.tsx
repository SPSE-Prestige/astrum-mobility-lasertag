"use client";

import { useMemo } from "react";
import { LogOut, Settings, SlidersHorizontal, Users } from "lucide-react";
import { GameControls } from "@/components/race-control/GameControls";
import { Leaderboard } from "@/components/race-control/Leaderboard";
import { LoginPanel } from "@/components/race-control/LoginPanel";
import { MatchHistoryManager } from "@/components/race-control/MatchHistoryManager";
import { PhaseSidebar } from "@/components/race-control/PhaseSidebar";
import { PlayerCard } from "@/components/race-control/PlayerCard";
import { useGameData } from "@/hooks/useGameData";

const formatRaceClock = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
    .toString()
    .padStart(2, "0");
  const secs = Math.max(seconds % 60, 0)
    .toString()
    .padStart(2, "0");
  return `${mins}:${secs}`;
};

export default function Home() {
  const {
    config,
    state,
    auth,
    leaderboard,
    matchHistory,
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

  const winner = useMemo(() => {
    return [...state.teamResults].sort((a, b) => b.score - a.score)[0];
  }, [state.teamResults]);

  if (!auth.isAuthenticated) {
    return (
      <div className="min-h-screen bg-[radial-gradient(circle_at_20%_20%,rgba(0,255,0,0.14),transparent_35%),radial-gradient(circle_at_80%_0%,rgba(255,0,0,0.12),transparent_30%),#020303] p-4 md:p-8">
        <div className="mx-auto flex min-h-[80vh] w-full max-w-xl items-center justify-center">
          <div className="w-full rounded-2xl border border-zinc-800 bg-zinc-950/70 p-5 shadow-[0_0_40px_rgba(0,0,0,0.35)] backdrop-blur">
            <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">Race Control Access</p>
            <h1 className="mt-2 text-3xl font-semibold text-zinc-100">Přihlášení do dashboardu</h1>
            <p className="mt-2 text-sm text-zinc-400">Bez přihlášení není přístup k nastavení ani live obrazovce.</p>
            <div className="mt-4">
              <LoginPanel
                isAuthenticated={auth.isAuthenticated}
                username={auth.username}
                error={auth.error}
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
      <div className="mx-auto grid w-full max-w-350 gap-4 md:grid-cols-[288px_1fr]">
        <PhaseSidebar phase={state.phase} onPhaseChange={updatePhase} />

        <main className="rounded-2xl border border-zinc-800 bg-zinc-950/60 p-4 md:p-6">
          <div className="mb-4 flex items-center justify-end gap-3 rounded-xl border border-zinc-800 bg-black/30 px-3 py-2">
            <span className="text-xs uppercase tracking-[0.16em] text-zinc-400">{auth.username}</span>
            <button
              type="button"
              onClick={logout}
              className="inline-flex items-center gap-2 rounded-md border border-zinc-700 bg-zinc-900 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.14em] text-zinc-200 transition hover:border-zinc-500"
            >
              <LogOut className="h-3.5 w-3.5" />
              Odhlásit
            </button>
          </div>

          {state.phase === "setup" && (

            <section className="space-y-6">
              <header>
                <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">Phase 1</p>
                <h1 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">Setup Control</h1>
              </header>


              <div className="grid gap-4 xl:grid-cols-2">
                <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                  <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                    <Settings className="h-4 w-4 text-[#00ff00]" />
                    Game Parameters
                  </h3>
                  <div className="grid gap-3">
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      Name
                      <input
                        value={config.gameName}
                        onChange={(event) => updateConfig({ ...config, gameName: event.target.value })}
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                      />
                    </label>
                    <div className="grid grid-cols-2 gap-3">
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        Duration (min)
                        <input
                          type="number"
                          min={5}
                          value={config.durationMinutes}
                          onChange={(event) => updateConfig({ ...config, durationMinutes: Number(event.target.value) })}
                          className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 outline-none ring-[#00ff00] transition focus:ring-1"
                        />
                      </label>
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        Herní režim
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
                          <option value="team">Týmová hra</option>
                          <option value="ffa">Všichni proti všem</option>
                        </select>
                      </label>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                      <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                        Teams
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
                          ? "V režimu FFA se týmy nepoužívají."
                          : "V týmovém režimu je možné přiřadit vozíky do týmů."}
                      </div>
                    </div>
                    <label className="mt-1 inline-flex items-center justify-between rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-200">
                      Friendly Fire
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
                    Weapon & Player Tuning
                  </h3>
                  <div className="grid grid-cols-2 gap-3">
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      Damage
                      <input
                        type="number"
                        value={config.weaponTuning.damage}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            weaponTuning: { ...config.weaponTuning, damage: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      Fire Rate
                      <input
                        type="number"
                        value={config.weaponTuning.fireRate}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            weaponTuning: { ...config.weaponTuning, fireRate: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      Reload Time
                      <input
                        type="number"
                        value={config.weaponTuning.reloadTime}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            weaponTuning: { ...config.weaponTuning, reloadTime: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                    <label className="text-xs uppercase tracking-[0.14em] text-zinc-500">
                      HP
                      <input
                        type="number"
                        value={config.playerTuning.hp}
                        onChange={(event) =>
                          updateConfig({
                            ...config,
                            playerTuning: { ...config.playerTuning, hp: Number(event.target.value) },
                          })
                        }
                        className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-100"
                      />
                    </label>
                    <label className="col-span-2 text-xs uppercase tracking-[0.14em] text-zinc-500">
                      Respawn Delay
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
                  </div>
                </article>
              </div>

              <article className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
                <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
                  <Users className="h-4 w-4 text-zinc-200" />
                  Cart Lobby & Assignment
                </h3>
                <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                  {state.players.map((player) => (
                    <div key={player.id} className="rounded-lg border border-zinc-800 bg-zinc-900/70 p-3">
                      <p className="text-sm font-semibold text-zinc-100">{player.name}</p>
                      <p className="mt-1 text-xs text-zinc-500">
                        Cart: {player.cartConnected ? "Connected" : "Offline"}
                      </p>
                      <select
                        value={config.gameMode === "ffa" ? "Unassigned" : player.team}
                        disabled={config.gameMode === "ffa"}
                        className="mt-2 w-full rounded-md border border-zinc-700 bg-black px-2 py-1 text-xs text-zinc-200"
                        onChange={(event) => assignPlayerTeam(player.id, event.target.value)}
                      >
                        <option>Neon Green</option>
                        <option>Neon Red</option>
                        <option>Unassigned</option>
                      </select>
                    </div>
                  ))}
                </div>
              </article>

              <MatchHistoryManager
                items={matchHistory}
                onUpdate={updateMatchHistoryItem}
                onDelete={deleteMatchHistoryItem}
              />
            </section>
          )}

          {state.phase === "live" && (
            <section className="space-y-4">
              <GameControls
                raceTime={formatRaceClock(state.raceTimeSeconds)}
                raceStatus={state.raceStatus}
                onStart={startRace}
                onPause={pauseRace}
                onStop={stopRace}
              />

              <div className="grid gap-4 xl:grid-cols-[2fr_1fr]">
                <div className="grid gap-3 md:grid-cols-2">
                  {state.players.map((player) => (
                    <PlayerCard key={player.id} player={player} />
                  ))}
                </div>

                <div className="space-y-4">
                  <Leaderboard players={leaderboard} gameMode={config.gameMode} />
                  <section className="rounded-2xl border border-zinc-800 bg-zinc-950/70 p-4">
                    <p className="mb-3 text-xs uppercase tracking-[0.2em] text-zinc-500">Kill Feed</p>
                    <div className="space-y-2">
                      {state.killFeed.map((entry) => (
                        <div key={entry.id} className="rounded-md border border-zinc-800 bg-black/40 px-3 py-2 text-xs text-zinc-300">
                          <span className="mr-2 text-[#00ff00]">[{entry.timestamp}]</span>
                          {entry.message}
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
                <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">Phase 3</p>
                <h2 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">Final Results</h2>
              </header>

              <article className="rounded-2xl border border-[#00ff00]/40 bg-[#00ff00]/5 p-5">
                <p className="text-xs uppercase tracking-[0.18em] text-zinc-400">
                  {config.gameMode === "ffa" ? "Winning Player" : "Winning Team"}
                </p>
                <h3 className="mt-2 text-2xl font-semibold text-[#00ff00]">{winner?.team ?? "TBD"}</h3>
              </article>

              <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {state.teamResults.map((teamResult) => (
                  <article key={teamResult.team} className="rounded-xl border border-zinc-800 bg-black/40 p-4 text-zinc-200">
                    <h4 className="text-lg font-semibold">{teamResult.team}</h4>
                    <div className="mt-3 space-y-2 text-sm">
                      <p>
                        Score: <span className="text-[#ff0000]">{teamResult.score}</span>
                      </p>
                      <p>
                        Accuracy: <span className="text-[#00ff00]">{teamResult.accuracy}%</span>
                      </p>
                      <p>
                        Damage Dealt: <span className="text-[#ff0000]">{teamResult.damageDealt}</span>
                      </p>
                    </div>
                  </article>
                ))}
              </div>
            </section>
          )}
        </main>
      </div>
    </div>
  );
}
