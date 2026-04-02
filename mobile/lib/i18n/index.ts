export type Language = "cs" | "en";

const translations = {
  cs: {
    "app.title": "LASER TAG",
    "app.subtitle": "Výsledky hry",
    "login.title": "VÝSLEDKY HRY",
    "login.subtitle": "Zadejte kód relace pro zobrazení výsledků",
    "login.placeholder": "Zadejte kód",
    "login.button": "Zobrazit výsledky",
    "login.loading": "Načítání...",
    "login.error.empty": "Zadejte prosím kód relace",
    "login.error.notFound": "Relace nenalezena",
    "login.error.network": "Chyba připojení k serveru",
    "login.error.generic": "Došlo k chybě",

    "results.title": "Výsledky",
    "results.gameCode": "Kód hry",
    "results.team": "Tým",
    "results.status.running": "ŽIVĚ",
    "results.status.finished": "DOKONČENO",
    "results.status.lobby": "ČEKÁNÍ",

    "stats.score": "Skóre",
    "stats.kills": "Zásahy",
    "stats.deaths": "Smrti",
    "stats.shots": "Výstřely",
    "stats.accuracy": "Přesnost",
    "stats.kd": "K/D Poměr",
    "stats.streak": "Série",
    "stats.weapon": "Zbraň",
    "stats.bestStreak": "Nejlepší série",
    "stats.weaponLevel": "Úroveň zbraně",
    "stats.shotsFired": "Celkem výstřelů",

    "heatmap.title": "Zóny zásahů",
    "heatmap.head": "Hlava",
    "heatmap.chest": "Hrudník",
    "heatmap.back": "Záda",
    "heatmap.shoulders": "Ramena",
    "heatmap.weapon": "Zbraň",

    "leaderboard.title": "Žebříček",
    "leaderboard.rank": "#",
    "leaderboard.player": "Hráč",
    "leaderboard.score": "Skóre",
    "leaderboard.kd": "K/D",

    "killfeed.title": "Poslední události",
    "killfeed.killed": "eliminoval/a",
    "killfeed.empty": "Zatím žádné události",

    "common.refresh": "Obnovit",
    "common.back": "Zpět",
    "common.loading": "Načítání...",
    "common.noData": "Žádná data",
    "common.you": "VY",

    "lang.cs": "CZ",
    "lang.en": "EN",
  },
  en: {
    "app.title": "LASER TAG",
    "app.subtitle": "Game Results",
    "login.title": "GAME RESULTS",
    "login.subtitle": "Enter your session code to view results",
    "login.placeholder": "Enter code",
    "login.button": "View Results",
    "login.loading": "Loading...",
    "login.error.empty": "Please enter a session code",
    "login.error.notFound": "Session not found",
    "login.error.network": "Server connection error",
    "login.error.generic": "An error occurred",

    "results.title": "Results",
    "results.gameCode": "Game Code",
    "results.team": "Team",
    "results.status.running": "LIVE",
    "results.status.finished": "FINISHED",
    "results.status.lobby": "LOBBY",

    "stats.score": "Score",
    "stats.kills": "Kills",
    "stats.deaths": "Deaths",
    "stats.shots": "Shots",
    "stats.accuracy": "Accuracy",
    "stats.kd": "K/D Ratio",
    "stats.streak": "Streak",
    "stats.weapon": "Weapon",
    "stats.bestStreak": "Best Streak",
    "stats.weaponLevel": "Weapon Level",
    "stats.shotsFired": "Total Shots",

    "heatmap.title": "Hit Zones",
    "heatmap.head": "Head",
    "heatmap.chest": "Chest",
    "heatmap.back": "Back",
    "heatmap.shoulders": "Shoulders",
    "heatmap.weapon": "Weapon",

    "leaderboard.title": "Leaderboard",
    "leaderboard.rank": "#",
    "leaderboard.player": "Player",
    "leaderboard.score": "Score",
    "leaderboard.kd": "K/D",

    "killfeed.title": "Recent Events",
    "killfeed.killed": "eliminated",
    "killfeed.empty": "No events yet",

    "common.refresh": "Refresh",
    "common.back": "Back",
    "common.loading": "Loading...",
    "common.noData": "No data",
    "common.you": "YOU",

    "lang.cs": "CZ",
    "lang.en": "EN",
  },
} as const;

export type TranslationKey = keyof (typeof translations)["cs"];

export function t(key: TranslationKey, lang: Language): string {
  return translations[lang][key] ?? key;
}

export default translations;
