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

void IRSender::setTeamId(int team) {
    teamId_ = team;
}

void IRSender::setGameActive(bool active) {
    gameActive_ = active;
}

void IRSender::setCooldown(unsigned long ms) {
    cooldownMs_ = ms;
    Serial.printf("[IR TX] Cooldown set to %lums\n", ms);
}

void IRSender::enableSecondTx(int txPin) {
    if (!irsend2_) {
        irsend2_ = new IRsend(txPin);
        irsend2_->begin();
        Serial.printf("[IR TX] Second emitter enabled on pin %d\n", txPin);
    }
}

void IRSender::onShoot(std::function<void()> cb) {
    shootCallback_ = cb;
}

void IRSender::onCooldown(std::function<void()> cb) {
    cooldownCallback_ = cb;
}

void IRSender::loop() {
    
#if DEBUG
    if (!gameActive_) return;

    bool pressed = false;

    if (buttonReader_) {
        pressed = buttonReader_();
    }

    bool pressed = buttonReader_();

    if (pressed && !wasPressed_) {
        unsigned long elapsed = millis() - lastShootMs_;
        Serial.print("[IR TX] btn=1  cooldown_remaining=");
        Serial.println(elapsed < cooldownMs_ ? cooldownMs_ - elapsed : 0);
    }
#endif
    bool pressed = false;

    if (buttonReader_) {
        pressed = buttonReader_();
    }
    wasPressed_ = pressed;

    if (pressed) {
        unsigned long now = millis();
        if (now - lastShootMs_ < cooldownMs_) {
            if (cooldownCallback_) cooldownCallback_();
            return;
        }
        lastShootMs_ = now;
        uint32_t code = ((uint32_t)(teamId_ & 0xF) << 4) | (playerId_ & 0xF);
        send(code);
        if (shootCallback_) shootCallback_();
        Serial.printf("[IR TX] Fired team=%d player=%d code=0x%X\n", teamId_, playerId_, code);
    }
}

void IRSender::send(uint32_t code) {
    irsend_.sendNEC(code, 32);
    if (irsend2_) irsend2_->sendNEC(code, 32);
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

void IRReceiver::setTeamId(int team) {
    teamId_ = team;
}

void IRReceiver::lockout(unsigned long ms) {
    lockoutUntilMs_ = millis() + ms;
}

int IRReceiver::loop() {
    uint32_t ownCode = ((uint32_t)(teamId_ & 0xF) << 4) | (playerId_ & 0xF);
    bool locked = millis() < lockoutUntilMs_;

    if (irrecv1_.decode(&results_)) {
        uint32_t code = results_.value;
        irrecv1_.resume();
        if (!locked && code != 0xFFFFFFFF && code != ownCode) {
            Serial.printf("IR RX1 hit: team=%d player=%d\n", (code >> 4) & 0xF, code & 0xF);
            return (int)code;
        }
    }
    if (irrecv2_.decode(&results_)) {
        uint32_t code = results_.value;
        irrecv2_.resume();
        if (!locked && code != 0xFFFFFFFF && code != ownCode) {
            Serial.printf("IR RX2 hit: team=%d player=%d\n", (code >> 4) & 0xF, code & 0xF);
            return (int)code;
        }
    }
    return -1;
}

} // namespace lt
