// internal/service/shortener_service.go
package service

import (
	"context"
	"crypto/rand" // 用于生成更安全的随机串
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/your_org/shortlink/internal/store"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultCodeLength     = 7       // 短码默认长度
	maxGenerationAttempts = 5       // 生成唯一短码的最大尝试次数
	defaultExpiryHours    = 24 * 30 // 默认30天过期 (示例)
)

// ErrIDGenerationFailed 表示在多次尝试后仍无法生成唯一的短码
var ErrIDGenerationFailed = errors.New("service: failed to generate unique short ID after multiple attempts")

// IDGenerator defines the interface for generating short codes.
type IDGenerator interface {
	Generate(ctx context.Context, length int) (string, error)
}

type secureRandomIDGenerator struct{} // 默认实现
func (g *secureRandomIDGenerator) Generate(ctx context.Context, length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("service: failed to generate random bytes: %w", err)
	}
	// 使用URLEncoding可以避免/和+字符，使其更适合用在URL路径中
	// 并取其前length个字符
	shortCode := base64.URLEncoding.EncodeToString(randomBytes)
	shortCode = strings.ReplaceAll(shortCode, "_", "A") // 替换下划线
	shortCode = strings.ReplaceAll(shortCode, "-", "B") // 替换连字符
	if len(shortCode) < length {                        // 理论上不太可能，除非length非常小
		err := fmt.Errorf("generated base64 string too short")
		return "", err
	}
	finalCode := shortCode[:length]
	return finalCode, nil
}

func NewSecureRandomIDGenerator() IDGenerator { return &secureRandomIDGenerator{} }

// ShortenerService 定义了短链接服务的核心业务逻辑接口
type ShortenerService interface {
	// CreateShortLink 为给定的长URL创建一个短链接。
	// userID 和 originalURL 用于防止同一用户为同一原始URL重复创建。
	// expireAt 用于设置链接的过期时间，零值表示永不过期或使用默认过期。
	CreateShortLink(ctx context.Context, longURL string, userID string, originalURL string, expireAt time.Time) (shortCode string, err error)

	// GetOriginalURL 根据短码查找并返回原始长链接及其相关信息。
	// 如果短码不存在或已过期，应返回 store.ErrNotFound。
	// 此方法也可能负责增加访问计数。
	GetOriginalURL(ctx context.Context, shortCode string) (linkEntry *store.LinkEntry, err error)
}

type shortenerServiceImpl struct {
	store       store.Store  // 存储接口的依赖
	logger      *slog.Logger // 日志记录器
	tracer      trace.Tracer // OpenTelemetry Tracer，用于手动创建子Span
	idGenerator IDGenerator
}

// NewShortenerService 创建一个新的ShortenerService实例
func NewShortenerService(s store.Store, logger *slog.Logger, idGen IDGenerator) ShortenerService {
	if idGen == nil {
		idGen = NewSecureRandomIDGenerator()
	}

	return &shortenerServiceImpl{
		store:  s,
		logger: logger,
		// 获取一个Tracer实例。Tracer的名称通常使用其所属的库或模块的导入路径。
		tracer:      otel.Tracer("github.com/your_org/shortlink/internal/service"),
		idGenerator: idGen,
	}
}

// generateSecureRandomCode 生成一个指定长度的、URL安全的随机字符串作为短码
func (s *shortenerServiceImpl) generateSecureRandomCode(ctx context.Context, length int) (string, error) {
	// 为这个内部操作创建一个子Span
	_, span := s.tracer.Start(ctx, "service.generateSecureRandomCode")
	defer span.End()

	shortCode, err := s.idGenerator.Generate(ctx, length)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "generate random string failed")
		return "", fmt.Errorf("service: failed to generate random bytes: %w", err)
	}

	span.SetAttributes(attribute.String("generated_code", shortCode))
	return shortCode, nil
}

