#include "mqtt_client.h"
#include <Arduino.h>
#include <ArduinoJson.h>
#include "can_bus.h"

namespace lt {
int playerId_ = 0;
String gameId = "";
static const unsigned long heartbeatInterval_ = 6000; // 6 seconds
static unsigned long lastHeartbeat_ = 0;
lt::CanBus can_;

// ── Begin ──
void MqttClient::begin(const char* host, int port, int playerId, CanBus can) {
    client_.setServer(host, port);
    playerId_ = playerId;
    can_ = can;

    // Set internal callback
    client_.setCallback([this](char* topic, byte* payload, unsigned int len){
        this->messageReceived(topic, payload, len);
    });
    
    // Attempt initial connection
    Serial.println("[MQTT] Connected to broker");
}

// ── Loop ──
void MqttClient::loop(int code) {
    if (!client_.connected()) {
        String clientId = "esp32-lasertag-" + String((uint32_t)ESP.getEfuseMac(), HEX);
        if (client_.connect(clientId.c_str())) {
            Serial.println("[MQTT] Connected to broker");
        } else {
            Serial.println("[MQTT] Connection failed");
        }

         // ── Send register message ──
        String topic = "device/" + (String)playerId_ + "/register";
        client_.publish(topic.c_str(), "{}");
        Serial.println("[MQTT] Sent register message to " + topic);

         // ── Send command message ──
        String cmdTopic = "device/" + (String)playerId_ + "/command";
        client_.publish(topic.c_str(), "{}");
        Serial.println("[MQTT] Sent command message to " + cmdTopic);

        // ── Subscribe to command topic ──
        client_.subscribe(cmdTopic.c_str());
        Serial.println("[MQTT] Subscribed to command topic: " + cmdTopic);
    }   
    client_.loop();

    if (code != -1 && gameId != "") {
        String eventTopic = "device/" + (String)code + "/event";
        String payload = "{\"game_id\":\"" + String(gameId) + "\",\"victim_id\":\"" + String(playerId_) + "\"}";
        client_.publish(eventTopic.c_str(), payload.c_str());
        Serial.println("[MQTT] Sent payload " + payload);
        Serial.println("[MQTT] Sent command message to " + eventTopic);
    }

    // Heartbeat every 6 seconds
    unsigned long now = millis();
    if (now - lastHeartbeat_ >= heartbeatInterval_) {
        String hbTopic = "device/" + String(playerId_) + "/heartbeat";
        client_.publish(hbTopic.c_str(), "{}");
        Serial.println("[MQTT] Sent heartbeat to " + hbTopic);
        lastHeartbeat_ = now; // update timestamp
    }
}

// ── Publish ──
void MqttClient::publish(const char* topic, const char* payload) {
    if (client_.connected()) {
        client_.publish(topic, payload);
    }
}

// ── Register message callback ──
void MqttClient::onMessage(MqttHandler cb) {
    handler_ = cb;
}

// ── Internal message wrapper ──
void MqttClient::messageReceived(char* topic, byte* payload, unsigned int len) {
    String msg;
    for (unsigned int i = 0; i < len; i++) {
        msg += (char)payload[i]; // build full payload string
    }

    // Print everything
    Serial.print("[MQTT] Received on topic: ");
    Serial.println(topic);
    Serial.print("[MQTT] Payload: ");
    Serial.println(msg);


    StaticJsonDocument<256> doc; // adjust size based on your payload

    // Try parsing JSON
    DeserializationError error = deserializeJson(doc, payload);
    if (error) {
        Serial.print("JSON parse error: ");
        Serial.println(error.c_str());
        return;
    }

    // Call user handler if set
    if(doc.containsKey("action") && doc.containsKey("game_id")) {
        String action = doc["action"].as<String>();

        if(action == "game_start") {
            // Send CAN frame: 0x200 = GAME(2) + START(0) + broadcast(0)
            can_.send(0x200, nullptr, 0);
            Serial.println("[CAN] Sent GAME_START frame");
        } else if(action == "game_end") {
            // Send CAN frame: 0x210 = GAME(2) + END(1) + broadcast(0)
            can_.send(0x210, nullptr, 0);
            Serial.println("[CAN] Sent GAME_END frame");
        }
    }


    if (handler_) {
        handler_(String(topic), msg);
    }
}

} // namespace lt