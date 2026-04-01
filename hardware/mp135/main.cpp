/*
 * LaserTag MP135 — Main Application
 *
 * Event-driven game display for M5Stack MP135 with:
 *   - CAN bus listener on can1 (1Mbps)
 *   - RGB565 framebuffer rendering on /dev/fb1
 *   - ALSA audio playback
 *   - Animation engine with damage flash, kill popup, death screen
 *
 * Architecture:
 *   CAN frames → EventBus → GameState → Animator + Audio → Framebuffer
 *
 * Usage:
 *   sudo ./lasertag [player_id] [can_iface]
 *   sudo ./lasertag 3 can1
 */

#include "src/animator.h"
#include "src/audio.h"
#include "src/can_receiver.h"
#include "src/event_bus.h"
#include "src/framebuffer.h"
#include "src/game_state.h"

#include <chrono>
#include <csignal>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <string>
#include <thread>

// ── Global shutdown flag ──
static volatile sig_atomic_t g_running = 1;

static void signal_handler(int /*sig*/) {
    g_running = 0;
}

// ── Monotonic clock in milliseconds ──
static uint32_t now_ms() {
    auto tp = std::chrono::steady_clock::now();
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(tp.time_since_epoch());
    return static_cast<uint32_t>(ms.count());
}

// ══════════════════════════════════════════════════
//  Main
// ══════════════════════════════════════════════════

int main(int argc, char** argv) {
    // ── Parse args ──
    uint8_t player_id = 0;  // 0 = auto-detect from ESP via CAN bus
    const char* can_iface = "can1";
    const char* fb_path   = "/dev/fb0";
    const char* media_dir = "media";

    for (int i = 1; i < argc; i++) {
        if (std::strcmp(argv[i], "--help") == 0 || std::strcmp(argv[i], "-h") == 0) {
            std::printf("Usage: %s [OPTIONS]\n\n", argv[0]);
            std::printf("Options:\n");
            std::printf("  -p, --player ID    Local player ID (1-15, default: auto from CAN)\n");
            std::printf("  -c, --can IFACE    CAN interface (default: can1)\n");
            std::printf("  -f, --fb PATH      Framebuffer device (default: /dev/fb1)\n");
            std::printf("  -m, --media DIR    Media directory (default: media)\n");
            std::printf("  -h, --help         Show this help\n");
            return 0;
        }
        if ((std::strcmp(argv[i], "-p") == 0 || std::strcmp(argv[i], "--player") == 0) && i + 1 < argc) {
            int id = std::atoi(argv[++i]);
            if (id >= 1 && id <= 15) player_id = static_cast<uint8_t>(id);
        }
        if ((std::strcmp(argv[i], "-c") == 0 || std::strcmp(argv[i], "--can") == 0) && i + 1 < argc) {
            can_iface = argv[++i];
        }
        if ((std::strcmp(argv[i], "-f") == 0 || std::strcmp(argv[i], "--fb") == 0) && i + 1 < argc) {
            fb_path = argv[++i];
        }
        if ((std::strcmp(argv[i], "-m") == 0 || std::strcmp(argv[i], "--media") == 0) && i + 1 < argc) {
            media_dir = argv[++i];
        }
    }

    std::fprintf(stderr, "╔══════════════════════════════════════╗\n");
    std::fprintf(stderr, "║     LASERTAG MP135 GAME DISPLAY     ║\n");
    std::fprintf(stderr, "╠══════════════════════════════════════╣\n");
    std::fprintf(stderr, "║  Player:     %s                      ║\n",
                 player_id > 0 ? std::to_string(player_id).c_str() : "AUTO");
    std::fprintf(stderr, "║  CAN:        %-8s                ║\n", can_iface);
    std::fprintf(stderr, "║  Framebuffer: %-10s              ║\n", fb_path);
    std::fprintf(stderr, "║  Media:      %-10s              ║\n", media_dir);
    std::fprintf(stderr, "╚══════════════════════════════════════╝\n");

    // ── Signal handling ──
    struct sigaction sa{};
    sa.sa_handler = signal_handler;
    sigemptyset(&sa.sa_mask);
    sigaction(SIGINT, &sa, nullptr);
    sigaction(SIGTERM, &sa, nullptr);

    // Ignore SIGCHLD to auto-reap audio child processes
    struct sigaction sa_chld{};
    sa_chld.sa_handler = SIG_IGN;
    sa_chld.sa_flags = SA_NOCLDWAIT;
    sigaction(SIGCHLD, &sa_chld, nullptr);

    // ── Initialize subsystems ──
    lt::EventBus bus;

    // Framebuffer
    lt::Framebuffer fb;
    if (!fb.open(fb_path)) {
        std::fprintf(stderr, "[FATAL] Cannot open framebuffer %s\n", fb_path);
        return 1;
    }

    // Animator
    lt::Animator animator(fb);

    // Audio
    lt::Audio audio(media_dir);

    // Game state
    lt::GameState game(bus, animator, audio);
    if (player_id > 0) {
        game.set_local_player(player_id);
    }
    // If player_id == 0, GameState starts in WAITING_FOR_DEVICE phase
    // and auto-assigns local_id from first DEVICE_REGISTER CAN frame
    game.init();

    // CAN receiver
    lt::CanReceiver can(bus, can_iface);
    if (!can.open()) {
        std::fprintf(stderr, "[WARN] CAN bus not available — running in display-only mode\n");
    } else {
        can.start();
    }

    std::fprintf(stderr, "[MAIN] All systems initialized. Running at ~30 FPS.\n");

    // ── Main loop (30 FPS target) ──
    constexpr uint32_t FRAME_MS = 33; // ~30 FPS

    while (g_running) {
        uint32_t frame_start = now_ms();

        // 1. Process all queued events
        bus.dispatch();

        // 2. Tick animations
        animator.tick(frame_start);

        // 3. Render HUD
        game.render_hud(fb, frame_start);

        // 4. Render animation overlays
        animator.render(frame_start);

        // 5. Present to screen
        fb.present();

        // 6. Frame timing
        uint32_t elapsed = now_ms() - frame_start;
        if (elapsed < FRAME_MS) {
            std::this_thread::sleep_for(std::chrono::milliseconds(FRAME_MS - elapsed));
        }
    }

    // ── Cleanup ──
    std::fprintf(stderr, "\n[MAIN] Shutting down...\n");
    can.stop();
    audio.stop_all();

    // Clear screen on exit
    fb.clear();
    fb.present();
    fb.close();

    std::fprintf(stderr, "[MAIN] Done.\n");
    return 0;
}
