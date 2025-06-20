package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	
	_ "github.com/lib/pq"
	"github.com/compareflow/compareflow/internal/connectors"
)

func init() {
	// Auto-register when package is imported
	connectors.Register("postgresql", func() connectors.Connector {
		return New()
	})
}

// Config represents PostgreSQL connection configuration
type Config struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port"`
	Database string `json:"database" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	SSLMode  string `json:"ssl_mode"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
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
		c.Port = 5432 // Set default
	}
	if c.SSLMode == "" {
		c.SSLMode = "prefer" // Set default
	}
	return nil
}

// Connector implements the Connector interface for PostgreSQL
type Connector struct{}

// New creates a new PostgreSQL connector
func New() *Connector {
	return &Connector{}
}

// Type returns the connector type
func (c *Connector) Type() string {
	return "postgresql"
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
		return nil, fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}
	
	if err := config.Validate(); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// Connect establishes a connection to PostgreSQL
func (c *Connector) Connect(config interface{}) (*sql.DB, error) {
	pgConfig, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected *Config, got %T", config)
	}
	
	connString := c.buildConnectionString(pgConfig)
	return sql.Open("postgres", connString)
}

// TestConnection verifies the PostgreSQL connection works
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

// GetTables returns a list of tables in the PostgreSQL database
func (c *Connector) GetTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT schemaname || '.' || tablename
		FROM pg_tables
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY schemaname, tablename
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
		schema = "public"
		table = tableName
	}
	
	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position
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

// buildConnectionString builds a PostgreSQL connection string
func (c *Connector) buildConnectionString(config *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, 
		config.Database, config.SSLMode)
}