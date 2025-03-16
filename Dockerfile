FROM golang:1.24 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o app .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x app

ENTRYPOINT ["./app"]