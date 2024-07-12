.PHONY: start build db-init db-migration-create db-seed-create db-migration-run db-seed-run

start:
	@go run cmd/api.go

build:
	@go build -o build/api cmd/api.go

db-init:
	@go run cmd/db/main.go init

db-migration-create:
	@go run cmd/db/main.go migration:create $(filter-out $@,$(MAKECMDGOALS))

db-migration-run:
	@go run cmd/db/main.go migration:run

db-migration-rollback:
	@go run cmd/db/main.go migration:rollback

db-seed-create:
	@go run cmd/db/main.go seed:create $(filter-out $@,$(MAKECMDGOALS))

db-seed-run-all:
	@go run cmd/db/main.go seed:runAll

db-seed-run:
	@go run cmd/db/main.go seed:run $(filter-out $@,$(MAKECMDGOALS))
