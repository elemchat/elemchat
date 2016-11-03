package codec

import "github.com/elemchat/elemchat/msg"

type Codec interface {
	Encode(msg.Message) ([]byte, error)
	Decode([]byte, msg.Message) error
}
