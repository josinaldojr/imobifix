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

test-integration:
	@set -a; [ -f .env.test ] && . ./.env.test; set +a; \
	go test -tags=integration ./... -v

cover:
	go test ./... -coverpkg=./... -coverprofile=cover.out
	go tool cover -func=cover.out

cover-html:
	go test ./... -coverpkg=./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	@echo "Gerado: $(CURDIR)/cover.html"
	@command -v open >/dev/null && open cover.html || true
	@command -v xdg-open >/dev/null && xdg-open cover.html || true

fmt:
	cd imobifix-api && gofmt -w .

migrate-up:
	docker run --rm \
	  -v $(PWD)/imobifix-api/migrations:/migrations \
	  --network imobifx_default \
	  migrate/migrate:v4.17.1 \
	  -path=/migrations \
	  -database "postgres://imobifx:imobifx@postgres:5432/imobifx?sslmode=disable" up