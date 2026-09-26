# GoShield

GoShield is a lightweight reverse proxy and load-balancing gateway built
from scratch in Go.

The project sits between clients and backend services and provides a
single entry point for forwarding HTTP and WebSocket traffic. It
includes JWT-based authentication, per-user rate limiting, round-robin
backend selection, backend health checking, Redis-backed server state,
PostgreSQL persistence, TLS support, structured logging, and graceful
shutdown.

## Architecture

``` text
                         ┌──────────────────┐
                         │      Client      │
                         └────────┬─────────┘
                                  │
                           HTTP / HTTPS / WS
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │         GoShield         │
                    │                          │
                    │  JWT Authentication      │
                    │  Rate Limiting           │
                    │  HTTP Proxy              │
                    │  WebSocket Proxy         │
                    │  Round-Robin Routing     │
                    │  TLS Support             │
                    └────────────┬─────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
             ┌──────────────┐          ┌──────────────┐
             │   Backend 1  │          │   Backend 2  │
             └──────────────┘          └──────────────┘
                    ▲                         ▲
                    └────────────┬────────────┘
                                 │
                          Health Checking
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
              PostgreSQL                   Redis
             Server metadata        Active/all server state
```

## Key Features

-   Reverse proxy for HTTP requests
-   WebSocket proxy with bidirectional message forwarding
-   Round-robin load balancing across active backend servers
-   Periodic backend health checks using a `/health` endpoint
-   Redis cache/state layer for active and registered backend
    information
-   PostgreSQL persistence for backend server and admin information
-   JWT authentication
-   Per-user in-memory rate limiting
-   TLS listener support
-   Request logging with rotating log files
-   Panic recovery middleware
-   Graceful application shutdown
-   Admin APIs for registering, listing, inspecting, and removing
    backend servers

## Technology Stack

  Component                       Technology
  ------------------------------- -------------------------------------
  Language                        Go 1.24.4
  Web Framework                   Fiber v2
  WebSocket                       Gorilla WebSocket + Fiber WebSocket
  Database                        PostgreSQL
  Cache / State                   Redis
  Authentication                  JWT
  Password Security               bcrypt
  Logging                         Uber Zap + rotating log files
  Database Access                 sqlx
  TLS                             Go TLS support through Fiber
  Build / Dependency Management   Go Modules

## Project Structure

``` text
Server/
├── cmd/
│   └── server/
│       └── main.go
│
├── configs/
│   └── env_config.conf
│
├── certs/
│   ├── cert.pem
│   └── key.pem
│
├── internal/
│   ├── config/              # Application configuration
│   ├── conn/                # PostgreSQL connection
│   ├── db/                  # Database repositories
│   ├── gbl/                 # Shared constants
│   ├── handler/
│   │   ├── proxyhandler/    # HTTP/WebSocket proxy logic
│   │   └── serverhandler/   # Admin/server management handlers
│   ├── log/                 # Application/access logging
│   ├── middleware/
│   │   ├── jwt_auth.go      # JWT authentication
│   │   └── rate_limiter.go  # Per-user rate limiting
│   ├── model/               # Data models
│   ├── redis/               # Redis initialization and server state
│   ├── routes/              # Admin and proxy routes
│   └── service/
│       ├── healthcheck/     # Backend health checking
│       └── routing/         # Round-robin load balancing
│
└── loaddb/
    └── shield_db_schema.sql
```

## Request Flow

### HTTP

A normal HTTP request follows this general path:

``` text
Client
  │
  ▼
GoShield
  │
  ├── JWT validation
  │
  ├── Rate-limit validation
  │
  ├── Select active backend
  │
  └── Forward request
          │
          ▼
      Backend Server
          │
          ▼
      GoShield
          │
          ▼
        Client
```

GoShield also adds standard proxy metadata such as:

-   `X-Forwarded-For`
-   `X-Forwarded-Proto`
-   `X-Forwarded-Host`

### WebSocket

WebSocket requests are detected separately and handled using a
persistent bidirectional connection:

``` text
Client WebSocket
       │
       ▼
   GoShield
    │     │
    │     └──────── Backend → Client
    │
    └────────────── Client → Backend
```

Two forwarding loops are used so messages can flow independently in both
directions.

## Load Balancing

GoShield uses a round-robin strategy.

For example, with three active backend servers:

``` text
Request 1 → Backend 1
Request 2 → Backend 2
Request 3 → Backend 3
Request 4 → Backend 1
Request 5 → Backend 2
...
```

The current routing position is protected with a mutex because multiple
requests can access the load-balancer concurrently.

Only servers present in the active-server list are eligible for routing.

## Backend Health Checking

