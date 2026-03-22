# LaserTag MP135 — Game Display

Event-driven herní displej pro M5Stack MP135 s CAN bus komunikací.

## Architektura

```
CAN bus (can1, 1Mbps)
    │
    ▼
┌──────────────┐     ┌───────────┐
│ CanReceiver  │────▶│ EventBus  │
└──────────────┘     └─────┬─────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌─────────┐
        │GameState │ │ Animator │ │  Audio  │
        └────┬─────┘ └────┬─────┘ └─────────┘
             │            │
             ▼            ▼
        ┌────────────────────┐
        │   Framebuffer      │
        │   (/dev/fb1)       │
        └────────────────────┘
```

### Moduly

| Modul | Popis |
|-------|-------|
| `event_bus.h` | Thread-safe pub/sub event systém |
| `can_receiver` | SocketCAN listener, dekóduje CAN framy → eventy |
| `framebuffer` | RGB565 double-buffered rendering na `/dev/fb1` |
| `animator` | Animační engine (damage flash, kill popup, death screen...) |
| `audio` | WAV přehrávač přes ALSA (`aplay`) |
| `game_state` | Herní logika, state machine (IDLE→COUNTDOWN→ACTIVE→ENDING) |
| `font.h` | 5x7 bitmap font (A-Z, 0-9, symboly) |

### Eventy

Všechny subsystémy komunikují přes `EventBus` — žádné přímé závislosti:

- `GAME_START` / `GAME_END` — řízení hry
- `PLAYER_HIT` — hráč byl trefený (→ damage flash + zvuk)
- `PLAYER_KILL` — +1 kill popup + zvuk
- `PLAYER_DEATH` — death screen + zvuk
- `SCORE_UPDATE` / `AMMO_UPDATE` — HUD update
- `HEARTBEAT` — device online tracking

### Animace

- **Damage Flash** — červený záblesk obrazovky s vignette efektem
- **+1 Kill Popup** — text stoupající nahoru s fade-out
- **Death Screen** — pulsující "DEAD" text přes červený overlay
- **Respawn Countdown** — 3..2..1 s expandujícími čísly
- **Game Countdown** — 5..4..3..2..1..GO! s progress barem
- **Game Over** — dramatický reveal s pulsující hranicí
- **Low Ammo Blink** — oranžový blikající rámeček

## Build

### Na MP135 (nativně)

```bash
sudo bash setup.sh   # instalace závislostí + CAN konfigurace
make native
sudo ./lasertag -p 1
```

### Cross-compile (z macOS/Linux)

```bash
# Potřebuje arm-linux-gnueabihf-g++ toolchain
make
make deploy MP135_HOST=192.168.0.213
```

### Argumenty

```
./lasertag [OPTIONS]
  -p, --player ID    Player ID (1-15, default: 1)
  -c, --can IFACE    CAN rozhraní (default: can1)
  -f, --fb PATH      Framebuffer (default: /dev/fb1)
  -m, --media DIR    Složka se zvuky (default: media)
```

### Systemd služba

```bash
sudo systemctl enable lasertag
sudo systemctl start lasertag
```

## Zvuky

Uložte WAV soubory do `media/`:

| Zvuk | Soubor | Kdy hraje |
|------|--------|-----------|
| Death | `gta-death.wav` | Smrt hráče |
| Kill | `mario-death.wav` | +1 kill |

Vlastní zvuky lze přidat v `src/audio.cpp` → `Audio::Audio()`.

## CAN Protokol

Viz `../esp/can_protocol.h` — ID schéma `0xABC`:
- `A` = kategorie (0=system, 1=combat, 2=game, 3=status)
- `B` = funkce
- `C` = hráč (0=broadcast, 1-15=player)
