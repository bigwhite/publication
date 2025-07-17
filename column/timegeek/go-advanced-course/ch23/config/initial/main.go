package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

const defaultPort = ":8080"

var databaseDSN = "user:password@tcp(localhost:3306)/mydb?parseTime=true" // 全局变量

func main() {
	fmt.Printf("Application starting on port %s...\n", defaultPort)

	db, err := sql.Open("mysql", databaseDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database (this is expected if DB is not running): %v", err)
	} else {
		fmt.Println("Successfully connected to database (or DSN is valid).")
	}

	http.ListenAndServe(defaultPort, nil) // Commented out to run without actual server

	fmt.Println("Application finished (simulated run).")
}
