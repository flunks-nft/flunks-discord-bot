build:
	@go build -o bin/discord-bot

run: build
	@./bin/discord-bot -t MTEyMTU2MDAzMzYwMDIwODkzNg.Gw4b-P.QDEvulL6XNY4vRJOBI3HdvBvA8IiTtAtxt1lOM

test:
	@go test ./... -v