# CompareFlow Technical Design Document

## 1. Introduction

This document provides the detailed technical design for CompareFlow, a single-binary data validation platform. It covers the implementation details, algorithms, data structures, and technical decisions.

## 2. Technology Stack

### 2.1 Backend Technologies
- **Language**: Go 1.21+
- **Web Framework**: Gin v1.9+
- **ORM**: GORM v1.25+
- **Database**: PostgreSQL 13+ / SQLite 3.36+
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **SQL Server Driver**: github.com/denisenkom/go-mssqldb
- **Databricks Driver**: Custom REST client

### 2.2 Frontend Technologies
- **Framework**: React 18.2+
- **Language**: TypeScript 5.0+
- **UI Library**: Material-UI v5
- **State Management**: Redux Toolkit
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **Build Tool**: Vite 4.4+

### 2.3 Development Tools
- **Version Control**: Git
- **Container Runtime**: Podman/Docker
- **Testing**: Go testing package, Jest, React Testing Library
- **Linting**: golangci-lint, ESLint
- **Formatting**: gofmt, Prettier

## 3. System Design

### 3.1 Application Structure

```
compareflow/
├── cmd/
│   └── compareflow/          # Application entry point
│       ├── main.go          # Main function
│       └── web/             # Embedded frontend assets
│           └── dist/        # Built React app
├── internal/                # Private application code
│   ├── api/                # HTTP layer
│   │   ├── handlers/       # Request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   └── routes.go       # Route definitions
│   ├── config/             # Configuration
│   ├── database/           # Database setup
│   ├── models/             # Domain models
│   ├── services/           # Business logic
│   └── validation/         # Validation engine
├── frontend/               # React application
│   ├── src/
│   │   ├── components/    # Reusable components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API clients
│   │   ├── store/        # Redux store
│   │   └── types/        # TypeScript types
│   └── dist/             # Built assets
├── migrations/            # Database migrations
├── scripts/              # Utility scripts
└── tests/               # Integration tests
```

### 3.2 Data Flow Architecture

```
User Request → Gin Router → Middleware → Handler → Service → Database
                                ↓                      ↓
                            Validation ←─────────────→ External DBs
                                ↓
                            Response ← Result Processing
```

### 3.3 Database Design

#### 3.3.1 Entity Relationship Diagram

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   users     │     │ connections  │     │ validations  │
├─────────────┤     ├──────────────┤     ├──────────────┤
│ id          │←────┤ user_id      │     │ id           │
│ username    │     │ id           │←┬───┤ source_conn  │
│ email       │     │ name         │ └───┤ target_conn  │
│ password    │     │ type         │     │ user_id      │──┐
│ created_at  │     │ config       │     │ config       │  │
│ updated_at  │     │ created_at   │     │ status       │  │
└─────────────┘     │ updated_at   │     │ results      │  │
                    └──────────────┘     │ created_at   │  │
                                        │ updated_at   │  │
                                        └──────────────┘  │
                                                          │
                    ┌──────────────┐                     │
                    │ audit_logs   │                     │
                    ├──────────────┤                     │
                    │ id           │                     │
                    │ user_id      │←────────────────────┘
                    │ action       │
                    │ entity_type  │
                    │ entity_id    │
                    │ timestamp    │
                    └──────────────┘
