package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	dbHost := os.Getenv("DB_HOST_APP") // 从环境变量读取
	dbPort := os.Getenv("DB_PORT_APP")
	dbUser := os.Getenv("DB_USER_APP")
	dbPassword := os.Getenv("DB_PASSWORD_APP")
	dbName := os.Getenv("DB_NAME_APP")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var db *sql.DB
	var err error
	// Retry logic for DB connection
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to PostgreSQL via Docker Compose!")
				break
			}
		}
		log.Printf("DB conn attempt %d failed: %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go App! DB Host from env: %s\n", dbHost)
		// ... (可以加一个简单的DB查询)
	})

	log.Printf("Go app listening on port %s...", appPort)
	http.ListenAndServe(":"+appPort, nil)
}
