package observability

import (
	"context"
	"log"
	"time"

	ctxpkg "github.com/yourname/go-tiny-claw/internal/context"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/schema"
)

var PricingModel = map[string]struct {
	InputPrice  float64
	OutputPrice float64
}{
	"glm-4.5-air": {InputPrice: 0.15, OutputPrice: 0.15},
}

type CostTracker struct {
	nextProvider provider.LLMProvider
	modelName    string
	session      *ctxpkg.Session
}

func NewCostTracker(next provider.LLMProvider, modelName string, session *ctxpkg.Session) *CostTracker {
	return &CostTracker{
		nextProvider: next,
		modelName:    modelName,
		session:      session,
	}
}

func (t *CostTracker) Generate(ctx context.Context, msgs []schema.Message, availableTools []schema.ToolDefinition) (*schema.Message, error) {
	startTime := time.Now()

	respMsg, err := t.nextProvider.Generate(ctx, msgs, availableTools)

	latency := time.Since(startTime)

	if err != nil {
		log.Printf("[Tracker] ❌ API 调用失败，耗时: %v\n", latency)
		return respMsg, err
	}

	if respMsg.Usage != nil {
		promptTokens := respMsg.Usage.PromptTokens
		completionTokens := respMsg.Usage.CompletionTokens

		var cost float64
		if price, exists := PricingModel[t.modelName]; exists {
			cost = (float64(promptTokens)*price.InputPrice + float64(completionTokens)*price.OutputPrice) / 1000000.0
		}

		log.Printf("[Tracker] 📊 API 调用完成 | 耗时: %v | 输入: %d tk | 输出: %d tk | 花费: ¥%.6f\n",
			latency, promptTokens, completionTokens, cost)

		if t.session != nil {
			t.session.RecordUsage(promptTokens, completionTokens, cost)
			log.Printf("[Tracker] 💰 当前会话 (%s) 累计花费: ¥%.6f\n", t.session.ID, t.session.TotalCostCNY)
		}
	} else {
		log.Printf("[Tracker] ⚠️ API 调用完成，但未返回 Usage 数据 | 耗时: %v\n", latency)
	}

	return respMsg, nil
}
