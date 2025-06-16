include .envrc

MIGRATION_NAME =

.PHONY: run
run:
	go run . -port=${PORT} -env=${ENV} -db-url=${DB_URL}

.PHONY: start-db
start-db:
	docker-compose up

.PHONY: migrate/create
migrate/create:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Error: MIGRATION_NAME is required."; \
		exit 1; \
	fi
	goose -dir ./sql/migrations -s create $(MIGRATION_NAME) sql

.PHONY: migrate/up
migrate/up:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL is required."; \
		exit 1; \
	fi
	goose -dir ./sql/migrations postgres $(DB_URL) up 