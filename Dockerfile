FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tpu-practice-searcher ./cmd/tpu-practice-searcher

FROM alpine:latest

WORKDIR /app/
COPY --from=builder /app/config ./config
COPY --from=builder /app/tpu-practice-searcher .

ENV CONFIG_PATH=./config/config.yaml
CMD ["./tpu-practice-searcher"]


