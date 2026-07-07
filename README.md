# GO-Eats

A food delivery backend built as **three independent Go microservices** communicating over **gRPC**, fronted by an **API Gateway**. Each service owns its own PostgreSQL database and runs as a separate process.

---

## Architecture

```
                        ┌─────────────────────────────────────┐
  React Frontend  ──►   │     API Gateway  :8080              │
  (or any client)       │   (reverse proxy, routes by path)   │
                        └────────┬────────┬────────┬──────────┘
                                 │        │        │
                    /user/*      │  /restaurant/*  │  /cart/* /delivery/*
                                 ▼        ▼        ▼  /review/* /notify/*
                        ┌────────────┐ ┌──────────────────┐ ┌─────────────┐
                        │ User Svc   │ │ Restaurant Svc   │ │  Order Svc  │
                        │ :8081 HTTP │ │ :8082 HTTP       │ │ :8083 HTTP  │
                        │ :9091 gRPC │ │ :9092 gRPC       │ │             │
                        └────────────┘ └──────────────────┘ └──────┬──────┘
                               ▲                ▲                   │
                               │  ValidateToken │  GetMenuItem      │
                               └───────────────┴───────────────────┘
                                         gRPC calls from Order Svc

                        ┌────────────┐  ┌──────────────────┐  ┌─────────────┐
                        │goeats_user │  │goeats_restaurant │  │goeats_order │
                        │ PostgreSQL │  │   PostgreSQL      │  │ PostgreSQL  │
                        └────────────┘  └──────────────────┘  └─────────────┘
```

### Service Responsibilities

| Service | HTTP | gRPC | Database | Owns |
|---|---|---|---|---|
| **User Service** | :8081 | :9091 | `goeats_user` | Registration, login, JWT issuance |
| **Restaurant Service** | :8082 | :9092 | `goeats_restaurant` | Restaurants, menu items |
| **Order Service** | :8083 | — | `goeats_order` | Cart, orders, delivery, reviews, notifications |
| **API Gateway** | :8080 | — | — | Routes requests to the three services |

### gRPC Cross-Service Calls

Two real inter-service gRPC calls happen at runtime:

1. **Token validation** — On every authenticated request to Order Service, it calls `User Service :9091 → ValidateToken(jwt)`. Order Service never reads the JWT secret directly.

2. **Menu price lookup** — When a user places an order (`POST /cart/order/new`), Order Service calls `Restaurant Service :9092 → GetMenuItem(menu_id)` for each cart item to get current prices. It never queries the restaurant database directly.

```
User places order
  → Order Svc :8083
      → User Svc :9091  gRPC ValidateToken(jwt)        ← auth
      → Restaurant Svc :9092  gRPC GetMenuItem(id) × N  ← pricing
      → writes to goeats_order
      → publishes to NATS [order.created]
          → WebSocket push to user
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP Framework | Gin |
| Inter-service | gRPC + Protocol Buffers |
| ORM | Uptrace/bun |
| Database | PostgreSQL 16 (3 separate DBs) |
| Messaging | NATS (async order events) |
| Real-time | WebSockets (gorilla/websocket) |
| Auth | JWT (golang-jwt) |
| 2FA | TOTP (pquerna/otp) |
| Testing | testify + testcontainers |

---

## Prerequisites

- **Go 1.24+**
- **PostgreSQL 16** running on `localhost:5432`
- **NATS Server** (optional — app degrades gracefully without it)
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` (only needed to regenerate proto files — pre-generated code is already committed)

---

## Getting Started

### 1. Clone

```bash
git clone https://github.com/Ayocodes24/GO-Eats.git
cd GO-Eats
```

### 2. Create the three databases

Run once before starting services for the first time:

```bash
make setup-dbs
```

This runs `scripts/setup-dbs.sh` which creates `goeats_user`, `goeats_restaurant`, and `goeats_order` in your local Postgres. Each service migrates its own tables on startup — no manual schema setup needed.

> **Note:** The script assumes `postgres` user with password `yourpassword` on `localhost:5432`. Edit `scripts/setup-dbs.sh` if your credentials differ, then update the three `.env.*` files to match.

### 3. Configure environment

Three env files are provided — one per service:

| File | Used by |
|---|---|
| `.env.user` | User Service |
| `.env.restaurant` | Restaurant Service |
| `.env.order` | Order Service |

The defaults work out of the box if your Postgres credentials match. Change `DB_PASSWORD` and `JWT_SECRET_KEY` in all three files if needed (the JWT secret must be identical across all three).

### 4. Start NATS (optional)

```bash
# via Docker
docker run -d --name nats -p 4222:4222 nats:latest

# or via brew
brew install nats-server && nats-server &
```

If NATS is unavailable, all services start normally — real-time notifications are disabled with a warning log.

### 5. Start all services

```bash
make run-all
```

This starts all four processes (gateway + 3 services) in the background. Logs are written to `./logs/<service>.log`.

