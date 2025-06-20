package databricks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	
	_ "github.com/databricks/databricks-sql-go"
	"github.com/compareflow/compareflow/internal/connectors"
)

func init() {
	// Register this connector when the package is imported
	connectors.Register("databricks", func() connectors.Connector {
		return New()
	})
}

// Config represents Databricks connection configuration
type Config struct {
	Workspace   string `json:"workspace" binding:"required"`
	HTTPPath    string `json:"http_path" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Workspace == "" {
		return fmt.Errorf("workspace URL is required")
	}
	if c.HTTPPath == "" {
		return fmt.Errorf("HTTP path is required")
	}
	if c.AccessToken == "" {
		return fmt.Errorf("access token is required")
	}
	return nil
}

// Connector implements the base.Connector interface for Databricks
type Connector struct{}

// New creates a new Databricks connector
func New() *Connector {
	return &Connector{}
}

// Type returns the connector type
func (c *Connector) Type() string {
	return "databricks"
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
		return nil, fmt.Errorf("failed to parse Databricks config: %w", err)
	}
	
	if err := config.Validate(); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// Connect establishes a connection to Databricks
func (c *Connector) Connect(config interface{}) (*sql.DB, error) {
	dbConfig, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected *Config, got %T", config)
	}
	
	connString := c.buildConnectionString(dbConfig)
	return sql.Open("databricks", connString)
}

// TestConnection verifies the Databricks connection works
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

// GetTables returns a list of tables in the Databricks database
func (c *Connector) GetTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT 
			table_schema || '.' || table_name AS full_table_name
		FROM information_schema.tables
		WHERE table_type = 'TABLE'
		  AND table_schema NOT IN ('information_schema', 'system')
		ORDER BY table_schema, table_name
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
		schema = "default"
		table = tableName
	}
	
	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
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

// buildConnectionString builds a Databricks connection string
func (c *Connector) buildConnectionString(config *Config) string {
	// Clean up workspace URL
	host := strings.TrimPrefix(config.Workspace, "https://")
	host = strings.TrimPrefix(host, "http://")
	
	// Format: databricks://token:<access_token>@<workspace>:<port><http_path>
	return fmt.Sprintf("databricks://token:%s@%s:443%s", 
		config.AccessToken, host, config.HTTPPath)
}