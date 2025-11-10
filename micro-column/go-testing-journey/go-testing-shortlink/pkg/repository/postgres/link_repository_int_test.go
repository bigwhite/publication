////go:build integration

package postgres_test

import (
	"context"
	"database/sql" // <-- 引入原生 sql 包
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/repository"
	"github.com/bigwhite/shortlink/pkg/repository/postgres"
)

var dbPool *sql.DB // <-- 关键变更：使用原生的 *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine3.22",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_USER": "test", "POSTGRES_PASSWORD": "password", "POSTGRES_DB": "testdb"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5 * time.Minute),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		log.Fatalf("无法启动 PG 容器: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("无法终止 PG 容器: %s", err)
		}
	}()

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://test:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	// 关键变更：使用 sql.Open，并手动 Ping
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("无法打开数据库连接: %s", err)
	}
	if err := pool.Ping(); err != nil {
		log.Fatalf("无法连接到 PG 容器: %s", err)
	}
	dbPool = pool
	defer dbPool.Close()

	migrator, err := migrate.New("file://../../../migrations", dsn)
	if err != nil {
		log.Fatalf("无法创建 migrator: %s", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("无法执行 up migration: %s", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLinkRepository(t *testing.T) {
	t.Run("SaveAndFindByCode", func(t *testing.T) {
		testCases := []struct {
			name        string
			linkToSave  *domain.Link
			codeToFind  string
			expectFound bool
			expectError bool
		}{
			// ... 测试用例数据维持不变 ...
			{
				name:       "成功保存和查找",
				linkToSave: &domain.Link{OriginalURL: "https://example.com/success", ShortCode: "success"},
				codeToFind: "success", expectFound: true, expectError: false,
			},
			{
				name:       "查找不存在的链接",
				linkToSave: nil, codeToFind: "not-found", expectFound: false, expectError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()

				// 关键变更：使用原生的 BeginTx
				tx, err := dbPool.BeginTx(ctx, nil)
				if err != nil {
					t.Fatalf("无法开始事务: %v", err)
				}
				defer tx.Rollback()

				// *sql.Tx 隐式地满足了我们自定义的 repository.DBTX 接口
				txRepo := postgres.NewPgLinkRepository(tx)

				// ... (后续的测试逻辑完全不变)
				if tc.linkToSave != nil {
					err := txRepo.Save(ctx, tc.linkToSave)
					if (err != nil) != tc.expectError {
						t.Fatalf("Save() error = %v, expectError %v", err, tc.expectError)
					}
					if err != nil {
						return
					}
				}

				foundLink, err := txRepo.FindByCode(ctx, tc.codeToFind)
				if err != nil {
					t.Fatalf("FindByCode() returned an unexpected error: %v", err)
				}

				if (foundLink != nil) != tc.expectFound {
					t.Fatalf("FindByCode() found link = %v, expectFound %v", (foundLink != nil), tc.expectFound)
				}

				if tc.expectFound {
					if foundLink.OriginalURL != tc.linkToSave.OriginalURL {
						t.Errorf("OriginalURL mismatch: got %s, want %s", foundLink.OriginalURL, tc.linkToSave.OriginalURL)
					}
				}
			})
		}
	})

	t.Run("SaveConflict", func(t *testing.T) {
		ctx := context.Background()
		tx, err := dbPool.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("无法开始事务: %v", err)
		}
		defer tx.Rollback()

		txRepo := postgres.NewPgLinkRepository(tx)

		// ... (后续的测试逻辑完全不变)
		link1 := &domain.Link{OriginalURL: "https://test.com", ShortCode: "conflict"}
		err = txRepo.Save(ctx, link1)
		if err != nil {
			t.Fatalf("第一次保存不应返回错误: %v", err)
		}

		link2 := &domain.Link{OriginalURL: "https://another.com", ShortCode: "conflict"}
		err = txRepo.Save(ctx, link2)

		if err == nil {
			t.Fatal("第二次保存应该返回错误，但得到了 nil")
		}
		if !errors.Is(err, repository.ErrCodeConflict) {
			t.Errorf("错误类型应该是 ErrCodeConflict，但得到了 %v", err)
		}
	})
}
