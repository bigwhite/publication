//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"crypto/rand"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/repository/fakes"
)

// 这个测试文件现在是完全独立的，不依赖任何包级别的变量篡改

func TestShortenerService_CreateLink(t *testing.T) {
	testCases := []struct {
		name           string
		originalURL    string
		generatedCodes []string // 模拟短码生成函数依次返回的值
		setupRepo      func(*fakes.FakeLinkRepository)
		wantErrMsg     string // 期望的错误信息，如果为空则表示不期望错误
	}{
		{
			name:           "成功创建",
			originalURL:    "https://example.com/a-very-long-url-that-needs-to-be-shortened",
			generatedCodes: []string{"abcde1"},
			setupRepo:      func(repo *fakes.FakeLinkRepository) {},
			wantErrMsg:     "",
		},
		{
			name:           "URL 无效",
			originalURL:    "not-a-valid-url",
			generatedCodes: []string{}, // 不会调用
			setupRepo:      func(repo *fakes.FakeLinkRepository) {},
			wantErrMsg:     "invalid URL",
		},
		{
			name:           "一次冲突后成功",
			originalURL:    "https://google.com",
			generatedCodes: []string{"g00gle", "g00gl1"}, // 第一次冲突，第二次成功
			setupRepo: func(repo *fakes.FakeLinkRepository) {
				// 预设一个冲突的短码
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "g00gle", OriginalURL: "https://some-other-url.com"})
			},
			wantErrMsg: "",
		},
		{
			name:           "重试耗尽后失败",
			originalURL:    "https://bing.com",
			generatedCodes: []string{"bing1", "bing2", "bing3", "bing4", "bing5"}, // 所有生成的码都将冲突
			setupRepo: func(repo *fakes.FakeLinkRepository) {
				// 预设所有生成的码都冲突
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "bing1"})
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "bing2"})
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "bing3"})
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "bing4"})
				_ = repo.Save(context.Background(), &domain.Link{ShortCode: "bing5"})
			},
			wantErrMsg: "failed to create a unique short code after multiple retries",
		},
		{
			name:           "短码生成器本身出错",
			originalURL:    "https://duck.com",
			generatedCodes: []string{}, // 模拟生成器立即出错
			setupRepo:      func(repo *fakes.FakeLinkRepository) {},
			wantErrMsg:     "generator failed", // 自定义一个错误
		},
	}

	for _, tc := range testCases {
		// 将 tc 捕获到循环变量中，以确保在并行测试中每个 goroutine 都能获取到正确的 tc
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // 标记此子测试可以与其他子测试并行运行

			// 1. 准备依赖
			fakeRepo := fakes.NewFakeLinkRepository()
			tc.setupRepo(fakeRepo)

			fakeCache := fakes.FakeLinkCache{}
			service := NewShortenerService(fakeRepo, &fakeCache)

			// 2. 注入可预测的、专属于此子测试的短码生成函数
			codes := make([]string, len(tc.generatedCodes))
			copy(codes, tc.generatedCodes)

			service.generateShortCodeFunc = func(length int) (string, error) {
				if tc.wantErrMsg == "generator failed" {
					return "", errors.New("generator failed")
				}
				if len(codes) == 0 {
					return "", errors.New("test setup error: not enough codes provided")
				}
				code := codes[0]
				codes = codes[1:]
				return code, nil
			}

			// 3. 执行被测方法
			createdLink, err := service.CreateLink(context.Background(), tc.originalURL)

			// 4. 断言错误
			if tc.wantErrMsg != "" {
				if err == nil {
					t.Fatalf("Expected an error but got nil")
				}
				if err.Error() != tc.wantErrMsg {
					t.Fatalf("CreateLink() error = %q, wantErrMsg %q", err.Error(), tc.wantErrMsg)
				}
				return // 错误符合预期，测试结束
			}

			if err != nil {
				t.Fatalf("CreateLink() returned an unexpected error: %v", err)
			}

			// 5. 断言成功时的状态
			if createdLink == nil {
				t.Fatal("Expected a link to be created, but got nil")
			}
			assertLinkSaved(t, fakeRepo, createdLink.ShortCode, tc.originalURL)
		})
	}
}

// assertLinkSaved 是一个测试辅助函数
func assertLinkSaved(t *testing.T, repo *fakes.FakeLinkRepository, code, originalURL string) {
	t.Helper() // 标记为辅助函数，错误报告会指向调用处
	found, err := repo.FindByCode(context.Background(), code)
	if err != nil {
		t.Fatalf("FindByCode returned an unexpected error: %v", err)
	}
	if found == nil {
		t.Fatalf("Link with code %q was not found in the fake repository", code)
	}
	if found.OriginalURL != originalURL {
		t.Errorf("Link with code %q has incorrect OriginalURL: got %q, want %q", code, found.OriginalURL, originalURL)
	}
}

// a mock reader that always returns an error
type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated reader error")
}

// TestGenerateShortCode_Coverage 用于专门提升 generateShortCode 的测试覆盖率
func TestGenerateShortCode_Coverage(t *testing.T) {
	// 这是一个白盒测试，因为它需要篡改包内的未导出函数所依赖的全局变量

	t.Run("当 rand.Reader 返回错误时", func(t *testing.T) {
		// --- 狸猫换太子 ---

		// 1. 保存原始的 rand.Reader
		originalReader := rand.Reader

		// 2. 将 rand.Reader 替换为我们自己的 errReader
		rand.Reader = &errReader{}

		// 3. 使用 defer 确保在测试结束后，无论发生什么，都将原始的 Reader 恢复回去
		//    这是至关重要的，避免污染其他测试！
		defer func() {
			rand.Reader = originalReader
		}()

		// --- 执行测试 ---
		_, err := generateShortCode(6)

		// --- 断言 ---
		if err == nil {
			t.Fatal("期望一个错误，但得到了 nil")
		}
		if err.Error() != "simulated reader error" {
			t.Errorf("期望的错误信息不匹配: got %v", err)
		}
	})

	// 我们可以再加一个成功路径的测试，以确保我们的篡改和恢复逻辑是正确的
	t.Run("成功生成（验证恢复）", func(t *testing.T) {
		// 在这个子测试中，rand.Reader 应该已经被恢复为原始版本
		code, err := generateShortCode(6)
		if err != nil {
			t.Fatalf("不期望的错误: %v", err)
		}
		if len(code) != 8 { // 6 bytes -> 8 base64 chars
			t.Errorf("生成的 code 长度不符合预期: got %d", len(code))
		}
	})
}
