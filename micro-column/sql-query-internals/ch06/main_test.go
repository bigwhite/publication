package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch06_db?parseTime=true"

// initPagingDB 初始化百万级数据
func initPagingDB(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS paging_demo")
	schema := `
	CREATE TABLE paging_demo (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		create_time DATETIME,
		payload VARCHAR(100),
		KEY idx_create_time (create_time)
	) ENGINE=InnoDB;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	fmt.Println("🚀 Start inserting 1M rows (this may take a while)...")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO paging_demo (create_time, payload) VALUES (NOW(), ?)")
	
	payload := strings.Repeat("A", 100) // 模拟真实负载

	// 插入 100万 行
	for i := 0; i < 1000000; i++ {
		_, err := stmt.Exec(payload)
		if err != nil {
			log.Fatal(err)
		}
		if (i+1)%5000 == 0 {
			tx.Commit()
			tx, _ = db.Begin()
			stmt, _ = tx.Prepare("INSERT INTO paging_demo (create_time, payload) VALUES (NOW(), ?)")
            if (i+1)%100000 == 0 {
			    fmt.Printf("Inserted %d rows...\n", i+1)
            }
		}
	}
	tx.Commit()
    db.Exec("ANALYZE TABLE paging_demo")
	fmt.Println("✅ Insert done.")
}

var globalDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	globalDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
    // 首次运行需解开注释初始化数据
	initPagingDB(globalDB)
	defer globalDB.Close()
	m.Run()
}

// 场景一：Offset 分页 (LIMIT N, 10)
// 随着 N 增大，性能会急剧下降
func benchmarkOffsetPaging(b *testing.B, offset int) {
	query := "SELECT id, payload FROM paging_demo ORDER BY id LIMIT ?, 10"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query, offset)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var id int
			var payload string
			rows.Scan(&id, &payload)
		}
		rows.Close()
	}
}

// 场景二：Seek 分页 (WHERE id > last_id LIMIT 10)
// 无论 last_id 是多少，性能应该保持稳定
func benchmarkSeekPaging(b *testing.B, lastID int) {
	query := "SELECT id, payload FROM paging_demo WHERE id > ? ORDER BY id LIMIT 10"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query, lastID)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var id int
			var payload string
			rows.Scan(&id, &payload)
		}
		rows.Close()
	}
}

// ------ 定义不同深度的 Benchmark ------

// 深度 0 (第 1 页)
func BenchmarkOffset_Page1(b *testing.B) { benchmarkOffsetPaging(b, 0) }
func BenchmarkSeek_Page1(b *testing.B)   { benchmarkSeekPaging(b, 0) }

// 深度 50,000 (第 5000 页)
func BenchmarkOffset_Page5k(b *testing.B) { benchmarkOffsetPaging(b, 50000) }
func BenchmarkSeek_Page5k(b *testing.B)   { benchmarkSeekPaging(b, 50000) }

// 深度 900,000 (第 9万 页) - 接近表尾部
func BenchmarkOffset_Page90k(b *testing.B) { benchmarkOffsetPaging(b, 900000) }
func BenchmarkSeek_Page90k(b *testing.B)   { benchmarkSeekPaging(b, 900000) }
