include .envrc

test:
	@go test -v ./...
start_dev:
	docker compose up -d
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" $(GOOSE_PATH) -dir $(GOOSE_MIGRATION_DIR) up
stop_dev:
	docker compose up -d
	$(go env GOPATH)/bin/goose down