// Atom S3 firmware for laser tag.
// IR TX and RX both on Port A. WiFi + MQTT + CAN bus.
#include <M5Unified.h>
#include <Preferences.h>
#include "pins.h"
#include "credentials.h"
#include "src/ir_controller.h"
#include "src/can_bus.h"
#include "src/can_protocol.h"
#include "src/mqtt_client.h"

Preferences prefs;

// Port A — IR TX on GPIO2, IR RX on GPIO1, button via M5Unified or G41 (active LOW)
lt::IRController ir(PIN_IR_TX, PIN_IR_RX, []() {
    return M5.BtnA.isPressed() || (digitalRead(PIN_BTN_SHOOT) == LOW);
});
lt::CanBus can;
lt::MqttClient mqtt;

int  playerId    = 0;
bool isDead      = false;
unsigned long flashUntilMs = 0;

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

    // WiFi
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

    // IR — Port A (TX + RX)
    ir.begin();
    ir.setPlayerId(playerId);
    ir.onShoot([]() {
        M5.Display.fillScreen(RED);
        flashUntilMs = millis() + 200;
        mqtt.publishShoot();
    });
    ir.onCooldown([]() {
        M5.Display.fillScreen(GREEN);
        flashUntilMs = millis() + 200;
    });
    Serial.println("IR ready (Port A)");

    // CAN bus
    can.begin(PIN_CAN_TX, PIN_CAN_RX, 500000);
    if (can.registerPlayer(playerId)) {
        Serial.println("Player registered on CAN bus");
    } else {
        Serial.println("Failed to register player");
    }

    // MQTT
    mqtt.begin(MQTT_HOST, MQTT_PORT, playerId, can);
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
        ir.setTeamId(mqtt.teamId());
        ir.setGameActive(true);
        M5.Display.fillScreen(CYAN);
        M5.Display.setTextColor(BLACK);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        M5.Display.println("GAME START");
        delay(1000);
        M5.Display.fillScreen(BLACK);
    });
    mqtt.onRejoin([](bool isAlive, uint8_t weaponLevel) {
        isDead = !isAlive;
        ir.setTeamId(mqtt.teamId());
        ir.setGameActive(isAlive);
        M5.Display.fillScreen(BLACK);
        M5.Display.setTextColor(WHITE);
        M5.Display.setTextSize(2);
        M5.Display.setCursor(10, 20);
        if (isAlive) {
            M5.Display.printf("REJOINED\nLvl: %d", weaponLevel);
        } else {
            M5.Display.printf("REJOINED\nDEAD\nLvl: %d", weaponLevel);
        }
        delay(2000);
        if (!isAlive) {
            M5.Display.fillScreen(M5.Display.color565(80, 0, 0));
            M5.Display.setTextColor(WHITE);
            M5.Display.setTextSize(2);
            M5.Display.setCursor(10, 20);
            M5.Display.println("DEAD");
        } else {
            M5.Display.fillScreen(BLACK);
        }
    });
    mqtt.onGameEnd([]() {
        ir.setGameActive(false);
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
        ir.setGameActive(false);
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
        ir.setGameActive(true);
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

    if (flashUntilMs && millis() >= flashUntilMs && !isDead) {
        flashUntilMs = 0;
        M5.Display.fillScreen(BLACK);
    }

    ir.loop();
    int code = ir.receive();

    if (code != -1 && !isDead) {
        Serial.print("[HIT] Received hit from player: 0x");
        Serial.println(code & 0xF, HEX);
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
