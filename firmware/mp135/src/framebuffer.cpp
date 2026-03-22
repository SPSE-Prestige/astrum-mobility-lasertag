#include "framebuffer.h"
#include "font.h"

#include <algorithm>
#include <cstdio>
#include <cstring>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <unistd.h>

#include <linux/fb.h>

namespace lt {

Framebuffer::~Framebuffer() {
    close();
}

bool Framebuffer::open(const char* path) {
    fd_ = ::open(path, O_RDWR);
    if (fd_ < 0) {
        std::fprintf(stderr, "[FB] cannot open %s\n", path);
        return false;
    }

    fb_var_screeninfo vinfo{};
    if (ioctl(fd_, FBIOGET_VSCREENINFO, &vinfo) < 0) {
        std::fprintf(stderr, "[FB] FBIOGET_VSCREENINFO failed\n");
        ::close(fd_);
        fd_ = -1;
        return false;
    }

    if (vinfo.bits_per_pixel != 16) {
        std::fprintf(stderr, "[FB] expected 16bpp, got %u\n", vinfo.bits_per_pixel);
        ::close(fd_);
        fd_ = -1;
        return false;
    }

    width_   = static_cast<int>(vinfo.xres);
    height_  = static_cast<int>(vinfo.yres);
    fb_size_ = static_cast<size_t>(width_) * height_ * 2;

    fb_ = static_cast<uint16_t*>(
        mmap(nullptr, fb_size_, PROT_READ | PROT_WRITE, MAP_SHARED, fd_, 0));
    if (fb_ == MAP_FAILED) {
        std::fprintf(stderr, "[FB] mmap failed\n");
        fb_ = nullptr;
        ::close(fd_);
        fd_ = -1;
        return false;
    }

    back_.resize(static_cast<size_t>(width_) * height_, Color::BLACK);

    std::fprintf(stderr, "[FB] %s open: %dx%d 16bpp\n", path, width_, height_);
    return true;
}

void Framebuffer::close() {
    if (fb_) {
        munmap(fb_, fb_size_);
        fb_ = nullptr;
    }
    if (fd_ >= 0) {
        ::close(fd_);
        fd_ = -1;
    }
}

// ── Drawing primitives ──

void Framebuffer::clear(uint16_t color) {
    std::fill(back_.begin(), back_.end(), color);
}

void Framebuffer::pixel(int x, int y, uint16_t color) {
    if (x >= 0 && x < width_ && y >= 0 && y < height_) {
        back_[y * width_ + x] = color;
    }
}

void Framebuffer::rect(int x, int y, int w, int h, uint16_t color) {
    int x0 = std::max(0, x);
    int y0 = std::max(0, y);
    int x1 = std::min(width_,  x + w);
    int y1 = std::min(height_, y + h);

    for (int py = y0; py < y1; py++) {
        for (int px = x0; px < x1; px++) {
            back_[py * width_ + px] = color;
        }
    }
}

void Framebuffer::rect_outline(int x, int y, int w, int h, uint16_t color, int thickness) {
    rect(x, y, w, thickness, color);                   // top
    rect(x, y + h - thickness, w, thickness, color);   // bottom
    rect(x, y, thickness, h, color);                   // left
    rect(x + w - thickness, y, thickness, h, color);   // right
}

void Framebuffer::gradient_h(int x, int y, int w, int h, uint16_t cl, uint16_t cr) {
    if (w <= 0) return;
    // Decompose to R5 G6 B5
    auto r0 = (cl >> 11) & 0x1F, g0 = (cl >> 5) & 0x3F, b0 = cl & 0x1F;
    auto r1 = (cr >> 11) & 0x1F, g1 = (cr >> 5) & 0x3F, b1 = cr & 0x1F;

    for (int px = 0; px < w; px++) {
        float t = static_cast<float>(px) / (w - 1);
        auto r = static_cast<uint16_t>(r0 + t * (r1 - r0));
        auto g = static_cast<uint16_t>(g0 + t * (g1 - g0));
        auto b = static_cast<uint16_t>(b0 + t * (b1 - b0));
        uint16_t c = (r << 11) | (g << 5) | b;
        for (int py = 0; py < h; py++) {
            pixel(x + px, y + py, c);
        }
    }
}

void Framebuffer::draw_char(int x, int y, char ch, int scale, uint16_t color) {
    uint8_t idx = static_cast<uint8_t>(ch);
    if (idx >= 128) idx = '?';

    const uint8_t* glyph = FONT_5x7[idx];

    for (int row = 0; row < GLYPH_H; row++) {
        for (int col = 0; col < GLYPH_W; col++) {
            if (!((glyph[row] >> (4 - col)) & 1)) continue;

            for (int sy = 0; sy < scale; sy++) {
                for (int sx = 0; sx < scale; sx++) {
                    pixel(x + col * scale + sx, y + row * scale + sy, color);
                }
            }
        }
    }
}

int Framebuffer::draw_text(int x, int y, const char* text, int scale, uint16_t color) {
    int cx = x;
    const int char_w = (GLYPH_W + 1) * scale; // 1px spacing

    while (*text) {
        draw_char(cx, y, *text, scale, color);
        cx += char_w;
        text++;
    }
    return cx - x;
}

void Framebuffer::draw_text_centered(int y, const char* text, int scale, uint16_t color) {
    int len = static_cast<int>(std::strlen(text));
    int text_w = len * (GLYPH_W + 1) * scale - scale; // subtract trailing space
    int x = (width_ - text_w) / 2;
    draw_text(x, y, text, scale, color);
}

void Framebuffer::draw_bar(int x, int y, int w, int h, float fraction,
                           uint16_t fg, uint16_t bg) {
    fraction = std::clamp(fraction, 0.0f, 1.0f);
    int filled = static_cast<int>(w * fraction);

    rect(x, y, w, h, bg);
    if (filled > 0) {
        rect(x, y, filled, h, fg);
    }
    rect_outline(x, y, w, h, Color::WHITE, 1);
}

void Framebuffer::tint(uint16_t color, uint8_t intensity) {
    // Fast RGB565 alpha-blend approx using bit shifting
    auto cr = (color >> 11) & 0x1F;
    auto cg = (color >> 5)  & 0x3F;
    auto cb = color & 0x1F;

    // intensity 0=no tint, 255=full tint
    int alpha = intensity;
    int inv   = 256 - alpha;

    int total = width_ * height_;
    for (int i = 0; i < total; i++) {
        uint16_t px = back_[i];
        auto pr = (px >> 11) & 0x1F;
        auto pg = (px >> 5)  & 0x3F;
        auto pb = px & 0x1F;

        auto nr = static_cast<uint16_t>((pr * inv + cr * alpha) >> 8);
        auto ng = static_cast<uint16_t>((pg * inv + cg * alpha) >> 8);
        auto nb = static_cast<uint16_t>((pb * inv + cb * alpha) >> 8);

        back_[i] = (nr << 11) | (ng << 5) | nb;
    }
}

void Framebuffer::present() {
    if (fb_) {
        std::memcpy(fb_, back_.data(), fb_size_);
    }
}

} // namespace lt
