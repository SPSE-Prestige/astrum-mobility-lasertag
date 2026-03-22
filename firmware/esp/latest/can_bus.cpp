#include "can_bus.h"
#include <driver/twai.h>

namespace lt {

bool CanBus::registerPlayer(uint8_t playerId) {
    twai_message_t msg{};
    // CAN ID = 0x00P  (System=0, Register=0, Serie=playerId)
    msg.identifier = (0x0 << 8) | (0x0 << 4) | (playerId & 0xF);
    msg.data_length_code = 0;
    return send
}

bool CanBus::begin(int tx, int rx, int bitrate) {
    // General config: loopback mode = 3
    twai_general_config_t g_config =
        TWAI_GENERAL_CONFIG_DEFAULT((gpio_num_t)tx, (gpio_num_t)rx, (twai_mode_t)3);

    // Timing config
    twai_timing_config_t t_config = TWAI_TIMING_CONFIG_500KBITS();
    if (bitrate == 1000000)
        t_config = TWAI_TIMING_CONFIG_1MBITS();

    // Accept all IDs
    twai_filter_config_t f_config = TWAI_FILTER_CONFIG_ACCEPT_ALL();

    // Install driver
    if (twai_driver_install(&g_config, &t_config, &f_config) != ESP_OK) {
        Serial.println("TWAI driver install failed");
        return false;
    }

    // Start driver
    if (twai_start() != ESP_OK) {
        Serial.println("TWAI driver start failed");
        return false;
    }

    Serial.println("CAN bus initialized (loopback ON)");
    return true;
}

void CanBus::loop() {

    twai_message_t msg;

    if (twai_receive(&msg, 0) == ESP_OK) {

        if (!handler_) return;

        CanFrame f;
        f.id = msg.identifier;
        f.len = msg.data_length_code;
        memcpy(f.data, msg.data, f.len);

        handler_(f);
    }
}

bool CanBus::send(uint32_t id, const uint8_t* data, uint8_t len) {
    twai_message_t msg{};
    msg.identifier = id;
    msg.data_length_code = len;
    if(data && len > 0) {
        memcpy(msg.data, data, len);
    }
    return twai_transmit(&msg, pdMS_TO_TICKS(100)) == ESP_OK;
}

void CanBus::onReceive(CanHandler cb) {
    handler_ = cb;
}

}