```
Gateway        → http://localhost:8080  (entry point for all clients)
User Service   → http://localhost:8081  |  gRPC :9091
Restaurant Svc → http://localhost:8082  |  gRPC :9092
Order Service  → http://localhost:8083
```

**Always use the gateway at `:8080`** — that is the single entry point, just like production.

### 6. Verify everything is up

```bash
curl http://localhost:8080/healthz
# {"status":"ok","service":"gateway"}

curl http://localhost:8081/healthz
# {"service":"user-svc"}

curl http://localhost:8082/healthz
# {"service":"restaurant-svc"}

curl http://localhost:8083/healthz
# {"service":"order-svc"}
```

---

## Running Services Individually

If you want to start services one by one (useful for debugging):

```bash
# Terminal 1 — User Service
go run ./cmd/user-svc/...

# Terminal 2 — Restaurant Service
go run ./cmd/restaurant-svc/...

# Terminal 3 — Order Service (needs User + Restaurant gRPC up first)
go run ./cmd/order-svc/...

# Terminal 4 — Gateway
go run ./cmd/gateway/...
```

> Order Service dials User Service `:9091` and Restaurant Service `:9092` at startup. Start User and Restaurant services before Order Service.

---

## Running the Legacy Monolith

The original single-binary monolith is preserved at `cmd/api/` and still compiles. It uses the original `food_delivery` database and `.env` file:

```bash
go run ./cmd/api/main.go
```

---

## Makefile Reference

```bash
make proto        # Regenerate gRPC Go code from .proto files
make build-all    # Compile all 4 binaries
make setup-dbs    # Create the 3 PostgreSQL databases
make run-all      # Start all 4 services
```

---

## API Reference

All requests go through the gateway at `http://localhost:8080`.

### Users `/user`

| Method | Path | Auth | Body / Notes |
|---|---|---|---|
| POST | `/user/` | — | `{"name":"","email":"","password":""}` |
| POST | `/user/login` | — | `{"email":"","password":""}` → `{"token":"<jwt>"}` |
| DELETE | `/user/:id` | — | Deletes user by ID |

### Restaurants `/restaurant`

| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/restaurant/` | — | `multipart/form-data`: name, description, address, city, state, image (file or URL) |
| GET | `/restaurant/` | — | List all restaurants |
| GET | `/restaurant/:id` | — | Get restaurant by ID |
| DELETE | `/restaurant/:id` | — | Delete restaurant |
| POST | `/restaurant/menu` | — | `multipart/form-data`: name, description, price, category, restaurant_id, photo (optional) |
| GET | `/restaurant/menu?restaurant_id=<id>` | — | List menu items for a restaurant |
| DELETE | `/restaurant/menu/:restaurant_id/:menu_id` | — | Delete menu item |

### Cart & Orders `/cart`

Requires `Authorization: Bearer <jwt>` on all routes.

| Method | Path | Body / Notes |
|---|---|---|
| POST | `/cart/add` | `{"item_id":1,"restaurant_id":1,"quantity":2}` |
| GET | `/cart/list` | Returns current cart items |
| DELETE | `/cart/remove/:cart_item_id` | Remove single item from cart |
| POST | `/cart/order/new` | `{"delivery_address":"123 Main St"}` — places order, calls gRPC to get prices |
| GET | `/cart/orders` | Order history |
| GET | `/cart/orders/:id` | Items in a specific order |
| GET | `/cart/orders/deliveries/:id` | Delivery info for an order |

### Delivery `/delivery`

| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/delivery/add` | — | Register a delivery person |
| POST | `/delivery/login` | — | TOTP-based login |
| POST | `/delivery/update-order` | JWT | Update order status (`on_the_way`, `delivered`, etc.) |
| GET | `/delivery/deliveries/:order_id` | JWT | List deliveries for an order |

### Reviews `/review`

Requires `Authorization: Bearer <jwt>`.

| Method | Path | Body |
|---|---|---|
| POST | `/review/` | `{"restaurant_id":1,"rating":5,"comment":"Great food"}` |
| DELETE | `/review/:id` | — |

### Real-time Notifications `/notify`

```
GET /notify/ws?token=<jwt>
```

WebSocket connection. Messages are pushed when:
- A new order is placed → `orders.new.<userID>`
- Order status changes → `orders.status.<orderID>`

Message format:
```
USER_ID:42|MESSAGE:Your order number 7 has been successfully placed...
```

### Announcements `/announcements`

| Method | Path | Notes |
|---|---|---|
| GET | `/announcements/events` | SSE stream of flash announcements |

---

## Complete Flow Example

