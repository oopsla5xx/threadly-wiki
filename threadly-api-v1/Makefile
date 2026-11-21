-include .env
export

GOOSE_DRIVER := postgres
GOOSE_DBSTRING := $(DATABASE_URL)
GOOSE_MIGRATION_DIR := database/migrations

export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

.PHONY: help db-up db-down db-status db-logs \
        migrate-create migrate-up migrate-down migrate-down-to migrate-status \
        sqlc-generate sqlc-vet

help:
	@echo "Database:"
	@echo "  make db-up           Start PostgreSQL"
	@echo "  make db-down         Stop PostgreSQL"
	@echo "  make db-status       Show container status"
	@echo "  make db-logs         Show PostgreSQL logs"
	@echo ""
	@echo "Migrations:"
	@echo "  make migrate-create name=<name>"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-down-to version=<version>"
	@echo "  make migrate-status"
	@echo ""
	@echo "SQLC:"
	@echo "  make sqlc-generate"
	@echo "  make sqlc-vet"

# Database
db-up:
	docker compose up -d --wait postgres

db-down:
	docker compose stop postgres

db-status:
	docker compose ps

db-logs:
	docker compose logs -f postgres

# Goose migrations
migrate-create:
	goose create $(name) sql

migrate-up:
	goose up

migrate-down:
	goose down

migrate-down-to:
	goose down-to $(version)

migrate-status:
	goose status

# SQLC
sqlc-generate:
	sqlc generate

sqlc-vet:
	sqlc vet
