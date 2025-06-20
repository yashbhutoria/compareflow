# CompareFlow Deployment Guide

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Local Development](#local-development)
3. [Production Deployment](#production-deployment)
4. [Configuration](#configuration)
5. [Security Hardening](#security-hardening)
6. [Monitoring and Maintenance](#monitoring-and-maintenance)
7. [Troubleshooting](#troubleshooting)
8. [Backup and Recovery](#backup-and-recovery)

## 1. Prerequisites

### System Requirements
- **OS**: Linux (Ubuntu 20.04+, RHEL 8+), macOS, Windows (WSL2)
- **CPU**: 2+ cores recommended
- **RAM**: 4GB minimum, 8GB recommended
- **Disk**: 10GB free space
- **Network**: Outbound HTTPS access

### Software Dependencies
- **Go**: 1.21+ (for building from source)
- **Node.js**: 18+ (for frontend development)
- **PostgreSQL**: 13+ (or Podman/Docker for containerized deployment)
- **Git**: For source code management

## 2. Local Development

### 2.1 Quick Start with Podman

```bash
# Clone the repository
git clone https://github.com/compareflow/compareflow.git
cd compareflow

# Start PostgreSQL and build the application
./start-local.sh

# Access the application
open http://localhost:8080

# Default credentials
# Username: admin
# Password: admin123
```

### 2.2 Manual Development Setup

#### Step 1: Database Setup
```bash
# Using Podman
podman run -d \
  --name compareflow-postgres \
  -e POSTGRES_USER=compareflow \
  -e POSTGRES_PASSWORD=compareflow123 \
  -e POSTGRES_DB=compareflow \
  -p 5432:5432 \
  postgres:15-alpine

# Or using existing PostgreSQL
createdb compareflow
createuser compareflow
```

#### Step 2: Environment Configuration
```bash
# Create .env file
cat > .env << EOF
DATABASE_URL=postgresql://compareflow:compareflow123@localhost:5432/compareflow?sslmode=disable
JWT_SECRET=your-development-secret-key-change-in-production
PORT=8080
GIN_MODE=debug
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
EOF
```

#### Step 3: Build Frontend
```bash
cd frontend
npm install
npm run build
cp -r dist/* ../cmd/compareflow/web/dist/
cd ..
```

#### Step 4: Run the Application
```bash
# Run with hot reload (requires Air)
air

# Or run directly
go run cmd/compareflow/main.go
```

### 2.3 Development with Frontend Hot Reload

```bash
# Terminal 1: Backend
go run cmd/compareflow/main.go

# Terminal 2: Frontend dev server
cd frontend
npm run dev
# Access at http://localhost:5173
```

## 3. Production Deployment

### 3.1 Building for Production

```bash
# Build optimized binary
make build-prod

# Or manually
cd frontend
npm ci --production
npm run build
cd ..
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o compareflow \
  cmd/compareflow/main.go
```

### 3.2 Deployment Options

#### Option 1: Systemd Service (Recommended for VPS/Bare Metal)

```bash
# Copy binary to server
scp compareflow user@server:/opt/compareflow/

# Create systemd service
sudo tee /etc/systemd/system/compareflow.service << EOF
[Unit]
Description=CompareFlow Data Validation Service
After=network.target postgresql.service

[Service]
Type=simple
User=compareflow
Group=compareflow
WorkingDirectory=/opt/compareflow
ExecStart=/opt/compareflow/compareflow
Restart=always
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/compareflow/data

# Environment
Environment="GIN_MODE=release"
Environment="DATABASE_URL=postgresql://compareflow:password@localhost/compareflow?sslmode=require"
Environment="JWT_SECRET=your-production-secret-key"
Environment="ALLOWED_ORIGINS=https://compareflow.yourdomain.com"

[Install]
WantedBy=multi-user.target
EOF

# Create user and directories
sudo useradd -r -s /bin/false compareflow
sudo mkdir -p /opt/compareflow/data
sudo chown -R compareflow:compareflow /opt/compareflow

# Start service
sudo systemctl daemon-reload
sudo systemctl enable compareflow
sudo systemctl start compareflow
```

#### Option 2: Docker Deployment

```bash
# Build Docker image
docker build -t compareflow:latest .

# Run with Docker Compose
docker-compose -f docker-compose.prod.yml up -d
```

**docker-compose.prod.yml:**
```yaml
version: '3.8'

services:
  compareflow:
    image: compareflow:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://compareflow:password@postgres:5432/compareflow?sslmode=disable
      - JWT_SECRET=${JWT_SECRET}
      - GIN_MODE=release
    depends_on:
      - postgres
    restart: always

  postgres:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=compareflow
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=compareflow
    restart: always

volumes:
  postgres_data:
```

#### Option 3: Kubernetes Deployment

```yaml
# compareflow-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: compareflow
spec:
  replicas: 3
  selector:
    matchLabels:
      app: compareflow
  template:
    metadata:
      labels:
        app: compareflow
    spec:
      containers:
      - name: compareflow
        image: compareflow:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: compareflow-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: compareflow-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: compareflow-service
spec:
  selector:
    app: compareflow
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

### 3.3 Reverse Proxy Configuration

#### Nginx Configuration
```nginx
server {
    listen 80;
    server_name compareflow.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name compareflow.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/compareflow.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/compareflow.yourdomain.com/privkey.pem;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self' https:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';" always;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts for long-running validations
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }

    # WebSocket support (future)
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

#### Apache Configuration
```apache
<VirtualHost *:80>
    ServerName compareflow.yourdomain.com
    Redirect permanent / https://compareflow.yourdomain.com/
</VirtualHost>

<VirtualHost *:443>
    ServerName compareflow.yourdomain.com
    
    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/compareflow.yourdomain.com/cert.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/compareflow.yourdomain.com/privkey.pem
    
    ProxyPreserveHost On
    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/
    
    # WebSocket support
    RewriteEngine on
    RewriteCond %{HTTP:Upgrade} websocket [NC]
    RewriteCond %{HTTP:Connection} upgrade [NC]
    RewriteRule ^/?(.*) "ws://localhost:8080/$1" [P,L]
</VirtualHost>
```

## 4. Configuration

### 4.1 Environment Variables

```bash
# Required
DATABASE_URL=postgresql://user:pass@host:port/db?sslmode=require
JWT_SECRET=your-very-long-random-secret-key

# Optional
PORT=8080                                    # Server port
GIN_MODE=release                            # Gin mode (debug/release)
ALLOWED_ORIGINS=https://app.domain.com      # CORS origins
LOG_LEVEL=info                              # Log level
MAX_CONNECTIONS=100                         # DB connection pool size
JWT_EXPIRATION_HOURS=168                    # Token expiration (7 days)
ENCRYPTION_KEY=32-byte-hex-key              # For encrypting sensitive data
```

### 4.2 Configuration File (config.yaml)

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

database:
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 30m
  log_level: error

auth:
  jwt_expiration: 168h
  bcrypt_cost: 10
  session_timeout: 24h

validation:
  max_concurrent: 10
  default_timeout: 300s
  max_result_size: 100MB
  stream_batch_size: 1000

features:
  enable_scheduling: false
  enable_notifications: false
  enable_api_keys: false
```

### 4.3 Database Configuration

```sql
-- Performance tuning for PostgreSQL
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

-- Create indexes for performance
CREATE INDEX idx_validations_user_status ON validations(user_id, status);
CREATE INDEX idx_validations_created_at ON validations(created_at DESC);
CREATE INDEX idx_connections_user_id ON connections(user_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
```

## 5. Security Hardening

### 5.1 SSL/TLS Configuration

```bash
# Generate self-signed certificate (development only)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout compareflow.key -out compareflow.crt

# Or use Let's Encrypt (production)
certbot certonly --standalone -d compareflow.yourdomain.com
```

### 5.2 Firewall Rules

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP (redirect to HTTPS)
sudo ufw allow 443/tcp     # HTTPS
sudo ufw enable

# firewalld (RHEL/CentOS)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 5.3 Security Checklist

- [ ] Change default passwords
- [ ] Use strong JWT secret (min 32 characters)
- [ ] Enable SSL/TLS
- [ ] Configure firewall
- [ ] Disable debug mode
- [ ] Set secure CORS origins
- [ ] Enable rate limiting
- [ ] Regular security updates
- [ ] Implement backup strategy
- [ ] Monitor logs for suspicious activity

### 5.4 Rate Limiting

```go
// Add to main.go or middleware
import "github.com/ulule/limiter/v3"

// Configure rate limiter
rate := limiter.Rate{
    Period: 1 * time.Minute,
    Limit:  60,
}
store := memory.NewStore()
instance := limiter.New(store, rate)

// Apply to routes
router.Use(gin_limiter.Middleware(instance))
```

## 6. Monitoring and Maintenance

### 6.1 Health Checks

```bash
# Basic health check
curl https://compareflow.yourdomain.com/health

# Detailed health check (with auth)
curl -H "Authorization: Bearer $TOKEN" \
  https://compareflow.yourdomain.com/api/v1/system/health
```

### 6.2 Logging

```bash
# View systemd logs
sudo journalctl -u compareflow -f

# Configure log rotation
sudo tee /etc/logrotate.d/compareflow << EOF
/var/log/compareflow/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 compareflow compareflow
    sharedscripts
    postrotate
        systemctl reload compareflow
    endscript
}
EOF
```

### 6.3 Monitoring with Prometheus

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'compareflow'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### 6.4 Performance Monitoring

```sql
-- Database query performance
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
WHERE query LIKE '%validations%'
ORDER BY total_time DESC
LIMIT 20;

-- Connection pool stats
SELECT 
    datname,
    numbackends,
    xact_commit,
    xact_rollback,
    blks_hit,
    blks_read
FROM pg_stat_database
WHERE datname = 'compareflow';
```

## 7. Troubleshooting

### 7.1 Common Issues

#### Application Won't Start
```bash
# Check if port is in use
sudo lsof -i :8080

# Check systemd status
sudo systemctl status compareflow
sudo journalctl -u compareflow --no-pager | tail -50

# Verify database connection
psql "postgresql://compareflow:password@localhost/compareflow?sslmode=disable" -c "SELECT 1"
```

#### Database Connection Issues
```bash
# Test connection
pg_isready -h localhost -p 5432

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Verify user permissions
psql -U postgres -c "\du compareflow"
```

#### High Memory Usage
```bash
# Check memory usage
ps aux | grep compareflow
htop

# Analyze Go memory profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

### 7.2 Debug Mode

```bash
# Enable debug logging
export GIN_MODE=debug
export LOG_LEVEL=debug

# Run with verbose output
./compareflow -v

# Enable pprof endpoints (development only)
./compareflow -enable-pprof
```

## 8. Backup and Recovery

### 8.1 Database Backup

```bash
# Manual backup
pg_dump -U compareflow -d compareflow -f compareflow_backup_$(date +%Y%m%d).sql

# Automated daily backup script
cat > /opt/compareflow/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/compareflow/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_NAME="compareflow"
DB_USER="compareflow"

mkdir -p $BACKUP_DIR
pg_dump -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/backup_$TIMESTAMP.sql.gz

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
EOF

chmod +x /opt/compareflow/backup.sh

# Add to crontab
echo "0 2 * * * /opt/compareflow/backup.sh" | crontab -
```

### 8.2 Restore Procedure

```bash
# Stop application
sudo systemctl stop compareflow

# Restore database
gunzip -c backup_20240115.sql.gz | psql -U compareflow -d compareflow

# Start application
sudo systemctl start compareflow
```

### 8.3 Configuration Backup

```bash
# Backup critical files
tar -czf compareflow_config_$(date +%Y%m%d).tar.gz \
  /etc/systemd/system/compareflow.service \
  /opt/compareflow/.env \
  /opt/compareflow/config.yaml \
  /etc/nginx/sites-available/compareflow
```

## 9. Scaling Considerations

### 9.1 Vertical Scaling
- Increase CPU cores for parallel validations
- Add RAM for larger dataset comparisons
- Use SSD storage for better I/O performance

### 9.2 Horizontal Scaling
- Deploy multiple instances behind load balancer
- Use shared PostgreSQL cluster
- Implement distributed locking for validations
- Consider message queue for job distribution

### 9.3 Database Scaling
- Read replicas for reporting
- Connection pooling with PgBouncer
- Partitioning for large tables
- Archive old validation results

## 10. Upgrade Procedure

```bash
# 1. Backup current installation
./backup.sh

# 2. Download new version
wget https://github.com/compareflow/compareflow/releases/latest/download/compareflow

# 3. Stop service
sudo systemctl stop compareflow

# 4. Replace binary
sudo cp compareflow /opt/compareflow/compareflow
sudo chmod +x /opt/compareflow/compareflow
sudo chown compareflow:compareflow /opt/compareflow/compareflow

# 5. Run migrations (if any)
/opt/compareflow/compareflow migrate

# 6. Start service
sudo systemctl start compareflow

# 7. Verify
curl https://compareflow.yourdomain.com/health
```

## Conclusion

This deployment guide covers the essential aspects of deploying CompareFlow in various environments. Always test thoroughly in a staging environment before deploying to production, and maintain regular backups of both the database and configuration files.