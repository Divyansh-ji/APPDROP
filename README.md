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
   JWT_SECRET=your-32-byte-secret-change-in-production
   LOGIN_PASSWORD=admin123
   JWT_COOKIE_NAME=appdrop_session
   ```

   - `DATABASE_URL`: adjust if you use different postgres user/password/db from `docker-compose.yml`.
   - `JWT_SECRET`: required for signing JWTs; use a long random string in production.
   - `LOGIN_PASSWORD`: password accepted by `POST /login` (no user DB; any email with this password works).
   - `JWT_COOKIE_NAME`: name of the HTTP-only session cookie (optional; default used if unset).

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

Server runs at **http://localhost:8090**.

## Multi-tenant flow (brands + auth)

- **Public:** `GET /health`, `POST /brands` — no brand or auth.
- **Brand-scoped:** All other routes need the current brand. Send **`X-Brand-Domain: <domain>`** (e.g. `acme`) on every request, or use a subdomain (e.g. `acme.localhost:8090`).
- **Login:** `POST /login` with brand domain and password → server sets an **HTTP-only session cookie**. Use the same `X-Brand-Domain` and send the cookie on subsequent requests (Postman/browser do this automatically).
- **Protected:** Pages, widgets, `GET /brands/me`, `GET /brands/:id` require a valid session (cookie) and that the token’s brand matches the request’s brand.

## API overview

| Method | Path                         | Description                             |
| ------ | ---------------------------- | --------------------------------------- |
| GET    | `/health`                    | Health check                            |
| POST   | `/brands`                    | Create a brand (public)                 |
| POST   | `/login`                     | Login (brand-scoped; sets cookie)       |
| POST   | `/logout`                    | Logout (brand-scoped; clears cookie)   |
| GET    | `/brands/me`                 | Current brand (protected)              |
| GET    | `/brands/:id`                | Brand by ID, same brand only (protected)|
| POST   | `/pages`                     | Create a page (protected)              |
| GET    | `/pages`                     | List pages (protected)                 |
| GET    | `/pages/:id`                 | Get page by ID (protected)             |
| PUT    | `/pages/:id`                 | Update a page (protected)              |
| DELETE | `/pages/:id`                 | Delete a page (protected)              |
| POST   | `/pages/:id/widgets`         | Add widget (protected)                  |
| PUT    | `/widgets/:id`               | Update a widget (protected)            |
| DELETE | `/widgets/:id`               | Delete a widget (protected)            |
| POST   | `/pages/:id/widgets/reorder` | Reorder widgets (protected)            |

- **GET /pages** – Optional `?page=1&limit=10` for paginated response `{ "data", "total", "page", "limit" }`.
- **GET /pages/:id** – Optional `?widget_type=banner` to filter widgets by type.

**Widget types:** `banner`, `product_grid`, `text`, `image`, `spacer`

## Quick run-through (local)

Run these in order. Base URL: `http://localhost:8090`. Use `-c cookies.txt` to save the session cookie and `-b cookies.txt` to send it. If "brand domain already exists", skip step 2 and use that domain (e.g. `acme`) in steps 3–5.

```bash
# 1. Health
curl -s http://localhost:8090/health

# 2. Create a brand (public; include email and password)
curl -s -X POST http://localhost:8090/brands \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Store","logo":"https://example.com/logo.png","office_address":"123 Main St","domain":"acme","email":"admin@acme.com","password":"secret123"}'

# 3. Login (use X-Brand-Domain; cookie is set in response)
curl -s -c cookies.txt -X POST http://localhost:8090/login \
  -H "Content-Type: application/json" \
  -H "X-Brand-Domain: acme" \
  -d '{"email":"test@example.com","password":"admin123"}'

# 4. Create a page (send cookie + X-Brand-Domain)
curl -s -b cookies.txt -X POST http://localhost:8090/pages \
  -H "Content-Type: application/json" \
  -H "X-Brand-Domain: acme" \
  -d '{"name":"Home","route":"/home","is_home":true}'

# 5. List pages (use PAGE_ID from step 4 if needed for other calls)
curl -s -b cookies.txt http://localhost:8090/pages \
  -H "X-Brand-Domain: acme"
```
If step 4 fails with "duplicate key" or "route already exists", use a different `route` (e.g. `"/about"`) in the create-page payload.

Replace `PAGE_ID` and `WIDGET_ID` in the examples below with IDs from responses.

## Example API requests (with brand + auth)

For protected routes, include **`X-Brand-Domain: acme`** and send the session cookie (e.g. `-b cookies.txt` in curl, or use Postman’s cookie handling).

```bash
# Get current brand
curl -s -b cookies.txt http://localhost:8090/brands/me -H "X-Brand-Domain: acme"

# List pages with pagination
curl -s -b cookies.txt "http://localhost:8090/pages?page=1&limit=10" -H "X-Brand-Domain: acme"

# Get page with widgets
curl -s -b cookies.txt "http://localhost:8090/pages/PAGE_ID" -H "X-Brand-Domain: acme"

# Add widget
curl -s -b cookies.txt -X POST http://localhost:8090/pages/PAGE_ID/widgets \
  -H "Content-Type: application/json" \
  -H "X-Brand-Domain: acme" \
  -d '{"type":"banner","position":0,"config":{"title":"Welcome"}}'

# Update page
curl -s -b cookies.txt -X PUT http://localhost:8090/pages/PAGE_ID \
  -H "Content-Type: application/json" \
  -H "X-Brand-Domain: acme" \
  -d '{"name":"Home Updated","route":"/home"}'

# Logout
curl -s -b cookies.txt -X POST http://localhost:8090/logout -H "X-Brand-Domain: acme"
```

## Tests

Run all tests:

```bash
go test -v ./...
```

- With `DATABASE_URL` set (e.g. from `.env`), full API tests run against the database.
- Without a database, tests that need the DB are skipped; health and validation tests still run.
- 