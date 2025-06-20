package sqlserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/compareflow/compareflow/internal/connectors"
)

func init() {
	// Auto-register when package is imported
	connectors.Register("sqlserver", func() connectors.Connector {
		return New()
	})
}

// Config represents SQL Server connection configuration
type Config struct {
	Server                 string `json:"server" binding:"required"`
	Port                   int    `json:"port"`
	Database               string `json:"database" binding:"required"`
	Username               string `json:"username" binding:"required"`
	Password               string `json:"password" binding:"required"`
	Encrypt                bool   `json:"encrypt"`
	TrustServerCertificate bool   `json:"trust_server_certificate"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server == "" {
		return fmt.Errorf("server is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database is required")
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	if c.Port == 0 {
		c.Port = 1433 // Set default
	}
	return nil
}

// Connector implements the Connector interface for SQL Server
type Connector struct{}

// New creates a new SQL Server connector
func New() *Connector {
	return &Connector{}
}

// Type returns the connector type
func (c *Connector) Type() string {
	return "sqlserver"
}

// ParseConfig parses raw config into typed config
func (c *Connector) ParseConfig(configMap map[string]interface{}) (interface{}, error) {
	var config Config
	
	// Convert map to JSON and back to struct
	jsonBytes, err := json.Marshal(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse SQL Server config: %w", err)
	}
	
	if err := config.Validate(); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// Connect establishes a connection to SQL Server
func (c *Connector) Connect(config interface{}) (*sql.DB, error) {
	sqlConfig, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected *Config, got %T", config)
	}
	
	connString := c.buildConnectionString(sqlConfig)
	return sql.Open("sqlserver", connString)
}

// TestConnection verifies the SQL Server connection works
func (c *Connector) TestConnection(config interface{}) error {
	db, err := c.Connect(config)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()
	
	// Ping the database
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Run a simple query
	var result int
	if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}
	
	return nil
}

// GetTables returns a list of tables in the SQL Server database
func (c *Connector) GetTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT TABLE_SCHEMA + '.' + TABLE_NAME
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_SCHEMA, TABLE_NAME
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()
	
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	
	return tables, rows.Err()
}

// GetColumns returns column information for a specific table
func (c *Connector) GetColumns(db *sql.DB, tableName string) ([]connectors.ColumnInfo, error) {
	// Parse schema and table name
	parts := strings.Split(tableName, ".")
	var schema, table string
	if len(parts) == 2 {
		schema = parts[0]
		table = parts[1]
	} else {
		schema = "dbo"
		table = tableName
	}
	
	query := `
		SELECT 
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = @p1 AND TABLE_NAME = @p2
		ORDER BY ORDINAL_POSITION
	`
	
	rows, err := db.Query(query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()
	
	var columns []connectors.ColumnInfo
	for rows.Next() {
		var col connectors.ColumnInfo
		var nullable string
		if err := rows.Scan(&col.Name, &col.DataType, &nullable); err != nil {
			return nil, err
		}
		col.Nullable = nullable == "YES"
		columns = append(columns, col)
	}
	
	return columns, rows.Err()
}

// buildConnectionString builds a SQL Server connection string
func (c *Connector) buildConnectionString(config *Config) string {
	encrypt := "false"
	if config.Encrypt {
		encrypt = "true"
	}
	
	trustServerCertificate := "true"
	if !config.TrustServerCertificate {
		trustServerCertificate = "false"
	}
	
	return fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=%s;TrustServerCertificate=%s",
		config.Server, config.Port, config.Database, config.Username, config.Password, 
		encrypt, trustServerCertificate)
}