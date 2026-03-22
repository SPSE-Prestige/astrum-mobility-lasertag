#include <Preferences.h>
Preferences prefs;

const int PLAYER_IDD = 14;

void savePlayerId(uint32_t playerId) {
    Serial.begin(9600);
    prefs.begin("lasertag", false); // "lasertag" is namespace
    prefs.putUInt("playerId", playerId);
    prefs.end();
}

uint32_t getPlayerId() {
    prefs.begin("lasertag", true); // read-only
    uint32_t id = prefs.getUInt("playerId", 0); // default 0 if not set
    prefs.end();
    return id;
}

void setup() {
    Serial.begin(115200);

    savePlayerId(PLAYER_IDD);
    Serial.println(getPlayerId());
}

void loop () {}