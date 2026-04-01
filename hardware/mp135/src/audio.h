/*
 * Audio Engine — ALSA PCM WAV player for LaserTag MP135
 *
 * Plays WAV files asynchronously via ALSA for game sound effects.
 * Supports concurrent playback on separate threads.
 * Uses aplay as backend (reliable, no library dependency).
 */

#pragma once

#include <cstdint>
#include <string>
#include <unordered_map>

namespace lt {

/// Sound effect identifiers.
enum class SoundId : uint8_t {
    GAME_START,
    GAME_END,
    COUNTDOWN_TICK,
    COUNTDOWN_GO,
    HIT_RECEIVED,
    HIT_DEALT,
    KILL,
    DEATH,
    RESPAWN,
    AMMO_LOW,
    AMMO_EMPTY,
};

class Audio {
public:
    /// Initialize audio engine.
    /// @param media_dir  path to directory with WAV files
    explicit Audio(const std::string& media_dir = "media");

    /// Register a WAV file for a sound ID.
    void register_sound(SoundId id, const std::string& filename);

    /// Play a sound effect (non-blocking, fire & forget via fork+aplay).
    void play(SoundId id);

    /// Stop all currently playing sounds.
    void stop_all();

    /// Set master volume (0-100).
    void set_volume(int percent);

private:
    std::string media_dir_;
    std::unordered_map<SoundId, std::string> sounds_;
    int volume_ = 80;
};

} // namespace lt
