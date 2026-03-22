#include "ir_controller.h"
#include "can_bus.h"
#include "mqtt_client.h"
#include <Preferences.h>

Preferences prefs;

lt::IRController ir(2,1,41);
lt::CanBus can;
lt::MqttClient mqtt;

int playerId = 0;

uint32_t getPlayerId() {
    prefs.begin("lasertag", true); // read-only
    uint32_t id = prefs.getUInt("playerId", 0); // default 0 if not set
    prefs.end();
    return id;
}

void setup() {
    Serial.begin(750);
    delay(2000);
    Serial.println("BOOT OK");

    // Get player id
    playerId = getPlayerId();
    Serial.print("PlayerID: ");
    Serial.println(playerId);

    // WIFI
    WiFi.mode(WIFI_STA);          // <-- important (client mode)
    WiFi.begin("MERCURSYS_4F0A", "qwertzuiop01");
    Serial.print("Connecting to WiFi");

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("WiFi connected");
    Serial.print("IP: ");
    Serial.println(WiFi.localIP());
    Serial.println("Setted up wifi");

    // WIFI

    ir.begin();
    Serial.println("Setted up ir");

    can.begin(5,6,500000);
    Serial.println("Setted up can");

    mqtt.begin("192.168.1.10",1883);
    Serial.println("Setted up mqtt");


    can.onReceive([](const lt::CanFrame& f){
        Serial.println("CAN RX");
    });

    mqtt.onMessage([](const String& topic, const String& payload){
        Serial.println("MQTT RX");
    });
}

void loop() {

    ir.loop(playerId);
    can.loop();
    // mqtt.loop();
}