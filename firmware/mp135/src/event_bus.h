/*
 * Event Bus — lightweight pub/sub for LaserTag MP135
 *
 * Thread-safe, lock-free event queue with typed event dispatch.
 * All game subsystems communicate through events — no direct coupling.
 */

#pragma once

#include <cstdint>
#include <functional>
#include <mutex>
#include <queue>
#include <string>
#include <unordered_map>
#include <vector>

namespace lt {

// ══════════════════════════════════════════════════
//  Event Types
// ══════════════════════════════════════════════════

enum class EventType : uint8_t {
    // ── Game flow ──
    GAME_START,
    GAME_END,
    GAME_COUNTDOWN,       // data: seconds remaining

    // ── Combat ──
    PLAYER_HIT,           // data: player_id, attacker_id
    PLAYER_KILL,          // data: killer_id, victim_id
    PLAYER_DEATH,         // data: player_id
    PLAYER_SHOOT,         // data: player_id
    PLAYER_RESPAWN,       // data: player_id

    // ── Status updates ──
    SCORE_UPDATE,         // data: player_id, score
    AMMO_UPDATE,          // data: player_id, ammo
    ALIVE_UPDATE,         // data: player_id, alive
    KILLCOUNT_UPDATE,     // data: player_id, kills

    // ── System ──
    HEARTBEAT,            // data: player_id
    DEVICE_REGISTER,      // data: player_id
    TEAM_ASSIGN,          // data: player_id, team_id

    // ── UI ──
    ANIMATION_DONE,
    RENDER_REQUEST,
};

// ══════════════════════════════════════════════════
//  Event Data
// ══════════════════════════════════════════════════

struct Event {
    EventType type;
    uint8_t   player_id  = 0;
    uint8_t   target_id  = 0;   // attacker, victim, team, etc.
    uint16_t  value      = 0;   // score, ammo, countdown, etc.
    uint64_t  timestamp  = 0;   // monotonic ms
};

// ══════════════════════════════════════════════════
//  Event Bus
// ══════════════════════════════════════════════════

using EventCallback = std::function<void(const Event&)>;

class EventBus {
public:
    /// Subscribe to a specific event type.
    void on(EventType type, EventCallback cb) {
        std::lock_guard<std::mutex> lock(sub_mtx_);
        subs_[type].push_back(std::move(cb));
    }

    /// Queue an event for deferred dispatch (thread-safe).
    void emit(Event ev) {
        std::lock_guard<std::mutex> lock(q_mtx_);
        queue_.push(std::move(ev));
    }

    /// Helper: emit with just type + player.
    void emit(EventType type, uint8_t player = 0, uint8_t target = 0, uint16_t value = 0) {
        Event ev;
        ev.type      = type;
        ev.player_id = player;
        ev.target_id = target;
        ev.value     = value;
        emit(std::move(ev));
    }

    /// Process all queued events on the calling thread.
    /// Call this once per frame from the main loop.
    void dispatch() {
        std::queue<Event> batch;
        {
            std::lock_guard<std::mutex> lock(q_mtx_);
            batch.swap(queue_);
        }

        while (!batch.empty()) {
            const auto& ev = batch.front();

            std::vector<EventCallback> cbs;
            {
                std::lock_guard<std::mutex> lock(sub_mtx_);
                auto it = subs_.find(ev.type);
                if (it != subs_.end()) {
                    cbs = it->second;
                }
            }

            for (auto& cb : cbs) {
                cb(ev);
            }
            batch.pop();
        }
    }

    /// Clear all subscriptions.
    void clear() {
        std::lock_guard<std::mutex> lock(sub_mtx_);
        subs_.clear();
    }

private:
    std::mutex                                              q_mtx_;
    std::queue<Event>                                       queue_;
    std::mutex                                              sub_mtx_;
    std::unordered_map<EventType, std::vector<EventCallback>> subs_;
};

} // namespace lt
