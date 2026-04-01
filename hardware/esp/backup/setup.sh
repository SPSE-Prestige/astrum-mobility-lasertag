#!/usr/bin/env bash
set -euo pipefail

echo "=== LaserTag ESP adapter dependencies ==="

apt-get update -qq
apt-get install -y --no-install-recommends \
    build-essential \
    g++ \
    make \
    libmosquitto-dev \
    can-utils \
    iproute2

echo "=== Done. Run 'make' to build liblt.a ==="