/*
 * ═══════════════════════════════════════════════════
 *  LaserTag CAN Protocol — TEST
 * ═══════════════════════════════════════════════════
 *
 * Testuje adresovací schéma CAN ID: 0xABC
 *   A = informace (kategorie)
 *   B = funkce
 *   C = série (hráč/zařízení)
 *
 * Build:   make example
 * Run:     sudo ./main
 *
 * Bez CAN hardware simuluje zprávy v paměti.
 * S CAN hardware posílá a přijímá reálné framy.
 */

#include "can_bus.h"
#include "can_protocol.h"
#include "mqtt_client.h"

#include <csignal>
#include <cstdio>
#include <cstring>
#include <atomic>
#include <thread>
#include <chrono>

static std::atomic<bool> g_running{true};
static void on_signal(int) { g_running = false; }

// ══════════════════════════════════════════════════
//  TEST 1: Protocol ID scheme verification
// ══════════════════════════════════════════════════

static void test_protocol_ids() {
    std::printf("\n════ TEST 1: CAN ID adresování ════\n\n");

    // Ukázka: can1  0x103  [1]  05
    //         info=1(combat), func=0(hit), serie=3(hráč 3)
    std::printf("Formát: can1  ID    [len]  data\n");
    std::printf("        can1  0xABC [n]    ...\n");
    std::printf("        A=informace, B=funkce, C=série\n\n");

    // Hit: hráč 3 byl trefený útočníkem 7
    uint32_t hit_id = lt::id_hit(3);
    std::printf("Hráč 3 byl trefený:     can1  %03X   [1]  %02X\n",
                hit_id, 0x07);
    lt::print_can_id(hit_id);

    // Kill: hráč 5 dal kill hráči 2
    uint32_t kill_id = lt::id_kill(5);
    std::printf("\nHráč 5 dal kill:        can1  %03X   [1]  %02X\n",
                kill_id, 0x02);
    lt::print_can_id(kill_id);

    // Kill count: hráč 5 má 12 killů
    uint32_t kc_id = lt::id_killcount(5);
    std::printf("\nKill count hráče 5:     can1  %03X   [2]  %02X %02X\n",
                kc_id, 0x00, 0x0C);
    lt::print_can_id(kc_id);

    // Smrt: hráč 2 umřel
    uint32_t death_id = lt::id_death(2);
    std::printf("\nHráč 2 umřel:           can1  %03X   [0]\n",
                death_id);
    lt::print_can_id(death_id);

    // Game start broadcast
    uint32_t start_id = lt::id_game_start();
    std::printf("\nGame START (broadcast):  can1  %03X   [0]\n",
                start_id);
    lt::print_can_id(start_id);

    // Game end broadcast
    uint32_t end_id = lt::id_game_end();
    std::printf("\nGame END (broadcast):    can1  %03X   [0]\n",
                end_id);
    lt::print_can_id(end_id);

    // Custom: make_can_id test
    uint32_t custom = lt::make_can_id(0x3, 0x0, 0xA);
    std::printf("\nCustom (skóre hráč 10): can1  %03X   [4]  00 00 00 2A\n",
                custom);
    lt::print_can_id(custom);

    std::printf("\n════ Všechny ID testy OK ════\n");
}

// ══════════════════════════════════════════════════
//  TEST 2: Decode incoming CAN frame
// ══════════════════════════════════════════════════

