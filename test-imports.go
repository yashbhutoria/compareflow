// +build ignore

package main

import (
	"fmt"
	"log"
	
	"github.com/compareflow/compareflow/internal/connectors"
)

func main() {
	fmt.Println("Testing import cycle fix...")
	
	types := connectors.List()
	fmt.Printf("Available connectors: %v\n", types)
	
	if len(types) == 0 {
		log.Println("No connectors registered yet")
	} else {
		fmt.Println("Import cycle fixed!")
	}
}