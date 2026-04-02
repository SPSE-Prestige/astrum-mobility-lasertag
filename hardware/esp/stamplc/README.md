# ESP Firmware — STAMPLC (StampS3A, ESP32-S3FN8)

Firmware for the M5Stack STAMPLC unit.
IR shooting (Port A + second emitter Port B at level 5), IR hit detection (Port C dual sensors), CAN bus, MQTT, WiFi, weapon level system, team support.

---

## Project structure

```
stamplc/
├── stamplc.ino              ← main sketch
├── pins.h                   ← hardware pin config
├── credentials.h            ← WiFi & MQTT credentials
├── README.md
├── src/
│   ├── ir_controller.h/cpp  ← IRSender (Port A/B TX) + IRReceiver (Port C RX)
│   ├── can_bus.h/cpp
│   ├── can_protocol.h
│   ├── can_debug.h
│   ├── game_state.h         ← kills, weapon level, callbacks
│   └── mqtt_client.h/cpp
├── test/
│   └── test.ino             ← AUnit tests
└── setupid/
    └── setupid.ino          ← flash player ID utility
```

---

## First-time setup — save player ID to flash

Before flashing the main firmware, write a unique player ID (0–15) into NVS flash. This survives power cycles and firmware updates.

1. Open `setupid/setupid.ino` in Arduino IDE
2. Set `PLAYER_IDD` to the device's unique number:
   ```cpp
   const int PLAYER_IDD = 1; // unique per device
   ```
3. Flash to device — saves the ID to NVS and prints confirmation over Serial
4. Only needed **once per device**

Then flash `stamplc.ino`.

---

## Credentials

Edit `credentials.h` before flashing:

```cpp
#define WIFI_SSID     "your_network"
#define WIFI_PASSWORD "your_password"
#define MQTT_HOST     "192.168.x.x"
#define MQTT_PORT     1883
```

---

## How to flash

**Requirements:** Arduino IDE 2.x

**Step 1 — Add ESP32 board support**

`File → Preferences` → Additional boards manager URLs:
```
https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json
```
`Tools → Board → Boards Manager` → search `esp32` → install **esp32 by Espressif**.

**Step 2 — Install libraries**

`Tools → Manage Libraries`, install:
- `M5Unified`
- `IRremoteESP8266`
- `PubSubClient`
- `ArduinoJson`

**Step 3 — Select board**

`Tools → Board → esp32 → ESP32S3 Dev Module`

Settings:
- Flash Size: `8MB`
- Partition Scheme: `8M Flash`
- USB CDC On Boot: `Enabled`

**Step 4 — Flash**

`File → Open` → select `stamplc.ino`, connect via USB-C, select port under `Tools → Port`, click **Upload**.

> If the board is not detected, hold `BOOT` while plugging in USB, then upload.

---

## Hardware pin mapping

| Function                       | Pin    | Port   |
|--------------------------------|--------|--------|
| IR TX primary (shoot)          | GPIO2  | Port A |
| Hit enable/control             | GPIO1  | Port A |
| IR TX secondary (level 5+)     | GPIO8  | Port B |
| IR RX sensor 1                 | GPIO5  | Port C |
| IR RX sensor 2                 | GPIO4  | Port C |
| CAN TX                         | GPIO42 | —      |
| CAN RX                         | GPIO43 | —      |
| I2C SDA (IO expander)          | GPIO13 | —      |
| I2C SCL (IO expander)          | GPIO15 | —      |
| External shoot button (active LOW) | GPIO41 | —  |

Button A is accessed via `M5.BtnA` (M5Unified, on I2C IO expander).
Both BtnA and GPIO41 trigger a shot — connect external button between G41 and GND.

---

## IR protocol

- Protocol: **NEC 32-bit**
- Transmitted value encoding:
  ```
  bits [7:4] = team_id   (4 bits)
  bits [3:0] = player_id (4 bits)
  ```
