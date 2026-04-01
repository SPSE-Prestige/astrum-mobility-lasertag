#pragma once
#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <IRrecv.h>
#include <IRutils.h>
#include <functional>

namespace lt {

using ButtonReader = std::function<bool()>;

// IR controller for Atom S3 — Port A handles both TX and RX
class IRController {
public:
    IRController(int txPin, int rxPin, ButtonReader btnReader);

    void begin();
    void loop();              // handle shoot button
    int  receive();           // returns attacker player ID, or -1 if none
    void send(uint32_t code);
    void setPlayerId(uint32_t id);
    void onShoot(std::function<void()> cb);
    void onCooldown(std::function<void()> cb);

private:
    IRsend irsend_;
    IRrecv irrecv_;
    decode_results results_;

    ButtonReader buttonReader_;
    std::function<void()> shootCallback_;
    std::function<void()> cooldownCallback_;
    uint32_t playerId_ = 0;
    unsigned long lastShootMs_ = -2000UL;
    bool wasPressed_ = false;
};

} // namespace lt
