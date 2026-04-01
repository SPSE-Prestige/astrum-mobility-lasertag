#include "mqtt_client.h"

#include <cstdio>
#include <cstring>
#include <mosquitto.h>

namespace lt {

// ── Init / cleanup ──

static bool g_mosquitto_init = false;

MqttClient::MqttClient() {
    if (!g_mosquitto_init) {
        mosquitto_lib_init();
        g_mosquitto_init = true;
    }
}

MqttClient::~MqttClient() {
    disconnect();
}

// ── Connect / Disconnect ──

bool MqttClient::connect(const std::string& host, int port,
                          const std::string& client_id, int keepalive) {
    auto* m = mosquitto_new(client_id.c_str(), true, this);
    if (!m) {
        std::fprintf(stderr, "[MQTT] mosquitto_new failed\n");
        return false;
    }
    mosq_ = m;

    // Automatic reconnect
    mosquitto_reconnect_delay_set(m, 1, 10, true);

    mosquitto_connect_callback_set(m,
        [](struct mosquitto*, void* obj, int rc) {
            auto* self = static_cast<MqttClient*>(obj);
            self->connected_ = (rc == 0);
            if (rc == 0) {
                std::fprintf(stderr, "[MQTT] connected\n");
                // Re-subscribe after reconnect
                auto* mosq = static_cast<struct mosquitto*>(self->mosq_);
                for (auto& [topic, qos] : self->subs_)
                    mosquitto_subscribe(mosq, nullptr, topic.c_str(), qos);
            } else {
                std::fprintf(stderr, "[MQTT] connect failed: rc=%d\n", rc);
            }
            if (self->conn_handler_)
                self->conn_handler_(rc == 0);
        });

    mosquitto_message_callback_set(m,
        [](struct mosquitto*, void* obj, const struct mosquitto_message* msg) {
            auto* self = static_cast<MqttClient*>(obj);
            if (!self->msg_handler_ || !msg->topic) return;
            std::string topic(msg->topic);
            std::string payload(static_cast<char*>(msg->payload),
                                static_cast<size_t>(msg->payloadlen));
            self->msg_handler_(topic, payload);
        });

    mosquitto_disconnect_callback_set(m,
        [](struct mosquitto*, void* obj, int rc) {
            auto* self = static_cast<MqttClient*>(obj);
            self->connected_ = false;
            std::fprintf(stderr, "[MQTT] disconnected (rc=%d)\n", rc);
        });

    int rc = mosquitto_connect(m, host.c_str(), port, keepalive);
    if (rc != MOSQ_ERR_SUCCESS) {
        std::fprintf(stderr, "[MQTT] connect error: %s\n", mosquitto_strerror(rc));
        mosquitto_destroy(m);
        mosq_ = nullptr;
        return false;
    }

    mosquitto_loop_start(m);
    std::fprintf(stderr, "[MQTT] connecting to %s:%d ...\n", host.c_str(), port);
    return true;
}

void MqttClient::disconnect() {
    if (!mosq_) return;
    auto* m = static_cast<struct mosquitto*>(mosq_);
    mosquitto_loop_stop(m, true);
    mosquitto_disconnect(m);
    mosquitto_destroy(m);
    mosq_ = nullptr;
    connected_ = false;
}

// ── Publish / Subscribe ──

bool MqttClient::publish(const std::string& topic, const std::string& payload,
                          int qos, bool retain) {
    if (!mosq_) return false;
    auto* m = static_cast<struct mosquitto*>(mosq_);
    int rc = mosquitto_publish(m, nullptr, topic.c_str(),
                               static_cast<int>(payload.size()),
                               payload.c_str(), qos, retain);
    return rc == MOSQ_ERR_SUCCESS;
}

bool MqttClient::subscribe(const std::string& topic, int qos) {
    subs_.emplace_back(topic, qos);
    if (!mosq_) return false;
    auto* m = static_cast<struct mosquitto*>(mosq_);
    int rc = mosquitto_subscribe(m, nullptr, topic.c_str(), qos);
    return rc == MOSQ_ERR_SUCCESS;
}

// ── Handlers ──

void MqttClient::on_message(MqttHandler handler) {
    msg_handler_ = std::move(handler);
}

void MqttClient::on_connect(MqttConnectHandler handler) {
    conn_handler_ = std::move(handler);
}

bool MqttClient::is_connected() const {
    return connected_;
}

} // namespace lt