// CreateShortLink 为给定的长URL创建一个短链接
func (s *shortenerServiceImpl) CreateShortLink(ctx context.Context, longURL string, userID string, originalURL string, expireAt time.Time) (string, error) {
	// 1. 从传入的ctx启动一个新的子Span，用于追踪这个方法的执行
	ctx, span := s.tracer.Start(ctx, "ShortenerService.CreateShortLink", trace.WithAttributes(
		attribute.String("long_url", longURL),
		attribute.String("user_id", userID),
		attribute.String("original_url", originalURL),
	))
	defer span.End() // 确保Span在函数退出时被结束

	s.logger.DebugContext(ctx, "Service: Attempting to create short link.",
		slog.String("longURL", longURL), slog.String("userID", userID))

	// 校验输入
	if strings.TrimSpace(longURL) == "" {
		err := errors.New("long URL cannot be empty")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.WarnContext(ctx, "Validation failed for CreateShortLink: long URL empty")
		return "", err
	}
	// (可以添加更多对longURL格式的校验)

	// 如果originalURL为空，则使用longURL
	if originalURL == "" {
		originalURL = longURL
	}

	// 如果expireAt是零值，设置一个默认的过期时间
	if expireAt.IsZero() {
		expireAt = time.Now().Add(time.Hour * defaultExpiryHours)
		span.SetAttributes(attribute.String("default_expiry_applied", expireAt.Format(time.RFC3339)))
	}
	span.SetAttributes(attribute.String("expire_at", expireAt.Format(time.RFC3339)))

	var shortCode string
	var err error

	for attempt := 0; attempt < maxGenerationAttempts; attempt++ {
		// 为生成和检查短码的每次尝试创建一个更细粒度的子Span
		attemptCtx, attemptSpan := s.tracer.Start(ctx, fmt.Sprintf("ShortenerService.CreateShortLink.Attempt%d", attempt+1))

		// shortCode, err = s.generateSecureRandomCode(attemptCtx, defaultCodeLength) // 使用内部方法
		shortCode, err = s.idGenerator.Generate(attemptCtx, defaultCodeLength)
		if err != nil {
			// 这个错误通常是内部生成随机串失败，比较严重，直接返回
			s.logger.ErrorContext(attemptCtx, "Service: Failed to generate random string for short code.", slog.Any("error", err))
			attemptSpan.RecordError(err)
			attemptSpan.SetStatus(codes.Error, "short code internal generation failed")
			attemptSpan.End()

			span.RecordError(err) // 在父Span也记录这个严重错误
			span.SetStatus(codes.Error, "short code internal generation failed")
			return "", fmt.Errorf("service: internal error generating short code: %w", err)
		}
		attemptSpan.SetAttributes(attribute.String("generated_code_attempt", shortCode))
		s.logger.DebugContext(attemptCtx, "Service: Generated short code attempt.", slog.String("short_code", shortCode))

		// 检查短码是否已存在 (这个store调用也应该被追踪)
		// 我们将在Store的实现中添加Tracing (如果需要更细粒度)
		// 这里，我们假设FindByShortCode是Store接口的一部分
		existingEntry, findErr := s.store.FindByShortCode(attemptCtx, shortCode) // 将attemptCtx传递下去

		if errors.Is(findErr, store.ErrNotFound) {
			s.logger.InfoContext(attemptCtx, "Service: Generated short code is unique.", slog.String("short_code", shortCode))

			entryToSave := &store.LinkEntry{
				ShortCode:   shortCode,
				LongURL:     longURL,
				OriginalURL: originalURL,
				UserID:      userID,
				CreatedAt:   time.Now(), // 在Service层设置创建时间
				ExpireAt:    expireAt,
			}
			saveErr := s.store.Save(attemptCtx, entryToSave)

			if saveErr == nil {
				s.logger.InfoContext(ctx, "Service: Successfully created and saved short link.",
					slog.String("short_code", shortCode),
					slog.String("long_url", longURL),
				)
				span.SetAttributes(attribute.String("final_short_code", shortCode))
				span.SetStatus(codes.Ok, "short link created")
				attemptSpan.SetStatus(codes.Ok, "short code unique and saved")
				attemptSpan.End()
				return shortCode, nil
			}
			s.logger.WarnContext(attemptCtx, "Service: Failed to save unique short code, retrying if possible.",
				slog.String("short_code", shortCode), slog.Any("save_error", saveErr))
			err = saveErr
			attemptSpan.RecordError(saveErr)
			attemptSpan.SetStatus(codes.Error, "failed to save short code")
		} else if findErr != nil {
			s.logger.ErrorContext(attemptCtx, "Service: Store error checking short code existence.",
				slog.String("short_code", shortCode), slog.Any("find_error", findErr))
			err = findErr
			attemptSpan.RecordError(findErr)
			attemptSpan.SetStatus(codes.Error, "store error checking short code")
		} else {
			// findErr is nil and existingEntry is not nil, means shortCode already exists
			s.logger.DebugContext(attemptCtx, "Service: Short code collision.", slog.String("short_code", shortCode), slog.String("existing_long_url", existingEntry.LongURL))
			err = fmt.Errorf("short code %s collision (attempt %d)", shortCode, attempt+1)
			attemptSpan.SetAttributes(attribute.Bool("collision", true))
			attemptSpan.SetStatus(codes.Error, "short code collision")
		}
		attemptSpan.End()

		// 如果是存储错误（而不是未找到或冲突），则不应继续重试
		if !errors.Is(findErr, store.ErrNotFound) && findErr != nil {
			span.RecordError(err) // 将store错误记录到父span
			span.SetStatus(codes.Error, "failed due to store error during creation")
			return "", fmt.Errorf("service: failed to create short link after store error: %w", err)
		}
	}

	// 如果循环结束仍未成功 (通常是多次冲突)
	finalErr := ErrIDGenerationFailed
	s.logger.ErrorContext(ctx, "Service: Failed to generate a unique short code after multiple attempts.",
		slog.Int("max_attempts", maxGenerationAttempts),
	)
	span.RecordError(finalErr)
	span.SetStatus(codes.Error, finalErr.Error())
	return "", finalErr
}

