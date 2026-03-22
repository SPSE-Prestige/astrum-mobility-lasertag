#include "audio.h"

#include <cstdio>
#include <cstdlib>
#include <fcntl.h>
#include <sys/wait.h>
#include <unistd.h>

namespace lt {

Audio::Audio(const std::string& media_dir)
    : media_dir_(media_dir) {
    // Register default sound mappings
    register_sound(SoundId::DEATH,      "gta-death.wav");
    register_sound(SoundId::KILL,       "mario-death.wav");
}

void Audio::register_sound(SoundId id, const std::string& filename) {
    sounds_[id] = media_dir_ + "/" + filename;
}

void Audio::play(SoundId id) {
    auto it = sounds_.find(id);
    if (it == sounds_.end()) return;

    const std::string& path = it->second;

    // Check file exists
    if (access(path.c_str(), R_OK) != 0) {
        std::fprintf(stderr, "[AUDIO] file not found: %s\n", path.c_str());
        return;
    }

    // Fork + exec aplay for async non-blocking playback
    pid_t pid = fork();
    if (pid < 0) {
        std::fprintf(stderr, "[AUDIO] fork failed\n");
        return;
    }

    if (pid == 0) {
        // Child process — redirect stdout/stderr to /dev/null
        int devnull = ::open("/dev/null", O_WRONLY);
        if (devnull >= 0) {
            dup2(devnull, STDOUT_FILENO);
            dup2(devnull, STDERR_FILENO);
            ::close(devnull);
        }

        // Execute aplay
        execlp("aplay", "aplay", "-q", path.c_str(), nullptr);
        // If execlp fails
        _exit(1);
    }

    // Parent — reap child asynchronously (no zombie)
    // Use SIGCHLD ignore or waitpid in background
    // We'll do a non-blocking waitpid sweep
    int status;
    while (waitpid(-1, &status, WNOHANG) > 0) {
        // Reap any finished children
    }
}

void Audio::stop_all() {
    // Kill all aplay processes owned by this user
    std::system("killall -q aplay 2>/dev/null");
}

void Audio::set_volume(int percent) {
    volume_ = std::max(0, std::min(100, percent));
    char cmd[128];
    std::snprintf(cmd, sizeof(cmd), "amixer set Master %d%% 2>/dev/null", volume_);
    std::system(cmd);
}

} // namespace lt
