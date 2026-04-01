#include "can_bus.h"

#include <cerrno>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <thread>
#include <unistd.h>

#include <linux/can.h>
#include <linux/can/raw.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <sys/socket.h>

namespace lt {

// ── Lifecycle ──

CanBus::~CanBus() {
    close();
}

bool CanBus::open(const std::string& iface, int bitrate) {
    iface_ = iface;

    // Configure interface via ip link
    char cmd[512];
    std::snprintf(cmd, sizeof(cmd),
        "ip link set %s down 2>/dev/null; "
        "ip link set %s up type can bitrate %d restart-ms 100 loopback off 2>&1",
        iface.c_str(), iface.c_str(), bitrate);

    if (std::system(cmd) != 0) {
        std::fprintf(stderr, "[CAN] failed to configure %s\n", iface.c_str());
        return false;
    }

    // Open SocketCAN
    sock_ = socket(PF_CAN, SOCK_RAW, CAN_RAW);
    if (sock_ < 0) {
        std::fprintf(stderr, "[CAN] socket(): %s\n", std::strerror(errno));
        return false;
    }

    struct ifreq ifr{};
    std::strncpy(ifr.ifr_name, iface.c_str(), IFNAMSIZ - 1);
    if (ioctl(sock_, SIOCGIFINDEX, &ifr) < 0) {
        std::fprintf(stderr, "[CAN] ioctl SIOCGIFINDEX: %s\n", std::strerror(errno));
        ::close(sock_);
        sock_ = -1;
        return false;
    }

    struct sockaddr_can addr{};
    addr.can_family  = AF_CAN;
    addr.can_ifindex = ifr.ifr_ifindex;

    if (bind(sock_, reinterpret_cast<struct sockaddr*>(&addr), sizeof(addr)) < 0) {
        std::fprintf(stderr, "[CAN] bind: %s\n", std::strerror(errno));
        ::close(sock_);
        sock_ = -1;
        return false;
    }

    // Read timeout 1s (so reader thread can check running_ flag)
    struct timeval tv{1, 0};
    setsockopt(sock_, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

    std::fprintf(stderr, "[CAN] %s open (bitrate=%d)\n", iface.c_str(), bitrate);
    return true;
}

void CanBus::close() {
    stop();
    if (sock_ >= 0) {
        ::close(sock_);
        sock_ = -1;
    }
}

bool CanBus::is_open() const {
    return sock_ >= 0;
}

// ── Send ──

bool CanBus::send(uint32_t id, const uint8_t* data, uint8_t len) {
    if (sock_ < 0) return false;

    struct can_frame frame{};
    frame.can_id  = id;
    frame.can_dlc = (len > 8) ? 8 : len;
    if (data && frame.can_dlc > 0)
        std::memcpy(frame.data, data, frame.can_dlc);

    ssize_t n = write(sock_, &frame, sizeof(frame));
    if (n < 0) {
        std::fprintf(stderr, "[CAN] write: %s\n", std::strerror(errno));
        return false;
    }

    return true;
}

// ── Receive ──

void CanBus::set_handler(CanHandler handler) {
    handler_ = std::move(handler);
}

void CanBus::start() {
    if (running_ || sock_ < 0) return;
    running_ = true;
    thread_  = new std::thread(&CanBus::reader_loop, this);
}

void CanBus::stop() {
    running_ = false;
    if (thread_) {
        auto* t = static_cast<std::thread*>(thread_);
        if (t->joinable()) t->join();
        delete t;
        thread_ = nullptr;
    }
}

void CanBus::reader_loop() {
    struct can_frame frame{};

    while (running_) {
        ssize_t nbytes = read(sock_, &frame, sizeof(frame));
        if (nbytes < 0) {
            if (errno == EAGAIN || errno == EWOULDBLOCK)
                continue;
            std::fprintf(stderr, "[CAN] read: %s\n", std::strerror(errno));
            break;
        }
        if (nbytes != sizeof(struct can_frame)) continue;
        if (!handler_) continue;

        CanFrame f{};
        f.id  = frame.can_id & CAN_EFF_MASK;
        f.len = frame.can_dlc;
        std::memcpy(f.data, frame.data, f.len);

        handler_(f);
    }
}

} // namespace lt
