# 🎯 Web Frontend — Enterprise Audit Report

**Projekt:** Astrum Mobility Laser Tag — Web Dashboard
**Stack:** Next.js 16.2.1 · React 19.2.4 · TypeScript 5.9.3 · Tailwind CSS v4
**Celkové hodnocení:** **5.5 / 10**
**Datum:** 2025-07

---

## Obsah

1. [Executive Summary](#1-executive-summary)
2. [Kritické bugy](#2-kritické-bugy)
3. [Architektura &amp; Clean Code](#3-architektura--clean-code)
4. [State Management &amp; Persistence](#4-state-management--persistence)
5. [API vrstva](#5-api-vrstva)
6. [WebSocket](#6-websocket)
7. [Error Handling &amp; UX Feedback](#7-error-handling--ux-feedback)
8. [UI/UX &amp; Responsivita](#8-uiux--responsivita)
9. [Accessibility (a11y)](#9-accessibility-a11y)
10. [Bezpečnost](#10-bezpečnost)
11. [Performance](#11-performance)
12. [Konfigurace &amp; DevOps](#12-konfigurace--devops)
13. [Testy](#13-testy)
14. [Shrnutí nálezů](#14-shrnutí-nálezů)
15. [Doporučený plán oprav](#15-doporučený-plán-oprav)

---

## 1. Executive Summary

Frontend je funkční prototyp s rozumnou volbou stacku (Next.js + TypeScript + Tailwind). Kód má několik silných stránek — striktní TypeScript, konzistentní Tailwind styling, oddělená API vrstva a komponentová struktura. Nicméně z hlediska enterprise standardů existuje řada kritických problémů:

**Co funguje dobře:**

- ✅ TypeScript strict mode
- ✅ Oddělení API klienta do `lib/api.ts`
- ✅ Komponenty v `components/race-control/`
- ✅ Custom hook `useGameData` (i když je příliš velký)
- ✅ Dark theme s gaming estetikou
- ✅ Podpora češtiny i angličtiny

**Co je potřeba opravit:**

- 🔴 Ztráta gameId při refreshi prohlížeče (kritický bug)
- 🔴 Race condition v timer synchronizaci
- 🟠 God-hook anti-pattern (`useGameData` = 478 řádků)
- 🟠 Žádné loading states na async akcích
- 🟠 Žádný error feedback pro uživatele (tiché selhání)
- 🟠 API klient bez timeoutů, retry, typed errors
- 🟡 Nulová accessibility
- 🟡 Chybějící testy

---

## 2. Kritické bugy

### 2.1 🔴 CRITICAL — Ztráta `gameId` při refreshi prohlížeče

**Soubor:** `hooks/useGameData.ts`
**Řádky:** 40–42

Když probíhá živá hra a uživatel refreshne prohlížeč, `gameId` se inicializuje na `null`. Celý game state (hráči, týmy, skóre) je ztracen. WebSocket polling se obnoví s `gameId: null` → žádné real-time updaty.

```ts
// PROBLÉM: Hardcoded null při každém načtení
const initialState: GameState = {
  gameId: null, // Ztraceno při refreshi!
  // ...
};
```

**Dopad:** Hráč uprostřed hry ztratí přístup k probíhající hře.

**Řešení:**

- Persistovat `gameId` a `phase` do `localStorage`
- Při inicializaci číst z `localStorage`
- Přidat `useEffect` pro synchronizaci

### 2.2 🔴 CRITICAL — Race condition v timer synchronizaci

**Soubor:** `hooks/useGameData.ts`
**Řádky:** 197–242

Dva konkurenční mechanismy updatují čas:

1. Server poll (každé 2s) nastaví `raceTimeSeconds` z API
2. Lokální timer (každou 1s) dekrementuje `raceTimeSeconds`

**Scénář:**

1. Server vrátí `raceTimeSeconds: 120`
2. Lokální timer okamžitě nastaví `119`
3. Poll aktualizuje na `118`
4. Timer přepíše na `117` (špatná synchronizace)

**Řešení:** Timer by měl být čistě vizuální (computed value z `raceStartTime`), ne state modifier.

---

## 3. Architektura & Clean Code

### 3.1 🟠 God-Hook anti-pattern

**Soubor:** `hooks/useGameData.ts` — **478 řádků**

Jeden hook spravuje VŠECHNO:

- Auth state (login, logout, token persistence)
- Game config (nastavení hry)
- Game state (fáze, hráči, týmy, skóre)
- Polling logika (interval, WebSocket trigger)
- Timer logika (odpočet)
- Match history (historie her)
- 17+ action funkcí

**Problémy:**

- Nemožné testovat izolovaně
- Změna jednoho concern ovlivní všechny consumers
- ~100 řádků provázaných `useEffect` chains
- Stale closure rizika

**Doporučené rozdělení:**

| Hook                | Odpovědnost                        |
| ------------------- | ----------------------------------- |
| `useAuth`         | Login, logout, token persistence    |
| `useGameState`    | Game fáze, hráči, týmy          |
| `useGamePolling`  | Polling interval + WS trigger       |
| `useRaceTimer`    | Countdown display (computed)        |
| `useMatchHistory` | Historie her                        |
| `useGameActions`  | createGame, addPlayer, startRace... |

### 3.2 🟠 Monolitická stránka

**Soubor:** `app/page.tsx` — **~390 řádků**

Jedna komponenta renderuje:

- Login UI (řádky ~54–106)
- Dashboard header (duplikovaný kód pro mobile/desktop)
- Setup fáze (~151–278)
- Players fáze (~281–294)
- Live fáze (~297–336)
- Results fáze (~339–392)

**Problémy:**

- Language buttons duplikovány (identický kód na 2 místech)
- Logout button duplikován
- Massive ternary nesting pro fáze
- Obtížné code review a údržba

**Doporučená struktura:**

```
app/page.tsx (orchestrátor, ~50 řádků)
├── components/TopBar.tsx (jazyk, logout)
├── components/phases/SetupPhase.tsx
├── components/phases/PlayersPhase.tsx
├── components/phases/LivePhase.tsx
└── components/phases/ResultsPhase.tsx
```

### 3.3 🟡 Duplicitní kód

- Top bar vykreslován 2× (mobile + desktop) s identickým kódem
- Language toggle duplicitní na obou verzích
- Logout button duplikovaný

### 3.4 🟡 Inline i18n ternary

Celý frontend používá inline ternary pro překlad:

```ts
language === "cs" ? "Vytvořit hru" : "Create Game"
```

Toto je rozeseto po celém kódu. Škáluje špatně — přidání třetího jazyka = upravit každý řetězec.

**Řešení:** Extrahovat do translation souboru:

```ts
const t = translations[language];
// pak: t.createGame
```

---

## 4. State Management & Persistence

### 4.1 🔴 Neúplná state persistence

| State         | Persistováno? | Dopad ztráty                    |
| ------------- | :-------------: | -------------------------------- |
| Auth token    | ✅ localStorage | —                               |
| Username      | ✅ localStorage | —                               |
| `gameId`    |       ❌       | Ztráta přístupu k živé hře |
| Game phase    |       ❌       | Neví se, v jaké fázi hra je   |
| Game config   |       ❌       | Ztráta nastavení               |
| Players/Teams |       ❌       | Re-fetch potřeba                |
| Match history |       ❌       | Re-fetch potřeba                |

### 4.2 🟠 Stale closures v akcích

```ts
const addPlayer = async (...) => {
  if (!state.gameId) return; // gameId z closure
  await api.addPlayer(state.gameId, ...);
  // Pokud se gameId změní mezi řádky, fetchneme špatné hráče
  const players = await api.listPlayers(state.gameId);
};
```

**Řešení:** Používat `useRef` pro gameId nebo lokální proměnnou:

```ts
const gameIdRef = useRef(state.gameId);
useEffect(() => { gameIdRef.current = state.gameId }, [state.gameId]);
```

### 4.3 🟠 Memory leak z polling intervalů

`pollGameState` se mění s každým renderem (závisí na `state.gameId`). Když se změní, `useEffect` se re-spustí, ale starý interval nemusí být správně vyčištěn.

**Řešení:** Stabilizovat callback přes `useRef`.

---

## 5. API vrstva

### 5.1 🟠 Žádné typed errors

**Soubor:** `lib/api.ts`

```ts
// Aktuálně:
throw new Error(`HTTP ${res.status}`);
// Nelze rozlišit 401 vs 404 vs 500
```

**Řešení:**

```ts
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    public data?: unknown,
  ) {
    super(`[${status}] ${code}`);
    this.name = "ApiError";
  }
}
```

### 5.2 🟠 Žádný request timeout

Fetch volání nemají timeout. Pokud backend zamrzne, frontend čeká nekonečně.

**Řešení:** `AbortController` s 10s timeout.

### 5.3 🟡 Žádný retry mechanismus

Jednorázový pokus — síťový glitch = permanentní selhání. Pro 5xx chyby by měl být exponential backoff retry.

### 5.4 🟡 Žádná runtime validace responses

Odpovědi z API jsou typované, ale ne validované. Pokud backend změní schéma, frontend tiše spadne.

**Řešení:** Zod/Valibot pro runtime validaci.

### 5.5 🟡 API klient není dělen dle domény

Jeden monolitický `ApiClient` se všemi metodami. Senior přístup:

```
lib/api/
├── client.ts        (base HTTP klient s interceptors)
├── auth.ts          (login, refresh)
├── games.ts         (CRUD hry)
├── players.ts       (CRUD hráčů)
├── devices.ts       (správa zařízení)
└── types.ts         (response/request types)
```

---

## 6. WebSocket

### 6.1 🟠 Chybějící error handling

**Soubor:** `hooks/useWebSocket.ts`

```ts
ws.onerror = () => {
  ws.close(); // Žádné logování, žádný callback
};
```

Uživatel neví, že WebSocket spadl. Žádný connection status indikátor.

### 6.2 🟡 Naivní reconnect

Fixní 2s delay bez exponential backoff. Při výpadku sítě → neustálé pokusy každé 2s → zbytečná zátěž.

### 6.3 🟡 Chybí connection status

Žádný indikátor připojení/odpojení. Uživatel neví, zda dostává real-time data.

---

## 7. Error Handling & UX Feedback

### 7.1 🟠 Tiché selhání všech akcí

**Aktuální stav:** Všechny `catch` bloky logují do `console.error()` nebo nic nedělají.

| Akce            | Co vidí uživatel při chybě |
| --------------- | ------------------------------ |
| Create Game     | Nic (console.error)            |
| Add Player      | Nic (console.error)            |
| Start Race      | Nic (console.error)            |
| Stop Race       | Nic (console.error)            |
| Refresh Devices | Nic (tiché selhání)         |
| Logout          | Nic (prázdný catch)          |

**Řešení:** Toast notification systém + error state per akce.

### 7.2 🟠 Žádné loading states

| Akce            |    Loading indikátor    |
| --------------- | :-----------------------: |
| Create Game     |            ❌            |
| Add Player      | ⚠️ Disabled button only |
| Remove Player   |            ❌            |
| Start Race      |            ❌            |
| Stop Race       |            ❌            |
| Refresh Devices |            ❌            |
| Login           |            ❌            |

Uživatel klikne → nic se neděje → klikne znovu → duplikátní requesty.

### 7.3 🟡 Žádné Error Boundaries

Chybí React Error Boundary. Pokud komponenta crashne, celá stránka spadne bez recovery.

### 7.4 🟡 Chybí empty states

- Leaderboard renderuje prázdnou tabulku (žádná zpráva)
- Kill feed prázdný stav OK ✅
- Player list prázdný stav OK ✅

---

## 8. UI/UX & Responsivita

### 8.1 🟡 Tabulka není scrollable na mobilu

**Soubor:** `components/race-control/Leaderboard.tsx`

`overflow-hidden` ale žádný horizontální scroll. Na malých obrazovkách text přetéká.

**Řešení:** `overflow-x-auto` wrapper.

### 8.2 🟡 Nestandardní max-width

```ts
max-w-350  // Nestandardní Tailwind třída
```

Mělo by být `max-w-7xl` nebo `max-w-[1400px]`.

### 8.3 🟡 Chybějící animace a transitions

- Žádné přechody mezi fázemi (setup → players → live → results)
- Kill feed bez vstupní animace
- Weapon upgrade bez vizuálního efektu
- Score update bez animace

Pro herní dashboard: animace jsou kritické pro engagement.

### 8.4 💡 Návrhy na vylepšení UX

- **Sound effects** pro kill confirmed / weapon upgrade
- **Confetti/particles** při weapon upgrade
- **Pulse animace** na aktivním hráči
- **Countdown overlay** při startu hry (3, 2, 1, GO!)
- **Kill feed** s fade-in animací

---

## 9. Accessibility (a11y)

### 9.1 🟡 Nulové ARIA labels

- Icon-only buttony bez `aria-label`
- Žádné `aria-live` regiony pro real-time updaty (kill feed, leaderboard)
- Formové inputy bez proper `id` + `htmlFor` asociace
- Language toggle není v `role="group"`

### 9.2 🟡 Rozlišení pouze barvou

- Zelená = alive, Červená = dead — ale žádný textový indikátor
- Tým rozlišen pouze barvou (červená/modrá)

### 9.3 🟡 Keyboard navigace

- Žádný focus management mezi fázemi
- Tab order neoptimalizován
- Žádné skip links

---

## 10. Bezpečnost

### 10.1 🟡 Token v localStorage

Bearer token v localStorage je zranitelný vůči XSS. Pro SPA je to akceptabilní, ale mělo by být doplněno o:

- Content Security Policy (CSP) headers
- Token expiry check na klientu

### 10.2 🟡 Žádná kontrola token expiry

Token se persistuje, ale nikdy se nekontroluje `expires_at`. Pokud token vyprší, frontend posílá neplatné requesty bez feedback.

**Řešení:** Auto-logout + refresh token flow.

### 10.3 🟡 Demo credentials v kódu

`admin / admin123` zobrazeno v `LoginPanel.tsx`. OK pro demo, ale musí být odstraněno před produkcí.

### 10.4 🟡 Chybějící security headers

`next.config.ts` je prázdný — chybí:

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `Content-Security-Policy`
- `poweredByHeader: false`

---

## 11. Performance

### 11.1 🟡 N+1 API calls

Po každé akci (addPlayer, removePlayer) se fetchnou VŠICHNI hráči + všechna zařízení:

```ts
const [players, devices] = await Promise.all([
  api.listPlayers(state.gameId),    // ALL players
  api.listAvailableDevices(),       // ALL devices
]);
```

Při 20+ hráčích = 3 API calls na každou akci.

**Řešení:** Optimistic updates nebo WS push pro player updates.

### 11.2 🟡 Žádný code splitting

Celá aplikace v jednom bundle. Chybí:

- `React.lazy()` pro fázové komponenty
- Dynamic imports
- Route-level splitting (jen jedna route ale budoucí rozšíření)

### 11.3 🟡 Hardcoded polling interval

2000ms fixní interval. Měl by být konfigurovatelný a adaptivní (pomalejší síť = delší interval).

---

## 12. Konfigurace & DevOps

### 12.1 🟡 Prázdný next.config.ts

```ts
const nextConfig: NextConfig = {};
```

Chybí: `reactStrictMode`, `compress`, `poweredByHeader`, security headers, image optimization.

### 12.2 🟡 Chybí `.env.example`

`NEXT_PUBLIC_API_URL` se používá v kódu, ale není dokumentováno. Nový vývojář neví, co nastavit.

### 12.3 🟡 Chybí ESLint konfigurace

Žádný `.eslintrc` nebo `eslint.config.js`. Default Next.js lint je minimální.

---

## 13. Testy

### 13.1 🟠 Nulové testy

Žádné unit testy, integration testy, ani E2E testy:

- ❌ Hook testy (useGameData, useWebSocket)
- ❌ Komponent testy (PlayerCard, Leaderboard)
- ❌ API klient testy
- ❌ E2E testy (Playwright/Cypress)

---

## 14. Shrnutí nálezů

| #  | Problém                     | Severita    | Soubor          |
| -- | ---------------------------- | ----------- | --------------- |
| 1  | Ztráta gameId při refreshi | 🔴 CRITICAL | useGameData.ts  |
| 2  | Race condition timer         | 🔴 CRITICAL | useGameData.ts  |
| 3  | God-hook (478 řádků)      | 🟠 HIGH     | useGameData.ts  |
| 4  | Tiché selhání akcí       | 🟠 HIGH     | všude          |
| 5  | Žádné loading states      | 🟠 HIGH     | komponenty      |
| 6  | API bez typed errors         | 🟠 HIGH     | api.ts          |
| 7  | API bez timeout              | 🟠 HIGH     | api.ts          |
| 8  | WS error handling            | 🟠 HIGH     | useWebSocket.ts |
| 9  | Memory leak z pollingu       | 🟠 HIGH     | useGameData.ts  |
| 10 | Nulové testy                | 🟠 HIGH     | —              |
| 11 | Monolitická page.tsx        | 🟡 MEDIUM   | page.tsx        |
| 12 | Inline i18n ternary          | 🟡 MEDIUM   | všude          |
| 13 | Žádné Error Boundaries    | 🟡 MEDIUM   | —              |
| 14 | Nulová accessibility        | 🟡 MEDIUM   | všude          |
| 15 | Chybí security headers      | 🟡 MEDIUM   | next.config.ts  |
| 16 | Token expiry check           | 🟡 MEDIUM   | useGameData.ts  |
| 17 | Tabulka nescrollable         | 🟡 MEDIUM   | Leaderboard.tsx |
| 18 | Chybí .env.example          | 🟡 MEDIUM   | —              |
| 19 | N+1 API calls                | 🟡 MEDIUM   | useGameData.ts  |
| 20 | Žádné animace             | 🔵 LOW      | komponenty      |
| 21 | Unused dependency (clsx)     | 🔵 LOW      | package.json    |
| 22 | Demo credentials v kódu     | 🔵 LOW      | LoginPanel.tsx  |
| 23 | Hardcoded polling interval   | 🔵 LOW      | useGameData.ts  |

---

## 15. Doporučený plán oprav

### Fáze 1 — Kritické bugy (IHNED)

- [ ] Persistovat `gameId` + `phase` do localStorage
- [ ] Opravit race condition v timeru (computed time)
- [ ] Opravit memory leak v polling intervalech

### Fáze 2 — API vrstva

- [ ] `ApiError` class s typed errors
- [ ] Request timeout (AbortController)
- [ ] Rozdělit API klient dle domény (`auth.ts`, `games.ts`, `players.ts`, `devices.ts`)
- [ ] `.env.example` s dokumentací

### Fáze 3 — Hook dekompozice

- [ ] Extrahovat `useAuth` (login, logout, token persistence, token expiry)
- [ ] Extrahovat `useGameState` (fáze, hráči, týmy)
- [ ] Extrahovat `useGamePolling` (interval + WS trigger, stabilní ref)
- [ ] Extrahovat `useRaceTimer` (computed countdown)
- [ ] Extrahovat `useMatchHistory`
- [ ] Extrahovat `useGameActions` (createGame, addPlayer, startRace...)

### Fáze 4 — Error Handling & UX

- [ ] Toast notification systém
- [ ] Loading states na všechny async akce
- [ ] Error Boundary wrapper
- [ ] Connection status indikátor (WS)
- [ ] Error messages viditelné uživateli

### Fáze 5 — Komponenty & Layout

- [ ] Extrahovat `TopBar` (deduplikace)
- [ ] Extrahovat fázové komponenty (`SetupPhase`, `PlayersPhase`, `LivePhase`, `ResultsPhase`)
- [ ] Translation systém (nahradit inline ternary)
- [ ] Empty states pro všechny seznamy

### Fáze 6 — WebSocket vylepšení

- [ ] Exponential backoff reconnect
- [ ] Error callback + logování
- [ ] Connection status state
- [ ] Typed event systém

### Fáze 7 — Bezpečnost & Konfigurace

- [ ] Security headers v next.config.ts
- [ ] Token expiry check + auto-logout
- [ ] CSP headers
- [ ] `reactStrictMode: true`

### Fáze 8 — UI/UX Polish

- [ ] Responsivní tabulky (overflow-x-auto)
- [ ] Animace přechodů mezi fázemi
- [ ] Kill feed animace
- [ ] Weapon upgrade vizuální efekt
- [ ] Accessibility (ARIA labels, keyboard nav, aria-live)

### Fáze 9 — Performance

- [ ] Optimistic updates (addPlayer, removePlayer)
- [ ] Konfigurovatelný polling interval
- [ ] Code splitting (lazy loading fázových komponent)

### Fáze 10 — Testy

- [ ] Unit testy pro hooky (vitest + @testing-library/react)
- [ ] Komponent testy
- [ ] API klient testy (MSW)
- [ ] E2E testy (Playwright)

---

---

## Stav implementace

- [X] Fáze 1 — Kritické bugy
- [X] Fáze 2 — API vrstva
- [X] Fáze 3 — Hook dekompozice
- [X] Fáze 4 — Error Handling & UX
- [X] Fáze 5 — Komponenty & Layout
- [X] Fáze 6 — WebSocket vylepšení
- [X] Fáze 7 — Bezpečnost & Konfigurace
- [X] Fáze 8 — UI/UX Polish
- [X] Fáze 9 — Performance
- [X] Fáze 10 — Testy

---

*Audit provedl: AI Senior Frontend Developer*
*Methodology: Manual code review, line-by-line analysis, Clean Architecture principles*
