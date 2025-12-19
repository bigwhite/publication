package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:your_strong_password@tcp(127.0.0.1:3306)/ch04_db?parseTime=true"

// setupJoinDB 初始化表结构和数据
func setupJoinDB() *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 1. 重建表
	db.Exec("DROP TABLE IF EXISTS orders")
	db.Exec("DROP TABLE IF EXISTS users")
	
	schemaUsers := `
	CREATE TABLE users (
		id INT PRIMARY KEY,
		name VARCHAR(50)
	) ENGINE=InnoDB;`
	
	schemaOrders := `
	CREATE TABLE orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		uid INT,
		amount INT,
		KEY idx_uid (uid) -- 关键：被驱动表的连接字段必须有索引
	) ENGINE=InnoDB;`

	db.Exec(schemaUsers)
	db.Exec(schemaOrders)

	// 2. 插入数据
	// 5000 用户
	fmt.Println("🚀 Inserting 5000 users...")
	tx, _ := db.Begin()
	stmtUser, _ := tx.Prepare("INSERT INTO users (id, name) VALUES (?, ?)")
	for i := 1; i <= 5000; i++ {
		stmtUser.Exec(i, fmt.Sprintf("User-%d", i))
	}
	tx.Commit()

	// 10万 订单 (平均每个用户 20 单)
	fmt.Println("🚀 Inserting 100k orders...")
	tx, _ = db.Begin()
	stmtOrder, _ := tx.Prepare("INSERT INTO orders (uid, amount) VALUES (?, ?)")
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++ {
		uid := rand.Intn(5000) + 1
		stmtOrder.Exec(uid, rand.Intn(1000))
		if (i+1)%5000 == 0 {
			tx.Commit()
			tx, _ = db.Begin()
			stmtOrder, _ = tx.Prepare("INSERT INTO orders (uid, amount) VALUES (?, ?)")
		}
	}
	tx.Commit()
    
    // 强制更新统计信息
    db.Exec("ANALYZE TABLE users")
    db.Exec("ANALYZE TABLE orders")
	fmt.Println("✅ Data ready.")
	return db
}

var globalDB *sql.DB

func TestMain(m *testing.M) {
	globalDB = setupJoinDB()
	defer globalDB.Close()
	m.Run()
}


// 场景一：N+1 查询 (应用层 Loop)
// 模拟：先查出前 100 个用户，再遍历查询他们的订单
func BenchmarkNPlusOne(b *testing.B) {
	// 为了公平，我们只查前 100 个用户，避免 N+1 慢到跑不完
    limit := 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
        // Step 1: 查用户
		userRows, _ := globalDB.Query("SELECT id FROM users LIMIT ?", limit)
		var uids []int
		for userRows.Next() {
			var id int
			userRows.Scan(&id)
			uids = append(uids, id)
		}
		userRows.Close()

        // Step 2: 循环查订单 (N 次查询)
		for _, uid := range uids {
			orderRows, _ := globalDB.Query("SELECT id, amount FROM orders WHERE uid = ?", uid)
			for orderRows.Next() {
				var oid, amt int
				orderRows.Scan(&oid, &amt)
			}
			orderRows.Close()
		}
	}
}

// 场景二：SQL Join
// 模拟：一条 SQL 一次性查出 100 个用户及其订单
func BenchmarkSQLJoin(b *testing.B) {
    limit := 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 这里的写法利用了子查询 limit，确保业务语义一致
		query := `
		SELECT u.id, o.id, o.amount
		FROM (SELECT id FROM users LIMIT ?) u
		JOIN orders o ON u.id = o.uid`

		rows, err := globalDB.Query(query, limit)
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var uid, oid, amt int
			rows.Scan(&uid, &oid, &amt)
		}
		rows.Close()
	}
}
