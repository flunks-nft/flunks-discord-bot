# Start from the official golang alpine image
FROM golang:1.20.1-alpine AS builder

WORKDIR /app

# Install make and other dependencies
RUN apk update && apk add --no-cache make git

COPY . .

# Use the Makefile to build the discord runner
RUN make build-discord-runner

# Final stage
FROM debian:buster-slim

WORKDIR /root

# Install ca-certificates using apt-get (since we are using a Debian-based image)
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/discord-runner .

COPY --from=builder /app/templates ./templates

EXPOSE 8080