```bash
# 1. Register a user
curl -X POST http://localhost:8080/user/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Ayush","email":"ayush@example.com","password":"secret123"}'

# 2. Login → get JWT
TOKEN=$(curl -s -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ayush@example.com","password":"secret123"}' | jq -r .token)

# 3. Add a restaurant
curl -X POST http://localhost:8080/restaurant/ \
  -F "name=Pizza Palace" \
  -F "description=[Italian] Authentic Neapolitan pizza" \
  -F "address=12 MG Road" -F "city=Bangalore" -F "state=KA" \
  -F "image_url=https://images.unsplash.com/photo-1565299624946-b28f40a0ae38"

# 4. Add a menu item
curl -X POST http://localhost:8080/restaurant/menu \
  -F "name=Margherita" -F "description=Classic tomato and mozzarella" \
  -F "price=299" -F "category=Pizza" -F "restaurant_id=1"

# 5. Add to cart (calls User Svc gRPC to validate token)
curl -X POST http://localhost:8080/cart/add \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"item_id":1,"restaurant_id":1,"quantity":2}'

# 6. Place order (calls User Svc gRPC + Restaurant Svc gRPC)
curl -X POST http://localhost:8080/cart/order/new \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"delivery_address":"42 Indiranagar, Bangalore"}'
```

---

## Project Structure

```
GO-Eats/
├── cmd/
│   ├── api/                  # Legacy monolith (preserved)
│   ├── gateway/              # API Gateway — routes to 3 services
│   ├── user-svc/             # User Service (HTTP :8081 + gRPC :9091)
│   ├── restaurant-svc/       # Restaurant Service (HTTP :8082 + gRPC :9092)
│   └── order-svc/            # Order Service (HTTP :8083)
│       └── middleware/       # gRPC-backed auth middleware
│
├── proto/
│   ├── user/user.proto       # ValidateToken RPC definition
│   └── restaurant/           # GetMenuItem RPC definition
│
├── pkg/
│   ├── proto/                # Generated gRPC Go stubs
│   │   ├── user/
│   │   └── restaurant/
│   ├── grpc/
│   │   ├── userserver/       # gRPC server impl for User Service
│   │   └── restaurantserver/ # gRPC server impl for Restaurant Service
│   ├── auth/                 # Shared JWT claims + validation
│   ├── database/             # DB wrapper, models, migrations
│   ├── handler/              # HTTP handlers (shared across services)
│   ├── service/              # Business logic (shared across services)
│   ├── nats/                 # NATS pub/sub wrapper
│   ├── wsclients/            # Thread-safe WebSocket registry (sync.RWMutex)
│   └── storage/              # Pluggable image storage
│
├── scripts/
│   ├── setup-dbs.sh          # Create 3 PostgreSQL databases
│   └── run-all.sh            # Start all 4 processes
│
├── .env.user                 # User Service env config
├── .env.restaurant           # Restaurant Service env config
├── .env.order                # Order Service env config
└── Makefile
```

---

## Concurrency Model

- **Goroutines** — Gin spawns a goroutine per HTTP request automatically. The gRPC server in each service runs on its own goroutine alongside the HTTP server.
- **sync.RWMutex** — WebSocket client registry uses `RWMutex` so multiple NATS subscriber goroutines can read concurrently while writes are exclusive.
- **sync.Mutex** — Background image update after menu insert holds a mutex for the DB write, keeping it off the request path.
- **NATS** — Order events (`orders.new`, `orders.status`) are published asynchronously; subscriber goroutines managed by the NATS library push to WebSocket clients.

---

## Running Tests

Tests use an in-memory SQLite DB (no Postgres needed) and testcontainers for NATS-dependent tests.

```bash
# All tests
go test ./...

# Specific packages
go test ./pkg/tests/user/...
go test ./pkg/tests/cart/...
go test ./pkg/tests/restaurant/...
go test ./pkg/tests/delivery/...
```

---

## Regenerating Proto Files

Only needed if you modify `.proto` definitions:

```bash
# Install tools (one-time)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
brew install protobuf   # macOS

# Regenerate
make proto
```

---

## Environment Variables Reference

### Shared across all three services

| Variable | Description |
|---|---|
| `APP_ENV` | `dev` or `TEST` |
| `DB_HOST` | Postgres host (default: `localhost`) |
| `DB_PORT` | Postgres port (default: `5432`) |
| `DB_USERNAME` | Postgres user |
| `DB_PASSWORD` | Postgres password |
| `JWT_SECRET_KEY` | HMAC key — must be identical across all services |

### Per-service

| Variable | Service | Description |
|---|---|---|
| `DB_NAME` | all | `goeats_user` / `goeats_restaurant` / `goeats_order` |
| `NATS_URL` | order-svc | NATS connection string (default: `nats://127.0.0.1:4222`) |
| `PASSWORD_SALT` | user-svc, order-svc | bcrypt password salt |
| `STORAGE_TYPE` | user-svc, restaurant-svc | `local` |
| `LOCAL_STORAGE_PATH` | user-svc, restaurant-svc | Filesystem path for uploaded images |
| `UNSPLASH_API_KEY` | restaurant-svc | Optional — for auto menu item photos |