```

#### 3.3.2 JSON Column Schemas

**Connection Config (JSONB):**
```json
{
  // SQL Server
  "server": "string",
  "port": "number",
  "database": "string",
  "username": "string",
  "password": "encrypted_string",
  "encrypt": "boolean",
  "trust_server_certificate": "boolean",
  
  // Databricks
  "workspace": "string",
  "http_path": "string",
  "access_token": "encrypted_string"
}
```

**Validation Config (JSONB):**
```json
{
  "comparison_type": "row_count|data_match|schema",
  "source_query": "string",
  "target_query": "string",
  "key_columns": ["string"],
  "comparison_options": {
    "check_row_count": "boolean",
    "check_column_count": "boolean",
    "check_data_types": "boolean",
    "check_nulls": "boolean",
    "case_sensitive": "boolean",
    "decimal_precision": "number",
    "date_precision": "string"
  },
  "error_margin": {
    "type": "absolute|percentage",
    "value": "number"
  },
  "performance": {
    "batch_size": "number",
    "timeout_seconds": "number",
    "max_differences": "number"
  }
}
```

**Validation Results (JSONB):**
```json
{
  "execution_id": "uuid",
  "start_time": "timestamp",
  "end_time": "timestamp",
  "duration_ms": "number",
  "status": "success|failure|error",
  "summary": {
    "source_row_count": "number",
    "target_row_count": "number",
    "matched_rows": "number",
    "mismatched_rows": "number",
    "missing_in_target": "number",
    "extra_in_target": "number",
    "success_rate": "number"
  },
  "details": {
    "differences": [{
      "key": "object",
      "type": "missing|extra|mismatch",
      "source_data": "object",
      "target_data": "object",
      "columns": ["string"]
    }],
    "column_stats": {
      "column_name": {
        "source_min": "any",
        "source_max": "any",
        "source_avg": "number",
        "target_min": "any",
        "target_max": "any",
        "target_avg": "number"
      }
    }
  },
  "errors": [{
    "timestamp": "string",
    "message": "string",
    "details": "string"
  }]
}
```

### 3.4 API Design

#### 3.4.1 Request/Response Formats

**Standard Response Envelope:**
```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

type MetaInfo struct {
    Page       int `json:"page,omitempty"`
    PerPage    int `json:"per_page,omitempty"`
    TotalCount int `json:"total_count,omitempty"`
}
```

#### 3.4.2 Authentication Flow

```go
// JWT Claims Structure
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// Token Generation
func GenerateToken(user *models.User) (string, error) {
    expirationTime := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "compareflow",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.JWTSecret))
}
```

### 3.5 Validation Engine Design

#### 3.5.1 Core Interfaces

```go
// Validator interface for different validation types
type Validator interface {
    Validate(ctx context.Context, config ValidationConfig) (*ValidationResult, error)
    GetType() string
}

// Connection interface for database operations
type DBConnection interface {
    Connect(ctx context.Context) error
    ExecuteQuery(ctx context.Context, query string) (*QueryResult, error)
    GetSchema(ctx context.Context, table string) (*TableSchema, error)
    Close() error
}

// Result processor interface
type ResultProcessor interface {
    Process(source, target *QueryResult) (*ComparisonResult, error)
}
```

#### 3.5.2 Validation Algorithms

**Row Count Validation:**
```go
func (v *RowCountValidator) Validate(ctx context.Context, config ValidationConfig) (*ValidationResult, error) {
    // 1. Execute count queries
    sourceCount := v.executeCountQuery(ctx, config.SourceConn, config.SourceQuery)
    targetCount := v.executeCountQuery(ctx, config.TargetConn, config.TargetQuery)
    
    // 2. Apply error margin
    difference := abs(sourceCount - targetCount)
    var passed bool
    
    if config.ErrorMargin.Type == "percentage" {
        threshold := sourceCount * config.ErrorMargin.Value / 100
        passed = difference <= threshold
    } else {
        passed = difference <= config.ErrorMargin.Value
    }
    
    // 3. Build result
    return &ValidationResult{
        Status: determineStatus(passed),
        Summary: map[string]interface{}{
            "source_count": sourceCount,
            "target_count": targetCount,
            "difference":   difference,
            "passed":       passed,
        },
    }, nil
}
```

**Data Match Validation:**
```go
func (v *DataMatchValidator) Validate(ctx context.Context, config ValidationConfig) (*ValidationResult, error) {
    // 1. Stream data from both sources
    sourceChan := v.streamData(ctx, config.SourceConn, config.SourceQuery)
    targetChan := v.streamData(ctx, config.TargetConn, config.TargetQuery)
    
    // 2. Build hash maps using key columns
    sourceMap := make(map[string]Row)
    targetMap := make(map[string]Row)
    
    // 3. Compare data
    differences := []Difference{}
    
    // Check for missing/mismatched rows
    for key, sourceRow := range sourceMap {
        if targetRow, exists := targetMap[key]; exists {
            if diff := v.compareRows(sourceRow, targetRow, config); diff != nil {
                differences = append(differences, diff)
            }
        } else {
            differences = append(differences, Difference{
                Type: "missing_in_target",
                Key:  key,
                SourceRow: sourceRow,
            })
        }
    }
    
    // Check for extra rows in target
    for key, targetRow := range targetMap {
        if _, exists := sourceMap[key]; !exists {
            differences = append(differences, Difference{
                Type: "extra_in_target",
                Key:  key,
                TargetRow: targetRow,
            })
        }
    }
    
    // 4. Calculate summary statistics
    return v.buildResult(len(sourceMap), len(targetMap), differences), nil
}
```

#### 3.5.3 Streaming and Memory Management

```go
type StreamProcessor struct {
    batchSize int
    maxMemory int64
}

