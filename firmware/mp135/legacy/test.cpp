#include <iostream>
#include <cstdint>
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#ifdef __linux__
#include <linux/fb.h>
#endif

// Definice barev (RGB565)
#define COLOR_BLACK   0x0000
#define COLOR_YELLOW  0xFFE0
#define COLOR_WHITE   0xFFFF

// Jednoduchá matice pro číslice (5x7)
const unsigned char digits[10][7] = {
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

void drawBigDigit(uint16_t* fbp, int startX, int startY, int digit, int scale, int screenWidth, uint16_t color) {
    for (int i = 0; i < 7; i++) { // řádky
        for (int j = 0; j < 5; j++) { // sloupce
            if ((digits[digit][i] >> (4 - j)) & 1) {
                for (int py = 0; py < scale; py++) {
                    for (int px = 0; px < scale; px++) {
                        int x = startX + j * scale + px;
                        int y = startY + i * scale + py;
                        fbp[y * screenWidth + x] = color;
                    }
                }
            }
        }
    }
}

#ifdef __linux__
int main() {
    // Otevření zařízení fb1
    int fbfd = open("/dev/fb1", O_RDWR);
    if (fbfd == -1) {
        std::cerr << "Chyba: Nelze otevrit /dev/fb1. Zkus sudo." << std::endl;
        return 1;
    }

    struct fb_var_screeninfo vinfo;
    ioctl(fbfd, FBIOGET_VSCREENINFO, &vinfo);

    // Mapování paměti
    long screensize = vinfo.xres * vinfo.yres * 2; // 2 byty na pixel (16-bit)
    uint16_t* fbp = (uint16_t*)mmap(0, screensize, PROT_READ | PROT_WRITE, MAP_SHARED, fbfd, 0);

    // Vyčistit obrazovku (Černá)
    for (int i = 0; i < vinfo.xres * vinfo.yres; i++) fbp[i] = COLOR_BLACK;

    // Nastavení velikosti a pozice pro "69"
    int scale = 35; // Velikost číslic
    int spacing = 10; // Mezera mezi nimi
    int totalWidth = (5 * scale * 2) + (spacing * scale);
    int startX = (vinfo.xres - totalWidth) / 2;
    int startY = (vinfo.yres - (7 * scale)) / 2;

    // Vykreslení
    drawBigDigit(fbp, startX, startY, 6, scale, vinfo.xres, COLOR_YELLOW);
    drawBigDigit(fbp, startX + (6 * scale), startY, 9, scale, vinfo.xres, COLOR_YELLOW);

    std::cout << "Displej MP135: Vykresleno 69 na IP 192.168.0.213" << std::endl;

    munmap(fbp, screensize);
    close(fbfd);
    return 0;
}
#else
int main() {
    std::cout << "Tento program pouziva Linux framebuffer API (/dev/fb1, linux/fb.h)." << std::endl;
    std::cout << "Na macOS se prelozi, ale framebuffer cast nelze spustit." << std::endl;
    return 0;
}
#endif