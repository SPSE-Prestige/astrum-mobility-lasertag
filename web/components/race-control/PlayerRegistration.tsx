"use client";

import { useState } from "react";
import { Cpu, Play, Plus, RefreshCw, Trash2, Users } from "lucide-react";
import type { Device, GameMode, Player, Team } from "@/types/game";
import type { Language } from "@/types/i18n";

interface PlayerRegistrationProps {
  players: Player[];
  devices: Device[];
  teams: Team[];
  gameMode: GameMode;
  language: Language;
  actionLoading?: Record<string, boolean>;
  onAddPlayer: (deviceId: string, nickname: string, teamId?: string) => Promise<void>;
  onRemovePlayer: (playerId: string) => Promise<void>;
  onAssignTeam: (playerId: string, teamId: string | null) => Promise<void>;
  onRefreshDevices: () => Promise<void>;
  onStartGame: () => Promise<void>;
}

export const PlayerRegistration = ({
  players,
  devices,
  teams,
  gameMode,
  language,
  actionLoading = {},
  onAddPlayer,
  onRemovePlayer,
  onAssignTeam,
  onRefreshDevices,
  onStartGame,
}: PlayerRegistrationProps) => {
  const [nickname, setNickname] = useState("");
  const [selectedDevice, setSelectedDevice] = useState("");
  const [selectedTeam, setSelectedTeam] = useState("");
  const [loading, setLoading] = useState(false);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!nickname.trim() || !selectedDevice) return;
    setLoading(true);
    try {
      await onAddPlayer(selectedDevice, nickname.trim(), selectedTeam || undefined);
      setNickname("");
      setSelectedDevice("");
      setSelectedTeam("");
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="space-y-6">
      <header className="flex items-center justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">
            {language === "cs" ? "Fáze 2" : "Phase 2"}
          </p>
          <h1 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">
            {language === "cs" ? "Hráči a Lobby" : "Players & Lobby"}
          </h1>
        </div>
        <button
          type="button"
          onClick={onStartGame}
          disabled={players.length < 2 || actionLoading.startRace}
          className="rounded-xl border-2 border-[#00ff00] bg-[#00ff00]/20 px-6 py-3 text-lg font-semibold uppercase tracking-[0.16em] text-[#00ff00] shadow-[0_0_20px_rgba(0,255,0,0.3)] transition hover:bg-[#00ff00]/30 hover:shadow-[0_0_30px_rgba(0,255,0,0.4)] disabled:opacity-40 disabled:cursor-not-allowed"
        >
          {actionLoading.startRace ? (
            <RefreshCw className="mr-2 inline h-5 w-5 animate-spin" />
          ) : (
            <Play className="mr-2 inline h-5 w-5" />
          )}
          {language === "cs" ? "Spustit hru" : "Start Game"}
        </button>
      </header>

      <div className="grid gap-6 xl:grid-cols-2">
        {/* Add Player */}
        <div className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
          <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <Plus className="h-4 w-4 text-[#00ff00]" />
            {language === "cs" ? "Přidat hráče" : "Add Player"}
          </h3>

          <form onSubmit={handleAdd} className="space-y-3">
            <input
              type="text"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              placeholder={language === "cs" ? "Přezdívka hráče" : "Player Nickname"}
              className="w-full rounded-lg border border-zinc-800 bg-zinc-950 px-4 py-2 text-zinc-100 placeholder-zinc-600 focus:border-zinc-700 focus:outline-none focus:ring-1 focus:ring-zinc-700"
            />

            <div className="flex items-end gap-2">
              <label className="flex-1 text-xs uppercase tracking-[0.14em] text-zinc-500">
                {language === "cs" ? "Zařízení" : "Device"}
                <select
                  value={selectedDevice}
                  onChange={(e) => setSelectedDevice(e.target.value)}
                  className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:ring-1 focus:ring-zinc-700"
                >
                  <option value="">{language === "cs" ? "Vyberte zařízení" : "Select device"}</option>
                  {devices.map((d) => (
                    <option key={d.id} value={d.deviceId}>
                      {d.deviceId} ({d.status})
                    </option>
                  ))}
                </select>
              </label>
              <button
                type="button"
                onClick={onRefreshDevices}
                disabled={actionLoading.refreshDevices}
                className="mb-0.5 rounded-lg border border-zinc-700 bg-zinc-900 p-2 text-zinc-400 transition hover:border-zinc-500 hover:text-zinc-200 disabled:opacity-50"
                title={language === "cs" ? "Obnovit zařízení" : "Refresh devices"}
              >
                <RefreshCw className={`h-4 w-4 ${actionLoading.refreshDevices ? "animate-spin" : ""}`} />
              </button>
            </div>

            {gameMode === "team" && teams.length > 0 && (
              <label className="block text-xs uppercase tracking-[0.14em] text-zinc-500">
                {language === "cs" ? "Tým" : "Team"}
                <select
                  value={selectedTeam}
                  onChange={(e) => setSelectedTeam(e.target.value)}
                  className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:ring-1 focus:ring-zinc-700"
                >
                  <option value="">{language === "cs" ? "Bez týmu" : "No team"}</option>
                  {teams.map((t) => (
                    <option key={t.id} value={t.id}>
                      {t.name}
                    </option>
                  ))}
                </select>
              </label>
            )}

            <button
              type="submit"
              disabled={!nickname.trim() || !selectedDevice || loading}
              className="flex w-full items-center justify-center gap-2 rounded-lg bg-zinc-100 px-4 py-2 font-medium text-zinc-950 hover:bg-zinc-300 disabled:opacity-50 transition-all active:scale-95"
            >
              <Plus className="h-4 w-4" />
              {language === "cs" ? "Přidat do hry" : "Add to Game"}
            </button>
          </form>

          {devices.length === 0 && (
            <p className="mt-3 text-xs text-zinc-500 italic">
              {language === "cs"
                ? "Žádná dostupná zařízení. Zapněte M5Stack jednotky."
                : "No available devices. Power on M5Stack units."}
            </p>
          )}
        </div>

        {/* Current Players */}
        <div className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
          <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <Users className="h-4 w-4 text-zinc-200" />
            {language === "cs" ? `Hráči ve hře (${players.length})` : `Players in Game (${players.length})`}
          </h3>

          <div className="space-y-2 max-h-80 overflow-y-auto pr-2">
            {players.length === 0 && (
              <div className="flex h-32 flex-col items-center justify-center border-2 border-dashed border-zinc-800/50 rounded-lg text-zinc-600">
                <Cpu className="mb-2 h-6 w-6 opacity-30" />
                <p className="text-xs italic">
                  {language === "cs" ? "Zatím žádní hráči" : "No players yet"}
                </p>
              </div>
            )}
            {players.map((player) => (
              <div
                key={player.id}
                className="group flex items-center justify-between rounded-lg border border-zinc-800 bg-zinc-900/50 p-3 text-sm hover:border-zinc-700 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="h-2 w-2 rounded-full bg-green-500" />
                  <div>
                    <span className="font-medium text-zinc-200">{player.name}</span>
                    <span className="ml-2 text-xs text-zinc-500">{player.deviceId}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {gameMode === "team" && (
                    <select
                      value={player.teamId ?? ""}
                      onChange={(e) => onAssignTeam(player.id, e.target.value || null)}
                      className="rounded border border-zinc-700 bg-zinc-900 px-2 py-1 text-xs text-zinc-300"
                    >
                      <option value="">{language === "cs" ? "Bez týmu" : "No team"}</option>
                      {teams.map((t) => (
                        <option key={t.id} value={t.id}>
                          {t.name}
                        </option>
                      ))}
                    </select>
                  )}
                  <button
                    onClick={() => onRemovePlayer(player.id)}
                    className="flex h-7 w-7 items-center justify-center rounded border border-zinc-800 bg-zinc-900 text-zinc-500 hover:border-red-900/50 hover:bg-red-900/10 hover:text-red-400 transition-colors"
                    title={language === "cs" ? "Odebrat hráče" : "Remove player"}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
};