func (sp *StreamProcessor) ProcessStream(ctx context.Context, query string) <-chan Batch {
    results := make(chan Batch, 10)
    
    go func() {
        defer close(results)
        
        rows, err := sp.conn.Query(ctx, query)
        if err != nil {
            results <- Batch{Error: err}
            return
        }
        defer rows.Close()
        
        batch := make([]Row, 0, sp.batchSize)
        for rows.Next() {
            row := sp.scanRow(rows)
            batch = append(batch, row)
            
            if len(batch) >= sp.batchSize {
                select {
                case results <- Batch{Rows: batch}:
                    batch = make([]Row, 0, sp.batchSize)
                case <-ctx.Done():
                    return
                }
            }
        }
        
        if len(batch) > 0 {
            results <- Batch{Rows: batch}
        }
    }()
    
    return results
}
```

### 3.6 Connection Pool Management

```go
type ConnectionPool struct {
    mu          sync.RWMutex
    connections map[string]*DBConnection
    maxSize     int
    timeout     time.Duration
}

func (cp *ConnectionPool) GetConnection(connID string) (*DBConnection, error) {
    cp.mu.RLock()
    if conn, exists := cp.connections[connID]; exists {
        cp.mu.RUnlock()
        return conn, nil
    }
    cp.mu.RUnlock()
    
    // Create new connection
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    // Double-check after acquiring write lock
    if conn, exists := cp.connections[connID]; exists {
        return conn, nil
    }
    
    // Load connection config and create
    conn, err := cp.createConnection(connID)
    if err != nil {
        return nil, err
    }
    
    cp.connections[connID] = conn
    return conn, nil
}
```

### 3.7 Security Implementation

#### 3.7.1 Encryption Service

```go
type EncryptionService struct {
    key []byte
}

