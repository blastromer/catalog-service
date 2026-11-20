# Catalog Service

Go microservices for product catalog reads: **products**, **categories**, and **attributes**.

Built as production-minded read APIs over MongoDB — layered packages, typed errors, pagination metadata, health checks, graceful shutdown, Docker, and CI.

**Author:** Romer Necesario

---

## Why this exists

Ecommerce catalog data is split across collections that change independently. This repo keeps three small services with a shared platform layer so each API can scale and deploy on its own without duplicating HTTP, config, or Mongo wiring.

## Architecture

```
cmd/
  product/     # :2001
  category/    # :2002
  attribute/   # :2003
internal/
  config/      # env-based configuration
  database/    # Mongo connect, ping, close
  platform/    # pagination, API errors, response envelope, HTTP server
  product|category|attribute/
               # handler → service → repository
```

**Design choices that matter:**

| Concern | Approach |
| --- | --- |
| Boundaries | Handlers parse HTTP; services own validation/rules; repositories own Mongo |
| Safety | Category “smart query” only accepts a **whitelist** of fields (`categoryId`, `parent`, `sort`, `name`, `domain`, `meta.slug`) |
| Contracts | Consistent `{ data, pagination? }` / `{ error: { code, message } }` envelopes |
| Ops | `/healthz`, `/readyz` (DB ping), request logging, SIGINT/SIGTERM graceful shutdown |
| Testability | Repository interfaces + table-style unit tests without a live database |

## API

### Product service (`:2001`)

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/api/products?page=&limit=` | Paginated list + total count |
| `GET` | `/api/products/:productId` | By business `productId` |
| `GET` | `/healthz` | Liveness |
| `GET` | `/readyz` | Readiness (Mongo ping) |

### Category service (`:2002`)

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/api/categories?page=&limit=` | Paginated list |
| `GET` | `/api/categories/query?field=&value=` | Whitelisted filter |
| `GET` | `/api/categories/:categoryId` | By `categoryId` |

### Attribute service (`:2003`)

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/api/attributes?page=&limit=` | Paginated list |
| `GET` | `/api/attributes/:attributeId` | By attribute `_id` |

### Example response

```json
{
  "data": [{ "productId": 1001, "title": "Brake Pad Set" }],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1284,
    "totalPages": 26
  }
}
```

## Quick start

```bash
cp .env.example .env
# set MONGO_URI / MONGO_DATABASE

go mod tidy
go test ./...

go run ./cmd/product
go run ./cmd/category
go run ./cmd/attribute
```

Or with Docker:

```bash
docker compose up --build
```

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `MONGO_URI` | _(required)_ | Mongo connection string |
| `MONGO_DATABASE` | `catalog` | Database name |
| `MONGO_COLLECTION` | per-service default | Override collection |
| `HTTP_ADDR` | `:2001` / `:2002` / `:2003` | Listen address |
| `DEFAULT_PAGE_LIMIT` | `50` | Default page size |
| `MAX_PAGE_LIMIT` | `200` | Hard cap on `limit` |
| `MONGO_QUERY_TIMEOUT` | `10s` | Per-query deadline |

## Development

```bash
make test
make build
make docker-up
```

Bitbucket Pipelines runs `go test ./...` and builds all three binaries on every push.

## Engineering notes

- Route order for categories registers `/query` **before** `/:categoryId` so the static path is reachable.
- Empty lists return `data: []` with pagination, not a fake “no products” success message.
- Mongo client uses `Connect` + `Ping` on startup and disconnects cleanly on shutdown.
- `.env` is gitignored; use `.env.example` as the contract for local setup.
