include .env

CDN="$(DB_PROVIDER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"

build:
	@go build -o $(BIN_PATH) $(MAIN_PATH)

run: build
	@./$(BIN_PATH)

test:
	@go test -v ./...

swag:
	swag init -d ./cmd,./services/songs

db-up:
	goose $(DB_PROVIDER) $(CDN) -dir $(MIGRATIONS_PATH) up

db-down:
	goose $(DB_PROVIDER) $(CDN) -dir $(MIGRATIONS_PATH) down

db-status:
	goose $(DB_PROVIDER) $(CDN) -dir $(MIGRATIONS_PATH) status

clean:
	@rm -rf $(BIN_PATH)
