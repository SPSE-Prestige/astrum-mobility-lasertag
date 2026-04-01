#pragma once
#include "can_bus.h"

// CAN ID format: 0xABC
//   A = category   (0=System, 1=Combat, 2=Game, 3=Status, 4=Hardware)
//   B = function
//   C = player/device ID (0=broadcast, 1-F=player)

// ─── Category 0: System ───────────────────────────────────────────────────────
#define CAN_SYS_REGISTER(p)     (0x000 | ((p) & 0xF))  // player P online
#define CAN_SYS_HEARTBEAT(p)    (0x010 | ((p) & 0xF))  // alive ping
#define CAN_SYS_CONFIG(p)       (0x020 | ((p) & 0xF))  // config data

// ─── Category 1: Combat ───────────────────────────────────────────────────────
#define CAN_COMBAT_HIT(p)       (0x100 | ((p) & 0xF))  // player P received a hit
#define CAN_COMBAT_KILL(p)      (0x110 | ((p) & 0xF))  // player P scored a kill
#define CAN_COMBAT_KILL_CNT(p)  (0x120 | ((p) & 0xF))  // kill count for player P
#define CAN_COMBAT_DEATH(p)     (0x130 | ((p) & 0xF))  // player P died

// ─── Category 2: Game control ────────────────────────────────────────────────
#define CAN_GAME_START          (0x200)                 // start game broadcast
#define CAN_GAME_END            (0x210)                 // end game broadcast
#define CAN_GAME_RESPAWN(p)     (0x220 | ((p) & 0xF))  // respawn player P
#define CAN_GAME_TEAM(p)        (0x230 | ((p) & 0xF))  // assign team for player P

// ─── Category 3: Status ───────────────────────────────────────────────────────
#define CAN_STATUS_SCORE(p)     (0x300 | ((p) & 0xF))  // score update
#define CAN_STATUS_ALIVE(p)     (0x310 | ((p) & 0xF))  // alive/dead state
#define CAN_STATUS_AMMO(p)      (0x320 | ((p) & 0xF))  // ammo level

// ─── Category 4: Hardware ─────────────────────────────────────────────────────
#define CAN_HW_MOTOR_ON(p)      (0x400 | ((p) & 0xF))  // turn motors ON  for device P
#define CAN_HW_MOTOR_OFF(p)     (0x410 | ((p) & 0xF))  // turn motors OFF for device P

// ─── Helpers ──────────────────────────────────────────────────────────────────

// Decode ID fields
#define CAN_ID_CATEGORY(id)     (((id) >> 8) & 0xF)
#define CAN_ID_FUNCTION(id)     (((id) >> 4) & 0xF)
#define CAN_ID_PLAYER(id)       ((id) & 0xF)

// ─── Send helpers ─────────────────────────────────────────────────────────────

// Send player registration
// can.send(CAN_SYS_REGISTER(playerId), nullptr, 0);

// Send hit received (1 data byte = attacker ID)
// uint8_t attacker = attackerId;
// can.send(CAN_COMBAT_HIT(playerId), &attacker, 1);

// Send kill scored (1 data byte = victim ID)
// uint8_t victim = victimId;
// can.send(CAN_COMBAT_KILL(playerId), &victim, 1);

// Send motor ON for device P
// can.send(CAN_HW_MOTOR_ON(deviceId), nullptr, 0);

// Send motor OFF for device P
// can.send(CAN_HW_MOTOR_OFF(deviceId), nullptr, 0);
