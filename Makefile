include .env

build:
	@echo "Building..."
	@go build -o bin/main cmd/http/*.go

clean:
	@echo "Cleaning..."
	@rm -rf bin

codegen:
	@sqlc generate

TABLE_NAME=$(wordlist 2, 2, $(MAKECMDGOALS))
migration:
	@if [ -z "$(TABLE_NAME)" ]; then \
		echo "Error: Please specify the table name (e.g., make migration my_table_name)"; \
		exit 1; \
	fi
	@echo "Creating migration..."
	@cd internal/database/migrations && goose create $(TABLE_NAME) sql

migration-down:
	@cd internal/database/migrations && goose postgres "$(DB_ADDR)" down

migration-up:
	@cd internal/database/migrations && goose postgres "$(DB_ADDR)" up

run:
	@go run cmd/http/*.go


.PHONY:  build clean codegen migration migration-down migration-up run

