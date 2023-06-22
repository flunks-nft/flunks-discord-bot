build:
	@go build -o bin/discord-bot

run: build
	@./bin/discord-bot

test:
	@go test ./... -v