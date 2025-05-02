# Go Fiber-GORM Boilerplate with JWT Authentication

A modern, structured Go application boilerplate using Fiber as the web framework and GORM for database interactions. Features a complete JWT authentication system with access/refresh tokens following clean architecture principles.

## Features

- **Structured Project Layout**: Following Go best practices for project organization
- **JWT Authentication**: Complete authentication flow with access and refresh tokens
- **Database Integration**: Using GORM with SQLite (easily configurable for other databases)
- **Config Management**: Environment-based configuration with Viper
- **Validation**: Request validation with custom validators
- **Testing**: Comprehensive test setup for authentication and endpoints
- **Makefile**: Various commands for development and testing
- **Docker Support**: Ready for containerized deployment

## Project Structure

```
.
├── cmd/
│   └── main.go                # Application entry point
├── internal/                  # Private application code
│   ├── config/                # Configuration management
│   ├── database/              # Database connection and setup
│   ├── handlers/              # HTTP request handlers
│   ├── logger/                # Logging setup
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Data models
│   ├── repository/            # Data access layer
│   ├── services/              # Business logic
│   ├── tests/                 # Test helpers and tests
│   └── validators/            # Input validation
├── Makefile                   # Development and build commands
├── go.mod                     # Go module definition
└── .env                       # Environment variables (create this)
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Make (optional, for using the Makefile commands)
- Docker and Docker Compose (optional, for containerized deployment)

### Setup

#### Local Development

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd fiber-gorm
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the project root:
   ```
   ENVIRONMENT=development
   DB_DRIVER=sqlite
   DB_SOURCE=app.db
   SERVER_PORT=3000
   LOG_LEVEL=info
   JWT_SECRET=your-secret-key-change-this-in-production
   ```

4. Run the application:
   ```bash
   make run
   ```

#### Docker Deployment

1. Build the Docker image:
   ```bash
   docker build -t fiber-gorm-app .
   ```

2. Run the container:
   ```bash
   docker run -p 3000:3000 -v $(pwd)/data:/app/data fiber-gorm-app
   ```

3. For a more complete setup, create a `docker-compose.yml` file:
   ```yaml
   version: '3.8'
   
   services:
     app:
       build: .
       ports:
         - "3000:3000"
       volumes:
         - ./data:/app/data
       environment:
         - ENVIRONMENT=production
         - JWT_SECRET=your-secure-secret-key
   ```

4. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

## Authentication Flow

### Register
```bash
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123!"
}
```

### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePassword123!"
}
```

### Refresh Token
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

### Protected Route
```bash
GET /api/me
Authorization: Bearer your-access-token
```

## Development

### Available Commands

```bash
# Run with hot reload
make dev

# Build the application
make build

# Run without hot reload
make run

# Clean build artifacts
make clean

# Clean up dependencies
make tidy

# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run tests matching a pattern
make test-filter

# Run tests in watch mode (requires fswatch)
make test-watch
```

## Authentication Details

The authentication system uses JWT tokens with the following characteristics:

- **Access Token**: Short-lived (15 minutes), contains user details
- **Refresh Token**: Long-lived (7 days), used to obtain new access tokens
- **Custom Claims**: Includes email, username, and role information
- **Secure Password Storage**: Using bcrypt for password hashing

## Testing

The project includes a comprehensive test suite for the authentication flow. Run the tests with:

```bash
make test
```

Generate a coverage report with:

```bash
make test-coverage
```

## Extending the Boilerplate

### Adding New Endpoints

1. Create a new handler in `internal/handlers/`
2. Add the route in `cmd/main.go`
3. Add tests in `internal/tests/`

### Adding Database Migrations

1. Create migration files in `internal/database/migrations/`
2. Add the migration logic to `internal/database/database.go`

### Customizing Docker Setup

1. Modify the `Dockerfile` to add additional dependencies
2. Update environment variables in the `docker-compose.yml`
3. Add additional services (like Redis or PostgreSQL) as needed

## Security Considerations

- Change the JWT secret key in production
- Store refresh tokens securely (HTTP-only cookies recommended)
- Consider adding rate limiting for authentication endpoints
- Implement token blacklisting for better security
- Use appropriate Docker security practices:
  - Don't run containers as root
  - Use container scanning tools
  - Keep base images up to date
  - Use multi-stage builds to minimize attack surface

## Deployment Strategies

### Basic Deployment

For simple deployments, the included Docker setup provides a great starting point:

```bash
docker build -t fiber-gorm-app .
docker run -d -p 3000:3000 fiber-gorm-app
```

### Production Deployment

For production, consider these additional steps:

1. Use a reverse proxy like Nginx in front of the Go application
2. Set up proper TLS/SSL certificates
3. Use a more robust database like PostgreSQL or MySQL
4. Implement health checks for container orchestration
5. Set up monitoring and logging

### Kubernetes Deployment

The application can be deployed to Kubernetes using the Docker image. Create a basic deployment with:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fiber-gorm-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fiber-gorm-app
  template:
    metadata:
      labels:
        app: fiber-gorm-app
    spec:
      containers:
      - name: fiber-gorm-app
        image: fiber-gorm-app:latest
        ports:
        - containerPort: 3000
        env:
        - name: ENVIRONMENT
          value: production
```

## License

[MIT License](LICENSE)
