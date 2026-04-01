/*
 * Game State — state machine + player tracking for LaserTag MP135
 *
 * Manages the complete game lifecycle and all player state.
 * Listens to EventBus events and triggers animations + sounds.
 */

#pragma once

#include "animator.h"
#include "audio.h"
#include "event_bus.h"
#include "framebuffer.h"

#include <array>
#include <cstdint>

namespace lt {

// ══════════════════════════════════════════════════
//  Game Phase (state machine)
// ══════════════════════════════════════════════════

enum class GamePhase : uint8_t {
    WAITING_FOR_DEVICE, // No device ID yet — waiting for ESP CAN registration
    IDLE,               // Device registered, waiting for game start
    COUNTDOWN,          // Pre-game countdown
    ACTIVE,             // Game in progress
    ENDING,             // Game over animation
};

// ══════════════════════════════════════════════════
//  Player Data
// ══════════════════════════════════════════════════

struct Player {
    uint8_t  id        = 0;
    uint8_t  team      = 0;
    uint16_t score     = 0;
    uint16_t kills     = 0;
    uint16_t deaths    = 0;
    uint16_t ammo      = 30;
    bool     alive     = true;
    bool     connected = false;
    uint32_t last_seen = 0;   // heartbeat timestamp
};

constexpr int MAX_PLAYERS = 16;  // 0=broadcast, 1-15=players

// ══════════════════════════════════════════════════
//  Game State Manager
// ══════════════════════════════════════════════════

class GameState {
public:
    GameState(EventBus& bus, Animator& animator, Audio& audio);

    /// Wire up all event handlers.
    void init();

    /// Get current game phase.
    GamePhase phase() const { return phase_; }

    /// Get player data (0-15). Player 0 = local/broadcast.
    const Player& player(uint8_t id) const { return players_[id & 0xF]; }

    /// Set the local player ID (0 = unassigned, waiting for CAN).
    void set_local_player(uint8_t id) { local_id_ = id; if (id > 0) phase_ = GamePhase::IDLE; }
    uint8_t local_player() const { return local_id_; }

    /// Whether we have a device ID assigned.
    bool has_device() const { return local_id_ > 0; }

    /// Get the local player ref.
    const Player& local() const { return players_[local_id_]; }

    /// Render the HUD for the current state.
    void render_hud(Framebuffer& fb, uint32_t now_ms);

    /// Get list of connected player IDs (for leaderboard).
    int connected_count() const;

private:
    // Event handlers
    void on_game_start(const Event& ev);
    void on_game_end(const Event& ev);
    void on_player_hit(const Event& ev);
    void on_player_kill(const Event& ev);
    void on_player_death(const Event& ev);
    void on_player_shoot(const Event& ev);
    void on_player_respawn(const Event& ev);
    void on_score_update(const Event& ev);
    void on_ammo_update(const Event& ev);
    void on_alive_update(const Event& ev);
    void on_killcount_update(const Event& ev);
    void on_heartbeat(const Event& ev);
    void on_device_register(const Event& ev);
    void on_team_assign(const Event& ev);

    // HUD rendering helpers
    void render_waiting_screen(Framebuffer& fb, uint32_t now_ms);
    void render_idle_screen(Framebuffer& fb, uint32_t now_ms);
    void render_active_hud(Framebuffer& fb, uint32_t now_ms);
    void render_game_over(Framebuffer& fb, uint32_t now_ms);

    EventBus&  bus_;
    Animator&  animator_;
    Audio&     audio_;
    GamePhase  phase_ = GamePhase::WAITING_FOR_DEVICE;
    uint8_t    local_id_ = 0;  // 0 = unassigned, set by ESP via CAN

    std::array<Player, MAX_PLAYERS> players_{};
};

} // namespace lt
