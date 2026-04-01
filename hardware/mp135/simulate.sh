#!/bin/bash
# ══════════════════════════════════════════════════════════════
#  LaserTag MP135 — Full Game Simulation
#
#  Sends CAN frames via cansend to simulate a complete game.
#  Uses CAN loopback mode so frames are echoed back without
#  needing a physical device on the bus.
#
#  Usage: bash simulate.sh [can_iface] [player_id]
#         bash simulate.sh can1 3
#
#  CAN ID: 0xABC  (A=category, B=function, C=player)
#    A: 0=sys, 1=combat, 2=game, 3=status
#    B: see can_protocol.h
#    C: player/device ID (1-F, 0=broadcast)
# ══════════════════════════════════════════════════════════════

set -euo pipefail

CAN="${1:-can1}"
PLAYER="${2:-3}"
ENEMY1="5"
ENEMY2="7"
TEAMMATE="2"

hex() { printf '%01X' "$1"; }

P=$(hex "$PLAYER")
E1=$(hex "$ENEMY1")
E2=$(hex "$ENEMY2")
T=$(hex "$TEAMMATE")

LASERTAG_PID=""

log() {
    echo ""
    echo "═══════════════════════════════════════"
    echo "  $1"
    echo "═══════════════════════════════════════"
}

send() {
    local id="$1"
    local data="${2:-}"
    if [ -n "$data" ]; then
        cansend "$CAN" "${id}#${data}"
    else
        cansend "$CAN" "${id}#"
    fi
    echo "  → 0x${id} ${data:+[${data}]}"
    sleep 0.05  # small delay to prevent TX buffer overflow
}

pause() {
    echo "  ⏳ ${1}s..."
    sleep "$1"
}

cleanup() {
    echo ""
    echo "Cleaning up..."
    if [ -n "$LASERTAG_PID" ] && kill -0 "$LASERTAG_PID" 2>/dev/null; then
        kill "$LASERTAG_PID" 2>/dev/null
        wait "$LASERTAG_PID" 2>/dev/null || true
    fi
    # Restore normal CAN config (loopback off)
    ip link set "$CAN" down 2>/dev/null || true
    ip link set "$CAN" up type can bitrate 1000000 restart-ms 100 loopback off 2>/dev/null || true
    # Restart lasertag service
    systemctl start lasertag 2>/dev/null || true
    echo "Done. Lasertag service restarted."
}

trap cleanup EXIT

echo "╔══════════════════════════════════════════════╗"
echo "║     LASERTAG FULL GAME SIMULATION           ║"
echo "╠══════════════════════════════════════════════╣"
echo "║  CAN:      $CAN                             "
echo "║  Player:   $PLAYER (0x${P})                  "
echo "║  Enemies:  $ENEMY1 (0x${E1}), $ENEMY2 (0x${E2})"
echo "║  Teammate: $TEAMMATE (0x${T})                "
echo "╚══════════════════════════════════════════════╝"

# ── Setup: stop service, enable CAN loopback, start binary ──
echo ""
echo "Setting up CAN loopback mode for simulation..."
systemctl stop lasertag 2>/dev/null || true
sleep 1

ip link set "$CAN" down 2>/dev/null || true
ip link set "$CAN" up type can bitrate 1000000 restart-ms 100 loopback on
echo "  CAN $CAN: loopback ON"

echo "Starting lasertag binary..."
cd /lasertag
./lasertag -c "$CAN" -f /dev/fb0 -m media &
LASERTAG_PID=$!
sleep 2

echo "  lasertag PID: $LASERTAG_PID"
echo ""
echo "Press ENTER to start game simulation..."
read -r

# ══════════════════════════════════════════════════
#  PHASE 1: Device Registration
#    CAN ID: 0x00C = SYS(0) + REGISTER(0) + player(C)
# ══════════════════════════════════════════════════
log "PHASE 1: DEVICE REGISTRATION"

echo "  Registering local player $PLAYER..."
send "00${P}"
pause 3

