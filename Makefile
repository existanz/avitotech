include .env

# Build the application
build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Run application from docker
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi
# Create DB container
docker-db:
	@if docker compose -f docker-compose-db.yml up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f docker-compose-db.yml up --build; \
	fi

# Run goose migrations up
migrate-up:
	goose -dir ./migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} search_path=${DB_SCHEMA} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} search_path=${DB_SCHEMA} sslmode=disable" down
