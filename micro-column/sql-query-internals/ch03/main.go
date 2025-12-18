package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch03_db?parseTime=true"

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// 初始化数据 (第一次运行打开，后续注释掉以节省时间)
	initSensorDB(db)

	fmt.Println("--- 🧪 Experiment: The Range Trap ---")

	// 场景 1: 使用 idx_device_type (device_id, type)
	// device_id > 100 (范围), type = 5 (等值)
	// 预期：device_id 命中 Access，type 沦为 Filter
	explainAndPrint(db, "USE INDEX (idx_device_type)", 
		"SELECT * FROM sensor_data USE INDEX (idx_device_type) WHERE device_id > 100 AND type = 5")

	// 场景 2: 使用 idx_type_device (type, device_id)
	// type = 5 (等值), device_id > 100 (范围)
	// 预期：type 命中 Access，device_id 也命中 Access
	explainAndPrint(db, "USE INDEX (idx_type_device)", 
		"SELECT * FROM sensor_data USE INDEX (idx_type_device) WHERE device_id > 100 AND type = 5")
}

func explainAndPrint(db *sql.DB, title, sqlStr string) {
	var output string
	// 使用 JSON 格式获取最详细信息
	err := db.QueryRow("EXPLAIN FORMAT=JSON " + sqlStr).Scan(&output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n📋 Plan for [%s]:\n%s\n", title, output)
}


// 辅助函数：初始化表结构和数据
func initSensorDB(db *sql.DB) {
	// 1. 重建表
    // 我们将创建两个不同的索引来对比
	db.Exec("DROP TABLE IF EXISTS sensor_data")
	schema := `
	CREATE TABLE sensor_data (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		device_id INT,    -- 范围查询字段
		type INT,         -- 等值查询字段
		value VARCHAR(50),
		KEY idx_device_type (device_id, type), -- 索引 A: (范围, 等值)
		KEY idx_type_device (type, device_id)  -- 索引 B: (等值, 范围)
	) ENGINE=InnoDB;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	// 2. 插入 100万 行数据
	fmt.Println("🚀 Start inserting 1M rows into sensor_data...")
	// 开启事务
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO sensor_data (device_id, type, value) VALUES (?, ?, ?)")

    // 构造数据分布：
    // device_id: 1 ~ 10000
    // type: 1 ~ 100
	for i := 0; i < 1000000; i++ {
		deviceID := i%10000 + 1
		sensorType := i%100 + 1
		_, err := stmt.Exec(deviceID, sensorType, "data-payload")
		if err != nil {
			log.Fatal(err)
		}
        // 每 2000 行提交一次，避免事务过大
		if (i+1)%2000 == 0 {
			tx.Commit()
			tx, _ = db.Begin()
			stmt, _ = tx.Prepare("INSERT INTO sensor_data (device_id, type, value) VALUES (?, ?, ?)")
		}
	}
	tx.Commit()

    // 3. 刷新统计信息
    db.Exec("ANALYZE TABLE sensor_data")
	fmt.Println("✅ Insert done.")
}