- Port A (GPIO2) transmits on button press with level-dependent cooldown
- Port B (GPIO8) transmits simultaneously at level 5
- Port C (GPIO5, GPIO4) receives — own `(team << 4) | playerId` filtered out
- NEC repeat codes (`0xFFFFFFFF`) are ignored

---

## Weapon level system

Level advances when the server sends `weapon_upgrade`. Resets to 0 on `die`.

| Level | Shoot cooldown | Second emitter (Port B GPIO8) |
|-------|---------------|-------------------------------|
| 0     | 2000ms        | off                           |
| 1     | 1000ms        | off                           |
| 2     | 500ms         | off                           |
| 3     | 200ms         | off                           |
| 4     | 100ms         | off                           |
| 5     | 100ms         | on                            |

---

## Game modes

```cpp
enum GameType {
    GAME_DEATHMATCH      = 0,
    GAME_TEAM_DEATHMATCH = 1,
};
```

Received via `game_start` / `game_rejoin` fields: `type_of_the_game`, `team_id`, `is_friendlyfire`.  
Team assignment is forwarded to CAN bus via `CAN_GAME_TEAM`. Friendly fire is resolved server-side.

---

## Display feedback

| Event                          | Color    | Text                                     | Duration              |
|--------------------------------|----------|------------------------------------------|-----------------------|
| Connecting to WiFi             | Black    | "Connecting to game server..." (white)   | Until connected       |
| WiFi connected                 | Black    | "Connected!" (white)                     | 1s, clears            |
| MQTT reconnecting              | Black    | "Reconnecting to game server..." (white) | Until reconnected     |
| MQTT reconnected               | Black    | "Reconnected!" (white)                   | 1s, clears            |
| Shot fired                     | Red      | *(none)*                                 | 200ms, clears         |
| Button pressed during cooldown | Green    | *(none)*                                 | 200ms, clears         |
| Hit received                   | Black    | "HITTED" (white)                         | 200ms, clears         |
| Kill confirmed                 | Black    | "KILL!\nKills: N\nLevel: N" (white)      | 1s, clears            |
| Weapon level up                | Magenta  | "LEVELUP N" (white)                      | 3s, clears            |
| Rejoined (alive)               | Black    | "REJOINED\nLvl: N" (white)               | 2s, clears            |
| Rejoined (dead)                | Dark red | "REJOINED\nDEAD\nLvl: N" (white)        | 2s, then stays "DEAD" |
| Player dead                    | Dark red | "DEAD" (white)                           | Stays until respawn   |
| Player respawned               | Blue     | "RESPAWNED" (white)                      | 2s, clears            |
| Game started                   | Cyan     | "GAME START" (black)                     | 1s, clears            |
| Game over                      | Yellow   | "GAME OVER" (black)                      | 3s, clears            |

---

## MQTT

- Subscribes to: `device/{playerId}/command`
- Publishes events to: `device/{playerId}/event`
- Heartbeat every 6s to: `device/{playerId}/heartbeat`
- Reconnects automatically with 2s cooldown
- On connect/reconnect: sends `CAN_SYS_REGISTER` + `CAN_STATUS_ALIVE` on CAN bus

### Incoming commands (`action` field on `command` topic)

| action           | Effect                                                                   |
|------------------|--------------------------------------------------------------------------|
| `game_start`     | Reset game state, set team/mode, CAN `GAME_START` + `GAME_TEAM`         |
| `game_end`       | Clear game/team state, CAN `GAME_END`                                    |
| `game_rejoin`    | Restore kills/level/alive/team, apply cooldown, motor on/off             |
| `kill_confirmed` | Add kill, CAN `KILL` + `KILL_COUNT` + `SCORE`                           |
| `weapon_upgrade` | Level up, update cooldown + second emitter, CAN `SCORE`                 |
| `die`            | Reset level to 0, CAN `DEATH` + motor OFF                              |
| `respawn`        | CAN `RESPAWN` + motor ON                                                 |

