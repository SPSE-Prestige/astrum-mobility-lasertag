#include <Preferences.h>
Preferences prefs;

const int PLAYER_IDD = 1;  // change this per device (1–15)

void savePlayerId(uint32_t playerId) {
    prefs.begin("lasertag", false);
    prefs.putUInt("playerId", playerId);
    prefs.end();
}

uint32_t getPlayerId() {
    prefs.begin("lasertag", true);
    uint32_t id = prefs.getUInt("playerId", 0);
    prefs.end();
    return id;
}

void setup() {
    Serial.begin(115200);
    delay(2000);
    savePlayerId(PLAYER_IDD);
    Serial.print("Saved PlayerID: ");
    Serial.println(getPlayerId());
}

void loop() {}
