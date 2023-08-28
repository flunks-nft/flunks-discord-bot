# Start from the official golang image
FROM golang:1.20.1 AS builder

WORKDIR /app

# Install make
RUN apt-get update && apt-get install -y make

COPY . .

# Use the Makefile to build the discord raid runner
RUN make build-oauth-server

# Final stage
FROM debian:buster-slim

WORKDIR /root

# Install ca-certificates using apt-get (since we are using a Debian-based image)
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/oauth-server .

CMD ["./oauth-server"]

# Expose ports (for server)
EXPOSE 8080