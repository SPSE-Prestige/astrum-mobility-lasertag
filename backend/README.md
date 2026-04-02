# Laser Tag Backend

Go backend server for the Astrum Mobility laser tag system.

## Architecture

Clean Architecture with four layers:

```
cmd/server        → Entry point, graceful shutdown
internal/
  domain/         → Entities, errors, repository & use-case port interfaces
  config/         → ENV-based configuration with validation
  repository/     → PostgreSQL repository implementations
  usecase/        → Business logic (depends only on domain ports)
  delivery/
    http/         → REST handlers, middleware, DTOs, router
    ws/           → WebSocket hub with per-client goroutines
  infrastructure/
    mqtt/         → MQTT client for ESP32 device communication
  di/             → Dependency injection container
```

## Quick Start

```bash
# Prerequisites: Go 1.24+, PostgreSQL, MQTT broker

# Copy environment template and configure
cp .env.example .env

# Build & run
make run

# Or with Docker
make docker-run
```

## Available Make Targets

| Target         | Description                     |
| -------------- | ------------------------------- |
| `make build`   | Compile binary to `./bin/`      |
| `make run`     | Build and run                   |
| `make test`    | Run tests with race detector    |
| `make lint`    | Run golangci-lint               |
| `make fmt`     | Format code                     |
| `make vet`     | Run go vet                      |
| `make tidy`    | Tidy modules                    |
| `make clean`   | Remove build artifacts          |

## Environment Variables

| Variable                  | Default     | Description                        |
| ------------------------- | ----------- | ---------------------------------- |
| `PORT`                    | `8080`      | HTTP server port                   |
| `DATABASE_URL`            | —           | PostgreSQL connection string       |
| `MQTT_BROKER_URL`         | —           | MQTT broker address                |
| `CORS_ALLOWED_ORIGINS`    | `*`         | Comma-separated allowed origins    |
| `SESSION_TTL_HOURS`       | `24`        | Auth session lifetime in hours     |
| `RATE_LIMIT_RPS`          | `5`         | Global rate limit (req/sec)        |
| `RATE_LIMIT_BURST`        | `10`        | Rate limit burst size              |

## API

Swagger UI available at `/swagger/` when the server is running.

### Auth
- `POST /api/auth/login` — Login (public)
- `POST /api/auth/logout` — Logout (auth required)

### Devices
- `GET /api/devices` — List all devices
- `GET /api/devices/available` — List available devices

### Games
- `POST /api/games` — Create game
- `GET /api/games` — List games
- `GET /api/games/{id}` — Get game detail
- `GET /api/games/{id}/full` — Get full game state
- `PATCH /api/games/{id}/settings` — Update settings
- `POST /api/games/{id}/start` — Start game
- `POST /api/games/{id}/end` — End game

### Teams & Players
- CRUD operations under `/api/games/{id}/teams` and `/api/games/{id}/players`

### WebSocket
- `GET /ws` — Real-time game events
