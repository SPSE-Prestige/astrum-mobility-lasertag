# ESP Firmware — STAMPLC (StampS3A, ESP32-S3FN8)

This directory contains the firmware for the M5Stack STAMPLC unit.
Responsibilities: IR shooting, IR hit detection, CAN bus communication, MQTT, display feedback.

---

## Project structure

```
stamplc/
├── stamplc.ino        ← main sketch
├── pins.h             ← hardware pin config
├── credentials.h      ← WiFi & MQTT credentials
├── README.md
├── src/               ← all library code
│   ├── can_bus.h/cpp
│   ├── can_protocol.h
│   ├── can_debug.h
│   ├── ir_controller.h/cpp
│   └── mqtt_client.h/cpp
├── test/
│   └── test.ino       ← AUnit tests
└── setupid/
    └── setupid.ino    ← flash player ID utility
```

---

## First-time setup — save player ID to flash

Before flashing the main firmware, you must write the player ID into the ESP's flash memory. Each device needs a unique ID (1–15).

1. Open `setupid/setupid.ino` in Arduino IDE
2. Change `PLAYER_IDD` to the desired player number:
   ```cpp
   const int PLAYER_IDD = 1; // change this per device
   ```
3. Flash it to the device — it will save the ID and print it over Serial
4. You only need to do this **once per device**. The ID survives power cycles and firmware updates.

After the ID is saved, flash the main `stamplc.ino` firmware.

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

**Requirements:** Arduino IDE 2.x — https://www.arduino.cc/en/software

**Step 1 — Add ESP32 board support**

Go to `File → Preferences` and add to "Additional boards manager URLs":
```
https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json
```
Then `Tools → Board → Boards Manager` → search `esp32` → install **esp32 by Espressif**.

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
- Partition Scheme: `Default 4MB with spiffs` (or `8M Flash`)
- USB CDC On Boot: `Enabled`

**Step 4 — Open the sketch**

`File → Open` → select `stamplc.ino`

**Step 5 — Flash**

Connect STAMPLC via USB-C, select port under `Tools → Port`, click **Upload**.

> If the board is not detected, hold the `BOOT` button while plugging in USB, then upload.

---

## Hardware pin mapping

| Function              | Pin     | Port   |
|-----------------------|---------|--------|
| IR TX (shoot)         | GPIO2   | Port A |
| Hit enable/control    | GPIO1   | Port A |
| IR RX sensor 1        | GPIO5   | Port C |
| IR RX sensor 2        | GPIO4   | Port C |
| CAN TX                | GPIO42  | —      |
| CAN RX                | GPIO43  | —      |
| I2C SDA (expander)    | GPIO13  | —      |
| I2C SCL (expander)    | GPIO15  | —      |

Button A is connected through an I2C IO expander — access via `M5.BtnA`.

---

## Display feedback

| Event                          | Color       | Text                         | Duration            |
|--------------------------------|-------------|------------------------------|---------------------|
| Connecting to WiFi             | Black       | "Connecting to game server..." (white) | Until connected |
| WiFi connected                 | Black       | "Connected!" (white)         | 1s, clears          |
| Shot fired                     | Red         | *(none)*                     | 200ms, clears       |
| Button pressed during cooldown | Green       | *(none)*                     | 200ms, clears       |
| Hit received                   | Black       | "HITTED" (white)             | 1s, clears          |
| Player dead                    | Dark red    | "DEAD" (white)               | Stays until respawn |
| Player respawned               | Blue        | "RESPAWNED" (white)          | 1s, clears          |
| Game started                   | Cyan        | "GAME START" (black)         | 1s, clears          |
| Game over                      | Yellow      | "GAME OVER" (black)          | Stays               |
| MQTT reconnecting              | Black       | "Reconnecting to game server..." (white) | Until reconnected |
| MQTT reconnected               | Black       | "Reconnected!" (white)       | 1s, clears          |

---

## IR protocol

