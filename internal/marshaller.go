package openai

import (
	"bytes"
	"encoding/json"
)

type Marshaller interface {
	Marshal(value any) ([]byte, error)
}

type JSONMarshaller struct{}

func (jm *JSONMarshaller) Marshal(value any) ([]byte, error) {
	// 禁用 HTML 转义，避免将 < > & 等字符转义为 \u003c \u003e \u0026。
	// 否则像 stop=["</output>"] 这类包含特殊字符的参数会被转义，
	// 在部分下游处理链路中无法正确匹配，导致 stop 不生效。
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	// Encoder.Encode 会追加一个换行符，去掉它以保持与 json.Marshal 一致。
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
