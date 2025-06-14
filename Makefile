include .envrc

test:
	@go test -v ./...
start_dev:
	docker compose up -d
stop_dev:
	docker compose down
seed:
	@DB_CONNECTION_STRING=$(DB_CONNECTION_STRING) go run ./migrations/seed/main.go
