import type { Language } from "@/types/i18n";

const translations = {
  // ── Phases ──
  "phase.1": { cs: "Fáze 1", en: "Phase 1" },
  "phase.setup": { cs: "Nastavení hry", en: "Game Setup" },
  "phase.players": { cs: "Registrace hráčů", en: "Player Registration" },
  "phase.live": { cs: "Živá hra", en: "Live Game" },
  "phase.results": { cs: "Finální výsledky", en: "Final Results" },

  // ── Setup ──
  "setup.createGame": { cs: "Vytvořit hru", en: "Create Game" },
  "setup.creating": { cs: "Vytváření...", en: "Creating..." },
  "setup.params": { cs: "Parametry hry", en: "Game Parameters" },
  "setup.additional": { cs: "Dodatečná nastavení", en: "Additional Settings" },
  "setup.name": { cs: "Název", en: "Name" },
  "setup.duration": { cs: "Délka (min)", en: "Duration (min)" },
  "setup.mode": { cs: "Herní režim", en: "Game Mode" },
  "setup.modeTeam": { cs: "Týmová hra", en: "Team Game" },
  "setup.modeFFA": { cs: "Každý proti každému", en: "Free For All" },
  "setup.teams": { cs: "Týmy", en: "Teams" },
  "setup.respawnDelay": { cs: "Prodleva respawnu (s)", en: "Respawn Delay (s)" },
  "setup.maxPlayers": { cs: "Max hráčů", en: "Max Players" },
  "setup.killsPerUpgrade": { cs: "Killů na upgrade zbraně", en: "Kills Per Weapon Upgrade" },
  "setup.killsPerUpgradeHint": { cs: "0 = upgrade zbraně vypnut", en: "0 = weapon upgrades disabled" },

  // ── Live ──
  "live.killFeed": { cs: "Feed událostí", en: "Kill Feed" },
  "live.noEvents": { cs: "Zatím žádné události", en: "No events yet" },

  // ── Results ──
  "results.winningTeam": { cs: "Vítězný tým", en: "Winning Team" },
  "results.winningPlayer": { cs: "Vítězný hráč", en: "Winning Player" },
  "results.tbd": { cs: "Neurčeno", en: "TBD" },
  "results.score": { cs: "Skóre", en: "Score" },
  "results.newGame": { cs: "Nová hra", en: "New Game" },
  "results.history": { cs: "Historie her", en: "Match History" },

  // ── Auth ──
  "auth.access": { cs: "Přístup do řízení závodu", en: "Race Control Access" },
  "auth.login": { cs: "Přihlášení do dashboardu", en: "Dashboard Login" },
  "auth.loginHint": { cs: "Bez přihlášení není přístup k nastavení ani živé obrazovce.", en: "Without login, settings and live screens are not accessible." },
  "auth.logout": { cs: "Odhlásit", en: "Logout" },
} as const;

type TranslationKey = keyof typeof translations;

export function t(key: TranslationKey, lang: Language): string {
  return translations[key][lang];
}
