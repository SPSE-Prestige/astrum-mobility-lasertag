#include "ir_controller.h"
#include <Arduino.h>

namespace lt {

// Constructor only sets pins
IRController::IRController(int txPin, int rxPin, int btnPin)
: irsend_(txPin), irrecv_(rxPin), buttonPin_(btnPin) {}

// Set player ID
void IRController::setPlayerId(uint32_t id) {
    playerId_ = id;
}

void IRController::begin() {
    irsend_.begin();
    irrecv_.enableIRIn();
    pinMode(buttonPin_, INPUT_PULLUP);  // button connects to GND
}

void IRController::sendloop() {
    if (digitalRead(buttonPin_) == LOW) {
        // Send this player's ID as an NEC IR code.
        // Other devices interpret the received code as the shooter/player who fired.
        send(playerId_);
        Serial.print("Sent IR code: 0x");
        Serial.println(playerId_, HEX);
        delay(200);  // simple debounce
    }
}
int IRController::reciveloop() {
    if (irrecv_.decode(&results_)) {
        int code = results_.value;
        irrecv_.resume();
        if (code != playerId_) {  // ignore own code
            Serial.print("Received IR code: 0x");
            Serial.println(code, HEX);
            return code; // returned code is the attacker/player ID
        }
    }
    return -1;
}

void IRController::send(uint32_t code) {
    // NEC sends a 32-bit value. Here `code` is the player ID.
    irsend_.sendNEC(code, 32);
}

} // namespace lt