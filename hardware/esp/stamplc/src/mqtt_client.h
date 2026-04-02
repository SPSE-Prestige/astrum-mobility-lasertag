#pragma once
#include <WiFi.h>
#include <PubSubClient.h>
#include <functional>
#include "can_bus.h"
#include "game_state.h"
#include "ir_controller.h"

namespace lt {

enum GameType {
    GAME_DEATHMATCH      = 0,
    GAME_TEAM_DEATHMATCH = 1,
};

using MqttHandler = std::function<void(const String&, const String&)>;

class MqttClient {
public:
    MqttClient() = default;

    void begin(const char* host, int port, int playerId, CanBus can, IRSender& irTx);
    void loop(int code);
    void publishShoot();
    void publish(const char* topic, const char* payload);
    void onMessage(MqttHandler cb);
    void onReconnecting(std::function<void()> cb);
    void onReconnected(std::function<void()> cb);
    void onDie(std::function<void()> cb);
    void onRespawn(std::function<void()> cb);
    void onGameStart(std::function<void()> cb);
    void onGameEnd(std::function<void()> cb);
    void onRejoin(std::function<void(bool isAlive, uint8_t weaponLevel)> cb);

    GameState& gameState() { return gameState_; }
    int teamId() const { return teamId_; }

private:
    void messageReceived(char* topic, byte* payload, unsigned int len);

    WiFiClient wifi_;
    PubSubClient client_{wifi_};
    MqttHandler handler_;
    std::function<void()> reconnectingCallback_;
    std::function<void()> reconnectedCallback_;
    std::function<void()> dieCallback_;
    std::function<void()> respawnCallback_;
    std::function<void()> gameStartCallback_;
    std::function<void()> gameEndCallback_;
    std::function<void(bool, uint8_t)> rejoinCallback_;
    GameState gameState_;

    int playerId_ = 0;
    String gameId_;
    bool registered_ = false;
    int teamId_ = -1;
    bool isFriendlyFire_ = false;
    GameType gameType_ = GAME_DEATHMATCH;
    CanBus can_;
    IRSender* irTx_ = nullptr;

    unsigned long lastHeartbeat_ = 0;
    const unsigned long heartbeatInterval_ = 6000;
    unsigned long lastReconnectMs_ = 0;
    const unsigned long reconnectCooldown_ = 2000;
};

} // namespace lt