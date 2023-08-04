# Flunks Discord Bot

## How to run the bot

- Create a Discord Bot and add it to your personal server
- Create an `.env` file in the root directory according to the variables in `.env.example`
- Run `make run-discord` in your command line to run the interactive Discord bot
- Run `make run-raider` in your command line to run the raid matching and conclusion worker

## Discord OAuth server

This runs on port 8080, it handles all Discord OAuth2 logics.

```
make run-oauth-server
```
