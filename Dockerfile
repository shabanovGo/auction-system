FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auction ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/auction .
COPY .env .env

RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser

CMD ["./auction"]