static void test_decode_frame() {
    std::printf("\n════ TEST 2: Dekódování CAN framů ════\n\n");

    // Simulace příchozích framů
    struct TestFrame {
        uint32_t id;
        uint8_t  data[8];
        uint8_t  len;
        const char* desc;
    };

    TestFrame frames[] = {
        // Hráč 3 trefený útočníkem 7
        { lt::id_hit(3),        {0x07},          1, "Hráč 3 trefený" },
        // Hráč 5 dal kill hráči 2
        { lt::id_kill(5),       {0x02},          1, "Hráč 5 dal kill" },
        // Kill count hráče 5 = 12
        { lt::id_killcount(5),  {0x00, 0x0C},    2, "Kill count hráče 5" },
        // Hráč 2 umřel
        { lt::id_death(2),      {},              0, "Hráč 2 umřel" },
        // Game start
        { lt::id_game_start(),  {},              0, "Game START" },
        // Heartbeat hráče 1
        { lt::make_can_id(lt::INFO_SYSTEM, lt::FN_SYS_HEARTBEAT, 1),
                                {0xFF},          1, "Heartbeat hráč 1" },
    };

    for (auto& f : frames) {
        uint8_t info  = lt::can_info(f.id);
        uint8_t func  = lt::can_func(f.id);
        uint8_t serie = lt::can_serie(f.id);

        std::printf("  can1  %03X   [%d]  ", f.id, f.len);
        for (int i = 0; i < f.len; i++)
            std::printf("%02X ", f.data[i]);

        std::printf("\n");
        std::printf("    → info=%X func=%X serie=%X  (%s)\n",
                    info, func, serie, f.desc);

        // Decode by category
        switch (info) {
        case lt::INFO_COMBAT:
            switch (func) {
            case lt::FN_CMB_HIT:
                std::printf("    → COMBAT: Hráč %d byl trefený", serie);
                if (f.len > 0) std::printf(" útočníkem %d", f.data[0]);
                std::printf("\n");
                break;
            case lt::FN_CMB_KILL:
                std::printf("    → COMBAT: Hráč %d dal kill", serie);
                if (f.len > 0) std::printf(" hráči %d", f.data[0]);
                std::printf("\n");
                break;
            case lt::FN_CMB_KILLCOUNT: {
                uint16_t kills = (f.len >= 2) ? (f.data[0] << 8 | f.data[1]) : 0;
                std::printf("    → COMBAT: Hráč %d má %d killů\n", serie, kills);
                break;
            }
            case lt::FN_CMB_DEATH:
                std::printf("    → COMBAT: Hráč %d umřel\n", serie);
                break;
            }
            break;
        case lt::INFO_GAME:
            if (func == lt::FN_GAME_START)
                std::printf("    → GAME: Start hry!\n");
            else if (func == lt::FN_GAME_END)
                std::printf("    → GAME: Konec hry!\n");
            break;
        case lt::INFO_SYSTEM:
            if (func == lt::FN_SYS_HEARTBEAT)
                std::printf("    → SYSTEM: Heartbeat od zařízení %d\n", serie);
            break;
        default:
            std::printf("    → UNKNOWN\n");
            break;
        }
        std::printf("\n");
    }

    std::printf("════ Dekódování OK ════\n");
}

// ══════════════════════════════════════════════════
//  TEST 3: Real CAN + MQTT (only with hardware)
// ══════════════════════════════════════════════════

