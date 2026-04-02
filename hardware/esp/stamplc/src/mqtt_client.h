#pragma once
#include <WiFi.h>
#include <PubSubClient.h>
#include <functional>
#include "can_bus.h"
#include "game_state.h"
#include "ir_controller.h"

namespace lt {

using MqttHandler = std::function<void(const String&, const String&)>;

class MqttClient {
public:
    MqttClient() = default;

    void begin(const char* host, int port, int playerId, CanBus can, IRSender& irTx);
    void loop(int code);
    void publish(const char* topic, const char* payload);
    void onMessage(MqttHandler cb);
    void onReconnecting(std::function<void()> cb);
    void onReconnected(std::function<void()> cb);
    void onDie(std::function<void()> cb);
    void onRespawn(std::function<void()> cb);
    void onGameStart(std::function<void()> cb);
    void onGameEnd(std::function<void()> cb);

    GameState& gameState() { return gameState_; }

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
    GameState gameState_;

    int playerId_ = 0;
    String gameId_;
    bool registered_ = false;
    CanBus can_;
    IRSender* irTx_ = nullptr;

    unsigned long lastHeartbeat_ = 0;
    const unsigned long heartbeatInterval_ = 6000;
    unsigned long lastReconnectMs_ = 0;
    const unsigned long reconnectCooldown_ = 2000;
};

} // namespace lt