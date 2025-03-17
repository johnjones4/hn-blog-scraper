FROM golang:1.24 AS builder

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
# ENV CC=x86_64-linux-gnu-gcc

RUN apt-get update && apt-get install -y \
    build-essential
    # gcc-x86-64-linux-gnu \
    # libc6-dev-amd64-cross \
RUN rm -rf /var/lib/apt/lists/*
    
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o app .

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y ca-certificates libc6
RUN rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x app

ENTRYPOINT ["./app"]