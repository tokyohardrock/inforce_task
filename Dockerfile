FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/app ./cmd/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

RUN mkdir -p internal/config

COPY internal/config/config.yaml ./internal/config/config.yaml

COPY --from=builder /app/bin/app /app/app

RUN adduser -D -g '' appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 9090

ENTRYPOINT ["/app/app"]
