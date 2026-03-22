# ══════════════════════════════════════════════════
#  LaserTag MP135 — Build System
#
#  Targets:
#    make              — cross-compile for ARM (MP135)
#    make native       — native build (for dev on MP135 directly)
#    make clean        — remove build artifacts
#    make deploy       — scp binary to MP135
#    make run          — deploy & run on remote MP135
# ══════════════════════════════════════════════════

# ── Cross-compilation (ARM Cortex-A7, MP135) ──
CROSS_PREFIX ?= arm-linux-gnueabihf-
CXX_CROSS    := $(CROSS_PREFIX)g++
CXX_NATIVE   := g++

CXX       = $(CXX_CROSS)
# Include path: look for can_protocol.h in esp/ (sibling) or ../esp/ (repo layout)
CXXFLAGS  := -std=c++17 -Wall -Wextra -O2 -march=armv7-a -mfpu=neon -mfloat-abi=hard -Iesp -I../esp
LDFLAGS   := -lpthread

# ── Target binary ──
TARGET    := lasertag

# ── Sources ──
SRCS := main.cpp \
        src/can_receiver.cpp \
        src/framebuffer.cpp \
        src/animator.cpp \
        src/audio.cpp \
        src/game_state.cpp

OBJS := $(SRCS:.cpp=.o)

# ── MP135 deploy config ──
MP135_HOST  ?= mp135
MP135_USER  ?= root
MP135_DIR   ?= /opt/lasertag

# ══════════════════════════════════════════════════
.PHONY: all native clean deploy run stop help

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LDFLAGS)
	@echo "Built: $(TARGET) (cross-compiled for ARM)"

# Native build (run directly on MP135 Debian)
native: CXX = $(CXX_NATIVE)
native: CXXFLAGS = -std=c++17 -Wall -Wextra -O2 -Iesp -I../esp
native: $(TARGET)
	@echo "Built: $(TARGET) (native)"

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c -o $@ $<

# Dependencies (header changes trigger rebuild)
src/can_receiver.o: src/can_receiver.cpp src/can_receiver.h src/event_bus.h
src/framebuffer.o:  src/framebuffer.cpp src/framebuffer.h src/font.h
src/animator.o:     src/animator.cpp src/animator.h src/framebuffer.h
src/audio.o:        src/audio.cpp src/audio.h
src/game_state.o:   src/game_state.cpp src/game_state.h src/event_bus.h src/animator.h src/audio.h src/framebuffer.h
main.o:             main.cpp src/event_bus.h src/can_receiver.h src/framebuffer.h src/animator.h src/audio.h src/game_state.h

clean:
	rm -f $(TARGET) $(OBJS)

# Deploy to MP135 via SCP
deploy: $(TARGET)
	ssh $(MP135_USER)@$(MP135_HOST) "mkdir -p $(MP135_DIR)/media"
	scp $(TARGET) $(MP135_USER)@$(MP135_HOST):$(MP135_DIR)/
	scp -r media/*.wav $(MP135_USER)@$(MP135_HOST):$(MP135_DIR)/media/ 2>/dev/null || true
	@echo "Deployed to $(MP135_HOST):$(MP135_DIR)"

# Run on MP135 (deploy + execute)
run: deploy
	ssh $(MP135_USER)@$(MP135_HOST) "cd $(MP135_DIR) && ./$(TARGET) -p 1"

# Stop running instance
stop:
	ssh $(MP135_USER)@$(MP135_HOST) "killall $(TARGET) 2>/dev/null || true"

help:
	@echo "LaserTag MP135 Build System"
	@echo ""
	@echo "Usage: make [TARGET] [VARS]"
	@echo ""
	@echo "Targets:"
	@echo "  all      Cross-compile for ARM (default)"
	@echo "  native   Build natively on MP135"
	@echo "  clean    Remove build artifacts"
	@echo "  deploy   SCP binary + media to MP135"
	@echo "  run      Deploy & run on MP135"
	@echo "  stop     Kill running instance on MP135"
	@echo ""
	@echo "Variables:"
	@echo "  CROSS_PREFIX  Cross-compiler prefix (default: arm-linux-gnueabihf-)"
	@echo "  MP135_HOST    SSH hostname (default: mp135)"
	@echo "  MP135_USER    SSH user (default: root)"
	@echo "  MP135_DIR     Remote install dir (default: /opt/lasertag)"
