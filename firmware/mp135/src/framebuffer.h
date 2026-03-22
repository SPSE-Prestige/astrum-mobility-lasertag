/*
 * Framebuffer — RGB565 Linux framebuffer abstraction for MP135
 *
 * Provides double-buffered rendering via /dev/fb1 with
 * primitive drawing operations: pixel, rect, text, fill.
 */

#pragma once

#include <cstdint>
#include <string>
#include <vector>

namespace lt {

// ── RGB565 color helpers ──
constexpr uint16_t rgb565(uint8_t r, uint8_t g, uint8_t b) {
    return static_cast<uint16_t>(
        ((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3));
}

// ── Predefined colors ──
namespace Color {
    constexpr uint16_t BLACK      = 0x0000;
    constexpr uint16_t WHITE      = 0xFFFF;
    constexpr uint16_t RED        = rgb565(255, 0, 0);
    constexpr uint16_t GREEN      = rgb565(0, 255, 0);
    constexpr uint16_t BLUE       = rgb565(0, 0, 255);
    constexpr uint16_t YELLOW     = rgb565(255, 255, 0);
    constexpr uint16_t CYAN       = rgb565(0, 255, 255);
    constexpr uint16_t MAGENTA    = rgb565(255, 0, 255);
    constexpr uint16_t ORANGE     = rgb565(255, 165, 0);
    constexpr uint16_t DARK_RED   = rgb565(139, 0, 0);
    constexpr uint16_t DARK_GREEN = rgb565(0, 100, 0);
    constexpr uint16_t GRAY       = rgb565(128, 128, 128);
    constexpr uint16_t DARK_GRAY  = rgb565(64, 64, 64);
} // namespace Color

class Framebuffer {
public:
    Framebuffer() = default;
    ~Framebuffer();

    Framebuffer(const Framebuffer&) = delete;
    Framebuffer& operator=(const Framebuffer&) = delete;

    /// Open /dev/fb1 and mmap the framebuffer.
    bool open(const char* path = "/dev/fb1");

    /// Close and unmap.
    void close();

    int width() const  { return width_; }
    int height() const { return height_; }

    // ── Drawing primitives (draw to back buffer) ──

    void clear(uint16_t color = Color::BLACK);

    void pixel(int x, int y, uint16_t color);

    void rect(int x, int y, int w, int h, uint16_t color);

    void rect_outline(int x, int y, int w, int h, uint16_t color, int thickness = 1);

    /// Fill a horizontal gradient from color_l to color_r.
    void gradient_h(int x, int y, int w, int h, uint16_t color_l, uint16_t color_r);

    /// Draw a single character at (x, y) with given scale.
    void draw_char(int x, int y, char ch, int scale, uint16_t color);

    /// Draw a null-terminated string. Returns the total pixel width drawn.
    int draw_text(int x, int y, const char* text, int scale, uint16_t color);

    /// Draw text centered horizontally.
    void draw_text_centered(int y, const char* text, int scale, uint16_t color);

    /// Draw a progress/health bar.
    void draw_bar(int x, int y, int w, int h, float fraction,
                  uint16_t fg, uint16_t bg = Color::DARK_GRAY);

    /// Tint the entire screen (alpha blend effect via bit shift).
    void tint(uint16_t color, uint8_t intensity);

    // ── Buffer management ──

    /// Swap back buffer → framebuffer (present).
    void present();

    /// Direct access to back buffer for custom effects.
    uint16_t* buffer() { return back_.data(); }
    const uint16_t* buffer() const { return back_.data(); }

private:
    int          fd_     = -1;
    int          width_  = 0;
    int          height_ = 0;
    uint16_t*    fb_     = nullptr;   // mmap'd framebuffer
    size_t       fb_size_ = 0;
    std::vector<uint16_t> back_;      // back buffer
};

} // namespace lt
