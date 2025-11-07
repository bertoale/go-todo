# Go Todo REST API

A simple RESTful API for a Todo application built with Go, Fiber, and GORM (MySQL). This project implements user authentication (JWT), task management, and follows clean architecture principles (repository, service, controller).

## Features

- User registration & login (JWT authentication)
- CRUD tasks (create, read, update, delete)
- Each task belongs to a user
- Secure endpoints (protected with JWT)
- Clean code structure (separation of concerns)
- Error handling & validation

## Tech Stack

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/) (web framework)
- [GORM](https://gorm.io/) (ORM for MySQL)
- [MySQL](https://www.mysql.com/) (database)
- JWT for authentication

## Project Structure

```
internal/
  controllers/   # HTTP handlers
  services/      # Business logic
  repositories/  # Database access
  models/        # Data models
  middlewares/   # Fiber middlewares (auth, error)
  dto/           # Request/response DTOs
  routes/        # Route definitions
  database/      # DB connection & migration
config/          # App configuration
cmd/             # Main entrypoint
```

## Getting Started

### Prerequisites

- Go 1.18+
- MySQL server

### Installation

1. Clone the repo:
   ```bash
   git clone <repo-url>
   cd go-todo
   ```
2. Copy `.env` and set your database credentials:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=yourpassword
   DB_NAME=your_db_name
   JWT_SECRET=your_jwt_secret
   JWT_EXPIRES_IN=168h
   PORT=5000
   NODE_ENV=development
   CORS_ORIGIN=http://localhost:3000
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the app:
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### Auth

- `POST /api/auth/register` — Register new user
- `POST /api/auth/login` — Login, returns JWT token

### User

- `GET /api/users/` — Get current user profile (JWT required)
- `PUT /api/users/:id` — Update user profile (JWT required)

### Tasks

- `POST /api/tasks/` — Create new task (JWT required)
- `GET /api/tasks/` — List all tasks for current user (JWT required)
- `GET /api/tasks/:id` — Get task by ID (JWT required)
- `PUT /api/tasks/:id` — Update task (JWT required)
- `DELETE /api/tasks/:id` — Delete task (JWT required)

## License

MIT
