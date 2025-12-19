
package main

import (
	"database/sql"
	"testing"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch05_db?parseTime=true"

// initSortDB 初始化数据
func initSortDB(db *sql.DB) {
	// 1. 重建表
	db.Exec("DROP TABLE IF EXISTS users_sort_demo")
	schema := `
	CREATE TABLE users_sort_demo (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		age INT,          -- 有索引，用于 Pipeline 演示
		score INT,        -- 无索引，用于 Filesort 演示
		payload VARCHAR(200),
		KEY idx_age (age)
	) ENGINE=InnoDB;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	// 2. 插入 50万 行数据
    // 数量必须足够大，才能让 Sort Buffer 溢出或体现 CPU 差距
	fmt.Println("🚀 Start inserting 500k rows...")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO users_sort_demo (name, age, score, payload) VALUES (?, ?, ?, ?)")

	rand.Seed(time.Now().UnixNano())
	payload := strings.Repeat("X", 100) // 增加行宽，让排序更占内存

	for i := 0; i < 500000; i++ {
		age := rand.Intn(100)
		score := rand.Intn(10000)
		_, err := stmt.Exec(fmt.Sprintf("User-%d", i), age, score, payload)
		if err != nil {
			log.Fatal(err)
		}
		if (i+1)%5000 == 0 {
			tx.Commit()
			tx, _ = db.Begin()
			stmt, _ = tx.Prepare("INSERT INTO users_sort_demo (name, age, score, payload) VALUES (?, ?, ?, ?)")
			fmt.Printf("\rInserted %d rows...", i+1)
		}
	}
	tx.Commit()
    fmt.Println("\n✅ Insert done. Analyzing table...")
    db.Exec("ANALYZE TABLE users_sort_demo")
}


var globalDB *sql.DB

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	globalDB = db
	initSortDB(db)
	defer globalDB.Close()
	m.Run()
}

// 场景一：Pipeline (利用索引有序性)
// SQL: SELECT ... ORDER BY age LIMIT 1000
func BenchmarkPipelineSort(b *testing.B) {
	// age 上有索引，MySQL 直接扫索引的前 1000 条
	query := "SELECT id, age, score FROM users_sort_demo ORDER BY age LIMIT 1000"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var id, age, score int
			rows.Scan(&id, &age, &score)
		}
		rows.Close()
	}
}

// 场景二：Filesort (内存排序)
// SQL: SELECT ... ORDER BY score LIMIT 1000
func BenchmarkFilesort(b *testing.B) {
	// score 上无索引，MySQL 必须全表扫描 50万行 -> 放入 Sort Buffer -> 排序 -> 取前 1000
	query := "SELECT id, age, score FROM users_sort_demo ORDER BY score LIMIT 1000"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var id, age, score int
			rows.Scan(&id, &age, &score)
		}
		rows.Close()
	}
}
