#include <cstdint>
#include <cstring>
#include <fcntl.h>
#include <iostream>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <unistd.h>

#ifdef __linux__
#include <linux/fb.h>
#endif

// RGB565 colors
static constexpr uint16_t COLOR_BLACK = 0x0000;
static constexpr uint16_t COLOR_YELLOW = 0xFFE0;

// 5x7 bitmap digits 0-9
static const uint8_t DIGITS[10][7] = {
    {0x1F, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1F}, // 0
    {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}, // 1
    {0x1F, 0x01, 0x01, 0x1F, 0x10, 0x10, 0x1F}, // 2
    {0x1F, 0x01, 0x01, 0x1F, 0x01, 0x01, 0x1F}, // 3
    {0x11, 0x11, 0x11, 0x1F, 0x01, 0x01, 0x01}, // 4
    {0x1F, 0x10, 0x10, 0x1F, 0x01, 0x01, 0x1F}, // 5
    {0x1F, 0x10, 0x10, 0x1F, 0x11, 0x11, 0x1F}, // 6
    {0x1F, 0x01, 0x01, 0x02, 0x04, 0x08, 0x10}, // 7
    {0x1F, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x1F}, // 8
    {0x1F, 0x11, 0x11, 0x1F, 0x01, 0x01, 0x1F}  // 9
};

#ifdef __linux__
static void draw_big_digit(
    uint16_t* fb,
    int width,
    int height,
    int start_x,
    int start_y,
    int digit,
    int scale,
    uint16_t color
) {
    for (int row = 0; row < 7; row++) {
        for (int col = 0; col < 5; col++) {
            if (((DIGITS[digit][row] >> (4 - col)) & 1U) == 0U) {
                continue;
            }

            for (int sy = 0; sy < scale; sy++) {
                for (int sx = 0; sx < scale; sx++) {
                    const int x = start_x + col * scale + sx;
                    const int y = start_y + row * scale + sy;
                    if (x >= 0 && x < width && y >= 0 && y < height) {
                        fb[y * width + x] = color;
                    }
                }
            }
        }
    }
}

int main() {
    const char* fb_path = "/dev/fb1";
    const int fd = open(fb_path, O_RDWR);
    if (fd < 0) {
        std::cerr << "Error: cannot open " << fb_path << " (try sudo).\n";
        return 1;
    }

    fb_var_screeninfo vinfo{};
    if (ioctl(fd, FBIOGET_VSCREENINFO, &vinfo) < 0) {
        std::cerr << "Error: FBIOGET_VSCREENINFO failed.\n";
        close(fd);
        return 1;
    }

    if (vinfo.bits_per_pixel != 16) {
        std::cerr << "Error: expected 16bpp framebuffer, got " << vinfo.bits_per_pixel << "bpp.\n";
        close(fd);
        return 1;
    }

    const size_t fb_size = static_cast<size_t>(vinfo.xres) * static_cast<size_t>(vinfo.yres) * 2U;
    auto* fb = static_cast<uint16_t*>(mmap(nullptr, fb_size, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0));
    if (fb == MAP_FAILED) {
        std::cerr << "Error: mmap failed.\n";
        close(fd);
        return 1;
    }

    // Clear screen
    std::memset(fb, 0, fb_size);

    const int width = static_cast<int>(vinfo.xres);
    const int height = static_cast<int>(vinfo.yres);
    const int scale = 35;
    const int digit_width = 5 * scale;
    const int digit_height = 7 * scale;
    const int gap = scale;
    const int text_width = digit_width * 2 + gap;

    const int start_x = (width - text_width) / 2;
    const int start_y = (height - digit_height) / 2;

    draw_big_digit(fb, width, height, start_x, start_y, 6, scale, COLOR_YELLOW);
    draw_big_digit(fb, width, height, start_x + digit_width + gap, start_y, 9, scale, COLOR_YELLOW);

    std::cout << "Done: rendered 69 on " << fb_path << "\n";

    munmap(fb, fb_size);
    close(fd);
    return 0;
}
#else
int main() {
    std::cout << "This file is for Linux framebuffer on MP135 (/dev/fb1).\n";
    std::cout << "Build and run it directly on MP135 Linux.\n";
    return 0;
}
#endif
