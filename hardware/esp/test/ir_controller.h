#pragma once
#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <IRrecv.h>
#include <IRutils.h>

namespace lt {

class IRController {
public:
    IRController(int txPin, int rxPin, int btnPin);

    void begin();
    void loop(int playerId);

private:
    void send(uint32_t code);

    IRsend irsend_;
    IRrecv irrecv_;
    decode_results results_;
    int buttonPin_;
};

} // namespace lt