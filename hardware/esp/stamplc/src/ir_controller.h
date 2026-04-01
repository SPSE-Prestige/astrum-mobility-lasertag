#pragma once
#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <IRrecv.h>
#include <IRutils.h>
#include <functional>

namespace lt {

using ButtonReader = std::function<bool()>;

// IR emitter — Port A, transmits player ID on button press
class IRSender {
public:
    IRSender(int txPin, ButtonReader btnReader);

    void begin();
    void loop();
    void send(uint32_t code);
    void setPlayerId(uint32_t id);
    void onShoot(std::function<void()> cb);
    void onCooldown(std::function<void()> cb);

private:
    IRsend irsend_;
    ButtonReader buttonReader_;
    std::function<void()> shootCallback_;
    std::function<void()> cooldownCallback_;
    uint32_t playerId_ = 0;
    unsigned long lastShootMs_ = -2000UL;
    bool wasPressed_ = false;
};

// IR receiver — Port C, two independent sensors
class IRReceiver {
public:
    // rx1Pin, rx2Pin: Port C GPIO5 and GPIO4
    IRReceiver(int rx1Pin, int rx2Pin);

    void begin();
    int loop();  // returns attacker player ID, or -1 if no signal

    void setPlayerId(uint32_t id);  // to ignore own codes

private:
    IRrecv irrecv1_;
    IRrecv irrecv2_;
    decode_results results_;
    uint32_t playerId_ = 0;
};

} // namespace lt
