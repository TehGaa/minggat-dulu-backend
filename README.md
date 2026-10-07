# Minggat Dulu Backend

This is the backend repository for the Minggat Dulu application, built with Go. It follows a microservices architecture using the standard Go project layout.

## Microservices

- **User Service:** Handles user authentication and management.
- **Booking Service:** Handles the booking flow.

## Project Structure

```text
minggat-dulu-backend/
├── api/                  # OpenAPI/Swagger specs, Protocol Buffers, JSON schema files
├── cmd/                  # Main applications for this project
│   ├── booking/          # Entry point for the Booking Service
│   │   └── main.go
│   └── user/             # Entry point for the User Service
│       └── main.go
├── configs/              # Configuration file templates or default configs
├── deployments/          # Deployment configurations (e.g., Dockerfiles, docker-compose)
├── internal/             # Private application and library code
│   ├── booking/          # Core logic for the Booking Service
│   │   ├── handler/      # HTTP/gRPC handlers
│   │   ├── model/        # Domain models/entities
│   │   ├── repository/   # Database access layer
│   │   └── service/      # Business logic layer
│   ├── pkg/              # Internal packages shared across microservices
│   │   ├── config/       # Shared configuration parsing logic
│   │   ├── database/     # Shared database connections/utilities
│   │   └── logger/       # Shared logging setup
│   └── user/             # Core logic for the User Service
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── pkg/                  # Library code that is okay to use by external applications
├── scripts/              # Scripts for build, install, analysis, DB migrations, etc.
├── go.mod                # Go module definition
└── go.sum                # Go dependencies checksums
```

## Running the Services

You can start the services locally by running the `main.go` files from the root directory:

**Start the User Service:**
```powershell
go run cmd/user/main.go
```
*(The user service runs on `:8081` with API routes under `/api/v1`)*

**Start the Booking Service:**
```powershell
go run cmd/booking/main.go
```
*(The booking service runs on `:8082` with API routes under `/api/v1`)*
