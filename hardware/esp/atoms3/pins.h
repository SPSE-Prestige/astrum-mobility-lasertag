#pragma once

// Set to 1 to enable debug logging over Serial
#define DEBUG 1

// Atom S3 (ESP32-S3) pin definitions

// Port A — IR TX and RX on the same port
#define PIN_IR_TX   2   // Port A SDA
#define PIN_IR_RX   1   // Port A SCL

// External shoot button (active LOW)
#define PIN_BTN_SHOOT  41

// CAN bus (external transceiver e.g. TJA1050)
#define PIN_CAN_TX  5
#define PIN_CAN_RX  6
