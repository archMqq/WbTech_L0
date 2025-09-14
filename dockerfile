FROM golang:1.24.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/demo/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config ./config

COPY --from=builder /app/internal ./internal

ENV CONFIG_PATH=./config/config.yaml

EXPOSE 8080
CMD ["./main"]