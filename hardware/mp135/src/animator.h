/*
 * Animator — visual effects engine for LaserTag MP135
 *
 * Manages timed visual effects (damage flash, kill popup, death screen,
 * respawn countdown, game start/end sequences).
 * Driven by the main loop tick.
 */

#pragma once

#include "framebuffer.h"
#include <cstdint>
#include <functional>
#include <vector>

namespace lt {

// ══════════════════════════════════════════════════
//  Animation Types
// ══════════════════════════════════════════════════

enum class AnimType : uint8_t {
    DAMAGE_FLASH,      // Red screen flash on hit
    KILL_POPUP,        // "+1 KILL" text popup
    DEATH_SCREEN,      // Full red overlay with DEAD text
    RESPAWN_COUNTDOWN, // 3..2..1 countdown
    GAME_COUNTDOWN,    // Pre-game countdown
    GAME_START_FLASH,  // Green flash on game start
    GAME_END_SCREEN,   // Game over screen
    SCORE_POPUP,       // Score change popup
    LOW_AMMO_BLINK,    // Ammo indicator blinking
};

// ══════════════════════════════════════════════════
//  Animation Instance
// ══════════════════════════════════════════════════

struct Animation {
    AnimType type;
    uint32_t start_ms    = 0;
    uint32_t duration_ms = 0;
    uint16_t value       = 0;   // type-specific: player_id, countdown, score delta
    uint8_t  param       = 0;   // secondary param
    bool     active      = true;
};

// ══════════════════════════════════════════════════
//  Animator Engine
// ══════════════════════════════════════════════════

using AnimDoneCallback = std::function<void(AnimType)>;

class Animator {
public:
    explicit Animator(Framebuffer& fb);

    /// Queue a new animation.
    void start(AnimType type, uint16_t value = 0, uint8_t param = 0);

    /// Tick all active animations (call once per frame).
    /// @param now_ms  current monotonic time in ms
    void tick(uint32_t now_ms);

    /// Render all active animation overlays onto the framebuffer.
    /// Call AFTER the HUD has been rendered.
    void render(uint32_t now_ms);

    /// Whether any blocking animation is active (death screen, etc).
    bool is_blocking() const;

    /// Callback when an animation finishes.
    void on_done(AnimDoneCallback cb) { done_cb_ = std::move(cb); }

private:
    void render_damage_flash(const Animation& a, float t);
    void render_kill_popup(const Animation& a, float t);
    void render_death_screen(const Animation& a, float t);
    void render_respawn_countdown(const Animation& a, float t);
    void render_game_countdown(const Animation& a, float t);
    void render_game_start_flash(const Animation& a, float t);
    void render_game_end_screen(const Animation& a, float t);
    void render_score_popup(const Animation& a, float t);
    void render_low_ammo_blink(const Animation& a, float t);

    uint32_t duration_for(AnimType type) const;
    uint32_t now_ms_ = 0;

    Framebuffer&          fb_;
    std::vector<Animation> anims_;
    AnimDoneCallback       done_cb_;
};

} // namespace lt
