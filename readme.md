# Food Delivery API

A REST API for browsing food delivery orders. Built with Go, Gin, and PostgreSQL, with JWT authentication and interactive API documentation in Swagger UI.

## Features

- Register, log in, log out, and change passwords.
- Protect order endpoints with JWT bearer authentication.
- Revoke the current access token on logout.
- Browse orders with pagination and optional filters.
- Retrieve detailed order information by `order_id`.
- Load sample orders from CSV into PostgreSQL.

## Run with Docker Compose

### 1. Configure the environment

Install Docker Desktop or Docker Engine with the Compose plugin. Run the following commands from the repository root in PowerShell.

```powershell
if (-not (Test-Path .env)) { Copy-Item .env.example .env }
```

Set `POSTGRES_PASSWORD` and `JWT_SECRET` in `.env`. Replace the example values with your own. The JWT secret must be at least 32 bytes. To generate a random value in PowerShell, run the following and paste the output into `JWT_SECRET`.

```powershell
$bytes = [byte[]]::new(32)
$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
$rng.GetBytes($bytes)
$rng.Dispose()
[Convert]::ToBase64String($bytes)
```

`API_PORT` defaults to `8080` and `POSTGRES_PORT` defaults to `5432`. Keep `.env` out of version control.

### 2. Start PostgreSQL and initialize the schema

```powershell
docker compose up -d --wait db
docker compose exec -T db psql -v ON_ERROR_STOP=1 -U food_delivery -d food_delivery -f /setup/schema.sql
```

To load sample orders, run the seed script after the schema is ready.

```powershell
docker compose exec -T db psql -v ON_ERROR_STOP=1 -U food_delivery -d food_delivery -f /setup/seed.sql
```

### 3. Build and start the API

```powershell
docker compose up -d --build api-go
docker compose ps
```

- [Swagger UI](http://localhost:8080/swagger/)
- [Health check](http://localhost:8080/health)

If you changed `API_PORT`, use that port in both URLs.

## Endpoints

| Method | Path                    | JWT required | Purpose                                        |
| ------ | ----------------------- | ------------ | ---------------------------------------------- |
| `GET`  | `/health`               | No           | Check API and database connectivity            |
| `POST` | `/auth/register`        | No           | Create a user account                          |
| `POST` | `/auth/login`           | No           | Return an access token                         |
| `POST` | `/auth/logout`          | Yes          | Revoke the current access token                |
| `POST` | `/auth/change-password` | Yes          | Update the authenticated user's password       |
| `GET`  | `/orders`               | Yes          | List and filter orders                         |
| `GET`  | `/orders/{id}`          | Yes          | Get details using the `order_id` from the list |

### Try it in Swagger UI

1. Call `POST /auth/register` with `email`, `password`, and `password_confirm`.
2. Call `POST /auth/login` and copy the returned `access_token`.
3. Click **Authorize** and paste the token. Swagger adds the `Bearer` prefix.
4. Call `GET /orders`. Copy an `order_id` from the response to try `GET /orders/{id}`.

### Filter orders

`pageindex` starts at `1`. `pagesize` defaults to `20` and accepts values from `1` to `100`. The optional filters are `cuisine_category`, `order_id`, `metro_city`, and `order_channel`. Filters use exact matches and are combined with `AND`. Leave unused fields empty.

```http
GET /orders?pageindex=1&pagesize=20&metro_city=Bangkok
```

## Project layout

| Path                          | Contents                                                       |
| ----------------------------- | -------------------------------------------------------------- |
| `api-go/cmd/api/`             | API entry point and route registration                         |
| `api-go/internal/auth/`       | Authentication handlers, service, repository, and JWT handling |
| `api-go/internal/middleware/` | JWT validation and revocation checks                           |
| `api-go/internal/orders/`     | Order handlers and database queries                            |
| `api-go/docs/`                | Swagger UI and OpenAPI specification                           |
| `sql/`                        | PostgreSQL schema and seed scripts                             |
| `data/`                       | Sample order CSV                                               |
| `compose.yml`                 | Local API and PostgreSQL services                              |

## View logs and stop services

```powershell
docker compose logs --tail 50 api-go
docker compose down
```

`docker compose down` preserves PostgreSQL data in the named volume.
