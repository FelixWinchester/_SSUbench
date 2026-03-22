.PHONY: run build test lint migrate-up migrate-down docker-up docker-down

# ── App ──────────────────────────────────────────────
run:
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

test:
	go test ./...

lint:
	golangci-lint run ./...

# ── Docker ───────────────────────────────────────────
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# ── Migrations ───────────────────────────────────────
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/ssubench?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/ssubench?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)