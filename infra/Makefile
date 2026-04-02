.PHONY: help dev prod up down restart build \
       db-migrate db-seed db-reset db-fresh db-studio \
       logs logs-postgres logs-redis logs-backend logs-web \
       ps push clean

COMPOSE_DEV   = docker compose -f docker/docker-compose.dev.yml
COMPOSE_PROD  = docker compose -f docker/docker-compose.prod.yml

# ──────────────────────────────────────────────
#  Výchozí cíl
# ──────────────────────────────────────────────

help: ## Zobrazí nápovědu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ──────────────────────────────────────────────
#  Development
# ──────────────────────────────────────────────

dev: ## Spustí DEV prostředí (build + logs v popředí)
	$(COMPOSE_DEV) up --build

dev-d: ## Spustí DEV prostředí na pozadí
	$(COMPOSE_DEV) up --build -d

dev-down: ## Zastaví DEV prostředí
	$(COMPOSE_DEV) down

# ──────────────────────────────────────────────
#  Production
# ──────────────────────────────────────────────

prod: ## Spustí PROD prostředí (build + detached)
	$(COMPOSE_PROD) up --build -d

prod-down: ## Zastaví PROD prostředí
	$(COMPOSE_PROD) down

prod-restart: prod-down prod ## Restart PROD služeb

# ──────────────────────────────────────────────
#  Obecné
# ──────────────────────────────────────────────

up: dev-d ## Alias pro dev-d

down: ## Zastaví kontejnery (dev)
	$(COMPOSE_DEV) down

restart: down up ## Restart služeb (dev)

build: ## Sestaví všechny image (dev)
	$(COMPOSE_DEV) build

# ──────────────────────────────────────────────
#  Databáze
# ──────────────────────────────────────────────

db-migrate: ## Spustí migrace (prisma migrate deploy)
	$(COMPOSE_DEV) run --build --rm prisma-migrate

db-seed: ## Spustí seedování databáze
	$(COMPOSE_DEV) run --rm prisma-seed

db-reset: ## Smaže volume a spustí vše od nuly (migrate + seed)
	$(COMPOSE_DEV) down -v
	$(COMPOSE_DEV) up --build

db-fresh: ## Smaže volume, spustí v detach módu (migrate + seed)
	$(COMPOSE_DEV) down -v
	$(COMPOSE_DEV) up --build -d

db-studio: ## Otevře Prisma Studio (GUI pro databázi)
	$(COMPOSE_DEV) run --rm -p 5555:5555 prisma-migrate \
		sh -c "npx prisma studio --port 5555 --browser none"

# ──────────────────────────────────────────────
#  Git push (commit + push)
# ──────────────────────────────────────────────

push: ## Git add + commit + push (použij: make push m="zpráva")
	cd .. && git add -A && git commit -m "$(m)" && git push

# ──────────────────────────────────────────────
#  Logy
# ──────────────────────────────────────────────

logs: ## Zobrazí logy všech služeb (dev)
	$(COMPOSE_DEV) logs -f

logs-postgres: ## Logy PostgreSQL
	$(COMPOSE_DEV) logs -f postgres

logs-redis: ## Logy Redis
	$(COMPOSE_DEV) logs -f redis

logs-backend: ## Logy Go backend
	$(COMPOSE_DEV) logs -f backend

logs-web: ## Logy Next.js web
	$(COMPOSE_DEV) logs -f web

# ──────────────────────────────────────────────
#  Stav a úklid
# ──────────────────────────────────────────────

ps: ## Zobrazí stav běžících služeb
	$(COMPOSE_DEV) ps -a

clean: ## Zastaví vše, smaže volumes, images a sítě
	$(COMPOSE_DEV) down -v --rmi local --remove-orphans
	$(COMPOSE_PROD) down -v --rmi local --remove-orphans
