package protocol

import (
	"bufio"
	"bytes"
	"io"
)

// Decoder RLE 解码器
type Decoder struct {
	r *bufio.Reader // 使用 bufio 以便 Peek 和 ReadByte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decode() ([]byte, error) {
	var buf bytes.Buffer

	for {
		b, err := d.r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if b == EscapeByte {
			// 读到了转义符，读取 Count
			countByte, err := d.r.ReadByte()
			if err != nil {
				return nil, err // 数据截断，格式错误
			}
			count := int(countByte)

			// 读取 Value
			val, err := d.r.ReadByte()
			if err != nil {
				return nil, err
			}

			// 写入 Count 个 Value
			for k := 0; k < count; k++ {
				buf.WriteByte(val)
			}
		} else {
			// 普通字符，直接写入
			buf.WriteByte(b)
		}
	}
	return buf.Bytes(), nil
}
