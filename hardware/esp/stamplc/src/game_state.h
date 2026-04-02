#pragma once
#include <Arduino.h>
#include <functional>

namespace lt {

class GameState {
public:
    void reset() {
        kills_ = 0;
        level_ = 0;
    }

    // Call on each confirmed kill (server decides level-up separately).
    void addKill() {
        kills_++;
        Serial.printf("[GAME] Kill #%d  level=%d\n", kills_, level_);
        if (killCallback_) killCallback_(kills_, level_);
    }

    // Call when server sends action == "upgrade".
    void upgrade() {
        level_++;
        Serial.printf("[GAME] LEVEL UP → %d\n", level_);
        if (levelUpCallback_) levelUpCallback_(level_);
    }

    uint8_t kills() const { return kills_; }
    uint8_t level() const { return level_; }

    void onKill(std::function<void(uint8_t kills, uint8_t level)> cb) {
        killCallback_ = cb;
    }
    void onLevelUp(std::function<void(uint8_t level)> cb) {
        levelUpCallback_ = cb;
    }

private:
    uint8_t kills_ = 0;
    uint8_t level_ = 0;
    std::function<void(uint8_t, uint8_t)> killCallback_;
    std::function<void(uint8_t)> levelUpCallback_;
};

} // namespace lt
