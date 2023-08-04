db-up:
	@docker-compose up db -d --force-recreate --no-deps

db-down:
	@docker-compose down db

buid-oauth-server:
	@go build -o bin/oauth-server ./cmd/oauth-server

build-discord-runner:
	@go build -o bin/discord-runner ./cmd/discord-runner

build-raid-runner:
	@go build -o bin/raid-runner ./cmd/raid-runner

run-discord: build-discord-runner
	@./bin/discord-runner

run-raider: build-raid-runner
	@./bin/raid-runner

run-oauth-server: buid-oauth-server
	@./bin/oauth-server

test:
	go test ./... -v