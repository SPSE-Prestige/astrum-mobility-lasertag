#include "ir_controller.h"
#include <Arduino.h>

namespace lt {

IRController::IRController(int txPin, int rxPin, int btnPin)
: irsend_(txPin), irrecv_(rxPin), buttonPin_(btnPin) {}

void IRController::begin() {
    irsend_.begin();
    irrecv_.enableIRIn();
    pinMode(buttonPin_, INPUT_PULLUP);  // button connects to GND
}

void IRController::loop(int playerId) {
    // Receive IR codes
    if (irrecv_.decode(&results_)) {
        Serial.print("Received IR code: 0x");
        Serial.println(results_.value, HEX);
        irrecv_.resume();  // ready to receive next code
    }

    // Send code if button pressed
    if (digitalRead(buttonPin_) == LOW) {
        send(playerId);  // send playerId as code
        Serial.print("Sent IR code: 0x");
        Serial.println(playerId, HEX);
        delay(200);  // simple debounce
    }
}

void IRController::send(uint32_t code) {
    irsend_.sendNEC(code, 32);
}

} // namespace lt