MARIADB_USER=rikka
MARIADB_DATABASE=rikka
MARIADB_PASSWORD=rikka-password
MARIADB_ROOT_PASSWORD=password

GO_RIKKA_CMD_PATH=cmd/rikka
GO_RITTO_CMD_PATH=cmd/ritto
GO_RUN=go run $(GO_RIKKA_CMD_PATH)/main.go $(GO_RIKKA_CMD_PATH)/config.go
GO_RUN_RITTO=go run $(GO_RITTO_CMD_PATH)/main.go $(GO_RITTO_CMD_PATH)/config.go

.PHONY: run run-ritto run-docker ps up down mariadb mariadb-drop-db mariadb-create-db mariadb-reset-db redis-cli self-signed-cert-and-key
run:
	$(GO_RUN) --config $(GO_RIKKA_CMD_PATH)/config.yaml

run-ritto:
	$(GO_RUN_RITTO) --config $(GO_RITTO_CMD_PATH)/config.yaml

run-docker:
	docker compose run --rm --service-ports go $(GO_RUN) --config $(GO_RIKKA_CMD_PATH)/config.docker.yaml

ps:
	docker compose ps

up:
	docker compose up -d

down:
	docker compose down

build:
	go build -o rikka ${GO_RIKKA_CMD_PATH}/*.go

build-ritto:
	go build -o ritto ${GO_RITTO_CMD_PATH}/*.go

mariadb:
	docker compose exec mariadb mysql -u $(MARIADB_USER) --password=$(MARIADB_PASSWORD) $(MARIADB_DATABASE)

mariadb-drop-db:
	docker compose exec mariadb mysql -u root --password=$(MARIADB_ROOT_PASSWORD) -e "drop database $(MARIADB_DATABASE)"

mariadb-create-db:
	docker compose exec mariadb mysql -u root --password=$(MARIADB_ROOT_PASSWORD) -e "create database $(MARIADB_DATABASE)"

mariadb-reset-db:
	make mariadb-drop-db
	make mariadb-create-db

redis-cli:
	docker compose exec redis redis-cli

self-signed-cert-and-key:
	openssl req -x509 -nodes -new -keyout $(GO_RIKKA_CMD_PATH)/privkey.pem -out $(GO_RIKKA_CMD_PATH)/cert.pem -days 365
