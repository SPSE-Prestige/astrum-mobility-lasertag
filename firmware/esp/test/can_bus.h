#pragma once
#include <Arduino.h>
#include <functional>

namespace lt {

struct CanFrame {
    uint32_t id;
    uint8_t len;
    uint8_t data[8];
};

using CanHandler = std::function<void(const CanFrame&)>;

class CanBus {
public:
    bool begin(int tx, int rx, int bitrate = 500000);
    void loop();
    bool send(uint32_t id, const uint8_t* data, uint8_t len);
    void onReceive(CanHandler cb);

private:
    CanHandler handler_;
};

}