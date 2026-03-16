FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY . .
RUN go build -o wallet-app .
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/wallet-app .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migration ./migration
EXPOSE 9091
CMD ["./wallet-app"]