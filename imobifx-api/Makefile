.PHONY: up down logs ps test fmt imobifix-api-test frontend-test

up:
	docker compose up --build

down:
	docker compose down -v

logs:
	docker compose logs -f imobifix-api

ps:			
	docker compose ps

test:
	go test ./... -v

fmt:
	cd imobifix-api && gofmt -w .

migrate-up:
	docker run --rm \
	  -v $(PWD)/imobifix-api/migrations:/migrations \
	  --network imobifx_default \
	  migrate/migrate:v4.17.1 \
	  -path=/migrations \
	  -database "postgres://imobifx:imobifx@postgres:5432/imobifx?sslmode=disable" up