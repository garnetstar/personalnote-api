# Simple Go API with React Frontend

A clean, modular REST API built with Go and a React TypeScript frontend that demonstrates best practices for full-stack development.

## � Quick Start

### Prerequisites
- **For Local Development**: Go 1.22+, Node.js 18+, npm
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

3. **Access the services:**
   - **React Frontend**: http://localhost:3000
   - **API**: http://localhost:8080
   - **Database Admin (phpMyAdmin)**: http://localhost:8081

### Option 2: Local Development

#### Backend (Go API)
1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the API:**
   ```bash
   go run main.go
   ```

3. **Or build and run:**
   ```bash
   go build -o server main.go
   ./server
   ```

#### Frontend (React)
1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start development server:**
   ```bash
   npm start
   ```

4. **Access the application:**
   - **React Frontend**: http://localhost:3000
   - **Go API**: http://localhost:8080

## 🧪 Test the API

### Using the React Frontend
Visit http://localhost:3000 to use the interactive web interface that provides:
- API health checking
- Article browsing and filtering
- Article details view
- User creation form

### Using Direct API Calls

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
```

### Local Development
- **React Frontend**: Press `Ctrl+C` in the frontend terminal
- **Go API**: Press `Ctrl+C` in the API terminal

---

**That's it! Your API is ready to use. 🚀**