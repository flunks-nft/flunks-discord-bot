db-up:
	@docker-compose up db -d --force-recreate --no-deps

db-down:
	@docker-compose down db

build:
	@go build -o bin/discord-bot

run: build
	@./bin/discord-bot

test:
	@go test ./... -v