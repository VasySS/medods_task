
MAIN_FILE := ./cmd/server/main.go
PG_URL := postgres://postgres:postgres@localhost:5432/medods?sslmode=disable
DOCKER_YML := ./docker-compose.yml

run:
	go run ${MAIN_FILE}

test:
	go test ./internal/usecase/... --covermode=atomic

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "${PG_URL}" create hw-migration sql

goose-up:
	goose -dir ./migrations postgres "${PG_URL}" up

goose-down:
	goose -dir ./migrations postgres "${PG_URL}" down-to 0

goose-status:
	goose -dir ./migrations postgres "${PG_URL}" status

squawk:
	squawk ./migrations/*

swag-install:
	go get -u github.com/swaggo/swag  
	go install github.com/swaggo/swag/cmd/swag@latest

swag-gen:
	swag init --parseDependency --parseInternal -g ${MAIN_FILE}

compose-up:
	docker-compose -f $(DOCKER_YML) up -d

compose-down:
	docker-compose -f $(DOCKER_YML) down

compose-stop:
	docker-compose -f $(DOCKER_YML) stop

compose-start:
	docker-compose -f $(DOCKER_YML) start

compose-ps:
	docker-compose -f $(DOCKER_YML) ps

.PHONY: run test goose-install goose-add goose-up goose-down goose-status squawk \
	compose-up compose-down compose-stop compose-start compose-ps
