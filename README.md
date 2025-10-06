# Simple Go API

A clean, modular REST API built with Go that demonstrates best practices for code organization and database integration.

## � Quick Start

### Prerequisites
- **For Local Development**: Go 1.22+
- **For Docker**: Docker & Docker Compose

### Option 1: Using Docker Compose (Recommended)

1. **Clone and navigate to the project:**
   ```bash
   cd simple-go-api
   ```

2. **Start the application with database:**
   ```bash
   docker compose up --build
   ```

   **Alternative commands:**
   ```bash
   # Start in background (detached mode)
   docker compose up -d --build
   
   # Start without building (if already built)
   docker compose up
   
   # Build only (without starting)
   docker compose build
   ```

   To include the optional React frontend container, add the `frontend` profile:
   ```bash
   docker compose --profile frontend up --build
   ```

3. **Access the services:**
   - **API**: http://localhost:8080
   - **Database Admin (phpMyAdmin)**: http://localhost:8081
   - **Frontend UI** (if started with the `frontend` profile): http://localhost:3000

### Option 2: Local Development

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the application:**
   ```bash
   go run main.go
   ```

3. **Or build and run:**
   ```bash
   go build -o server main.go
   ./server
   ```

## 🧪 Test the API

### Health Check
```bash
curl http://localhost:8080/
```

### Get All Articles
```bash
curl http://localhost:8080/articles
```

### Get Article by ID
```bash
curl http://localhost:8080/article/1
```

### Filter Articles
```bash
curl http://localhost:8080/article/filter/category/technology
```

### Create User
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","id":42}'
```

## 🗄️ Database Access

When running with Docker Compose:
- **phpMyAdmin**: http://localhost:8081
- **MySQL Direct**: localhost:3306
- **Credentials**: root/root_password or api_user/api_password

## 🛑 Stop the Application

### Docker Compose
```bash
# Stop and remove containers
docker compose down

# Stop and remove containers + volumes (removes database data)
docker compose down -v

# Stop containers (keep them for restart)
docker compose stop

# Restart stopped containers
docker compose start

# Start with frontend
docker compose --profile frontend up -d

```

### Local Development
Press `Ctrl+C` in the terminal

---

**That's it! Your API is ready to use. 🚀**

## 🌐 CORS configuration

The API now includes built-in CORS handling so the React frontend (or any external client) can call it directly.

- **Allow all origins (default):** no extra configuration required.
- **Restrict origins:** set `CORS_ALLOWED_ORIGINS` to a comma-separated list (e.g. `http://localhost:3000,https://example.com`).

Remember to restart the API container or process after changing the environment variable so the new policy is applied.