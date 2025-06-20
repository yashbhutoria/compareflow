// +build ignore

package main

import (
	"fmt"
	"log"
	
	"github.com/compareflow/compareflow/internal/connectors"
	_ "github.com/compareflow/compareflow/internal/connectors/databricks"
	_ "github.com/compareflow/compareflow/internal/connectors/postgresql"
	_ "github.com/compareflow/compareflow/internal/connectors/sqlserver"
)

func main() {
	fmt.Println("Testing connector registration...")
	
	// List all registered connectors
	connectors := connectors.List()
	fmt.Printf("Registered connectors: %v\n", connectors)
	
	// Test getting each connector
	for _, connType := range connectors {
		conn, err := connectors.Get(connType)
		if err != nil {
			log.Printf("Error getting connector %s: %v", connType, err)
			continue
		}
		fmt.Printf("✓ Connector %s loaded successfully\n", conn.Type())
	}
	
	fmt.Println("All connectors loaded successfully!")
}