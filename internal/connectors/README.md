# Database Connectors

This package provides a pluggable architecture for adding new database connectors to CompareFlow.

## Architecture

The connector system uses:
- **Package-based modularity**: Each connector is in its own package
- **Self-registration**: Connectors register themselves via `init()` functions
- **Interface-based design**: All connectors implement the `Connector` interface
- **Strongly-typed configurations**: Each connector has its own config struct

## Benefits

1. **Complete Isolation**: Each connector is completely independent
2. **Easy Testing**: Test each connector in isolation
3. **No Central Config File**: Each connector manages its own configuration
4. **Parallel Development**: Multiple developers can work on different connectors
5. **Clean Dependencies**: Each connector only imports what it needs

## Project Structure

```
connectors/
├── types.go           # Core interfaces and registry
├── connectors.go      # Imports all connector packages
├── README.md          # This file
├── sqlserver/         # SQL Server connector package
│   ├── sqlserver.go
│   └── sqlserver_test.go
├── databricks/        # Databricks connector package
│   ├── databricks.go
│   └── databricks_test.go
└── postgresql/        # PostgreSQL connector package
    ├── postgresql.go
    └── postgresql_test.go
```

## Adding a New Connector

### 1. Create a new package

Create a new directory under `connectors/`:

```bash
mkdir internal/connectors/mysql
```

### 2. Implement the connector

Create `mysql/mysql.go`:

```go
package mysql

import (
    "database/sql"
    "fmt"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/compareflow/compareflow/internal/connectors"
)

func init() {
    // Auto-register when package is imported
    connectors.Register(New())
}

// Config represents MySQL connection configuration
type Config struct {
    Host     string `json:"host" binding:"required"`
    Port     int    `json:"port"`
    Database string `json:"database" binding:"required"`
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Charset  string `json:"charset"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("host is required")
    }
    // ... other validations
    
    // Set defaults
    if c.Port == 0 {
        c.Port = 3306
    }
    if c.Charset == "" {
        c.Charset = "utf8mb4"
    }
    return nil
}

// Connector implements the Connector interface for MySQL
type Connector struct{}

// New creates a new MySQL connector
func New() *Connector {
    return &Connector{}
}

// Type returns the connector type
func (c *Connector) Type() string {
    return "mysql"
}

// ParseConfig parses raw config into typed config
func (c *Connector) ParseConfig(configMap map[string]interface{}) (interface{}, error) {
    var config Config
    if err := connectors.ParseJSON(configMap, &config); err != nil {
        return nil, fmt.Errorf("failed to parse MySQL config: %w", err)
    }
    
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return &config, nil
}

// Connect establishes a connection to MySQL
func (c *Connector) Connect(config interface{}) (*sql.DB, error) {
    myConfig, ok := config.(*Config)
    if !ok {
        return nil, fmt.Errorf("invalid config type: expected *Config, got %T", config)
    }
    
    connString := c.buildConnectionString(myConfig)
    return sql.Open("mysql", connString)
}

// TestConnection verifies the MySQL connection works
func (c *Connector) TestConnection(config interface{}) error {
    // Implementation
}

// GetTables returns a list of tables
func (c *Connector) GetTables(db *sql.DB) ([]string, error) {
    // Implementation
}

// GetColumns returns column information
func (c *Connector) GetColumns(db *sql.DB, tableName string) ([]connectors.ColumnInfo, error) {
    // Implementation
}

// buildConnectionString builds a MySQL connection string
func (c *Connector) buildConnectionString(config *Config) string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
        config.Username, config.Password, config.Host, 
        config.Port, config.Database, config.Charset)
}
```

### 3. Add tests

Create `mysql/mysql_test.go`:

```go
package mysql

import (
    "testing"
)

func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: Config{
                Host:     "localhost",
                Database: "testdb",
                Username: "root",
                Password: "password",
            },
            wantErr: false,
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 4. Register the import

Add to `connectors/connectors.go`:

```go
import (
    _ "github.com/compareflow/compareflow/internal/connectors/mysql"
)
```

### 5. Update models

Add to `internal/models/connection.go`:

```go
const (
    ConnectionTypeMySQL ConnectionType = "mysql"
)
```

## Connector Interface

```go
type Connector interface {
    // Type returns the unique identifier for this connector
    Type() string
    
    // ParseConfig parses a raw configuration map into a validated config struct
    ParseConfig(configMap map[string]interface{}) (interface{}, error)
    
    // Connect establishes a connection to the database
    Connect(config interface{}) (*sql.DB, error)
    
    // TestConnection verifies that the connection works
    TestConnection(config interface{}) error
    
    // GetTables returns a list of tables in the database
    GetTables(db *sql.DB) ([]string, error)
    
    // GetColumns returns column information for a specific table
    GetColumns(db *sql.DB, tableName string) ([]ColumnInfo, error)
}
```

## Testing

Each connector can be tested independently:

```bash
# Test a specific connector
go test ./internal/connectors/mysql

# Test all connectors
go test ./internal/connectors/...

# Test with coverage
go test -cover ./internal/connectors/mysql
```

## Configuration Examples

### SQL Server
```json
{
    "type": "sqlserver",
    "config": {
        "server": "localhost",
        "port": 1433,
        "database": "mydb",
        "username": "sa",
        "password": "password",
        "encrypt": true,
        "trust_server_certificate": false
    }
}
```

### Databricks
```json
{
    "type": "databricks", 
    "config": {
        "workspace": "https://workspace.databricks.com",
        "http_path": "/sql/1.0/endpoints/abc123",
        "access_token": "dapi..."
    }
}
```

### PostgreSQL
```json
{
    "type": "postgresql",
    "config": {
        "host": "localhost",
        "port": 5432,
        "database": "mydb",
        "username": "postgres",
        "password": "password",
        "ssl_mode": "prefer"
    }
}
```

## Best Practices

1. **Keep connectors independent**: Don't import from other connector packages
2. **Validate early**: Validate config in the `Validate()` method
3. **Set sensible defaults**: Set default ports, SSL modes, etc.
4. **Handle errors properly**: Wrap errors with context
5. **Test thoroughly**: Include unit tests for all methods
6. **Document config fields**: Use struct tags and comments
7. **Use prepared statements**: For security and performance

## Security Considerations

1. Never log passwords or sensitive configuration
2. Always use parameterized queries
3. Validate all configuration inputs
4. Use secure defaults (e.g., SSL/TLS enabled)
5. Sanitize error messages to avoid leaking sensitive info