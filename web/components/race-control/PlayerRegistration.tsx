"use client";

import { useState } from "react";
import { Copy, Plus, Trash2, UserPlus, Users, User, Zap } from "lucide-react";
import type { GameMode, Player, RegisteredPlayer } from "@/types/game";
import type { Language } from "@/types/i18n";

interface PlayerRegistrationProps {
  activePlayers: Player[];
  registeredPlayers: RegisteredPlayer[];
  activeRoster: string[];
  gameMode: GameMode;
  language: Language;
  onRegister: (name: string, type?: "guest" | "registered") => void;
  onDeleteRegistry: (id: string) => void;
  onToggleRoster: (playerId: string) => void;
}

export const PlayerRegistration = ({
  activePlayers = [],
  registeredPlayers = [],
  activeRoster = [],
  gameMode,
  language,
  onRegister,
  onDeleteRegistry,
  onToggleRoster,
}: PlayerRegistrationProps) => {
  const [newName, setNewName] = useState("");
  const [registrationMode, setRegistrationMode] = useState<"registered" | "guest">("registered");

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    if (registrationMode === "registered") {
      if (!newName.trim()) return;
      onRegister(newName.trim(), "registered");
      setNewName("");
    } else {
      onRegister("", "guest");
    }
  };

  const copyCode = (code: string) => {
    navigator.clipboard.writeText(code);
  };

  return (
    <section className="space-y-6">
      <header>
        <p className="text-xs uppercase tracking-[0.2em] text-zinc-500">
          {language === "cs" ? "Správa hráčů" : "Player Management"}
        </p>
        <h1 className="mt-1 text-3xl font-semibold text-zinc-100 md:text-4xl">
          {language === "cs" ? "Hráči a Lobby" : "Players & Lobby"}
        </h1>
      </header>

      <div className="grid gap-6 xl:grid-cols-2">
        {/* Registration Section */}
        <div className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
          <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <UserPlus className="h-4 w-4 text-zinc-200" />
            {language === "cs" ? "Registrace nových hráčů" : "Player Registration"}
          </h3>

          <div className="mb-4 flex w-full rounded-lg bg-zinc-900/50 p-1">
            <button
              onClick={() => setRegistrationMode("registered")}
              className={`flex-1 rounded-md py-1.5 text-xs font-medium transition ${
                registrationMode === "registered"
                  ? "bg-zinc-800 text-zinc-100 shadow-sm"
                  : "text-zinc-500 hover:text-zinc-300"
              }`}
            >
              <span className="flex items-center justify-center gap-2">
                <User className="h-3.5 w-3.5" />
                {language === "cs" ? "Registrace" : "Standard"}
              </span>
            </button>
            <button
              onClick={() => setRegistrationMode("guest")}
              className={`flex-1 rounded-md py-1.5 text-xs font-medium transition ${
                registrationMode === "guest"
                  ? "bg-zinc-800 text-zinc-100 shadow-sm"
                  : "text-zinc-500 hover:text-zinc-300"
              }`}
            >
              <span className="flex items-center justify-center gap-2">
                <Zap className="h-3.5 w-3.5" />
                {language === "cs" ? "Host (Superhrdina)" : "Guest (Superhero)"}
              </span>
            </button>
          </div>

          <form onSubmit={handleRegister} className="flex gap-3">
            {registrationMode === "registered" ? (
              <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder={language === "cs" ? "Jméno hráče" : "Player Name"}
                className="flex-1 rounded-lg border border-zinc-800 bg-zinc-950 px-4 py-2 text-zinc-100 placeholder-zinc-600 focus:border-zinc-700 focus:outline-none focus:ring-1 focus:ring-zinc-700"
              />
            ) : (
              <div className="flex-1 rounded-lg border border-zinc-800 bg-zinc-900/40 px-4 py-2 text-sm text-zinc-500 italic flex items-center">
                {language === "cs"
                  ? "Bude vygenerována přezdívka superhrdiny"
                  : "Superhero nickname will be generated"}
              </div>
            )}
            <button
              type="submit"
              disabled={registrationMode === "registered" && !newName.trim()}
              className="flex items-center gap-2 rounded-lg bg-zinc-100 px-4 py-2 font-medium text-zinc-950 hover:bg-zinc-300 disabled:opacity-50 transition-all active:scale-95"
            >
              <Plus className="h-4 w-4" />
              {language === "cs" ? "Přidat" : "Add"}
            </button>
          </form>
        </div>

        {/* Lobby Section */}
        <div className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
          <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
            <Users className="h-4 w-4 text-zinc-200" />
            {language === "cs" ? "Lobby aktivních hráčů" : "Active Players Lobby"}
          </h3>
          <div className="space-y-2 max-h-55 overflow-y-auto pr-2 custom-scrollbar">
             {!activePlayers.length && (
                 <div className="flex h-32 flex-col items-center justify-center border-2 border-dashed border-zinc-800/50 rounded-lg text-zinc-600">
                     <p className="text-xs italic">
                         {language === "cs" ? "Lobby je prázdné" : "Lobby empty"}
                     </p>
                 </div>
             )}
             {activePlayers.map((player) => (
                <div key={player.id} className="group flex items-center justify-between rounded-lg border border-zinc-800 bg-zinc-900/50 p-3 text-sm hover:border-zinc-700 transition-colors">
                    <div className="flex items-center gap-3">
                        <div className={`h-2 w-2 rounded-full ${player.cartConnected ? "bg-green-500" : "bg-red-500"}`} />
                        <span className="font-medium text-zinc-200">{player.name}</span>
                    </div>
                    
                    <button 
                        onClick={() => onToggleRoster(player.id)} 
                        className="flex h-7 w-7 items-center justify-center rounded border border-zinc-800 bg-zinc-900 text-zinc-500 hover:border-red-900/50 hover:bg-red-900/10 hover:text-red-400 transition-colors"
                        title={language === "cs" ? "Odebrat z lobby" : "Remove from lobby"}
                    >
                        <Trash2 className="h-3.5 w-3.5" />
                    </button>
                </div>
             ))}
          </div>
        </div>
      </div>

      {/* Registered Players List */}
      <div className="rounded-2xl border border-zinc-800 bg-black/40 p-4">
        <h3 className="mb-4 inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.15em] text-zinc-300">
          <Users className="h-4 w-4 text-zinc-200" />
          {language === "cs" ? "Seznam registrovaných hráčů" : "Registered Players List"}
        </h3>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {registeredPlayers.map((player) => {
            const isActive = activeRoster.includes(player.id);
            return (
              <div
                key={player.id}
                onClick={() => onToggleRoster(player.id)}
                className={`group relative flex cursor-pointer flex-col gap-3 rounded-xl border p-4 transition-all active:scale-[0.98] ${
                  isActive
                    ? "border-[#00ff00] bg-[#00ff00]/5 ring-1 ring-[#00ff00] shadow-[0_0_15px_rgba(0,255,0,0.1)]"
                    : "border-zinc-800 bg-zinc-900/30 hover:bg-zinc-900/60 hover:border-zinc-700"
                }`}
              >
                <div className="flex items-start justify-between">
                  <div>
                    <div className="flex flex-wrap items-center gap-2">
                        <p className={`font-semibold ${isActive ? "text-[#00ff00]" : "text-zinc-100"}`}>{player.name}</p>
                        {player.type === "guest" && (
                            <span className="rounded border border-yellow-500/20 bg-yellow-500/10 px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wider text-yellow-500">
                                {language === "cs" ? "Host" : "Guest"}
                            </span>
                        )}
                         {player.type === "registered" && (
                            <span className="rounded border border-blue-500/20 bg-blue-500/10 px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wider text-blue-500">
                                {language === "cs" ? "Reg" : "Reg"}
                            </span>
                        )}
                    </div>
                    <p className="mt-1 text-[10px] text-zinc-600 uppercase tracking-widest">
                        {new Date(player.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </p>
                  </div>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onDeleteRegistry(player.id);
                    }}
                    className="rounded p-1.5 text-zinc-600 opacity-0 transition-opacity hover:bg-red-500/10 hover:text-red-500 group-hover:opacity-100"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>

                <div className="mt-auto flex items-center justify-between rounded-lg border border-zinc-800/50 bg-black/40 px-3 py-2">
                  <span className="text-[10px] font-bold tracking-widest text-zinc-600">CODE</span>
                  <div className="flex items-center gap-3">
                    <span className="font-mono text-xl font-bold tracking-[0.2em] text-[#00ff00]">{player.code}</span>
                    <button
                      type="button"
                      onClick={(e) => {
                        e.stopPropagation();
                        copyCode(player.code);
                      }}
                      className="text-zinc-600 hover:text-zinc-300 transition-colors"
                    >
                      <Copy className="h-3.5 w-3.5" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}

          {registeredPlayers.length === 0 && (
            <div className="col-span-full flex flex-col items-center justify-center py-12 text-zinc-600 border-2 border-dashed border-zinc-800/50 rounded-xl">
              <Users className="mb-2 h-8 w-8 opacity-20" />
              <p>{language === "cs" ? "Zatím žádní registrovaní hráči" : "No registered players yet"}</p>
            </div>
          )}
        </div>
      </div>
    </section>
  );
};