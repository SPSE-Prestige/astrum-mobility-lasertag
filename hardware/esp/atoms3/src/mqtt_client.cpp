#include "mqtt_client.h"
#include <Arduino.h>
#include <ArduinoJson.h>
#include "can_protocol.h"

namespace lt {

void MqttClient::begin(const char* host, int port, int playerId, CanBus can) {
    client_.setServer(host, port);
    playerId_ = playerId;
    can_      = can;
    registered_ = false;

    client_.setCallback([this](char* topic, byte* payload, unsigned int len){
        this->messageReceived(topic, payload, len);
    });
    Serial.println("[MQTT] Initialized");
}

void MqttClient::loop(int code) {
    unsigned long now = millis();

    if (!client_.connected()) {
        if (now - lastReconnectMs_ < reconnectCooldown_) return;
        lastReconnectMs_ = now;

        if (reconnectingCallback_) reconnectingCallback_();

        String clientId = "esp32-lasertag-" + String((uint32_t)ESP.getEfuseMac(), HEX);
        Serial.println("[MQTT] Reconnecting...");
        if (client_.connect(clientId.c_str())) {
            Serial.println("[MQTT] Connected");
            String cmdTopic = "device/" + String(playerId_) + "/command";
            client_.subscribe(cmdTopic.c_str());
            client_.publish(("device/" + String(playerId_) + "/register").c_str(), "{}");
            registered_ = true;
            if (reconnectedCallback_) reconnectedCallback_();

            uint8_t pid = playerId_ & 0xF;
            can_.send(CAN_SYS_REGISTER(pid), nullptr, 0);
            uint8_t alive = 1;
            can_.send(CAN_STATUS_ALIVE(pid), &alive, 1);
        } else {
            Serial.println("[MQTT] Connection failed");
            return;
        }
    }

    client_.loop();

    if (code != -1 && gameId_ != "") {
        int attackerId   = code & 0xF;
        int attackerTeam = (code >> 4) & 0xF;
        String eventTopic = "device/" + String(attackerId) + "/event";
        StaticJsonDocument<128> doc;
        doc["game_id"]       = gameId_;
        doc["victim_id"]     = playerId_;
        doc["attacker_team"] = attackerTeam;
        if (teamId_ >= 0) doc["team_id"] = teamId_;
        String payload;
        serializeJson(doc, payload);
        client_.publish(eventTopic.c_str(), payload.c_str());
        Serial.printf("[MQTT] Hit: attacker=%d team=%d victim=%d\n", attackerId, attackerTeam, playerId_);
    }

    if (now - lastHeartbeat_ >= heartbeatInterval_) {
        client_.publish(("device/" + String(playerId_) + "/heartbeat").c_str(), "{}");
        lastHeartbeat_ = now;
    }
}

void MqttClient::publishShoot() {
    if (!client_.connected() || gameId_ == "") return;
    StaticJsonDocument<128> doc;
    doc["action"]  = "weapon_shoot";
    doc["game_id"] = gameId_;
    String payload;
    serializeJson(doc, payload);
    String topic = "device/" + String(playerId_) + "/event";
    client_.publish(topic.c_str(), payload.c_str());
    Serial.printf("[MQTT] weapon_shoot player=%d\n", playerId_);
}

void MqttClient::publish(const char* topic, const char* payload) {
    if (client_.connected()) client_.publish(topic, payload);
}

void MqttClient::onMessage(MqttHandler cb)                        { handler_ = cb; }
void MqttClient::onReconnecting(std::function<void()> cb)         { reconnectingCallback_ = cb; }
void MqttClient::onReconnected(std::function<void()> cb)          { reconnectedCallback_ = cb; }
void MqttClient::onDie(std::function<void()> cb)                  { dieCallback_ = cb; }
void MqttClient::onRespawn(std::function<void()> cb)              { respawnCallback_ = cb; }
void MqttClient::onGameStart(std::function<void()> cb)            { gameStartCallback_ = cb; }
void MqttClient::onGameEnd(std::function<void()> cb)              { gameEndCallback_ = cb; }
void MqttClient::onRejoin(std::function<void(bool, uint8_t)> cb)  { rejoinCallback_ = cb; }

void MqttClient::messageReceived(char* topic, byte* payload, unsigned int len) {
    String msg;
    for (unsigned int i = 0; i < len; i++) msg += (char)payload[i];

    Serial.printf("[MQTT] RX %s: %s\n", topic, msg.c_str());

    StaticJsonDocument<384> doc;
    if (deserializeJson(doc, payload, len) || !doc.containsKey("action")) return;

    String action = doc["action"].as<String>();
    uint8_t pid   = playerId_ & 0xF;

    if (action == "game_rejoin") {
        if (doc.containsKey("game_id")) gameId_ = doc["game_id"].as<String>();

        uint8_t weaponLevel  = doc.containsKey("weapon_level")      ? doc["weapon_level"].as<uint8_t>()      : 0;
        uint8_t kills        = doc.containsKey("kills")             ? doc["kills"].as<uint8_t>()             : 0;
        bool    isAlive      = doc.containsKey("is_alive")          ? doc["is_alive"].as<bool>()             : true;
        teamId_         = doc.containsKey("team_id")           ? doc["team_id"].as<int>()              : -1;
        isFriendlyFire_ = doc.containsKey("is_friendlyfire")   ? doc["is_friendlyfire"].as<bool>()     : false;
        gameType_       = (GameType)(doc.containsKey("type_of_the_game") ? doc["type_of_the_game"].as<int>() : 0);

        Serial.printf("[GAME] Rejoin type=%d team=%d ff=%d\n", gameType_, teamId_, isFriendlyFire_);

        gameState_.setLevel(weaponLevel);
        gameState_.setKills(kills);

        uint8_t ldata[2] = { kills, weaponLevel };
        can_.send(CAN_STATUS_SCORE(pid), ldata, 2);

        if (isAlive) {
            can_.send(CAN_HW_MOTOR_ON(pid), nullptr, 0);
        } else {
            can_.send(CAN_HW_MOTOR_OFF(pid), nullptr, 0);
        }

        Serial.printf("[GAME] Rejoin game=%s alive=%d kills=%d weaponLevel=%d\n",
            gameId_.c_str(), isAlive, kills, weaponLevel);

        if (rejoinCallback_) rejoinCallback_(isAlive, weaponLevel);
    }
    else if (action == "game_start") {
        if (doc.containsKey("game_id")) gameId_ = doc["game_id"].as<String>();
        teamId_         = doc.containsKey("team_id")           ? doc["team_id"].as<int>()              : -1;
        isFriendlyFire_ = doc.containsKey("is_friendlyfire")   ? doc["is_friendlyfire"].as<bool>()     : false;
        gameType_       = (GameType)(doc.containsKey("type_of_the_game") ? doc["type_of_the_game"].as<int>() : 0);
        gameState_.reset();
        can_.send(CAN_GAME_START, nullptr, 0);
        if (teamId_ >= 0) {
            uint8_t tdata[1] = { (uint8_t)teamId_ };
            can_.send(CAN_GAME_TEAM(pid), tdata, 1);
        }
        Serial.printf("[GAME] Start type=%d team=%d ff=%d\n", gameType_, teamId_, isFriendlyFire_);
        if (gameStartCallback_) gameStartCallback_();
    }
    else if (action == "game_end") {
        gameId_ = "";
        teamId_ = -1;
        isFriendlyFire_ = false;
        gameType_ = GAME_DEATHMATCH;
        can_.send(CAN_GAME_END, nullptr, 0);
        if (gameEndCallback_) gameEndCallback_();
    }
    else if (action == "kill_confirmed") {
        uint8_t victim = doc.containsKey("victim_id") ? (doc["victim_id"].as<int>() & 0xF) : 0;
        uint8_t vdata[1] = { victim };
        can_.send(CAN_COMBAT_KILL(pid), vdata, 1);

        gameState_.addKill();
        uint8_t kdata[1] = { gameState_.kills() };
        can_.send(CAN_COMBAT_KILL_CNT(pid), kdata, 1);
        uint8_t ldata[2] = { gameState_.kills(), gameState_.level() };
        can_.send(CAN_STATUS_SCORE(pid), ldata, 2);
        Serial.printf("[GAME] kills=%d level=%d\n", gameState_.kills(), gameState_.level());
    }
    else if (action == "weapon_upgrade") {
        gameState_.upgrade();
        uint8_t ldata[2] = { gameState_.kills(), gameState_.level() };
        can_.send(CAN_STATUS_SCORE(pid), ldata, 2);
        Serial.printf("[GAME] Upgrade → level=%d\n", gameState_.level());
    }
    else if (action == "die") {
        gameState_.reset();
        can_.send(CAN_COMBAT_DEATH(pid), nullptr, 0);
        can_.send(CAN_HW_MOTOR_OFF(pid), nullptr, 0);
        if (dieCallback_) dieCallback_();
    }
    else if (action == "respawn") {
        can_.send(CAN_GAME_RESPAWN(pid), nullptr, 0);
        can_.send(CAN_HW_MOTOR_ON(pid), nullptr, 0);
        if (respawnCallback_) respawnCallback_();
    }

    if (handler_) handler_(String(topic), msg);
}

} // namespace lt
