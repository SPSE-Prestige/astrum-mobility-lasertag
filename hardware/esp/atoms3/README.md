# ESP Firmware — Atom S3 (ESP32-S3)

Laser tag firmware for M5Stack Atom S3.
WiFi + MQTT + CAN bus. IR TX and RX both on Port A.

---

## Project structure

```
atoms3/
├── atoms3.ino              ← main sketch
├── pins.h                  ← hardware pin config
├── credentials.h           ← WiFi & MQTT credentials
├── README.md
├── src/
│   ├── ir_controller.h/cpp ← combined TX+RX on Port A
│   ├── can_bus.h/cpp
│   ├── can_protocol.h
│   ├── can_debug.h
│   ├── mqtt_client.h/cpp
│   └── game_state.h
└── setupid/
    └── setupid.ino         ← flash player ID utility
```

---

## First-time setup — save player ID to flash

1. Open `setupid/setupid.ino`
2. Set `PLAYER_IDD` to the device's unique number (1–15)
3. Flash to device — saves ID to NVS flash
4. Only needed **once per device**

Then flash `atoms3.ino`.

---

## Hardware pin mapping

| Function       | Pin    | Port   |
|----------------|--------|--------|
| IR TX (shoot)  | GPIO2  | Port A |
| IR RX (hit)    | GPIO1  | Port A |
| CAN TX         | GPIO5  | —      |
| CAN RX         | GPIO6  | —      |

Button A is accessed via `M5.BtnA` (M5Unified).

---

## How to flash

**Board:** `Tools → Board → ESP32S3 Dev Module`

Settings:
- Flash Size: `8MB`
- USB CDC On Boot: `Enabled`

**Libraries required:**
- `M5Unified`
- `IRremoteESP8266`
- `PubSubClient`
- `ArduinoJson`

---

## Display feedback

| Event                          | Color                  | Text              | Duration     |
|--------------------------------|------------------------|-------------------|--------------|
| Shot fired                     | Red                    | *(none)*          | 200ms, clears |
| Button pressed during cooldown | Green                  | *(none)*          | 200ms, clears |
| Hit received                   | Black                  | "HITTED" (white)  | 1s, clears   |
| Player dead                    | Dark red               | "DEAD" (white)    | Stays until respawn |
| Player respawned               | Blue                   | "RESPAWNED" (white) | 1s, clears |
| Game started                   | Cyan                   | "GAME START" (black) | 1s, clears |
| Game over                      | Yellow                 | "GAME OVER" (black) | Stays       |

---

## CAN messages sent

| Trigger      | CAN ID         | Data        |
|--------------|----------------|-------------|
| Boot         | `CAN_SYS_REGISTER(p)` `0x00P` | none |
| Hit received | `CAN_COMBAT_HIT(p)` `0x10P`   | attacker ID |
