#include <algorithm>
#include <cstdint>
#include <cstdio>
#include <iostream>
#include <string>
#include <vector>

#include "../esp/can_protocol.h"
#include <sqlite3.h>

namespace {

struct PlayerStat {
    uint8_t player_id = 0;
    uint16_t kills = 0;
    uint16_t deaths = 0;
    uint16_t score = 0;
    uint16_t ammo = 0;
    bool alive = true;
};

struct CanOutFrame {
    uint32_t id = 0;
    uint8_t len = 0;
    uint8_t data[8]{};
    std::string label;
};

bool in_player_range(int v) {
    return v >= 1 && v <= 15;
}

uint16_t clamp_u16(int v) {
    if (v < 0) {
        return 0;
    }
    if (v > 65535) {
        return 65535;
    }
    return static_cast<uint16_t>(v);
}

void print_usage(const char* exe) {
    std::cout << "Usage: " << exe << " <db_path>\n";
    std::cout << "\nExpected table: player_stats\n";
    std::cout << "Columns: player_id INTEGER, kills INTEGER, deaths INTEGER, score INTEGER, ammo INTEGER, alive INTEGER\n";
}

bool load_player_stats(sqlite3* db, std::vector<PlayerStat>& out) {
    constexpr const char* kSql =
        "SELECT player_id, kills, deaths, score, ammo, alive "
        "FROM player_stats ORDER BY player_id ASC;";

    sqlite3_stmt* stmt = nullptr;
    int rc = sqlite3_prepare_v2(db, kSql, -1, &stmt, nullptr);
    if (rc != SQLITE_OK) {
        std::cerr << "SQLite prepare failed: " << sqlite3_errmsg(db) << "\n";
        return false;
    }

    while ((rc = sqlite3_step(stmt)) == SQLITE_ROW) {
        const int player_id = sqlite3_column_int(stmt, 0);
        if (!in_player_range(player_id)) {
            continue;
        }

        PlayerStat s;
        s.player_id = static_cast<uint8_t>(player_id);
        s.kills = clamp_u16(sqlite3_column_int(stmt, 1));
        s.deaths = clamp_u16(sqlite3_column_int(stmt, 2));
        s.score = clamp_u16(sqlite3_column_int(stmt, 3));
        s.ammo = clamp_u16(sqlite3_column_int(stmt, 4));
        s.alive = sqlite3_column_int(stmt, 5) != 0;
        out.push_back(s);
    }

    if (rc != SQLITE_DONE) {
        std::cerr << "SQLite step failed: " << sqlite3_errmsg(db) << "\n";
        sqlite3_finalize(stmt);
        return false;
    }

    sqlite3_finalize(stmt);
    return true;
}

void push_u16_frame(std::vector<CanOutFrame>& out, uint32_t id, uint16_t value, const char* label) {
    CanOutFrame f;
    f.id = id;
    f.len = 2;
    f.data[0] = static_cast<uint8_t>((value >> 8) & 0xFF);
    f.data[1] = static_cast<uint8_t>(value & 0xFF);
    f.label = label;
    out.push_back(f);
}

void push_u8_frame(std::vector<CanOutFrame>& out, uint32_t id, uint8_t value, const char* label) {
    CanOutFrame f;
    f.id = id;
    f.len = 1;
    f.data[0] = value;
    f.label = label;
    out.push_back(f);
}

std::vector<CanOutFrame> build_frames(const std::vector<PlayerStat>& stats) {
    std::vector<CanOutFrame> out;
    out.reserve(stats.size() * 4);

    for (const auto& s : stats) {
        push_u16_frame(out, lt::id_killcount(s.player_id), s.kills, "killcount");
        push_u16_frame(out, lt::make_can_id(lt::INFO_STATUS, lt::FN_STS_SCORE, s.player_id), s.score, "score");
        push_u16_frame(out, lt::make_can_id(lt::INFO_STATUS, lt::FN_STS_AMMO, s.player_id), s.ammo, "ammo");
        push_u8_frame(out, lt::make_can_id(lt::INFO_STATUS, lt::FN_STS_ALIVE, s.player_id), s.alive ? 1 : 0, "alive");
    }

    return out;
}

void print_ascii_graph(const std::vector<PlayerStat>& stats) {
    std::cout << "\n=== SCORE GRAPH ===\n";
    uint16_t max_score = 1;
    for (const auto& s : stats) {
        max_score = std::max(max_score, s.score);
    }

    for (const auto& s : stats) {
        const int bar_len = static_cast<int>((40.0 * s.score) / max_score);
        std::printf("P%02u |", s.player_id);
        for (int i = 0; i < bar_len; i++) {
            std::printf("#");
        }
        for (int i = bar_len; i < 40; i++) {
            std::printf(" ");
        }
        std::printf("| score=%u kills=%u deaths=%u ammo=%u %s\n",
                    s.score, s.kills, s.deaths, s.ammo, s.alive ? "ALIVE" : "DEAD");
    }
}

void print_frames(const std::vector<CanOutFrame>& frames) {
    std::cout << "\n=== GENERATED CAN FRAMES ===\n";
    for (const auto& f : frames) {
        std::printf("[%s] ", f.label.c_str());
        lt::print_can_id(f.id);
        std::printf("  data:");
        for (uint8_t i = 0; i < f.len; i++) {
            std::printf(" %02X", f.data[i]);
        }
        std::printf("\n");
    }
}

} // namespace

int main(int argc, char** argv) {
    if (argc < 2) {
        print_usage(argv[0]);
        return 1;
    }

    const char* db_path = argv[1];
    sqlite3* db = nullptr;

    const int rc = sqlite3_open(db_path, &db);
    if (rc != SQLITE_OK) {
        std::cerr << "Cannot open DB: " << db_path << "\n";
        if (db) {
            std::cerr << "SQLite error: " << sqlite3_errmsg(db) << "\n";
            sqlite3_close(db);
        }
        return 1;
    }

    std::vector<PlayerStat> stats;
    if (!load_player_stats(db, stats)) {
        sqlite3_close(db);
        return 1;
    }

    sqlite3_close(db);

    if (stats.empty()) {
        std::cerr << "No rows loaded from player_stats.\n";
        return 1;
    }

    std::cout << "Loaded players: " << stats.size() << "\n";

    print_ascii_graph(stats);

    const auto frames = build_frames(stats);
    print_frames(frames);

    std::cout << "\nDone: DB data processed and CAN payloads generated.\n";
    return 0;
}
