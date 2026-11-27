.PHONY: up down logs build test api

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

build:
	docker compose build

test:
	cd backend && go test ./...
	cd backend && go vet ./...

api:
	curl -fsS http://localhost:8080/api/health
