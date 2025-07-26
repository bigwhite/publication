// internal/service/shortener_service_test.go
package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	// 导入真实的Store接口定义和错误
	"github.com/your_org/shortlink/internal/store"

	// 使用testify进行Mock和断言
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// MockStore 是 store.Store 接口的一个Mock实现。
// 它内嵌了mock.Mock，用于记录方法调用和返回预设值。
type MockStore struct {
	mock.Mock
}

// Save 为Store接口的Save方法实现Mock。
func (m *MockStore) Save(ctx context.Context, entry *store.LinkEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

// FindByShortCode 为Store接口的FindByShortCode方法实现Mock。
func (m *MockStore) FindByShortCode(ctx context.Context, shortCode string) (*store.LinkEntry, error) {
	args := m.Called(ctx, shortCode)
	// Get(0)获取第一个返回值，并尝试类型断言。
	entry, _ := args.Get(0).(*store.LinkEntry)
	return entry, args.Error(1)
}

// 为其他Store接口方法提供满足接口的Mock实现。
func (m *MockStore) FindByOriginalURLAndUserID(ctx context.Context, originalURL, userID string) (*store.LinkEntry, error) {
	// 在这个测试中不使用，返回nil即可。
	args := m.Called(ctx, originalURL, userID)
	entry, _ := args.Get(0).(*store.LinkEntry)
	return entry, args.Error(1)
}
func (m *MockStore) IncrementVisitCount(ctx context.Context, shortCode string) error {
	return m.Called(ctx, shortCode).Error(0)
}
func (m *MockStore) GetVisitCount(ctx context.Context, shortCode string) (int64, error) {
	args := m.Called(ctx, shortCode)
	val, _ := args.Get(0).(int64)
	return val, args.Error(1)
}
func (m *MockStore) Close() error { return m.Called().Error(0) }

// MockIDGenerator 是我们为IDGenerator接口创建的Mock实现。
type MockIDGenerator struct {
	mock.Mock
}

// Generate 为IDGenerator接口的Generate方法实现Mock。
func (m *MockIDGenerator) Generate(ctx context.Context, length int) (string, error) {
	args := m.Called(ctx, length)
	return args.String(0), args.Error(1)
}

// getTestLogger 返回一个用于测试的、不输出任何内容的slog.Logger。
func getTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestShortenerServiceImpl_CreateShortLink(t *testing.T) {
	logger := getTestLogger()
	// 为测试设置一个NoOp TracerProvider，这样service中的otel.Tracer()调用不会panic。
	otel.SetTracerProvider(trace.NewNoopTracerProvider())

	// 定义表驱动测试的用例结构体。
	testCases := []struct {
		name              string
		longURL           string
		userID            string
		mockStoreSetup    func(mockStore *MockStore)       // 用于设置Mock Store的行为。
		mockIdGenSetup    func(mockIdGen *MockIDGenerator) // 用于设置Mock ID生成器的行为。
		expectedShortCode string
		expectError       bool
		expectedErrorType error
	}{
		{
			name:    "SuccessfulCreation_FirstAttempt",
			longURL: "https://example.com/a-very-long-url",
			userID:  "user-123",
			mockStoreSetup: func(mockStore *MockStore) {
				// 期望FindByShortCode被调用一次，参数是"abcdefg"，并返回ErrNotFound。
				mockStore.On("FindByShortCode", mock.Anything, "abcdefg").Return(nil, store.ErrNotFound).Once()
				// 期望Save被调用一次，参数是一个*store.LinkEntry，并返回nil (无错误)。
				// 使用mock.MatchedBy进行更灵活的参数匹配，检查entry的关键字段。
				mockStore.On("Save", mock.Anything, mock.MatchedBy(func(e *store.LinkEntry) bool {
					return e.ShortCode == "abcdefg" && e.LongURL == "https://example.com/a-very-long-url"
				})).Return(nil).Once()
			},
			mockIdGenSetup: func(mockIdGen *MockIDGenerator) {
				// 期望ID生成器的Generate方法被调用一次，并返回"abcdefg"。
				mockIdGen.On("Generate", mock.Anything, defaultCodeLength).Return("abcdefg", nil).Once()
			},
			expectedShortCode: "abcdefg",
			expectError:       false,
		},
		{
			name:    "Collision_RetryOnce_ThenSuccess",
			longURL: "https://another-example.com",
			userID:  "user-456",
			mockStoreSetup: func(mockStore *MockStore) {
				// 第一次FindByShortCode，模拟冲突。
				mockStore.On("FindByShortCode", mock.Anything, "collide").Return(&store.LinkEntry{}, nil).Once()
				// 第二次FindByShortCode，模拟成功。
				mockStore.On("FindByShortCode", mock.Anything, "unique1").Return(nil, store.ErrNotFound).Once()
				// 随后的Save应该成功。
				mockStore.On("Save", mock.Anything, mock.MatchedBy(func(e *store.LinkEntry) bool { return e.ShortCode == "unique1" })).Return(nil).Once()
			},
			mockIdGenSetup: func(mockIdGen *MockIDGenerator) {
				// 期望Generate被调用两次，并按顺序返回不同的值。
				mockIdGen.On("Generate", mock.Anything, defaultCodeLength).Return("collide", nil).Once()
				mockIdGen.On("Generate", mock.Anything, defaultCodeLength).Return("unique1", nil).Once()
			},
			expectedShortCode: "unique1",
			expectError:       false,
		},
		{
			name:    "AllAttemptsCollide_Fails",
			longURL: "https://collision.com",
			userID:  "user-789",
			mockStoreSetup: func(mockStore *MockStore) {
				// 期望FindByShortCode被调用maxGenerationAttempts次，每次都返回已存在。
				mockStore.On("FindByShortCode", mock.Anything, mock.AnythingOfType("string")).Return(&store.LinkEntry{}, nil).Times(maxGenerationAttempts)
			},
			mockIdGenSetup: func(mockIdGen *MockIDGenerator) {
				// 期望Generate也被调用maxGenerationAttempts次。
				mockIdGen.On("Generate", mock.Anything, defaultCodeLength).Return("any-colliding-code", nil).Times(maxGenerationAttempts)
			},
			expectError:       true,
			expectedErrorType: ErrIDGenerationFailed,
		},
		// 可以添加更多测试用例，如Store.Save失败、输入校验失败等。
	}

	for _, tc := range testCases {
		currentTC := tc // 捕获range变量。
		t.Run(currentTC.name, func(t *testing.T) {
			// t.Parallel() // 如果测试用例之间完全独立，可以并行。

			// 1. 创建Mock实例
			mockStore := new(MockStore)
			mockIdGen := new(MockIDGenerator)

			// 2. 设置Mock期望
			currentTC.mockStoreSetup(mockStore)
			currentTC.mockIdGenSetup(mockIdGen)

			// 3. 创建被测Service实例，并注入Mock依赖
			serviceImpl := NewShortenerService(mockStore, logger, mockIdGen)

			// 4. 执行被测方法
			// 注意：CreateShortLink的完整签名可能需要更多参数，这里简化调用。
			shortCode, err := serviceImpl.CreateShortLink(context.Background(), currentTC.longURL, currentTC.userID, "", time.Time{})

			// 5. 断言结果
			if currentTC.expectError {
				assert.Error(t, err, "Expected an error for test case: %s", currentTC.name)
				if currentTC.expectedErrorType != nil {
					assert.ErrorIs(t, err, currentTC.expectedErrorType, "Error type mismatch for test case: %s", currentTC.name)
				}
			} else {
				assert.NoError(t, err, "Did not expect an error for test case: %s", currentTC.name)
			}

			if !currentTC.expectError && currentTC.expectedShortCode != "" {
				assert.Equal(t, currentTC.expectedShortCode, shortCode, "Short code mismatch for test case: %s", currentTC.name)
			}

			// 6. 验证所有Mock期望都已满足
			mockStore.AssertExpectations(t)
			mockIdGen.AssertExpectations(t)
		})
	}
}
