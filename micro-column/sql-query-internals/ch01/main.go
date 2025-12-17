package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 定义 DSN (请修改为你本地的配置)
// 关键参数：parseTime=true 用于正确解析时间类型的字段
const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch01_db?parseTime=true"

func main() {
	// 1. 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. 环境初始化 (调用上一节定义的函数)
	initDB(db)

	// 3. 准备我们要透视的 SQL
	// 这是一个典型的范围查询，我们期望它走 idx_age 索引
	targetSQL := "SELECT * FROM users_explain_demo WHERE age > 20"

	// 4. 包装成 EXPLAIN 语句
	// FORMAT=JSON 是 MySQL 5.6+ 引入的神器，能提供比表格形式更详细的成本数据
	explainSQL := "EXPLAIN FORMAT=JSON " + targetSQL

	// 5. 执行查询并读取结果
	// Explain 的 JSON 结果通常是一个巨大的字符串，存在第一行第一列中
	var explainOutput string
	err = db.QueryRow(explainSQL).Scan(&explainOutput)
	if err != nil {
		log.Fatalf("Explain execution failed: %v", err)
	}

	// 6. 打印结果
	fmt.Println("=== 🌟 MySQL Execution Plan (JSON) 🌟 ===")
	// 直接打印 JSON 字符串，后续我们将对其进行解读
	fmt.Println(explainOutput)
}

// initDB 负责重置环境，确保每次实验结果一致
func initDB(db *sql.DB) {
	// 1. 清理旧表
	_, err := db.Exec("DROP TABLE IF EXISTS users_explain_demo")
	if err != nil {
		log.Fatal(err)
	}

	// 2. 建表：包含主键和 age 索引
	// 这是一个典型的 InnoDB 表结构
	schema := `
	CREATE TABLE users_explain_demo (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		age INT,
		created_at DATETIME,
		KEY idx_age (age) -- 关键：我们在 age 上建了二级索引
	) ENGINE=InnoDB;
	`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 预置数据
	// 插入一批数据，让优化器认为走索引比全表扫描更划算
	// 我们构造一些 age > 20 和 age <= 20 的混合数据
	values := []string{}
	for i := 1; i <= 20; i++ {
		// 构造数据：User1 (age 11), User2 (age 12) ... User20 (age 30)
		// 这样 age > 20 的数据大概占一半
		values = append(values, fmt.Sprintf("('User%d', %d, NOW())", i, 10+i))
	}
	insertSQL := "INSERT INTO users_explain_demo (name, age, created_at) VALUES " + strings.Join(values, ",")
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Fatal(err)
	}

    // 4. 强制刷新统计信息 (Analyze Table)
    // 在生产环境不需要手动做，但在这种瞬时创建的小表中，
    // 这一步能帮助优化器更准确地感知数据分布，避免它误判为全表扫描
    _, _ = db.Exec("ANALYZE TABLE users_explain_demo")

	fmt.Println("✅ Environment initialized: Table created and 20 rows inserted.")
}



