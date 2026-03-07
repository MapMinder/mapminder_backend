# MapMinder Backend

MapMinder is a location-based reminder application designed to help you remember tasks tied to specific places.
The name comes from a combination of "Map" and "Reminder" — exactly what the app offers.

Have you ever forgotten to buy groceries on your way home from work because you were too caught up in a podcast,
anime, or just your thoughts? You only realize it when you're already home — and it keeps happening.
MapMinder solves this by letting you set reminders that trigger when you arrive near a specific location.
Whether it's picking up groceries, collecting laundry, or stopping by a friend's place, MapMinder ensures you
get timely, location-aware reminders — right when and where you need them.

---

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Setup & Installation](#setup--installation)
- [Database Migrations](#database-migrations)
- [Running the Server](#running-the-server)
- [Testing](#testing)

---

## Architecture Overview

This backend follows a layered architecture to keep concerns separated and the codebase easy to maintain:

```
cmd/
  main.go                   # Entry point
internal/
  config/                   # App configuration (env loading)
  logger/                   # Zap logger setup
  status/                   # Shared status/error types
  httpServer/               # Router setup
  server/                   # Server bootstrap
  infrastructure/
    database/               # DB connection and transaction manager
feature/
  auth/                     # Google Sign-In (OAuth2 + JWT issuance)
    domain/                 # Domain models
    dto/                    # Request/response DTOs
    handler/                # HTTP handlers
    repository/             # DB + Google API access
    usecase/                # Business logic
  reminder/                 # Location-based reminder CRUD + event tracking
    domain/
    dto/
    handler/
    mapper/                 # DTO <-> domain mapping
    repository/
    usecase/
  health/                   # Health check endpoint
  user/                     # User domain and repository
shared/
  appError/                 # Centralized app error types
  middleware/               # JWT auth + error handler middleware
  response/                 # Standardized HTTP response helpers
  tx/                       # Transaction manager abstraction
  validator/                # Request validation (global + reminder-specific)
  uuid_manager/             # UUID generation abstraction
  time/                     # Time provider abstraction (mockable for tests)
```

Each feature is self-contained with its own `repository`, `usecase`, and `handler` layers. Dependencies flow inward — handlers depend on usecases, usecases depend on repositories, repositories depend on the database.

Shared utilities like transaction management, UUID generation, and time are abstracted behind interfaces so they can be easily mocked in tests.

Authentication is handled via **Google Sign-In** (OAuth2) and **JWT** tokens. Every protected route passes through the JWT middleware which injects the user ID into the request context.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| Web Framework | Gin |
| ORM | GORM |
| Database | MySQL |
| Logger | Zap |
| Auth | Google OAuth2 + JWT |
| Containerization | Docker / Docker Compose |
| Migrations | golang-migrate |

---

## Prerequisites

Make sure you have the following installed:

- [Go 1.24+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
- MySQL 8+ (if running without Docker)

---

## Environment Variables

Create a `.env` file in the project root. Use the following keys:

```env
PORT=
ENVIRONMENT=

LOG_LEVEL=

# MySQL
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASSWORD=

# MySQL Root (only for init / debugging)
MYSQL_ROOT_PASSWORD=

# Google OAuth2
GOOGLE_CLIENT_ID=

# JWT
JWT_SIGN_KEY=
```

---

## Setup & Installation

**1. Clone the repository**
```bash
git clone https://github.com/MapMinder/mapminder_backend.git
cd mapminder_backend
```

**2. Copy and fill in environment variables**
```bash
cp .env.example .env
# edit .env with your values
```

**3. Install Go dependencies**
```bash
go mod download
```

---

## Database Migrations

Migrations are managed in a separate repository: [MapMinder/migration](https://github.com/MapMinder/migration)

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) to manage database schema changes.

**1. Clone the migration repository**
```bash
git clone https://github.com/MapMinder/migration.git
cd migration
```

**2. Run migrations**
```bash
migrate -path . -database "mysql://DB_USER:DB_PASSWORD@tcp(DB_HOST:DB_PORT)/DB_NAME" up
```

**3. Rollback last migration**
```bash
migrate -path . -database "mysql://DB_USER:DB_PASSWORD@tcp(DB_HOST:DB_PORT)/DB_NAME" down 1
```

> Make sure your database is running before applying migrations.

---

## Running the Server

**Option 1: Run with Go directly**
```bash
go run cmd/main.go
```

**Option 2: Run with Docker Compose**
```bash
docker compose up
```

Docker Compose will spin up both the application and the MySQL database together.

> The server runs on the port defined by `PORT` in your `.env` (default: `8080`).

---

## Testing

Tests are written using Go's standard `testing` package along with `gomock` for mocking and `sqlmock` for database mocking.

**Run all tests**
```bash
go test ./...
```

**Run tests for a specific feature**
```bash
go test ./feature/reminder/...
```

**Run a specific test**
```bash
go test ./feature/reminder/handler -run TestReminderHandler_DeleteReminder
```

Each feature has tests at all three layers — repository, usecase, and handler — to ensure full coverage from the database up to the HTTP response.
