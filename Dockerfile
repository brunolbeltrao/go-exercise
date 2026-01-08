# Build stage
FROM golang:1.22 AS builder
WORKDIR /src
COPY . .
# Static build for small runtime image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=builder /out/app /app/app
EXPOSE 8080
# Healthcheck hitting /healthz
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["/app/app"]

