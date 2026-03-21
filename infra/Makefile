.PHONY: help dev prod up down restart build \
       db-migrate db-seed db-reset db-fresh db-studio \
       logs logs-postgres logs-redis \
       ps push clean

# ──────────────────────────────────────────────
#  Výchozí cíl
# ──────────────────────────────────────────────

help: ## Zobrazí nápovědu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ──────────────────────────────────────────────
#  Dev / Prod
# ──────────────────────────────────────────────

dev: ## Spustí infra v dev módu (build + logs v popředí)
	docker compose up --build

prod: ## Spustí infra v prod módu (build + detached)
	docker compose up --build -d

up: ## Spustí služby (bez rebuildu)
	docker compose up -d

down: ## Zastaví a odstraní kontejnery
	docker compose down

restart: down up ## Restart služeb

build: ## Sestaví všechny image
	docker compose build

# ──────────────────────────────────────────────
#  Databáze
# ──────────────────────────────────────────────

db-migrate: ## Spustí migrace (prisma migrate deploy)
	docker compose run --rm prisma-migrate

db-seed: ## Spustí seedování databáze
	docker compose run --rm prisma-seed

db-reset: ## Smaže volume a spustí vše od nuly (migrate + seed)
	docker compose down -v
	docker compose up --build

db-fresh: ## Smaže volume, spustí v detach módu (migrate + seed)
	docker compose down -v
	docker compose up --build -d

db-studio: ## Otevře Prisma Studio (GUI pro databázi)
	docker compose run --rm -p 5555:5555 prisma-migrate \
		sh -c "npx prisma studio --port 5555 --browser none"

# ──────────────────────────────────────────────
#  Git push (commit + push)
# ──────────────────────────────────────────────

push: ## Git add + commit + push (použij: make push m="zpráva")
	cd .. && git add -A && git commit -m "$(m)" && git push

# ──────────────────────────────────────────────
#  Logy
# ──────────────────────────────────────────────

logs: ## Zobrazí logy všech služeb
	docker compose logs -f

logs-postgres: ## Logy PostgreSQL
	docker compose logs -f postgres

logs-redis: ## Logy Redis
	docker compose logs -f redis

# ──────────────────────────────────────────────
#  Stav a úklid
# ──────────────────────────────────────────────

ps: ## Zobrazí stav běžících služeb
	docker compose ps -a

clean: ## Zastaví vše, smaže volumes, images a sítě
	docker compose down -v --rmi local --remove-orphans