A background health-check service periodically checks registered backend
servers.

Each backend is expected to expose:

``` text
GET /health
```

A successful HTTP 200 response marks the backend as active; otherwise it
is marked inactive.

The current server information is maintained in PostgreSQL and
synchronized into Redis for routing.

## Redis Usage

Redis maintains two server-state entries:

``` text
all_servers
active_servers
```

The application loads server information from PostgreSQL into Redis and
uses the active-server list during request routing.

This avoids querying PostgreSQL for every proxied request.

## PostgreSQL

The included SQL schema creates two main tables:

``` text
serverinfo
admin
```

`serverinfo` stores backend host, port, protocol/method, and status.

`admin` stores administrator information used by the admin
authentication flow.

Initialize the database using:

``` bash
psql -U postgres -d load_balancer -f loaddb/shield_db_schema.sql
```

Adjust the database name, credentials, and connection settings to match
your environment.

## Admin API

The application exposes administrative endpoints under `/admin`.

  Method   Endpoint                   Purpose
  -------- -------------------------- -------------------------------
  POST     `/admin/login`             Authenticate an administrator
  GET      `/admin/server/list`       List registered servers
  GET      `/admin/server/details`    Retrieve server details
  POST     `/admin/register/server`   Register a backend server
  POST     `/admin/delete/server`     Remove a backend server

The authenticated token is then used for protected application requests.

## Configuration

Application configuration is stored in:

``` text
configs/env_config.conf
```

It contains configuration for:

-   Proxy host and port
-   TLS port
-   PostgreSQL
-   Redis
-   TLS certificate and private-key paths

Example structure:

``` json
{
  "proxy": {
    "host": "0.0.0.0",
    "port": 8080,
    "tls_port": 8081,
    "protocol": "http"
  },
  "dbconfig": {
    "host": "localhost",
    "port": 45432,
    "database": "load_balancer",
    "username": "postgres",
    "password": "postgres"
  },
  "redis": {
    "host": "localhost",
    "port": 6379
  }
}
```

Do not use development credentials or committed private keys in a
production deployment.

## Running Locally

### Prerequisites

Install:

-   Go 1.24.4 or compatible Go version
-   PostgreSQL
-   Redis

### 1. Clone the repository

``` bash
git clone <repository-url>
cd GoShield/Server
```

### 2. Configure PostgreSQL

Create the database and run:

``` bash
psql -U postgres -d load_balancer -f loaddb/shield_db_schema.sql
```

Update `configs/env_config.conf` with your PostgreSQL settings.

### 3. Start Redis

Make sure Redis is running on the configured host and port.

Default configuration:

``` text
localhost:6379
```

### 4. Configure TLS

Provide a certificate and private key and update:

``` text
configs/env_config.conf
```

with the correct certificate paths.

For local development, you can use a locally generated self-signed
certificate.

### 5. Download dependencies

``` bash
go mod download
```

### 6. Start GoShield

From the `Server` directory:

``` bash
go run ./cmd/server
```

The application starts:

-   HTTP proxy listener on the configured proxy port
-   HTTPS/TLS listener on the configured TLS port
-   Background backend health checking

## Security Notes

This repository is intended as a learning/engineering project and should
be hardened before production use.

In particular:

-   Store JWT secrets outside source code.
-   Never commit real passwords, private keys, or production
    certificates.
-   Use environment variables or a secrets manager for sensitive
    configuration.
-   Use production-issued TLS certificates.
-   Apply authentication and authorization rules appropriate to the
    deployment.
-   Consider distributed rate limiting when running multiple GoShield
    instances.

## Engineering Considerations / Future Improvements

Potential improvements include:

-   Distributed rate limiting using Redis
-   Streaming large HTTP responses instead of buffering the complete
    response
-   Stronger configuration and secret management
-   More comprehensive request-header forwarding
-   Improved health-check synchronization and cancellation
-   Backend connection pooling and HTTP transport tuning
-   Retry and circuit-breaker policies
-   Prometheus metrics and distributed tracing
-   Automated unit, integration, and concurrency tests
-   Containerization and CI/CD deployment
-   Multiple GoShield instances with shared routing state

## What This Project Demonstrates

GoShield was built to explore practical backend and networking concepts
in Go, including:

-   HTTP request proxying
-   WebSocket bidirectional communication
-   Reverse-proxy architecture
-   Load balancing
-   Health checking
-   Concurrent request handling
-   Mutexes and goroutines
-   JWT authentication
-   Rate limiting
-   Redis-backed application state
-   PostgreSQL persistence
-   TLS
-   Logging
-   Graceful shutdown

## Author

Built from scratch as a Go backend and networking project.
