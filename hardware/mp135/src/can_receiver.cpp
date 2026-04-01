#include "can_receiver.h"
#include "can_protocol.h"

#include <cerrno>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <unistd.h>

#include <linux/can.h>
#include <linux/can/raw.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <sys/socket.h>

namespace lt {

CanReceiver::CanReceiver(EventBus& bus, const std::string& iface)
    : bus_(bus), iface_(iface) {}

CanReceiver::~CanReceiver() {
    stop();
    if (sock_ >= 0) {
        ::close(sock_);
    }
}

bool CanReceiver::open(int /*bitrate*/) {
    // CAN interface should already be configured externally
    // (by systemd ExecStartPre or simulate.sh)
    // We only open the socket here.

    sock_ = socket(PF_CAN, SOCK_RAW, CAN_RAW);
    if (sock_ < 0) {
        std::fprintf(stderr, "[CAN] socket(): %s\n", std::strerror(errno));
        return false;
    }

    struct ifreq ifr{};
    std::strncpy(ifr.ifr_name, iface_.c_str(), IFNAMSIZ - 1);
    if (ioctl(sock_, SIOCGIFINDEX, &ifr) < 0) {
        std::fprintf(stderr, "[CAN] ioctl: %s\n", std::strerror(errno));
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

    struct timeval tv{1, 0};
    setsockopt(sock_, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

    std::fprintf(stderr, "[CAN] %s open\n", iface_.c_str());
    return true;
}

void CanReceiver::start() {
    if (running_.load() || sock_ < 0) return;
    running_.store(true);
    thread_ = std::thread(&CanReceiver::reader_loop, this);
}

void CanReceiver::stop() {
    running_.store(false);
    if (thread_.joinable()) {
        thread_.join();
    }
}

bool CanReceiver::send(uint32_t id, const uint8_t* data, uint8_t len) {
    if (sock_ < 0) return false;

    struct can_frame frame{};
    frame.can_id  = id;
    frame.can_dlc = (len > 8) ? 8 : len;
    if (data && frame.can_dlc > 0) {
        std::memcpy(frame.data, data, frame.can_dlc);
    }

    return write(sock_, &frame, sizeof(frame)) > 0;
}

void CanReceiver::reader_loop() {
    struct can_frame frame{};

    while (running_.load()) {
        ssize_t nbytes = read(sock_, &frame, sizeof(frame));
        if (nbytes < 0) {
            if (errno == EAGAIN || errno == EWOULDBLOCK) continue;
            std::fprintf(stderr, "[CAN] read: %s\n", std::strerror(errno));
            break;
        }
        if (nbytes != sizeof(struct can_frame)) continue;

        uint32_t id = frame.can_id & CAN_EFF_MASK;
        decode_frame(id, frame.data, frame.can_dlc);
    }
}

void CanReceiver::decode_frame(uint32_t id, const uint8_t* data, uint8_t len) {
    const uint8_t info  = can_info(id);
    const uint8_t func  = can_func(id);
    const uint8_t serie = can_serie(id);

    switch (info) {
    case INFO_COMBAT:
        switch (func) {
        case FN_CMB_HIT:
            // Player was hit — data[0] = attacker_id
            if (len >= 1) {
                bus_.emit(EventType::PLAYER_HIT, serie, data[0]);
            }
            break;
        case FN_CMB_KILL:
            // Player got a kill — data[0] = victim_id
            if (len >= 1) {
                bus_.emit(EventType::PLAYER_KILL, serie, data[0]);
            }
            break;
        case FN_CMB_KILLCOUNT:
            // Kill count update — data[0..1] = kills (big-endian)
            if (len >= 2) {
                uint16_t kills = (static_cast<uint16_t>(data[0]) << 8) | data[1];
                bus_.emit(EventType::KILLCOUNT_UPDATE, serie, 0, kills);
            }
            break;
        case FN_CMB_DEATH:
            bus_.emit(EventType::PLAYER_DEATH, serie);
            break;
        }
        break;

    case INFO_GAME:
        switch (func) {
        case FN_GAME_START:
            bus_.emit(EventType::GAME_START);
            break;
        case FN_GAME_END:
            bus_.emit(EventType::GAME_END);
            break;
        case FN_GAME_RESPAWN:
            bus_.emit(EventType::PLAYER_RESPAWN, serie);
            break;
        case FN_GAME_TEAM:
            if (len >= 1) {
                bus_.emit(EventType::TEAM_ASSIGN, serie, data[0]);
            }
            break;
        }
        break;

    case INFO_STATUS:
        switch (func) {
        case FN_STS_SCORE:
            if (len >= 2) {
                uint16_t score = (static_cast<uint16_t>(data[0]) << 8) | data[1];
                bus_.emit(EventType::SCORE_UPDATE, serie, 0, score);
            }
            break;
        case FN_STS_ALIVE:
            if (len >= 1) {
                bus_.emit(EventType::ALIVE_UPDATE, serie, 0, data[0]);
            }
            break;
        case FN_STS_AMMO:
            if (len >= 2) {
                uint16_t ammo = (static_cast<uint16_t>(data[0]) << 8) | data[1];
                bus_.emit(EventType::AMMO_UPDATE, serie, 0, ammo);
            }
            break;
        }
        break;

    case INFO_SYSTEM:
        switch (func) {
        case FN_SYS_HEARTBEAT:
            bus_.emit(EventType::HEARTBEAT, serie);
            break;
        case FN_SYS_REGISTER:
            bus_.emit(EventType::DEVICE_REGISTER, serie);
            break;
        }
        break;
    }
}

} // namespace lt
