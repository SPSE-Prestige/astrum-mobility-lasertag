#pragma once
#include "../pins.h"
#include "can_protocol.h"
#include <Arduino.h>

#if DEBUG

// Decode a CAN ID into a human-readable label
static const char* canIdLabel(uint32_t id) {
    switch (id & 0xFF0) {
        // System
        case 0x000: return "SYS/REGISTER";
        case 0x010: return "SYS/HEARTBEAT";
        case 0x020: return "SYS/CONFIG";
        // Combat
        case 0x100: return "COMBAT/HIT";
        case 0x110: return "COMBAT/KILL";
        case 0x120: return "COMBAT/KILL_COUNT";
        case 0x130: return "COMBAT/DEATH";
        // Game
        case 0x200: return "GAME/START";
        case 0x210: return "GAME/END";
        case 0x220: return "GAME/RESPAWN";
        case 0x230: return "GAME/TEAM";
        // Status
        case 0x300: return "STATUS/SCORE";
        case 0x310: return "STATUS/ALIVE";
        case 0x320: return "STATUS/AMMO";
        // Hardware
        case 0x400: return "HW/MOTOR_ON";
        case 0x410: return "HW/MOTOR_OFF";
        default:    return "UNKNOWN";
    }
}

static void canDebugTx(uint32_t id, const uint8_t* data, uint8_t len) {
    Serial.print("[CAN TX] ID=0x");
    Serial.print(id, HEX);
    Serial.print("  label=");
    Serial.print(canIdLabel(id));
    Serial.print("  player=");
    Serial.print(CAN_ID_PLAYER(id));
    if (len > 0 && data) {
        Serial.print("  data=[ ");
        for (uint8_t i = 0; i < len; i++) {
            Serial.print("0x");
            Serial.print(data[i], HEX);
            Serial.print(" ");
        }
        Serial.print("]");
    }
    Serial.println();
}

static void canDebugRx(uint32_t id, const uint8_t* data, uint8_t len) {
    Serial.print("[CAN RX] ID=0x");
    Serial.print(id, HEX);
    Serial.print("  label=");
    Serial.print(canIdLabel(id));
    Serial.print("  player=");
    Serial.print(CAN_ID_PLAYER(id));
    if (len > 0 && data) {
        Serial.print("  data=[ ");
        for (uint8_t i = 0; i < len; i++) {
            Serial.print("0x");
            Serial.print(data[i], HEX);
            Serial.print(" ");
        }
        Serial.print("]");
    }
    Serial.println();
}

#define CAN_DEBUG_TX(id, data, len) canDebugTx(id, data, len)
#define CAN_DEBUG_RX(id, data, len) canDebugRx(id, data, len)

#else

#define CAN_DEBUG_TX(id, data, len)
#define CAN_DEBUG_RX(id, data, len)

#endif
