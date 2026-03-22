/*
 * LaserTag CAN Protocol v2.0
 *
 * Classic CAN, 1M bitrate, 11-bit IDs, max 8 bytes
 *
 * ┌─────────────────────────────────────────────────┐
 * │  CAN ID: 0xABC  (3 hex digits)                  │
 * │    A = Informace (kategorie zprávy)              │
 * │    B = Funkce (typ akce v kategorii)             │
 * │    C = Série (ID zařízení/hráče, 0x0–0xF)       │
 * └─────────────────────────────────────────────────┘
 *
 * Příklad: can1  0x120  [3]  05 00 03
 *          ^^^^  ^^^^^  ^^^  ^^^^^^^^
 *          iface  ID   bytes  data
 *
 *   ID 0x120 → A=1 (combat), B=2 (kill count), C=0 (hráč 0)
 *   Data: [kill_count_hi, kill_count_lo, death_count]
 *
 * ── Kategorie (A) ──
 *   0 = Systém / device management
 *   1 = Combat (boj – hity, killy)
 *   2 = Game control (start, stop, respawn)
 *   3 = Stav (heartbeat, skóre)
 *
 * ── Série (C) ──
 *   0x0 = broadcast / všichni
 *   0x1–0xF = konkrétní hráč/zařízení (1–15)
 */

#pragma once

#include <cstdint>
#include <cstdio>

namespace lt {

// ══════════════════════════════════════════════════
//  CAN ID construction:  make_can_id(info, func, serie)
// ══════════════════════════════════════════════════

/// Build a CAN ID from its 3 components.
///   info  = 0x0–0xF  (kategorie)
///   func  = 0x0–0xF  (funkce)
///   serie = 0x0–0xF  (hráč / zařízení)
constexpr uint32_t make_can_id(uint8_t info, uint8_t func, uint8_t serie) {
    return (static_cast<uint32_t>(info & 0xF) << 8)
         | (static_cast<uint32_t>(func & 0xF) << 4)
         | (static_cast<uint32_t>(serie & 0xF));
}

/// Extract info (A) from CAN ID.
constexpr uint8_t can_info(uint32_t id)  { return (id >> 8) & 0xF; }
/// Extract function (B) from CAN ID.
constexpr uint8_t can_func(uint32_t id)  { return (id >> 4) & 0xF; }
/// Extract serie (C) from CAN ID — player/device ID.
constexpr uint8_t can_serie(uint32_t id) { return id & 0xF; }

// ══════════════════════════════════════════════════
//  Kategorie (A = informace)
// ══════════════════════════════════════════════════

constexpr uint8_t INFO_SYSTEM  = 0x0;   // device management
constexpr uint8_t INFO_COMBAT  = 0x1;   // hity, killy
constexpr uint8_t INFO_GAME    = 0x2;   // game control
constexpr uint8_t INFO_STATUS  = 0x3;   // heartbeat, skóre

// ══════════════════════════════════════════════════
//  Funkce (B) pro jednotlivé kategorie
// ══════════════════════════════════════════════════

// ── System (0x0_x) ──
constexpr uint8_t FN_SYS_REGISTER  = 0x0;  // 0x00C  registrace zařízení
constexpr uint8_t FN_SYS_HEARTBEAT = 0x1;  // 0x01C  heartbeat
constexpr uint8_t FN_SYS_CONFIG    = 0x2;  // 0x02C  konfigurace

// ── Combat (0x1_x) ──
constexpr uint8_t FN_CMB_HIT       = 0x0;  // 0x10C  "byl jsem trefený"
constexpr uint8_t FN_CMB_KILL      = 0x1;  // 0x11C  "dal jsem kill"
constexpr uint8_t FN_CMB_KILLCOUNT = 0x2;  // 0x12C  aktuální počet killů
constexpr uint8_t FN_CMB_DEATH     = 0x3;  // 0x13C  "umřel jsem"

// ── Game control (0x2_x) ──
constexpr uint8_t FN_GAME_START    = 0x0;  // 0x200  start hry (broadcast)
constexpr uint8_t FN_GAME_END      = 0x1;  // 0x210  konec hry (broadcast)
constexpr uint8_t FN_GAME_RESPAWN  = 0x2;  // 0x22C  respawn hráče
constexpr uint8_t FN_GAME_TEAM     = 0x3;  // 0x23C  přiřazení týmu

// ── Status (0x3_x) ──
constexpr uint8_t FN_STS_SCORE     = 0x0;  // 0x30C  skóre hráče
constexpr uint8_t FN_STS_ALIVE     = 0x1;  // 0x31C  alive/dead stav
constexpr uint8_t FN_STS_AMMO      = 0x2;  // 0x32C  munice

// ══════════════════════════════════════════════════
//  Broadcast / all players
// ══════════════════════════════════════════════════

constexpr uint8_t SERIE_BROADCAST = 0x0;

// ══════════════════════════════════════════════════
//  Pre-built CAN IDs for common messages
// ══════════════════════════════════════════════════

/// "Hráč X byl trefený" — gun→gateway   (data: [attacker_id])
constexpr uint32_t id_hit(uint8_t player) {
    return make_can_id(INFO_COMBAT, FN_CMB_HIT, player);
}

/// "Hráč X dal kill" — gun→gateway   (data: [victim_id])
constexpr uint32_t id_kill(uint8_t player) {
    return make_can_id(INFO_COMBAT, FN_CMB_KILL, player);
}

/// "Kill count hráče X" — gateway→gun   (data: [kills_hi, kills_lo])
constexpr uint32_t id_killcount(uint8_t player) {
    return make_can_id(INFO_COMBAT, FN_CMB_KILLCOUNT, player);
}

/// "Hráč X umřel" — gun→gateway
constexpr uint32_t id_death(uint8_t player) {
    return make_can_id(INFO_COMBAT, FN_CMB_DEATH, player);
}

/// Game start broadcast
constexpr uint32_t id_game_start() {
    return make_can_id(INFO_GAME, FN_GAME_START, SERIE_BROADCAST);
}

/// Game end broadcast
constexpr uint32_t id_game_end() {
    return make_can_id(INFO_GAME, FN_GAME_END, SERIE_BROADCAST);
}

// ══════════════════════════════════════════════════
//  Debug: print CAN ID breakdown
// ══════════════════════════════════════════════════

inline void print_can_id(uint32_t id) {
    std::printf("CAN ID 0x%03X → info=%X func=%X serie=%X\n",
                id, can_info(id), can_func(id), can_serie(id));
}

} // namespace lt
