FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tpu-practice-searcher ./cmd/tpu-practice-searcher

FROM alpine:latest

WORKDIR /app/
COPY --from=builder /app/config ./config
COPY --from=builder /app/tpu-practice-searcher .

EXPOSE 8000
CMD ["./tpu-practice-searcher"]