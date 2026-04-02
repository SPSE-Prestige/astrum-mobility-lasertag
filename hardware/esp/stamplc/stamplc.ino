// Latest ESP firmware for laser tag.
// See hardware/esp/stamplc/README.md for IR and CAN bus code semantics.
#include <M5Unified.h>
#include <Preferences.h>
#include "pins.h"
#include "credentials.h"
#include "src/ir_controller.h"
#include "src/can_bus.h"
#include "src/can_protocol.h"
#include "src/mqtt_client.h"

Preferences prefs;

lt::IRSender   irTx(PIN_IR_TX, []() {
    return M5.BtnA.isPressed() || (digitalRead(PIN_BTN_SHOOT) == LOW);
});
lt::IRReceiver irRx(PIN_IR_RX_1, PIN_IR_RX_2);

lt::CanBus can;
lt::MqttClient mqtt;
int  playerId = 0;
bool isDead   = false;

uint32_t getPlayerId() {
    prefs.begin("lasertag", true);
    uint32_t id = prefs.getUInt("playerId", 0);
    prefs.end();
    return id;
}

void setup() {
    M5.begin();
    Serial.begin(115200);
    pinMode(PIN_BTN_SHOOT, INPUT_PULLUP);
    delay(2000);
    Serial.println("BOOT OK");

    playerId = getPlayerId();
    Serial.print("PlayerID: ");
    Serial.println(playerId);

    M5.Display.fillScreen(BLACK);
    M5.Display.setTextColor(WHITE);
    M5.Display.setTextSize(2);
    M5.Display.setCursor(10, 20);
    M5.Display.println("Connecting to");
    M5.Display.println("game server...");

    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
    Serial.print("Connecting to WiFi");
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.println();
    Serial.println("WiFi connected");
    Serial.print("IP: ");
    Serial.println(WiFi.localIP());

    M5.Display.fillScreen(BLACK);
    M5.Display.setCursor(10, 20);
    M5.Display.println("Connected!");
    delay(1000);
    M5.Display.fillScreen(BLACK);

    irTx.begin();
    irTx.setPlayerId(playerId);
    irTx.onShoot([]() {
        M5.Display.fillScreen(RED);
    });
    irTx.onCooldown([]() {
        M5.Display.fillScreen(GREEN);
    });
    Serial.println("IR TX ready (Port A)");

    irRx.begin();
    irRx.setPlayerId(playerId);
    Serial.println("IR RX ready (Port C)");

    can.begin(PIN_CAN_TX, PIN_CAN_RX, 500000);
    if (can.registerPlayer(playerId)) {
        Serial.println("Player registered on CAN bus");
    } else {
        Serial.println("Failed to register player");
    }

    mqtt.begin(MQTT_HOST, MQTT_PORT, playerId, can, irTx);
    mqtt.onMessage([](const String& topic, const String& payload) {
        Serial.println("MQTT RX");
    });
    mqtt.onReconnecting([]() {
        M5.Display.fillScreen(BLACK);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("Reconnecting to");
        M5.Display.println("game server...");
    });
    mqtt.onReconnected([]() {
        M5.Display.fillScreen(BLACK);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("Reconnected!");
        delay(1000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.gameState().onKill([](uint8_t kills, uint8_t level) {
        M5.Display.fillScreen(BLACK);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.printf("KILL!\nKills: %d\nLevel: %d", kills, level);
        delay(1000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.gameState().onLevelUp([](uint8_t level) {
        M5.Display.fillScreen(MAGENTA);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.printf("LEVELUP %d", level);
        delay(3000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.onGameStart([]() {
        isDead = false;
        M5.Display.fillScreen(CYAN);
        M5.Display.setTextColor(BLACK);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("GAME START");
        delay(1000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.onGameEnd([]() {
        M5.Display.fillScreen(YELLOW);
        M5.Display.setTextColor(BLACK);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("GAME OVER");
        delay(3000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.onDie([]() {
        isDead = true;
        Serial.println("[GAME] Player died");
        // Motors off handled by MqttClient via CAN
        M5.Display.fillScreen(M5.Display.color565(80, 0, 0));
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("DEAD");
    });
    mqtt.onRespawn([]() {
        isDead = false;
        Serial.println("[GAME] Player respawned");
        // Motors on handled by MqttClient via CAN
        M5.Display.fillScreen(BLUE);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("RESPAWNED");
        delay(2000);
        M5.Display.fillScreen(BLACK);
    });
}

void loop() {
    M5.update();

    irTx.loop();
    int code = irRx.loop();

    if (code != -1 && !isDead) {
        Serial.print("[HIT] Received hit from player: 0x");
        Serial.println(code, HEX);
        M5.Display.fillScreen(BLACK);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("HITTED");
        delay(200);
        M5.Display.fillScreen(BLACK);
    }

    mqtt.loop(code);
}
