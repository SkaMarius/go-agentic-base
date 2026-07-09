MIGRATIONS_DIR ?= migrations
DATABASE_URL ?=

# Requires the golang-migrate CLI (https://github.com/golang-migrate/migrate):
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: sandbox-up sandbox-down sandbox-ports

sandbox-up:
	docker compose up --build -d

sandbox-down:
	docker compose down -v

sandbox-ports:
	docker compose port app 8080

.PHONY: install-hooks

install-hooks:
	pre-commit install
