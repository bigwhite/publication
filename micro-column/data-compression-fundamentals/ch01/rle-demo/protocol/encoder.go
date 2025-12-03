package protocol

import (
	"io"
)

const (
	EscapeByte = 0xFF // 定义转义符
	MinRunLen  = 4    // 只有连续出现4次以上才压缩，否则开销划不来
)

// Encoder RLE 编码器
type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(data []byte) error {
	length := len(data)
	if length == 0 {
		return nil
	}

	for i := 0; i < length; {
		// 1. 寻找连续重复的字符长度
		runLen := 1
		for i+runLen < length && data[i+runLen] == data[i] && runLen < 255 {
			runLen++
		}

		val := data[i]

		// 2. 判断是否触发压缩逻辑
		// 条件：重复次数 >= 阈值，或者 遇到转义字符本身
		if runLen >= MinRunLen || val == EscapeByte {
			// 写入：[Escape] [Count] [Value]
			// 特殊情况：如果 val 是 EscapeByte，runLen 可能是 1。
			// 协议规定：0xFF 后面的 0x00 代表 原文是 0xFF，0x04 代表 4个0xFF
			// 这里为了简化，统一用 [FF] [Count] [Value]

			// 优化：如果只是单个 EscapeByte，为了区分，我们可以规定
			// FF 00 -> 原文 FF
			// FF N X -> N个X (N >= 4)
			// 但为了演示最通用的逻辑，我们这里简单粗暴：
			// 只要是 EscapeByte 或者 runLen >= 4，都用三字节表示。

			header := []byte{EscapeByte, byte(runLen), val}
			if _, err := e.w.Write(header); err != nil {
				return err
			}
		} else {
			// 3. 不压缩，直接写入原文（runLen 次）
			// 注意：这里必须把这 runLen 个字符都写进去，因为没触发压缩
			if _, err := e.w.Write(data[i : i+runLen]); err != nil {
				return err
			}
		}

		i += runLen
	}
	return nil
}
