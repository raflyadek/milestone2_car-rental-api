# Stage 1: Build the Go application binary
FROM golang:1.24.0 AS builder
WORKDIR /cmd
COPY . .
WORKDIR /app/cmd
RUN go build -o car-rental-api main.go

# Stage 2: Create a minimal runtime image and running the application
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY .env .
COPY --from=builder /app/cmd/car-rental-api /car-rental-api
CMD ["/car-rental-api"]