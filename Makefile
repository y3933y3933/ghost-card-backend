MIGRATION_NAME =

.PHONY: run
run:
	go run .

.PHONY: start-db
start-db:
	docker-compose up

.PHONY: migrate/create
migrate/create:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Error: MIGRATION_NAME is required."; \
		exit 1; \
	fi
	goose -dir ./migrations -s create $(MIGRATION_NAME) sql
