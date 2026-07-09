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
