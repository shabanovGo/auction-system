.PHONY: up down migrate postgres recreate-db build logs

DC=docker compose
DB_USER=auction_user
DB_PASS=auction_password
DB_NAME=auction_db
DB_HOST=localhost
DB_PORT=5432
MIGRATIONS_DIR=migrations
PROJECT_NAME=auction-system
NETWORK=$(PROJECT_NAME)_auction-network

up:
	$(DC) up

down:
	$(DC) down -v

build:
	$(DC) build

logs:
	$(DC) logs -f

postgres:
	$(DC) up -d postgres
	@until docker exec $$(docker ps -q -f name=postgres) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for postgres..."; \
		sleep 1; \
	done

recreate-db: postgres
	docker exec $$(docker ps -q -f name=postgres) dropdb -U $(DB_USER) --if-exists $(DB_NAME)
	docker exec $$(docker ps -q -f name=postgres) createdb -U $(DB_USER) $(DB_NAME)

migrate: recreate-db
	@for file in $(MIGRATIONS_DIR)/*.up.sql; do \
		echo "Applying $$file..."; \
		docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done

restart-app:
	$(DC) restart app

app-logs:
	$(DC) logs -f app

db-logs:
	$(DC) logs -f postgres

start: build up migrate

.DEFAULT_GOAL := start
