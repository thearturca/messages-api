FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dist/message-http-service ./cmd/message-http-service/
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dist/message-consumer-service ./cmd/message-consumer-service/
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dist/message-processor ./cmd/message-processor/

FROM alpine:latest AS consumer-production
WORKDIR /app
COPY --from=builder /app/dist/message-consumer-service .
CMD ["./message-consumer-service"]

FROM alpine:latest AS api-production
WORKDIR /app
COPY --from=builder /app/dist/message-http-service .
EXPOSE ${PORT}
CMD ["./message-http-service"]

FROM alpine:latest AS processor-production
WORKDIR /app
COPY --from=builder /app/dist/message-processor .
CMD ["./message-processor"]

FROM gomicro/goose AS migrator
WORKDIR /migrations
COPY ./db/migrations/*.sql ./
COPY entrypoint_migrator.sh ./entrypoint.sh

CMD ["sh", "./entrypoint.sh"]
