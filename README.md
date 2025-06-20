# CompareFlow

A data validation and comparison platform built as a single Go binary with an embedded React frontend. CompareFlow helps you validate data consistency between different databases and data sources.

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

- Go 1.21 or higher
- Node.js 18+ and npm (for frontend development)
- Podman or Docker (for PostgreSQL)
- PostgreSQL 13+ (or use the provided container)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/compareflow/compareflow.git
   cd compareflow
   ```

2. **Start PostgreSQL with Podman**
   ```bash
   ./start-local.sh
   ```

3. **Create admin user (optional)**
   ```bash
   psql -h localhost -U compareflow -d compareflow -f scripts/create_admin_user.sql
   # Default admin credentials: admin/admin123
   ```

4. **Run the application**
   ```bash
   go run cmd/compareflow/main.go
   ```

   The application will be available at http://localhost:8080

### Building for Production

1. **Build frontend**
   ```bash
   cd frontend
   npm install
   npm run build
   cp -r dist/* ../cmd/compareflow/web/dist/
   ```

2. **Build Go binary**
   ```bash
   go build -o compareflow cmd/compareflow/main.go
   ```

3. **Run with custom database**
   ```bash
   ./compareflow -db "postgresql://user:pass@host:5432/dbname?sslmode=require"
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
compareflow-go/
├── cmd/
│   └── compareflow/      # Main application entry point
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── models/          # Data models
│   └── services/        # Business logic
├── frontend/            # React application
│   ├── src/
│   │   ├── components/  # Reusable UI components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API client services
│   │   └── store/       # Redux store and slices
│   └── dist/           # Built frontend assets
├── scripts/            # Utility scripts
└── migrations/         # Database migrations
```

### Frontend Development

For frontend development with hot reload:

```bash
# Terminal 1: Run backend
go run cmd/compareflow/main.go

# Terminal 2: Run frontend dev server
cd frontend
npm run dev
```

The frontend dev server will proxy API requests to the backend.

### Adding New Features

1. **New API Endpoint**: Add handler in `internal/api/handlers/`
2. **New Model**: Add to `internal/models/`
3. **New Service**: Add to `internal/services/`
4. **Frontend Page**: Add to `frontend/src/pages/`
5. **Redux Slice**: Add to `frontend/src/store/slices/`

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
   ```

2. Verify connection string:
   ```bash
   psql "postgresql://compareflow:compareflow123@localhost:5432/compareflow?sslmode=disable"
   ```

3. Check logs:
   ```bash
   tail -f app.log
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

## Migration Status

- [x] Basic Go project structure
- [x] JWT authentication
- [x] User management
- [x] Connection CRUD operations
- [x] Validation CRUD operations
- [x] PostgreSQL integration with Podman
- [x] Local development scripts
- [x] Frontend embedding
- [x] Connection testing implementation (SQL Server)
- [x] Connection testing implementation (Databricks)
- [ ] Validation engine
- [ ] WebSocket support
- [ ] Data streaming
- [x] Proper password hashing (bcrypt)
- [x] React frontend with Material-UI
- [x] Redux state management

## Roadmap

- [ ] Implement data validation engine
- [ ] Add more database connectors (MySQL, PostgreSQL, Snowflake)
- [ ] WebSocket support for real-time updates
- [ ] Scheduling system for automated validations
- [ ] Email notifications
- [ ] Export results to CSV/Excel
- [ ] Data lineage visualization

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