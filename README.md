# CompareFlow

A data validation and comparison platform built as a single Go binary with an embedded React frontend. CompareFlow helps you validate data consistency between different databases and data sources.

## Documentation

📚 **Complete documentation is available in the [docs/](docs/) folder:**

### Getting Started
- [🚀 Run Instructions](docs/RUN_INSTRUCTIONS.md) - How to build and run the application
- [🛠️ Development Guide](docs/DEVELOPMENT_GUIDE.md) - Development setup and workflow  
- [📦 Deployment Guide](docs/DEPLOYMENT_GUIDE.md) - Production deployment instructions
- [🐳 Podman Quickstart](docs/PODMAN_QUICKSTART.md) - Container deployment with Podman

### Technical Documentation  
- [🏗️ Architecture](docs/ARCHITECTURE.md) - System architecture and design patterns
- [⚙️ Technical Design](docs/TECHNICAL_DESIGN.md) - Detailed technical specifications
- [📋 Functional Requirements](docs/FUNCTIONAL_REQUIREMENTS.md) - Feature requirements and specifications
- [📖 API Reference](docs/API_REFERENCE.md) - REST API endpoints and usage

### Reference
- [📚 Documentation Index](docs/DOCUMENTATION_INDEX.md) - Complete documentation overview

## Features

- **Multi-Database Support**: SQL Server, Databricks, and PostgreSQL with pluggable architecture for easy extension
- **Flexible Validations**: Row count, data matching, schema comparison
- **Single Binary Deployment**: Everything packaged in one executable
- **Modern UI**: React with Material-UI and Redux
- **Secure**: JWT authentication with bcrypt password hashing
- **Fast**: Built with Go for high performance
- **Extensible**: Clean connector interface makes adding new databases simple

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Node.js 18+ and npm (for frontend development)
- Podman (recommended) or Docker
- PostgreSQL 15+ (or use the provided container)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd compareflow
   ```

2. **Start PostgreSQL and run the application**
   ```bash
   ./start.sh
   ```
   This script will:
   - Start PostgreSQL container with Podman
   - Wait for database to be ready
   - Create necessary tables
   - Start the CompareFlow application

   The application will be available at http://localhost:8080
   Default login: admin / admin123

### Alternative Setup (Manual)

1. **Install dependencies**
   ```bash
   make deps
   ```

2. **Start database manually**
   ```bash
   podman run -d \
     --name compareflow-postgres \
     -e POSTGRES_USER=compareflow \
     -e POSTGRES_PASSWORD=compareflow123 \
     -e POSTGRES_DB=compareflow \
     -p 5432:5432 \
     postgres:15
   ```

3. **Run the application**
   ```bash
   make run
   ```

### Building for Production

1. **Build complete application (frontend + backend)**
   ```bash
   make build-all
   ```

2. **Create release builds for multiple platforms**
   ```bash
   make release
   ```

3. **Run with custom database**
   ```bash
   ./bin/compareflow -db "postgresql://user:pass@host:5432/dbname?sslmode=require"
   ```

## Configuration

### Environment Variables

Create a `.env` file or set these environment variables:

```bash
# Database
DATABASE_URL=postgresql://compareflow:compareflow123@localhost:5432/compareflow?sslmode=disable

# JWT Settings  
JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRATION_DAYS=7

# Server
PORT=8080
GIN_MODE=release  # Set to "debug" for development

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Command Line Flags

```bash
./compareflow -h
  -db string
        Database connection string (default "postgresql://...")
  -port string
        Server port (default "8080")
  -jwt-secret string
        JWT secret key (default "your-secret-key")
```

## API Documentation

### Authentication

**Login**
```bash
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin123"
}
```

