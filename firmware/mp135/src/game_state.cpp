#include "game_state.h"

#include <algorithm>
#include <cmath>
#include <cstdio>
#include <cstring>

namespace lt {

GameState::GameState(EventBus& bus, Animator& animator, Audio& audio)
    : bus_(bus), animator_(animator), audio_(audio) {
    for (int i = 0; i < MAX_PLAYERS; i++) {
        players_[i].id = static_cast<uint8_t>(i);
    }
}

void GameState::init() {
    bus_.on(EventType::GAME_START, [this](const Event& e) { on_game_start(e); });
    bus_.on(EventType::GAME_END,   [this](const Event& e) { on_game_end(e); });
    bus_.on(EventType::PLAYER_HIT, [this](const Event& e) { on_player_hit(e); });
    bus_.on(EventType::PLAYER_KILL,[this](const Event& e) { on_player_kill(e); });
    bus_.on(EventType::PLAYER_DEATH, [this](const Event& e) { on_player_death(e); });
    bus_.on(EventType::PLAYER_SHOOT, [this](const Event& e) { on_player_shoot(e); });
    bus_.on(EventType::PLAYER_RESPAWN, [this](const Event& e) { on_player_respawn(e); });
    bus_.on(EventType::SCORE_UPDATE, [this](const Event& e) { on_score_update(e); });
    bus_.on(EventType::AMMO_UPDATE, [this](const Event& e) { on_ammo_update(e); });
    bus_.on(EventType::ALIVE_UPDATE, [this](const Event& e) { on_alive_update(e); });
    bus_.on(EventType::KILLCOUNT_UPDATE, [this](const Event& e) { on_killcount_update(e); });
    bus_.on(EventType::HEARTBEAT, [this](const Event& e) { on_heartbeat(e); });
    bus_.on(EventType::DEVICE_REGISTER, [this](const Event& e) { on_device_register(e); });
    bus_.on(EventType::TEAM_ASSIGN, [this](const Event& e) { on_team_assign(e); });
}

// ══════════════════════════════════════════════════
//  Event Handlers
// ══════════════════════════════════════════════════

void GameState::on_game_start(const Event& /*ev*/) {
    std::fprintf(stderr, "[GAME] START\n");
    phase_ = GamePhase::COUNTDOWN;

    // Reset all player stats
    for (auto& p : players_) {
        p.score  = 0;
        p.kills  = 0;
        p.deaths = 0;
        p.ammo   = 30;
        p.alive  = true;
    }

    animator_.start(AnimType::GAME_COUNTDOWN);
    audio_.play(SoundId::GAME_START);

    // After countdown, transition to ACTIVE
    animator_.on_done([this](AnimType type) {
        if (type == AnimType::GAME_COUNTDOWN) {
            phase_ = GamePhase::ACTIVE;
            animator_.start(AnimType::GAME_START_FLASH);
            audio_.play(SoundId::COUNTDOWN_GO);
        }
    });
}

void GameState::on_game_end(const Event& /*ev*/) {
    std::fprintf(stderr, "[GAME] END\n");
    phase_ = GamePhase::ENDING;
    animator_.start(AnimType::GAME_END_SCREEN);
    audio_.play(SoundId::GAME_END);

    animator_.on_done([this](AnimType type) {
        if (type == AnimType::GAME_END_SCREEN) {
            phase_ = GamePhase::IDLE;
        }
    });
}

void GameState::on_player_hit(const Event& ev) {
    uint8_t player_id  = ev.player_id;
    uint8_t attacker   = ev.target_id;

    std::fprintf(stderr, "[COMBAT] Player %u hit by %u\n", player_id, attacker);

    // If local player was hit → damage flash + sound
    if (player_id == local_id_) {
        animator_.start(AnimType::DAMAGE_FLASH);
        audio_.play(SoundId::HIT_RECEIVED);
    }

    // If local player dealt the hit → hit marker sound
    if (attacker == local_id_) {
        audio_.play(SoundId::HIT_DEALT);
    }
}

void GameState::on_player_kill(const Event& ev) {
    uint8_t killer = ev.player_id;
    uint8_t victim = ev.target_id;

    std::fprintf(stderr, "[COMBAT] Player %u killed %u\n", killer, victim);

    if (killer == local_id_) {
        players_[killer].kills++;
        animator_.start(AnimType::KILL_POPUP);
        audio_.play(SoundId::KILL);
    }

    if (victim > 0 && victim < MAX_PLAYERS) {
        players_[victim].deaths++;
    }
}

void GameState::on_player_shoot(const Event& ev) {
    uint8_t player_id = ev.player_id;
    std::fprintf(stderr, "[COMBAT] Player %u shoot\n", player_id);
    audio_.play(SoundId::SHOOT);
}

void GameState::on_player_death(const Event& ev) {
    uint8_t player_id = ev.player_id;

    std::fprintf(stderr, "[COMBAT] Player %u died\n", player_id);

    if (player_id < MAX_PLAYERS) {
        players_[player_id].alive = false;
    }

    if (player_id == local_id_) {
        animator_.start(AnimType::DEATH_SCREEN);
        audio_.play(SoundId::DEATH);
    }
}

void GameState::on_player_respawn(const Event& ev) {
    uint8_t player_id = ev.player_id;

    std::fprintf(stderr, "[GAME] Player %u respawn\n", player_id);

    if (player_id < MAX_PLAYERS) {
        players_[player_id].alive = true;
        players_[player_id].ammo  = 30;
    }

    if (player_id == local_id_) {
        animator_.start(AnimType::RESPAWN_COUNTDOWN);
        audio_.play(SoundId::RESPAWN);
    }
}

void GameState::on_score_update(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        uint16_t old_score = players_[pid].score;
        players_[pid].score = ev.value;

        if (pid == local_id_ && ev.value > old_score) {
            uint16_t delta = ev.value - old_score;
            animator_.start(AnimType::SCORE_POPUP, delta);
        }
    }
}

