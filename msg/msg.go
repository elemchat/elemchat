package msg

import (
	"time"

	"github.com/elemchat/elemchat/magic"
)

var (
	_ Message = &Init{}
	_ Message = &WaitChat{}
	_ Message = &WaitMagic{}
	_ Message = &Chat{}
	_ Message = &Magic{}
	_ Message = &Effect{}
	_ Message = &Dualover{}
)

type Message interface {
	Type() Type
}

type Init struct {
}

type WaitChat struct {
	Deadline time.Time `json:"deadline"`
}

type WaitMagic struct {
	Deadline time.Time `json:"deadline"`
}

type Chat struct {
	Text string `json:"text"`
}

type Magic struct {
	Magic magic.Magic `json:"magic"`
}

type Effect struct {
}

type Dualover struct {
}

func (_ *Init) Type() Type {
	return INIT
}
func (_ *WaitChat) Type() Type {
	return WAIT_CHAT
}
func (_ *WaitMagic) Type() Type {
	return WAIT_MAGIC
}
func (_ *Chat) Type() Type {
	return CHAT
}
func (_ *Magic) Type() Type {
	return MAGIC
}
func (_ *Effect) Type() Type {
	return EFFECT
}
func (_ *Dualover) Type() Type {
	return DUALOVER
}
