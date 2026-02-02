.PHONY: run migrate-up migrate-down migrate-reset create-migration tidy build

DB_URL=postgres://postgres:1234@localhost:5432/dental_clinic?sslmode=disable

# Run the application
run:
	go run ./cmd/main.go

# Build the application binary
build:
	go build -o app ./cmd/main.go

# Apply all pending database migrations
migrate-up:
	goose -dir ./migrations postgres "$(DB_URL)" up

# Roll back the last applied migration
migrate-down:
	goose -dir ./migrations postgres "$(DB_URL)" down

# Roll back all migrations
migrate-reset:
	goose -dir ./migrations postgres "$(DB_URL)" down-to 0

# Create a new migration file
create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-migration name=add_field"; \
	else \
		goose -dir ./migrations create $(name) sql; \
	fi

# Cleaning
tidy:
	go mod tidy