void GameState::on_ammo_update(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        players_[pid].ammo = ev.value;

        if (pid == local_id_ && ev.value <= 5 && ev.value > 0) {
            animator_.start(AnimType::LOW_AMMO_BLINK);
            audio_.play(SoundId::AMMO_LOW);
        }
        if (pid == local_id_ && ev.value == 0) {
            audio_.play(SoundId::AMMO_EMPTY);
        }
    }
}

void GameState::on_alive_update(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        bool was_alive = players_[pid].alive;
        players_[pid].alive = (ev.value != 0);

        // Trigger death animation if status changed
        if (was_alive && !players_[pid].alive && pid == local_id_) {
            animator_.start(AnimType::DEATH_SCREEN);
            audio_.play(SoundId::DEATH);
        }
    }
}

void GameState::on_killcount_update(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        players_[pid].kills = ev.value;
    }
}

void GameState::on_heartbeat(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        players_[pid].connected = true;
        players_[pid].last_seen = ev.timestamp;
    }
}

void GameState::on_device_register(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        players_[pid].connected = true;
        std::fprintf(stderr, "[SYS] Device %u registered\n", pid);
    }

    // Auto-assign local player from first DEVICE_REGISTER if not yet set
    if (local_id_ == 0 && pid > 0 && pid < MAX_PLAYERS) {
        local_id_ = pid;
        phase_ = GamePhase::IDLE;
        std::fprintf(stderr, "[SYS] Local player assigned: %u\n", pid);
    }
}

void GameState::on_team_assign(const Event& ev) {
    uint8_t pid = ev.player_id;
    if (pid < MAX_PLAYERS) {
        players_[pid].team = ev.target_id;
    }
}

// ══════════════════════════════════════════════════
//  HUD Rendering
// ══════════════════════════════════════════════════

int GameState::connected_count() const {
    int count = 0;
    for (int i = 1; i < MAX_PLAYERS; i++) {
        if (players_[i].connected) count++;
    }
    return count;
}

