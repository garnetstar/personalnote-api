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
 
   # start with frontend and do build
   docker compose --profile frontend up --build

   # start with frontend in debug mode
   docker compose up -d && cd ./frontend/app && npm start dev
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

## 🔐 Google OAuth Authentication

The application uses Google OAuth 2.0 for authentication. All article creation, editing, and deletion operations require authentication.

### Setup Google OAuth

1. **Create a Google Cloud Project:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one

2. **Enable Google+ API:**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google+ API" and enable it

3. **Create OAuth 2.0 Credentials:**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Web application"
   - Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
   - For production, add your production URL

4. **Configure Environment Variables:**
   Create a `.env` file in the project root (or set environment variables):
   ```bash
   GOOGLE_CLIENT_ID=your_google_client_id_here
   GOOGLE_CLIENT_SECRET=your_google_client_secret_here
   GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
   FRONTEND_URL=http://localhost:3000
   JWT_SECRET=your_super_secret_jwt_key_min_32_characters
   ```

5. **Restart the Application:**
   ```bash
   docker compose down
   docker compose --profile frontend up -d
   ```

### How Authentication Works

1. Users click "Sign in with Google" on the login page
2. They're redirected to Google's OAuth consent screen
3. After approval, Google redirects back to the API with an auth code
4. The API exchanges the code for user info and generates a JWT token
5. The frontend stores the token and includes it in all protected requests
6. Article operations (create/edit/delete) require valid authentication

### Protected Endpoints

- **POST** `/articles` - Create new article (requires auth)
- **PUT** `/article/{id}` - Update article (requires auth)
- **DELETE** `/article/{id}` - Delete article (requires auth)

### Public Endpoints

- **GET** `/articles` - List all articles
- **GET** `/article/{id}` - Get article details
- **GET** `/article/filter/{mode}/{keyword}` - Search articles

## 🌐 CORS configuration

The API now includes built-in CORS handling so the React frontend (or any external client) can call it directly.

- **Allow all origins (default):** no extra configuration required.
- **Restrict origins:** set `CORS_ALLOWED_ORIGINS` to a comma-separated list (e.g. `http://localhost:3000,https://example.com`).

Remember to restart the API container or process after changing the environment variable so the new policy is applied.