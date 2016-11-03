package msg

import (
	"time"

	"github.com/elemchat/elemchat/magic"
)

type Type int

const (
	ILLEGAL Type = iota
)

type Message interface {
	Type() Type
}

type WaitChat struct {
	deadline time.Time
}

type WaitMagic struct {
	deadline time.Time
}

type Chat struct {
	Text string
}

type Magic struct {
	Magic magic.Magic
}
