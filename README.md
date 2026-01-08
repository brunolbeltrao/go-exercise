---

## Build and Run (Local)

Prerequisites: Go 1.22+

```bash
# build
go build -o bin/app ./cmd/server

# run (defaults: PORT=8080, CACHE_TTL=60s, HTTP_TIMEOUT=3s)
PORT=8080 CACHE_TTL=60s HTTP_TIMEOUT=3s ./bin/app
```

## API Usage

```bash
# All supported pairs
curl -s "http://localhost:8080/api/v1/ltp"

# Specific pairs (repeated param)
curl -s "http://localhost:8080/api/v1/ltp?pair=BTC/USD&pair=BTC/EUR"

# Specific pairs (CSV)
curl -s "http://localhost:8080/api/v1/ltp?pairs=BTC/USD,BTC/CHF"
```

Response:

```json
{
  "ltp": [
    {"pair":"BTC/CHF","amount":49000.12},
    {"pair":"BTC/EUR","amount":50000.12},
    {"pair":"BTC/USD","amount":52000.12}
  ]
}
```

## Environment Variables

- `PORT`: server port (default `8080`)
- `CACHE_TTL`: per-pair cache TTL, e.g. `60s` (default `60s`)
- `HTTP_TIMEOUT`: upstream Kraken timeout, e.g. `3s` (default `3s`)

Optional .env for development:

```bash
# Copy the example and edit values as needed
cp env.example .env

# Load into current shell (bash/zsh)
set -a; source .env; set +a
```

## Tests

```bash
# run all tests
go test ./...

# run only integration tests
go test -v ./test/integration

# optional live test against Kraken
TEST_LIVE_KRAKEN=1 go test -v ./test/integration -run Live
```

## Docker

```bash
# build image
docker build -t go-exercise:latest .

# run container (inline env)
docker run --rm -p 8080:8080 -e PORT=8080 -e CACHE_TTL=60s -e HTTP_TIMEOUT=3s go-exercise:latest

# or using env file in deploy/env
cp deploy/env/dev.env.example deploy/env/dev.env
docker run --rm -p 8080:8080 --env-file deploy/env/dev.env go-exercise:latest
```

## Health Endpoints

- Liveness: `GET /healthz` → 200 OK
- Readiness: `GET /readyz` → 200 OK
