#include "mqtt_client.h"
#include <Arduino.h>
#include <ArduinoJson.h>
#include "can_bus.h"

namespace lt {
int playerId_ = 0;
String gameId = "";
static const unsigned long heartbeatInterval_ = 6000; // 6 seconds
static unsigned long lastHeartbeat_ = 0;
static bool registered_ = false;
lt::CanBus can_;

// ── Begin ──
void MqttClient::begin(const char* host, int port, int playerId, CanBus can) {
    client_.setServer(host, port);
    playerId_ = playerId;
    can_ = can;
    registered_ = false;

    // Set internal callback
    client_.setCallback([this](char* topic, byte* payload, unsigned int len){
        this->messageReceived(topic, payload, len);
    });
    
    Serial.println("[MQTT] Initialized");
}

// ── Loop ──
void MqttClient::loop(int code) {
    if (!client_.connected()) {
        String clientId = "esp32-lasertag-" + String((uint32_t)ESP.getEfuseMac(), HEX);
        if (client_.connect(clientId.c_str())) {
            Serial.println("[MQTT] Connected to broker");

            // Subscribe to command topic
            String cmdTopic = "device/" + String(playerId_) + "/command";
            client_.subscribe(cmdTopic.c_str());
            Serial.println("[MQTT] Subscribed to: " + cmdTopic);

            // Send register (once per boot, or on reconnect)
            String regTopic = "device/" + String(playerId_) + "/register";
            client_.publish(regTopic.c_str(), "{}");
            Serial.println("[MQTT] Sent register to " + regTopic);
            registered_ = true;
        } else {
            Serial.println("[MQTT] Connection failed");
            return;
        }
    }
    client_.loop();

    // IR hit event: code = attacker's player ID received via IR
    if (code != -1 && gameId != "") {
        // Publish as: device/{attacker}/event  payload: {game_id, victim_id: self}
        String eventTopic = "device/" + String(code) + "/event";
        StaticJsonDocument<128> doc;
        doc["game_id"] = gameId;
        doc["victim_id"] = String(playerId_);
        String payload;
        serializeJson(doc, payload);
        client_.publish(eventTopic.c_str(), payload.c_str());
        Serial.println("[MQTT] Hit event: attacker=" + String(code) + " victim=" + String(playerId_));
    }

    // Heartbeat every 6 seconds
    unsigned long now = millis();
    if (now - lastHeartbeat_ >= heartbeatInterval_) {
        String hbTopic = "device/" + String(playerId_) + "/heartbeat";
        client_.publish(hbTopic.c_str(), "{}");
        lastHeartbeat_ = now;
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
        msg += (char)payload[i];
    }

    Serial.print("[MQTT] RX topic: ");
    Serial.println(topic);
    Serial.print("[MQTT] RX payload: ");
    Serial.println(msg);

    StaticJsonDocument<256> doc;
    DeserializationError error = deserializeJson(doc, payload, len);
    if (error) {
        Serial.print("[MQTT] JSON parse error: ");
        Serial.println(error.c_str());
        return;
    }

    if (!doc.containsKey("action")) {
        return;
    }

    String action = doc["action"].as<String>();
    uint8_t pid = playerId_ & 0xF;

    // ── Forward MQTT commands to CAN bus for MP135 ──

    if (action == "game_start") {
        // Store game ID for hit events
        if (doc.containsKey("game_id")) {
            gameId = doc["game_id"].as<String>();
        }
        // CAN 0x200 = GAME(2) + START(0) + broadcast(0)
        can_.send(0x200, nullptr, 0);
        Serial.println("[CAN] >> GAME_START 0x200");
    }
    else if (action == "game_end") {
        gameId = "";
        // CAN 0x210 = GAME(2) + END(1) + broadcast(0)
        can_.send(0x210, nullptr, 0);
        Serial.println("[CAN] >> GAME_END 0x210");
    }
    else if (action == "kill_confirmed") {
        // Attacker got a kill — forward to MP135
        // CAN 0x11P = COMBAT(1) + KILL(1) + player(P), data=[victim_id]
        uint8_t victim = 0;
        if (doc.containsKey("victim_id")) {
            victim = doc["victim_id"].as<int>() & 0xF;
        }
        uint8_t data[1] = { victim };
        can_.send(0x110 | pid, data, 1);
        Serial.printf("[CAN] >> KILL 0x%03X victim=%d\n", 0x110 | pid, victim);

        // Forward score update if included
        if (doc.containsKey("score")) {
            uint16_t score = doc["score"].as<int>();
            uint8_t sdata[2] = { (uint8_t)(score >> 8), (uint8_t)(score & 0xFF) };
            can_.send(0x300 | pid, sdata, 2);
            Serial.printf("[CAN] >> SCORE 0x%03X score=%d\n", 0x300 | pid, score);
        }

        // Forward kill count if included
        if (doc.containsKey("kills")) {
            uint16_t kills = doc["kills"].as<int>();
            uint8_t kdata[2] = { (uint8_t)(kills >> 8), (uint8_t)(kills & 0xFF) };
            can_.send(0x120 | pid, kdata, 2);
            Serial.printf("[CAN] >> KILLCOUNT 0x%03X kills=%d\n", 0x120 | pid, kills);
        }
    }
    else if (action == "die") {
        // This player died — forward to MP135
        // CAN 0x13P = COMBAT(1) + DEATH(3) + player(P)
        can_.send(0x130 | pid, nullptr, 0);
        Serial.printf("[CAN] >> DEATH 0x%03X\n", 0x130 | pid);
    }
    else if (action == "respawn") {
        // This player respawns — forward to MP135
        // CAN 0x22P = GAME(2) + RESPAWN(2) + player(P)
        can_.send(0x220 | pid, nullptr, 0);
        Serial.printf("[CAN] >> RESPAWN 0x%03X\n", 0x220 | pid);
    }

    // Call user handler if set
    if (handler_) {
        handler_(String(topic), msg);
    }
}

} // namespace lt