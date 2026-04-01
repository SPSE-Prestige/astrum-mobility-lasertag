/*
 * MQTT Client Adapter — libmosquitto wrapper
 *
 * Usage:
 *   #include "mqtt_client.h"
 *   lt::MqttClient mqtt;
 *   mqtt.on_message([](const std::string& topic, const std::string& payload) { ... });
 *   mqtt.connect("192.168.1.100", 1883, "my-client-id");
 *   mqtt.subscribe("device/+/command");
 *   mqtt.publish("device/gun-001/event", R"({"game_id":"abc","victim_id":"gun-002"})");
 */

#pragma once

#include <functional>
#include <string>
#include <vector>

namespace lt {

/// Callback: void(topic, payload)
using MqttHandler = std::function<void(const std::string& topic, const std::string& payload)>;

/// Callback: void(success)
using MqttConnectHandler = std::function<void(bool success)>;

class MqttClient {
public:
    MqttClient();
    ~MqttClient();

    MqttClient(const MqttClient&) = delete;
    MqttClient& operator=(const MqttClient&) = delete;

    /// Connect to MQTT broker (non-blocking, starts internal loop).
    /// @return true if initial connect call succeeds
    bool connect(const std::string& host, int port = 1883,
                 const std::string& client_id = "lasertag-gw",
                 int keepalive = 60);

    /// Disconnect from broker.
    void disconnect();

    /// Publish a message (QoS 1).
    bool publish(const std::string& topic, const std::string& payload, int qos = 1, bool retain = false);

    /// Subscribe to a topic pattern (QoS 1).
    bool subscribe(const std::string& topic, int qos = 1);

    /// Register handler for incoming messages.
    void on_message(MqttHandler handler);

    /// Register handler for connect/reconnect events.
    void on_connect(MqttConnectHandler handler);

    /// Is currently connected?
    bool is_connected() const;

private:
    static void cb_connect(void* mosq, void* obj, int rc);
    static void cb_message(void* mosq, void* obj, const void* msg);
    static void cb_disconnect(void* mosq, void* obj, int rc);

    void*              mosq_ = nullptr;   // struct mosquitto*
    MqttHandler        msg_handler_;
    MqttConnectHandler conn_handler_;
    bool               connected_ = false;
    // Topics to re-subscribe on reconnect
    std::vector<std::pair<std::string, int>> subs_;
};

} // namespace lt
