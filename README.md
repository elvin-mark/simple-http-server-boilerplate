# Simple HTTP Server Boilerplate

A simple boilerplate for creating an HTTP server in Go, with a focus on simplicity and ease of use.

## Configuration

The application can be configured using the `config.yaml` file. The following options are available:

```yaml
server:
  port: 8080
log_level: info # can be debug, info, warn, or error

database:
  host: localhost
  port: 5432
  user: postgres
  password: 
  dbname: demo
```

## Usage

A `Makefile` is provided to simplify the development process. The following commands are available:

* `make build`: Build the application.
* `make run`: Run the application.
* `make test`: Run the tests.
* `make clean`: Clean the build artifacts.
* `make help`: Display the help message.

## Docker

To build and run the application using Docker, use the following commands:

```bash
docker build -t simple-http-server .
docker run -p 8080:8080 simple-http-server
```

## Endpoints

The following endpoints are available:

* `GET /`: Home page.
* `GET /health`: Health check.
* `GET /users`: List users.
* `POST /users`: Create a new user.
* `GET /users/{id}`: Get a user by ID.
* `DELETE /users/{id}`: Delete a user by ID.
