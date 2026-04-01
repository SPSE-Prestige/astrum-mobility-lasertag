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

void IRController::onShoot(std::function<void()> cb) {
    shootCallback_ = cb;
}

void IRController::onCooldown(std::function<void()> cb) {
    cooldownCallback_ = cb;
}

void IRController::loop() {
    if (!buttonReader_) return;

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
        send(playerId_);
        if (shootCallback_) shootCallback_();
        Serial.print("[IR TX] Fired code: 0x");
        Serial.println(playerId_, HEX);
    }
}

int IRController::receive() {
    if (irrecv_.decode(&results_)) {
        int code = results_.value;
        irrecv_.resume();
        if ((uint32_t)code != playerId_) {
            Serial.print("[IR RX] Hit from: 0x");
            Serial.println(code, HEX);
            return code;
        }
    }
    return -1;
}

void IRController::send(uint32_t code) {
    irsend_.sendNEC(code, 32);
}

} // namespace lt