**Register**
```bash
POST /api/v1/auth/register
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

### Connections

**List Connections**
```bash
GET /api/v1/connections
Authorization: Bearer <token>
```

**Create Connection**
```bash
POST /api/v1/connections
Authorization: Bearer <token>
{
  "name": "Production SQL Server",
  "type": "sqlserver",
  "config": {
    "server": "prod-sql.company.com",
    "port": 1433,
    "database": "SalesDB",
    "username": "readonly",
    "password": "password",
    "encrypt": true,
    "trust_server_certificate": false
  }
}
```

**Test Connection**
```bash
POST /api/v1/connections/:id/test
Authorization: Bearer <token>
```

### Validations

**Create Validation**
```bash
POST /api/v1/validations
Authorization: Bearer <token>
{
  "name": "Daily Sales Check",
  "source_connection_id": 1,
  "target_connection_id": 2,
  "config": {
    "comparison_type": "row_count",
    "source_query": "SELECT COUNT(*) FROM orders WHERE date = CURRENT_DATE",
    "target_query": "SELECT COUNT(*) FROM staging.orders WHERE date = CURRENT_DATE"
  }
}
```

**Run Validation**
```bash
POST /api/v1/validations/:id/run
Authorization: Bearer <token>
```

## Development

### Project Structure

```
compareflow/
├── cmd/compareflow/     # Main application entry point
├── internal/
│   ├── api/            # HTTP handlers, middleware, and routes
│   ├── config/         # Configuration management
│   ├── connectors/     # Database connector framework
│   ├── database/       # Database connection and migrations
│   ├── models/         # Data models
│   └── services/       # Business logic
├── frontend/           # React application
│   ├── src/
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Page components
│   │   ├── services/   # API client services
│   │   ├── store/      # Redux store and slices
│   │   └── types/      # TypeScript type definitions
│   └── dist/          # Built frontend assets
├── docs/              # Documentation
├── scripts/           # Utility scripts
├── migrations/        # Database migrations
└── bin/              # Build outputs
```

### Available Make Commands

```bash
make deps          # Install Go dependencies
make run           # Run the application in development mode
make build         # Build the Go binary
make build-all     # Build frontend and backend together
make frontend      # Build frontend and copy to web/dist
make test          # Run Go tests
make clean         # Clean build artifacts
make fmt           # Format Go code
make lint          # Run linter (requires golangci-lint)
make release       # Create release builds for multiple platforms
make dev           # Run with hot reload (requires air)
```

### Frontend Development

For frontend development with hot reload:

```bash
# Terminal 1: Run backend
make run

# Terminal 2: Run frontend dev server
cd frontend
npm install
npm start
```

The frontend dev server will proxy API requests to the backend.

### Adding New Features

1. **New Database Connector**: Add to `internal/connectors/`
2. **New API Endpoint**: Add handler in `internal/api/handlers/`
3. **New Model**: Add to `internal/models/`
4. **New Service**: Add to `internal/services/`
5. **Frontend Page**: Add to `frontend/src/pages/`
6. **Redux Slice**: Add to `frontend/src/store/slices/`

## Deployment

### Using Systemd (Linux)

Create `/etc/systemd/system/compareflow.service`:

```ini
[Unit]
Description=CompareFlow Data Validation Service
After=network.target

[Service]
Type=simple
User=compareflow
WorkingDirectory=/opt/compareflow
ExecStart=/opt/compareflow/compareflow
Restart=always
Environment="GIN_MODE=release"
Environment="DATABASE_URL=postgresql://..."

[Install]
WantedBy=multi-user.target
```

### Using Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o compareflow cmd/compareflow/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/compareflow .
EXPOSE 8080
CMD ["./compareflow"]
```

## Troubleshooting

### Database Connection Issues

1. Check PostgreSQL is running:
   ```bash
   podman ps
   # Should show compareflow-postgres container
   ```

2. Check container logs:
   ```bash
   podman logs compareflow-postgres
   ```

3. Verify connection string:
   ```bash
   podman exec compareflow-postgres psql -U compareflow -d compareflow -c "SELECT 1;"
   ```

4. Restart database:
   ```bash
   ./stop.sh
   ./start.sh
   ```

### Frontend Build Issues

1. Clear npm cache:
   ```bash
   npm cache clean --force
   ```

2. Reinstall dependencies:
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

### Authentication Issues

1. Verify JWT secret is set
2. Check token expiration
3. Ensure CORS origins are configured

## Current Status

### Completed Features
- ✅ Go project structure with clean architecture
- ✅ JWT authentication system
- ✅ User management (registration, login)
- ✅ Database connection management (CRUD)
- ✅ Validation job management (CRUD)
- ✅ Multi-database connector framework
- ✅ Database connectors: SQL Server, PostgreSQL, Databricks
- ✅ Connection testing and schema discovery
- ✅ PostgreSQL integration with Podman
- ✅ Development and deployment scripts
- ✅ React frontend with TypeScript
- ✅ Material-UI design system
- ✅ Redux state management
- ✅ Frontend asset embedding

### In Progress
- 🔄 Data validation engine implementation
- 🔄 Validation execution and reporting

### Planned Features
- 📋 WebSocket support for real-time updates
- 📋 Scheduling system for automated validations
- 📋 Email/Slack notifications
- 📋 Export results to CSV/Excel/PDF
- 📋 Data lineage visualization
- 📋 Additional database connectors (MySQL, Snowflake, Oracle)
- 📋 Advanced validation types (schema comparison, data profiling)
- 📋 Role-based access control
- 📋 Audit logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and feature requests, please use the GitHub issue tracker.