package codec

import "github.com/elemchat/elemchat/msg"

var _ Codec = (*jsonCodec)(nil)

type jsonCodec struct{}

func (c *jsonCodec) Encode(m msg.Message) ([]byte, error) {
	return nil, nil
}

func (c *jsonCodec) Decode(data []byte, m msg.Message) error {
	return nil
}
