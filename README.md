# Simple Go API

A clean, modular REST API built with Go that demonstrates best practices for code organization and structure.

## 🚀 Features

- **Clean Architecture** - Well-organized code structure with separation of concerns
- **JSON REST API** - RESTful endpoints with JSON request/response
- **Input Validation** - Request validation with detailed error messages
- **Docker Support** - Containerized application with Docker Compose
- **Logging** - Comprehensive request logging to console
- **Health Check** - Simple health check endpoint
- **User Management** - User data processing with validation

## 📁 Project Structure

```
simple-go-api/
├── main.go                 # 🎯 Application entry point
├── models/
│   └── models.go          # 📊 Data structures (User, Response, ErrorResponse)
├── handlers/
│   └── handlers.go        # 🎛️ HTTP request handlers
├── utils/
│   ├── validation.go      # ✅ Input validation logic
│   ├── response.go        # 📤 JSON response utilities
│   └── database.go        # 🗄️ Database connection utilities
├── router/
│   └── router.go          # 🚏 Route configuration
├── docker-compose.yml     # 🐳 Docker Compose configuration
├── Dockerfile             # 🐳 Docker build instructions
├── init.sql               # 🗄️ Database initialization script
├── .env.example           # ⚙️ Environment variables example
└── go.mod                 # 📦 Go module dependencies
```

## 🛠️ Technology Stack

- **Language**: Go 1.22+
- **Framework**: Net/HTTP (standard library)
- **Database**: MySQL 8.0
- **Database Driver**: go-sql-driver/mysql
- **Container**: Docker & Docker Compose
- **Base Image**: Alpine Linux (lightweight)
- **Database Admin**: phpMyAdmin

## 📚 API Endpoints

### 🟢 GET `/` - Health Check
Simple health check endpoint that returns a greeting message.

**Response:**
```json
{
  "message": "Hallo, from Go!"
}
```

**Headers:**
- `Counter`: Request counter (increments with each request)

### 🟡 POST `/user` - User Management
Accepts user data, validates it, and processes the request.

**Request Body:**
```json
{
  "name": "jan",
  "id": 1
}
```

**Success Response (200):**
```json
{
  "message": "User jan with ID 1 has been processed successfully"
}
```

**Validation Error Response (400):**
```json
{
  "error": "Validation failed",
  "message": "Validation errors: name is required, id must be a positive integer"
}
```

**Method Error Response (405):**
```json
{
  "error": "Method not allowed",
  "message": "Only POST requests are accepted"
}
```

## 🗄️ Database Services

The application includes a complete database setup with:

### MySQL Database
- **Host**: `localhost:3306`
- **Database**: `simple_go_api`
- **User**: `api_user`
- **Password**: `api_password`
- **Root Password**: `root_password`

### phpMyAdmin
- **URL**: `http://localhost:8081`
- **Login**: Use root credentials or api_user credentials
- **Features**: 
  - Database management interface
  - SQL query execution
  - Data visualization
  - Table structure management

### Database Schema
```sql
-- Users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    external_id INT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Request logs table
CREATE TABLE request_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    request_body TEXT,
    response_status INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🚀 Quick Start

### Prerequisites
- Go 1.22+ (for local development)
- Docker & Docker Compose (recommended)

### Option 1: Using Docker Compose (Recommended)

1. **Clone and navigate to the project:**
   ```bash
   cd simple-go-api
   ```

2. **Start the application:**
   ```bash
   docker compose up --build
   ```

3. **Access the services:**
   ```bash
   # API Health check
   curl http://localhost:8080/
   
   # API User endpoint
   curl -X POST http://localhost:8080/user \
     -H "Content-Type: application/json" \
     -d '{"name":"jan","id":1}'
   
   # phpMyAdmin (Database Management)
   open http://localhost:8081
   ```

4. **Database connection:**
   - **phpMyAdmin**: `http://localhost:8081`
   - **MySQL Direct**: `localhost:3306`
   - **Credentials**: root/root_password or api_user/api_password

### Option 2: Local Development

1. **Run locally:**
   ```bash
   go run main.go
   ```

2. **Build binary:**
   ```bash
   go build -o server main.go
   ./server
   ```

## 🧪 Testing Examples

### Valid User Request
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","id":42}'
```
**Response:** `{"message":"User Alice with ID 42 has been processed successfully"}`

### Invalid Requests

**Empty name:**
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"","id":1}'
```
**Response:** `{"error":"Validation failed","message":"Validation errors: name is required"}`

**Invalid ID:**
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"jan","id":0}'
```
**Response:** `{"error":"Validation failed","message":"Validation errors: id must be a positive integer"}`

**Wrong HTTP method:**
```bash
curl -X GET http://localhost:8080/user
```
**Response:** `{"error":"Method not allowed","message":"Only POST requests are accepted"}`

## 📝 Console Logging

The application provides detailed logging to help with debugging:

```
🚀 Server starting on :8080
📡 Endpoints available:
   GET  / - Health check
   POST /user - User management

Received request #1
Query params: map[]
URI: /
✅ Received valid user data: Name=jan, ID=1
📝 Full request body: {"name":"jan","id":1}
Validation failed: [name is required]
```

## 🏗️ Code Architecture

### Models Package (`models/`)
Contains all data structures used throughout the application:
- `User` - User entity with name and ID
- `Response` - Standard success response
- `ErrorResponse` - Error response with details

### Handlers Package (`handlers/`)
HTTP request handlers that process incoming requests:
- `HelloHandler` - Handles health check requests
- `UserHandler` - Processes user management requests

### Utils Package (`utils/`)
Utility functions for common operations:
- `validation.go` - Input validation logic
- `response.go` - JSON response helpers

### Router Package (`router/`)
Route configuration and setup:
- `SetupRoutes()` - Configures all application routes

## 🐳 Docker Configuration

The application uses a multi-stage Docker build for optimal image size:

1. **Build Stage**: Uses `golang:1.22.2-alpine` to compile the application
2. **Runtime Stage**: Uses `alpine:latest` for a minimal runtime environment
3. **Security**: Runs as non-root user `app`
4. **Size**: Final image is ~15MB

## 🔧 Development

### Adding New Endpoints

1. **Add route** in `router/router.go`:
   ```go
   http.HandleFunc("/new-endpoint", handlers.NewHandler)
   ```

2. **Create handler** in `handlers/handlers.go`:
   ```go
   func NewHandler(w http.ResponseWriter, r *http.Request) {
       // Handler logic
   }
   ```

3. **Add models** if needed in `models/models.go`

### Adding Validation

Add new validation rules in `utils/validation.go`:
```go
func ValidateNewEntity(entity NewEntity) []string {
    var errors []string
    // Add validation logic
    return errors
}
```

## 📊 Benefits of This Structure

✅ **Maintainable** - Each package has a single responsibility  
✅ **Testable** - Easy to unit test individual components  
✅ **Scalable** - Simple to add new features and endpoints  
✅ **Readable** - Clear separation makes code easy to understand  
✅ **Reusable** - Utility functions can be used across handlers  

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

**Happy coding! 🚀**