.PHONY: start build lint test format db-init db-migration-create db-seed-create db-migration-run db-seed-run

format:
	@gofmt -e -s -w -l ./

lint:
	@golangci-lint run -v ./... --timeout 3m0s

start:
	@docker-compose --env-file .env up -d

build_api:
	@go build -o build/api cmd/api.go

build_db_cmd:
	@go build -o build/db cmd/db/main.go

test:
	@docker-compose --env-file .env.test -f docker-compose.test.yml up -d
	@ENVIRONMENT=test go run cmd/db/main.go init
	@PROJECT_ROOT=$(PWD) ENVIRONMENT=test go run github.com/onsi/ginkgo/v2/ginkgo -r ./... --race -coverpkg=./internal/...

ENVIRONMENT ?= production

ifeq ($(ENVIRONMENT),development)
	DB_CMD := go run cmd/db/main.go
else
  DB_CMD := docker-compose --env-file .env run --rm app ./db
endif

db-init:
	@$(DB_CMD) init

db-migration-create:
	@$(DB_CMD) migration:create $(filter-out $@,$(MAKECMDGOALS))

db-migration-run:
	@$(DB_CMD) migration:run

db-migration-rollback:
	@$(DB_CMD) migration:rollback

db-seed-create:
	@$(DB_CMD) seed:create $(filter-out $@,$(MAKECMDGOALS))

db-seed-run-all:
	@$(DB_CMD) seed:runAll

db-seed-run:
	@$(DB_CMD) seed:run $(filter-out $@,$(MAKECMDGOALS))
