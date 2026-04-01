// Unit tests for STAMPLC firmware logic
// Framework: AUnit  (install via Library Manager: "AUnit" by Brian T. Park)
//
// Flash to STAMPLC and open Serial Monitor at 115200.
// Each test prints PASSED or FAILED.

#include <AUnit.h>
#include "../src/can_protocol.h"
#include "../src/can_debug.h"

// ─── CAN ID macro tests ───────────────────────────────────────────────────────

test(SysRegister) {
    assertEqual((uint32_t)CAN_SYS_REGISTER(3),  (uint32_t)0x003);
    assertEqual((uint32_t)CAN_SYS_REGISTER(15), (uint32_t)0x00F);
    assertEqual((uint32_t)CAN_SYS_REGISTER(0),  (uint32_t)0x000);
}

test(SysHeartbeat) {
    assertEqual((uint32_t)CAN_SYS_HEARTBEAT(1), (uint32_t)0x011);
    assertEqual((uint32_t)CAN_SYS_HEARTBEAT(5), (uint32_t)0x015);
}

test(SysConfig) {
    assertEqual((uint32_t)CAN_SYS_CONFIG(2), (uint32_t)0x022);
}

test(CombatHit) {
    assertEqual((uint32_t)CAN_COMBAT_HIT(3),  (uint32_t)0x103);
    assertEqual((uint32_t)CAN_COMBAT_HIT(15), (uint32_t)0x10F);
}

test(CombatKill) {
    assertEqual((uint32_t)CAN_COMBAT_KILL(4), (uint32_t)0x114);
}

test(CombatKillCount) {
    assertEqual((uint32_t)CAN_COMBAT_KILL_CNT(7), (uint32_t)0x127);
}

test(CombatDeath) {
    assertEqual((uint32_t)CAN_COMBAT_DEATH(2), (uint32_t)0x132);
}

test(GameStart) {
    assertEqual((uint32_t)CAN_GAME_START, (uint32_t)0x200);
}

test(GameEnd) {
    assertEqual((uint32_t)CAN_GAME_END, (uint32_t)0x210);
}

test(GameRespawn) {
    assertEqual((uint32_t)CAN_GAME_RESPAWN(6), (uint32_t)0x226);
}

test(GameTeam) {
    assertEqual((uint32_t)CAN_GAME_TEAM(1), (uint32_t)0x231);
}

test(StatusScore) {
    assertEqual((uint32_t)CAN_STATUS_SCORE(3), (uint32_t)0x303);
}

test(StatusAlive) {
    assertEqual((uint32_t)CAN_STATUS_ALIVE(2), (uint32_t)0x312);
}

test(StatusAmmo) {
    assertEqual((uint32_t)CAN_STATUS_AMMO(9), (uint32_t)0x329);
}

test(HwMotorOn) {
    assertEqual((uint32_t)CAN_HW_MOTOR_ON(1),  (uint32_t)0x401);
    assertEqual((uint32_t)CAN_HW_MOTOR_ON(15), (uint32_t)0x40F);
    assertEqual((uint32_t)CAN_HW_MOTOR_ON(0),  (uint32_t)0x400);
}

test(HwMotorOff) {
    assertEqual((uint32_t)CAN_HW_MOTOR_OFF(1),  (uint32_t)0x411);
    assertEqual((uint32_t)CAN_HW_MOTOR_OFF(15), (uint32_t)0x41F);
    assertEqual((uint32_t)CAN_HW_MOTOR_OFF(0),  (uint32_t)0x410);
}

// ─── Decode helper tests ──────────────────────────────────────────────────────

test(DecodeCategory) {
    assertEqual((uint32_t)CAN_ID_CATEGORY(0x103), (uint32_t)1);
    assertEqual((uint32_t)CAN_ID_CATEGORY(0x200), (uint32_t)2);
    assertEqual((uint32_t)CAN_ID_CATEGORY(0x400), (uint32_t)4);
    assertEqual((uint32_t)CAN_ID_CATEGORY(0x000), (uint32_t)0);
}

test(DecodeFunction) {
    assertEqual((uint32_t)CAN_ID_FUNCTION(0x103), (uint32_t)0);
    assertEqual((uint32_t)CAN_ID_FUNCTION(0x113), (uint32_t)1);
    assertEqual((uint32_t)CAN_ID_FUNCTION(0x413), (uint32_t)1);
}

test(DecodePlayer) {
    assertEqual((uint32_t)CAN_ID_PLAYER(0x103), (uint32_t)3);
    assertEqual((uint32_t)CAN_ID_PLAYER(0x10F), (uint32_t)15);
    assertEqual((uint32_t)CAN_ID_PLAYER(0x200), (uint32_t)0);
}

// ─── Label tests ──────────────────────────────────────────────────────────────

