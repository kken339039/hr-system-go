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

db-init:
	@docker-compose --env-file .env run --rm app ./db init

db-migration-create:
	@docker-compose --env-file .env run --rm app ./db migration:create $(filter-out $@,$(MAKECMDGOALS))

db-migration-run:
	@docker-compose --env-file .env run --rm app ./db migration:run

db-migration-rollback:
	@docker-compose --env-file .env run --rm app ./db migration:rollback

db-seed-create:
	@docker-compose --env-file .env run --rm app ./db seed:create $(filter-out $@,$(MAKECMDGOALS))

db-seed-run-all:
	@docker-compose --env-file .env run --rm app ./db seed:runAll

db-seed-run:
	@docker-compose --env-file run --rm app ./db seed:run $(filter-out $@,$(MAKECMDGOALS))
