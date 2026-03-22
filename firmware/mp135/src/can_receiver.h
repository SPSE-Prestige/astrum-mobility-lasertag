/*
 * CAN Receiver — SocketCAN listener for LaserTag MP135
 *
 * Runs a background thread reading CAN frames from can1,
 * decodes them via can_protocol.h and emits typed events
 * onto the EventBus.
 */

#pragma once

#include "event_bus.h"

#include <atomic>
#include <cstdint>
#include <string>
#include <thread>

namespace lt {

class CanReceiver {
public:
    explicit CanReceiver(EventBus& bus, const std::string& iface = "can1");
    ~CanReceiver();

    CanReceiver(const CanReceiver&) = delete;
    CanReceiver& operator=(const CanReceiver&) = delete;

    bool open(int bitrate = 1000000);
    void start();
    void stop();
    bool send(uint32_t id, const uint8_t* data, uint8_t len);

private:
    void reader_loop();
    void decode_frame(uint32_t id, const uint8_t* data, uint8_t len);

    EventBus&        bus_;
    std::string      iface_;
    int              sock_ = -1;
    std::atomic<bool> running_{false};
    std::thread      thread_;
};

} // namespace lt
