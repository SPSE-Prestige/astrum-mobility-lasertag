#include "ir_controller.h"
#include <Arduino.h>

namespace lt {

// --- IRSender ---

IRSender::IRSender(int txPin, ButtonReader btnReader)
: irsend_(txPin), buttonReader_(btnReader) {}

void IRSender::begin() {
    irsend_.begin();
}

void IRSender::setPlayerId(uint32_t id) {
    playerId_ = id;
}

void IRSender::onShoot(std::function<void()> cb) {
    shootCallback_ = cb;
}

void IRSender::onCooldown(std::function<void()> cb) {
    cooldownCallback_ = cb;
}

void IRSender::loop() {
    bool pressed = false;
    
    if (buttonReader_) {
        pressed = buttonReader_();
    }

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

void IRSender::send(uint32_t code) {
    irsend_.sendNEC(code, 32);
}


// --- IRReceiver ---

IRReceiver::IRReceiver(int rx1Pin, int rx2Pin)
: irrecv1_(rx1Pin), irrecv2_(rx2Pin) {}

void IRReceiver::begin() {
    irrecv1_.enableIRIn();
    irrecv2_.enableIRIn();
}

void IRReceiver::setPlayerId(uint32_t id) {
    playerId_ = id;
}

int IRReceiver::loop() {
    if (irrecv1_.decode(&results_)) {
        int code = results_.value;
        irrecv1_.resume();
        if ((uint32_t)code != playerId_) {
            Serial.print("IR RX1 hit: 0x");
            Serial.println(code, HEX);
            return code;
        }
    }
    if (irrecv2_.decode(&results_)) {
        int code = results_.value;
        irrecv2_.resume();
        if ((uint32_t)code != playerId_) {
            Serial.print("IR RX2 hit: 0x");
            Serial.println(code, HEX);
            return code;
        }
    }
    return -1;
}

} // namespace lt
