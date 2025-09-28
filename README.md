# Go Todo API

 Todo REST API written in Go with Gin, GORM, Cobra, and Uber Fx. The service exposes batch create, update, and list endpoints and persists data in MySQL (SQLite fallback is available for quick tests).

## Features

- Gin HTTP server with JSON responses
- Batch `POST /todos`, `PATCH /todos`, and paginated `GET /todos`
- MySQL persistence via GORM (automatic migrations)
- Cobra CLI with `api` and `migrate` commands
- Dockerfile and docker-compose for running the API plus MySQL

## Prerequisites

- Go 1.20 or newer (module targets 1.25)
- Docker Desktop (optional but recommended for MySQL)
- Make sure port `3306` is free if you use the bundled MySQL container

## Quick Start (Docker Compose)

```bash
docker compose up --build
```

This builds the API image (`Dockerfile`) and launches:

- `api` service – runs `./main api` inside the container
- `db` service – MySQL 8.0 with credentials in `docker-compose.yaml`

The API becomes available at <http://localhost:8080>.

To stop everything:

```bash
docker compose down
```

## Running Locally Without Docker

1. Ensure MySQL is running on your machine (create a `todo_db` database).
2. Copy `.env.example` to `.env` (see below) and adjust credentials.
3. Export the environment variables before starting the API:

   ```powershell
   Get-Content .env | ForEach-Object {
     if ($_ -match '^(?<k>[^#=]+)=(?<v>.*)$') {
       $env:$Matches['k'].Trim() = $Matches['v'].Trim()
     }
   }
   go run . api
   ```

If you prefer SQLite, set `DB_TYPE=sqlite`

## CLI Commands (Cobra)

```bash
go run . --help
go run . api
go run . migrate
```

When running inside Docker, the container executes `./main api`.

## Environment Variables

```
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_PASSWORD=password
DB_NAME=todo_db
DB_DSN=mysql://root:password@tcp(localhost:3306)/todo_db?parseTime=true&loc=Local
```

## API Overview

### POST /todos

Create one or more todos.

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{
        "todos": [
          {
            "title": "Buy groceries",
            "description": "Milk, eggs, bread",
            "due_date": "2025-09-30T17:00:00Z"
          }
        ]
      }'
```

### PATCH /todos

Update existing todos by ID.

```bash
curl -X PATCH http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{
        "todos": [
          {
            "id": 1,
            "complete": true
          }
        ]
      }'
```

### DELETE /todos/:id

Delete a todo by ID.

```bash
curl -X DELETE http://localhost:8080/todos/1
```

### GET /todos

List todos with pagination.

```bash
curl "http://localhost:8080/todos?page=1&limit=10"
```

Response structure:

```json
{
  "todos": [
    {
      "id": 1,
      "title": "Buy groceries",
      "description": "Milk, eggs, bread",
      "due_date": "2025-09-30T17:00:00Z",
      "complete": true,
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

## Migrations

To apply MySQL migrations from the host:

The `migrate` command expects `DB_DSN` (e.g. `mysql://user:pass@tcp(host:port)/db?parseTime=true&loc=Local`).

## Troubleshooting

- **SQLite requires CGO**: install a C compiler (MSYS2 + MinGW on Windows) and set `CGO_ENABLED=1`, or stick with MySQL.
- **Port 3306 already in use**: stop the existing MySQL service or change the compose mapping and `DB_PORT`.
- **Access denied**: verify the MySQL user/password and update both `DB_PASS` and `DB_DSN`.
- **Docker Desktop + WSL issues**: restarting `wsl --shutdown` and re-opening Docker Desktop usually clears the error; ensure WSL distros are initialized.

## Next Steps

- Add and GET-by-ID routes
- Add unit/integration tests and wire a CI workflow
- Improve validation and error handling (e.g., handle duplicate titles gracefully)
- Harden configuration (structured logging, graceful shutdown, CORS, health checks)



