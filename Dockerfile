FROM golang:1.22.1-alpine3.18 AS builder

RUN go version

COPY ./ /minqyy-api-build
WORKDIR /minqyy-api-build

RUN go build -o ./.bin/minqyy-api ./cmd/api/main.go

# Lightweight container with executables
FROM alpine:latest

WORKDIR /app

COPY --from=builder /minqyy-api-build/.bin/ ./.bin
COPY --from=builder /minqyy-api-build/config/ ./config

CMD ["./.bin/minqyy-api"]