echo "  Registering other players..."
send "00${E1}"
pause 0.5
send "00${E2}"
pause 0.5
send "00${T}"
pause 2

echo "  Assigning teams..."
# CAN ID: 0x23C = GAME(2) + TEAM(3) + player(C)
send "23${P}" "01"                # player → team 1
pause 0.3
send "23${T}" "01"                # teammate → team 1
pause 0.3
send "23${E1}" "02"               # enemy1 → team 2
pause 0.3
send "23${E2}" "02"               # enemy2 → team 2
pause 1

echo "  Sending heartbeats..."
# CAN ID: 0x01C = SYS(0) + HEARTBEAT(1) + player(C)
send "01${P}"
pause 0.1
send "01${E1}"
pause 0.1
send "01${E2}"
pause 0.1
send "01${T}"
pause 3

# ══════════════════════════════════════════════════
#  PHASE 2: Game Start
#    CAN ID: 0x200 = GAME(2) + START(0) + broadcast(0)
# ══════════════════════════════════════════════════
log "PHASE 2: GAME START (countdown 5..4..3..2..1..GO!)"

send "200"
echo "  Countdown animation playing..."
pause 8

# ══════════════════════════════════════════════════
#  PHASE 3: Combat
# ══════════════════════════════════════════════════
log "PHASE 3: COMBAT"

# -- Initial ammo: 0x32C = STATUS(3) + AMMO(2) + player --
echo "  Setting initial ammo (30)..."
send "32${P}" "001E"
pause 1.5

# -- Enemy hits us: 0x10C = COMBAT(1) + HIT(0) + player --
echo ""
echo "  >>> Enemy $ENEMY1 hits us!"
send "10${P}" "0${E1}"
pause 3

# -- Ammo decreasing --
echo "  >>> Shooting... ammo decreasing"
send "32${P}" "0019"              # 25
pause 0.8
send "32${P}" "0014"              # 20
pause 0.8
send "32${P}" "000F"              # 15
pause 1.5

# -- We hit enemy: 0x10C = COMBAT(1) + HIT(0) + enemy --
echo "  >>> We hit enemy $ENEMY1!"
send "10${E1}" "0${P}"
pause 2

# -- Score update: 0x30C = STATUS(3) + SCORE(0) + player --
echo "  >>> Score: +100"
send "30${P}" "0064"
pause 2

# -- Another hit on us --
echo "  >>> Enemy $ENEMY2 hits us!"
send "10${P}" "0${E2}"
pause 3

# -- Kill! --
# 0x11C = COMBAT(1) + KILL(1) + killer
# 0x13C = COMBAT(1) + DEATH(3) + victim
# 0x12C = COMBAT(1) + KILLCOUNT(2) + player
echo "  >>> WE KILLED enemy $ENEMY1!"
send "11${P}" "0${E1}"
pause 1
send "13${E1}"
pause 0.5
send "12${P}" "0001"
pause 0.5
send "30${P}" "00C8"              # Score=200
pause 4

# -- Low ammo --
echo "  >>> Ammo getting low..."
send "32${P}" "000A"              # 10
pause 1.5
send "32${P}" "0005"              # 5 → LOW AMMO BLINK!
pause 3
send "32${P}" "0002"              # 2
pause 2

# -- Empty --
echo "  >>> AMMO EMPTY!"
send "32${P}" "0000"              # 0
pause 3

# -- Refill --
echo "  >>> Ammo refilled!"
send "32${P}" "001E"              # 30
pause 2

# ══════════════════════════════════════════════════
#  PHASE 4: Death & Respawn
# ══════════════════════════════════════════════════
log "PHASE 4: DEATH & RESPAWN"

echo "  >>> Hit..."
send "10${P}" "0${E2}"
pause 1.5

echo "  >>> DEAD!"
# 0x13C = COMBAT(1) + DEATH(3) + player
send "13${P}"
pause 6

echo "  >>> Respawning..."
# 0x22C = GAME(2) + RESPAWN(2) + player
send "22${P}"
pause 5

