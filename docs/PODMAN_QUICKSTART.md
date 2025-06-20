# CompareFlow with Podman Quick Start

## Prerequisites
- Podman installed
- Go 1.21+

## Start Database

Make the scripts executable first:
```bash
chmod +x start-podman.sh stop-podman.sh
```

Start PostgreSQL:
```bash
./start-podman.sh
```

## Run CompareFlow

In a new terminal:
```bash
# Run directly
go run cmd/compareflow/main.go

# Or build and run
go build -o compareflow cmd/compareflow/main.go
./compareflow
```

## Access Application

Open http://localhost:8080 in your browser
- Default login: `admin` / `admin123`

## Stop Database

```bash
./stop-podman.sh
```

## Troubleshooting

### Permission denied when running scripts
```bash
# Make executable
chmod +x start-podman.sh stop-podman.sh

# Or run with bash
bash start-podman.sh
```

### Port 5432 already in use
```bash
# Check what's using the port
sudo lsof -i :5432

# Or use a different port
podman run -d \
    --name compareflow-postgres \
    -e POSTGRES_USER=compareflow \
    -e POSTGRES_PASSWORD=compareflow123 \
    -e POSTGRES_DB=compareflow \
    -p 5433:5432 \
    postgres:15

# Then run CompareFlow with custom DB URL
go run cmd/compareflow/main.go -db "postgresql://compareflow:compareflow123@localhost:5433/compareflow?sslmode=disable"
```

### Container fails to start
```bash
# Check logs
podman logs compareflow-postgres

# Remove and retry
podman rm -f compareflow-postgres
./start-podman.sh
```

### Database connection refused
```bash
# Check if container is running
podman ps

# Check if PostgreSQL is ready
podman exec compareflow-postgres pg_isready -U compareflow

# Test connection directly
podman exec -it compareflow-postgres psql -U compareflow -d compareflow
```

## Manual Podman Commands

If the scripts don't work, run these commands manually:

1. **Start PostgreSQL**:
```bash
podman run -d \
    --name compareflow-postgres \
    -e POSTGRES_USER=compareflow \
    -e POSTGRES_PASSWORD=compareflow123 \
    -e POSTGRES_DB=compareflow \
    -p 5432:5432 \
    postgres:15
```

2. **Wait for it to be ready**:
```bash
# Check status
podman exec compareflow-postgres pg_isready -U compareflow
```

3. **Run CompareFlow**:
```bash
go run cmd/compareflow/main.go
```

4. **Stop when done**:
```bash
podman stop compareflow-postgres
podman rm compareflow-postgres
```