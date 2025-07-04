package idgen

import (
	"context"
	"errors"
)

// Generator 定义了短码生成器的能力。
// 实现应尽可能保证生成的短码在一定概率下是唯一的，并符合业务对长度、字符集的要求。
type Generator interface {
	// GenerateShortCode 为给定的输入（通常是长URL）生成一个短码。
	// - ctx: 用于传递超时或取消信号，例如ID生成依赖外部服务时。
	// - input: 用于生成短码的原始数据，通常是长URL。
	// - 返回生成的短码和可能的错误（如生成超时、内部错误等）。
	GenerateShortCode(ctx context.Context, input string) (string, error)
}

var ErrGeneratorUnavailable = errors.New("idgen: generator service unavailable")
var ErrInputTooLongForGenerator = errors.New("idgen: input too long to generate short code")
var ErrInputIsEmpty = errors.New("idgen: input is empty")