`game_start` and `game_rejoin` optional fields:

| Field               | Type   | Description                        |
|---------------------|--------|------------------------------------|
| `game_id`           | string | Game session ID                    |
| `team_id`           | int    | Player's team (–1 = no team)       |
| `is_friendlyfire`   | bool   | Friendly fire enabled              |
| `type_of_the_game`  | int    | 0 = deathmatch, 1 = team DM        |

`game_rejoin` additional fields:

| Field           | Type   | Description                    |
|-----------------|--------|--------------------------------|
| `is_alive`      | bool   | Whether player is currently alive |
| `kills`         | int    | Current kill count             |
| `weapon_level`  | int    | Current weapon level (0–5)     |

### Outgoing events (published to `event` topic)

**Shot fired** — published to `device/{playerId}/event` on every button press:
```json
{ "action": "weapon_shoot", "game_id": "..." }
```

**Hit received** — published to `device/{attackerId}/event` on IR hit:
```json
{
  "game_id": "...",
  "victim_id": 3,
  "attacker_team": 1,
  "team_id": 2
}
```

---

## CAN bus protocol

CAN ID format: `0xABC` — A = category, B = function, C = player ID (0 = broadcast)

### Category 0 — System

| Macro                  | CAN ID  | Data        | Meaning                  |
|------------------------|---------|-------------|--------------------------|
| `CAN_SYS_REGISTER(p)`  | `0x00P` | none        | Player P is online       |
| `CAN_SYS_HEARTBEAT(p)` | `0x01P` | none        | Alive ping               |
| `CAN_SYS_CONFIG(p)`    | `0x02P` | config data | Config for player P      |

### Category 1 — Combat

| Macro                    | CAN ID  | Data (1 byte) | Meaning                         |
|--------------------------|---------|---------------|---------------------------------|
| `CAN_COMBAT_HIT(p)`      | `0x10P` | attacker ID   | Player P received a hit         |
| `CAN_COMBAT_KILL(p)`     | `0x11P` | victim ID     | Player P scored a kill          |
| `CAN_COMBAT_KILL_CNT(p)` | `0x12P` | kill count    | Current kill count for player P |
| `CAN_COMBAT_DEATH(p)`    | `0x13P` | none          | Player P died                   |

### Category 2 — Game control

| Macro                 | CAN ID  | Data          | Meaning                   |
|-----------------------|---------|---------------|---------------------------|
| `CAN_GAME_START`      | `0x200` | none          | Start game (broadcast)    |
| `CAN_GAME_END`        | `0x210` | none          | End game (broadcast)      |
| `CAN_GAME_RESPAWN(p)` | `0x22P` | none          | Respawn player P          |
| `CAN_GAME_TEAM(p)`    | `0x23P` | team_id       | Assign player P to team   |

### Category 3 — Status

| Macro                 | CAN ID  | Data (2 bytes)      | Meaning                   |
|-----------------------|---------|---------------------|---------------------------|
| `CAN_STATUS_SCORE(p)` | `0x30P` | [kills, level]      | Score update for player P |
| `CAN_STATUS_ALIVE(p)` | `0x31P` | [0/1]               | Alive or dead             |
| `CAN_STATUS_AMMO(p)`  | `0x32P` | [ammo count]        | Ammo level                |

### Category 4 — Hardware

| Macro                 | CAN ID  | Data | Meaning                      |
|-----------------------|---------|------|------------------------------|
| `CAN_HW_MOTOR_ON(p)`  | `0x40P` | none | Turn motors ON for device P  |
| `CAN_HW_MOTOR_OFF(p)` | `0x41P` | none | Turn motors OFF for device P |

### Decode helpers

```cpp
CAN_ID_CATEGORY(id)  // bits 8–11
CAN_ID_FUNCTION(id)  // bits 4–7
CAN_ID_PLAYER(id)    // bits 0–3
```
