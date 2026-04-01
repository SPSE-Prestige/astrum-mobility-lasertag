# ESP Latest Firmware

This directory contains the current ESP32 firmware for the latest hardware revision.
The firmware is responsible for: Wi-Fi setup, IR shooting, and CAN bus registration.

## What `latest.ino` does

- Starts serial debug output.
- Reads `playerId` from non-volatile preferences (`Preferences`).
- Connects to Wi-Fi as a station.
- Initializes IR and CAN hardware.
- Registers the player on the CAN bus using ID `0x00P`.
- Sends IR when the button is pressed.
- Receives IR codes from other devices and forwards them to MQTT.

## IR codes sent by ESP

The ESP uses the NEC IR protocol.

### Meaning of IR codes

- `ir.send(playerId)` sends the current player ID as a 32-bit NEC code.
- The value is the shooter/player identifier.
- When another device receives this code, it interprets it as "hit by player X".
- The ESP ignores its own IR code on receive.

### What the code means

- If your `playerId` is `5`, the ESP sends `0x00000005` over IR.
- The receiver prints `Received IR code: 0x5` and may use that as the attacker ID.

## CAN bus codes sent by ESP to MP135

The ESP `CanBus` implementation uses standard 11-bit CAN identifiers.
The ID format is:

```
0xABC
```

Where:

- `A` = information category
- `B` = function within that category
- `C` = device / player series ID

### Categories (`A`)

- `0` = System / device management
- `1` = Combat (hits, kills)
- `2` = Game control (start/end/respawn)
- `3` = Status (heartbeat, score, ammo)

### Series (`C`)

- `0x0` = broadcast / all devices
- `0x1`–`0xF` = player or device ID 1–15

## ESP messages on CAN bus

### Registration

- `CAN ID = 0x00P`
- `P` = player ID low nibble (`playerId & 0xF`)
- No data bytes

This is the only CAN message sent by the current `esp/latest` firmware.
It is a registration packet telling the network "player P is online".

### Example

- `playerId == 3` → CAN ID `0x003`
- `playerId == 12` → CAN ID `0x00C`

## Full CAN code meanings

| CAN ID | Category | Function | Target | Meaning |
|--------|----------|----------|--------|---------|
| `0x00P` | System | Register | player P | register player P on the bus |
| `0x01P` | System | Heartbeat | player P | periodic alive ping |
| `0x02P` | System | Config | player P | configuration data for player P |
| `0x10P` | Combat | Hit | player P | player P reports a hit received |
| `0x11P` | Combat | Kill | player P | player P reports it scored a kill |
| `0x12P` | Combat | Kill count | player P | current kill count for player P |
| `0x13P` | Combat | Death | player P | player P reports it died |
| `0x20P` | Game | Start | broadcast | start game broadcast |
| `0x21P` | Game | End | broadcast | end game broadcast |
| `0x22P` | Game | Respawn | player P | respawn request for player P |
| `0x23P` | Game | Team | player P | assign player P to team |
| `0x30P` | Status | Score | player P | score update for player P |
| `0x31P` | Status | Alive | player P | alive/dead state for player P |
| `0x32P` | Status | Ammo | player P | ammo level for player P |

## How MP135 should read the codes

The MP135-side CAN decoder should split the ID into three parts:

- `info = (id >> 8) & 0xF`
- `func = (id >> 4) & 0xF`
- `serie = id & 0xF`

That gives you the category, function, and player/device target.

## Notes

- The current ESP firmware sends only registration on CAN bus.
- The IR layer sends full player ID values.
- If you need more CAN messages from the ESP, add them in `firmware/esp/latest/can_bus.cpp` and `latest.ino`.
