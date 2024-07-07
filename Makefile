.PHONY: start build

start:
	@go run cmd/api.go

build:
	@go build -o build/api cmd/api.go