void GameState::render_hud(Framebuffer& fb, uint32_t now_ms) {
    switch (phase_) {
    case GamePhase::WAITING_FOR_DEVICE:
        render_waiting_screen(fb, now_ms);
        break;
    case GamePhase::IDLE:
        render_idle_screen(fb, now_ms);
        break;
    case GamePhase::COUNTDOWN:
        // Countdown is handled by animator
        break;
    case GamePhase::ACTIVE:
        render_active_hud(fb, now_ms);
        break;
    case GamePhase::ENDING:
        // Game over handled by animator
        break;
    }
}

void GameState::render_waiting_screen(Framebuffer& fb, uint32_t now_ms) {
    fb.clear(Color::BLACK);

    // Title
    fb.draw_text_centered(30, "LASERTAG", 5, Color::CYAN);

    // Pulsing "CONNECTING..." text
    float pulse = 0.5f + 0.5f * std::sin(static_cast<float>(now_ms) / 400.0f);
    uint16_t color = (pulse > 0.5f) ? Color::ORANGE : Color::DARK_GRAY;
    fb.draw_text_centered(110, "CONNECTING", 3, color);
    fb.draw_text_centered(145, "TO DEVICE...", 2, Color::GRAY);

    // CAN bus status hint
    fb.draw_text_centered(200, "WAITING FOR ESP", 2, Color::YELLOW);

    // Border
    float border_pulse = 0.5f + 0.5f * std::sin(static_cast<float>(now_ms) / 600.0f);
    uint16_t border_color = (border_pulse > 0.5f) ? Color::ORANGE : Color::DARK_GRAY;
    fb.rect_outline(0, 0, fb.width(), fb.height(), border_color, 2);
}

void GameState::render_idle_screen(Framebuffer& fb, uint32_t now_ms) {
    fb.clear(Color::BLACK);

    // Title
    fb.draw_text_centered(30, "LASERTAG", 5, Color::CYAN);

    // Pulsing "WAITING..." text
    float pulse = 0.5f + 0.5f * std::sin(static_cast<float>(now_ms) / 500.0f);
    uint16_t wait_color = (pulse > 0.5f) ? Color::WHITE : Color::GRAY;
    fb.draw_text_centered(120, "WAITING...", 2, wait_color);

    // Player info
    char buf[32];
    std::snprintf(buf, sizeof(buf), "PLAYER %u", local_id_);
    fb.draw_text_centered(170, buf, 3, Color::YELLOW);

    // Connected players count
    std::snprintf(buf, sizeof(buf), "ONLINE: %d", connected_count());
    fb.draw_text_centered(220, buf, 2, Color::GREEN);

    // Border
    fb.rect_outline(0, 0, fb.width(), fb.height(), Color::CYAN, 2);
}

void GameState::render_active_hud(Framebuffer& fb, uint32_t /*now_ms*/) {
    fb.clear(Color::BLACK);

    const auto& lp = local();
    int w = fb.width();
    int h = fb.height();
    char buf[32];

    // ── Device ID (top left, small) ──
    std::snprintf(buf, sizeof(buf), "ID:%u", local_id_);
    fb.draw_text(8, 8, buf, 2, Color::GRAY);

    // ── Score (big, centered) ──
    std::snprintf(buf, sizeof(buf), "%u", lp.score);
    fb.draw_text_centered(50, buf, 10, Color::YELLOW);

    // ── "SCORE" label ──
    fb.draw_text_centered(50 + 10 * 7 + 10, "SCORE", 2, Color::GRAY);

    // ── Kills (below score) ──
    std::snprintf(buf, sizeof(buf), "KILLS: %u", lp.kills);
    fb.draw_text_centered(h - 60, buf, 4, Color::WHITE);

    // ── Border ──
    uint16_t border_color = lp.alive ? Color::GREEN : Color::RED;
    fb.rect_outline(0, 0, w, h, border_color, 1);
}

void GameState::render_game_over(Framebuffer& fb, uint32_t /*now_ms*/) {
    fb.clear(Color::BLACK);
    // Handled by animator GAME_END_SCREEN
    (void)fb;
}

} // namespace lt
