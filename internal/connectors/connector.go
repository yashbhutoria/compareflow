package connectors

import (
	"database/sql"
)

// Connector defines the interface that all database connectors must implement
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

// ColumnInfo represents information about a database column
type ColumnInfo struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Nullable bool   `json:"nullable"`
}