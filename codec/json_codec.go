package codec

import (
	"encoding/json"
	"fmt"

	"github.com/elemchat/elemchat/msg"
)

var _ Codec = (*jsonCodec)(nil)

type jsonCodec struct{}

func JsonCodec() Codec {
	return (*jsonCodec)(nil)
}

func (c *jsonCodec) Encode(m msg.Message) ([]byte, error) {
	if m == nil {
		return nil, ErrMessageNil
	}

	j := make(map[string]interface{})
	j["type"] = msg.GetType(m).String()
	j["msg"] = m
	return json.Marshal(j)
}

func (c *jsonCodec) Decode(data []byte) (msg.Message, error) {
	j := make(map[string]interface{})
	err := json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}

	t, ok := j["type"]
	if !ok {
		return nil, ErrNoType
	}

	m, ok := j["msg"]
	if !ok {
		return nil, ErrNoMsg
	}
	mdata, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}

	return msg.Decode(msg.ToType(fmt.Sprint(t)), func(m msg.Message) error {
		return json.Unmarshal(mdata, m)
	})
}

func unmarshalMsg(data []byte, m msg.Message) (msg.Message, error) {
	if m == nil {
		return nil, ErrMessageNil
	}

	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
