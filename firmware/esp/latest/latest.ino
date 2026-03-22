#include "ir_controller.h"
#include "can_bus.h"
#include "mqtt_client.h"
#include <Preferences.h>

Preferences prefs;

lt::IRController ir(2,1,41);
lt::CanBus can;
lt::MqttClient mqtt;
int code = -1;
int eventCode = -1;
int playerId = 0;

uint32_t getPlayerId() {
    prefs.begin("lasertag", true); // read-only
    uint32_t id = prefs.getUInt("playerId", 0); // default 0 if not set
    prefs.end();
    return id;
}

void setup() {
    Serial.begin(115200);
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
    ir.setPlayerId(playerId);
    Serial.println("Setted up ir");

    can.begin(5,6,1000000);
    if(can.registerPlayer(playerId)) {
        Serial.println("Player registered on CAN bus");
    } else {
        Serial.println("Failed to register player");
    }
    Serial.println("Setted up can");

    mqtt.begin("192.168.0.196",1883,playerId, can);
    mqtt.onMessage([](const String& topic, const String& msg){
        Serial.println("Setted up mqtt");
    });

    mqtt.onMessage([](const String& topic, const String& payload){
        Serial.println("MQTT RX");
    });
}

void loop() {

    ir.sendloop();
    code = ir.reciveloop();

    mqtt.loop(code);
    code = -1;

}