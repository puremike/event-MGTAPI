include .env
export $(shell sed 's/=.*//' .env)

.PHONY: migrate-up migrate-down

migrate-up:
	@echo "Migrating up"
	migrate -database "$(DB_URL)" -path cmd/migrate/migrations up

migrate-down:
	@echo "Migrating down"
	migrate -database "$(DB_URL)" -path cmd/migrate/migrations down

