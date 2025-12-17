
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch02_db?parseTime=true"

// setupDB 初始化表和数据
func setupDB() *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 1. 重建表
	db.Exec("DROP TABLE IF EXISTS users_index_demo")
	schema := `
	CREATE TABLE users_index_demo (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		age INT,
		padding VARCHAR(200), -- 增加行大小，让回表 I/O 成本更显著
		KEY idx_age (age)
	) ENGINE=InnoDB;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	// 2. 批量插入 10万 行数据
	// 这是一个耗时操作，通常只在第一次运行时做，或者在 bench 前做
	fmt.Println("Start inserting 100k rows...")
	batchSize := 1000
	totalRows := 100000

    // 开启事务加速插入
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO users_index_demo (name, age, padding) VALUES (?, ?, ?)")

	rand.Seed(time.Now().UnixNano())
	paddingData := strings.Repeat("A", 200)

	for i := 0; i < totalRows; i++ {
		age := rand.Intn(100) // age 0-99
		name := fmt.Sprintf("User-%d", i)
		_, err := stmt.Exec(name, age, paddingData)
		if err != nil {
			log.Fatal(err)
		}
		if (i+1)%batchSize == 0 {
            // 分批提交
			tx.Commit()
			tx, _ = db.Begin()
			stmt, _ = tx.Prepare("INSERT INTO users_index_demo (name, age, padding) VALUES (?, ?, ?)")
		}
	}
	tx.Commit()

    // 3. 强制分析表，更新统计信息
    db.Exec("ANALYZE TABLE users_index_demo")
	fmt.Println("Insert done.")
	return db
}

var globalDB *sql.DB

func TestMain(m *testing.M) {
	globalDB = setupDB()
    defer globalDB.Close()
	m.Run()
}

// 场景一：覆盖索引 (只查 id)
// 不需要回表，直接在 idx_age 索引树上就能拿全数据
func BenchmarkCoveringIndex(b *testing.B) {
	query := "SELECT id FROM users_index_demo WHERE age = 50"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query)
		if err != nil {
			b.Fatal(err)
		}
		// 必须遍历完 rows，确保数据库真的把数据发完了
		for rows.Next() {
			var id int
			_ = rows.Scan(&id)
		}
		rows.Close()
	}
}

// 场景二：回表查询 (查 *)
// 需要从 idx_age 拿到 id，再跳回 Cluster Index 拿 name 和 padding
func BenchmarkTableAccess(b *testing.B) {
	query := "SELECT * FROM users_index_demo WHERE age = 50"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := globalDB.Query(query)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var id int
			var name, padding string
			var age int
			_ = rows.Scan(&id, &name, &age, &padding)
		}
		rows.Close()
	}
}