func (es *EncryptionService) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(es.key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (es *EncryptionService) Decrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(es.key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}
```

#### 3.7.2 SQL Injection Prevention

```go
// Safe query builder
type QueryBuilder struct {
    query  strings.Builder
    params []interface{}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
    qb.query.WriteString("SELECT ")
    for i, col := range columns {
        if i > 0 {
            qb.query.WriteString(", ")
        }
        qb.query.WriteString(qb.sanitizeIdentifier(col))
    }
    return qb
}

func (qb *QueryBuilder) Where(condition string, params ...interface{}) *QueryBuilder {
    qb.query.WriteString(" WHERE ")
    qb.query.WriteString(condition)
    qb.params = append(qb.params, params...)
    return qb
}

func (qb *QueryBuilder) sanitizeIdentifier(identifier string) string {
    // Remove any characters that aren't alphanumeric or underscore
    return regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(identifier, "")
}
```

### 3.8 Performance Optimizations

#### 3.8.1 Query Optimization

```go
// Parallel query execution
func ExecuteQueriesParallel(ctx context.Context, queries []Query) ([]Result, error) {
    results := make([]Result, len(queries))
    errChan := make(chan error, len(queries))
    
    var wg sync.WaitGroup
    for i, query := range queries {
        wg.Add(1)
        go func(idx int, q Query) {
            defer wg.Done()
            
            result, err := q.Execute(ctx)
            if err != nil {
                errChan <- err
                return
            }
            results[idx] = result
        }(i, query)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }
    
    return results, nil
}
```

#### 3.8.2 Caching Strategy

```go
type Cache struct {
    data map[string]CacheEntry
    mu   sync.RWMutex
    ttl  time.Duration
}

type CacheEntry struct {
    Value      interface{}
    Expiration time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, exists := c.data[key]
    if !exists || time.Now().After(entry.Expiration) {
        return nil, false
    }
    
    return entry.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = CacheEntry{
        Value:      value,
        Expiration: time.Now().Add(c.ttl),
    }
}
```

### 3.9 Error Handling

```go
// Custom error types
type ValidationError struct {
    Code    string
    Message string
    Details interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error codes
const (
    ErrConnectionFailed = "CONNECTION_FAILED"
    ErrQueryTimeout     = "QUERY_TIMEOUT"
    ErrInvalidConfig    = "INVALID_CONFIG"
    ErrUnauthorized     = "UNAUTHORIZED"
)

// Error handler middleware
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            var status int
            var response Response
            
            switch e := err.Err.(type) {
            case ValidationError:
                status = http.StatusBadRequest
                response = Response{
                    Success: false,
                    Error: &ErrorInfo{
                        Code:    e.Code,
                        Message: e.Message,
                        Details: fmt.Sprintf("%v", e.Details),
                    },
                }
            default:
                status = http.StatusInternalServerError
                response = Response{
                    Success: false,
                    Error: &ErrorInfo{
                        Code:    "INTERNAL_ERROR",
                        Message: "An unexpected error occurred",
                    },
                }
            }
            
            c.JSON(status, response)
        }
    }
}
```

### 3.10 Testing Strategy

#### 3.10.1 Unit Tests

```go
func TestPasswordHashing(t *testing.T) {
    user := &User{Username: "testuser"}
    password := "testpass123"
    
    err := user.SetPassword(password)
    assert.NoError(t, err)
    assert.NotEmpty(t, user.PasswordHash)
    
    valid := user.CheckPassword(password)
    assert.True(t, valid)
    
    invalid := user.CheckPassword("wrongpass")
    assert.False(t, invalid)
}
```

#### 3.10.2 Integration Tests

```go
func TestValidationEndToEnd(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    // Create test connections
    sourceConn := createTestConnection(t, db, "source")
    targetConn := createTestConnection(t, db, "target")
    
    // Create validation
    validation := &Validation{
        Name:               "Test Validation",
        SourceConnectionID: sourceConn.ID,
        TargetConnectionID: targetConn.ID,
        Config: ValidationConfig{
            ComparisonType: "row_count",
            SourceQuery:    "SELECT COUNT(*) FROM users",
            TargetQuery:    "SELECT COUNT(*) FROM users",
        },
    }
    
    // Run validation
    engine := NewValidationEngine(db)
    result, err := engine.RunValidation(context.Background(), validation)
    
    assert.NoError(t, err)
    assert.Equal(t, "success", result.Status)
}
```

### 3.11 Monitoring and Logging

```go
// Structured logging
type Logger struct {
    *zap.Logger
}

func (l *Logger) LogValidationStart(validationID uint, userID uint) {
    l.Info("Validation started",
        zap.Uint("validation_id", validationID),
        zap.Uint("user_id", userID),
        zap.Time("timestamp", time.Now()),
    )
}

// Metrics collection
type Metrics struct {
    ValidationDuration  prometheus.Histogram
    ValidationCount     prometheus.Counter
    ConnectionPoolSize  prometheus.Gauge
    ActiveValidations   prometheus.Gauge
}

func (m *Metrics) RecordValidation(duration time.Duration, status string) {
    m.ValidationDuration.Observe(duration.Seconds())
    m.ValidationCount.With(prometheus.Labels{"status": status}).Inc()
}
```

## 4. Deployment Architecture

### 4.1 Build Process

```dockerfile
# Multi-stage build
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./cmd/compareflow/web/dist
RUN go build -o compareflow cmd/compareflow/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/compareflow .
EXPOSE 8080
CMD ["./compareflow"]
```

### 4.2 Configuration Management

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Auth     AuthConfig     `mapstructure:"auth"`
    Features FeatureFlags   `mapstructure:"features"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/compareflow/")
    
    // Environment variable overrides
    viper.SetEnvPrefix("COMPAREFLOW")
    viper.AutomaticEnv()
    
    // Defaults
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("auth.jwt_expiration", "168h")
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 5. Future Technical Considerations

### 5.1 Horizontal Scaling
- Implement distributed locking for validations
- Use message queue for job distribution
- Shared cache layer (Redis)

### 5.2 Performance Improvements
- Query result streaming
- Columnar data processing
- GPU acceleration for comparisons

### 5.3 Advanced Features
- Plugin architecture for custom validators
- Webhook system for integrations
- GraphQL API option

## 6. Conclusion

This technical design provides a robust foundation for CompareFlow, balancing simplicity with power, and maintaining flexibility for future enhancements while delivering immediate value through the single-binary deployment model.