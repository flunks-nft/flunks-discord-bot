# Start from the official golang image
FROM golang:1.20.1 AS builder

WORKDIR /app

# Install make
RUN apt-get update && apt-get install -y make

COPY . .

# Use the Makefile to build the oauth server
RUN make buid-oauth-server

# Final stage
FROM debian:buster-slim

WORKDIR /root

COPY --from=builder /app/bin/oauth-server .

CMD ["./oauth-server"]

# Expose ports (for server)
EXPOSE 8080