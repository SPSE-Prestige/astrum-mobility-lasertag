#include "mqtt_client.h"
#include <Arduino.h>


namespace lt {

static MqttHandler globalHandler;

static void callback(char* topic, byte* payload, unsigned int len) {
    String msg;
    for (unsigned int i=0;i<len;i++) msg += (char)payload[i];
    if (globalHandler) globalHandler(topic, msg);
}

void MqttClient::begin(const char* host, int port) {
    client_.setServer(host, port);
    client_.setCallback(callback);
}

void MqttClient::loop() {
    if (!client_.connected()) {
        Serial.print("not connected to MQTT");
        client_.connect("esp32-lasertag");
    }
    
    client_.loop();
}

void MqttClient::publish(const char* topic, const char* payload) {
    client_.publish(topic, payload);
}

void MqttClient::onMessage(MqttHandler cb) {
    globalHandler = cb;
}

}