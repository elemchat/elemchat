package codec

import (
	"errors"

	"github.com/elemchat/elemchat/msg"
)

type Codec interface {
	Encode(msg.Message) ([]byte, error)
	Decode([]byte) (msg.Message, error)
}

var (
	ErrMessageNil  = errors.New("message is nil")
	ErrNoType      = errors.New("no type")
	ErrNoMsg       = errors.New("no msg")
	ErrIllegalType = errors.New("illegal type")
)
