.PHONY: help setup start stop restart status logs \
       logs-nginx logs-backend logs-web logs-postgres logs-redis \
       update clean ps db-backup db-restore

COMPOSE = docker compose -f ../docker/docker-compose.yml -f ../docker/docker-compose.prod.yml

# ──────────────────────────────────────────────
#  Nápověda
# ──────────────────────────────────────────────

help: ## Zobrazí nápovědu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ──────────────────────────────────────────────
#  Prvotní nastavení
# ──────────────────────────────────────────────

setup: ## Prvotní setup RPi (sudo)
	sudo bash setup.sh

# ──────────────────────────────────────────────
#  Životní cyklus
# ──────────────────────────────────────────────

start: ## Spustí produkční server
	$(COMPOSE) --env-file .env up --build -d

stop: ## Zastaví server
	$(COMPOSE) --env-file .env down

restart: stop start ## Restart serveru

update: ## Stáhne poslední kód a restartuje
	cd ../.. && git pull
	$(COMPOSE) --env-file .env up --build -d --force-recreate

# ──────────────────────────────────────────────
#  Monitoring
# ──────────────────────────────────────────────

status: ## Stav všech služeb
	$(COMPOSE) ps -a

ps: status ## Alias pro status

logs: ## Logy všech služeb
	$(COMPOSE) logs -f --tail=100

logs-nginx: ## Logy Nginx
	$(COMPOSE) logs -f nginx

logs-backend: ## Logy Go backendu
	$(COMPOSE) logs -f backend

logs-web: ## Logy Next.js webu
	$(COMPOSE) logs -f web

logs-postgres: ## Logy PostgreSQL
	$(COMPOSE) logs -f postgres

logs-redis: ## Logy Redis
	$(COMPOSE) logs -f redis

# ──────────────────────────────────────────────
#  Údržba
# ──────────────────────────────────────────────

clean: ## Zastaví, smaže volumes a images
	$(COMPOSE) --env-file .env down -v --rmi local --remove-orphans
	docker system prune -f

db-backup: ## Záloha databáze do souboru
	@mkdir -p backups
	$(COMPOSE) exec postgres pg_dump -U $${POSTGRES_USER:-lasertag} $${POSTGRES_DB:-lasertag} \
		| gzip > backups/lasertag_$$(date +%Y%m%d_%H%M%S).sql.gz
	@echo "Záloha uložena do backups/"

db-restore: ## Obnova databáze ze zálohy (použij: make db-restore f=soubor.sql.gz)
	@test -n "$(f)" || (echo "Použití: make db-restore f=backups/soubor.sql.gz" && exit 1)
	gunzip -c $(f) | $(COMPOSE) exec -T postgres psql -U $${POSTGRES_USER:-lasertag} $${POSTGRES_DB:-lasertag}
