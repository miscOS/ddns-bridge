# Builder
FROM golang:1.26.1-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ddns-bridge .

# Runtime
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

WORKDIR /app
COPY --from=builder /app/ddns-bridge .

CMD ["./ddns-bridge"]