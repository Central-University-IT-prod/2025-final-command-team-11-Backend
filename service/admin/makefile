
run:
	go build -o ./bin/admin ./cmd/admin
	./bin/admin

gen-swagger-docs:
	swag init -g cmd/admin/main.go

e2e-test:
	go test ./tests -count=1

unit-test:
	go test ./internal/... -count=1
