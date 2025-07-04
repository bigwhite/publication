package simplehash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/your_org/shortlink/internal/idgen"
)

// 确保 *Generator 实现了 idgen.Generator 接口 (编译时检查)
var _ idgen.Generator = (*Generator)(nil)

const defaultCodeLength = 7

// Generator 是基于简单哈希的短码生成器具体实现类型
type Generator struct {
}

// New 创建一个新的简单哈希生成器实例
func New() *Generator {
	rand.Seed(time.Now().UnixNano())
	return &Generator{}
}

// GenerateShortCode 为给定的长URL生成一个短码
func (g *Generator) GenerateShortCode(ctx context.Context, longURL string) (string, error) {
	if longURL == "" {
		return "", idgen.ErrInputIsEmpty
	}
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	hasher.Write([]byte(fmt.Sprintf("%d", rand.Int63())))
	hashBytes := hasher.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(hashBytes)

	if len(encoded) < defaultCodeLength {
		return encoded, nil
	}
	return encoded[:defaultCodeLength], nil
}
