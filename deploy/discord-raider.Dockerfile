# Start from the official golang image
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install make
RUN apt-get update && apt-get install -y make

COPY . .

# Use the Makefile to build the discord raid runner
RUN make build-raid-runner

# Final stage
FROM debian:buster-slim

WORKDIR /root

# Install ca-certificates using apt-get (since we are using a Debian-based image)
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/raid-runner ./raid-runner

CMD ["./raid-runner"]