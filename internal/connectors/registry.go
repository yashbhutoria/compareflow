package connectors

import (
	"fmt"
	"sync"
)

// Registry holds all available connectors
var (
	registry = make(map[string]func() Connector)
	mu       sync.RWMutex
)

// Register adds a new connector to the registry
func Register(name string, factory func() Connector) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = factory
}

// Get returns a connector instance for the given type
func Get(connectorType string) (Connector, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	factory, exists := registry[connectorType]
	if !exists {
		return nil, fmt.Errorf("unsupported connector type: %s", connectorType)
	}
	return factory(), nil
}

// List returns all registered connector types
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	
	types := make([]string, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}