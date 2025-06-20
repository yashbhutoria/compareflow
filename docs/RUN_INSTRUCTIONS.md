# Running CompareFlow

## Prerequisites

1. **Go 1.21+** - [Install Go](https://golang.org/doc/install)
2. **Node.js 18+** - [Install Node.js](https://nodejs.org/)
3. **PostgreSQL** or **Podman/Docker** - For the database

## Quick Start

### Option 1: Using Podman/Docker (Recommended)

1. **Start the database**:
   ```bash
   # Using the provided script
   ./start-local.sh
   
   # Or manually with docker-compose
   docker-compose up -d postgres
   ```

2. **Build and run the application**:
   ```bash
   # Install dependencies
   go mod tidy
   
   # Build frontend (if needed)
   cd frontend
   npm install
   npm run build
   cd ..
   
   # Copy frontend to embedded location
   rm -rf cmd/compareflow/web/dist
   mkdir -p cmd/compareflow/web/dist
   cp -r frontend/dist/* cmd/compareflow/web/dist/
   
   # Run the application
   go run cmd/compareflow/main.go
   ```

3. **Access the application**:
   - Open http://localhost:8080 in your browser
   - Default admin credentials: `admin` / `admin123`

### Option 2: Using External PostgreSQL

1. **Set database URL**:
   ```bash
   export DATABASE_URL="postgresql://user:password@host:5432/compareflow?sslmode=disable"
   ```

2. **Run the application**:
   ```bash
   go run cmd/compareflow/main.go -db "$DATABASE_URL"
   ```

## Build Commands

### Using Make

```bash
# Install dependencies
make deps

# Build frontend
make frontend

# Build backend
make build

# Build everything
make build-all

# Run in development mode
make run

# Run tests
make test
```

### Manual Build

```bash
# Clean up old files (if any)
rm -f internal/connectors/config.go
rm -f internal/connectors/connector.go
rm -f internal/connectors/sqlserver.go
rm -f internal/connectors/databricks.go
rm -f internal/connectors/postgresql.go

# Build frontend
cd frontend
npm install
npm run build
cd ..

# Copy frontend files
rm -rf cmd/compareflow/web/dist
mkdir -p cmd/compareflow/web/dist
cp -r frontend/dist/* cmd/compareflow/web/dist/

# Build Go binary
go build -o bin/compareflow cmd/compareflow/main.go

# Run the binary
./bin/compareflow
```

## Command Line Options

```bash
# Run with custom port
go run cmd/compareflow/main.go -port 3000

# Run with custom database URL
go run cmd/compareflow/main.go -db "postgresql://user:pass@host:5432/db"

# Run in debug mode
go run cmd/compareflow/main.go -mode debug
```

## Troubleshooting

### Build Errors

1. **Missing embed files**:
   ```bash
   # Ensure web/dist exists
   mkdir -p cmd/compareflow/web/dist
   echo "<html><body>CompareFlow</body></html>" > cmd/compareflow/web/dist/index.html
   ```

2. **Module errors**:
   ```bash
   # Clean module cache
   go clean -modcache
   go mod download
   go mod tidy
   ```

3. **Database connection errors**:
   ```bash
   # Check PostgreSQL is running
   docker ps | grep postgres
   
   # Test connection
   psql -h localhost -U compareflow -d compareflow
   ```

### Runtime Errors

1. **Port already in use**:
   ```bash
   # Find process using port 8080
   lsof -i :8080
   
   # Run on different port
   go run cmd/compareflow/main.go -port 8081
   ```

2. **Database migrations fail**:
   ```bash
   # Check database exists
   psql -h localhost -U postgres -c "CREATE DATABASE compareflow;"
   
   # Check user permissions
   psql -h localhost -U postgres -c "GRANT ALL ON DATABASE compareflow TO compareflow;"
   ```

## Development Workflow

1. **Start database**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Run backend with hot reload**:
   ```bash
   # Install air
   go install github.com/cosmtrek/air@latest
   
   # Run with hot reload
   air
   ```

3. **Run frontend dev server**:
   ```bash
   cd frontend
   npm run dev
   ```

4. **Access development servers**:
   - Backend API: http://localhost:8080
   - Frontend Dev: http://localhost:5173

## Production Deployment

1. **Build release binary**:
   ```bash
   # Linux
   CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o compareflow-linux-amd64 cmd/compareflow/main.go
   
   # macOS Intel
   CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o compareflow-darwin-amd64 cmd/compareflow/main.go
   
   # macOS Apple Silicon
   CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o compareflow-darwin-arm64 cmd/compareflow/main.go
   ```

2. **Set production environment variables**:
   ```bash
   export DATABASE_URL="postgresql://user:pass@host:5432/compareflow?sslmode=require"
   export JWT_SECRET="your-secure-secret-key"
   export GIN_MODE="release"
   ```

3. **Run the binary**:
   ```bash
   ./compareflow-linux-amd64 -port 8080
   ```

## Testing Connections

After starting the application, you can test database connections:

1. **Login** at http://localhost:8080
2. **Create a connection**:
   - SQL Server example:
     ```json
     {
       "name": "Test SQL Server",
       "type": "sqlserver",
       "config": {
         "server": "localhost",
         "port": 1433,
         "database": "testdb",
         "username": "sa",
         "password": "YourPassword123!",
         "encrypt": false,
         "trust_server_certificate": true
       }
     }
     ```
   
   - Databricks example:
     ```json
     {
       "name": "Test Databricks",
       "type": "databricks",
       "config": {
         "workspace": "https://your-workspace.databricks.com",
         "http_path": "/sql/1.0/endpoints/your-endpoint",
         "access_token": "dapi..."
       }
     }
     ```

3. **Test the connection** using the Test button

## Notes

- The application embeds the frontend files, so you need to rebuild when frontend changes
- Database migrations run automatically on startup
- All passwords are encrypted in the database
- JWT tokens expire after 7 days by default