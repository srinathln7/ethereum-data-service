package main

import (
	"log"

	"ethereum-data-service/cmd"
)

func main() {
	// Execute the Cobra commands
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal("error:", err)
	}
}
