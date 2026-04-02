#include "ir_controller.h"
#include <Arduino.h>

namespace lt {

IRController::IRController(int txPin, int rxPin, ButtonReader btnReader)
: irsend_(txPin), irrecv_(rxPin), buttonReader_(btnReader) {}

void IRController::begin() {
    irsend_.begin();
    irrecv_.enableIRIn();
}

void IRController::setPlayerId(uint32_t id) {
    playerId_ = id;
}

void IRController::setTeamId(int team) {
    teamId_ = team;
}

void IRController::setGameActive(bool active) {
    gameActive_ = active;
}

void IRController::onShoot(std::function<void()> cb) {
    shootCallback_ = cb;
}

void IRController::onCooldown(std::function<void()> cb) {
    cooldownCallback_ = cb;
}

void IRController::loop() {
    if (!buttonReader_ || !gameActive_) return;

    bool pressed = buttonReader_();

#if DEBUG
    if (pressed && !wasPressed_) {
        unsigned long elapsed = millis() - lastShootMs_;
        Serial.print("[IR TX] btn=1  cooldown_remaining=");
        Serial.println(elapsed < 2000 ? 2000 - elapsed : 0);
    }
#endif
    wasPressed_ = pressed;

    if (pressed) {
        unsigned long now = millis();
        if (now - lastShootMs_ < 2000) {
            if (cooldownCallback_) cooldownCallback_();
            return;
        }
        lastShootMs_ = now;
        lockoutUntilMs_ = now + 50;
        uint32_t code = ((uint32_t)(teamId_ & 0xF) << 4) | (playerId_ & 0xF);
        send(code);
        if (shootCallback_) shootCallback_();
        Serial.printf("[IR TX] Fired team=%d player=%d code=0x%X\n", teamId_, playerId_, code);
    }
}

int IRController::receive() {
    uint32_t ownCode = ((uint32_t)(teamId_ & 0xF) << 4) | (playerId_ & 0xF);
    bool locked = millis() < lockoutUntilMs_;
    if (irrecv_.decode(&results_)) {
        uint32_t code = results_.value;
        irrecv_.resume();
        if (!locked && code != 0xFFFFFFFF && code != ownCode) {
            Serial.printf("[IR RX] Hit from: team=%d player=%d\n", (code >> 4) & 0xF, code & 0xF);
            return (int)code;
        }
    }
    return -1;
}

void IRController::send(uint32_t code) {
    irsend_.sendNEC(code, 32);
}

} // namespace lt
