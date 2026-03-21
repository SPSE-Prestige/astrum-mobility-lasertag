# Laser Tag — Raspberry Pi Production Server

Produkční nasazení laser-tag serveru na Raspberry Pi s Ubuntu.
Používá standardní `docker/docker-compose.yml` + `docker/docker-compose.prod.yml` (s Nginx).

## Architektura

```
        ┌──────────────────────┐
   :80  │       Nginx          │
        │   (reverse proxy)    │
        └──┬────────┬──────────┘
           │        │
    /api/  │        │  /
    /ws/   │        │  /_next/
           ▼        ▼
     ┌──────────┐  ┌──────────┐
     │ Backend  │  │   Web    │
     │ Go :8080 │  │ Next.js  │
     └────┬─────┘  │  :3000   │
          │        └──────────┘
    ┌─────┴─────┐
    │           │
┌───▼───┐  ┌───▼───┐
│Postgres│  │ Redis │
│ :5432  │  │ :6379 │
└────────┘  └───────┘
```

## Požadavky

- Raspberry Pi 4/5 (min. 2 GB RAM, doporučeno 4 GB)
- Ubuntu Server 22.04+ (arm64)
- Připojení k síti

## Rychlý start

### 1. Klonování repozitáře

```bash
git clone https://github.com/SPSE-Prestige/aimtec2026-lasertag.git
cd aimtec2026-lasertag/infra/rpi
```

### 2. Prvotní setup

```bash
make setup
# Nebo ručně:
# sudo bash setup.sh
```

Setup provede:
- Aktualizaci systému
- Instalaci Dockeru a potřebných nástrojů
- Konfiguraci firewallu (UFW) — povolí porty 22, 80, 1883
- Nastavení Fail2Ban pro ochranu SSH
- Vytvoření swap souboru (2 GB)
- Optimalizaci kernelu pro kontejnery
- Systemd službu pro automatický start

### 3. Konfigurace

```bash
# Uprav .env soubor — POVINNĚ změň hesla!
nano .env

# Vygeneruj JWT secret:
openssl rand -base64 32
```

### 4. Spuštění

```bash
make start
```

Server bude dostupný na `http://<IP_ADRESA_RPI>`.

## Příkazy (Makefile)

| Příkaz | Popis |
|---|---|
| `make help` | Zobrazí nápovědu |
| `make setup` | Prvotní setup RPi |
| `make start` | Spustí server |
| `make stop` | Zastaví server |
| `make restart` | Restart serveru |
| `make update` | Git pull + restart |
| `make status` | Stav služeb |
| `make logs` | Logy všech služeb |
| `make logs-nginx` | Logy Nginx |
| `make logs-backend` | Logy backendu |
| `make db-backup` | Záloha databáze |
| `make db-restore f=soubor.sql.gz` | Obnova databáze |
| `make clean` | Smazání všeho |

## Automatický start

Po setupu se server spustí automaticky po restartu RPi díky systemd službě:

```bash
# Stav služby
sudo systemctl status lasertag

# Ruční ovládání
sudo systemctl start lasertag
sudo systemctl stop lasertag
sudo systemctl restart lasertag
```

## Zálohy

```bash
# Vytvoření zálohy
make db-backup

# Obnova ze zálohy
make db-restore f=backups/lasertag_20260321_120000.sql.gz
```

## Řešení problémů

```bash
# Všechny logy
make logs

# Stav kontejnerů
make status

# Restart po problémech
make restart

# Kompletní reset (smaže data!)
make clean && make start
```