test(LabelSysRegister)   { assertEqual(strcmp(canIdLabel(0x003), "SYS/REGISTER"),    0); }
test(LabelSysHeartbeat)  { assertEqual(strcmp(canIdLabel(0x010), "SYS/HEARTBEAT"),   0); }
test(LabelSysConfig)     { assertEqual(strcmp(canIdLabel(0x020), "SYS/CONFIG"),      0); }
test(LabelCombatHit)     { assertEqual(strcmp(canIdLabel(0x103), "COMBAT/HIT"),      0); }
test(LabelCombatKill)    { assertEqual(strcmp(canIdLabel(0x114), "COMBAT/KILL"),     0); }
test(LabelCombatKillCnt) { assertEqual(strcmp(canIdLabel(0x120), "COMBAT/KILL_COUNT"), 0); }
test(LabelCombatDeath)   { assertEqual(strcmp(canIdLabel(0x130), "COMBAT/DEATH"),    0); }
test(LabelGameStart)     { assertEqual(strcmp(canIdLabel(0x200), "GAME/START"),      0); }
test(LabelGameEnd)       { assertEqual(strcmp(canIdLabel(0x210), "GAME/END"),        0); }
test(LabelGameRespawn)   { assertEqual(strcmp(canIdLabel(0x220), "GAME/RESPAWN"),    0); }
test(LabelGameTeam)      { assertEqual(strcmp(canIdLabel(0x230), "GAME/TEAM"),       0); }
test(LabelStatusScore)   { assertEqual(strcmp(canIdLabel(0x300), "STATUS/SCORE"),    0); }
test(LabelStatusAlive)   { assertEqual(strcmp(canIdLabel(0x310), "STATUS/ALIVE"),    0); }
test(LabelStatusAmmo)    { assertEqual(strcmp(canIdLabel(0x320), "STATUS/AMMO"),     0); }
test(LabelHwMotorOn)     { assertEqual(strcmp(canIdLabel(0x401), "HW/MOTOR_ON"),     0); }
test(LabelHwMotorOff)    { assertEqual(strcmp(canIdLabel(0x410), "HW/MOTOR_OFF"),    0); }
test(LabelUnknown)       { assertEqual(strcmp(canIdLabel(0xFFF), "UNKNOWN"),         0); }

// ─── Player ID mask tests ─────────────────────────────────────────────────────

test(PlayerIdMaskedTo4Bits) {
    assertEqual((uint32_t)CAN_COMBAT_HIT(0x1F),   (uint32_t)0x10F);
    assertEqual((uint32_t)CAN_HW_MOTOR_ON(0xFF),  (uint32_t)0x40F);
    assertEqual((uint32_t)CAN_HW_MOTOR_OFF(0xFF), (uint32_t)0x41F);
    assertEqual((uint32_t)CAN_SYS_REGISTER(0x10), (uint32_t)0x000); // 0x10 & 0xF = 0
}

// ─── Shoot cooldown logic tests ───────────────────────────────────────────────
// Tests the 2-second cooldown without hardware by simulating timestamps.

static bool     fakeButton   = false;
static bool     shootFired   = false;
static bool     cooldownFired = false;
static unsigned long fakeNow = 0;
static unsigned long lastShot = 0;

static bool runCooldownLogic() {
    if (!fakeButton) return false;
    if (fakeNow - lastShot < 2000) {
        cooldownFired = true;
        return false;
    }
    lastShot = fakeNow;
    shootFired = true;
    return true;
}

test(ShootFirstPressAlwaysFires) {
    fakeButton = true; fakeNow = 0; lastShot = 0;
    shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertTrue(shootFired);
    assertFalse(cooldownFired);
}

test(ShootBlockedWithin2Seconds) {
    lastShot = 0; fakeNow = 1000;  // 1s after last shot
    fakeButton = true; shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertFalse(shootFired);
    assertTrue(cooldownFired);
}

test(ShootAllowedAfter2Seconds) {
    lastShot = 0; fakeNow = 2000;  // exactly 2s
    fakeButton = true; shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertTrue(shootFired);
    assertFalse(cooldownFired);
}

test(ShootNotFiredWhenButtonNotPressed) {
    fakeButton = false; fakeNow = 5000; lastShot = 0;
    shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertFalse(shootFired);
    assertFalse(cooldownFired);
}

test(ShootAllowedAfterCooldownExpires) {
    // First shot at t=0
    fakeButton = true; fakeNow = 0; lastShot = 0;
    shootFired = false; runCooldownLogic();
    assertTrue(shootFired);
    // Second attempt at t=1500 — blocked
    fakeNow = 1500; shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertFalse(shootFired);
    // Third attempt at t=2001 — allowed
    fakeNow = 2001; shootFired = false; cooldownFired = false;
    runCooldownLogic();
    assertTrue(shootFired);
}

// ─── Runner ───────────────────────────────────────────────────────────────────

void setup() {
    Serial.begin(115200);
    delay(2000);
}

void loop() {
    aunit::TestRunner::run();
}
