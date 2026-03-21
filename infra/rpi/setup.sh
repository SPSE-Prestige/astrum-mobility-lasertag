#!/usr/bin/env bash
## ─────────────────────────────────────────────
##  Laser Tag — Prvotní setup RPi (Ubuntu)
##  Spuštění:  chmod +x setup.sh && sudo ./setup.sh
## ─────────────────────────────────────────────
set -euo pipefail

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }

# ── Kontrola root ──
if [[ $EUID -ne 0 ]]; then
    echo "Spusť jako root: sudo ./setup.sh"
    exit 1
fi

info "=== Laser Tag RPi Setup ==="

# ── 1. Aktualizace systému ──
info "Aktualizace systému..."
apt-get update -y
apt-get upgrade -y

# ── 2. Instalace základních nástrojů ──
info "Instalace základních balíčků..."
apt-get install -y \
    curl \
    git \
    htop \
    ufw \
    fail2ban \
    unattended-upgrades \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    make

# ── 3. Instalace Dockeru ──
if command -v docker &>/dev/null; then
    info "Docker je již nainstalován: $(docker --version)"
else
    info "Instalace Dockeru..."
    curl -fsSL https://get.docker.com | sh
    # Přidání aktuálního uživatele do skupiny docker
    SUDO_USER_NAME="${SUDO_USER:-$(whoami)}"
    if [[ "$SUDO_USER_NAME" != "root" ]]; then
        usermod -aG docker "$SUDO_USER_NAME"
        info "Uživatel $SUDO_USER_NAME přidán do skupiny docker"
    fi
fi

# ── 4. Docker Compose plugin (součást moderního Dockeru) ──
if docker compose version &>/dev/null; then
    info "Docker Compose plugin: $(docker compose version)"
else
    info "Instalace Docker Compose plugin..."
    apt-get install -y docker-compose-plugin
fi

# ── 5. Firewall (UFW) ──
info "Konfigurace firewallu..."
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh          # SSH (port 22)
ufw allow 80/tcp       # HTTP (Nginx)
ufw allow 1883/tcp     # MQTT (pro laser-tag zařízení)
echo "y" | ufw enable
ufw status verbose

# ── 6. Fail2Ban ──
info "Konfigurace Fail2Ban..."
cat > /etc/fail2ban/jail.local <<'EOF'
[DEFAULT]
bantime  = 1h
findtime = 10m
maxretry = 5

[sshd]
enabled = true
port    = ssh
filter  = sshd
logpath = /var/log/auth.log
EOF
systemctl enable fail2ban
systemctl restart fail2ban

# ── 7. Automatické bezpečnostní aktualizace ──
info "Zapínám automatické bezpečnostní aktualizace..."
cat > /etc/apt/apt.conf.d/20auto-upgrades <<'EOF'
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
APT::Periodic::AutocleanInterval "7";
EOF

# ── 8. Swap (RPi má málo RAM) ──
SWAP_SIZE="2G"
if [[ ! -f /swapfile ]]; then
    info "Vytváření swapu ($SWAP_SIZE)..."
    fallocate -l "$SWAP_SIZE" /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
else
    info "Swap už existuje."
fi

# ── 9. Optimalizace pro Docker na RPi ──
info "Optimalizace kernelu pro kontejnery..."
cat > /etc/sysctl.d/99-lasertag.conf <<'EOF'
# Lepší výkon sítě
net.core.somaxconn = 1024
net.ipv4.tcp_tw_reuse = 1

# Menší swappiness (RPi)
vm.swappiness = 10

# Zvýšení limitu inotify (pro Docker)
fs.inotify.max_user_watches = 65536
fs.inotify.max_user_instances = 256
EOF
sysctl --system

# ── 10. Docker daemon config ──
info "Konfigurace Docker daemonu..."
mkdir -p /etc/docker
cat > /etc/docker/daemon.json <<'EOF'
{
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "10m",
        "max-file": "3"
    },
    "storage-driver": "overlay2",
    "live-restore": true
}
EOF
systemctl restart docker

# ── 11. Systemd service pro auto-start ──
info "Vytváření systemd služby..."
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_CMD="/usr/bin/docker compose -f ${PROJECT_DIR}/../docker/docker-compose.yml -f ${PROJECT_DIR}/../docker/docker-compose.prod.yml --env-file ${PROJECT_DIR}/.env"
cat > /etc/systemd/system/lasertag.service <<EOF
[Unit]
Description=Laser Tag Production Server
Requires=docker.service
After=docker.service network-online.target
Wants=network-online.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${PROJECT_DIR}
ExecStart=${COMPOSE_CMD} up --build -d
ExecStop=${COMPOSE_CMD} down
ExecReload=${COMPOSE_CMD} up --build -d --force-recreate
TimeoutStartSec=300

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable lasertag.service

info "=== Setup dokončen! ==="
echo ""
echo "Další kroky:"
echo "  1. Uprav .env soubor — změň hesla a JWT_SECRET"
echo "     Vygeneruj JWT_SECRET: openssl rand -base64 32"
echo "  2. Spusť server: make start"
echo "  3. Nebo restartuj RPi — server se spustí automaticky"
echo ""
echo "Server bude dostupný na: http://$(hostname -I | awk '{print $1}')"
echo ""
if [[ "${SUDO_USER:-}" != "" ]]; then
    warn "Odhlásit se a znovu přihlásit (aby fungovala skupina docker bez sudo)"
fi