// GetOriginalURL 根据短码查找原始长链接
func (s *shortenerServiceImpl) GetOriginalURL(ctx context.Context, shortCode string) (*store.LinkEntry, error) {
	ctx, span := s.tracer.Start(ctx, "ShortenerService.GetOriginalURL", trace.WithAttributes(
		attribute.String("short_code", shortCode),
	))
	defer span.End()

	s.logger.InfoContext(ctx, "Service: Attempting to retrieve original URL for short code.", slog.String("short_code", shortCode))

	if strings.TrimSpace(shortCode) == "" {
		err := errors.New("short code cannot be empty")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.WarnContext(ctx, "Validation failed for GetOriginalURL: short code empty")
		return nil, err
	}

	entry, err := s.store.FindByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.WarnContext(ctx, "Service: Failed to find original URL for short code.",
			slog.String("short_code", shortCode),
			slog.Any("error", err), // 这个err可能是store.ErrNotFound或者其他DB错误
		)
		span.RecordError(err)
		if errors.Is(err, store.ErrNotFound) {
			span.SetStatus(codes.Error, "short code not found") // 更具体的错误状态
		} else {
			span.SetStatus(codes.Error, "store error retrieving short code")
		}
		return nil, err
	}

	// 检查是否过期
	if !entry.ExpireAt.IsZero() && time.Now().After(entry.ExpireAt) {
		s.logger.InfoContext(ctx, "Service: Short code found but has expired.",
			slog.String("short_code", shortCode),
			slog.Time("expire_at", entry.ExpireAt),
		)
		// (可选) 在这里可以从存储中删除过期的条目 (异步或同步)
		// s.store.Delete(ctx, shortCode)
		span.SetAttributes(attribute.Bool("expired", true))
		span.SetStatus(codes.Error, "short code expired")
		return nil, store.ErrNotFound // 对外表现为未找到
	}

	s.logger.InfoContext(ctx, "Service: Successfully retrieved original URL.",
		slog.String("short_code", shortCode),
		slog.String("original_url", entry.OriginalURL), // 假设LinkEntry有OriginalURL
	)
	span.SetAttributes(attribute.String("retrieved_original_url", entry.OriginalURL))
	span.SetStatus(codes.Ok, "original URL retrieved")

	// 可以在这里异步增加访问计数 (如果Store.IncrementVisitCount是异步安全的，或者用另一个goroutine)
	// go s.store.IncrementVisitCount(context.Background(), shortCode) // 简单示例，实际要考虑错误处理
	// 或者，如果IncrementVisitCount是快速的，也可以同步调用
	// visitCtx, visitSpan := s.tracer.Start(ctx, "ShortenerService.IncrementVisitCount")
	// if errVisit := s.store.IncrementVisitCount(visitCtx, shortCode); errVisit != nil {
	//  s.logger.ErrorContext(visitCtx, "Failed to increment visit count", slog.String("short_code", shortCode), slog.Any("error", errVisit))
	//  visitSpan.RecordError(errVisit)
	//  visitSpan.SetStatus(codes.Error, "failed to increment visit count")
	// }
	// visitSpan.End()

	return entry, nil
}
