#pragma once
#include <WiFi.h>
#include <PubSubClient.h>
#include <functional>
#include "can_bus.h"

namespace lt {

using MqttHandler = std::function<void(const String&, const String&)>;

class MqttClient {
public:
    MqttClient() = default;

    void begin(const char* host, int port, int playerId, CanBus can);
    void loop(int code);
    void publish(const char* topic, const char* payload);
    void onMessage(MqttHandler cb);

private:
    void messageReceived(char* topic, byte* payload, unsigned int len);

    WiFiClient wifi_;
    PubSubClient client_{wifi_};
    MqttHandler handler_;

    String deviceId_;
    unsigned long lastHeartbeat_ = 0;
    const unsigned long heartbeatInterval_ = 10000; // 10 seconds
};

} // namespace lt