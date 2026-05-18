package context

import (
	"fmt"
	"log"

	"github.com/yourname/go-tiny-claw/internal/schema"
)

type Compactor struct {
	MaxChars       int
	RetainLastMsgs int
}

func NewCompactor(maxChars int, retainLastMsgs int) *Compactor {
	return &Compactor{
		MaxChars:       maxChars,
		RetainLastMsgs: retainLastMsgs,
	}
}

func (c *Compactor) Compact(msgs []schema.Message) []schema.Message {
	currentLength := c.estimateLength(msgs)

	if currentLength < c.MaxChars {
		return msgs
	}

	log.Printf("[Compactor] ⚠️ 内存告警：当前上下文长度 (%d 字符) 超过阈值 (%d)，触发压缩清理...\n", currentLength, c.MaxChars)

	var compacted []schema.Message
	msgCount := len(msgs)

	protectStartIndex := msgCount - c.RetainLastMsgs
	if protectStartIndex < 0 {
		protectStartIndex = 0
	}

	for i, msg := range msgs {
		if msg.Role == schema.RoleSystem {
			compacted = append(compacted, msg)
			continue
		}

		newMsg := msg
		isInWorkingMemory := i >= protectStartIndex

		if msg.Role == schema.RoleUser && msg.ToolCallID != "" {
			if !isInWorkingMemory {
				if len(msg.Content) > 200 {
					newMsg.Content = fmt.Sprintf("...[为了节省内存，早期的工具输出已被系统强制清理。原始长度: %d 字节]...", len(msg.Content))
				}
			} else {
				const maxKeep = 1000
				if len(msg.Content) > maxKeep {
					head := msg.Content[:500]
					tail := msg.Content[len(msg.Content)-500:]
					newMsg.Content = fmt.Sprintf("%s\n\n...[内容过长，中间 %d 字节已被系统截断]...\n\n%s", head, len(msg.Content)-maxKeep, tail)
				}
			}
		} else if msg.Role == schema.RoleAssistant && msg.Content != "" {
			if !isInWorkingMemory && len(msg.Content) > 200 {
				newMsg.Content = "...[早期的推理思考过程已折叠]..."
			}
		}

		compacted = append(compacted, newMsg)
	}

	newLength := c.estimateLength(compacted)
	log.Printf("[Compactor] ✅ 压缩完成。上下文长度从 %d 降至 %d 字符。\n", currentLength, newLength)

	return compacted
}

func (c *Compactor) estimateLength(msgs []schema.Message) int {
	length := 0
	for _, msg := range msgs {
		length += len(msg.Content)
		for _, tc := range msg.ToolCalls {
			length += len(tc.Name) + len(tc.Arguments)
		}
	}
	return length
}
