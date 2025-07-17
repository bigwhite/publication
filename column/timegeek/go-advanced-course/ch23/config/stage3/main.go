package main

import (
	"ch23/config/stage3/config"
	"fmt"
	"log"
)

func main() {
	err := config.Load("./app.yaml") // Adjust path
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appName, ok := config.GetString("appName")
	if ok {
		fmt.Printf("appName from GetString: %s\n", appName)
	} else {
		fmt.Println("appName not found or not a string.")
	}

	port, ok := config.GetInt("server.port")
	if ok {
		fmt.Printf("server.port from GetInt: %d\n", port)
	} else {
		fmt.Println("server.port not found or not an int.")
	}

	_, ok = config.GetString("server.host")
	if !ok {
		fmt.Println("server.host (non-existent) correctly not found.")
	}
}
