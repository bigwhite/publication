//go:build integration

package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

      "github.com/bigwhite/shortlink/pkg/repository/postgres"
      redis_repo "github.com/bigwhite/shortlink/pkg/repository/redis"
      "github.com/bigwhite/shortlink/pkg/service"
)

var (
	dbPool      *sql.DB
	redisClient *redis.Client
	testService *service.ShortenerService
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. 定义网络名称并创建共享网络
	networkName := "shortlink-test-network" // <-- 关键变更：将网络名称保存到变量中
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: networkName},
	})
	if err != nil {
		log.Fatalf("无法创建共享网络: %s", err)
	}
	defer network.Remove(ctx)

	// 2. 启动 PostgreSQL 容器，并加入网络
	pgContainer := setupPostgres(ctx, networkName) // <-- 传递网络名称
	defer pgContainer.Terminate(ctx)

	// 3. 启动 Redis 容器，并加入网络
	redisContainer := setupRedis(ctx, networkName) // <-- 传递网络名称
	defer redisContainer.Terminate(ctx)

	// ... (后续的连接和初始化逻辑完全不变)
	pgHost, _ := pgContainer.Host(ctx)
	pgPort, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://test:password@%s:%s/testdb?sslmode=disable", pgHost, pgPort.Port())
	pool, err := sql.Open("postgres", dsn)
	if err != nil { log.Fatalf("无法连接 PG: %s", err) }
	dbPool = pool
	defer dbPool.Close()

	redisHost, _ := redisContainer.Host(ctx)
	redisPort, _ := redisContainer.MappedPort(ctx, "6379")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := rdb.Ping(ctx).Err(); err != nil { log.Fatalf("无法连接 Redis: %s", err) }
	redisClient = rdb
	defer redisClient.Close()
	
	migrator, _ := migrate.New("file://../../migrations", dsn)
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("无法执行 migration: %s", err)
	}

	linkRepo := postgres.NewPgLinkRepository(dbPool)
	linkCache := redis_repo.NewRedisLinkCache(redisClient)
	testService = service.NewShortenerService(linkRepo, linkCache)

	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupPostgres 接收网络名称字符串，而不是 network 对象
func setupPostgres(ctx context.Context, networkName string) testcontainers.Container { // <-- 关键变更
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine3.22",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_USER": "test", "POSTGRES_PASSWORD": "password", "POSTGRES_DB": "testdb"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		Networks:     []string{networkName}, // <-- 使用网络名称
		NetworkAliases: map[string][]string{networkName: {"postgres-db"}}, // <-- 使用网络名称
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil { log.Fatalf("PG 启动失败: %s", err) }
	return container
}

// setupRedis 接收网络名称字符串，而不是 network 对象
func setupRedis(ctx context.Context, networkName string) testcontainers.Container { // <-- 关键变更
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
		Networks:     []string{networkName}, // <-- 使用网络名称
		NetworkAliases: map[string][]string{networkName: {"redis-cache"}}, // <-- 使用网络名称
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil { log.Fatalf("Redis 启动失败: %s", err) }
	return container
}

// --- TestCreateAndRedirect_HappyPath 函数及其辅助函数 assertEventually 完全不变 ---

func TestCreateAndRedirect_HappyPath(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://www.google.com/very-long-path"

	tx, err := dbPool.BeginTx(ctx, nil)
	if err != nil { t.Fatal(err) }
	defer tx.Rollback()
	
	txRepo := postgres.NewPgLinkRepository(tx)
	txTestService := service.NewShortenerService(txRepo, redis_repo.NewRedisLinkCache(redisClient))
	
	createdLink, err := txTestService.CreateLink(ctx, originalURL)
	if err != nil { t.Fatalf("CreateLink 不应返回错误: %v", err) }

	redirectedLink, err := txTestService.Redirect(ctx, createdLink.ShortCode)
	if err != nil { t.Fatalf("Redirect 不应返回错误: %v", err) }
	if redirectedLink == nil || redirectedLink.OriginalURL != originalURL {
		t.Fatalf("重定向的链接不正确")
	}

	assertEventually(t, func() bool {
		count, err := redisClient.Get(ctx, fmt.Sprintf("link:visits:%s", createdLink.ShortCode)).Int64()
		if err != nil { return false }
		return count == 1
	}, "访问计数应该在 Redis 中变为 1")

	_, err = txTestService.Redirect(ctx, createdLink.ShortCode)
	if err != nil { t.Fatalf("第二次 Redirect 不应返回错误: %v", err) }
	
	assertEventually(t, func() bool {
		count, err := redisClient.Get(ctx, fmt.Sprintf("link:visits:%s", createdLink.ShortCode)).Int64()
		if err != nil { return false }
		return count == 2
	}, "再次访问后，计数应该在 Redis 中变为 2")

	redisClient.Del(ctx, fmt.Sprintf("link:visits:%s", createdLink.ShortCode))
}

func assertEventually(t *testing.T, condition func() bool, msgAndArgs ...interface{}) {
	t.Helper()
	const (
		timeout  = 2 * time.Second
		interval = 50 * time.Millisecond
	)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-timer.C:
			t.Fatalf("Condition was not met within %v: %s", timeout, fmt.Sprint(msgAndArgs...))
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}
