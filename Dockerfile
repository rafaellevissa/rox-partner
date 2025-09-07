FROM golang:1.24.5-alpine AS builder-base
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM builder-base AS ingestor-builder
RUN go build -o b3-ingestor ./cmd/b3ingestor

FROM builder-base AS api-builder
RUN go build -o b3-api ./cmd/api

FROM alpine:3.19 AS ingestor
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=ingestor-builder /app/b3-ingestor /usr/local/bin/b3-ingestor
COPY configs/config.yaml /app/config.yaml
ENTRYPOINT ["b3-ingestor"]

FROM alpine:3.19 AS api
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=api-builder /app/b3-api /usr/local/bin/b3-api
COPY configs/config.yaml /app/config.yaml
ENTRYPOINT ["b3-api"]
