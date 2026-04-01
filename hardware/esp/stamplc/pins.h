#pragma once

// Set to 1 to enable debug logging over Serial
#define DEBUG 1

// STAMPLC (StampS3A, ESP32-S3FN8) pin definitions

// Port A — IR transmitter (shoot / hit signal)
#define PIN_IR_TX       2   // Port A SDA — hit TX
#define PIN_HIT_ENABLE  1   // Port A SCL — hit enable/control

// Port C — IR receivers (input only)
#define PIN_IR_RX_1     5   // Port C
#define PIN_IR_RX_2     4   // Port C

// CAN bus
#define PIN_CAN_TX      42
#define PIN_CAN_RX      43

// I2C IO expander (Button A, relays)
#define PIN_I2C_SDA     13
#define PIN_I2C_SCL     15

// External shoot button (active LOW)
#define PIN_BTN_SHOOT  41