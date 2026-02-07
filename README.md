# AppDrop

A Go API for managing pages and widgets, built with [Gin](https://gin-gonic.com/) and [GORM](https://gorm.io/) (PostgreSQL).

## Prerequisites

- Go 1.24+
- Docker & Docker Compose (for PostgreSQL)

## Setup

1. **Start PostgreSQL**

   ```bash
   docker-compose up -d postgres
   ```

2. **Create a `.env` file** in the project root:

   ```
   DATABASE_URL=postgres://appdrop:postgres@localhost:5432/appdropdb?sslmode=disable
   ```

   (Adjust if you use different postgres user/password/db from `docker-compose.yml`.)

3. **Apply the schema** (if not using GORM auto-migrate). With `psql` or any PostgreSQL client, run:

   ```bash
   psql "$DATABASE_URL" -f db/schema.sql
   ```

4. **Install dependencies**

   ```bash
   go mod download
   ```

## Run the server

```bash
go run main.go
```

Server runs at **http://localhost:8082**.

## API overview

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/pages` | Create a page |
| GET | `/pages` | List all pages |
| GET | `/pages/:id` | Get page by ID |
| DELETE | `/pages/:id` | Delete a page (cannot delete home page) |
| POST | `/pages/:id/widgets` | Add widget to a page |
| PUT | `/widgets/:id` | Update a widget |
| DELETE | `/widgets/:id` | Delete a widget |
| POST | `/pages/:id/widgets/reorder` | Reorder widgets |

**Widget types:** `banner`, `product_grid`, `text`, `image`, `spacer`

## Tests

Run all tests:

```bash
go test -v ./...
```

- With `DATABASE_URL` set (e.g. from `.env`), full API tests run against the database.
- Without a database, tests that need the DB are skipped; health and validation tests still run.
