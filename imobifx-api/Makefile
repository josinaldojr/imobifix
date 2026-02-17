.PHONY: up down logs ps \
	test test-integration cover cover-html fmt \
	migrate-up \
	e2e-up e2e-down e2e \
	e2e-logs e2e-ps

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)

COMPOSE := docker compose --project-directory $(ROOT_DIR) \
	-f $(ROOT_DIR)/docker-compose.yml

COMPOSE_E2E := docker compose --project-directory $(ROOT_DIR) \
	-f $(ROOT_DIR)/docker-compose.yml \
	-f $(ROOT_DIR)/docker-compose.e2e.yml

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f imobifx-api

ps:
	$(COMPOSE) ps

test:
	go test ./... -v

test-integration:
	@set -a; [ -f .env.test ] && . ./.env.test; set +a; \
	go test -tags=integration ./... -v

cover:
	go test ./... -coverpkg=./... -coverprofile=cover.out
	go tool cover -func=cover.out | tail -n 1

cover-html:
	go test ./... -coverpkg=./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	@echo "Gerado: $(CURDIR)/cover.html"
	@command -v open >/dev/null && open cover.html || true
	@command -v xdg-open >/dev/null && xdg-open cover.html || true

fmt:
	gofmt -w .

migrate-up:
	$(COMPOSE) run --rm migrate

e2e-up:
	$(COMPOSE_E2E) up -d --build --force-recreate

e2e-down:
	$(COMPOSE_E2E) down -v

e2e-logs:
	$(COMPOSE_E2E) logs -f imobifx-api viacep-mock

e2e-ps:
	$(COMPOSE_E2E) ps
	
e2e:
	@set -e; \
	$(MAKE) e2e-up; \
	trap '$(MAKE) e2e-down' EXIT; \
	set -a; [ -f .env.test ] && . ./.env.test; set +a; \
	go test -tags=e2e ./test/e2e -v