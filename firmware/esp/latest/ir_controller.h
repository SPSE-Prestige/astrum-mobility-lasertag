#pragma once
#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <IRrecv.h>
#include <IRutils.h>

namespace lt {

class IRController {
public:
    // Constructor only sets pins
    IRController(int txPin, int rxPin, int btnPin);

    void begin();
    void sendloop();
    int reciveloop();

    // Set player ID after construction
    void setPlayerId(uint32_t id);

    bool available() { return hasCode_; }
    uint32_t getCode() { 
        hasCode_ = false; 
        return code_; 
    }

    void send(uint32_t code);

private:
    IRsend irsend_;
    IRrecv irrecv_;
    decode_results results_;
    bool hasCode_ = false;
    uint32_t code_ = 0;

    int buttonPin_;
    uint32_t playerId_ = 0;  // player ID, set later
};

} // namespace lt