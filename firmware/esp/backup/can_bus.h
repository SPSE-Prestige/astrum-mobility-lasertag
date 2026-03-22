/*
 * CAN Bus Adapter — SocketCAN (Linux)
 *
 * Classic CAN, 1M bitrate, max 8 bytes per frame.
 * Usage:
 *   #include "can_bus.h"
 *   lt::CanBus bus;
 *   bus.open("can1");            // 1M bitrate
 *   bus.send(0x100, data, len);
 *   bus.set_handler([](const lt::CanFrame& f) { ... });
 *   bus.start();                 // background reader thread
 */

#pragma once

#include <cstdint>
#include <functional>
#include <string>
#include <vector>

namespace lt {

/// Maximum classic CAN payload.
constexpr int CAN_MAX_DLEN = 8;

/// Received CAN frame.
struct CanFrame {
    uint32_t id;
    uint8_t  data[CAN_MAX_DLEN];
    uint8_t  len;        ///< actual data length (0–8)
};

/// Callback signature: void(const CanFrame&)
using CanHandler = std::function<void(const CanFrame&)>;

class CanBus {
public:
    CanBus() = default;
    ~CanBus();

    CanBus(const CanBus&) = delete;
    CanBus& operator=(const CanBus&) = delete;

    /// Configure and open a SocketCAN interface.
    /// Runs `ip link set` to configure classic CAN, 1M, restart-ms 100.
    /// @param iface  e.g. "can0", "can1"
    /// @param bitrate  bitrate (default 1 000 000)
    /// @return true on success
    bool open(const std::string& iface, int bitrate = 1000000);

    /// Close socket and stop reader thread.
    void close();

    /// Send a CAN frame.
    /// @param id  CAN arbitration ID (11-bit)
    /// @param data  payload bytes
    /// @param len  number of bytes (max 8)
    /// @return true on success
    bool send(uint32_t id, const uint8_t* data, uint8_t len);

    /// Send a CAN frame with vector payload.
    bool send(uint32_t id, const std::vector<uint8_t>& data) {
        return send(id, data.data(), static_cast<uint8_t>(data.size()));
    }

    /// Register a handler called for every received frame (from reader thread).
    void set_handler(CanHandler handler);

    /// Start the background reader thread.
    void start();

    /// Stop the background reader thread.
    void stop();

    /// Is the bus open?
    bool is_open() const;

private:
    void reader_loop();

    int         sock_    = -1;
    std::string iface_;
    CanHandler  handler_;
    bool        running_ = false;
    // Reader thread handle stored as opaque pointer to avoid <thread> in header
    void*       thread_  = nullptr;
};

} // namespace lt
