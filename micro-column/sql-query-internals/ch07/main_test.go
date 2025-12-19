
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// 请修改为你本地的数据库连接串
const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch07_db?parseTime=true"

// InsertCount 设定为 30万，配合 padding 足以填满默认的 Buffer Pool (128MB)
// 从而体现出随机 IO 的劣势
const InsertCount = 300000

// BatchSize 分批提交的大小，避免单次事务过大
const BatchSize = 5000

var globalDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	globalDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer globalDB.Close()

	// ⚠️ 必须执行初始化，创建带 padding 的表结构
	initWriteDB(globalDB)

	// 运行测试
	m.Run()
}

// initWriteDB 初始化三张表，增加了 padding 字段以快速消耗内存
func initWriteDB(db *sql.DB) {
	tables := []string{"users_baseline", "users_heavy_index", "users_uuid_bad"}
	for _, t := range tables {
		db.Exec("DROP TABLE IF EXISTS " + t)
	}

	// 1. Baseline: 只有主键，带 padding
	_, err := db.Exec(`CREATE TABLE users_baseline (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		col1 INT, col2 INT, col3 VARCHAR(50), col4 VARCHAR(50), col5 TIMESTAMP,
		padding CHAR(200)
	) ENGINE=InnoDB`)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Heavy Index: 主键 + 5个二级索引 + padding
	_, err = db.Exec(`CREATE TABLE users_heavy_index (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		col1 INT, col2 INT, col3 VARCHAR(50), col4 VARCHAR(50), col5 TIMESTAMP,
		padding CHAR(200),
		KEY idx_1 (col1), KEY idx_2 (col2), KEY idx_3 (col3), KEY idx_4 (col4), KEY idx_5 (col5)
	) ENGINE=InnoDB`)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Bad UUID: UUID主键 (乱序) + 5个二级索引 + padding
	_, err = db.Exec(`CREATE TABLE users_uuid_bad (
		id VARCHAR(36) PRIMARY KEY,
		col1 INT, col2 INT, col3 VARCHAR(50), col4 VARCHAR(50), col5 TIMESTAMP,
		padding CHAR(200),
		KEY idx_1 (col1), KEY idx_2 (col2), KEY idx_3 (col3), KEY idx_4 (col4), KEY idx_5 (col5)
	) ENGINE=InnoDB`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Tables created (with padding).")
}

func runInsert(b *testing.B, tableName string, useUUID bool) {
	// 每次 Benchmark 前清空表，确保环境一致
	b.StopTimer()
	if _, err := globalDB.Exec("TRUNCATE TABLE " + tableName); err != nil {
		b.Fatalf("Truncate table failed: %v", err)
	}

	// 准备注水数据 (200字节)
	paddingData := strings.Repeat("A", 200)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// 计算需要多少个批次
		totalBatches := InsertCount / BatchSize
		if InsertCount%BatchSize != 0 {
			totalBatches++
		}

		for batch := 0; batch < totalBatches; batch++ {
			tx, err := globalDB.Begin()
			if err != nil {
				b.Fatal(err)
			}

			var stmt *sql.Stmt
			var prepErr error

			// 准备 SQL，包含 padding 列
			if useUUID {
				stmt, prepErr = tx.Prepare("INSERT INTO " + tableName + " VALUES (?, ?, ?, ?, ?, ?, ?)")
			} else {
				stmt, prepErr = tx.Prepare("INSERT INTO " + tableName + " (col1, col2, col3, col4, col5, padding) VALUES (?, ?, ?, ?, ?, ?)")
			}

			if prepErr != nil {
				tx.Rollback()
				b.Fatalf("Prepare failed: %v", prepErr)
			}

			// 批量插入
			for j := 0; j < BatchSize; j++ {
				var execErr error
				if useUUID {
					_, execErr = stmt.Exec(uuid.New().String(), j, j, "data", "data", time.Now(), paddingData)
				} else {
					_, execErr = stmt.Exec(j, j, "data", "data", time.Now(), paddingData)
				}

				if execErr != nil {
					stmt.Close()
					tx.Rollback()
					b.Fatalf("Exec failed: %v", execErr)
				}
			}

			stmt.Close()
			if err := tx.Commit(); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkInsert_Baseline(b *testing.B) {
	runInsert(b, "users_baseline", false)
}

func BenchmarkInsert_HeavyIndex(b *testing.B) {
	runInsert(b, "users_heavy_index", false)
}

func BenchmarkInsert_UUID_Bad(b *testing.B) {
	runInsert(b, "users_uuid_bad", true)
}
