# Personal Website Backend

A robust backend service for a personal portfolio website, built with Go (Golang) and the Gin web framework. This service manages content for blog posts, projects, skills, user profiles, and contact messages.

## 🚀 Tech Stack

- **Language**: Go (Golang) 1.21+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **Authentication**: JWT (JSON Web Tokens)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Logging**: [Zap](https://github.com/uber-go/zap)
- **API Documentation**: [Swagger](https://github.com/swaggo/swag)

## 📋 Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://go.dev/dl/) (version 1.21 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Git](https://git-scm.com/downloads)

## 🛠️ Setup & Installation

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd personal-backend
    ```

2.  **Environment Configuration**
    Create a `.env` file in the root directory (or copy a template if available).
    ```env
    PORT=8080
    ENV=development
    
    # Database
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=personal_db
    DB_SSLMODE=disable

    # JWT
    JWT_SECRET=your_super_secret_key
    JWT_EXPIRATION_HOURS=24
    ```

3.  **Database Setup**
    Ensure your PostgreSQL service is running and create the database:
    ```sql
    CREATE DATABASE personal_db;
    ```
    The application will automatically migrate the schema on startup using Gorm, or you can use the migration scripts in `migrations/`.

4.  **Install Dependencies**
    ```bash
    go mod download
    ```

## ▶️ Running the Application

To run the server in development mode:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env`).

## 📚 API Documentation

### Swagger UI
When the server is running, you can access the interactive API documentation at:
**[http://localhost:8080/api/public/docs/index.html](http://localhost:8080/api/public/docs/index.html)**

### Public vs. Protected Endpoints
- **Public**: `/api/public` (e.g., viewing posts, projects, profile). No authentication required.
- **Admin**: `/api/admin` (e.g., creating posts, updating profile). Requires a valid JWT token in the `Authorization` header (`Bearer <token>`).

### 👨‍💻 Developer Guide: Updating Swagger Docs

If you modify the API handlers and want to update the Swagger documentation, first install the `swag` CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Then, generate the documentation from the project root:

```bash
swag init -g cmd/server/main.go
```

This will update the files in the `docs/` folder.

## 📂 Project Structure

```
.
├── cmd/server/         # Entry point (main.go)
├── docs/               # Swagger documentation files
├── internal/
│   ├── config/         # Configuration loading
│   ├── database/       # Database connection
│   ├── middleware/     # HTTP middleware (Auth, CORS, etc.)
│   ├── modules/        # Feature modules (Auth, Posts, Projects, etc.)
│   ├── routes/         # API route definitions
│   └── utils/          # Utility functions (Response, Pagination)
├── migrations/         # SQL migration files
├── storage/            # Local storage for uploaded files
├── go.mod              # Go module definition
└── README.md           # Project documentation
```
