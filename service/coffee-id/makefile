include .env
LOCAL_BIN:=$(CURDIR)/bin

run:
	go build -o ./bin/coffee-id ./cmd/coffee-id
	./bin/coffee-id

gen-swagger-docs:
	swag init -g cmd/coffee-id/main.go

e2e-test:
	go test ./tests -count=1

unit-test:
	go test ./internal/... -count=1