static void test_real_can_mqtt() {
    std::printf("\n════ TEST 3: Reálný CAN + MQTT test ════\n\n");

    // ── MQTT ──
    lt::MqttClient mqtt;
    bool mqtt_ok = false;

    mqtt.on_connect([&](bool ok) {
        mqtt_ok = ok;
        if (ok) {
            std::printf("[MQTT] Připojeno! Subscribuju device/+/command\n");
            mqtt.subscribe("device/+/command");
        }
    });

    mqtt.on_message([&](const std::string& topic, const std::string& payload) {
        std::printf("[MQTT←] %s: %s\n", topic.c_str(), payload.c_str());
    });

    if (!mqtt.connect("127.0.0.1", 1883, "lt-test")) {
        std::printf("[MQTT] Nepodařilo se připojit (pokračuju bez MQTT)\n");
    }

    // ── CAN ──
    lt::CanBus can;

    if (!can.open("can1")) {
        std::printf("[CAN] can1 se nepodařilo otevřít (pokračuju bez CAN)\n");
        std::printf("[CAN] Pro test bez HW použij: sudo ip link add dev vcan0 type vcan && sudo ip link set up vcan0\n");

        // zkus vcan0
        if (!can.open("vcan0", 1000000)) {
            std::printf("[CAN] Ani vcan0 nefunguje. Pouze protocol testy.\n");
            mqtt.disconnect();
            return;
        }
    }

    can.set_handler([&](const lt::CanFrame& f) {
        uint8_t info  = lt::can_info(f.id);
        uint8_t func  = lt::can_func(f.id);
        uint8_t serie = lt::can_serie(f.id);

        std::printf("[CAN←]  can1  %03X   [%d]  ", f.id, f.len);
        for (int i = 0; i < f.len; i++)
            std::printf("%02X ", f.data[i]);
        std::printf("  (info=%X func=%X serie=%X)\n", info, func, serie);

        // Forward combat events to MQTT
        if (info == lt::INFO_COMBAT && mqtt_ok) {
            std::string device = "gun-" + std::to_string(serie);
            char json[256];

            if (func == lt::FN_CMB_HIT && f.len >= 1) {
                std::snprintf(json, sizeof(json),
                    R"({"event":"hit","player":%d,"attacker":%d})",
                    serie, f.data[0]);
                mqtt.publish("device/" + device + "/event", json);
                std::printf("[MQTT→] %s\n", json);
            }
            else if (func == lt::FN_CMB_KILL && f.len >= 1) {
                std::snprintf(json, sizeof(json),
                    R"({"event":"kill","player":%d,"victim":%d})",
                    serie, f.data[0]);
                mqtt.publish("device/" + device + "/event", json);
                std::printf("[MQTT→] %s\n", json);
            }
        }
    });

    can.start();

    // Pošli testovací framy
    std::printf("\n[TEST] Posílám testovací CAN framy...\n\n");

    // Hráč 3 byl trefený hráčem 7
    uint8_t hit_data[] = {0x07};
    can.send(lt::id_hit(3), hit_data, 1);
    std::printf("[CAN→]  can1  %03X   [1]  07    (hráč 3 trefený hráčem 7)\n",
                lt::id_hit(3));

    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    // Hráč 7 dal kill hráči 3
    uint8_t kill_data[] = {0x03};
    can.send(lt::id_kill(7), kill_data, 1);
    std::printf("[CAN→]  can1  %03X   [1]  03    (hráč 7 dal kill hráči 3)\n",
                lt::id_kill(7));

    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    // Kill count hráče 7 = 5 killů
    uint8_t kc_data[] = {0x00, 0x05};
    can.send(lt::id_killcount(7), kc_data, 2);
    std::printf("[CAN→]  can1  %03X   [2]  00 05 (hráč 7 má 5 killů)\n",
                lt::id_killcount(7));

    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    // Game start broadcast
    can.send(lt::id_game_start(), nullptr, 0);
    std::printf("[CAN→]  can1  %03X   [0]        (game start broadcast)\n",
                lt::id_game_start());

    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    // Game end broadcast
    can.send(lt::id_game_end(), nullptr, 0);
    std::printf("[CAN→]  can1  %03X   [0]        (game end broadcast)\n",
                lt::id_game_end());

    // Čekej na odpovědi
    std::printf("\n[TEST] Čekám na příchozí framy (Ctrl+C pro konec)...\n");
    while (g_running) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    can.stop();
    can.close();
    mqtt.disconnect();
    std::printf("\n════ CAN+MQTT test ukončen ════\n");
}

// ══════════════════════════════════════════════════
//  Main
// ══════════════════════════════════════════════════

int main(int argc, char** argv) {
    std::signal(SIGINT, on_signal);
    std::signal(SIGTERM, on_signal);

    std::printf("╔═══════════════════════════════════════╗\n");
    std::printf("║   LaserTag CAN Protocol Test Suite    ║\n");
    std::printf("╠═══════════════════════════════════════╣\n");
    std::printf("║  Adresování: 0xABC                   ║\n");
    std::printf("║    A = informace (kategorie)          ║\n");
    std::printf("║    B = funkce                         ║\n");
    std::printf("║    C = série (hráč/zařízení 0-F)     ║\n");
    std::printf("╚═══════════════════════════════════════╝\n");

    bool real_test = false;
    if (argc > 1 && std::strcmp(argv[1], "--real") == 0) {
        real_test = true;
    }

    // Vždy spusť testy protokolu
    test_protocol_ids();
    test_decode_frame();

    // Reálný HW test jen s --real
    if (real_test) {
        test_real_can_mqtt();
    } else {
        std::printf("\n──────────────────────────────────\n");
        std::printf("Pro reálný CAN+MQTT test spusť:  sudo ./main --real\n");
        std::printf("──────────────────────────────────\n");
    }

    return 0;
}