echo "  >>> Back alive!"
# 0x31C = STATUS(3) + ALIVE(1) + player
send "31${P}" "01"
pause 0.3
send "32${P}" "001E"              # full ammo
pause 3

# ══════════════════════════════════════════════════
#  PHASE 5: Kill Streak
# ══════════════════════════════════════════════════
log "PHASE 5: KILL STREAK"

echo "  >>> Kill #2: enemy $ENEMY2!"
send "11${P}" "0${E2}"
pause 1
send "13${E2}"
pause 0.5
send "12${P}" "0002"
pause 0.5
send "30${P}" "012C"              # Score=300
pause 4

echo "  >>> Enemies respawn..."
send "22${E1}"
pause 0.3
send "22${E2}"
pause 2

echo "  >>> Kill #3: enemy $ENEMY1!"
send "11${P}" "0${E1}"
pause 1
send "13${E1}"
pause 0.5
send "12${P}" "0003"
pause 0.5
send "30${P}" "01F4"              # Score=500
pause 4

echo "  >>> Teammate scoring too..."
send "30${T}" "0096"              # Teammate Score=150
pause 0.3
send "12${T}" "0001"
pause 1

echo "  >>> Enemy $ENEMY2 scoring..."
send "30${E2}" "00C8"             # Enemy Score=200
pause 0.3
send "12${E2}" "0002"
pause 2

echo "  >>> Final score"
send "30${P}" "0258"              # Score=600
pause 3

# ══════════════════════════════════════════════════
#  PHASE 6: Second Death
# ══════════════════════════════════════════════════
log "PHASE 6: SECOND DEATH"

echo "  >>> Hit..."
send "10${P}" "0${E1}"
pause 1
echo "  >>> Hit again..."
send "10${P}" "0${E2}"
pause 1
echo "  >>> DEAD AGAIN!"
send "13${P}"
pause 6

echo "  >>> Respawn"
send "22${P}"
pause 5
send "31${P}" "01"
pause 0.3
send "32${P}" "001E"
pause 3

# ══════════════════════════════════════════════════
#  PHASE 7: Game End
#    CAN ID: 0x210 = GAME(2) + END(1) + broadcast(0)
# ══════════════════════════════════════════════════
log "PHASE 7: GAME END"

echo "  Final scores:"
echo "    Player $PLAYER: 600 pts, 3 kills"
echo "    Enemy $ENEMY1:  100 pts, 1 kill"
echo "    Enemy $ENEMY2:  200 pts, 2 kills"
echo "    Teammate $TEAMMATE: 150 pts, 1 kill"

send "30${P}" "0258"
pause 0.3
send "30${E1}" "0064"
pause 0.3
send "30${E2}" "00C8"
pause 0.3
send "30${T}" "0096"
pause 2

echo "  >>> GAME OVER!"
send "210"
pause 10

# ══════════════════════════════════════════════════
#  PHASE 8: Back to Idle
# ══════════════════════════════════════════════════
log "PHASE 8: IDLE"

send "01${P}"
pause 0.2
send "01${E1}"
pause 0.2
send "01${E2}"
pause 0.2
send "01${T}"
pause 5

echo ""
echo "╔══════════════════════════════════════════════╗"
echo "║     SIMULATION COMPLETE!                    ║"
echo "╚══════════════════════════════════════════════╝"
echo ""
echo "Tested:"
echo "  - Device registration (WAITING → IDLE)"
echo "  - Team assignment"
echo "  - Game countdown (5..1..GO)"
echo "  - Damage flash + sound"
echo "  - Hit marker sound"
echo "  - Kill popup + sound"
echo "  - Score popup"
echo "  - Low ammo blink + sound"
echo "  - Ammo empty sound"
echo "  - Death screen + sound"
echo "  - Respawn countdown"
echo "  - Kill streak"
echo "  - Game over screen"
echo "  - Return to idle"
echo ""
echo "Cleanup: restoring normal CAN config + restarting service..."
# cleanup runs via trap EXIT
