# GO-Eats

A scalable backend for a food delivery platform built with Go, PostgreSQL, NATS, and WebSockets.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP Framework | Gin |
| ORM | Uptrace/bun |
| Database | PostgreSQL 16 |
| Messaging | NATS |
| Real-time | WebSockets (gorilla/websocket) |
| Auth | JWT (golang-jwt) |
| 2FA | TOTP (pquerna/otp) |
| Image Search | Unsplash API |
| Testing | testify + testcontainers |

---

## Architecture

```
cmd/api/          — entry point (main.go), middleware
pkg/
  abstract/       — domain interfaces (Cart, Order, Restaurant, User, ...)
  database/       — DB wrapper, models, migrations
  handler/        — HTTP handlers + route registration
  service/        — business logic
  nats/           — NATS pub/sub wrapper (nil-safe, graceful degradation)
  wsclients/      — thread-safe WebSocket connection registry
  storage/        — local / pluggable image storage
  tests/          — integration tests (testcontainers)
```

### Request flow

```
Client → Gin Router → Handler → Service → Database
                                        ↘ NATS (async notifications)
                                             ↘ WebSocket (real-time push)
```

---

## Features

- **Users** — register, login (JWT), delete
- **Restaurants** — CRUD with image upload; menu items auto-illustrated via Unsplash
- **Cart** — add/remove items; cart auto-created per user
- **Orders** — place order from cart, order history, per-order item breakdown
- **Delivery** — delivery person management, 2FA login (TOTP), order status updates
- **Reviews** — add/delete restaurant reviews
- **Notifications** — real-time order alerts via NATS → WebSocket
- **Announcements** — SSE-based rotating flash events

---

## Getting Started

### Prerequisites

- Go 1.24+
- Docker (for PostgreSQL + NATS via docker-compose)

### 1. Clone

```bash
git clone https://github.com/Ayocodes24/GO-Eats.git
cd GO-Eats
```

### 2. Start dependencies

```bash
docker-compose up -d
```

### 3. Configure environment

Create a `.env` file in the project root:

```env
# App
APP_ENV=development
SERVER_PORT=:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
DB_NAME=food_delivery

# NATS (optional — app starts without it, notifications disabled)
NATS_URL=nats://127.0.0.1:4222

# Auth
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRY_HOURS=2

# Storage
STORAGE_TYPE=local
STORAGE_DIRECTORY=uploads
LOCAL_STORAGE_PATH=./tmp/uploads

# Unsplash (optional — for auto menu item photos)
UNSPLASH_API_KEY=your-unsplash-client-id
```

### 4. Run

```bash
go run ./cmd/api/main.go
```

### 5. Health check

```
GET /healthz  →  { "status": "ok" }
```

---

## API Reference

### Users `/user`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/user/` | — | Register |
| POST | `/user/login` | — | Login → JWT |
| DELETE | `/user/:id` | — | Delete user |

### Restaurants `/restaurant`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/restaurant/` | — | Add restaurant (multipart) |
| GET | `/restaurant/` | — | List all restaurants |
| GET | `/restaurant/:id` | — | Get restaurant by ID |
| DELETE | `/restaurant/:id` | — | Delete restaurant |
| POST | `/restaurant/menu` | — | Add menu item |
| GET | `/restaurant/menu` | — | List all menu items |
| DELETE | `/restaurant/menu/:id` | — | Delete menu item |

### Cart & Orders `/cart`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/cart/add` | JWT | Add item to cart |
| GET | `/cart/list` | JWT | View cart items |
| DELETE | `/cart/:id` | JWT | Remove item from cart |
| POST | `/cart/order/new` | JWT | Place order `{"delivery_address":"..."}` |
| GET | `/cart/order/list` | JWT | Order history |
| GET | `/cart/order/:id/items` | JWT | Items in an order |
| GET | `/cart/order/:id/delivery` | JWT | Delivery info for an order |

### Delivery `/delivery`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/delivery/add` | — | Register delivery person |
| GET | `/delivery/list` | JWT | List delivery persons |
| POST | `/delivery/order` | JWT | Accept/update order status |
| POST | `/delivery/login` | — | Delivery person login (TOTP) |

### Reviews `/review`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/review/` | JWT | Add review |
| DELETE | `/review/:id` | JWT | Delete review |

### Notifications `/notify`

| Method | Path | Description |
|---|---|---|
| GET | `/notify/ws?token=<jwt>` | WebSocket — real-time order updates |

### Announcements `/announcements`

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/announcements/flash` | JWT | SSE stream of flash events |

---

## Real-time Notifications

Connect to `/notify/ws?token=<jwt>` via WebSocket. Messages are pushed when:

- A new order is placed (`orders.new.<userID>`)
- Order status changes — `on_the_way`, `delivered`, `failed`, `cancelled` (`orders.status.<orderID>`)

NATS is optional: if unavailable at startup, the app runs normally and notifications are skipped with a warning log.

---

## Running Tests

Tests use an in-memory SQLite DB and testcontainers (Docker required for NATS-dependent tests).

```bash
# All tests
go test ./...

# Specific package
go test ./pkg/tests/cart/...
go test ./pkg/tests/user/...
go test ./pkg/tests/restaurant/...
go test ./pkg/tests/delivery/...
```

### Test without NATS

`TestCartOrderWithoutNATS` demonstrates the full order flow with `nil` NATS — no container needed, no Docker required.

---

## Configuration Reference

| Env Variable | Default | Description |
|---|---|---|
| `APP_ENV` | — | `development` / `TEST` |
| `SERVER_PORT` | `:8080` | HTTP listen address |
| `DB_HOST` | — | Postgres host |
| `DB_PORT` | — | Postgres port |
| `DB_USERNAME` | — | Postgres user |
| `DB_PASSWORD` | — | Postgres password |
| `DB_NAME` | — | Postgres database name |
| `NATS_URL` | `nats://127.0.0.1:4222` | NATS server URL |
| `JWT_SECRET_KEY` | — | HMAC signing key |
| `JWT_EXPIRY_HOURS` | `2` | JWT token lifetime in hours |
| `STORAGE_TYPE` | — | `local` |
| `STORAGE_DIRECTORY` | — | URL path prefix for static files |
| `LOCAL_STORAGE_PATH` | — | Filesystem path for uploads |
| `UNSPLASH_API_KEY` | — | Unsplash Client-ID (optional) |
