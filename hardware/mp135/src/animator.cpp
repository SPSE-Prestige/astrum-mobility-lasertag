#include "animator.h"

#include <algorithm>
#include <cmath>
#include <cstdio>
#include <cstring>

namespace lt {

Animator::Animator(Framebuffer& fb) : fb_(fb) {}

uint32_t Animator::duration_for(AnimType type) const {
    switch (type) {
    case AnimType::DAMAGE_FLASH:      return 400;
    case AnimType::KILL_POPUP:        return 1500;
    case AnimType::DEATH_SCREEN:      return 3000;
    case AnimType::RESPAWN_COUNTDOWN: return 3000;
    case AnimType::GAME_COUNTDOWN:    return 5000;
    case AnimType::GAME_START_FLASH:  return 800;
    case AnimType::GAME_END_SCREEN:   return 5000;
    case AnimType::SCORE_POPUP:       return 1200;
    case AnimType::LOW_AMMO_BLINK:    return 2000;
    }
    return 1000;
}

void Animator::start(AnimType type, uint16_t value, uint8_t param) {
    // For some types, replace existing instead of stacking
    if (type == AnimType::DEATH_SCREEN || type == AnimType::GAME_END_SCREEN ||
        type == AnimType::GAME_COUNTDOWN || type == AnimType::GAME_START_FLASH) {
        anims_.erase(std::remove_if(anims_.begin(), anims_.end(),
            [type](const Animation& a) { return a.type == type; }),
            anims_.end());
    }

    Animation a;
    a.type        = type;
    a.start_ms    = now_ms_;
    a.duration_ms = duration_for(type);
    a.value       = value;
    a.param       = param;
    a.active      = true;
    anims_.push_back(a);
}

void Animator::tick(uint32_t now_ms) {
    now_ms_ = now_ms;

    // Collect completed animation types BEFORE modifying the vector,
    // because done_cb_ may call start() which invalidates iterators.
    std::vector<AnimType> completed;

    for (auto it = anims_.begin(); it != anims_.end();) {
        if (it->active && (now_ms - it->start_ms) >= it->duration_ms) {
            it->active = false;
            completed.push_back(it->type);
            it = anims_.erase(it);
        } else {
            ++it;
        }
    }

    // Fire callbacks after iteration is complete (safe to modify anims_ now)
    if (done_cb_) {
        for (auto type : completed) {
            done_cb_(type);
        }
    }
}

void Animator::render(uint32_t now_ms) {
    for (const auto& a : anims_) {
        if (!a.active) continue;

        uint32_t elapsed = now_ms - a.start_ms;
        float t = static_cast<float>(elapsed) / a.duration_ms;
        t = std::clamp(t, 0.0f, 1.0f);

        switch (a.type) {
        case AnimType::DAMAGE_FLASH:      render_damage_flash(a, t);      break;
        case AnimType::KILL_POPUP:        render_kill_popup(a, t);        break;
        case AnimType::DEATH_SCREEN:      render_death_screen(a, t);      break;
        case AnimType::RESPAWN_COUNTDOWN: render_respawn_countdown(a, t); break;
        case AnimType::GAME_COUNTDOWN:    render_game_countdown(a, t);    break;
        case AnimType::GAME_START_FLASH:  render_game_start_flash(a, t);  break;
        case AnimType::GAME_END_SCREEN:   render_game_end_screen(a, t);   break;
        case AnimType::SCORE_POPUP:       render_score_popup(a, t);       break;
        case AnimType::LOW_AMMO_BLINK:    render_low_ammo_blink(a, t);    break;
        }
    }
}

bool Animator::is_blocking() const {
    for (const auto& a : anims_) {
        if (a.active && (a.type == AnimType::DEATH_SCREEN ||
                         a.type == AnimType::GAME_END_SCREEN ||
                         a.type == AnimType::GAME_COUNTDOWN)) {
            return true;
        }
    }
    return false;
}

// ══════════════════════════════════════════════════
//  Individual animation renderers
// ══════════════════════════════════════════════════

void Animator::render_damage_flash(const Animation& /*a*/, float t) {
    // Fast red flash that fades out
    // t=0 → full red, t=1 → transparent
    float fade = 1.0f - t;
    // Ease out: quadratic
    fade = fade * fade;

    uint8_t intensity = static_cast<uint8_t>(fade * 180);
    fb_.tint(Color::RED, intensity);

    // Draw red vignette border
    int border = static_cast<int>(fade * 30);
    if (border > 0) {
        int w = fb_.width(), h = fb_.height();
        fb_.rect(0, 0, w, border, Color::DARK_RED);
        fb_.rect(0, h - border, w, border, Color::DARK_RED);
        fb_.rect(0, 0, border, h, Color::DARK_RED);
        fb_.rect(w - border, 0, border, h, Color::DARK_RED);
    }
}

void Animator::render_kill_popup(const Animation& /*a*/, float t) {
    // "+1 KILL" rises up and fades out
    int y_base = fb_.height() / 3;
    int y_offset = static_cast<int>(t * 60);
    int y = y_base - y_offset;

    // Fade: visible [0..0.7], fade out [0.7..1.0]
    if (t > 0.7f) {
        float fade_t = (t - 0.7f) / 0.3f;
        uint8_t alpha = static_cast<uint8_t>((1.0f - fade_t) * 255);
        if (alpha < 30) return;
    }

    fb_.draw_text_centered(y, "+1 KILL", 4, Color::YELLOW);
}

void Animator::render_death_screen(const Animation& /*a*/, float t) {
    // Full-screen red overlay with pulsing "DEAD" text
    uint8_t red_intensity = 120;
    fb_.tint(Color::RED, red_intensity);

    // Pulsing text
    float pulse = 0.5f + 0.5f * std::sin(t * 3.14159f * 4.0f);
    uint16_t text_color = (pulse > 0.5f) ? Color::WHITE : Color::RED;

    int scale = 8;
    int y = (fb_.height() - 7 * scale) / 2;
    fb_.draw_text_centered(y, "DEAD", scale, text_color);

    // Crosshair skull icon (simple X pattern)
    int cx = fb_.width() / 2;
    int cy = y + 7 * scale + 30;
    int sz = 20;
    for (int i = -sz; i <= sz; i++) {
        fb_.pixel(cx + i, cy + i, Color::WHITE);
        fb_.pixel(cx + i, cy - i, Color::WHITE);
        fb_.pixel(cx + i + 1, cy + i, Color::WHITE);
        fb_.pixel(cx + i + 1, cy - i, Color::WHITE);
    }
}

void Animator::render_respawn_countdown(const Animation& /*a*/, float t) {
    // 3..2..1 countdown with expanding numbers
    int sec_remaining = 3 - static_cast<int>(t * 3.0f);
    if (sec_remaining < 1) sec_remaining = 1;

    // Per-second animation phase
    float phase = std::fmod(t * 3.0f, 1.0f);

    // Number grows slightly then shrinks
    float scale_mod = 1.0f + 0.3f * std::sin(phase * 3.14159f);
    int scale = static_cast<int>(10 * scale_mod);

    char buf[12];
    std::snprintf(buf, sizeof(buf), "%d", sec_remaining);

    int y = (fb_.height() - 7 * scale) / 2;
    fb_.draw_text_centered(y, buf, scale, Color::CYAN);

    // "RESPAWN" label below
    fb_.draw_text_centered(y + 7 * scale + 20, "RESPAWN", 3, Color::WHITE);
}

void Animator::render_game_countdown(const Animation& /*a*/, float t) {
    // Full-screen countdown: 5..4..3..2..1..GO!
    fb_.clear(Color::BLACK);

    int total_sec = 5;
    int sec_remaining = total_sec - static_cast<int>(t * total_sec);
    if (sec_remaining < 0) sec_remaining = 0;

    float phase = std::fmod(t * total_sec, 1.0f);
    float scale_mod = 1.0f + 0.5f * (1.0f - phase);
    int scale = static_cast<int>(12 * scale_mod);

    char buf[12];
    if (sec_remaining > 0) {
        std::snprintf(buf, sizeof(buf), "%d", sec_remaining);
    } else {
        std::snprintf(buf, sizeof(buf), "GO");
    }

    uint16_t color = sec_remaining <= 1 ? Color::GREEN :
                     sec_remaining <= 3 ? Color::YELLOW : Color::WHITE;

    int y = (fb_.height() - 7 * scale) / 2;
    fb_.draw_text_centered(y, buf, scale, color);

    // Progress bar at bottom
    float progress = t;
    fb_.draw_bar(20, fb_.height() - 40, fb_.width() - 40, 20, progress,
                 Color::GREEN, Color::DARK_GRAY);
}

void Animator::render_game_start_flash(const Animation& /*a*/, float t) {
    // Bright green flash
    float fade = 1.0f - t;
    fade = fade * fade;
    uint8_t intensity = static_cast<uint8_t>(fade * 200);
    fb_.tint(Color::GREEN, intensity);
}

void Animator::render_game_end_screen(const Animation& /*a*/, float t) {
    fb_.clear(Color::BLACK);

    // "GAME OVER" with dramatic reveal
    float reveal = std::min(t * 3.0f, 1.0f);

    if (reveal >= 1.0f) {
        // Full text visible
        int scale = 6;
        int y = (fb_.height() - 7 * scale) / 2 - 30;
        fb_.draw_text_centered(y, "GAME", scale, Color::RED);
        fb_.draw_text_centered(y + 7 * scale + 10, "OVER", scale, Color::RED);
    } else {
        // Expanding horizontal line reveal
        int line_y = fb_.height() / 2;
        int line_w = static_cast<int>(fb_.width() * reveal);
        int line_x = (fb_.width() - line_w) / 2;
        fb_.rect(line_x, line_y - 2, line_w, 4, Color::RED);
    }

    // Subtle pulsing border
    float pulse = 0.5f + 0.5f * std::sin(t * 6.28f);
    uint16_t border_color = (pulse > 0.5f) ? Color::DARK_RED : Color::BLACK;
    fb_.rect_outline(0, 0, fb_.width(), fb_.height(), border_color, 3);
}

void Animator::render_score_popup(const Animation& a, float t) {
    // "+XX" score popup floating up
    int delta = static_cast<int>(a.value);
    int y_base = fb_.height() / 4;
    int y = y_base - static_cast<int>(t * 40);

    if (t > 0.8f) return; // fade away

    char buf[16];
    std::snprintf(buf, sizeof(buf), "+%d", delta);
    fb_.draw_text_centered(y, buf, 3, Color::GREEN);
}

void Animator::render_low_ammo_blink(const Animation& /*a*/, float t) {
    // Blinking orange border
    float blink = std::sin(t * 3.14159f * 8.0f);
    if (blink > 0) {
        fb_.rect_outline(0, 0, fb_.width(), fb_.height(), Color::ORANGE, 3);
    }
}

} // namespace lt
