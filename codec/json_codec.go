package codec

import (
	"encoding/json"

	"github.com/elemchat/elemchat/msg"
)

var _ Codec = (*jsonCodec)(nil)

type jsonCodec struct{}

func (c *jsonCodec) Encode(m msg.Message) ([]byte, error) {
	if m == nil {
		return nil, ErrMessageNil
	}

	j := make(map[string]interface{})
	j["type"] = m.Type().String()
	j["msg"] = m
	return json.Marshal(j)
}

func (c *jsonCodec) Decode(data []byte, m msg.Message) error {
	return nil
}
