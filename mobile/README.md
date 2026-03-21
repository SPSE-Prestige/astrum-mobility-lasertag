# LaserTag Mobile App

Mobilní aplikace pro LaserTag race control systém - postavená s **Expo** a **React Native**.

## Getting Started

```bash
# 1. Instalace závislostí
cd mobile
npm install

# 2. Spuštění na vybraném zařízení
npm start          # Zobrazí QR code, skenuj v Expo Go
npm run ios        # Spuštění na iOS simulátoru
npm run android    # Spuštění na Android emulátoru
npm run web        # Spuštění v prohlížeči
```

Po spuštění skenuj QR code v aplikaci **Expo Go** nebo výběr platformy.

### Struktura projektu

```
mobile/
├── README.md
├── package.json              # Dependencies a scripts
├── app.json                  # Expo konfiguraci
├── tsconfig.json             # TypeScript konfigurace
├── .eslintrc.json            # Linting pravidla
├── .gitignore
├── app/                      # Expo Router - hlavní aplikační kód
│   └── index.tsx            # Domovská stránka
└── assets/                  # Ikony, splash screeny (TBD)
```

## Technologie

- **Expo** – Framework pro React Native
- **React Native** – Cross-platform mobilní vývoj
- **TypeScript** – Typová bezpečnost
- **NativeWind** – Tailwind CSS pro React Native
- **Expo Router** – Navigace (app-based routing)

## Plány

- [x] Inicializace Expo projektu
- [ ] Nastavení navigace (Expo Router)
- [ ] Propojení s webovým backend (API)
- [ ] Vytvoření UI pro mobilní rozhraní (Login, Dashboard)
- [ ] Synchronizace stavu s webem (WebSocket / REST)
- [ ] Testing na fyzických zařízeních

## Poznámky

- Aplikace sdílí TypeScript konfiguraci a coding standards s [web/](../web/)
- Dark mode je výchozím nastavením (ref: `#020303` background)
- Pro development: `expo start` a pak `w` pro web, `i` pro iOS, `a` pro Android

Více informací viz [web/README.md](../web/README.md) pro detaily o architektuře projektu.