- Protocol: **NEC 32-bit**
- Sent value: player ID (`uint32_t`)
- Port A (GPIO2) transmits on button press, **2-second cooldown** between shots
- Port C (GPIO5, GPIO4) receives — own player ID is ignored

---

## MQTT

- Connects to broker on boot, reconnects automatically with **2-second cooldown** between attempts
- On connect/reconnect: sends `CAN_SYS_REGISTER` and `CAN_STATUS_ALIVE` on CAN bus
- Subscribes to: `device/{playerId}/command`
- Publishes hit events to: `device/{attackerId}/event`
- Heartbeat every 6 seconds to: `device/{playerId}/heartbeat`

### Supported MQTT `action` commands

| action             | Effect                                      |
|--------------------|---------------------------------------------|
| `game_start`       | CAN `GAME_START`, stores `game_id`          |
| `game_end`         | CAN `GAME_END`, clears `game_id`            |
| `kill_confirmed`   | CAN `KILL`, `SCORE`, `KILL_COUNT`           |
| `die`              | CAN `DEATH`                                 |
| `respawn`          | CAN `RESPAWN`                               |

---

## CAN bus protocol

CAN ID format: `0xABC`

```
A = category
B = function
C = player / device ID  (0 = broadcast, 1–F = player)
```

### Categories

| A | Name     |
|---|----------|
| 0 | System   |
| 1 | Combat   |
| 2 | Game     |
| 3 | Status   |
| 4 | Hardware |

### Category 0 — System

| Macro                  | CAN ID  | Data        | Meaning                         |
|------------------------|---------|-------------|---------------------------------|
| `CAN_SYS_REGISTER(p)`  | `0x00P` | none        | Player P is online              |
| `CAN_SYS_HEARTBEAT(p)` | `0x01P` | none        | Alive ping from player P        |
| `CAN_SYS_CONFIG(p)`    | `0x02P` | config data | Configuration data for player P |

### Category 1 — Combat

| Macro                    | CAN ID  | Data (1 byte) | Meaning                         |
|--------------------------|---------|---------------|---------------------------------|
| `CAN_COMBAT_HIT(p)`      | `0x10P` | attacker ID   | Player P received a hit         |
| `CAN_COMBAT_KILL(p)`     | `0x11P` | victim ID     | Player P scored a kill          |
| `CAN_COMBAT_KILL_CNT(p)` | `0x12P` | kill count    | Current kill count for player P |
| `CAN_COMBAT_DEATH(p)`    | `0x13P` | none          | Player P died                   |

### Category 2 — Game control

| Macro                 | CAN ID  | Data | Meaning                   |
|-----------------------|---------|------|---------------------------|
| `CAN_GAME_START`      | `0x200` | none | Start game (broadcast)    |
| `CAN_GAME_END`        | `0x210` | none | End game (broadcast)      |
| `CAN_GAME_RESPAWN(p)` | `0x22P` | none | Respawn player P          |
| `CAN_GAME_TEAM(p)`    | `0x23P` | team | Assign player P to a team |

### Category 3 — Status

| Macro                 | CAN ID  | Data       | Meaning                   |
|-----------------------|---------|------------|---------------------------|
| `CAN_STATUS_SCORE(p)` | `0x30P` | score      | Score update for player P |
| `CAN_STATUS_ALIVE(p)` | `0x31P` | 0/1        | Alive or dead state       |
| `CAN_STATUS_AMMO(p)`  | `0x32P` | ammo count | Ammo level for player P   |

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

### Usage examples

```cpp
// Register player on boot / reconnect
can.send(CAN_SYS_REGISTER(playerId), nullptr, 0);

// Report hit received from attacker
uint8_t attacker = attackerId;
can.send(CAN_COMBAT_HIT(playerId), &attacker, 1);

// Report kill scored
uint8_t victim = victimId;
can.send(CAN_COMBAT_KILL(playerId), &victim, 1);

// Turn motors on / off
can.send(CAN_HW_MOTOR_ON(deviceId), nullptr, 0);
can.send(CAN_HW_MOTOR_OFF(deviceId), nullptr, 0);
```
