#pragma once
#include <WiFi.h>
#include <PubSubClient.h>
#include <functional>

namespace lt {

using MqttHandler = std::function<void(const String&, const String&)>;

class MqttClient {
public:
    void begin(const char* host, int port);
    void loop();
    void publish(const char* topic, const char* payload);
    void onMessage(MqttHandler cb);

private:
    WiFiClient wifi_;
    PubSubClient client_{wifi_};
    MqttHandler handler_;
};

}