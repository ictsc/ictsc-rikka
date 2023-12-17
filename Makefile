#!make
include .env.dev

GO_RIKKA_CMD_PATH=cmd/rikka
GO_RITTO_CMD_PATH=cmd/ritto
GO_RUN=go run $(GO_RIKKA_CMD_PATH)/main.go $(GO_RIKKA_CMD_PATH)/config.go
GO_RUN_RITTO=go run $(GO_RITTO_CMD_PATH)/main.go $(GO_RITTO_CMD_PATH)/config.go

DOCKER_COMPOSE=docker compose --env-file .env.dev

.PHONY: run
run:
	$(GO_RUN) --config $(GO_RIKKA_CMD_PATH)/config.yaml

.PHONY: run-ritto
run-ritto:
	$(GO_RUN_RITTO) --config $(GO_RITTO_CMD_PATH)/config.yaml

.env.dev:
	if [ ! -f .env.dev ]; then cp .env .env.dev; fi

.PHONY: ps
ps: .env.dev
	$(DOCKER_COMPOSE) ps

.PHONY: up
up: .env.dev
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down: .env.dev
	$(DOCKER_COMPOSE) down

.PHONY: build
build:
	go build -o rikka ${GO_RIKKA_CMD_PATH}/*.go

.PHONY: test
test:
	go test ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: build-ritto
build-ritto:
	go build -o ritto ${GO_RITTO_CMD_PATH}/*.go

.PHONY: mariadb
mariadb:
	docker compose exec mariadb mysql -u $(MARIADB_USER) --password=$(MARIADB_PASSWORD) $(MARIADB_DATABASE)

.PHONY: mariadb-drop-db
mariadb-drop-db:
	docker compose exec mariadb mysql -u root --password=$(MARIADB_ROOT_PASSWORD) -e "drop database $(MARIADB_DATABASE)"

.PHONY: mariadb-create-db
mariadb-create-db:
	docker compose exec mariadb mysql -u root --password=$(MARIADB_ROOT_PASSWORD) -e "create database $(MARIADB_DATABASE)"

.PHONY: mariadb-reset-db
mariadb-reset-db:
	make mariadb-drop-db
	make mariadb-create-db

.PHONY: redis-cli
redis-cli:
	docker compose exec redis redis-cli

.PHONY: self-signed-cert-and-key
self-signed-cert-and-key:
	openssl req -x509 -nodes -new -keyout $(GO_RIKKA_CMD_PATH)/privkey.pem -out $(GO_RIKKA_CMD_PATH)/cert.pem -days 365
