# 🔍 Backend Audit — Enterprise-Grade Review

> **Datum:** 2026-04-01
> **Scope:** `backend/` — Go REST API, WebSocket hub, MQTT integration
> **Standard:** Clean Architecture, SOLID, enterprise-grade production readiness

---

## Shrnutí

Backend má **dobrý základ** — rozumnou adresářovou strukturu (`domain/usecase/delivery/repository/infrastructure`), oddělené vrstvy a definované doménové interfejsy pro repositáře. Nicméně pro skutečný enterprise-grade kód chybí řada **kritických** prvků: nulové pokrytí testy, porušení dependency rule, race conditions, bezpečnostní díry a absence observability.

**Celkové hodnocení: 4/10** — solidní prototyp, ale daleko od production-ready.

---

## Obsah

1. [🏗️ Architektura &amp; Clean Architecture](#1--architektura--clean-architecture)
2. [🧪 Testy](#2--testy)
3. [🔒 Bezpečnost](#3--bezpečnost)
4. [⚠️ Error Handling](#4--error-handling)
5. [🔀 Concurrency &amp; Race Conditions](#5--concurrency--race-conditions)
6. [📊 Observability &amp; Logging](#6--observability--logging)
7. [🗄️ Databáze &amp; Transakce](#7--databáze--transakce)
8. [🌐 WebSocket Hub](#8--websocket-hub)
9. [📡 MQTT Client](#9--mqtt-client)
10. [🎮 Game Logic &amp; Business Rules](#10--game-logic--business-rules)
11. [⚙️ Konfigurace](#11--konfigurace)
12. [📦 Deployment &amp; DevOps](#12--deployment--devops)
13. [📝 Code Quality &amp; Tooling](#13--code-quality--tooling)
14. [🚀 Chybějící Enterprise Features](#14--chybějící-enterprise-features)
15. [✅ Akční plán](#15--akční-plán)

---

## 1. 🏗️ Architektura & Clean Architecture

### Co je dobře ✅

- Adresářová struktura respektuje vrstvy: `domain → usecase → delivery/repository`
- Domain interface pro repositáře (`domain/repositories.go`) — správně definovaný kontrakt
- Doménové entity a errory žijí v `domain/` — správný přístup
- DI container jako single wiring point

### Co porušuje Clean Architecture ❌

#### 1.1 Use case vrstva nemá interface (porty)

**Soubory:** `internal/delivery/http/*.go`, `internal/infrastructure/mqtt/client.go`, `internal/di/container.go`

Handlery a MQTT client závisí na **konkrétních structech** (`*usecase.AuthUseCase`, `*usecase.GameUseCase`, ...), ne na interfejsech. Clean Architecture vyžaduje, aby vnější vrstvy závisely na **portech** (interfejsech) definovaných v doméně nebo use case vrstvě.

```go
// ❌ Aktuálně — handler závisí na konkrétní implementaci
type GameHandler struct {
    gameUC *usecase.GameUseCase  // konkrétní struct
    mqtt   *mqtt.Client          // konkrétní struct
}

// ✅ Správně — handler závisí na interface (portu)
type GameHandler struct {
    gameUC GameUseCasePort   // interface definovaný v domain/ nebo usecase/
    mqtt   MQTTPublisher     // interface
}
```

**Dopad:** Netestovatelné bez reálné DB/MQTT, nemožné mockovat.

#### 1.2 Infrastructure závisí na use case vrstvě (Dependency Rule violation)

**Soubor:** `internal/infrastructure/mqtt/client.go`

```go
import "github.com/.../internal/usecase"  // ❌ infra → usecase = špatný směr
```

Infrastructure vrstva **nesmí** importovat use case vrstvu. MQTT client by měl volat interface definovaný v domain:

```go
// ✅ domain/ports.go
type HitProcessor interface {
    ProcessHit(ctx context.Context, gameID, attackerDeviceID, victimDeviceID string) (*HitResult, error)
}
type DeviceRegistrar interface {
    Register(ctx context.Context, deviceID string) (*Device, error)
    Heartbeat(ctx context.Context, deviceID string) error
}
```

#### 1.3 `GameFull` struct žije v `usecase/` — patří do `domain/`

**Soubor:** `internal/usecase/game_usecase.go:223`

```go
type GameFull struct { ... }  // ❌ v usecase vrstvě
```

Doménový agregát by měl být v `domain/entities.go`.

#### 1.4 `HitResult` struct v `usecase/` — patří do `domain/`

**Soubor:** `internal/usecase/hit_usecase.go:26`

Výsledek business operace by měl být doménový typ.

#### 1.5 Delivery handler referencuje infrastrukturu

**Soubor:** `internal/delivery/http/game_handler.go:8`

```go
import "github.com/.../internal/infrastructure/mqtt"  // ❌ delivery → infrastructure
```

Handler by neměl znát MQTT. Publish logiku by měl řešit use case nebo přes event bus.

---

## 2. 🧪 Testy

### Nulové pokrytí testy ❌❌❌

**Neexistuje jediný `_test.go` soubor v celém backendu.**

Toto je absolutně kritický nedostatek. Senior dev by nikdy neodevzdal kód bez:

| Typ testu                       | Stav      | Priorita     |
| ------------------------------- | --------- | ------------ |
| Unit testy use cases            | ❌ Chybí | 🔴 Kritická |
| Unit testy domain logiky        | ❌ Chybí | 🔴 Kritická |
| Integration testy repository    | ❌ Chybí | 🟡 Vysoká   |
| Integration testy HTTP handlers | ❌ Chybí | 🟡 Vysoká   |
| E2E testy                       | ❌ Chybí | 🟠 Střední |
| Benchmark testy                 | ❌ Chybí | 🔵 Nízká   |

### Co je potřeba:

```
backend/
├── internal/
│   ├── domain/
│   │   └── entities_test.go         # test DefaultGameSettings, validace
│   ├── usecase/
│   │   ├── auth_usecase_test.go     # mock repos, test login/logout/validate
│   │   ├── game_usecase_test.go     # state machine testy
│   │   ├── hit_usecase_test.go      # kill logic, friendly fire, self-hit
│   │   └── device_usecase_test.go   # register, heartbeat, mark offline
│   ├── delivery/
│   │   └── http/
│   │       ├── auth_handler_test.go
│   │       ├── game_handler_test.go
│   │       └── testutil_test.go     # shared test helpers
│   ├── repository/
│   │   └── postgres/
│   │       └── integration_test.go  # testcontainers
│   └── mocks/                       # ← NOVÝ adresář
│       ├── mock_repos.go
│       └── mock_usecases.go
```

### Doporučený přístup:

1. Definovat use case interface (porty) → automaticky mockable
2. Použít `github.com/stretchr/testify` pro assertions
3. Použít `go.uber.org/mock` (nebo ruční mocky) pro unit testy
4. Použít `testcontainers-go` pro integration testy s PostgreSQL
5. Table-driven testy pro game state machine

---

## 3. 🔒 Bezpečnost

### 3.1 CORS — Wide Open 🔴

**Soubor:** `internal/delivery/http/middleware.go:46`

```go
w.Header().Set("Access-Control-Allow-Origin", "*")  // ❌ NIKDY v produkci
```

**Fix:** Konfigurovat povolené originy z env vars:

```go
allowedOrigins := cfg.CORSOrigins  // e.g. "https://dashboard.lasertag.com"
```

### 3.2 WebSocket — Žádná validace originu 🔴

**Soubor:** `internal/delivery/ws/hub.go:13`

```go
CheckOrigin: func(r *http.Request) bool { return true }  // ❌
```

### 3.3 Žádný Rate Limiting 🔴

Login endpoint nemá rate limiting → brute force je triviální.

**Fix:** Přidat rate limiter middleware (`golang.org/x/time/rate` nebo `github.com/ulule/limiter`).

### 3.4 Žádný limit velikosti request body 🔴

**Soubor:** `internal/delivery/http/helpers.go:21`

```go
func readJSON(r *http.Request, v any) error {
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(v)  // ❌ neomezená velikost → OOM DoS
}
```

**Fix:**

```go
func readJSON(r *http.Request, v any) error {
    defer r.Body.Close()
    r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1 MB limit
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    return dec.Decode(v)
}
```

### 3.5 Session management 🟡

- Token je opaque hex string — OK pro single-server, ale neškáluje horizontálně
- `DeleteExpired()` existuje ale **nikdy se nevolá** → expired sessions se hromadí
- Session TTL (24h) je hardcoded — měl by být konfigurovatelný
- Logout nevaliduje token formát

### 3.6 Žádná RBAC autorizace 🟡

Middleware ověří autentizaci, ale nikde se nekontroluje `Role`. Admin i user mají identický přístup ke všem endpointům.

### 3.7 MQTT bez autentizace 🟡

Broker `tcp://localhost:1883` bez TLS a credentials. Jakýkoliv klient se může připojit a posílat falešné hit events.

---

## 4. ⚠️ Error Handling

### 4.1 Silentně ignorované errory 🔴

Napříč celým kódem se ignorují errory pomocí `_ =`:

**`game_usecase.go:163` (StartGame):**

```go
players, _ := uc.players.ListByGame(ctx, gameID)  // ❌ error ignorován
for _, p := range players {
    _ = uc.devices.UpdateStatus(ctx, p.DeviceID, domain.DeviceInGame)  // ❌
}
```

**`game_usecase.go:190` (EndGame):**

```go
players, _ := uc.players.ListByGame(ctx, gameID)  // ❌ stejný problém
```

**`hit_usecase.go:102` (ProcessHit):**

```go
_ = uc.events.Create(ctx, event)  // ❌ kill event se nezaloguje = data loss
```

**`game_handler.go:161` (Start):**

```go
players, _ := h.gameUC.ListPlayers(r.Context(), gameID)  // ❌
```

### 4.2 Interní chyby leakují klientovi 🟡

```go
writeError(w, http.StatusInternalServerError, err.Error())  // ❌ stack trace/SQL chyba klientovi
```

**Fix:** Logovat interní chybu, klientovi vracet generickou zprávu:

```go
log.Printf("[ERROR] %v", err)
writeError(w, http.StatusInternalServerError, "internal server error")
```

### 4.3 Chybí error wrapping 🟡

Errory se nepropagují s kontextem:

```go
// ❌ Aktuálně
return nil, err

// ✅ Správně
return nil, fmt.Errorf("create game: %w", err)
```

### 4.4 Chybí error kódy v API responses 🟡

```json
// ❌ Aktuálně
{"error": "game not found"}

// ✅ Enterprise standard
{"error": {"code": "GAME_NOT_FOUND", "message": "Game with given ID was not found"}}
```

---

## 5. 🔀 Concurrency & Race Conditions

### 5.1 WebSocket Hub — Deadlock Risk 🔴

**Soubor:** `internal/delivery/ws/hub.go:68-80`

```go
func (h *Hub) Broadcast(msg any) {
    h.mu.RLock()         // ← drží RLock
    defer h.mu.RUnlock()
    for conn := range h.clients {
        if err := conn.WriteMessage(...); err != nil {
            conn.Close()
            go func(c *websocket.Conn) {
                h.mu.Lock()    // ← chce Write Lock zatímco hlavní goroutine drží RLock
                delete(h.clients, c)
                h.mu.Unlock()
            }(conn)
        }
    }
}
```

Spawned goroutine čeká na `Lock()` zatímco rodičovská goroutine drží `RLock()`. Pokud selže více connections najednou, `RLock` se neuvolní dokud neskončí iterace → goroutiny čekající na `Lock` zůstanou viset, ale `RUnlock` se nakonec zavolá. **Reálnější problém:** modifikace mapy concurrent s iterací.

**Fix:** Shromáždit failed connections a smazat je po uvolnění locku:

```go
func (h *Hub) Broadcast(msg any) {
    data, _ := json.Marshal(msg)
    h.mu.RLock()
    var failed []*websocket.Conn
    for conn := range h.clients {
        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            failed = append(failed, conn)
        }
    }
    h.mu.RUnlock()

    if len(failed) > 0 {
        h.mu.Lock()
        for _, c := range failed {
            delete(h.clients, c)
            c.Close()
        }
        h.mu.Unlock()
    }
}
```

### 5.2 ProcessHit — Race condition na player stats 🔴

**Soubor:** `internal/usecase/hit_usecase.go:36-111`

Dva současné hity na stejného hráče → oba přečtou `IsAlive=true` → oba zapíší kill → duplikát. Chybí databázový lock nebo optimistic concurrency.

**Fix:** Použít `SELECT ... FOR UPDATE` v transakci, nebo atomický update:

```sql
UPDATE players SET is_alive = false, deaths = deaths + 1
WHERE id = $1 AND is_alive = true
RETURNING is_alive;
```

### 5.3 Game state transitions — No locking 🟡

`StartGame` a `EndGame` nemají žádný lock. Dva requesty mohou současně přejít z lobby→running.

**Fix:** Optimistic locking (version field) nebo `SELECT ... FOR UPDATE`.

### 5.4 Respawn goroutiny leakují při shutdownu 🟡

**Soubor:** `internal/infrastructure/mqtt/client.go:156-174`

```go
go func() {
    time.Sleep(time.Duration(game.Settings.RespawnDelay) * time.Second)  // ❌ nekanceluje se
    ...
}()
```

Při graceful shutdown se tyto goroutiny nemohou zastavit.

**Fix:** Použít `context.WithCancel` nebo `time.AfterFunc` s cancel capability.

---

## 6. 📊 Observability & Logging

### 6.1 Nestrukturovaný logging 🔴

```go
log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)  // ❌ plain text, neparsable
```

**Fix:** Použít `log/slog` (stdlib od Go 1.21) nebo `zerolog`:

```go
slog.Info("request",
    "method", r.Method,
    "path", r.URL.Path,
    "request_id", requestID,
    "duration_ms", elapsed.Milliseconds(),
    "status", status,
)
```

### 6.2 Žádný Request ID / Correlation ID 🔴

Není možné traceovat request přes layers. Každý request by měl mít unikátní ID propagovaný přes context.

### 6.3 Logging middleware neloguje response 🟡

**Soubor:** `internal/delivery/http/middleware.go:61`

Aktuálně loguje jen method+path. Chybí:

- HTTP status code
- Response time (latency)
- Request/response body size
- Client IP
- User-Agent

### 6.4 Žádné metriky 🟡

Chybí Prometheus/OpenTelemetry metriky:

- Request count/duration histogram
- Active WebSocket connections gauge
- Active games gauge
- MQTT message count
- DB connection pool stats

### 6.5 Health check je nedostatečný 🟡

```go
mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})  // ❌ jen statický string
})
```

**Fix:** Zkontrolovat DB ping, MQTT connection, Redis:

```json
{
  "status": "healthy",
  "checks": {
    "database": {"status": "up", "latency_ms": 2},
    "mqtt": {"status": "up"},
    "redis": {"status": "down", "error": "connection refused"}
  }
}
```

---

## 7. 🗄️ Databáze & Transakce

### 7.1 Žádné transakce pro multi-step operace 🔴

`ProcessHit` aktualizuje 2 hráče + vytváří event — **bez transakce**:

```go
// hit_usecase.go — 3 nezávislé DB operace, žádná transakce
uc.players.Update(ctx, victim)    // krok 1
uc.players.Update(ctx, attacker)  // krok 2 — co když selže?
uc.events.Create(ctx, event)      // krok 3 — data inconsistency
```

`StartGame` aktualizuje game + N device statusů — **bez transakce**.

**Fix:** Repository interface by měl podporovat transakce:

```go
// domain/repositories.go
type UnitOfWork interface {
    Begin(ctx context.Context) (Transaction, error)
}
type Transaction interface {
    Commit() error
    Rollback() error
    PlayerRepo() PlayerRepository
    EventRepo() EventRepository
}
```

### 7.2 Connection pool — chybí `ConnMaxLifetime` 🟡

**Soubor:** `internal/repository/postgres/db.go`

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
// ❌ Chybí:
// db.SetConnMaxLifetime(5 * time.Minute)
// db.SetConnMaxIdleTime(1 * time.Minute)
```

Bez `ConnMaxLifetime` mohou stale connections způsobit problémy s load balancery a PostgreSQL restarts.

### 7.3 `sql.ErrNoRows` — nepodporuje wrapping 🟡

```go
if err == sql.ErrNoRows {  // ❌ fragile
    return nil, domain.ErrNotFound
}
// ✅ Bezpečnější:
if errors.Is(err, sql.ErrNoRows) {
```

### 7.4 Žádná paginace na databázové vrstvě 🟡

`ListAll()` metody načítají **všechny záznamy**. S 10 000+ hrami/eventy = OOM.

```go
// ❌ Aktuálně
ListAll(ctx context.Context) ([]Game, error)

// ✅ Enterprise
ListAll(ctx context.Context, opts PaginationOpts) (PagedResult[Game], error)
```

### 7.5 Nepoužitá funkce `timePtr` 🔵

**Soubor:** `internal/repository/postgres/game_repo.go:116` — mrtvý kód.

---

## 8. 🌐 WebSocket Hub

### 8.1 Žádná autentizace na WS 🔴

**Soubor:** `internal/delivery/http/router.go:33`

```go
mux.HandleFunc("GET /ws", wsHub.HandleWS)  // ❌ public endpoint, žádný auth
```

Kdokoli se může připojit a poslouchat game events.

**Fix:** Validovat token z query parametru nebo prvního message:

```
ws://host/ws?token=abc123
```

### 8.2 Žádný ping/pong heartbeat 🟡

Stale connections se nedetekují. Client se může odpojit (network drop) a server to nepozná.

**Fix:**

```go
conn.SetPongHandler(func(string) error {
    conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    return nil
})
// + periodický ping z writePump goroutiny
```

### 8.3 Žádný message type system 🟡

Broadcast posílá `any` bez typového systému:

```go
func (h *Hub) Broadcast(msg any)
```

**Fix:** Definovat typed events:

```go
type WSEvent struct {
    Type    string `json:"type"`
    GameID  string `json:"game_id,omitempty"`
    Payload any    `json:"payload"`
}
```

### 8.4 Žádný channel/room systém 🟡

Všichni klienti dostávají všechny zprávy. Mělo by se posílat jen klientům, kteří sledují konkrétní hru.

### 8.5 Concurrent write na WebSocket connection 🔴

`Broadcast()` drží `RLock` a volá `conn.WriteMessage()`. Ale gorilla/websocket **nepovoluje concurrent writes**. Pokud se Broadcast volá z více goroutin (MQTT handler + background task), může dojít ke corrupted frames.

**Fix:** Každý connection needs vlastní write goroutinu s channel-based message queue.

---

## 9. 📡 MQTT Client

### 9.1 Hardcoded client ID 🔴

**Soubor:** `internal/infrastructure/mqtt/client.go:41`

```go
SetClientID("lasertag-backend")  // ❌ nemožné spustit více instancí
```

**Fix:** `SetClientID(fmt.Sprintf("lasertag-backend-%s", uuid.New().String()[:8]))`

### 9.2 Respawn timer — fire and forget goroutine 🟡

**Soubor:** `internal/infrastructure/mqtt/client.go:156`

```go
go func() {
    time.Sleep(...)  // nekanceluje se při game end / server shutdown
    c.hitUC.Respawn(...)
}()
```

**Problémy:**

- Hra může skončit před respawnem → respawn po konci hry
- Server shutdown → goroutine leak
- Žádné sledování aktivních respawn timerů

**Fix:** Spravovat timery v `GameUseCase` s cancellation support.

### 9.3 QoS nekonzistence 🟡

Subscribes s QoS 0 (at most once), publishes s QoS 1 (at least once). Pro game commands (die, respawn) by měl být QoS 1+ obě strany.

### 9.4 Žádná validace MQTT payload 🟡

```go
var evt hitEvent
if err := json.Unmarshal(payload, &evt); err != nil { ... }
// ❌ Žádná validace: evt.GameID != "", evt.VictimID != ""
```

### 9.5 MQTT publish neověřuje úspěšnost 🟡

```go
c.client.Publish(topic, 1, false, data)  // token ignorován, no error check
```

---

## 10. 🎮 Game Logic & Business Rules

### 10.1 Game code collision 🟡

**Soubor:** `internal/usecase/game_usecase.go:213`

```go
func generateGameCode() string {
    const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    b := make([]byte, 6)
    for i := range b {
        b[i] = chars[rand.Intn(len(chars))]
    }
    return string(b)
}
```

31^6 ≈ 887M kombinací, ale žádná kontrola unikátnosti. S rostoucím počtem her se kolize stane pravděpodobnou.

**Fix:** Retry loop s DB kontrolou:

```go
for attempts := 0; attempts < 10; attempts++ {
    code := generateGameCode()
    _, err := uc.games.GetByCode(ctx, code)
    if errors.Is(err, domain.ErrNotFound) {
        return code, nil
    }
}
return "", errors.New("failed to generate unique game code")
```

### 10.2 Score hardcoded 🟡

```go
attacker.Score += 100  // ❌ hardcoded, mělo by být v GameSettings
```

### 10.3 `RemovePlayer` neuvolňuje device 🔴

**Soubor:** `internal/usecase/game_usecase.go:136`

```go
func (uc *GameUseCase) RemovePlayer(ctx context.Context, playerID string) error {
    return uc.players.Delete(ctx, playerID)  // ❌ device zůstane ve stavu "in_game"
}
```

### 10.4 `RemoveTeam` nekontroluje stav hry 🟡

Team lze smazat i v running stavu. Hráči v tomto teamu ztratí team assignment.

### 10.5 Žádná validace game settings 🟡

`MaxPlayers: -5`, `RespawnDelay: -100`, `GameDuration: 999999999` — vše projde.

```go
// ✅ Přidat validaci v domain
func (s GameSettings) Validate() error {
    if s.MaxPlayers < 2 || s.MaxPlayers > 100 {
        return errors.New("max_players must be between 2 and 100")
    }
    if s.RespawnDelay < 0 || s.RespawnDelay > 300 {
        return errors.New("respawn_delay must be between 0 and 300")
    }
    // ...
}
```

### 10.6 `AddPlayer` nekontroluje duplikát nickname 🔵

Dva hráči ve stejné hře mohou mít stejný nickname → matoucí v leaderboardu.

---

## 11. ⚙️ Konfigurace

### 11.1 Redis je konfigurovaný ale nikde nepoužívaný 🟡

```go
RedisAddr: getEnv("REDIS_ADDR", "localhost:6379")  // ❌ mrtvá konfigurace
```

### 11.2 Žádná validace konfigurace 🟡

```go
func Load() *Config {
    return &Config{
        ServerPort: getEnv("SERVER_PORT", "8080"),
        // ❌ co když je port "banana"? Žádná validace.
    }
}
```

**Fix:**

```go
func Load() (*Config, error) {
    cfg := &Config{...}
    if err := cfg.validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    return cfg, nil
}
```

### 11.3 Hardcoded hodnoty roztroušené v kódu 🟡

| Hodnota                         | Kde                    | Mělo by být |
| ------------------------------- | ---------------------- | ------------- |
| Session TTL:`24h`             | `auth_usecase.go:42` | Config        |
| Heartbeat interval:`10s`      | `main.go:47`         | Config        |
| Offline timeout:`30s`         | `main.go:50`         | Config        |
| Auto-end check:`5s`           | `main.go:63`         | Config        |
| HTTP timeouts:`15s/60s`       | `main.go:93`         | Config        |
| Shutdown timeout:`10s`        | `main.go:110`        | Config        |
| MQTT disconnect wait:`1000ms` | `mqtt/client.go:63`  | Config        |
| Score per kill:`100`          | `hit_usecase.go:78`  | GameSettings  |

---

## 12. 📦 Deployment & DevOps

### 12.1 Dockerfile — neexistující Go verze 🔴

```dockerfile
FROM golang:1.25-alpine  # ❌ Go 1.25 neexistuje (k datu auditu)
```

### 12.2 Chybí `.dockerignore` 🟡

Bez `.dockerignore` se do build kontextu kopíruje vše (docs, .git, etc.).

```
# ✅ .dockerignore
.git
docs/
*.md
.env*
```

### 12.3 Chybí `Makefile` 🟡

```makefile
# ✅ Doporučený Makefile
.PHONY: build test lint run swagger

build:
    go build -o bin/server ./cmd/server

test:
    go test -v -race -cover ./...

lint:
    golangci-lint run

swagger:
    swag init -g cmd/server/main.go -o docs

run:
    go run ./cmd/server
```

### 12.4 Žádný `.golangci.yml` 🟡

Chybí linter konfigurace.

### 12.5 Žádné README v backendu 🟡

Chybí dokumentace pro:

- Jak spustit lokálně
- Architektura diagram
- API overview
- Environment variables
- MQTT topic schema

### 12.6 Background goroutiny nemají graceful shutdown 🟡

**Soubor:** `cmd/server/main.go:45-86`

Ticker goroutiny nemají cancel signal. Při shutdownu tick callback může běžet po zavření DB.

**Fix:**

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // signal goroutinám

go func() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            // ... work
        case <-ctx.Done():
            return
        }
    }
}()
```

---

## 13. 📝 Code Quality & Tooling

### 13.1 `math/rand` místo `crypto/rand` pro game code 🟡

`generateGameCode` používá `math/rand` — pro game kódy je to OK (Go 1.20+ auto-seeds), ale `crypto/rand` by byl bezpečnější proti prediction attacks.

### 13.2 Swagger definitions obsahují internal package paths 🟡

```yaml
definitions:
  internal_delivery_http.AddPlayerRequest:  # ❌ leakuje interní strukturu
```

### 13.3 DTOs nemají validační tagy 🟡

```go
type LoginRequest struct {
    Username string `json:"username"`  // ❌ chybí validate:"required,min=3,max=50"
    Password string `json:"password"`
}
```

### 13.4 Žádné interface pro use cases = nemožné generovat mocky 🟡

Viz bod 1.1 — fundamentální problém pro testovatelnost.

### 13.5 Chybí `go vet` / `staticcheck` v CI 🔵

---

## 14. 🚀 Chybějící Enterprise Features

| Feature                  | Stav | Popis                                         |
| ------------------------ | ---- | --------------------------------------------- |
| API versioning           | ❌   | `/api/v1/games`                             |
| Pagination               | ❌   | Limit/offset nebo cursor-based                |
| Request validation       | ❌   | Validační library (go-playground/validator) |
| Structured logging       | ❌   | `log/slog` s JSON output                    |
| Request ID tracing       | ❌   | UUID per request v headers                    |
| Metrics                  | ❌   | Prometheus endpoint                           |
| Rate limiting            | ❌   | Per-IP a per-user                             |
| Graceful degradation     | ❌   | Fallback když MQTT/Redis je down             |
| Database transactions    | ❌   | Unit of Work pattern                          |
| Cache layer              | ❌   | Redis pro game state, sessions                |
| Audit trail              | ❌   | Kdo co kdy udělal                            |
| API key auth             | ❌   | Pro device-to-server (místo/doplnění MQTT) |
| Configuration hot-reload | ❌   | Bez restartu                                  |
| Circuit breaker          | ❌   | Pro external service calls                    |
| Idempotency              | ❌   | Idempotency-Key header pro POST               |
| OpenAPI validation       | ❌   | Runtime validace proti schématu              |

---

## 15. ✅ Akční plán

### Fáze 1 — Kritické opravy (musí být před production)

1. **Definovat use case interface (porty)** — základ pro testy i Clean Architecture
2. **Přidat unit testy** — minimálně pro use cases a domain
3. **Opravit race condition** v WebSocket hub (`Broadcast`)
4. **Přidat database transakce** pro `ProcessHit`, `StartGame`, `EndGame`
5. **Omezit request body size** (`http.MaxBytesReader`)
6. **Opravit `RemovePlayer`** — uvolnit device
7. **Error handling** — neposílat interní chyby klientovi
8. **CORS** — konfigurovat z env proměnných

### Fáze 2 — Vysoká priorita

9. **Structured logging** — `log/slog` s request ID
10. **Input validace** — game settings, DTOs
11. **Rate limiting** — na login endpoint minimálně
12. **WebSocket autentizace**
13. **Graceful shutdown** pro background goroutiny
14. **Session cleanup** — spustit `DeleteExpired` jako background task
15. **Game code uniqueness** — retry loop
16. **MQTT client ID** — unique per instance

### Fáze 3 — Enterprise polish

17. **API versioning** (`/api/v1/`)
18. **Pagination** na list endpointy
19. **Prometheus metriky**
20. **Makefile + `.golangci.yml`**
21. **Integration testy** s testcontainers
22. **Cache layer** (Redis) pro game state
23. **Audit logging**
24. **Backend README**

### Fáze 4 — Pokročilé

25. **Event-driven architecture** — decouple MQTT handler od use cases
26. **CQRS** pro read-heavy leaderboard endpoint
27. **OpenTelemetry tracing**
28. **Circuit breaker** pro external dependencies
29. **Horizontal scaling** support (unique MQTT client IDs, shared sessions)
30. **Load testing** s k6 nebo vegeta

---

## Appendix: Navrhovaná cílová adresářová struktura

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entity/           # Game, Player, Team, Device, User, Session, GameEvent
│   │   ├── valueobject/      # GameCode, DeviceID, Score
│   │   ├── port/             # ← NOVÉ: UseCase interfaces + Repository interfaces
│   │   │   ├── repository.go # UserRepository, GameRepository, ...
│   │   │   └── usecase.go    # GameUseCasePort, AuthUseCasePort, ...
│   │   └── error.go
│   ├── usecase/
│   │   ├── auth.go
│   │   ├── game.go
│   │   ├── hit.go
│   │   └── device.go
│   ├── adapter/
│   │   ├── handler/          # HTTP handlers (delivery)
│   │   │   ├── auth.go
│   │   │   ├── game.go
│   │   │   ├── device.go
│   │   │   ├── middleware.go
│   │   │   ├── router.go
│   │   │   └── dto/
│   │   ├── ws/               # WebSocket
│   │   ├── mqtt/             # MQTT adapter
│   │   └── repository/       # PostgreSQL implementations
│   │       └── postgres/
│   ├── config/
│   └── pkg/                  # Shared utilities (logger, validator, etc.)
├── test/
│   ├── integration/
│   ├── e2e/
│   └── mock/
├── docs/
├── Dockerfile
├── .dockerignore
├── .golangci.yml
├── Makefile
└── README.md
```

---

*Audit provedl AI code reviewer. Doporučuje se přezkoumání senior Go vývojářem.*
