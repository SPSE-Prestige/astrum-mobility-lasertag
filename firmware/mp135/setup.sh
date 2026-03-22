#!/bin/bash
# ══════════════════════════════════════════════════
#  LaserTag MP135 — Setup Script
#
#  Run this on the MP135 Debian system to install
#  all required dependencies and configure CAN bus.
#
#  Usage: sudo bash setup.sh
# ══════════════════════════════════════════════════

set -euo pipefail

echo "╔══════════════════════════════════════╗"
echo "║   LASERTAG MP135 SETUP              ║"
echo "╚══════════════════════════════════════╝"

# ── Install build tools ──
echo "[1/5] Installing build dependencies..."
apt-get update -qq
apt-get install -y -qq \
    build-essential \
    g++ \
    make \
    can-utils \
    iproute2 \
    alsa-utils \
    > /dev/null

echo "[2/5] Configuring CAN bus (can1, 1Mbps)..."
ip link set can1 down 2>/dev/null || true
ip link set can1 up type can bitrate 1000000 restart-ms 100 loopback off
echo "  can1 UP at 1Mbps"

echo "[3/5] Setting up audio..."
# Ensure ALSA is configured
if command -v amixer &>/dev/null; then
    amixer set Master 80% 2>/dev/null || true
    echo "  Volume set to 80%"
else
    echo "  amixer not found, skipping volume setup"
fi

echo "[4/5] Creating application directory..."
mkdir -p /lasertag/media
echo "  /lasertag ready"

echo "[5/5] Setting up systemd service..."
cat > /etc/systemd/system/lasertag.service << 'EOF'
[Unit]
Description=LaserTag MP135 Game Display
After=network.target

[Service]
Type=simple
WorkingDirectory=/lasertag
ExecStartPre=/sbin/ip link set can1 down
ExecStartPre=/sbin/ip link set can1 up type can bitrate 1000000 restart-ms 100 loopback off
ExecStart=/lasertag/lasertag -c can1 -f /dev/fb0
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable lasertag
echo "  Service installed and enabled (starts on boot)"

echo ""
echo "Setup complete! Build and deploy with:"
echo "  make native      # (build on MP135)"
echo "  make deploy      # (cross-compile & scp)"
echo "  sudo ./lasertag  # (run directly)"
