# Simple HTTP Server Boilerplate

A robust boilerplate for creating an HTTP server in Go, built with [Chi](https://go-chi.io/), focusing on simplicity, ease of use, and modern best practices.

## Features

This boilerplate comes packed with the following features:

- **Clean Architecture:** Organized into handlers, services, and storage layers for maintainability and scalability.
- **Basic Authentication Middleware:** Secure your routes with basic HTTP authentication.
- **CORS Middleware:** Configured for Cross-Origin Resource Sharing, allowing flexible frontend integration.
- **Rate Limiting Middleware:** Protect your API from abuse and ensure fair usage with request rate limiting.
- **Prometheus Metrics:** Exposes detailed application metrics (total requests, request duration, status codes) at the `/metrics` endpoint for robust monitoring.
- **Swagger (OpenAPI) Documentation:** Automatically generated and served at `/swagger/*` for easy API exploration and understanding.
- **Graceful Shutdown:** Ensures the server shuts down cleanly upon receiving termination signals, allowing active requests to complete without interruption.
- **Environment-based Configuration:** Utilizes `viper` to manage configurations, loading settings from `config.{environment}.yaml` files and environment variables, supporting `development` and `production` environments.
- **Docker Compose Setup:** Simplifies local development by providing a `docker-compose.yml` to spin up the application, PostgreSQL database, and Redis cache with a single command.
- **Redis Caching:** Integrated into the `UserService` to cache user data, significantly improving performance for read operations and including cache invalidation for write operations.
- **Kubernetes YAMLs:** Provides foundational Kubernetes Deployment, Service, and Secret definitions (`k8s/deployment.yaml`, `k8s/service.yaml`, `k8s/db-secret.yaml`) for seamless CI/CD integration and deployment to a Kubernetes cluster.

## Configuration

The application's configuration is managed via `config/{environment}.yaml` files and environment variables. The `APP_ENV` environment variable determines which configuration file is loaded (defaults to `development`).

Example `config.development.yaml`:

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  db_name: demo

redis:
  host: localhost
  port: 6379

log_level: debug # can be debug, info, warn, or error
```

## Usage

### Local Development with Docker Compose

To get the application, PostgreSQL, and Redis running locally using Docker Compose:

```bash
docker compose up --build
```

The application will be accessible at `http://localhost:8081` (or the port configured in `docker-compose.yml`).

### Running Locally (without Docker Compose)

1.  **Install Dependencies:**

    ```bash
    go mod tidy
    ```

2.  **Set Environment Variables:**

    Ensure your `APP_ENV` is set (e.g., `export APP_ENV=development`) and provide database and Redis connection details either in `config.development.yaml` or as environment variables (e.g., `DB_HOST`, `DB_PORT`, `REDIS_HOST`, `REDIS_PORT`).

3.  **Run the Application:**

    ```bash
    go run main.go
    ```

### Database Migrations

This project includes a custom Go-based system to manage database schema changes. Migration files are plain SQL located in the `migrations/` directory.

To apply all pending migrations, run:

```bash
make migrate
```

The command will:

1. Connect to the database.
2. Create a `schema_migrations` table if it doesn't exist.
3. Check which migrations have already been applied.
4. Run any new `.up.sql` migration files in sequential order.

### Kubernetes Deployment

Basic Kubernetes manifests are provided in the `k8s/` directory:

- `k8s/deployment.yaml`: Defines the application deployment.
- `k8s/service.yaml`: Defines the Kubernetes service for the application.
- `k8s/db-secret.yaml`: A template for creating a Kubernetes Secret for database credentials. **Remember to replace placeholder values and handle this file securely in a real CI/CD pipeline.**

To deploy to Kubernetes, you would typically apply these files:

```bash
kubectl apply -f k8s/db-secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## Endpoints

- `GET /`: Home page.
- `GET /health`: Health check.
- `GET /metrics`: Prometheus metrics endpoint.
- `GET /swagger/*`: Swagger UI for API documentation.
- `GET /users`: List users (requires Basic Auth: `admin:password`).
- `POST /users`: Create a new user (requires Basic Auth: `admin:password`).
- `GET /users/{id}`: Get a user by ID (requires Basic Auth: `admin:password`).
- `DELETE /users/{id}`: Delete a user by ID (requires Basic Auth: `admin:password`).

## API Testing with httpyac

A `requests.http` file is provided with sample HTTP requests to test the API endpoints. You can use extensions like "REST Client" for VS Code or "HTTP Client" for IntelliJ IDEA to run these requests directly from your editor.

To use it:

1.  Open the `requests.http` file in your IDE.
2.  Click on the "Send Request" or similar button above each request.

## Development

A `Makefile` is provided to simplify common development tasks:

- `make build`: Build the application.
- `make run`: Run the application.
- `make test`: Run the tests.
- `make clean`: Clean the build artifacts.
- `make migrate`: Run database migrations.
- `make help`: Display the help message.

## Contributing

Feel free to fork this repository and contribute!
