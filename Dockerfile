# Builder
FROM golang:1.23.6-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ddns-bridge .

# Runtime
FROM debian:bookworm-slim

RUN apt install ca-certificates

WORKDIR /app
COPY --from=builder /app/ddns-bridge .

CMD ["./ddns-bridge"]