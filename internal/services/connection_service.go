package services

import (
	"fmt"

	"github.com/compareflow/compareflow/internal/connectors"
	"github.com/compareflow/compareflow/internal/models"
)

// ConnectionService handles database connection operations
type ConnectionService struct {
	// Could add caching or other dependencies here
}

// NewConnectionService creates a new connection service instance
func NewConnectionService() *ConnectionService {
	return &ConnectionService{}
}

// TestConnection tests a database connection
func (s *ConnectionService) TestConnection(conn *models.Connection) error {
	// Get the appropriate connector
	connector, err := connectors.Get(string(conn.Type))
	if err != nil {
		return err
	}
	
	// Parse the config
	config, err := connector.ParseConfig(conn.Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	
	return connector.TestConnection(config)
}

// GetTables returns a list of tables from the connection
func (s *ConnectionService) GetTables(conn *models.Connection) ([]string, error) {
	// Get the appropriate connector
	connector, err := connectors.Get(string(conn.Type))
	if err != nil {
		return nil, err
	}
	
	// Parse the config
	config, err := connector.ParseConfig(conn.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Connect to database
	db, err := connector.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()
	
	return connector.GetTables(db)
}

// GetColumns returns column information for a specific table
func (s *ConnectionService) GetColumns(conn *models.Connection, tableName string) ([]connectors.ColumnInfo, error) {
	// Get the appropriate connector
	connector, err := connectors.Get(string(conn.Type))
	if err != nil {
		return nil, err
	}
	
	// Parse the config
	config, err := connector.ParseConfig(conn.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Connect to database
	db, err := connector.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()
	
	return connector.GetColumns(db, tableName)
}

// ValidateConnectionConfig validates the configuration for a connection
func (s *ConnectionService) ValidateConnectionConfig(connType models.ConnectionType, configMap map[string]interface{}) error {
	// Get the appropriate connector
	connector, err := connectors.Get(string(connType))
	if err != nil {
		return err
	}
	
	// Parse and validate the config
	_, err = connector.ParseConfig(configMap)
	return err
}

// ExecuteQuery executes a query on the connection and returns the results
func (s *ConnectionService) ExecuteQuery(conn *models.Connection, query string) ([]map[string]interface{}, error) {
	// Get the appropriate connector
	connector, err := connectors.Get(string(conn.Type))
	if err != nil {
		return nil, err
	}
	
	// Parse the config
	config, err := connector.ParseConfig(conn.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Connect to database
	db, err := connector.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()
	
	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	// Prepare result slice
	var results []map[string]interface{}
	
	// Create a slice of interface{} to hold each column value
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	
	// Iterate through rows
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			
			// Handle byte arrays (common for some SQL types)
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			
			rowMap[col] = val
		}
		
		results = append(results, rowMap)
	}
	
	return results, rows.Err()
}

// GetSupportedConnectors returns a list of all supported connector types
func (s *ConnectionService) GetSupportedConnectors() []string {
	return connectors.List()
}